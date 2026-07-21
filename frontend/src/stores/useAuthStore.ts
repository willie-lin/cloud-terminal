import { create } from 'zustand';
import type { User, LoginRequest, RegisterRequest } from '../types';
import { authApi } from '../api/auth';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;

  /** Login with email + password to the real backend. */
  login: (data: LoginRequest) => Promise<void>;
  /** Register a new account. */
  register: (data: RegisterRequest) => Promise<void>;
  /** Logout and clear session. */
  logout: () => Promise<void>;

  /** For development: use a preset mock user without calling the backend. */
  loginAsPreset: (account: 'admin' | 'user') => void;
  /** For development: quick login with a username (no real auth). */
  loginWithUsername: (username: string) => void;

  clearError: () => void;
  setUser: (user: User) => void;
}

export const PRESET_ADMIN: User = {
  id: 'usr-admin-01',
  username: 'Alex Vance',
  email: 'alex.vance@cloudterminal.io',
  avatarUrl: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150&auto=format&fit=crop&q=80',
  role: 'admin',
};

export const PRESET_USER: User = {
  id: 'usr-dev-02',
  username: 'Max Jensen',
  email: 'max.jensen@acme.co',
  avatarUrl: 'https://images.unsplash.com/photo-1539571696357-5a69c17a67c6?w=150&auto=format&fit=crop&q=80',
  role: 'user',
};

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,

  login: async (data) => {
    set({ isLoading: true, error: null });
    try {
      const res = await authApi.login(data);
      const u = res.user;
      const roleName = u.roleName || 'user';
      const mappedUser: User = {
        id: u.id,
        tenantId: u.tenantId,
        groupId: u.groupId,
        email: u.email,
        username: u.username,
        roleName: roleName,
        role: (roleName === 'super_admin' || roleName === 'admin' || roleName === 'tenant_admin')
          ? (roleName as User['role'])
          : 'user',
        isTenantAdmin: u.isTenantAdmin,
        isSuperAdmin: u.isSuperAdmin,
        avatarUrl: `https://ui-avatars.com/api/?name=${encodeURIComponent(u.username)}&background=random`,
      };
      set({ user: mappedUser, isAuthenticated: true, isLoading: false });
    } catch (err: any) {
      const msg = err?.message || 'Login failed';
      set({ error: msg, isLoading: false });
      throw err;
    }
  },

  register: async (data) => {
    set({ isLoading: true, error: null });
    try {
      const res = await authApi.register(data);
      const u = res.user;
      const roleName = u.roleName || 'user';
      const mappedUser: User = {
        id: u.id,
        tenantId: u.tenantId,
        groupId: u.groupId,
        email: u.email,
        username: u.username,
        roleName: roleName,
        role: (roleName === 'super_admin' || roleName === 'admin' || roleName === 'tenant_admin')
          ? (roleName as User['role'])
          : 'user',
        isTenantAdmin: u.isTenantAdmin,
        isSuperAdmin: u.isSuperAdmin,
        avatarUrl: `https://ui-avatars.com/api/?name=${encodeURIComponent(u.username)}&background=random`,
      };
      set({ user: mappedUser, isAuthenticated: true, isLoading: false });
    } catch (err: any) {
      const msg = err?.message || 'Registration failed';
      set({ error: msg, isLoading: false });
      throw err;
    }
  },

  logout: async () => {
    try {
      await authApi.logout();
    } catch {
      // ignore logout errors
    }
    set({ user: null, isAuthenticated: false, error: null });
  },

  loginAsPreset: (account) => {
    set({
      isAuthenticated: true,
      user: account === 'admin' ? PRESET_ADMIN : PRESET_USER,
      error: null,
    });
  },

  loginWithUsername: (username) => {
    const isMatchedAdmin = username.toLowerCase().includes('admin') || username.toLowerCase().includes('alex');
    set({
      isAuthenticated: true,
      user: {
        id: `usr-${Date.now()}`,
        username,
        email: `${username.toLowerCase().replace(/\s+/g, '.')}@cloudterminal.io`,
        role: isMatchedAdmin ? 'admin' : 'user',
        avatarUrl: `https://ui-avatars.com/api/?name=${encodeURIComponent(username)}&background=random`,
      },
      error: null,
    });
  },

  clearError: () => set({ error: null }),

  setUser: (user) => set({ user, isAuthenticated: true }),
}));
