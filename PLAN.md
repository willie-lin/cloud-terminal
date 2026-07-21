# Cloud Terminal 转型方案：舍弃 Guacamole → 接入 ContainerSSH

> 作者：项目分析
> 日期：2026-07-15

---

## 一、当前项目完整状态评估

### 1.1 总体完成度：约 65%

通过对 `app/database/ent/schema/` 下所有 9 个文件的读取，确认：

| 模块 | 完成度 | 说明 |
|:---|:---:|:---|
| Ent 数据模型 | 80% | 8 实体 + Mixins，但 AccessPolicy 和 AuditLog 缺业务字段 |
| Handler 业务逻辑 | 70% | 核心 CRUD 有代码，但 session/accesspolicy 逻辑薄弱 |
| 安全层 (JWT/CSRF/Casbin) | 90% | 完善且可用 |
| Guacamole 集成 | **0%** | 占位符代码，从未真正接入 |
| 前端 Vue 3 | 30% | 框架搭好，无实质业务页面 |
| 配置系统 | 20% | 无 `.env`/`config.yaml` 配置文件 |
| Docker 部署 | 50% | Dockerfile 存在但未验证 |

### 1.2 Ent 实体全景

| 实体 | 字段状态 | 边缘关系 | 完成度 |
|:---|:---|:---|:---:|
| **Platform** | id(uuid), name, description, region, version, status(enum), config(JSON) | → Tenant | ✅ |
| **Tenant** | id(uuid), name, display_name, description, status(enum) | ← Platform → User/Role/Resource/Account/AccessPolicy/AuditLog | ✅ |
| **User** | id(uuid), username, password, email, phone, avatar, description, status, is_locked, last_login, expires_at, totp_secret, totp_enabled, refresh_token | ← Tenant | ✅ |
| **Role** | id(uuid), name, description, status | ← Tenant | ✅ |
| **Resource** | id(uuid), name, host, port(22), type(ssh/rdp/vnc/telnet), description, status, metadata(JSON) | ← Tenant → Account | ✅ |
| **Account** | id(uuid), name, credential_type, credential, privilege | ← Tenant, ← Resource | ✅ |
| **AccessPolicy** | **仅 id(uuid)** | ← Tenant | ❌ |
| **AuditLog** | **仅 id(uuid)** | ← Tenant | ❌ |
| **Mixins** | TimeMixin(created_at/updated_at) | — | ✅ |

### 1.3 Handler 层状态

| Handler | 完成度 | 备注 |
|:---|:---:|:---|
| user.go | 80% | CRUD + 启用/禁用 + CSRF Token |
| tenant.go | 80% | CRUD |
| role.go | 70% | CRUD |
| resource.go | 70% | CRUD |
| session.go | 50% | Guacamole 桩代码（将删除） |
| accesspolicy.go | 40% | 骨架，待增强 |
| 2FA.go | 60% | TOTP 基础逻辑 |
| hello.go | 100% | 健康检查 |
| real-IP.go | 100% | 真实 IP |

### 1.4 安全中间件（均可复用）

| 组件 | 路径 | 状态 |
|:---|:---|:---:|
| **JWT** | `pkg/utils/jwt.go` | ✅ Token 签发/验证 |
| **CSRF** | `pkg/utils/csrf.go` | ✅ Double Submit Cookie |
| **Casbin** | `pkg/utils/casbinMiddleware.go` | ✅ RBAC 鉴权 |
| **Session** | `pkg/utils/session.go` | ✅ Gorilla Sessions |
| **Auth 中间件** | `app/middlewarers/Authenticate.go` | ✅ JWT 上下文注入 |
| **Zap 日志** | `app/logger/echozap.go` | ✅ Zap + Lumberjack |

### 1.5 Guacamole 层

`pkg/guacd/guacd.go` — **确认是占位符代码，从未真正接入**。可以直接删除，无任何影响。

### 1.6 前端 (Vue 3 + Vite)

