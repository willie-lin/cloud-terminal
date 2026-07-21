import { useState, useEffect } from 'react';
import { useLanguageStore } from '../stores/useLanguageStore';
import api from '../api/client';
import { Plus, Building2, Users, Mail, KeyRound } from 'lucide-react';

interface Tenant {
  id: string;
  name: string;
  description: string;
  created_at: string;
}

export const TenantManagementPage: React.FC = () => {
  const { t } = useLanguageStore();
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [isAdding, setIsAdding] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [error, setError] = useState('');

  // Create admin user modal
  const [adminModal, setAdminModal] = useState<{ tenantId: string; tenantName: string } | null>(null);
  const [adminEmail, setAdminEmail] = useState('');
  const [adminPassword, setAdminPassword] = useState('');
  const [adminUsername, setAdminUsername] = useState('');

  const fetchTenants = async () => {
    try {
      const data = await api.get<Tenant[]>('/admin/tenants');
      setTenants(data);
    } catch (err: any) {
      setError(err?.message || 'Failed to fetch tenants');
    }
  };

  useEffect(() => { fetchTenants(); }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      await api.post('/admin/tenants', { name, description: description || `${name} tenant` });
      setName('');
      setDescription('');
      setIsAdding(false);
      await fetchTenants();
    } catch (err: any) {
      setError(err?.message || 'Failed to create tenant');
    }
  };

  const handleCreateAdmin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!adminModal) return;
    setError('');
    try {
      await api.post(`/admin/tenants/${adminModal.tenantId}/admin-user`, {
        email: adminEmail,
        password: adminPassword,
        username: adminUsername || undefined,
      });
      setAdminModal(null);
      setAdminEmail('');
      setAdminPassword('');
      setAdminUsername('');
    } catch (err: any) {
      setError(err?.message || 'Failed to create admin user');
    }
  };

  return (
    <div className="space-y-6 py-2 text-left">
      <div className="flex items-center justify-between border-b border-zinc-200 dark:border-white/[0.06] pb-4">
        <div>
          <h1 className="text-xl font-bold text-zinc-900 dark:text-white tracking-tight">租户管理</h1>
          <p className="text-xs text-zinc-500 mt-0.5">超管专属 — 创建和管理租户</p>
        </div>
        <button onClick={() => setIsAdding(!isAdding)}
          className="inline-flex items-center space-x-1.5 px-3.5 py-1.5 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs hover:bg-zinc-800 dark:hover:bg-zinc-200 transition-all">
          <Plus className="w-3.5 h-3.5" />
          <span>{isAdding ? '取消' : '创建租户'}</span>
        </button>
      </div>

      {error && <div className="p-3 rounded-xl bg-red-500/10 border border-red-500/30 text-red-600 dark:text-red-400 text-xs">{error}</div>}

      {isAdding && (
        <form onSubmit={handleCreate} className="p-6 rounded-xl bg-white dark:bg-[#121318] border border-zinc-300 dark:border-white/10 space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">租户名称</label>
              <input type="text" value={name} onChange={(e) => setName(e.target.value)} placeholder="e.g. acme-corp" required
                className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
            </div>
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">描述</label>
              <input type="text" value={description} onChange={(e) => setDescription(e.target.value)} placeholder="可选"
                className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
            </div>
          </div>
          <div className="flex justify-end space-x-3">
            <button type="button" onClick={() => setIsAdding(false)} className="px-3.5 py-1.5 rounded-lg text-xs font-medium text-zinc-500 hover:text-zinc-900">取消</button>
            <button type="submit" className="px-4 py-1.5 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs">创建租户</button>
          </div>
        </form>
      )}

      <div className="rounded-xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/[0.07] overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-zinc-200 dark:border-white/[0.06] bg-zinc-50 dark:bg-black/40 text-[10px] uppercase font-mono text-zinc-500">
              <th className="py-3 px-5">名称</th>
              <th className="py-3 px-5">描述</th>
              <th className="py-3 px-5">创建时间</th>
              <th className="py-3 px-5 text-right">操作</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-zinc-200 dark:divide-white/[0.04] text-xs">
            {tenants.filter(t => t.name !== 'systemTenant').map((tenant) => (
              <tr key={tenant.id} className="hover:bg-zinc-50 dark:hover:bg-white/[0.02]">
                <td className="py-3 px-5">
                  <div className="flex items-center space-x-2">
                    <Building2 className="w-3.5 h-3.5 text-zinc-400" />
                    <span className="font-semibold text-zinc-900 dark:text-white">{tenant.name}</span>
                  </div>
                </td>
                <td className="py-3 px-5 text-zinc-500">{tenant.description || '-'}</td>
                <td className="py-3 px-5 text-zinc-500 font-mono">{new Date(tenant.created_at).toLocaleDateString()}</td>
                <td className="py-3 px-5 text-right">
                  <button onClick={() => setAdminModal({ tenantId: tenant.id, tenantName: tenant.name })}
                    className="inline-flex items-center space-x-1 px-2.5 py-1 rounded bg-zinc-100 dark:bg-white/5 hover:bg-zinc-200 dark:hover:bg-white/10 text-zinc-700 dark:text-zinc-300 text-xs transition-colors">
                    <Users className="w-3 h-3" />
                    <span>创建管理员</span>
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Create Admin Modal */}
      {adminModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
          <div className="relative w-full max-w-md rounded-xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/10 p-6 shadow-2xl text-left">
            <h2 className="text-base font-semibold text-zinc-900 dark:text-white mb-1">创建租户管理员</h2>
            <p className="text-xs text-zinc-500 mb-4">为租户 <strong>{adminModal.tenantName}</strong> 创建管理员账号</p>
            <form onSubmit={handleCreateAdmin} className="space-y-4">
              <div className="space-y-1.5">
                <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">邮箱</label>
                <div className="relative">
                  <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-400" />
                  <input type="email" value={adminEmail} onChange={(e) => setAdminEmail(e.target.value)} placeholder="admin@company.com" required
                    className="w-full pl-10 pr-4 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
                </div>
              </div>
              <div className="space-y-1.5">
                <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">密码</label>
                <div className="relative">
                  <KeyRound className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-400" />
                  <input type="password" value={adminPassword} onChange={(e) => setAdminPassword(e.target.value)} placeholder="••••••••" required
                    className="w-full pl-10 pr-4 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
                </div>
              </div>
              <div className="space-y-1.5">
                <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">用户名（可选）</label>
                <input type="text" value={adminUsername} onChange={(e) => setAdminUsername(e.target.value)} placeholder="留空自动生成"
                  className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
              </div>
              <div className="flex justify-end space-x-3 pt-2">
                <button type="button" onClick={() => setAdminModal(null)}
                  className="px-3.5 py-1.5 rounded-lg text-xs font-medium text-zinc-500 hover:text-zinc-900">取消</button>
                <button type="submit"
                  className="px-4 py-1.5 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs">创建管理员</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
