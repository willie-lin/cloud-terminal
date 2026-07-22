import { useAuthStore } from '../stores/useAuthStore';
import { useApprovalStore } from '../stores/useApprovalStore';
import { useResourceStore } from '../stores/useResourceStore';
import { useThemeStore } from '../stores/useThemeStore';
import { useLanguageStore } from '../stores/useLanguageStore';
import { 
  Search, 
  Layers, Building2, Users, 
  CheckSquare, 
  ShieldCheck,
  LogOut,
  Plus,
  Sun,
  Moon,
  Globe
} from 'lucide-react';

interface NavigationProps {
  activeTab: 'access' | 'approvals' | 'resources' | 'tenants' | 'users';
  onSelectTab: (tab: 'access' | 'approvals' | 'resources' | 'tenants' | 'users') => void;
  onOpenRequestModal: () => void;
}

export const Navigation: React.FC<NavigationProps> = ({
  activeTab,
  onSelectTab,
  onOpenRequestModal,
}) => {
  const { user, logout } = useAuthStore();
  const { theme, toggleTheme } = useThemeStore();
  const { lang, toggleLanguage, t } = useLanguageStore();

  const pendingCount = useApprovalStore((state) => 
    state.tasks.filter((t) => t.status === 'pending').length
  );
  const activeTokensCount = useResourceStore((state) => state.activeTokens.length);

  const isSuperAdmin = Boolean(user?.isSuperAdmin || user?.roleName === 'super_admin' || user?.role === 'super_admin');
  const isTenantAdmin = Boolean(user?.isTenantAdmin || user?.roleName === 'tenant_admin' || user?.roleName?.includes('tenant_admin'));
  const isAdminOrAbove = isSuperAdmin || isTenantAdmin || user?.role === 'admin';

  return (
    <header className="glass-header sticky top-0 z-40 w-full px-6 py-3 transition-all">
      <div className="max-w-7xl mx-auto flex items-center justify-between">
        {/* Left Brand & Tabs */}
        <div className="flex items-center space-x-8">
          {/* Unique Brand Logo */}
          <div 
            onClick={() => onSelectTab(isSuperAdmin ? 'tenants' : (isTenantAdmin ? 'users' : 'access'))}
            className="flex items-center space-x-2.5 cursor-pointer group"
          >
            <div className="w-8 h-8 rounded-lg bg-zinc-100 dark:bg-white/10 text-zinc-900 dark:text-white border border-zinc-200 dark:border-white/10 flex items-center justify-center group-hover:border-zinc-400 dark:group-hover:border-white/30 transition-all">
              <svg className="w-4 h-4 stroke-[2]" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <path d="M17.5 19H9a7 7 0 1 1 6.71-9h1.79a4.5 4.5 0 1 1 0 9Z" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="m8.5 12.5 1.5 1.5-1.5 1.5" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M12.5 15.5h2" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>
            </div>
            <div className="flex items-center space-x-2">
              <span className="font-semibold text-zinc-900 dark:text-white tracking-tight text-sm">
                Cloud Terminal
              </span>
              <span className="text-[10px] font-mono px-1.5 py-0.2 rounded bg-zinc-100 dark:bg-white/[0.06] text-zinc-600 dark:text-zinc-400 border border-zinc-200 dark:border-white/[0.08]">
                {isSuperAdmin ? 'SuperAdmin' : (isTenantAdmin ? 'TenantAdmin' : 'User')}
              </span>
            </div>
          </div>

          {/* Dynamic RBAC Tabs */}
          <nav className="hidden md:flex items-center space-x-6">
            {/* 超级管理员标签页：租户管理 */}
            {isSuperAdmin && (
              <button
                onClick={() => onSelectTab('tenants')}
                className={`flex items-center space-x-2 py-1 text-xs font-medium border-b-2 transition-all ${
                  activeTab === 'tenants'
                    ? 'border-zinc-900 dark:border-white text-zinc-900 dark:text-white font-semibold'
                    : 'border-transparent text-zinc-500 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200'
                }`}
              >
                <Building2 className="w-3.5 h-3.5" />
                <span>租户管理</span>
              </button>
            )}

            {/* 超管与租户管理员标签页：用户管理 */}
            {(isSuperAdmin || isTenantAdmin) && (
              <button
                onClick={() => onSelectTab('users')}
                className={`flex items-center space-x-2 py-1 text-xs font-medium border-b-2 transition-all ${
                  activeTab === 'users'
                    ? 'border-zinc-900 dark:border-white text-zinc-900 dark:text-white font-semibold'
                    : 'border-transparent text-zinc-500 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200'
                }`}
              >
                <Users className="w-3.5 h-3.5" />
                <span>用户管理</span>
              </button>
            )}

            {/* 管理员视角：工单审批中心 */}
            {isAdminOrAbove && (
              <button
                onClick={() => onSelectTab('approvals')}
                className={`flex items-center space-x-2 py-1 text-xs font-medium border-b-2 transition-all ${
                  activeTab === 'approvals'
                    ? 'border-zinc-900 dark:border-white text-zinc-900 dark:text-white font-semibold'
                    : 'border-transparent text-zinc-500 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200'
                }`}
              >
                <CheckSquare className="w-3.5 h-3.5" />
                <span>{t('approvals')}</span>
                {pendingCount > 0 && (
                  <span className="px-1.5 py-0.2 text-[10px] rounded-full bg-amber-500/10 dark:bg-amber-500/20 text-amber-600 dark:text-amber-400 font-mono border border-amber-500/30">
                    {pendingCount}
                  </span>
                )}
              </button>
            )}

            {/* 管理员视角：物理资源池 */}
            {isAdminOrAbove && (
              <button
                onClick={() => onSelectTab('resources')}
                className={`flex items-center space-x-2 py-1 text-xs font-medium border-b-2 transition-all ${
                  activeTab === 'resources'
                    ? 'border-zinc-900 dark:border-white text-zinc-900 dark:text-white font-semibold'
                    : 'border-transparent text-zinc-500 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200'
                }`}
              >
                <ShieldCheck className="w-3.5 h-3.5" />
                <span>{t('resources')}</span>
              </button>
            )}

            {/* 普通用户 / 租户管理员视角：我的访问权限 */}
            {(!isSuperAdmin || activeTab === 'access') && (
              <button
                onClick={() => onSelectTab('access')}
                className={`flex items-center space-x-2 py-1 text-xs font-medium border-b-2 transition-all ${
                  activeTab === 'access'
                    ? 'border-zinc-900 dark:border-white text-zinc-900 dark:text-white font-semibold'
                    : 'border-transparent text-zinc-500 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200'
                }`}
              >
                <Layers className="w-3.5 h-3.5" />
                <span>{t('myAccess')}</span>
                {activeTokensCount > 0 && (
                  <span className="px-1.5 py-0.2 text-[10px] rounded-full bg-emerald-500/10 dark:bg-emerald-500/20 text-emerald-600 dark:text-emerald-400 font-mono border border-emerald-500/30">
                    {activeTokensCount}
                  </span>
                )}
              </button>
            )}
          </nav>
        </div>

        {/* Right Actions */}
        <div className="flex items-center space-x-3">
          {/* Quick Search */}
          <div className="hidden lg:flex items-center space-x-2 px-3 py-1.5 rounded-lg bg-zinc-100 dark:bg-white/[0.04] border border-zinc-200 dark:border-white/[0.08] text-zinc-500 dark:text-zinc-400 text-xs hover:border-zinc-300 dark:hover:border-white/20 transition-colors cursor-pointer">
            <Search className="w-3.5 h-3.5 text-zinc-400" />
            <span>{t('searchPlaceholder')}</span>
            <kbd className="px-1.5 py-0.5 text-[10px] font-mono bg-zinc-200 dark:bg-white/10 rounded text-zinc-600 dark:text-zinc-300">⌘K</kbd>
          </div>

          {/* Primary CTA */}
          <button
            onClick={onOpenRequestModal}
            className="inline-flex items-center space-x-1.5 px-3.5 py-1.5 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white font-semibold text-xs shadow-sm active:scale-[0.98] transition-all"
          >
            <Plus className="w-3.5 h-3.5 stroke-[2.5]" />
            <span>{t('newRequest')}</span>
          </button>

          <div className="h-4 w-px bg-zinc-200 dark:bg-white/10" />

          {/* Light / Dark Mode Toggle Button */}
          <button
            onClick={toggleTheme}
            title={theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'}
            className="p-1.5 rounded-lg text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-white/10 transition-colors"
          >
            {theme === 'dark' ? <Sun className="w-4 h-4 text-amber-400" /> : <Moon className="w-4 h-4 text-zinc-700" />}
          </button>

          {/* Language Switcher Button (中 / EN) */}
          <button
            onClick={toggleLanguage}
            title="Switch Language (中 / EN)"
            className="inline-flex items-center space-x-1 px-2 py-1 rounded-lg text-xs font-mono font-medium border border-zinc-200 dark:border-white/10 bg-zinc-100 dark:bg-white/[0.04] text-zinc-700 dark:text-zinc-300 hover:border-zinc-300 dark:hover:border-white/20 transition-all"
          >
            <Globe className="w-3.5 h-3.5 text-zinc-500" />
            <span>{lang === 'zh' ? '中' : 'EN'}</span>
          </button>

          {/* User Profile & Role Badge */}
          <div className="flex items-center space-x-2.5 pl-1">
            <img
              src={user?.avatarUrl}
              alt={user?.username}
              className="w-7 h-7 rounded-full ring-1 ring-zinc-200 dark:ring-white/15 object-cover"
            />
            <div className="hidden sm:block text-left">
              <div className="flex items-center space-x-1.5">
                <span className="text-xs font-semibold text-zinc-900 dark:text-white leading-none">
                  {user?.username}
                </span>
                <span className={`px-1.5 py-0.2 text-[9px] font-mono rounded font-bold uppercase ${
                  user?.role === 'admin'
                    ? 'bg-purple-500/15 text-purple-600 dark:text-purple-300 border border-purple-500/30'
                    : 'bg-cyan-500/15 text-cyan-600 dark:text-cyan-300 border border-cyan-500/30'
                }`}>
                  {user?.role}
                </span>
              </div>
              <p className="text-[10px] text-zinc-500 font-mono mt-0.5">
                {user?.role === 'admin' ? t('securityAdmin') : t('devopsEngineer')}
              </p>
            </div>

            <button
              onClick={logout}
              title="Sign out"
              className="p-1 rounded text-zinc-400 hover:text-zinc-700 dark:hover:text-white transition-colors"
            >
              <LogOut className="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>
    </header>
  );
};