| 文件 | 状态 | 说明 |
|:---|:---:|:---|
| package.json | ✅ | Vue 3.5 + Pinia + Vue Router + Axios + Vite 6 |
| main.ts | ✅ | 入口完整 |
| router/index.ts | ⚠️ | 仅 Home + About |
| stores/user.ts | ⚠️ | Pinia 骨架 |
| api/index.ts | ⚠️ | Axios 封装 |
| views/ | ❌ | 无登录/终端/管理页面 |

---

## 二、ContainerSSH-0.6.0 源码分析

### 2.1 目录结构

```
ContainerSSH-0.6.0/
├── main.go                    # containerssh.Main() 入口
├── cmd/containerssh/          # CLI 命令
├── config/                    # 配置结构
│   ├── auth.go                # Auth 配置
│   ├── docker.go              # Docker 后端配置
│   ├── kubernetes.go          # K8s 后端配置
│   ├── ssh.go                 # SSH 服务配置
│   └── webhook.go             # Webhook 配置
├── message/                   # Webhook 消息结构
│   ├── auth.go                # AuthRequest / AuthResponse
│   └── configwebhook.go       # ConfigRequest / ConfigResponse
├── auth/webhook/              # Auth 客户端
├── config/webhook/            # Config 客户端
├── http/                      # HTTP 处理接口
├── internal/                  # SSH Server、代理等
└── service/                   # 服务管理
```

### 2.2 Webhook 认证流程

```
用户 SSH 连接
       │
       ▼
ContainerSSH (SSH Server)
       │
       ├── HTTP POST → Echo API: /api/webhook/auth
       │     {
       │       "username": "admin",
       │       "password": "xxx",
       │       "remote_address": "10.0.0.1",
       │       "connection_id": "abc123"
       │     }
       │     ← Response: { "authenticated": true, "user": "admin" }
       │
       └── HTTP POST → Echo API: /api/webhook/config
             {
               "username": "admin",
               "session_id": "sess_001",
               "remote_address": "10.0.0.1"
             }
             ← Response: { "config": { "docker": { "image": "ubuntu:22.04" } } }
```

### 2.3 核心消息结构（从源码提取）

```go
// Auth 请求/响应
type AuthRequest struct {
    Username      string `json:"username"`
    Password      string `json:"password,omitempty"`
    PublicKey     string `json:"public_key,omitempty"`
    RemoteAddress string `json:"remote_address"`
    ConnectionID  string `json:"connection_id"`
}

type AuthResponse struct {
    Authenticated bool   `json:"authenticated"`
    User          string `json:"user"`
}

// Config 请求/响应
type ConfigRequest struct {
    Username      string `json:"username"`
    SessionID     string `json:"session_id"`
    RemoteAddress string `json:"remote_address"`
}

type ConfigResponse struct {
    Config json.RawMessage `json:"config"`
}
```

---

## 三、转型方案：舍弃 Guacamole，接入 ContainerSSH

### 3.1 核心策略：做"一删三增"

```
删: pkg/guacd/                       ← Guacamole 占位符（从未接入）
增: POST /api/webhook/auth           ← ContainerSSH 认证
增: POST /api/webhook/config         ← ContainerSSH 容器配置
增: GET  /ws/terminal                ← 浏览器 WebSocket ↔ SSH 桥接
```

Go + Echo v5 + Ent 后端 **完整保留，不动骨架**。安全中间件（JWT/CSRF/Casbin）**全部复用**。

### 3.2 架构总览

```
Web Browser (React/xterm.js)                    SSH Client
         │                                          │
         │ WebSocket                                │ SSH
         ▼                                          ▼
┌────────────────────────────────────────────────────────┐
│                 Echo v5 API Server                     │
│                                                        │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ JWT Auth    │  │ Casbin RBAC  │  │ CSRF/Session │  │
│  └─────────────┘  └──────────────┘  └──────────────┘  │
│                                                        │
│  ┌──────────────────────────────────────────────────┐  │
│  │         新增的 3 个端点                            │  │
│  │                                                    │  │
│  │  POST /api/webhook/auth  ←── ContainerSSH 调用    │  │
│  │  POST /api/webhook/config ←── ContainerSSH 调用    │  │
│  │  GET  /ws/terminal        ──→ ContainerSSH:2222   │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────┬───────────────────────────────┘
                         │
                         ▼
                   ContainerSSH:2222
                         │
                         ├── Docker/K8s API
                         ▼
                   │  Guest Container │  (临时，用完即毁)
                         │
                         ▼
                     S3 录像 (.cast)
```

