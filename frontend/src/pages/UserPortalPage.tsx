import { useResourceStore } from '../stores/useResourceStore';
import { useLanguageStore } from '../stores/useLanguageStore';
import { ResourceCard } from '../components/ResourceCard';
import type { AccessPolicyToken } from '../types';
import { Plus, ArrowRight, ShieldCheck } from 'lucide-react';

interface UserPortalPageProps {
  onConnectToken: (token: AccessPolicyToken) => void;
  onOpenRequestModal: () => void;
}

export const UserPortalPage: React.FC<UserPortalPageProps> = ({
  onConnectToken,
  onOpenRequestModal,
}) => {
  const { activeTokens } = useResourceStore();
  const { t } = useLanguageStore();

  return (
    <div className="space-y-8 py-2 text-left">
      {/* Header Banner */}
      <div className="relative rounded-2xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/[0.08] p-8 shadow-sm">
        <div className="max-w-2xl space-y-3">
          <div className="inline-flex items-center space-x-2 px-2.5 py-0.5 rounded-md bg-zinc-100 dark:bg-white/[0.04] border border-zinc-200 dark:border-white/[0.08] text-zinc-700 dark:text-zinc-300 text-[11px] font-mono">
            <ShieldCheck className="w-3.5 h-3.5 text-zinc-500 dark:text-zinc-400" />
            <span>{t('zeroTrustJit')}</span>
          </div>
          
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-white tracking-tight">
            {t('portalTitle')}
          </h1>

          <p className="text-xs text-zinc-600 dark:text-zinc-400 leading-relaxed">
            {t('portalDesc')}
          </p>

          <div className="pt-2 flex items-center space-x-3">
            <button
              onClick={onOpenRequestModal}
              className="inline-flex items-center space-x-1.5 px-4 py-2 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs hover:bg-zinc-800 dark:hover:bg-zinc-200 active:scale-[0.98] transition-all"
            >
              <Plus className="w-3.5 h-3.5 stroke-[2.5]" />
              <span>{t('applyForAccess')}</span>
            </button>

            <a
              href="#active-resources"
              className="inline-flex items-center space-x-1 px-3 py-2 rounded-lg bg-zinc-100 dark:bg-white/[0.03] border border-zinc-200 dark:border-white/[0.06] text-zinc-700 dark:text-zinc-300 text-xs hover:bg-zinc-200 dark:hover:bg-white/[0.06] transition-colors"
            >
              <span>{t('activeSessionsCount')} ({activeTokens.length})</span>
              <ArrowRight className="w-3.5 h-3.5 text-zinc-400" />
            </a>
          </div>
        </div>
      </div>

      {/* Active Tokens Grid */}
      <div id="active-resources" className="space-y-4">
        <div className="flex items-center justify-between border-b border-zinc-200 dark:border-white/[0.06] pb-3">
          <div>
            <h2 className="text-base font-semibold text-zinc-900 dark:text-white tracking-tight">{t('activeStsConnections')}</h2>
            <p className="text-xs text-zinc-500">{t('stsTokensDesc')}</p>
          </div>
        </div>

        {activeTokens.length === 0 ? (
          <div className="p-12 rounded-xl bg-white dark:bg-[#121318]/50 border border-zinc-200 dark:border-white/[0.06] text-center space-y-3">
            <p className="text-xs text-zinc-500 dark:text-zinc-400">{t('noActiveTokens')}</p>
            <button
              onClick={onOpenRequestModal}
              className="inline-flex items-center space-x-1.5 px-3 py-1.5 rounded-lg bg-zinc-100 dark:bg-white/10 hover:bg-zinc-200 dark:hover:bg-white/20 text-zinc-900 dark:text-white font-semibold text-xs transition-all"
            >
              <Plus className="w-3.5 h-3.5" />
              <span>{t('applyForAccess')}</span>
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
            {activeTokens.map((token) => (
              <ResourceCard
                key={token.token}
                token={token}
                onConnect={onConnectToken}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
