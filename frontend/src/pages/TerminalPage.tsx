import { useEffect, useRef, useState, useCallback } from 'react';
import type { TerminalSession } from '../types';
import { useLanguageStore } from '../stores/useLanguageStore';
import { connectTerminalSession } from '../api/sessions';
import type { SessionController } from '../api/sessions';
import { 
  Clock, 
  Wifi, 
  ArrowLeft, 
  Box, 
  ShieldCheck, 
  RefreshCw, 
  Copy, 
  Check, 
  Server, 
  Key, 
  Cpu, 
  ShieldAlert,
} from 'lucide-react';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import '@xterm/xterm/css/xterm.css';

interface TerminalPageProps {
  session: TerminalSession;
  onDisconnect: () => void;
}

export const TerminalPage: React.FC<TerminalPageProps> = ({ session, onDisconnect }) => {
  const { t } = useLanguageStore();
  const terminalRef = useRef<HTMLDivElement>(null);
  const xtermInstance = useRef<Terminal | null>(null);
  const fitAddonRef = useRef<FitAddon | null>(null);
  const sessionControllerRef = useRef<SessionController | null>(null);
  
  const [containerId] = useState(`ct-sandbox-${session.sessionId.slice(-6)}`);
  const [copied, setCopied] = useState(false);
  const [wsConnected, setWsConnected] = useState(false);
  const [remainingSeconds, setRemainingSeconds] = useState<number>(() => {
    const diff = Math.floor((session.expiresAt * 1000 - Date.now()) / 1000);
    return diff > 0 ? diff : 3600;
  });

  const handleClear = () => {
    xtermInstance.current?.clear();
  };

  const handleCopyContainerId = useCallback(() => {
    navigator.clipboard.writeText(containerId);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }, [containerId]);

  // Real-time countdown timer
  useEffect(() => {
    const interval = setInterval(() => {
      setRemainingSeconds((prev) => {
        if (prev <= 1) {
          clearInterval(interval);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  const formatCountdown = (secs: number) => {
    const hours = Math.floor(secs / 3600);
    const mins = Math.floor((secs % 3600) / 60);
    const s = secs % 60;
    const pad = (n: number) => n.toString().padStart(2, '0');
    return `${pad(hours)}:${pad(mins)}:${pad(s)}`;
  };

  const isExpiringSoon = remainingSeconds <= 300;

  // Initialize xterm and WebSocket connection
  useEffect(() => {
    if (!terminalRef.current) return;

    // Warp / Raycast style modern dark terminal canvas (Always Dark)
    const term = new Terminal({
      cursorBlink: true,
      cursorStyle: 'bar',
      fontFamily: '"JetBrains Mono", "Fira Code", Consolas, monospace',
      fontSize: 13.5,
      lineHeight: 1.48,
      letterSpacing: 0.5,
      theme: {
        background: '#0d0e14',
        foreground: '#e4e4e7',
        cursor: '#10b981',
        selectionBackground: 'rgba(16, 185, 129, 0.25)',
        black: '#0d0e14',
        red: '#f87171',
        green: '#34d399',
        yellow: '#fbbf24',
        blue: '#60a5fa',
        magenta: '#c084fc',
        cyan: '#38bdf8',
        white: '#ffffff',
      },
    });

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.open(terminalRef.current);
    fitAddon.fit();

    xtermInstance.current = term;
    fitAddonRef.current = fitAddon;

    // Connect to real WebSocket backend
    const controller = connectTerminalSession(session.resourceUrn, session.token, {
      onData: (data) => {
        // Binary data from SSH → write to terminal
        const decoder = new TextDecoder();
        term.write(decoder.decode(data));
      },
      onError: () => {
        term.writeln('\r\n\x1b[1;31m[Connection Error] WebSocket connection failed.\x1b[0m');
        setWsConnected(false);
      },
      onClose: (code, reason) => {
        term.writeln(`\r\n\x1b[1;33m[Connection Closed] Code: ${code}, Reason: ${reason}\x1b[0m`);
        setWsConnected(false);
      },
    });

    sessionControllerRef.current = controller;
    setWsConnected(true);

    // Show connection banner
    term.writeln(`\x1b[1;36m[Cloud Terminal]\x1b[0m Connecting to ${session.resourceName}...`);
    term.writeln(`\x1b[90m> STS Token URN:   \x1b[37m${session.resourceUrn}\x1b[0m`);
    term.writeln(`\x1b[90m> Session ID:      \x1b[33m${session.sessionId}\x1b[0m`);
    term.writeln(`\x1b[32m✔ WebSocket Secure Tunnel Established.\x1b[0m`);
    term.writeln(`--------------------------------------------------------------------------------`);

    // Handle terminal input
    const disposable = term.onData((data) => {
      controller.send(data);
    });

    // Handle resize
    const resizeObserver = new ResizeObserver(() => {
      fitAddon.fit();
      const dims = fitAddon.proposeDimensions();
      if (dims) {
        controller.resize(dims.cols, dims.rows);
      }
    });

    if (terminalRef.current) {
      resizeObserver.observe(terminalRef.current);
    }

    return () => {
      disposable.dispose();
      resizeObserver.disconnect();
      controller.close();
      term.dispose();
    };
  }, [session]);

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-[#09090b] py-4 space-y-3 selection:bg-zinc-200 dark:selection:bg-white/20">
      {/* ... existing JSX ... */}
      <div className="max-w-[1600px] mx-auto">

        {/* Top Banner with Session Info */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 mb-4">

          {/* Left: Session Info Panel */}
          <div className="lg:col-span-4 xl:col-span-3 space-y-3">

            {/* Return Link */}
            <button
              onClick={onDisconnect}
              className="inline-flex items-center space-x-1.5 px-3 py-1.5 rounded-lg bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/10 text-zinc-700 dark:text-zinc-300 text-xs hover:bg-zinc-100 dark:hover:bg-white/5 transition-colors"
            >
              <ArrowLeft className="w-3.5 h-3.5" />
              <span>{t('returnToPortal')}</span>
            </button>

            {/* Session Details Card */}
            <div className="p-4 rounded-xl bg-white dark:bg-[#121318] border border-zinc-200 dark:border-white/10 space-y-4 shadow-sm">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-semibold text-zinc-900 dark:text-white">{t('sandboxInfo')}</h3>
                <span className={`inline-flex items-center space-x-1 px-2 py-0.5 rounded-full text-[10px] font-mono ${
                  wsConnected
                    ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/30'
                    : 'bg-zinc-500/10 text-zinc-500 border border-zinc-500/30'
                }`}>
                  <span className={`w-1.5 h-1.5 rounded-full ${wsConnected ? 'bg-emerald-500' : 'bg-zinc-500'}`} />
                  <span>{wsConnected ? t('wsConnected') : 'Disconnected'}</span>
                </span>
              </div>

              <div className="space-y-2.5 text-xs">
                {/* Target Resource */}
                <div className="flex items-center space-x-2.5">
                  <Server className="w-3.5 h-3.5 text-zinc-400 shrink-0" />
                  <div>
                    <p className="text-[10px] text-zinc-500">{t('targetResourceLabel')}</p>
                    <p className="font-semibold text-zinc-900 dark:text-white">{session.resourceName}</p>
                  </div>
                </div>

                {/* STS Token URN */}
                <div className="flex items-center space-x-2.5">
                  <Key className="w-3.5 h-3.5 text-amber-500 shrink-0" />
                  <div className="min-w-0">
                    <p className="text-[10px] text-zinc-500">{t('stsTokenUrn')}</p>
                    <p className="font-mono text-zinc-600 dark:text-zinc-400 truncate">{session.resourceUrn}</p>
                  </div>
                </div>

                {/* Session ID */}
                <div className="flex items-center space-x-2.5">
                  <Cpu className="w-3.5 h-3.5 text-zinc-400 shrink-0" />
                  <div>
                    <p className="text-[10px] text-zinc-500">{t('sessionCountdown')}</p>
                    <p className="font-mono font-bold text-zinc-900 dark:text-white">
                      {formatCountdown(remainingSeconds)}
                    </p>
                  </div>
                </div>

                {/* Copy Container ID */}
                <div className="flex items-center justify-between p-2.5 rounded-xl bg-zinc-50 dark:bg-white/[0.03] border border-zinc-200 dark:border-white/[0.08]">
                  <div className="flex items-center space-x-2.5">
                    <Box className="w-3.5 h-3.5 text-zinc-400 shrink-0" />
                    <div>
                      <p className="text-[10px] text-zinc-500">{t('containerId')}</p>
                      <p className="font-mono text-xs text-zinc-700 dark:text-zinc-300">{containerId}</p>
                    </div>
                  </div>
                  <button
                    onClick={handleCopyContainerId}
                    className="p-1.5 rounded-lg hover:bg-zinc-200 dark:hover:bg-white/10 text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors"
                    title={t('copyContainerId')}
                  >
                    {copied ? <Check className="w-3.5 h-3.5 text-emerald-500" /> : <Copy className="w-3.5 h-3.5" />}
                  </button>
                </div>
              </div>

              {/* Isolation Status */}
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-[10px] font-mono font-semibold text-zinc-500 uppercase">
                    {t('isolationStatus')}
                  </span>
                  <div className="w-2 h-2 rounded-full bg-sky-500" />
                </div>

                <div className="isolation-badge p-2.5 rounded-xl bg-zinc-50 dark:bg-white/[0.03] border border-zinc-200 dark:border-white/[0.08] flex items-center justify-between">
                  <div className="flex items-center space-x-2.5">
                    <ShieldAlert className="w-3.5 h-3.5 text-amber-600 dark:text-amber-400" />
                    <span className="text-xs font-medium text-zinc-800 dark:text-zinc-200">
                      {t('autoDestructOnExit')}
                    </span>
                  </div>
                  <span className="text-[10px] font-mono text-amber-600 dark:text-amber-400 font-bold">
                    rm -f
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Right Terminal Window (Always Dark macOS/Warp Style) */}
          <div className="lg:col-span-8 xl:col-span-9 rounded-2xl bg-[#0d0e14] border border-zinc-300 dark:border-white/10 shadow-2xl overflow-hidden flex flex-col h-[72vh] lg:h-auto">
            
            {/* Window Title Bar with macOS Traffic Light Buttons */}
            <div className="px-4 py-3 bg-[#13141f] border-b border-white/[0.08] flex items-center justify-between shrink-0 select-none">
              {/* Traffic Light Buttons */}
              <div className="flex items-center space-x-2">
                <button 
                  onClick={onDisconnect} 
                  title="Close session & destroy container" 
                  className="w-3 h-3 rounded-full bg-red-500/80 traffic-dot-close transition-all cursor-pointer" 
                />
                <button 
                  title="Minimize / Background" 
                  className="w-3 h-3 rounded-full bg-amber-500/80 traffic-dot-minimize transition-all cursor-pointer" 
                />
                <button 
                  onClick={() => fitAddonRef.current?.fit()}
                  title="Fit viewport size" 
                  className="w-3 h-3 rounded-full bg-emerald-500/80 traffic-dot-expand transition-all cursor-pointer" 
                />
              </div>

              {/* Window Title */}
              <div className="flex items-center space-x-2 text-xs font-mono text-zinc-300">
                <Box className="w-3.5 h-3.5 text-sky-400" />
                <span className="font-semibold text-white">{session.resourceName}</span>
                <span className="text-zinc-500">•</span>
                <span className="text-zinc-400">{containerId}</span>
              </div>

              {/* Window Quick Controls */}
              <div className="flex items-center space-x-2">
                <button
                  onClick={handleClear}
                  title="Clear Terminal Output"
                  className="p-1 rounded text-zinc-400 hover:text-white hover:bg-white/10 transition-colors"
                >
                  <RefreshCw className="w-3.5 h-3.5" />
                </button>
                <div className="flex items-center space-x-1.5 px-2 py-0.5 rounded bg-black/40 border border-white/10 text-[11px] font-mono text-emerald-400">
                  <Clock className="w-3 h-3" />
                  <span>{t('sessionValid')}</span>
                </div>
              </div>
            </div>

            {/* Terminal Canvas (Always Dark) */}
            <div className="flex-1 p-4 relative overflow-hidden bg-[#0d0e14]">
              <div ref={terminalRef} className="w-full h-full" />
            </div>

            {/* Footer Window Status Bar */}
            <div className="px-4 py-2 bg-[#0a0b10] border-t border-white/[0.06] flex items-center justify-between text-[11px] font-mono text-zinc-500 shrink-0">
              <div className="flex items-center space-x-4">
                <span className={`flex items-center space-x-1.5 ${wsConnected ? 'text-emerald-400' : 'text-zinc-500'}`}>
                  <Wifi className={`w-3 h-3 ${wsConnected ? 'animate-pulse' : ''}`} />
                  <span>{wsConnected ? t('wsConnected') : 'Disconnected'}</span>
                </span>
                <span className="hidden sm:inline text-zinc-400">STS: {session.token.substring(0, 16)}...</span>
                <span className="hidden md:inline text-zinc-500">{session.resourceUrn}</span>
              </div>
              <span className={`${isExpiringSoon ? 'text-amber-400' : 'text-zinc-500'}`}>
                {formatCountdown(remainingSeconds)}
              </span>
            </div>
          </div>

        </div>

        {/* Global Bottom Audit Banner */}
        <div className="terminal-footer-bar px-4 py-2.5 rounded-xl flex items-center justify-between text-xs font-mono text-zinc-600 dark:text-zinc-400 border border-zinc-200 dark:border-white/[0.08] shadow-sm">
          <div className="flex items-center space-x-2">
            <ShieldCheck className="w-4 h-4 text-emerald-600 dark:text-emerald-400 shrink-0" />
            <span className="text-zinc-800 dark:text-zinc-200 font-semibold">{t('auditActive')}</span>
          </div>
          <span className="hidden sm:inline text-zinc-500 text-[11px]">
            Session ID: {session.sessionId}
          </span>
        </div>
      </div>
    </div>
  );
};