### 3.3 详细的文件变动清单

| 操作 | 文件 | 改动内容 |
|:---:|:---|:---|
| 🗑️ **删除** | `pkg/guacd/` 整个目录 | Guacamole 占位符 |
| ✏️ **修改** | `app/database/ent/schema/accesspolicy.go` | 补充 name, description, policy_type, priority, rules, environment_id 等字段 |
| ✏️ **修改** | `app/database/ent/schema/auditlog.go` | 补充 session_id, user_id, username, action, result, started_at, ended_at, detail, s3_path 等字段 |
| ✏️ **修改** | `go.mod` | 移除 Guacamole 依赖 |
| ✏️ **修改** | `app/api/api.go` | 增加 3 个新路由 |
| ✏️ **修改** | `app/handler/session.go` | 移除 Guacamole 引用 |
| ✏️ **修改** | `app/database/ent/schema/tenant.go` | 增加 `Edge.To("environments", Environment.Type)` |
| ✏️ **修改** | `main.go` | 路由引用更新 |
| 📄 **新增** | `app/database/ent/schema/environment.go` | 新的环境模板实体 |
| 📄 **新增** | `app/handler/webhook.go` | Webhook auth + config Handler |
| 📄 **新增** | `app/handler/wsbridge.go` | WebSocket ↔ SSH 桥接 Handler |
| 📄 **新增** | `config/config.yaml` | 统一配置文件 |
| 📄 **新增** | `docker-compose.yml` (更新) | 加入 ContainerSSH 服务 |

### 3.4 AccessPolicy Schema 增强

```go
// 当前：仅 id(uuid)
// 需要补充：
field.String("name").NotEmpty(),
field.String("description").Optional(),
field.Enum("policy_type").Values("allow", "deny").Default("allow"),
field.Int("priority").Default(0),
field.JSON("rules", map[string]interface{}{}).Optional(),
field.UUID("environment_id", uuid.UUID{}).Optional(),
```

### 3.5 AuditLog Schema 增强

```go
// 当前：仅 id(uuid)
// 需要补充：
field.String("session_id").Unique(),
field.UUID("user_id", uuid.UUID{}),
field.String("username"),
field.String("action"),
field.String("result"),
field.Time("started_at"),
field.Time("ended_at").Optional().Nillable(),
field.JSON("detail", map[string]interface{}{}).Optional(),
field.String("s3_path").Optional(),
```

### 3.6 Environment Schema（新增）

```go
// 新的环境模板实体
type Environment struct { ent.Schema }

func (Environment) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
        field.String("name").Unique().NotEmpty(),
        field.String("description").Optional(),
        field.String("image").NotEmpty(),                          // 容器镜像
        field.Int("port").Default(22),                             // SSH 端口
        field.JSON("resource_limit", map[string]interface{}{}).Optional(),
        field.JSON("env_vars", map[string]interface{}{}).Optional(),
        field.JSON("volumes", []map[string]interface{}{}).Optional(),
        field.Enum("status").Values("active", "inactive").Default("active"),
    }
}

func (Environment) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("tenant", Tenant.Type).Ref("environments").Unique(),
        edge.To("access_policies", AccessPolicy.Type),
    }
}
```

---

## 四、分阶段执行计划

### Phase 1：清理与 Schema 增强（0.5 天）

| 步骤 | 内容 | 产出 |
|:---:|:---|:---:|
| 1.1 | 删除 `pkg/guacd/` 整个目录 | 清理无用的 Guacamole 代码 |
| 1.2 | 从 `go.mod` 移除 guacamole 依赖 | 减少依赖数量 |
| 1.3 | 增强 `AccessPolicy` Schema | 完善字段 |
| 1.4 | 增强 `AuditLog` Schema | 完善字段 |
| 1.5 | 新增 `Environment` Schema | 新的环境模板实体 |
| 1.6 | 在 `Tenant` 增加 `Edge.To("environments")` | 关联关系 |
| 1.7 | 运行 `go generate ./...` 重新生成 ent client | 编译通过 |

