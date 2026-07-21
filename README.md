# Cloud-Terminal

> **A cloud-native access control plane that manages secure, time-bound, and auditable access to infrastructure through abstract resources instead of exposing servers directly.**
> **一个云原生基础设施访问控制平面，通过“资源”而不是服务器，为用户提供安全、临时、可审计的访问能力。**

**Approved SSH Access** — Connect to target infrastructure securely without managing static SSH keys or IP permissions. Ephemeral tokens are issued on-demand and executed inside isolated Docker container sandboxes.

我们不是传统跳板机或 CMDB，也不做 Teleport 或 JumpServer 的复刻。Cloud-Terminal 是一个面向云原生时代的 **SSH Access Control Plane（SSH 控制平面）**：让用户通过“资源 (Resource)”访问基础设施，而不暴露底层服务器 IP 与静态凭据。

---

## 架构：两层设计

| 层面 | 职责 | 技术栈 |
|:---|:---|:---|
| **控制面** | 用户管理、资源管理、权限评估、Webhook/Config 生成 | Echo API (Go) |
| **执行面** | 动态创建临时容器、执行用户操作、审计录屏 | ContainerSSH (CNCF) |

---

## 核心业务流程

![Architecture](https://via.placeholder.com/800x400?text=Architecture+Diagram)

### 场景一：用户通过 Web UI 访问 VPC 内 MySQL

```mermaid
sequenceDiagram
    actor User as 运维人员
    participant UI as Cloud-Terminal Web UI
    participant API as Cloud-Terminal API
    participant DB as 平台数据库
    participant CSSH as ContainerSSH (部署在 VPC-A)
    participant Container as 临时诊断容器 (VPC-A)
    participant MySQL as 目标 MySQL (VPC-A)

    User->>UI: 登录 Web UI
    UI->>API: 认证 (JWT + TOTP)
    API->>DB: 验证身份
    DB-->>API: 通过
    API-->>UI: 返回 Token + 资源列表

    User->>UI: 点击 "VPC-A / MySQL-主库"
    UI->>API: 申请访问资源
    API->>DB: 检查 AccessPolicy
    DB-->>API: 策略评估通过

    Note over API: 生成 ContainerSSH 配置
    Note over API:   - 镜像: mysql-client-tools:latest
    Note over API:   - 安全组: 允许 3306
    Note over API:   - 超时: 60 分钟后自动销毁

    API->>CSSH: 下发 Webhook Config
    CSSH-->>API: 确认

    API-->>UI: 返回 WebSocket 地址
    UI->>CSSH: WebSocket SSH 连接

    CSSH->>Container: Docker 创建临时容器
    Note over Container: 安全组: 仅允许 :3306<br/>资源: 0.5CPU / 512MB<br/>有效期: 60min
    Container-->>CSSH: 就绪
    CSSH-->>UI: SSH Session 建立
    UI-->>User: Web 终端 (xterm.js)

    User->>Container: 执行 mysql -h mysql.internal -u ops -p
    Container->>MySQL: TCP :3306
    MySQL-->>Container: 返回数据
    Container-->>User: 查询结果

    Note over Container: 60min 超时 / 主动退出
    CSSH->>Container: 停止并删除容器
    Note over CSSH: 审计日志 → S3
```

### 场景二：用户通过 Web UI 访问 K8s 集群内 Service

```mermaid
sequenceDiagram
    actor User as 运维人员
    participant UI as Cloud-Terminal Web UI
    participant API as Cloud-Terminal API
    participant CSSH as ContainerSSH (部署在 K8s 集群 B)
    participant Pod as 临时诊断 Pod (K8s Namespace ops-tools)
    participant SVC as 目标 Service (K8s Namespace-B)

    User->>UI: 登录，浏览资源
    User->>UI: 点击 "集群-B / redis-sentinel"

    API->>DB: 检查 AccessPolicy
    DB-->>API: 通过

    Note over API: 生成 ContainerSSH 配置
    Note over API:   - 镜像: redis-client-tools:latest
    Note over API:   - Namespace: ops-tools
    Note over API:   - NetworkPolicy: 允许 redis-sentinel:6379
    Note over API:   - 超时: 30 分钟

    API->>CSSH: 下发 Webhook Config
    CSSH->>Pod: K8s 创建 Pod (带 NetworkPolicy)
    Note over Pod: NetworkPolicy: 仅 6379<br/>资源: 0.5CPU / 256MB
    Pod-->>CSSH: 就绪
    API-->>UI: 返回 WebSocket 地址
    UI->>CSSH: WebSocket SSH 连接
    CSSH-->>UI: SSH Session 建立
    UI-->>User: Web 终端

    User->>Pod: redis-cli -h redis-sentinel
    Pod->>SVC: ClusterIP :6379
    SVC-->>Pod: 返回数据
    Pod-->>User: 结果

    Note over Pod: 用户退出 / 超时
    CSSH->>Pod: 删除 Pod
    Note over CSSH: 审计日志 → S3
```

### 场景三：高级用户直接 SSH（绕过 Web UI）

```mermaid
sequenceDiagram
    actor User as 高级用户
    participant SSH as SSH 客户端
    participant CSSH as ContainerSSH (目标网络内)
    participant API as Cloud-Terminal API
    participant Container as 临时容器

    User->>SSH: ssh ops@containerSSH-vpc-a:2222
    CSSH->>API: Webhook POST /auth
    Note over CSSH: { username, password, remote_addr, connection_id }
    API-->>CSSH: { authenticated: true, user: "ops" }

    CSSH->>API: Webhook POST /config
    Note over API: 根据 username + 策略决定容器配置
    API-->>CSSH: 返回容器配置

    CSSH->>Container: 创建临时容器
    Container-->>CSSH: 就绪
    CSSH-->>User: SSH Session（进入容器 Shell）
```

---

## 部署架构

```mermaid
graph TB
    subgraph "用户层"
        BROWSER[浏览器 Web UI]
        CLI[SSH 客户端]
    end

    subgraph "控制面"
        API[Cloud-Terminal API<br/>端口 :8080]
        DB[(MySQL)]
        API --> DB
    end

    subgraph "VPC-A 生产环境"
        CSSH_A[ContainerSSH<br/>端口 :2222]
        C_A[临时容器<br/>mysql-client<br/>安全组:3306]
        MYSQL[MySQL 目标服务]
        CSSH_A --> C_A --> MYSQL
    end

    subgraph "K8s 集群 B"
        CSSH_B[ContainerSSH<br/>端口 :2222]
        P_B[临时 Pod<br/>redis-client<br/>NetworkPolicy]
        SVC_B[Redis Service]
        CSSH_B --> P_B --> SVC_B
    end

    subgraph "VPC-C 测试环境"
        CSSH_C[ContainerSSH<br/>端口 :2222]
        C_C[临时容器<br/>通用工具<br/>安全组:自定义]
        TGT_C[测试服务]
        CSSH_C --> C_C --> TGT_C
    end

    BROWSER --- API
    API -- Webhook/Config --> CSSH_A
    API -- Webhook/Config --> CSSH_B
    API -- Webhook/Config --> CSSH_C
    BROWSER -- WebSocket --> CSSH_A
    BROWSER -- WebSocket --> CSSH_B
    BROWSER -- WebSocket --> CSSH_C
    CLI -- SSH :2222 --> CSSH_A
    CLI -- SSH :2222 --> CSSH_B
    CLI -- SSH :2222 --> CSSH_C
```

---

## 为什么不是 JumpServer / Teleport / 堡垒机？

传统平台（如 JumpServer / 传统堡垒机）的思路是围绕 **Host（服务器）** 与 CMDB 展开的：
```
CMDB → Host (服务器/IP) → 授权 → 直连登录
```
用户最终还是在管理和暴露底层主机。

**Cloud-Terminal 的定位不同，核心管理对象是抽象的 `Resource`（资源）与 `Task`（工单）：**
```
Task (工单) → Approval (审批) → Policy (动态授权) → Resource (目标服务) → ContainerSSH Runtime → 临时沙箱连接
```
用户管理的是“我要完成什么工作”（如：生产数据库、Redis 运维、线上日志），系统负责鉴权、下发临时凭据、并由 ContainerSSH 作为 Runtime 启动临时微隔离沙箱建立 SSH 桥接，用完即毁。

| 维度 | 传统平台（JumpServer / 堡垒机 / Teleport） | Cloud-Terminal（SSH 控制平面 + Cloud Gateway） |
|:---|:---|:---|
| **核心对象** | **Host / IP**（物理机、虚拟机、CMDB 主机） | **Resource**（抽象业务资源：MySQL、Redis、K8s Pod 等） |
| **凭据与授权** | 静态长效 SSH 密码 / 密钥授权 | **按需动态签发（STS 临时 Token）**，零静态凭据暴露 |
| **执行 Runtime** | SSH 代理直连目标持久化主机 | **ContainerSSH 驱动隔离沙箱**，会话结束自动销毁 (`rm -f`) |
| **网络边界** | 依赖静态安全组 / Host 打通 | **微隔离网络** + JIT 单次临时连接 |
| **用户认知** | “我要登录 10.10.10.23 主机” | **“我要维护生产 MySQL-主库”** |
| **多租户体系** | 弱（通常为多部门或资产授权分组） | **原生多租户隔离体系 (Tenant-Level Isolation)** |

---

## URN 标识体系

类似 AWS ARN，使用 URN 格式统一标识所有资源：

```
urn:cloud-term:{tenant_id}:{resource_type}:{resource_id}
```

示例：
- `urn:cloud-term:tenant-a:resource:mysql-001`
- `urn:cloud-term:tenant-a:user:zhangsan`
- `urn:cloud-term:tenant-a:platform:mysql`

---

## 数据模型

```mermaid
erDiagram
    Tenant ||--o{ Environment : has
    Tenant ||--o{ Resource : has
    Tenant ||--o{ AccessPolicy : has
    Account ||--o{ User : has
    Account ||--o{ Role : has
    Account ||--o{ AccessPolicy : has
    Account ||--o{ Resource : has
    User ||--o{ Role : has
    User ||--o{ AuditLog : has
    Role ||--o{ AccessPolicy : has
    AccessPolicy ||--o{ Environment : has
```

| 实体 | 说明 |
|:---|:---|
| **Tenant** | 租户（组织隔离单位） |
| **Account** | 账号（凭据管理） |
| **User** | 用户（平台登录者） |
| **Role** | 角色（权限集合） |
| **Resource** | 目标服务（MySQL、Redis 等） |
| **Environment** | 容器模板（镜像、资源限制、环境变量） |
| **AccessPolicy** | IAM 策略（谁、在何时、用什么工具、访问什么资源） |
| **Session** | 会话记录（用户 → 容器 → 目标服务的完整链路） |
| **AuditLog** | 操作审计日志 |

---

## 技术栈

| 层 | 技术 |
|:---|:---|
| 后端语言 | Go |
| Web 框架 | Echo v5 |
| ORM | Ent (entgo.io) |
| 数据库 | MySQL |
| 认证 | JWT + TOTP 双因素 |
| 授权 | Casbin RBAC + IAM 策略 (AccessPolicy) |
| SSH 网关 | ContainerSSH (CNCF) |
| 审计 | Session 录像 + S3 存储 |
| 前端 | React (UmiJS) |
| 终端 | xterm.js |

---

## 项目结构

```
.
├── api/            # API 路由（登录、注册、认证等）
├── ent/            # Ent ORM 代码（自动生成）
│   └── schema/     # Schema 定义
├── handler/        # 业务逻辑处理器
├── middlewarers/   # 认证中间件
├── config/         # 配置
├── rule/           # 权限规则
├── viewer/         # 视图上下文
├── pkg/
│   └── utils/      # 工具包（JWT、密码、CSRF 等）
└── web/            # 前端
```

---

## 当前实现状态

### ✅ 已完成
- [x] Ent ORM 全部 10 个 Schema
- [x] 多租户用户认证（JWT + TOTP）
- [x] Casbin RBAC 权限模型
- [x] 访问策略引擎
- [x] 资源 / 账号 / 平台 CRUD API
- [x] 会话管理与 WebSocket 终端
- [x] 审计日志记录
- [x] 前端框架（React + UmiJS + xterm.js）

### 🚧 待实现
- [ ] Webhook + Config 生成器（Cloud-Terminal → ContainerSSH）
- [ ] 动态容器编排（根据资源类型选择镜像）
- [ ] 安全组 / NetworkPolicy 动态生成
- [ ] 容器生命周期管理（限时自动销毁）
- [ ] 操作录播回放
- [ ] 多 VPC / 多集群管理
