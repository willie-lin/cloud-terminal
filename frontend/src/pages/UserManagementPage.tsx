import { useState, useEffect } from 'react';
import { useLanguageStore } from '../stores/useLanguageStore';
import { useAuthStore } from '../stores/useAuthStore';
import api from '../api/client';
import { Plus, Mail, KeyRound, User as UserIcon } from 'lucide-react';

interface ManagedUser {
  id: string;
  email: string;
  username: string;
  created_at: string;
}

export const UserManagementPage: React.FC = () => {
  const { t } = useLanguageStore();
  const { user } = useAuthStore();
  const [users, setUsers] = useState<ManagedUser[]>([]);
  const [isAdding, setIsAdding] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [username, setUsername] = useState('');
  const [error, setError] = useState('');

  const fetchUsers = async () => {
    try {
      const data = await api.get<ManagedUser[]>('/admin/users');
      setUsers(data);
    } catch (err: any) {
      setError(err?.message || 'Failed to fetch users');
    }
  };

  useEffect(() => { fetchUsers(); }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      await api.post('/admin/users', { email, password, username: username || undefined });
      setEmail('');
      setPassword('');
      setUsername('');
      setIsAdding(false);
      await fetchUsers();
    } catch (err: any) {
      setError(err?.message || 'Failed to create user');
    }
  };

  const isSuperAdmin = user?.roleName === 'super_admin';

  return (
    <div className="space-y-6 py-2 text-left">
      <div className="flex items-center justify-between border-b border-zinc-200 dark:border-white/[0.06] pb-4">
        <div>
          <h1 className="text-xl font-bold text-zinc-900 dark:text-white tracking-tight">用户管理</h1>
          <p className="text-xs text-zinc-500 mt-0.5">{isSuperAdmin ? '超管 — 管理所有用户' : '租户管理员 — 在本租户下创建用户'}</p>
        </div>
        <button onClick={() => setIsAdding(!isAdding)}
          className="inline-flex items-center space-x-1.5 px-3.5 py-1.5 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs hover:bg-zinc-800 dark:hover:bg-zinc-200 transition-all">
          <Plus className="w-3.5 h-3.5" />
          <span>{isAdding ? '取消' : '创建用户'}</span>
        </button>
      </div>

      {error && <div className="p-3 rounded-xl bg-red-500/10 border border-red-500/30 text-red-600 dark:text-red-400 text-xs">{error}</div>}

      {isAdding && (
        <form onSubmit={handleCreate} className="p-6 rounded-xl bg-white dark:bg-[#121318] border border-zinc-300 dark:border-white/10 space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">邮箱</label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-400" />
                <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="user@company.com" required
                  className="w-full pl-10 pr-4 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
              </div>
            </div>
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">密码</label>
              <div className="relative">
                <KeyRound className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-400" />
                <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="••••••••" required
                  className="w-full pl-10 pr-4 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
              </div>
            </div>
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">用户名（可选）</label>
              <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} placeholder="留空自动生成"
                className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400" />
            </div>
          </div>
          <div className="flex justify-end space-x-3 pt-2">
            <button type="button" onClick={() => setIsAdding(false)} className="px-3.5 py-1.5 rounded-lg text-xs font-medium text-zinc-500 hover:text-zinc-900">取消</button>
            <button type="submit" className="px-4 py-1.5 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs">创建用户</button>
          </div>
        </form>
      )}

      <div className="rounded-xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/[0.07] overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-zinc-200 dark:border-white/[0.06] bg-zinc-50 dark:bg-black/40 text-[10px] uppercase font-mono text-zinc-500">
              <th className="py-3 px-5">用户名</th>
              <th className="py-3 px-5">邮箱</th>
              <th className="py-3 px-5">用户 ID</th>
              <th className="py-3 px-5">创建时间</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-zinc-200 dark:divide-white/[0.04] text-xs">
            {users.map((u) => (
              <tr key={u.id} className="hover:bg-zinc-50 dark:hover:bg-white/[0.02]">
                <td className="py-3 px-5">
                  <div className="flex items-center space-x-2">
                    <UserIcon className="w-3.5 h-3.5 text-zinc-400" />
                    <span className="font-semibold text-zinc-900 dark:text-white">{u.username}</span>
                  </div>
                </td>
                <td className="py-3 px-5 text-zinc-600 dark:text-zinc-400">{u.email}</td>
                <td className="py-3 px-5 text-zinc-500 font-mono text-[10px]">{u.id.substring(0, 12)}...</td>
                <td className="py-3 px-5 text-zinc-500 font-mono">{new Date(u.created_at).toLocaleDateString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
