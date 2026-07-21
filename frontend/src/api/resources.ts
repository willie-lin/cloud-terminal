import api from './client';
import type { BackendResource } from '../types';

export interface CreateResourceData {
  urn: string;
  name: string;
  type?: string;
  ip: string;
  port?: number;
  env?: string;
  region?: string;
  description?: string;
  status?: string;
  details?: Record<string, any>;
  auth_data?: Record<string, any>;
  host_key?: string;
}

export interface ResourceApi {
  list(): Promise<BackendResource[]>;
  get(id: string): Promise<BackendResource>;
  create(data: CreateResourceData): Promise<BackendResource>;
  update(id: string, data: Partial<CreateResourceData>): Promise<BackendResource>;
  delete(id: string): Promise<void>;
}

export const resourceApi: ResourceApi = {
  list: () => api.get<BackendResource[]>('/admin/resources'),

  get: (id) => api.get<BackendResource>(`/admin/resources/${id}`),

  create: (data) => api.post<BackendResource>('/admin/resources', data),

  update: (id, data) => api.put<BackendResource>(`/admin/resources/${id}`, data),

  delete: (id) => api.delete<void>(`/admin/resources/${id}`),
};

export default resourceApi;
