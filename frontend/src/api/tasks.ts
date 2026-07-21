import api from './client';
import type { BackendTask, CreateTaskRequest } from '../types';

export interface TasksApi {
  list(params?: { status?: string; resource_id?: string }): Promise<BackendTask[]>;
  get(id: string): Promise<BackendTask>;
  create(data: CreateTaskRequest): Promise<BackendTask>;
  approve(id: string, comment?: string): Promise<BackendTask>;
  reject(id: string, comment?: string): Promise<BackendTask>;
  delete(id: string): Promise<void>;
}

export const tasksApi: TasksApi = {
  list: (params) => {
    const query = params
      ? '?' + new URLSearchParams(Object.entries(params).filter(([_, v]) => v != null) as [string, string][]).toString()
      : '';
    return api.get<BackendTask[]>(`/admin/tasks${query}`);
  },

  get: (id) => api.get<BackendTask>(`/admin/tasks/${id}`),

  create: (data) => api.post<BackendTask>('/admin/tasks', data),

  approve: (id, reviewer_comment) =>
    api.put<BackendTask>(`/admin/tasks/${id}/approve`, { reviewer_comment }),

  reject: (id, reviewer_comment) =>
    api.put<BackendTask>(`/admin/tasks/${id}/reject`, { reviewer_comment }),

  delete: (id) => api.delete<void>(`/admin/tasks/${id}`),
};

export default tasksApi;
