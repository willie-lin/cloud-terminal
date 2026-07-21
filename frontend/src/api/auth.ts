import api from './client';
import type { LoginRequest, LoginResponse, RegisterRequest } from '../types';

export interface AuthApi {
  login(data: LoginRequest): Promise<LoginResponse>;
  register(data: RegisterRequest): Promise<LoginResponse>;
  logout(): Promise<{ message: string }>;
  checkEmail(email: string): Promise<{ exists: boolean }>;
}

export const authApi: AuthApi = {
  login: (data) => api.post<LoginResponse>('/api/login', data),

  register: (data) => api.post<LoginResponse>('/api/register', data),

  logout: () => api.post<{ message: string }>('/api/logout'),

  checkEmail: (email) => api.post<{ exists: boolean }>('/api/check-email', { email }),
};

export default authApi;
