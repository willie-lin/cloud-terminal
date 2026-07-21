import { useState, useEffect } from 'react';
import type { AccessPolicyToken } from '../types';
import { useLanguageStore } from '../stores/useLanguageStore';
import { Terminal, Clock, Server, KeyRound, ExternalLink } from 'lucide-react';

interface ResourceCardProps {
  token: AccessPolicyToken;
  onConnect: (token: AccessPolicyToken) => void;
}

export const ResourceCard: React.FC<ResourceCardProps> = ({ token, onConnect }) => {
  const { t } = useLanguageStore();
  const [timeLeft, setTimeLeft] = useState<string>('');
  const [isExpired, setIsExpired] = useState<boolean>(false);

  useEffect(() => {
    const updateCountdown = () => {
      const now = Math.floor(Date.now() / 1000);
      const diff = token.expiresAt - now;

      if (diff <= 0) {
        setIsExpired(true);
        setTimeLeft('00:00:00');
        return;
      }

      const hours = Math.floor(diff / 3600);
      const minutes = Math.floor((diff % 3600) / 60);
      const seconds = diff % 60;

      if (hours > 24) {
        const days = Math.floor(hours / 24);
        const remainingHours = hours % 24;
        setTimeLeft(`${days}d ${remainingHours}h`);
      } else {
        const fmt = (num: number) => num.toString().padStart(2, '0');
        setTimeLeft(`${fmt(hours)}:${fmt(minutes)}:${fmt(seconds)}`);
      }
    };

    updateCountdown();
    const interval = setInterval(updateCountdown, 1000);
    return () => clearInterval(interval);
  }, [token]);

  return (
    <div 
      className={`relative rounded-xl p-5 border transition-all duration-200 ${
        isExpired 
          ? 'bg-zinc-100 dark:bg-[#121318]/40 border-zinc-200 dark:border-white/[0.04] opacity-60' 
          : 'bg-white dark:bg-[#121318] border-zinc-200 dark:border-white/[0.07] hover:border-zinc-400 dark:hover:border-white/20 shadow-sm'
      }`}
    >
      <div className="flex flex-col justify-between h-full space-y-5">
        {/* Header */}
        <div className="flex items-start justify-between">
          <div className="flex items-center space-x-3">
            <div className="p-2 rounded-lg bg-zinc-100 dark:bg-white/[0.04] border border-zinc-200 dark:border-white/[0.08] text-zinc-700 dark:text-zinc-300">
              <Server className="w-4 h-4 stroke-[1.8]" />
            </div>
            <div>
              <h3 className="font-semibold text-zinc-900 dark:text-white text-sm tracking-tight">
                {token.resourceName}
              </h3>
              <p className="text-[11px] font-mono text-zinc-500 mt-0.5">
                {token.resourceUrn.split(':').pop()}
              </p>
            </div>
          </div>

          {/* Status Badge */}
          <div className="flex items-center space-x-1.5 px-2 py-0.5 rounded-full text-[10px] font-mono border border-zinc-200 dark:border-white/[0.08] bg-zinc-100 dark:bg-white/[0.02]">
            <span className={`w-1.5 h-1.5 rounded-full ${isExpired ? 'bg-zinc-400' : 'bg-emerald-500'}`} />
            <span className={isExpired ? 'text-zinc-500' : 'text-emerald-600 dark:text-emerald-400 font-semibold'}>
              {isExpired ? t('expired') : t('online')}
            </span>
          </div>
        </div>

        {/* Expiration & Token Metadata */}
        <div className="p-3 rounded-lg bg-zinc-50 dark:bg-black/40 border border-zinc-200 dark:border-white/[0.05] space-y-2 text-xs">
          <div className="flex items-center justify-between">
            <span className="text-zinc-500 dark:text-zinc-400 flex items-center space-x-1.5 text-[11px]">
              <Clock className="w-3.5 h-3.5 text-zinc-400" />
              <span>{t('timeRemaining')}:</span>
            </span>
            <span className={`font-mono font-medium text-[11px] ${isExpired ? 'text-zinc-400' : 'text-zinc-800 dark:text-zinc-200'}`}>
              {timeLeft}
            </span>
          </div>

          <div className="flex items-center justify-between border-t border-zinc-200 dark:border-white/[0.04] pt-2 text-[11px]">
            <span className="text-zinc-500 dark:text-zinc-400 flex items-center space-x-1.5">
              <KeyRound className="w-3.5 h-3.5 text-zinc-400" />
              <span>{t('stsTokenId')}:</span>
            </span>
            <span className="font-mono text-zinc-600 dark:text-zinc-400 bg-zinc-200/60 dark:bg-white/[0.04] px-1.5 py-0.2 rounded border border-zinc-300 dark:border-white/[0.06]">
              {token.token.substring(0, 16)}...
            </span>
          </div>
        </div>

        {/* Action Button */}
        <button
          disabled={isExpired}
          onClick={() => onConnect(token)}
          className={`w-full py-2 px-3 rounded-lg font-semibold text-xs flex items-center justify-center space-x-2 transition-all ${
            isExpired
              ? 'bg-zinc-200 dark:bg-zinc-800 text-zinc-400 dark:text-zinc-600 cursor-not-allowed border border-zinc-300 dark:border-white/5'
              : 'bg-zinc-900 text-white dark:bg-white dark:text-black hover:bg-zinc-800 dark:hover:bg-zinc-200 active:scale-[0.99]'
          }`}
        >
          <Terminal className="w-3.5 h-3.5 stroke-[2.2]" />
          <span>{isExpired ? t('sessionExpired') : t('connectSshTerminal')}</span>
          {!isExpired && <ExternalLink className="w-3 h-3 ml-0.5 opacity-60" />}
        </button>
      </div>
    </div>
  );
};