### Phase 2：ContainerSSH Webhook Handler（1 天）

| 步骤 | 内容 | 产出 |
|:---:|:---|:---:|
| 2.1 | 新增 `app/handler/webhook.go` | Webhook 处理器 |
| 2.2 | 实现 `WebhookAuth()` — 查 User，验证密码，返回 {authenticated, user} | Auth 端点 |
| 2.3 | 实现 `WebhookConfig()` — 查 User → Tenant → AccessPolicy → Environment，返回 Docker 配置 JSON | Config 端点 |
| 2.4 | 在 `api.go` 注册路由（不经过 JWT，通过 ContainerSSH IP 白名单保护） | 路由就绪 |
| 2.5 | 创建 `config/config.yaml`（数据库连接、JWT 密钥、ContainerSSH 地址等） | 配置系统 |

### Phase 3：WebSocket ↔ SSH 桥接（1 天）

| 步骤 | 内容 | 产出 |
|:---:|:---|:---:|
| 3.1 | 新增 `app/handler/wsbridge.go` | WS 桥接处理器 |
| 3.2 | WebSocket 升级 + `golang.org/x/crypto/ssh` 连接 ContainerSSH:2222 | 双向桥接 |
| 3.3 | 两个 Goroutine 做 `io.Copy` 双向流绑定 | 数据流通 |
| 3.4 | `context.Context` + `defer` 确保断开时 Goroutine 正确回收 | 防泄露 |
| 3.5 | 在 `api.go` 注册 `/ws/terminal`（需 JWT 认证） | 路由就绪 |

### Phase 4：前端增强（2-3 天，可并行）

| 步骤 | 内容 | 产出 |
|:---:|:---|:---:|
| 4.1 | 增加登录页面 | 用户认证界面 |
| 4.2 | 增加终端页面（xterm.js 封装） | Web 终端 |
| 4.3 | 增加管理面板（用户/策略/环境管理） | 管理界面 |
| 4.4 | (远期) 增加录像回放页面（asciinema-player） | 审计回放 |

### Phase 5：部署配置（1 天）

| 步骤 | 内容 | 产出 |
|:---:|:---|:---:|
| 5.1 | `docker-compose.yml` 加入 ContainerSSH 服务 | 容器编排 |
| 5.2 | ContainerSSH 配置 Webhook URL 指向 Echo API | 认证通路 |
| 5.3 | 端到端测试：浏览器 → WS → Echo → SSH → ContainerSSH → Docker 容器 | MVP 验证 |

### 总投资

| Phase | 内容 | 工作量 |
|:---:|:---|:---:|
| P1 | 清理 + Schema 增强 | **0.5 天** |
| P2 | Webhook Handler | **1 天** |
| P3 | WebSocket 桥接 | **1 天** |
| P4 | 前端增强 | **2-3 天** |
| P5 | 部署配置 | **1 天** |
| | **MVP 总计** | **约 5-6 天** |

---

## 五、关键注意事项

1. **ContainerSSH 作为独立容器运行** — 不作为 Go 库导入，通过 Webhook URL 配置与 Echo API 交互
2. **AccessPolicy 和 AuditLog 补字段是 P1 必须做的** — 否则无法对接策略和审计
3. **前端建议保留 Vue 3 项目** — 先在现有项目中补页面，等后续再迁移 React + Tailwind v4
4. **go:embed 单二进制是远期目标** — MVP 阶段用 Docker Compose + Nginx 反向代理
5. **Goroutine 泄露风险** — WebSocket 桥接中必须用 `context.Context` 确保资源正确回收
6. **Webhook 端点需要 IP 白名单** — `/api/webhook/auth` 和 `/api/webhook/config` 应只允许 ContainerSSH 容器 IP 访问

---

## 六、总结

