// ─── Auth ────────────────────────────────────────────────────

export type UserRole = 'super_admin' | 'tenant_admin' | 'admin' | 'user';

export interface User {
  id: string;
  username: string;
  email: string;
  tenantId?: string;
  groupId?: string;
  roleName?: string;
  avatarUrl?: string;
  role: UserRole;
  isTenantAdmin?: boolean;
  isSuperAdmin?: boolean;
}

export interface LoginRequest {
  email: string;
  password: string;
  otp?: string;
}

export interface LoginResponse {
  accessToken: string;
  refreshToken: string;
  user: {
    id: string;
    tenantId: string;
    groupId: string;
    email: string;
    username: string;
    roleName: string;
    isTenantAdmin: boolean;
    isSuperAdmin: boolean;
  };
}

export interface RegisterRequest {
  email: string;
  password: string;
  tenant_name: string;
}

// ─── Resource ─────────────────────────────────────────────────

export type ResourceType = 'ssh' | 'mysql' | 'redis' | 'k8s-service' | 'rdp' | 'vnc' | 'telnet' | 'http' | 'custom';
export type ResourceEnv = 'prod' | 'staging' | 'dev' | 'test' | 'dr';
export type ResourceStatus = 'active' | 'inactive';

export interface Resource {
  id: string;
  urn: string;
  name: string;
  type: ResourceType;
  ip: string;
  port: number;
  env: ResourceEnv;
  region: string;
  description?: string;
  status: ResourceStatus;
  details?: Record<string, any>;
  auth_data?: Record<string, any>;
  host_key?: string;
  // For frontend display:
  targetHost?: string;
  targetPort?: number;
  authUsername?: string;
  os?: string;
}

// ─── STS Token ─────────────────────────────────────────────────

export interface AccessPolicyToken {
  token: string;
  resourceUrn: string;
  resourceName: string;
  resourceId?: string;
  issuedAt: number;
  expiresAt: number;
  durationMinutes: number;
}

// ─── Task / Approval ──────────────────────────────────────────

export type TaskStatus = 'pending' | 'approved' | 'rejected' | 'expired';

export interface AccessTask {
  id: string;
  userId: string;
  userName: string;
  userEmail: string;
  userAvatar?: string;
  resourceId: string;
  resourceUrn: string;
  resourceName: string;
  reason: string;
  durationMinutes: number;
  durationHours: number;
  status: TaskStatus;
  requestedAt: string;
  reviewedAt?: string;
  issuedToken?: AccessPolicyToken;
}

export interface CreateTaskRequest {
  resource_id: string;
  reason: string;
  duration_minutes: number;
}

// ─── Terminal Session ─────────────────────────────────────────

export interface TerminalSession {
  sessionId: string;
  resourceUrn: string;
  resourceName: string;
  token: string;
  connectedAt: number;
  expiresAt: number;
}

// ─── API response wrapper for paginated/array endpoints ───────

export interface ApiListResponse<T> {
  data: T[];
  total?: number;
}

// ─── Backend Task edge (with edges.requester + edges.resource) ─

export interface BackendTask {
  id: string;
  reason: string;
  duration_minutes: number;
  status: string;
  reviewed_at: string | null;
  reviewer_comment: string | null;
  issued_token: string | null;
  expires_at: string | null;
  created_at: string;
  updated_at: string;
  edges: {
    requester?: {
      id: string;
      username: string;
      email: string;
    };
    resource?: {
      id: string;
      urn: string;
      name: string;
    };
    reviewer?: {
      id: string;
      username: string;
    } | null;
  };
}

// ─── Backend Resource shape ───────────────────────────────────

export interface BackendResource {
  id: string;
  urn: string;
  name: string;
  type: string;
  ip: string;
  port: number;
  env: string;
  region: string;
  description: string | null;
  status: string;
  details: Record<string, any> | null;
  created_at: string;
  updated_at: string;
}
