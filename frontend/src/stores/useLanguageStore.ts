import { create } from 'zustand';
import { translations, type Language } from '../i18n/translations';

interface LanguageState {
  lang: Language;
  setLanguage: (lang: Language) => void;
  toggleLanguage: () => void;
  t: (key: keyof typeof translations['en']) => string;
}

export const useLanguageStore = create<LanguageState>((set, get) => ({
  lang: 'zh', // Default to Chinese as requested!
  setLanguage: (lang) => set({ lang }),
  toggleLanguage: () => set((state) => ({ lang: state.lang === 'zh' ? 'en' : 'zh' })),
  t: (key) => {
    const { lang } = get();
    return translations[lang]?.[key] || translations['en'][key] || String(key);
  },
}));