| 核心结论 | 说明 |
|:---|:---|
| **后端骨架完整可复用** | Echo + Ent + JWT/Casbin 全部保留 |
| **只需"一删三增"** | 删 Guacamole + 增 Auth/Config/WS 三个端点 |
| **Ent Schema 需局部增强** | AccessPolicy/AuditLog 补字段 + 新增 Environment |
| **ContainerSSH 做 sidecar** | 独立容器运行，通过 Webhook 交互 |
| **MVP 可在 5-6 天内交付** | 含完整认证 + 终端连接 + 审计记录 |

---

## 七、ContainerSSH 关系详解（基于 CNCF 官方用例理解）

### 7.1 ContainerSSH 是什么

**ContainerSSH** 是一个 CNCF Sandbox 项目（https://containerssh.io/），它本质是一个 **SSH 服务器 + 容器编排引擎**。它不是传统的跳板机或 CMDB 堡垒机，而是一个作为 Runtime 的“SSH 驱动临时容器管理器”。

### 7.2 和我们的项目是什么关系

```
┌──────────────────────────────────────────────────────────────┐
│                  Cloud Terminal (我们的项目)                   │
│                                                              │
│  角色：控制面 (Control Plane) + 管理面 (Management Plane)      │
│  - 多租户管理 (Platform → Tenant → User)                      │
│  - REST API (用户/策略/环境/审计 CRUD)                         │
│  - JWT + Casbin 认证授权                                       │
│  - 2FA (TOTP) 双因子认证                                       │
│  - WebSocket ↔ SSH 桥接                                       │
│  - S3 录像管理 + asciinema 回放                                │
│  - React 前端 Dashboard                                        │
└───────────────────────┬──────────────────────────────────────┘
                        │ Webhook (HTTP)
                        ▼
┌──────────────────────────────────────────────────────────────┐
│                  ContainerSSH (CNCF Sandbox)                  │
│                                                              │
│  角色：SSH 引擎 (Engine Layer)                                 │
│  - SSH Server (端口 2222)                                     │
│  - 认证 Webhook 客户端 (回调我们的 API)                         │
│  - 配置 Webhook 客户端 (回调我们的 API)                         │
│  - Docker/K8s 容器生命周期管理                                 │
│  - 会话桥接 (SSH ↔ Container)                                  │
│  - 资源限制 (CPU/内存/网络/磁盘 IO)                            │
└───────────────────────┬──────────────────────────────────────┘
                        │ Docker / K8s API
                        ▼
┌──────────────────────────────────────────────────────────────┐
│                    Guest Container (临时容器)                  │
│                                                              │
│  角色：用户操作环境                                           │
│  - 每个会话一个独立容器                                        │
│  - 用完即毁 (会话结束自动删除)                                 │
│  - 隔离的文件系统、网络、进程                                   │
│  - 可配置资源上限                                              │
└──────────────────────────────────────────────────────────────┘
```

### 7.3 关系本质：控制面 + 引擎层，松耦合

| 维度 | Cloud Terminal (我们) | ContainerSSH (CNCF) |
|:---|:---|:---|
| **职责** | 业务逻辑、鉴权决策、审计管理 | SSH 协议处理、容器编排 |
| **代码** | 我们编写的 Go 代码 | 官方维护的 CNCF 项目 |
| **部署** | Echo API 容器 | 独立 sidecar 容器 |
| **通信** | Webhook (ContainerSSH → Echo API) | Docker/K8s API |
| **升级** | 独立升级 | 换镜像 tag 即可 |
| **优势** | 多租户、RBAC、2FA、审计 | SSH 协议成熟、CNCF 社区 |

**一句话**：ContainerSSH 只关心"谁通过 SSH 连进来了，给他起个什么容器"，我们关心"这个用户属于哪个租户、有什么权限、用完留下什么审计记录"。两者通过 Webhook 协议松耦合通信。

---

## 八、ContainerSSH 官方 6 种用例映射到我们的项目

根据 ContainerSSH 官方文档（v0.6），它提供了 6 种经典用例场景，我们的项目可以一一对应：

### 8.1 实验室环境 (Lab Environment)

