/**
 * API client for Cloud Terminal backend communication.
 * Handles CSRF tokens, cookie-based auth, and consistent error handling.
 */

const API_BASE = ''; // uses Vite proxy in dev, same origin in production

// ─── CSRF Token management ────────────────────────────────────

let csrfToken: string | null = null;

/** Extract CSRF token from the cookie set by the backend. */
function getCsrfFromCookie(): string | null {
  if (typeof document === 'undefined') return null;
  const match = document.cookie.match(/(?:^|;\s*)_csrf=([^;]+)/);
  return match ? decodeURIComponent(match[1]) : null;
}

/** Ensure we have a CSRF token (fetch one if missing). */
export async function ensureCsrfToken(): Promise<string> {
  if (csrfToken) return csrfToken;

  const fromCookie = getCsrfFromCookie();
  if (fromCookie) {
    csrfToken = fromCookie;
    return csrfToken;
  }

  // Fetch a fresh CSRF token from the backend
  const res = await fetch(`${API_BASE}/api/csrf-token`, {
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Failed to fetch CSRF token');

  const fromCookieAfter = getCsrfFromCookie();
  if (fromCookieAfter) {
    csrfToken = fromCookieAfter;
    return csrfToken;
  }

  throw new Error('CSRF token not set after fetch');
}

// ─── HTTP helpers ──────────────────────────────────────────────

export class ApiError extends Error {
  status: number;
  body: any;

  constructor(status: number, body: any) {
    super(typeof body === 'object' && body?.error ? body.error : `Request failed (${status})`);
    this.name = 'ApiError';
    this.status = status;
    this.body = body;
  }
}

interface RequestOptions {
  method?: string;
  body?: any;
  skipCsrf?: boolean;
  headers?: Record<string, string>;
}

async function request<T = any>(path: string, opts: RequestOptions = {}): Promise<T> {
  const { method = 'GET', body, skipCsrf = false, headers = {} } = opts;

  const reqHeaders: Record<string, string> = {
    'Content-Type': 'application/json',
    ...headers,
  };

  if (!skipCsrf) {
    const token = await ensureCsrfToken();
    reqHeaders['X-CSRF-Token'] = token;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    method,
    headers: reqHeaders,
    credentials: 'include',
    body: body ? JSON.stringify(body) : undefined,
  });

  // No content
  if (res.status === 204) return undefined as T;

  const contentType = res.headers.get('content-type') || '';
  const data = contentType.includes('application/json') ? await res.json() : await res.text();

  if (!res.ok) {
    throw new ApiError(res.status, data);
  }

  return data as T;
}

// ─── Public API ────────────────────────────────────────────────

export const api = {
  get: <T = any>(path: string, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: 'GET' }),

  post: <T = any>(path: string, body?: any, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: 'POST', body }),

  put: <T = any>(path: string, body?: any, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: 'PUT', body }),

  delete: <T = any>(path: string, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: 'DELETE' }),
};

export default api;
