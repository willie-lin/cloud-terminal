import { ApprovalTable } from '../components/ApprovalTable';
import { useApprovalStore } from '../stores/useApprovalStore';
import { useLanguageStore } from '../stores/useLanguageStore';

export const AdminApprovalPage: React.FC = () => {
  const { tasks } = useApprovalStore();
  const { t } = useLanguageStore();
  
  const pendingCount = tasks.filter((t) => t.status === 'pending').length;
  const approvedCount = tasks.filter((t) => t.status === 'approved').length;

  return (
    <div className="space-y-6 py-2 text-left">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-zinc-200 dark:border-white/[0.06] pb-4">
        <div>
          <h1 className="text-xl font-bold text-zinc-900 dark:text-white tracking-tight">{t('approvalCenterTitle')}</h1>
          <p className="text-xs text-zinc-500 mt-0.5">
            {t('approvalCenterDesc')}
          </p>
        </div>

        <div className="flex items-center space-x-3 text-xs">
          <div className="px-3 py-1 rounded-lg bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/[0.08] font-mono">
            <span className="text-zinc-500">{t('pending')}: </span>
            <span className="text-amber-600 dark:text-amber-400 font-bold">{pendingCount}</span>
          </div>

          <div className="px-3 py-1 rounded-lg bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/[0.08] font-mono">
            <span className="text-zinc-500">{t('approved')}: </span>
            <span className="text-emerald-600 dark:text-emerald-400 font-bold">{approvedCount}</span>
          </div>
        </div>
      </div>

      <ApprovalTable />
    </div>
  );
};