| ContainerSSH 能力 | 我们的实现 |
|:---|:---|
| 认证 Webhook | `/api/webhook/auth` 查 User 表 - bcrypt 验证 |
| 资源限制 Webhook | `/api/webhook/config` 返回 Docker CPU/内存限制 |
| 容器清理 | ContainerSSH 自动删除退出容器 |
| 审计监控 | SessionLog 实体 + asciinema 录像 → S3 |

**对应我们的场景**：企业给实习生/外包提供 SSH 环境，资源隔离，用完自动回收。

### 8.2 调试生产系统 (Debug Production)

| ContainerSSH 能力 | 我们的实现 |
|:---|:---|
| 时间限制 | AccessPolicy 规则中的 `time_constraint` |
| 审计记录 | AuditLog 实体记录所有 SSH 命令 + 文件传输 |
| 只读容器 | Config Webhook 返回配置限制写入权限 |
| 动态权限 | Casbin RBAC 策略控制谁能访问 |

**对应我们的场景**：运维人员限时访问生产服务器调试，全程录像审计。

### 8.3 虚拟主机 (Virtual Hosting)

| ContainerSSH 能力 | 我们的实现 |
|:---|:---|
| 用户隔离 | Environment 实体 + Docker 容器隔离 |
| SFTP 支持 | ContainerSSH 原生支持 |
| 自定义镜像 | Environment 实体中的 `image` 字段 |
| 多服务器 | Platform 实体管理多个部署|

**对应我们的场景**：为每个开发者提供独立的开发环境容器。

### 8.4 学习环境 (Learning Environment)

| ContainerSSH 能力 | 我们的实现 |
|:---|:---|
| 按需启动 | Environment 模板 → 自动拉起容器 |
| 资源限制 | Docker 配置限制学生每个容器的资源 |
| 用完即毁 | 会话结束自动删除 |
| 自定义镜像 | 老师可预设镜像 (带 MySQL/MongoDB 等工具) |

**对应我们的场景**：培训平台、考试系统、在线实验室。

### 8.5 蜜罐 (Honeypot)

| ContainerSSH 能力 | 我们的实现 |
|:---|:---|
| 记录攻击者凭据 | Auth Webhook 记录失败登录 |
| 隔离容器 | 蜜罐容器无网络/受限资源 |
| 审计日志 | 记录所有攻击者操作 |
| 防火墙联动 | 可集成 IP 黑名单 |

**对应我们的场景**：安全团队监控 SSH 攻击，记录攻击者行为。

### 8.6 安全访问 (Secure Access with Vault)

| ContainerSSH 能力 | 我们的实现 |
|:---|:---|
| 短生命周期凭据 | Config Webhook 可集成 Vault |
| 密码/密钥认证 | User Schema 支持密码 + 公钥 |
| 审计日志 | 记录所有凭证使用 |

**对应我们的场景**：高安全环境，短期凭证 + 双因子认证。

---

## 九、升级后的完整部署架构

### 9.1 Sidecar 部署模式（推荐）

由于 ContainerSSH 是 CNCF 项目，我们 **不修改 ContainerSSH 源码**，而是将其作为独立 sidecar 容器运行：

```yaml
# docker-compose.yml 核心服务
services:
  postgres:      # PostgreSQL 数据库
  api:           # 我们的 Echo API + React 前端 (go:embed)
  containerssh:  # CNCF 项目，官方镜像
  minio:         # S3 兼容存储 (MinIO)
  nginx:         # TLS 终止 + 反向代理 (可选)
```

### 9.2 架构全景图

