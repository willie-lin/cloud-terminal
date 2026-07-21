import { useState } from 'react';
import { useApprovalStore } from '../stores/useApprovalStore';
import { useLanguageStore } from '../stores/useLanguageStore';
import type { TaskStatus } from '../types';
import { 
  CheckCircle, 
  XCircle, 
  Server, 
  Search
} from 'lucide-react';

export const ApprovalTable: React.FC = () => {
  const { tasks, approveTask, rejectTask } = useApprovalStore();
  const { t } = useLanguageStore();
  const [filterStatus, setFilterStatus] = useState<string>('all');
  const [searchQuery, setSearchQuery] = useState<string>('');

  const filteredTasks = tasks.filter((task) => {
    const matchesStatus = filterStatus === 'all' || task.status === filterStatus;
    const matchesSearch = 
      task.userName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      task.resourceName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      task.reason.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesStatus && matchesSearch;
  });

  const getStatusBadge = (status: TaskStatus) => {
    switch (status) {
      case 'pending':
        return (
          <span className="inline-flex items-center space-x-1.5 px-2 py-0.5 rounded text-[10px] font-mono border border-amber-500/30 bg-amber-500/10 text-amber-600 dark:text-amber-400">
            <span className="w-1 h-1 rounded-full bg-amber-500" />
            <span>{t('pending').toUpperCase()}</span>
          </span>
        );
      case 'approved':
        return (
          <span className="inline-flex items-center space-x-1 px-2 py-0.5 rounded text-[10px] font-mono border border-emerald-500/30 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
            <CheckCircle className="w-3 h-3" />
            <span>{t('approved').toUpperCase()}</span>
          </span>
        );
      case 'rejected':
        return (
          <span className="inline-flex items-center space-x-1 px-2 py-0.5 rounded text-[10px] font-mono border border-red-500/30 bg-red-500/10 text-red-600 dark:text-red-400">
            <XCircle className="w-3 h-3" />
            <span>{t('rejected').toUpperCase()}</span>
          </span>
        );
      default:
        return null;
    }
  };

  return (
    <div className="space-y-4 text-left">
      {/* Controls */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-3">
        {/* Search */}
        <div className="relative w-full sm:w-72">
          <input
            type="text"
            placeholder={t('searchPlaceholder')}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-9 pr-3 py-1.5 rounded-lg bg-zinc-100 dark:bg-white/[0.03] border border-zinc-300 dark:border-white/[0.08] text-zinc-900 dark:text-white text-xs placeholder-zinc-500 focus:outline-none focus:border-zinc-400 dark:focus:border-white/20 transition-all"
          />
          <Search className="absolute left-3 top-2 w-3.5 h-3.5 text-zinc-400 pointer-events-none" />
        </div>

        {/* Filter Pills */}
        <div className="flex items-center space-x-1 p-0.5 rounded-lg bg-zinc-100 dark:bg-white/[0.03] border border-zinc-300 dark:border-white/[0.06]">
          {['all', 'pending', 'approved', 'rejected'].map((status) => (
            <button
              key={status}
              onClick={() => setFilterStatus(status)}
              className={`px-2.5 py-1 rounded text-xs capitalize transition-all ${
                filterStatus === status
                  ? 'bg-zinc-900 text-white dark:bg-white/10 dark:text-white font-semibold shadow-sm'
                  : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200'
              }`}
            >
              {status === 'all' ? t('all') : status === 'pending' ? t('pending') : status === 'approved' ? t('approved') : t('rejected')}
            </button>
          ))}
        </div>
      </div>

      {/* Modern High-Contrast Table */}
      <div className="rounded-xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/[0.07] overflow-hidden shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-zinc-200 dark:border-white/[0.06] bg-zinc-50 dark:bg-black/40 text-[10px] uppercase font-mono tracking-wider text-zinc-500 dark:text-zinc-400">
                <th className="py-3 px-5">{t('applicant')}</th>
                <th className="py-3 px-5">{t('targetResource')}</th>
                <th className="py-3 px-5">{t('duration')}</th>
                <th className="py-3 px-5">{t('justification')}</th>
                <th className="py-3 px-5">{t('status')}</th>
                <th className="py-3 px-5 text-right">{t('actions')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-200 dark:divide-white/[0.04] text-xs">
              {filteredTasks.map((task) => (
                <tr key={task.id} className="hover:bg-zinc-50 dark:hover:bg-white/[0.02] transition-colors">
                  <td className="py-3 px-5">
                    <div className="flex items-center space-x-2.5">
                      <img
                        src={task.userAvatar}
                        alt={task.userName}
                        className="w-7 h-7 rounded-full ring-1 ring-zinc-300 dark:ring-white/10 object-cover"
                      />
                      <div>
                        <p className="font-semibold text-zinc-900 dark:text-white">{task.userName}</p>
                        <p className="text-[10px] text-zinc-500 font-mono">{task.userEmail}</p>
                      </div>
                    </div>
                  </td>

                  <td className="py-3 px-5">
                    <div className="flex items-center space-x-2">
                      <Server className="w-3.5 h-3.5 text-zinc-400" />
                      <div>
                        <p className="font-semibold text-zinc-900 dark:text-white">{task.resourceName}</p>
                        <p className="text-[10px] text-zinc-500 font-mono">
                          {task.resourceUrn.split(':').pop()}
                        </p>
                      </div>
                    </div>
                  </td>

                  <td className="py-3 px-5">
                    <span className="font-mono text-zinc-700 dark:text-zinc-300 text-xs">{task.durationHours}{t('hours')}</span>
                  </td>

                  <td className="py-3 px-5 max-w-xs">
                    <p className="text-zinc-600 dark:text-zinc-400 line-clamp-1 leading-relaxed" title={task.reason}>
                      {task.reason}
                    </p>
                  </td>

                  <td className="py-3 px-5">
                    {getStatusBadge(task.status)}
                  </td>

                  <td className="py-3 px-5 text-right">
                    {task.status === 'pending' ? (
                      <div className="flex items-center justify-end space-x-2">
                        <button
                          onClick={() => approveTask(task.id)}
                          className="px-2.5 py-1 rounded bg-zinc-900 text-white dark:bg-white dark:text-black hover:bg-zinc-800 dark:hover:bg-zinc-200 font-semibold text-xs transition-all"
                        >
                          {t('approve')}
                        </button>
                        <button
                          onClick={() => rejectTask(task.id)}
                          className="px-2.5 py-1 rounded bg-zinc-100 dark:bg-white/5 hover:bg-red-500 hover:text-white border border-zinc-200 dark:border-white/10 text-zinc-700 dark:text-zinc-300 text-xs transition-all"
                        >
                          {t('reject')}
                        </button>
                      </div>
                    ) : (
                      <span className="text-zinc-500 font-mono text-[10px]">
                        {t('reviewedAt')} ({task.reviewedAt})
                      </span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
