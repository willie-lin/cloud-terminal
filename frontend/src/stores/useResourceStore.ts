import { create } from 'zustand';
import type { Resource, AccessPolicyToken, TerminalSession, BackendResource } from '../types';
import { resourceApi } from '../api/resources';

interface ResourceState {
  resources: Resource[];
  activeTokens: AccessPolicyToken[];
  activeSession: TerminalSession | null;
  isLoading: boolean;
  error: string | null;

  /** Fetch resources from the backend API. */
  fetchResources: () => Promise<void>;
  /** Create a new resource via API. */
  createResource: (data: {
    name: string;
    urn: string;
    ip: string;
    port?: number;
    type?: string;
    env?: string;
    region?: string;
    description?: string;
  }) => Promise<void>;
  /** Delete a resource via API. */
  deleteResource: (id: string) => Promise<void>;

  /** Add a locally-issued STS token (from approved tasks). */
  addToken: (token: AccessPolicyToken) => void;
  /** Start a terminal session. */
  startSession: (token: AccessPolicyToken) => void;
  /** End the active terminal session. */
  endSession: () => void;
}

/** Map a BackendResource from the API to the frontend Resource type. */
function mapBackendResource(br: BackendResource): Resource {
  return {
    id: br.id,
    urn: br.urn,
    name: br.name,
    type: (br.type as Resource['type']) || 'ssh',
    ip: br.ip,
    port: br.port || 22,
    env: (br.env as Resource['env']) || 'dev',
    region: br.region || 'default',
    description: br.description || undefined,
    status: (br.status as Resource['status']) || 'active',
    details: br.details || undefined,
    // For frontend display compatibility:
    targetHost: br.ip,
    targetPort: br.port || 22,
    authUsername: br.details?.auth_username || 'root',
    os: br.details?.os || `Linux (${br.env})`,
  };
}

export const useResourceStore = create<ResourceState>((set) => ({
  resources: [],
  activeTokens: [],
  activeSession: null,
  isLoading: false,
  error: null,

  fetchResources: async () => {
    set({ isLoading: true, error: null });
    try {
      const backendResources = await resourceApi.list();
      const mapped = backendResources.map(mapBackendResource);
      set({ resources: mapped, isLoading: false });
    } catch (err: any) {
      set({ error: err?.message || 'Failed to fetch resources', isLoading: false });
    }
  },

  createResource: async (data) => {
    set({ isLoading: true, error: null });
    try {
      const created = await resourceApi.create({
        urn: data.urn,
        name: data.name,
        ip: data.ip,
        port: data.port || 22,
        type: data.type || 'ssh',
        env: data.env || 'dev',
        region: data.region || 'default',
        description: data.description || '',
        status: 'active',
      });
      set((state) => ({
        resources: [...state.resources, mapBackendResource(created)],
        isLoading: false,
      }));
    } catch (err: any) {
      set({ error: err?.message || 'Failed to create resource', isLoading: false });
      throw err;
    }
  },

  deleteResource: async (id) => {
    try {
      await resourceApi.delete(id);
      set((state) => ({
        resources: state.resources.filter((r) => r.id !== id),
      }));
    } catch (err: any) {
      set({ error: err?.message || 'Failed to delete resource' });
    }
  },

  addToken: (token) =>
    set((state) => ({
      activeTokens: [token, ...state.activeTokens],
    })),

  startSession: (token) =>
    set({
      activeSession: {
        sessionId: `sess-${Date.now()}`,
        resourceUrn: token.resourceUrn,
        resourceName: token.resourceName,
        token: token.token,
        connectedAt: Math.floor(Date.now() / 1000),
        expiresAt: token.expiresAt,
      },
    }),

  endSession: () => set({ activeSession: null }),
}));