```
                         ┌───────────────┐
                         │  用户浏览器    │
                         │  (xterm.js)   │
                         └───────┬───────┘
                                 │ WebSocket (wss://)
                                 ▼
┌─────────────────────────────────────────────────────────────┐
│                     Nginx (端口 443)                         │
│                     TLS 终止 + 反向代理                       │
└────────┬────────────────────────────────────┬────────────────┘
         │                                    │
         ▼                                    ▼
┌──────────────────┐                ┌────────────────────┐
│  React 前端      │                │  Echo API Server   │
│  (go:embed)      │                │  端口 :8080         │
│                  │                │                     │
│  - 登录页        │                │  REST API 管理端点  │
│  - Dashboard     │                │  /api/users/*       │
│  - 终端页        │                │  /api/tenants/*     │
│  - 回放页        │                │  /api/policies/*    │
│  - 审计页        │                │  /api/environments/*│
└──────────────────┘                │                     │
                                    │  Webhook 端点:      │
                                    │  POST /webhook/auth  │ ← ContainerSSH 调用
                                    │  POST /webhook/config│ ← ContainerSSH 调用
                                    │                     │
                                    │  WebSocket 桥接:    │
                                    │  GET  /ws/terminal   │ → SSH → ContainerSSH:2222
                                    └────────┬────────────┘
                                             │ Webhook (HTTP)
                                             ▼
                              ┌────────────────────────────┐
                              │  ContainerSSH (CNCF)       │
                              │  端口 :2222                │
                              │                            │
                              │  - SSH Server              │
                              │  - Auth Webhook Client     │
                              │  - Config Webhook Client   │
                              │  - Audit Logger            │
                              └────────┬───────────────────┘
                                       │ Docker / K8s API
                                       ▼
                              ┌────────────────────────────┐
                              │  Guest Container (临时)     │
                              │                            │
                              │  - 用完即毁                 │
                              │  - 资源隔离                 │
                              │  - 操作录像 → S3            │
                              └────────────────────────────┘
```

### 9.3 三种访问路径

| 路径 | 用户 | 流程 |
|:---:|:---|:---|
| **Web 终端** | 浏览器用户 | xterm.js → WebSocket → Echo API → SSH → ContainerSSH:2222 → Docker 容器 |
| **直接 SSH** | 高级用户 | ssh user@host -p 2222 → ContainerSSH (Webhook → Echo API) → Docker 容器 |
| **SFTP 文件** | 文件传输 | sftp user@host -p 2222 → ContainerSSH (原生支持) → 容器挂载卷 |

---

## 十、利用 ContainerSSH 作为 CNCF 项目的优势

| 优势 | 说明 | 对我们项目的影响 |
|:---|:---|:---|
| **社区维护** | CNCF Sandbox 项目，社区活跃 | 不需要自研 SSH 网关，聚焦业务层 |
| **Webhook 松耦合** | 认证/配置通过 HTTP 协议交互 | 我们的 API 和 ContainerSSH 独立部署、独立升级 |
| **多后端支持** | Docker / K8s / Podman | 初期 Docker MVP，后期 K8s 扩展 |
| **即用即毁** | 每个会话独立容器，用完删除 | 天然零信任模型，无需清理遗留 |
| **SFTP 原生支持** | 内置文件传输支持 | 无需额外实现文件传输功能 |
| **资源限制** | CPU/内存/磁盘 IO/网络 | 通过 Config Webhook 按用户/策略动态配置 |
| **审计日志** | 详细记录 SSH 命令和文件传输 | 配合 S3 + asciinema 实现审计回放 |
| **官方文档完善** | 6 种成熟用例参考 | 直接参照文档实现，降低设计风险 |

---

## 十一、总结：优化执行总览

| 项目 | 操作 | 工作量 |
|:---|:---|:---:|
| 删除 `pkg/guacd/` | 🗑️ 直接删除 | 5 分钟 |
| 增强 AccessPolicy Schema | ✏️ 补字段 | 30 分钟 |
| 增强 AuditLog Schema | ✏️ 补字段 | 30 分钟 |
| 新增 Environment Schema | 📄 新建 | 30 分钟 |
| 新增 Webhook Handler | 📄 新建 (auth + config) | 1 天 |
| 新增 WS 桥接 Handler | 📄 新建 | 1 天 |
| 前端增强 | ✏️ 补页面 (xterm.js) | 2-3 天 |
| 配置系统 | 📄 新建 config.yaml | 30 分钟 |
| Docker Compose | ✏️ 加入 ContainerSSH | 1 天 |
| **MVP 总计** | | **约 5-6 天** |

---

*文档版本：v1.1*
*补充内容：ContainerSSH 官方 6 种用例映射 + 部署架构 + CNCF 关系分析*
