import { useState, useEffect } from 'react';
import { useAuthStore } from './stores/useAuthStore';
import { useResourceStore } from './stores/useResourceStore';
import { useApprovalStore } from './stores/useApprovalStore';
import { useLanguageStore } from './stores/useLanguageStore';
import { Navigation } from './components/Navigation';
import { TaskRequestModal } from './components/TaskRequestModal';
import { UserPortalPage } from './pages/UserPortalPage';
import { AdminApprovalPage } from './pages/AdminApprovalPage';
import { ResourceManagerPage } from './pages/ResourceManagerPage';
import { TenantManagementPage } from './pages/TenantManagementPage';
import { UserManagementPage } from './pages/UserManagementPage';
import { TerminalPage } from './pages/TerminalPage';
import type { AccessPolicyToken } from './types';
import { ShieldCheck, ArrowRight, Lock, Mail, KeyRound } from 'lucide-react';
import './App.css';

export function App() {
  const { user, isAuthenticated, isLoading, error, login } = useAuthStore();
  const { activeSession, startSession, endSession, fetchResources } = useResourceStore();
  const { fetchTasks } = useApprovalStore();
  const { t } = useLanguageStore();

  const [activeTab, setActiveTab] = useState<'access' | 'approvals' | 'resources' | 'tenants' | 'users'>('access');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [registerMode, setRegisterMode] = useState(false);
  const [tenantName, setTenantName] = useState('');

  useEffect(() => {
    if (isAuthenticated) {
      fetchResources();
      fetchTasks();
    }
  }, [isAuthenticated, fetchResources, fetchTasks]);

  const handleConnect = (token: AccessPolicyToken) => {
    startSession(token);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) return;
    try {
      if (registerMode) {
        await useAuthStore.getState().register({ email, password, tenant_name: tenantName || email.split('@')[0] });
      } else {
        await login({ email, password });
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
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase tracking-wider">邮箱</label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-400" />
                <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="you@company.com" autoComplete="email"
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
            {registerMode && (
              <div className="space-y-2">
                <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase tracking-wider">公司名称</label>
                <input type="text" value={tenantName} onChange={(e) => setTenantName(e.target.value)} placeholder="your-company"
                  className="w-full px-4 py-2.5 rounded-xl bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/[0.08] text-zinc-900 dark:text-white text-sm focus:outline-none focus:border-zinc-400 transition-all" />
              </div>
            )}
            <div className="p-3.5 rounded-xl bg-zinc-100 dark:bg-white/[0.03] border border-zinc-200 dark:border-white/[0.06] flex items-start space-x-3 text-xs text-zinc-500 dark:text-zinc-400">
              <Lock className="w-4 h-4 text-zinc-400 shrink-0 mt-0.5" />
              <span>{t('loginSecurityNotice')}</span>
            </div>
            <button type="submit" disabled={!email || !password || isLoading}
              className="w-full py-2.5 px-5 rounded-xl bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-sm hover:bg-zinc-800 dark:hover:bg-zinc-200 disabled:opacity-40 disabled:cursor-not-allowed transition-all flex items-center justify-center space-x-2">
              <span>{isLoading ? '登录中...' : (registerMode ? '注册' : t('enterControlPlane'))}</span>
              <ArrowRight className="w-4 h-4 stroke-[2.5]" />
            </button>
            <div className="text-center">
              <button type="button" onClick={() => setRegisterMode(!registerMode)}
                className="text-xs text-zinc-500 hover:text-zinc-900 dark:hover:text-white transition-colors">
                {registerMode ? '已有账号？点此登录' : '没有账号？点此注册'}
              </button>
            </div>
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
      <Navigation activeTab={activeTab} onSelectTab={setActiveTab} onOpenRequestModal={() => setIsModalOpen(true)} />
      <main className="flex-1 max-w-7xl w-full mx-auto px-6 py-6 pb-20">
        {activeTab === 'access' && <UserPortalPage onConnectToken={handleConnect} onOpenRequestModal={() => setIsModalOpen(true)} />}
        {activeTab === 'approvals' && user?.role === 'admin' && <AdminApprovalPage />}
        {activeTab === 'resources' && user?.role === 'admin' && <ResourceManagerPage />}
        {activeTab === 'tenants' && user?.roleName === 'super_admin' && <TenantManagementPage />}
        {activeTab === 'users' && (user?.roleName === 'super_admin' || user?.roleName?.includes('tenant_admin')) && <UserManagementPage />}
      </main>
      <TaskRequestModal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />
      <footer className="border-t border-zinc-200 dark:border-white/[0.06] py-6 px-6 text-center text-xs text-zinc-500 bg-white dark:bg-[#0e0f14]">
        <div className="max-w-7xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center space-x-2">
            <ShieldCheck className="w-4 h-4 text-zinc-400" />
            <span>Cloud Terminal (MVP) • Secure Gateway Engine</span>
          </div>
          <div className="flex items-center space-x-6 font-mono text-[11px]">
            <span>Task → Approval → Policy → Resource → Session</span>
            <span className="text-emerald-500 font-semibold">● Operational</span>
          </div>
        </div>
      </footer>
    </div>
  );
}

export default App;
