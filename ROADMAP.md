# cloud-terminal 开发路线图

> 项目定位：自托管、多租户 SSH 控制平面 (SSH Access Control Plane) + IAM 权限平台
> 核心定位：通过抽象的“资源 (Resource)”而非原生服务器 IP 访问基础设施，结合 ContainerSSH 临时隔离沙箱完成 JIT 会话

---

## 总体架构

```
                    ┌─────────────────────────┐
                    │   RESTful CRUD API      │  ← handler/
                    ├─────────────────────────┤
                    │   WebSocket 终端        │  ← handler/terminal.go
                    ├─────────────────────────┤
                    │   ContainerSSH Webhook  │  ← handler/webhook_v2.go
                    ├─────────────────────────┤
                    │   STS 令牌服务          │  ← pkg/sts/
                    ├─────────────────────────┤
                    │ ★ IAM 策略引擎          │  ← pkg/iam/
                    ├─────────────────────────┤
                    │   pkg/connector/        │  ← 本次新增
                    └─────────────────────────┘
```

---

## 优先级总览

| 优先级 | 阶段 | 任务 | 工作量 |
|---|---|---|---|
| 🔴 P0 | Foundation | `pkg/connector/` 接口抽象 | ~200 行 |
| 🔴 P0 | Security | `auth_data` 列级加密 | ~150 行 |
| 🔴 P0 | Security | 密钥环境变量化 | ~50 行 |
| 🟡 P1 | Security | 生产级 HostKey 校验 | ~50 行 |
| 🟡 P1 | Connector | ContainerSSH proxy backend | ~80 行 |
| 🟡 P1 | Feature | 会话录像 + 回放 | ~400 行 |
| 🟢 P2 | Feature | 多协议 RDP/VNC | ~100 行 |
| 🟢 P2 | Deploy | Helm Chart + NetworkPolicy | ~300 行 |
| 🟢 P2 | Infra | Session store Redis/DB | ~100 行 |
| 🔵 P3 | Security | Webhook mTLS | ~200 行 |
| 🔵 P3 | Feature | LDAP/OIDC | ~500 行 |
| 🔵 P3 | Feature | Session operator（K8s） | ~1000 行 |

---

## Iteration 1 — Foundation（P0）

### 1.1 `pkg/connector/` 接口抽象

**目标**：将 `terminal.go` 和 `webhook_v2.go` 中硬编码的连接逻辑抽成统一接口，后续所有连接能力建立在此之上。

```go
// pkg/connector/connector.go
type Connector interface {
    Name() string
    Connect(ctx, *ConnectRequest) (Connection, error)
    ContainerSSHConfig(ctx, *ConnectRequest) (*ContainerSSHConfig, error)
}

type ConnectRequest struct {
    SessionID       string
    ResourceURN     string
    Protocol        string            // "ssh", "rdp", "mysql"...
    Target          TargetInfo        // IP + Port
    AuthData        map[string]interface{}
    ResourceDetails map[string]interface{}
}

type ContainerSSHConfig struct {
    Backend    string             // "docker" | "proxy" | "kubernetes"
    Proxy      *ProxyConfig       `json:",omitempty"`
    Docker     *DockerConfig      `json:",omitempty"`
    Kubernetes *KubernetesConfig  `json:",omitempty"`
}
```

**实现**：
1. 定义接口和类型 → `pkg/connector/connector.go`
2. 把 `terminal.go` 的 `dialSSH` 迁入 `direct` connector
3. 把 `webhook_v2.go` 的 `buildContainerConfig` 迁入 `container` connector
4. `terminal.go` 改为调用 `connector.Connect()`
5. 零行为变化，所有测试通过

**文件清单**：
- 新增 `pkg/connector/connector.go`（接口定义）
- 新增 `pkg/connector/direct.go`（direct 实现）
- 新增 `pkg/connector/container.go`（container 实现）
- 修改 `handler/terminal.go`（改用接口）
- 修改 `handler/webhook_v2.go`（改用接口）

