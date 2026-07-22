import { useState, useEffect } from 'react';
import { Routes, Route, Navigate, useNavigate, useLocation } from 'react-router-dom';
import { useAuthStore } from './stores/useAuthStore';
import { useResourceStore } from './stores/useResourceStore';
import { useApprovalStore } from './stores/useApprovalStore';
import { useLanguageStore } from './stores/useLanguageStore';
import { Navigation } from './components/Navigation';
import { ProtectedRoute } from './components/ProtectedRoute';
import { TaskRequestModal } from './components/TaskRequestModal';
import { UserPortalPage } from './pages/UserPortalPage';
import { AdminApprovalPage } from './pages/AdminApprovalPage';
import { ResourceManagerPage } from './pages/ResourceManagerPage';
import { TenantManagementPage } from './pages/TenantManagementPage';
import { UserManagementPage } from './pages/UserManagementPage';
import { ForbiddenPage } from './pages/ForbiddenPage';
import { TerminalPage } from './pages/TerminalPage';
import type { AccessPolicyToken } from './types';
import { ShieldCheck, ArrowRight, Mail, KeyRound } from 'lucide-react';
import './App.css';

export function App() {
  const { user, isAuthenticated, isLoading, error, login } = useAuthStore();
  const { activeSession, startSession, endSession, fetchResources } = useResourceStore();
  const { fetchTasks } = useApprovalStore();
  const { t } = useLanguageStore();

  const navigate = useNavigate();
  const location = useLocation();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const isSuperAdmin = Boolean(user?.isSuperAdmin || user?.roleName === 'super_admin' || user?.role === 'super_admin');
  const isTenantAdmin = Boolean(user?.isTenantAdmin || user?.roleName === 'tenant_admin' || user?.roleName?.includes('tenant_admin'));

  const getActiveTabFromPath = (path: string): 'access' | 'approvals' | 'resources' | 'tenants' | 'users' => {
    if (path.startsWith('/tenants')) return 'tenants';
    if (path.startsWith('/users')) return 'users';
    if (path.startsWith('/resources')) return 'resources';
    if (path.startsWith('/approvals')) return 'approvals';
    return 'access';
  };

  const activeTab = getActiveTabFromPath(location.pathname);

  useEffect(() => {
    if (isAuthenticated && user) {
      fetchResources();
      fetchTasks();
    }
  }, [isAuthenticated, user, fetchResources, fetchTasks]);

  const handleConnect = (token: AccessPolicyToken) => {
    startSession(token);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) return;
    try {
      await login({ email, password });
      // 登录成功后根据角色进行路由跳转
      const loggedUser = useAuthStore.getState().user;
      const superAdmin = Boolean(loggedUser?.isSuperAdmin || loggedUser?.roleName === 'super_admin' || loggedUser?.role === 'super_admin');
      const tenantAdmin = Boolean(loggedUser?.isTenantAdmin || loggedUser?.roleName === 'tenant_admin' || loggedUser?.roleName?.includes('tenant_admin'));
      
      if (superAdmin) {
        navigate('/tenants');
      } else if (tenantAdmin) {
        navigate('/users');
      } else {
        navigate('/access');
      }
    } catch {
      // error is stored in useAuthStore.error
    }
  };

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-zinc-50 dark:bg-[#09090b] flex items-center justify-center p-6 text-left">
        <div className="w-full max-w-md rounded-2xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/10 p-8 shadow-2xl">
          <div className="flex items-center space-x-3 mb-8">
            <div className="w-10 h-10 rounded-xl bg-zinc-900 text-white dark:bg-white dark:text-black flex items-center justify-center shadow-md">
              <svg className="w-5 h-5 stroke-[2]" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <path d="M17.5 19H9a7 7 0 1 1 6.71-9h1.79a4.5 4.5 0 1 1 0 9Z" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="m8.5 12.5 1.5 1.5-1.5 1.5" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M12.5 15.5h2" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>
            </div>
            <div>
              <h1 className="font-bold text-zinc-900 dark:text-white tracking-tight text-lg">Cloud Terminal</h1>
              <p className="text-xs text-zinc-500">Cloud Native SSH Control Plane</p>
            </div>
          </div>

          {error && (
            <div className="mb-4 p-3 rounded-xl bg-red-500/10 border border-red-500/30 text-red-600 dark:text-red-400 text-xs">{error}</div>
          )}

          <form onSubmit={handleSubmit} className="space-y-5">
            <div className="space-y-2">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase tracking-wider">邮箱 / 用户名</label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-400" />
                <input type="text" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="you@company.com 或 用户名" autoComplete="username"
                  className="w-full pl-10 pr-4 py-2.5 rounded-xl bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/[0.08] text-zinc-900 dark:text-white text-sm focus:outline-none focus:border-zinc-400 transition-all" />
              </div>
            </div>
            <div className="space-y-2">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase tracking-wider">密码</label>
              <div className="relative">
                <KeyRound className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-400" />
                <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="••••••••" autoComplete="current-password"
                  className="w-full pl-10 pr-4 py-2.5 rounded-xl bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/[0.08] text-zinc-900 dark:text-white text-sm focus:outline-none focus:border-zinc-400 transition-all" />
              </div>
            </div>
            <div className="p-3.5 rounded-xl bg-zinc-100 dark:bg-white/[0.03] border border-zinc-200 dark:border-white/[0.06] space-y-2 text-xs">
              <div className="flex items-center justify-between text-zinc-600 dark:text-zinc-400 font-semibold">
                <span>初始超级管理员账号</span>
                <button
                  type="button"
                  onClick={() => {
                    setEmail('superadmin@example.com');
                    setPassword('67727a41b5b1d4dfca981e4045b1bb2f1e7fef0e3e8825c028949d186cad4c00');
                  }}
                  className="text-[11px] text-emerald-600 dark:text-emerald-400 hover:underline font-mono font-bold"
                >
                  一键填入超管凭证
                </button>
              </div>
              <div className="text-[11px] font-mono text-zinc-500 dark:text-zinc-500 break-all leading-tight">
                账号: superadmin@example.com
              </div>
            </div>
            <button type="submit" disabled={!email || !password || isLoading}
              className="w-full py-2.5 px-5 rounded-xl bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-sm hover:bg-zinc-800 dark:hover:bg-zinc-200 disabled:opacity-40 disabled:cursor-not-allowed transition-all flex items-center justify-center space-x-2">
              <span>{isLoading ? '登录中...' : t('enterControlPlane')}</span>
              <ArrowRight className="w-4 h-4 stroke-[2.5]" />
            </button>
          </form>
        </div>
      </div>
    );
  }

  if (activeSession) {
    return (
      <div className="min-h-screen terminal-page-wrapper text-zinc-900 dark:text-zinc-100 flex flex-col font-sans px-4 sm:px-6">
        <TerminalPage session={activeSession} onDisconnect={endSession} />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-[#09090b] text-zinc-900 dark:text-zinc-100 flex flex-col font-sans">
      <Navigation activeTab={activeTab} onSelectTab={(tab) => navigate(`/${tab}`)} onOpenRequestModal={() => setIsModalOpen(true)} />
      <main className="flex-1 max-w-7xl w-full mx-auto px-6 py-6 pb-20">
        <Routes>
          <Route path="/" element={
            <Navigate to={isSuperAdmin ? '/tenants' : (isTenantAdmin ? '/users' : '/access')} replace />
          } />
          <Route path="/access" element={
            <ProtectedRoute>
              <UserPortalPage onConnectToken={handleConnect} onOpenRequestModal={() => setIsModalOpen(true)} />
            </ProtectedRoute>
          } />
          <Route path="/approvals" element={
            <ProtectedRoute requireAdminOrAbove>
              <AdminApprovalPage />
            </ProtectedRoute>
          } />
          <Route path="/resources" element={
            <ProtectedRoute requireAdminOrAbove>
              <ResourceManagerPage />
            </ProtectedRoute>
          } />
          <Route path="/tenants" element={
            <ProtectedRoute requireSuperAdmin>
              <TenantManagementPage />
            </ProtectedRoute>
          } />
          <Route path="/users" element={
            <ProtectedRoute requireTenantAdmin>
              <UserManagementPage />
            </ProtectedRoute>
          } />
          <Route path="/403" element={<ForbiddenPage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </main>
      <TaskRequestModal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />
      <footer className="border-t border-zinc-200 dark:border-white/[0.06] py-6 px-6 text-center text-xs text-zinc-500 bg-white dark:bg-[#0e0f14]">
        <div className="max-w-7xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center space-x-2">
            <ShieldCheck className="w-4 h-4 text-zinc-400" />
            <span>Cloud Terminal (MVP) • Secure Gateway Engine</span>
          </div>
          <div className="flex items-center space-x-6 font-mono text-[11px]">
            <span className="text-emerald-500 font-semibold">● Operational</span>
          </div>
        </div>
      </footer>
    </div>
  );
}

export default App;
