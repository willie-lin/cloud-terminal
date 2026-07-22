import React from 'react';
import { ShieldAlert, ArrowLeft } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../stores/useAuthStore';

export const ForbiddenPage: React.FC = () => {
  const navigate = useNavigate();
  const { user } = useAuthStore();

  const handleBackToHome = () => {
    const isSuperAdmin = Boolean(user?.isSuperAdmin || user?.roleName === 'super_admin' || user?.role === 'super_admin');
    const isTenantAdmin = Boolean(user?.isTenantAdmin || user?.roleName === 'tenant_admin' || user?.roleName?.includes('tenant_admin'));

    if (isSuperAdmin) {
      navigate('/tenants');
    } else if (isTenantAdmin) {
      navigate('/users');
    } else {
      navigate('/access');
    }
  };

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-[#09090b] flex items-center justify-center p-6 text-left">
      <div className="w-full max-w-md rounded-2xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/10 p-8 shadow-2xl text-center space-y-5">
        <div className="w-14 h-14 mx-auto rounded-2xl bg-red-500/10 text-red-600 dark:text-red-400 flex items-center justify-center border border-red-500/20">
          <ShieldAlert className="w-7 h-7" />
        </div>

        <div>
          <h1 className="text-xl font-bold text-zinc-900 dark:text-white tracking-tight">403 权限不足</h1>
          <p className="text-xs text-zinc-500 mt-1">
            您当前的身份级别为 <span className="font-mono font-bold text-zinc-700 dark:text-zinc-300">{user?.roleName || user?.role || '普通用户'}</span>，无权访问该模块。
          </p>
        </div>

        <div className="p-3.5 rounded-xl bg-zinc-100 dark:bg-white/[0.03] border border-zinc-200 dark:border-white/[0.06] text-xs text-zinc-500 dark:text-zinc-400 text-left font-mono space-y-1">
          <div>租户 ID: {user?.tenantId || 'systemTenant'}</div>
          <div>权限规则: RBAC Restricted Endpoint</div>
        </div>

        <button
          onClick={handleBackToHome}
          className="w-full py-2.5 px-5 rounded-xl bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs hover:bg-zinc-800 dark:hover:bg-zinc-200 transition-all flex items-center justify-center space-x-2"
        >
          <ArrowLeft className="w-4 h-4" />
          <span>返回我的专属控制台</span>
        </button>
      </div>
    </div>
  );
};