---

### 1.2 `auth_data` 列级加密

**目标**：Resource 表中的 `auth_data` 字段（含 SSH 私钥、密码）当前为明文 JSON 存储。改为 AES-256-GCM 列级加密。

**方案**：
1. 新增 `pkg/crypto/envelope.go` — 信封加密
   - 主密钥（KEK）从环境变量读取（`RESOURCE_ENCRYPTION_KEY`）
   - 每个资源生成独立的数据密钥（DEK）
   - `auth_data` 存储 `{encrypted_data, encrypted_dek, nonce}`
2. 在 Resource Create/Update handler 中加密
3. 在 Resource Get/List handler 中解密
4. `dialSSH` / container connector 读取时自动解密

**数据流**：
```
写入：
  auth_data(明文) → AES-256-GCM(DEK) → {ciphertext, encrypted_DEK, nonce} → DB

读取：
  DB → {ciphertext, encrypted_DEK, nonce} → AES-256-GCM(DEK) → auth_data(明文) → 使用
```

**文件清单**：
- 新增 `pkg/crypto/envelope.go`
- 修改 `handler/resource.go`（Create/Update/Get 时加解密）
- 修改 `pkg/connector/direct.go`（读取时解密）

---

### 1.3 密钥环境变量化

**目标**：消除所有硬编码密钥和 config.yaml 明文密钥。

**改动**：
| 当前 | 改为 |
|---|---|
| `cfg.Server.JWTSecret` 从 yaml 读 | `JWT_SECRET` 环境变量 |
| `handler/session.go` 硬编码 | `SESSION_SECRET` 环境变量 |
| `RESOURCE_ENCRYPTION_KEY` | 新增，用于 auth_data 加密 |

优先级：环境变量 > yaml > 默认值

**文件清单**：
- 修改 `cmd/cloud-terminal/main.go`
- 修改 `handler/session.go`
- 修改 `pkg/config/config.go`

---

## Iteration 2 — Security & ContainerSSH（P1）

### 2.1 生产级 HostKey 校验

**目标**：`dialSSH` 中的 `ssh.InsecureIgnoreHostKey()` 替换为生产级校验。

**方案**：
- 首次连接时记录目标主机密钥到 DB（`known_hosts` 表或 Resource 字段）
- 后续连接时校验是否匹配
- 支持 `HostKeyCallback: ssh.FixedHostKey()` 用于已知密钥

**文件**：`pkg/connector/direct.go`（HostKeyCallback）

---

### 2.2 ContainerSSH proxy backend

**目标**：ContainerSSH 配置支持直连代理模式（不启动容器）。

**现有流程**：
```
ContainerSSH → ConfigWebhook → buildContainerConfig → Docker 配置 → 启容器
```

**新增流程**：
```
ContainerSSH → ConfigWebhook → ContainerSSHConfig{Backend: "proxy"} → TCP 直连
```

由部署配置控制（环境变量 `SESSION_ISOLATION=proxy|container`）。

**文件**：`pkg/connector/container.go`

---

### 2.3 会话录像 + 回放

**目标**：记录所有 SSH 会话为 asciinema 格式，支持回放。

**方案**：
1. **录像写入**：在 connector 层拦截 I/O 流，写入 asciinema cast 格式
   - 每个会话一个 cast 文件
   - 存储到 S3/MinIO（租户可配置后端）
   - 记录路径到 `session.recording_path`
2. **录像回放**：前端嵌入 asciinema player
   - `GET /api/sessions/:id/recording` 返回 cast 文件
   - 前端用 asciinema-player 渲染

**ContainerSSH 集成**：ContainerSSH 内核已支持 asciinema 格式录像（`AuditLogFormatAsciinema`），webhook 配置中开启即可。

**文件清单**：
- 新增 `pkg/recording/writer.go`（录像写入）
- 新增 `pkg/recording/player.go`（回放 API）
- 修改 `pkg/connector/direct.go`（录像拦截层）
- 修改 `handler/session.go`（添加回放路由）
- 前端：asciinema-player 组件

