import { useState } from 'react';
import { useAuthStore } from '../stores/useAuthStore';
import { useResourceStore } from '../stores/useResourceStore';
import { useApprovalStore } from '../stores/useApprovalStore';
import { useLanguageStore } from '../stores/useLanguageStore';
import { X, Server, CheckCircle2, ShieldAlert } from 'lucide-react';

interface TaskRequestModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const TaskRequestModal: React.FC<TaskRequestModalProps> = ({ isOpen, onClose }) => {
  const { user } = useAuthStore();
  const { resources } = useResourceStore();
  const { createTaskRequest } = useApprovalStore();
  const { t } = useLanguageStore();

  const [selectedResourceId, setSelectedResourceId] = useState<string>(resources[0]?.id || '');
  const [durationHours, setDurationHours] = useState<number>(2);
  const [reason, setReason] = useState<string>('');
  const [submitted, setSubmitted] = useState<boolean>(false);

  if (!isOpen) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!user || !selectedResourceId || !reason.trim()) return;

    const resource = resources.find((r) => r.id === selectedResourceId);
    if (!resource) return;

    createTaskRequest(
      user.id,
      user.username,
      user.email,
      resource.id,
      resource.urn,
      resource.name,
      reason,
      durationHours
    );

    setSubmitted(true);
    setTimeout(() => {
      setSubmitted(false);
      setReason('');
      onClose();
    }, 1200);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 dark:bg-black/80 backdrop-blur-sm animate-in fade-in duration-150">
      <div 
        className="relative w-full max-w-lg rounded-xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/10 shadow-2xl overflow-hidden text-left"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="px-6 py-4 flex items-center justify-between border-b border-zinc-200 dark:border-white/[0.08]">
          <div>
            <h2 className="text-base font-semibold text-zinc-900 dark:text-white tracking-tight">{t('requestTitle')}</h2>
            <p className="text-xs text-zinc-500 dark:text-zinc-400">{t('requestDesc')}</p>
          </div>
          <button 
            onClick={onClose}
            className="p-1 rounded text-zinc-400 hover:text-zinc-900 dark:hover:text-white hover:bg-zinc-100 dark:hover:bg-white/10 transition-colors"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Modal Body */}
        {submitted ? (
          <div className="p-10 flex flex-col items-center justify-center text-center space-y-3">
            <div className="w-12 h-12 rounded-full bg-emerald-500/10 dark:bg-white/10 border border-emerald-500/30 dark:border-white/20 flex items-center justify-center text-emerald-600 dark:text-white">
              <CheckCircle2 className="w-6 h-6" />
            </div>
            <div>
              <h3 className="text-base font-semibold text-zinc-900 dark:text-white">{t('requestSubmitted')}</h3>
              <p className="text-xs text-zinc-500 dark:text-zinc-400 max-w-xs mt-1">
                {t('requestSubmittedDesc')}
              </p>
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="p-6 space-y-4">
            {/* Target Resource Picker */}
            <div className="space-y-1.5">
              <label className="block text-[11px] font-mono font-semibold text-zinc-600 dark:text-zinc-400 uppercase tracking-wider">
                {t('targetResource')}
              </label>
              <div className="relative">
                <select
                  value={selectedResourceId}
                  onChange={(e) => setSelectedResourceId(e.target.value)}
                  className="w-full pl-9 pr-4 py-2.5 rounded-lg bg-zinc-50 dark:bg-black/60 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400 dark:focus:border-white/30 transition-all appearance-none"
                  required
                >
                  {resources.map((res) => (
                    <option key={res.id} value={res.id} className="bg-white dark:bg-[#121318] text-zinc-900 dark:text-white">
                      {res.name} — ({res.os || res.type.toUpperCase()})
                    </option>
                  ))}
                </select>
                <Server className="absolute left-3 top-3 w-3.5 h-3.5 text-zinc-400 pointer-events-none" />
              </div>
            </div>

            {/* Time Duration Pills */}
            <div className="space-y-1.5">
              <label className="block text-[11px] font-mono font-semibold text-zinc-600 dark:text-zinc-400 uppercase tracking-wider flex items-center justify-between">
                <span>{t('requestedDuration')}</span>
                <span className="text-zinc-900 dark:text-white font-mono">{durationHours} {t('hours')}</span>
              </label>
              <div className="grid grid-cols-5 gap-2">
                {[1, 2, 4, 8, 24].map((hrs) => (
                  <button
                    key={hrs}
                    type="button"
                    onClick={() => setDurationHours(hrs)}
                    className={`py-1.5 rounded-lg text-xs font-mono transition-all border ${
                      durationHours === hrs
                        ? 'bg-zinc-900 text-white dark:bg-white dark:text-black font-bold border-zinc-900 dark:border-white'
                        : 'bg-zinc-100 dark:bg-white/[0.03] text-zinc-600 dark:text-zinc-400 border-zinc-200 dark:border-white/[0.08] hover:border-zinc-400 dark:hover:border-white/20'
                    }`}
                  >
                    {hrs}h
                  </button>
                ))}
              </div>
            </div>

            {/* Justification Reason */}
            <div className="space-y-1.5">
              <label className="block text-[11px] font-mono font-semibold text-zinc-600 dark:text-zinc-400 uppercase tracking-wider">
                {t('taskDescription')}
              </label>
              <textarea
                value={reason}
                onChange={(e) => setReason(e.target.value)}
                placeholder={t('reasonPlaceholder')}
                rows={3}
                required
                className="w-full p-3 rounded-lg bg-zinc-50 dark:bg-black/60 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs placeholder-zinc-400 focus:outline-none focus:border-zinc-400 dark:focus:border-white/30 transition-all resize-none"
              />
            </div>

            {/* Info notice */}
            <div className="p-3 rounded-lg bg-zinc-100 dark:bg-white/[0.03] border border-zinc-200 dark:border-white/[0.06] flex items-start space-x-2.5 text-[11px] text-zinc-600 dark:text-zinc-400">
              <ShieldAlert className="w-4 h-4 text-zinc-500 shrink-0 mt-0.5" />
              <span>{t('zeroTrustNotice')}</span>
            </div>

            {/* Actions */}
            <div className="pt-2 flex items-center justify-end space-x-3">
              <button
                type="button"
                onClick={onClose}
                className="px-3.5 py-2 rounded-lg text-xs font-medium text-zinc-500 hover:text-zinc-900 dark:hover:text-white transition-colors"
              >
                {t('cancel')}
              </button>
              <button
                type="submit"
                className="px-4 py-2 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs hover:bg-zinc-800 dark:hover:bg-zinc-200 active:scale-[0.98] transition-all"
              >
                {t('submitRequest')}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
};
