/*
Package iam — AWS 风格策略鉴权引擎

完整实现了 AWS IAM 核心概念，可用于任何 Go 项目的细粒度访问控制。

━━━ 核心概念 ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Policy（策略）→ Statement（语句）→ Effect（Allow / Deny）
完整的评估链：SCP → PermissionBoundary → IdentityPolicy → ResourcePolicy → SessionPolicy

━━━ 架构层次（AWS IAM 完全对应） ━━━━━━━━━━━━━━━━━━━━━━━━

  Organization SCP（组织级控制策略）
    ↓ 任何 Deny 立即拒绝
  Permission Boundary（权限边界，最大权限集）
    ↓ 超出边界的行为拒绝
  Identity Policy（身份策略：User → Group → Role）
    ↓ 身份层面允许
  Resource Policy（资源策略：附加在目标资源上）
    ↓ 资源层面允许（可跨账号）
  Session Policy（会话策略：Role Assume 时传入的限制）
    ↓ 仅能缩小权限
  → 最终决定

━━━ 三定律评估 ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  1. 默认拒绝（Default Deny）
  2. 显式拒绝最高优先（Explicit Deny）
  3. 显式允许（Explicit Allow）

━━━ 支持的 Condition 运算符（12+ 种） ━━━━━━━━━━━━━━━━━

  StringEquals / StringNotEquals / StringEqualsIgnoreCase
  StringLike / StringNotLike（支持 * ? 通配）
  ArnEquals / ArnLike（URN 匹配）
  IpAddress / NotIpAddress（CIDR）
  DateGreaterThan / DateLessThan / DateEquals / DateBetween
  Bool
  NumericEquals / NumericNotEquals / NumericLessThan / NumericGreaterThan
  Null（检查 Key 是否存在）
  RequireMFA
  ForAllValues / ForAnyValue（Set 运算符）

━━━ 扩展能力 ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  - NotAction / NotResource 反向匹配
  - Policy 变量替换：${aws:username}, ${aws:CurrentTime}, ${resource:Urn}
  - Tag-based ABAC：PrincipalTag / ResourceTag / RequestTag
  - Simulate API：批量评估多个 Action/Resource 组合
  - Provider 接口：可对接 ent / Redis / 文件 / gRPC 等任意数据源
*/
package iam