---

## Iteration 3 — 部署 & 多协议（P2）

### 3.1 Helm Chart + NetworkPolicy

**目标**：一键部署到 K8s。

**结构**：
```
deploy/helm/cloud-terminal/
├── templates/
│   ├── deployment.yaml       # API server
│   ├── statefulset.yaml      # PostgreSQL
│   ├── configmap.yaml        # 配置
│   ├── secret.yaml           # 密钥（需提前创建）
│   ├── ingress.yaml          # TLS + WebSocket
│   ├── networkpolicy.yaml    # 组件间隔离
│   └── service.yaml
└── values.yaml
```

**NetworkPolicy 策略**：
- API ↔ DB：只允许 5432
- API ↔ Internet：只允许 443
- ContainerSSH ↔ API：只允许 `/webhook/*`
- ContainerSSH ↔ 目标：允许所有出口（必需）

---

### 3.2 多协议 RDP/VNC

**目标**：补充 `buildContainerConfig`/`container connector` 中缺失的 RDP/VNC/Telnet/HTTP 协议处理。

**文件**：`pkg/connector/container.go` 的协议 switch

---

### 3.3 Session store Redis/DB

**目标**：`CookieStore` 替换为可水平扩展的存储。

**方案**：
- 新增 `pkg/session/` 接口
- `redis` 实现（生产）
- `database` 实现（DB 回退）
- `cookie` 实现（当前，单机保留）

---

## Iteration 4 — 企业特性（P3）

### 4.1 Webhook mTLS

**目标**：ContainerSSH ↔ cloud-terminal 的 webhook 调用从 HTTP 升级为 mTLS。

**方案**：
- API server 启动时加载 CA 证书
- ContainerSSH 配置中指定客户端证书
- 证书由 cert-manager 自动管理（K8s 部署）
- 支持双向证书验证

---

### 4.2 LDAP/OIDC 外部认证

**目标**：租户可绑定外部身份源。

**方案**：
- 新增 `pkg/idp/` 接口
- `ldap` 实现
- `oidc` 实现（支持 Keycloak/Azure AD/Okta）
- 登录时：校验外部密码 → 同步/映射本地用户 → 签发 JWT

---

### 4.3 Session operator（K8s 隔离）

**目标**：每个 SSH 会话一个独立 K8s Pod，实现 Namespace 级隔离。

**方案**：
- Session 创建时在租户 Namespace 创建临时 Pod
- Pod 使用租户的 ServiceAccount + NetworkPolicy
- 会话结束后自动清理 Pod
- 支持独立部署或内嵌 controller 模式

---

## 迭代依赖关系图

```
Iteration 1                     Iteration 2                 Iteration 3
─────────────────────           ─────────────────────       ─────────────────────
pkg/connector/  ←──────┬──────→ HostKey 校验                 多协议 RDP/VNC
                        │        proxy backend                Helm Chart
auth_data 加密  ←───────┤        Session 录像                 Session store
                        │
密钥环境变量    ←───────┘
                                       │
                                       │
                                       ▼
                                  Iteration 4
                                  ─────────────────────
                                  mTLS
                                  LDAP/OIDC
                                  Session operator
```

---

## 当前状态

- ✅ `pkg/iam/` — 完整 IAM 引擎（AWS 风格，5 层评估链）
- ✅ `pkg/sts/` — STS 令牌服务
- ✅ RESTful CRUD API — 所有实体完整
- ✅ WebSocket 终端 — 带 Session/AuditLog 记录
- ✅ ContainerSSH v2 Webhook — 认证 + 配置
- ✅ `api/init.go` — 初始化流程修复
- ⬜ `pkg/connector/` — **待实现 ← 当前迭代**
- ⬜ `auth_data` 加密 — **待实现**
- ⬜ 密钥环境变量 — **待实现**
