import React, { createContext, useContext, useState, useCallback } from 'react';
import { AlertCircle, CheckCircle2, Info, AlertTriangle, X } from 'lucide-react';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface ToastItem {
  id: string;
  type: ToastType;
  title?: string;
  message: string;
}

interface ToastContextType {
  toast: (message: string, type?: ToastType, title?: string) => void;
  showError: (message: string, title?: string) => void;
  showSuccess: (message: string, title?: string) => void;
  showInfo: (message: string, title?: string) => void;
  showWarning: (message: string, title?: string) => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export const ToastProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [toasts, setToasts] = useState<ToastItem[]>([]);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const toast = useCallback((message: string, type: ToastType = 'info', title?: string) => {
    const id = Math.random().toString(36).substring(2, 9);
    setToasts((prev) => [...prev, { id, type, title, message }]);

    setTimeout(() => {
      removeToast(id);
    }, 5000);
  }, [removeToast]);

  const showError = useCallback((message: string, title: string = 'Terjadi Kesalahan') => {
    toast(message, 'error', title);
  }, [toast]);

  const showSuccess = useCallback((message: string, title: string = 'Berhasil') => {
    toast(message, 'success', title);
  }, [toast]);

  const showInfo = useCallback((message: string, title: string = 'Informasi') => {
    toast(message, 'info', title);
  }, [toast]);

  const showWarning = useCallback((message: string, title: string = 'Peringatan') => {
    toast(message, 'warning', title);
  }, [toast]);

  return (
    <ToastContext.Provider value={{ toast, showError, showSuccess, showInfo, showWarning }}>
      {children}
      {/* Toast Render Container */}
      <div className="fixed bottom-5 right-5 z-[9999] flex flex-col space-y-3 max-w-md w-full pointer-events-none p-4">
        {toasts.map((t) => (
          <div
            key={t.id}
            className={`pointer-events-auto flex items-start space-x-3 p-4 rounded-2xl border shadow-2xl backdrop-blur-xl transition-all duration-300 transform translate-y-0 animate-in fade-in slide-in-from-bottom-5 ${
              t.type === 'error'
                ? 'bg-slate-900/95 border-red-500/40 text-red-200 shadow-red-950/40'
                : t.type === 'success'
                ? 'bg-slate-900/95 border-emerald-500/40 text-emerald-200 shadow-emerald-950/40'
                : t.type === 'warning'
                ? 'bg-slate-900/95 border-amber-500/40 text-amber-200 shadow-amber-950/40'
                : 'bg-slate-900/95 border-blue-500/40 text-blue-200 shadow-blue-950/40'
            }`}
          >
            {/* Icon */}
            <div className="flex-shrink-0 mt-0.5">
              {t.type === 'error' && <AlertCircle className="w-5 h-5 text-red-400" />}
              {t.type === 'success' && <CheckCircle2 className="w-5 h-5 text-emerald-400" />}
              {t.type === 'warning' && <AlertTriangle className="w-5 h-5 text-amber-400" />}
              {t.type === 'info' && <Info className="w-5 h-5 text-blue-400" />}
            </div>

            {/* Content */}
            <div className="flex-1 min-w-0">
              {t.title && <h4 className="text-xs font-black tracking-wide uppercase mb-0.5">{t.title}</h4>}
              <p className="text-xs font-medium leading-relaxed opacity-90 break-words">{t.message}</p>
            </div>

            {/* Close Button */}
            <button
              onClick={() => removeToast(t.id)}
              className="flex-shrink-0 p-1 rounded-lg opacity-60 hover:opacity-100 hover:bg-slate-800 transition"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
};

export const useToast = (): ToastContextType => {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context;
};
