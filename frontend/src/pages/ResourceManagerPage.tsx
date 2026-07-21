import { useState } from 'react';
import { useResourceStore } from '../stores/useResourceStore';
import { useLanguageStore } from '../stores/useLanguageStore';
import { Plus, Lock, Eye, EyeOff, Terminal, KeyRound } from 'lucide-react';

export const ResourceManagerPage: React.FC = () => {
  const { resources, createResource } = useResourceStore();
  const { t } = useLanguageStore();
  
  const [isAdding, setIsAdding] = useState(false);
  const [name, setName] = useState('');
  const [targetHost, setTargetHost] = useState('');
  const [targetPort, setTargetPort] = useState(22);
  const [authUsername, setAuthUsername] = useState('root');
  const [authCredential, setAuthCredential] = useState('');
  const [description, setDescription] = useState('');
  const [showCred, setShowCred] = useState(false);

  const handleAdd = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !targetHost) return;

    createResource({
      name,
      urn: `urn:ct:resource:ssh:${name.toLowerCase().replace(/\s+/g, '-')}`,
      type: 'ssh',
      ip: targetHost,
      port: targetPort,
      env: 'dev',
      region: 'AWS us-east-1',
      description,
    });

    setName('');
    setTargetHost('');
    setDescription('');
    setAuthCredential('');
    setIsAdding(false);
  };

  return (
    <div className="space-y-6 py-2 text-left">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-zinc-200 dark:border-white/[0.06] pb-4">
        <div>
          <h1 className="text-xl font-bold text-zinc-900 dark:text-white tracking-tight">{t('resourceDirectoryTitle')}</h1>
          <p className="text-xs text-zinc-500 mt-0.5">
            {t('resourceDirectoryDesc')}
          </p>
        </div>

        <button
          onClick={() => setIsAdding(!isAdding)}
          className="inline-flex items-center space-x-1.5 px-3.5 py-1.5 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs hover:bg-zinc-800 dark:hover:bg-zinc-200 transition-all self-start sm:self-auto"
        >
          <Plus className="w-3.5 h-3.5 stroke-[2.5]" />
          <span>{isAdding ? t('cancelRegistration') : t('registerResource')}</span>
        </button>
      </div>

      {/* Form */}
      {isAdding && (
        <form onSubmit={handleAdd} className="p-6 rounded-xl bg-white dark:bg-[#121318] border border-zinc-300 dark:border-white/10 shadow-lg space-y-4 animate-in fade-in duration-150">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">{t('resourceName')}</label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g. Production Payment Worker Node"
                required
                className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs focus:outline-none focus:border-zinc-400 dark:focus:border-white/30 transition-all"
              />
            </div>

            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">{t('targetIp')}</label>
              <input
                type="text"
                value={targetHost}
                onChange={(e) => setTargetHost(e.target.value)}
                placeholder="10.0.88.15"
                required
                className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs font-mono focus:outline-none focus:border-zinc-400 dark:focus:border-white/30 transition-all"
              />
            </div>

            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">{t('sshUsername')}</label>
              <input
                type="text"
                value={authUsername}
                onChange={(e) => setAuthUsername(e.target.value)}
                required
                className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs font-mono focus:outline-none focus:border-zinc-400 dark:focus:border-white/30 transition-all"
              />
            </div>

            <div className="space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase">{t('sshPort')}</label>
              <input
                type="number"
                value={targetPort}
                onChange={(e) => setTargetPort(Number(e.target.value))}
                required
                className="w-full px-3 py-2 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs font-mono focus:outline-none focus:border-zinc-400 dark:focus:border-white/30 transition-all"
              />
            </div>

            <div className="md:col-span-2 space-y-1.5">
              <label className="block text-xs font-semibold text-zinc-600 dark:text-zinc-400 uppercase flex items-center justify-between">
                <span>{t('sshCredential')}</span>
                <button
                  type="button"
                  onClick={() => setShowCred(!showCred)}
                  className="text-zinc-500 hover:text-zinc-900 dark:hover:text-white flex items-center space-x-1 font-normal text-[11px]"
                >
                  {showCred ? <EyeOff className="w-3.5 h-3.5" /> : <Eye className="w-3.5 h-3.5" />}
                  <span>{showCred ? 'Hide' : 'Reveal'}</span>
                </button>
              </label>
              <textarea
                value={authCredential}
                onChange={(e) => setAuthCredential(e.target.value)}
                placeholder="-----BEGIN OPENSSH PRIVATE KEY-----..."
                rows={3}
                className="w-full p-3 rounded-lg bg-zinc-50 dark:bg-black/50 border border-zinc-200 dark:border-white/10 text-zinc-900 dark:text-white text-xs font-mono focus:outline-none focus:border-zinc-400 dark:focus:border-white/30 transition-all resize-none"
              />
            </div>
          </div>

          <div className="flex items-center justify-end space-x-3 pt-2">
            <button
              type="button"
              onClick={() => setIsAdding(false)}
              className="px-3.5 py-1.5 rounded-lg text-xs font-medium text-zinc-500 hover:text-zinc-900 dark:hover:text-white transition-colors"
            >
              {t('cancel')}
            </button>
            <button
              type="submit"
              className="px-4 py-1.5 rounded-lg bg-zinc-900 text-white dark:bg-white dark:text-black font-semibold text-xs hover:bg-zinc-800 dark:hover:bg-zinc-200 transition-all"
            >
              {t('saveResource')}
            </button>
          </div>
        </form>
      )}

      {/* Inventory Table */}
      <div className="rounded-xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/[0.07] overflow-hidden shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-zinc-200 dark:border-white/[0.06] bg-zinc-50 dark:bg-black/40 text-[10px] uppercase font-mono tracking-wider text-zinc-500 dark:text-zinc-400">
                <th className="py-3 px-5">{t('resourceName')}</th>
                <th className="py-3 px-5">{t('targetIp')}</th>
                <th className="py-3 px-5">OS / Region</th>
                <th className="py-3 px-5">Security</th>
                <th className="py-3 px-5">{t('status')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-200 dark:divide-white/[0.04] text-xs">
              {resources.map((res) => (
                <tr key={res.id} className="hover:bg-zinc-50 dark:hover:bg-white/[0.02] transition-colors">
                  <td className="py-3 px-5">
                    <div className="flex items-center space-x-2.5">
                      <div className="p-1.5 rounded bg-zinc-100 dark:bg-white/[0.04] border border-zinc-200 dark:border-white/[0.08] text-zinc-700 dark:text-zinc-300">
                        <Terminal className="w-3.5 h-3.5" />
                      </div>
                      <div>
                        <p className="font-semibold text-zinc-900 dark:text-white">{res.name}</p>
                        <p className="text-[10px] text-zinc-500 font-mono">{res.urn}</p>
                      </div>
                    </div>
                  </td>

                  <td className="py-3 px-5 font-mono text-zinc-700 dark:text-zinc-300">
                    <div className="flex items-center space-x-1.5">
                      <Lock className="w-3.5 h-3.5 text-zinc-400" />
                      <span>{res.authUsername || 'root'}@{res.targetHost}:{res.targetPort || 22}</span>
                    </div>
                  </td>

                  <td className="py-3 px-5 text-zinc-500">
                    <p className="text-zinc-700 dark:text-zinc-300">{res.os || 'Linux x86_64'}</p>
                    <p className="text-[10px] font-mono text-zinc-400">{res.region || 'AWS'}</p>
                  </td>

                  <td className="py-3 px-5">
                    <span className="inline-flex items-center space-x-1 px-2 py-0.5 rounded bg-zinc-100 dark:bg-white/[0.04] border border-zinc-200 dark:border-white/[0.08] text-zinc-700 dark:text-zinc-300 font-mono text-[10px]">
                      <KeyRound className="w-3 h-3 text-zinc-400" />
                      <span>{t('encryptedDek')}</span>
                    </span>
                  </td>

                  <td className="py-3 px-5">
                    <span className="inline-flex items-center space-x-1.5 px-2 py-0.5 rounded-full text-[10px] font-mono border border-emerald-500/30 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
                      <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                      <span className="capitalize">{res.status}</span>
                    </span>
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
