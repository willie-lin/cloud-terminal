import { create } from 'zustand';
import type { AccessTask, AccessPolicyToken, BackendTask } from '../types';
import { useResourceStore } from './useResourceStore';
import { tasksApi } from '../api/tasks';

interface ApprovalState {
  tasks: AccessTask[];
  isLoading: boolean;
  error: string | null;

  /** Fetch tasks from the backend API. */
  fetchTasks: () => Promise<void>;
  /** Create a new access task request. */
  createTaskRequest: (
    userId: string,
    userName: string,
    userEmail: string,
    resourceId: string,
    resourceUrn: string,
    resourceName: string,
    reason: string,
    durationHours: number
  ) => Promise<void>;
  /** Approve a pending task. */
  approveTask: (taskId: string) => Promise<void>;
  /** Reject a pending task. */
  rejectTask: (taskId: string) => Promise<void>;
}

/** Map a BackendTask from the API to the frontend AccessTask type. */
function mapBackendTask(bt: BackendTask): AccessTask {
  const requester = bt.edges?.requester;
  const resource = bt.edges?.resource;
  const isApproved = bt.status === 'approved';
  const durationMinutes = bt.duration_minutes;

  let issuedToken: AccessPolicyToken | undefined;
  if (isApproved && bt.issued_token) {
    const now = Math.floor(Date.now() / 1000);
    const expiresAt = bt.expires_at
      ? Math.floor(new Date(bt.expires_at).getTime() / 1000)
      : now + durationMinutes * 60;
    issuedToken = {
      token: bt.issued_token,
      resourceUrn: resource?.urn || '',
      resourceName: resource?.name || '',
      resourceId: resource?.id,
      issuedAt: now,
      expiresAt,
      durationMinutes,
    };
  }

  const createdDate = new Date(bt.created_at);
  const now = new Date();
  const diffMs = now.getTime() - createdDate.getTime();
  const diffMins = Math.floor(diffMs / 60000);
  let requestedAt: string;
  if (diffMins < 1) requestedAt = 'Just now';
  else if (diffMins < 60) requestedAt = `${diffMins} mins ago`;
  else requestedAt = createdDate.toLocaleString();

  let reviewedAt: string | undefined;
  if (bt.reviewed_at) {
    const rd = new Date(bt.reviewed_at);
    const rDiffMs = now.getTime() - rd.getTime();
    const rDiffMins = Math.floor(rDiffMs / 60000);
    if (rDiffMins < 1) reviewedAt = 'Just now';
    else if (rDiffMins < 60) reviewedAt = `${rDiffMins} mins ago`;
    else reviewedAt = rd.toLocaleString();
  }

  return {
    id: bt.id,
    userId: requester?.id || '',
    userName: requester?.username || 'Unknown',
    userEmail: requester?.email || '',
    userAvatar: requester
      ? `https://ui-avatars.com/api/?name=${encodeURIComponent(requester.username)}&background=random`
      : undefined,
    resourceId: resource?.id || '',
    resourceUrn: resource?.urn || '',
    resourceName: resource?.name || '',
    reason: bt.reason,
    durationMinutes: bt.duration_minutes,
    durationHours: Math.ceil(bt.duration_minutes / 60),
    status: bt.status as AccessTask['status'],
    requestedAt,
    reviewedAt,
    issuedToken,
  };
}

export const useApprovalStore = create<ApprovalState>((set) => ({
  tasks: [],
  isLoading: false,
  error: null,

  fetchTasks: async () => {
    set({ isLoading: true, error: null });
    try {
      const backendTasks = await tasksApi.list();
      const mapped = backendTasks.map(mapBackendTask);
      set({ tasks: mapped, isLoading: false });
    } catch (err: any) {
      set({ error: err?.message || 'Failed to fetch tasks', isLoading: false });
    }
  },

  createTaskRequest: async (
    _userId,
    _userName,
    _userEmail,
    resourceId,
    _resourceUrn,
    _resourceName,
    reason,
    durationHours
  ) => {
    set({ isLoading: true, error: null });
    try {
      const created = await tasksApi.create({
        resource_id: resourceId,
        reason,
        duration_minutes: durationHours * 60,
      });
      const mapped = mapBackendTask(created);
      set((state) => ({
        tasks: [mapped, ...state.tasks],
        isLoading: false,
      }));
    } catch (err: any) {
      set({ error: err?.message || 'Failed to create task', isLoading: false });
      throw err;
    }
  },

  approveTask: async (taskId) => {
    try {
      const updated = await tasksApi.approve(taskId);
      const mapped = mapBackendTask(updated);

      // If approved and has an issued token, push it into the resource store
      if (mapped.issuedToken) {
        useResourceStore.getState().addToken(mapped.issuedToken);
      }

      set((state) => ({
        tasks: state.tasks.map((t) =>
          t.id === taskId ? mapped : t
        ),
      }));
    } catch (err: any) {
      set({ error: err?.message || 'Failed to approve task' });
    }
  },

  rejectTask: async (taskId) => {
    try {
      const updated = await tasksApi.reject(taskId);
      const mapped = mapBackendTask(updated);
      set((state) => ({
        tasks: state.tasks.map((t) =>
          t.id === taskId ? mapped : t
        ),
      }));
    } catch (err: any) {
      set({ error: err?.message || 'Failed to reject task' });
    }
  },
}));
