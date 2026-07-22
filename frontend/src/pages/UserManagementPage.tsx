import { useState, useEffect } from 'react';
import { useAuthStore } from '../stores/useAuthStore';
import api from '../api/client';
import { Plus, Mail, KeyRound, User as UserIcon, Shield, Pencil, Trash2, X } from 'lucide-react';

interface ManagedUser {
  id: string;
  email: string;
  username: string;
  nickname?: string;
  bio?: string;
  phone_number?: string;
  created_at: string;
  edges?: {
    roles?: Array<{ id: string; name: string; description?: string }>;
  };
}

export const UserManagementPage: React.FC = () => {
  const { user } = useAuthStore();
  const [users, setUsers] = useState<ManagedUser[]>([]);
  const [isAdding, setIsAdding] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [username, setUsername] = useState('');
  const [roleName, setRoleName] = useState('user');
  const [error, setError] = useState('');
  const [successMsg, setSuccessMsg] = useState('');

  // Edit modal state
  const [editModal, setEditModal] = useState<ManagedUser | null>(null);
  const [editNickname, setEditNickname] = useState('');
  const [editBio, setEditBio] = useState('');
  const [editPhone, setEditPhone] = useState('');
  const [editOnline, setEditOnline] = useState(true);
  const [editStatus, setEditStatus] = useState(true);

  // Delete modal state
  const [deleteModal, setDeleteModal] = useState<ManagedUser | null>(null);

  const fetchUsers = async () => {
    try {
      const data = await api.get<ManagedUser[]>('/admin/users');
      setUsers(data);
    } catch (err: any) {
      setError(err?.message || 'Failed to fetch users');
    }
  };

  useEffect(() => { fetchUsers(); }, []);

  const clearMessages = () => { setError(''); setSuccessMsg(''); };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    clearMessages();
    try {
      await api.post('/admin/users', {
        email: email.trim(),
        password,
        username: username.trim(),
        role_name: roleName,
      });
      setEmail(''); setPassword(''); setUsername(''); setRoleName('user');
      setIsAdding(false);
      setSuccessMsg(`用户 ${username} 创建成功，已赋予 ${roleName} 角色！`);
      await fetchUsers();
    } catch (err: any) {
      setError(err?.message || 'Failed to create user');
    }
  };

  const openEditModal = (u: ManagedUser) => {
    setEditModal(u);
    setEditNickname(u.nickname || '');
    setEditBio(u.bio || '');
    setEditPhone(u.phone_number || '');
    setEditOnline(true);
    setEditStatus(true);
    clearMessages();
  };

  const handleUpdateUser = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editModal) return;
    clearMessages();
    try {
      await api.put(`/admin/users/${editModal.id}`, {
        nickname: editNickname || undefined,
        bio: editBio || undefined,
        phone_number: editPhone || undefined,
        online: editOnline,
        status: editStatus,
      });
      setEditModal(null);
      setSuccessMsg(`用户 ${editModal.username} 信息更新成功！`);
      await fetchUsers();
    } catch (err: any) {
      setError(err?.message || 'Failed to update user');
    }
  };

  const handleDeleteUser = async () => {
    if (!deleteModal) return;
    clearMessages();
    try {
      await api.delete(`/admin/users/${deleteModal.id}`);
      const deletedName = deleteModal.username;
      setDeleteModal(null);
      setSuccessMsg(`用户 ${deletedName} 已成功删除！`);
      await fetchUsers();
    } catch (err: any) {
      setError(err?.message || 'Failed to delete user');
    }
  };

  const isSuperAdmin = user?.roleName === 'super_admin';

  return (
    <div className="space-y-6 py-2 text-left">
      <div className="flex items-center justify-between border-b border-zinc-200 dark:border-white/[0.06] pb-4">
        <div>
          <h1 className="text-xl font-bold text-zinc-900 dark:text-white tracking-tight">用户管理</h1>
          <p className="text-xs text-zinc-500 mt-0.5">
            {isSuperAdmin ? '超管 — 管理全平台用户与角色 (CRUD)' : '租户管理员 — 在本租户下创建成员与分配角色'}
          </p>
        </div>
        <button onClick={() => setIsAdding(!isAdding)}
          className="inline-flex items-center space-x-1.5 px-3.5 py-1.5 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs hover:bg-zinc-800 dark:hover:bg-zinc-200 transition-all">
          <Plus className="w-3.5 h-3.5" />
          <span>{isAdding ? '取消' : '创建用户'}</span>
        </button>
      </div>

      {error && <div className="p-3 rounded-xl bg-red-500/10 border border-red-500/30 text-red-600 dark:text-red-400 text-xs">{error}</div>}
      {successMsg && <div className="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-600 dark:text-emerald-400 text-xs">{successMsg}</div>}

      {isAdding && (
        <form onSubmit={handleCreate} className="p-6 rounded-xl bg-white dark:bg-[#121318] border border-zinc-300 dark:border-white/10 space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">邮箱 *</label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-400" />
                <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="user@company.com" required
                  className="w-full pl-10 pr-4 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
              </div>
            </div>
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">密码 *</label>
              <div className="relative">
                <KeyRound className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-400" />
                <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="••••••••" required minLength={6}
                  className="w-full pl-10 pr-4 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
              </div>
            </div>
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">用户名 *（至少6字符）</label>
              <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} placeholder="例如：user_john" required minLength={6}
                className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
            </div>
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">分配角色 *</label>
              <select value={roleName} onChange={(e) => setRoleName(e.target.value)}
                className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400">
                <option value="user">普通用户 (user)</option>
                <option value="tenant_admin">租户管理员 (tenant_admin)</option>
              </select>
            </div>
          </div>
          <div className="flex justify-end space-x-3 pt-2">
            <button type="button" onClick={() => setIsAdding(false)} className="px-3.5 py-1.5 rounded-lg text-xs font-medium text-zinc-500 hover:text-zinc-900">取消</button>
            <button type="submit" className="px-4 py-1.5 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs shadow">创建用户</button>
          </div>
        </form>
      )}

      <div className="rounded-xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/[0.07] overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-zinc-200 dark:border-white/[0.06] bg-zinc-50 dark:bg-black/40 text-[10px] uppercase font-mono text-zinc-500">
              <th className="py-3 px-5">用户名</th>
              <th className="py-3 px-5">邮箱</th>
              <th className="py-3 px-5">绑定角色</th>
              <th className="py-3 px-5">用户 ID</th>
              <th className="py-3 px-5">创建时间</th>
              <th className="py-3 px-5 text-right">操作 (CRUD)</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-zinc-200 dark:divide-white/[0.04] text-xs">
            {users.map((u) => {
              const rolesList = u.edges?.roles || [];
              return (
                <tr key={u.id} className="hover:bg-zinc-50 dark:hover:bg-white/[0.02]">
                  <td className="py-3 px-5">
                    <div className="flex items-center space-x-2">
                      <UserIcon className="w-3.5 h-3.5 text-zinc-400" />
                      <span className="font-semibold text-zinc-900 dark:text-white">{u.username}</span>
                    </div>
                  </td>
                  <td className="py-3 px-5 text-zinc-600 dark:text-zinc-400">{u.email}</td>
                  <td className="py-3 px-5">
                    <div className="flex items-center space-x-1.5 flex-wrap gap-y-1">
                      {rolesList.length > 0 ? (
                        rolesList.map(r => (
                          <span key={r.id} className="inline-flex items-center space-x-1 px-2 py-0.5 rounded text-[10px] font-mono font-semibold bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20">
                            <Shield className="w-2.5 h-2.5" />
                            <span>{r.name}</span>
                          </span>
                        ))
                      ) : (
                        <span className="inline-flex items-center space-x-1 px-2 py-0.5 rounded text-[10px] font-mono bg-zinc-100 dark:bg-white/5 text-zinc-500">
                          <span>user</span>
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="py-3 px-5 text-zinc-500 font-mono text-[10px]">{u.id.substring(0, 12)}...</td>
                  <td className="py-3 px-5 text-zinc-500 font-mono">{new Date(u.created_at).toLocaleDateString()}</td>
                  <td className="py-3 px-5 text-right space-x-2">
                    <button
                      onClick={() => openEditModal(u)}
                      className="inline-flex items-center space-x-1 px-2.5 py-1 rounded bg-zinc-100 dark:bg-white/5 hover:bg-zinc-200 dark:hover:bg-white/10 text-zinc-700 dark:text-zinc-300 text-xs transition-colors"
                    >
                      <Pencil className="w-3 h-3" />
                      <span>编辑</span>
                    </button>
                    <button
                      onClick={() => { setDeleteModal(u); clearMessages(); }}
                      className="inline-flex items-center space-x-1 px-2.5 py-1 rounded bg-red-500/10 hover:bg-red-500/20 text-red-600 dark:text-red-400 text-xs transition-colors"
                    >
                      <Trash2 className="w-3 h-3" />
                      <span>删除</span>
                    </button>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {/* Edit User Modal */}
      {editModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
          <div className="relative w-full max-w-md rounded-xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/10 p-6 shadow-2xl text-left">
            <button onClick={() => setEditModal(null)} className="absolute top-4 right-4 text-zinc-400 hover:text-zinc-700 dark:hover:text-white">
              <X className="w-4 h-4" />
            </button>
            <h2 className="text-base font-semibold text-zinc-900 dark:text-white mb-1">编辑用户信息</h2>
            <p className="text-xs text-zinc-500 mb-4">修改用户 <strong>{editModal.username}</strong> 的资料</p>
            <form onSubmit={handleUpdateUser} className="space-y-4">
              <div className="grid grid-cols-1 gap-4">
                <div className="space-y-1.5">
                  <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">昵称</label>
                  <input type="text" value={editNickname} onChange={(e) => setEditNickname(e.target.value)} placeholder="显示名称"
                    className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
                </div>
                <div className="space-y-1.5">
                  <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">手机号</label>
                  <input type="tel" value={editPhone} onChange={(e) => setEditPhone(e.target.value)} placeholder="+86 138 0000 0000"
                    className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
                </div>
                <div className="space-y-1.5">
                  <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">个人简介</label>
                  <textarea value={editBio} onChange={(e) => setEditBio(e.target.value)} placeholder="介绍一下..." rows={2}
                    className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400 resize-none" />
                </div>
                <div className="flex items-center space-x-6">
                  <label className="flex items-center space-x-2 cursor-pointer">
                    <input type="checkbox" checked={editOnline} onChange={(e) => setEditOnline(e.target.checked)}
                      className="rounded border-zinc-300 dark:border-white/20" />
                    <span className="text-xs text-zinc-600 dark:text-zinc-400">在线状态</span>
                  </label>
                  <label className="flex items-center space-x-2 cursor-pointer">
                    <input type="checkbox" checked={editStatus} onChange={(e) => setEditStatus(e.target.checked)}
                      className="rounded border-zinc-300 dark:border-white/20" />
                    <span className="text-xs text-zinc-600 dark:text-zinc-400">账号启用</span>
                  </label>
                </div>
              </div>
              <div className="flex justify-end space-x-3 pt-2 border-t border-zinc-200 dark:border-white/10">
                <button type="button" onClick={() => setEditModal(null)}
                  className="px-3.5 py-1.5 rounded-lg text-xs font-medium text-zinc-500 hover:text-zinc-900">取消</button>
                <button type="submit"
                  className="px-4 py-1.5 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs shadow">保存修改</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Delete User Confirmation Modal */}
      {deleteModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
          <div className="relative w-full max-w-md rounded-xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/10 p-6 shadow-2xl text-left">
            <h2 className="text-base font-bold text-red-600 dark:text-red-400 mb-1">确认删除用户？</h2>
            <p className="text-xs text-zinc-500 mb-4">确定要删除用户 <strong>{deleteModal.username}</strong>（{deleteModal.email}）吗？此操作不可撤销！</p>
            <div className="flex justify-end space-x-3 pt-2">
              <button type="button" onClick={() => setDeleteModal(null)}
                className="px-3.5 py-1.5 rounded-lg text-xs font-medium text-zinc-500 hover:text-zinc-900">取消</button>
              <button type="button" onClick={handleDeleteUser}
                className="px-4 py-1.5 rounded-lg bg-red-600 text-white font-semibold text-xs shadow hover:bg-red-700">确认彻底删除</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
