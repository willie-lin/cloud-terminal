import { useState, useEffect } from 'react';
import api from '../api/client';
import { Plus, Building2, Pencil, Trash2, Mail, KeyRound } from 'lucide-react';

interface Tenant {
  id: string;
  name: string;
  description: string;
  created_at: string;
}

export const TenantManagementPage: React.FC = () => {
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [isAdding, setIsAdding] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [initAdminEmail, setInitAdminEmail] = useState('');
  const [initAdminPassword, setInitAdminPassword] = useState('');
  const [initAdminUsername, setInitAdminUsername] = useState('');
  const [error, setError] = useState('');
  const [successMsg, setSuccessMsg] = useState('');

  // Edit tenant modal
  const [editModal, setEditModal] = useState<Tenant | null>(null);
  const [editName, setEditName] = useState('');
  const [editDesc, setEditDesc] = useState('');

  // Delete tenant modal
  const [deleteModal, setDeleteModal] = useState<Tenant | null>(null);

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
    setSuccessMsg('');
    try {
      const payload: any = {
        name: name.trim(),
        description: description || `${name} 企业租户空间`,
      };
      if (initAdminEmail && initAdminPassword) {
        payload.admin_email = initAdminEmail.trim();
        payload.admin_password = initAdminPassword;
        payload.admin_username = initAdminUsername.trim() || undefined;
      }

      const res: any = await api.post('/admin/tenants', payload);
      
      setName('');
      setDescription('');
      setInitAdminEmail('');
      setInitAdminPassword('');
      setInitAdminUsername('');
      setIsAdding(false);

      if (res?.admin_created) {
        setSuccessMsg(`租户 ${payload.name} 创建成功，且已自动部署租户管理员账号 (${payload.admin_email})！`);
      } else {
        setSuccessMsg(`租户 ${payload.name} 创建成功！`);
      }

      await fetchTenants();
    } catch (err: any) {
      setError(err?.message || 'Failed to create tenant');
    }
  };

  const handleUpdateTenant = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editModal) return;
    setError('');
    setSuccessMsg('');
    try {
      await api.put(`/admin/tenants/${editModal.id}`, {
        name: editName.trim(),
        description: editDesc.trim(),
      });
      setEditModal(null);
      setSuccessMsg(`租户 ${editName} 更新成功！`);
      await fetchTenants();
    } catch (err: any) {
      setError(err?.message || 'Failed to update tenant');
    }
  };

  const handleDeleteTenant = async () => {
    if (!deleteModal) return;
    setError('');
    setSuccessMsg('');
    try {
      await api.delete(`/admin/tenants/${deleteModal.id}`);
      const deletedName = deleteModal.name;
      setDeleteModal(null);
      setSuccessMsg(`租户 ${deletedName} 已成功删除！`);
      await fetchTenants();
    } catch (err: any) {
      setError(err?.message || 'Failed to delete tenant');
    }
  };

  return (
    <div className="space-y-6 py-2 text-left">
      <div className="flex items-center justify-between border-b border-zinc-200 dark:border-white/[0.06] pb-4">
        <div>
          <h1 className="text-xl font-bold text-zinc-900 dark:text-white tracking-tight">租户管理</h1>
          <p className="text-xs text-zinc-500 mt-0.5">超管专属 — 全生命周期开户与控制 (CRUD)</p>
        </div>
        <button onClick={() => setIsAdding(!isAdding)}
          className="inline-flex items-center space-x-1.5 px-3.5 py-1.5 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs hover:bg-zinc-800 dark:hover:bg-zinc-200 transition-all">
          <Plus className="w-3.5 h-3.5" />
          <span>{isAdding ? '取消' : '创建租户'}</span>
        </button>
      </div>

      {error && <div className="p-3 rounded-xl bg-red-500/10 border border-red-500/30 text-red-600 dark:text-red-400 text-xs">{error}</div>}
      {successMsg && <div className="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-600 dark:text-emerald-400 text-xs">{successMsg}</div>}

      {isAdding && (
        <form onSubmit={handleCreate} className="p-6 rounded-xl bg-white dark:bg-[#121318] border border-zinc-300 dark:border-white/10 space-y-5 text-left">
          <div className="border-b border-zinc-200 dark:border-white/10 pb-3">
            <h3 className="text-sm font-bold text-zinc-900 dark:text-white">1. 租户基本信息</h3>
            <p className="text-[11px] text-zinc-500">分配独立租户空间与资源隔离防线</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">租户标识 / 名称 *</label>
              <input type="text" value={name} onChange={(e) => setName(e.target.value)} placeholder="例如：acme-corp" required
                className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
            </div>
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">租户描述</label>
              <input type="text" value={description} onChange={(e) => setDescription(e.target.value)} placeholder="例如：Acme 企业版专属空间"
                className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
            </div>
          </div>

          <div className="border-b border-zinc-200 dark:border-white/10 pt-2 pb-3">
            <h3 className="text-sm font-bold text-zinc-900 dark:text-white">2. 同步初始化租户管理员 (必填)</h3>
            <p className="text-[11px] text-zinc-500">原子化生成该租户的首个 tenant_admin 管理账号</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">管理员邮箱 *</label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-zinc-400" />
                <input type="email" value={initAdminEmail} onChange={(e) => setInitAdminEmail(e.target.value)} placeholder="admin@acme.com" required
                  className="w-full pl-9 pr-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
              </div>
            </div>
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">管理员初始密码 *</label>
              <div className="relative">
                <KeyRound className="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-zinc-400" />
                <input type="password" value={initAdminPassword} onChange={(e) => setInitAdminPassword(e.target.value)} placeholder="••••••••" required minLength={6}
                  className="w-full pl-9 pr-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
              </div>
            </div>
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">管理员用户名 *（至少6字符）</label>
              <input type="text" value={initAdminUsername} onChange={(e) => setInitAdminUsername(e.target.value)} placeholder="例如：acme_admin" required minLength={6}
                className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
            </div>
          </div>

          <div className="flex justify-end space-x-3 pt-2 border-t border-zinc-200 dark:border-white/10">
            <button type="button" onClick={() => setIsAdding(false)} className="px-3.5 py-1.5 rounded-lg text-xs font-medium text-zinc-500 hover:text-zinc-900">取消</button>
            <button type="submit" className="px-4 py-1.5 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs shadow">
              一键建租户并初始化管理员
            </button>
          </div>
        </form>
      )}

      <div className="rounded-xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/[0.07] overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-zinc-200 dark:border-white/[0.06] bg-zinc-50 dark:bg-black/40 text-[10px] uppercase font-mono text-zinc-500">
              <th className="py-3 px-5">租户名称</th>
              <th className="py-3 px-5">描述</th>
              <th className="py-3 px-5">创建时间</th>
              <th className="py-3 px-5 text-right">操作 (CRUD)</th>
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
                <td className="py-3 px-5 text-right space-x-2">
                  <button
                    onClick={() => {
                      setEditModal(tenant);
                      setEditName(tenant.name);
                      setEditDesc(tenant.description || '');
                    }}
                    className="inline-flex items-center space-x-1 px-2.5 py-1 rounded bg-zinc-100 dark:bg-white/5 hover:bg-zinc-200 dark:hover:bg-white/10 text-zinc-700 dark:text-zinc-300 text-xs transition-colors"
                  >
                    <Pencil className="w-3 h-3" />
                    <span>编辑</span>
                  </button>
                  <button
                    onClick={() => setDeleteModal(tenant)}
                    className="inline-flex items-center space-x-1 px-2.5 py-1 rounded bg-red-500/10 hover:bg-red-500/20 text-red-600 dark:text-red-400 text-xs transition-colors"
                  >
                    <Trash2 className="w-3 h-3" />
                    <span>删除</span>
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Edit Tenant Modal */}
      {editModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
          <div className="relative w-full max-w-md rounded-xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/10 p-6 shadow-2xl text-left">
            <h2 className="text-base font-semibold text-zinc-900 dark:text-white mb-1">编辑租户信息</h2>
            <p className="text-xs text-zinc-500 mb-4">修改租户 <strong>{editModal.name}</strong> 的配置</p>
            <form onSubmit={handleUpdateTenant} className="space-y-4">
              <div className="space-y-1.5">
                <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">租户名称</label>
                <input type="text" value={editName} onChange={(e) => setEditName(e.target.value)} required
                  className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
              </div>
              <div className="space-y-1.5">
                <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">描述</label>
                <input type="text" value={editDesc} onChange={(e) => setEditDesc(e.target.value)} placeholder="例如：Acme 企业版空间"
                  className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
              </div>
              <div className="flex justify-end space-x-3 pt-2">
                <button type="button" onClick={() => setEditModal(null)}
                  className="px-3.5 py-1.5 rounded-lg text-xs font-medium text-zinc-500 hover:text-zinc-900">取消</button>
                <button type="submit"
                  className="px-4 py-1.5 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs shadow">保存修改</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Delete Tenant Modal */}
      {deleteModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
          <div className="relative w-full max-w-md rounded-xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/10 p-6 shadow-2xl text-left">
            <h2 className="text-base font-bold text-red-600 dark:text-red-400 mb-1">确认删除租户？</h2>
            <p className="text-xs text-zinc-500 mb-4">确定要删除租户 <strong>{deleteModal.name}</strong> 吗？此操作不可撤销，并将释放相关绑定的组空间！</p>
            <div className="flex justify-end space-x-3 pt-2">
              <button type="button" onClick={() => setDeleteModal(null)}
                className="px-3.5 py-1.5 rounded-lg text-xs font-medium text-zinc-500 hover:text-zinc-900">取消</button>
              <button type="button" onClick={handleDeleteTenant}
                className="px-4 py-1.5 rounded-lg bg-red-600 text-white font-semibold text-xs shadow hover:bg-red-700">确认彻底删除</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
