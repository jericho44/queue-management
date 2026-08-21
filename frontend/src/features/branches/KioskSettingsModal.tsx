import React, { useState } from 'react';
import { Branch } from '../../types';
import { fetchApi } from '../../api/client';
import { X, Printer, Settings2, Sliders, Check } from 'lucide-react';

import { useToast } from '../../components/Toast';

interface KioskSettingsModalProps {
  branch: Branch;
  onClose: () => void;
  onSaveSuccess: () => void;
}

export const KioskSettingsModal: React.FC<KioskSettingsModalProps> = ({
  branch,
  onClose,
  onSaveSuccess,
}) => {
  const { showError, showSuccess } = useToast();
  const [kioskEnabled, setKioskEnabled] = useState(branch.kiosk_enabled ?? true);
  const [kioskMode, setKioskMode] = useState<'DUAL' | 'PAPERLESS' | 'PHYSICAL'>(
    branch.kiosk_mode || 'DUAL'
  );
  const [paperSize, setPaperSize] = useState<'58mm' | '80mm'>(branch.paper_size || '58mm');
  const [receiptHeader, setReceiptHeader] = useState(branch.receipt_header || '');
  const [receiptFooter, setReceiptFooter] = useState(
    branch.receipt_footer || 'Terima kasih atas kunjungan Anda. Harap menunggu hingga dipanggil.'
  );
  const [autoPrint, setAutoPrint] = useState(branch.auto_print ?? false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await fetchApi(`/branches/${branch.id}/kiosk-settings`, {
        method: 'PUT',
        body: JSON.stringify({
          kiosk_enabled: kioskEnabled,
          kiosk_mode: kioskMode,
          paper_size: paperSize,
          receipt_header: receiptHeader,
          receipt_footer: receiptFooter,
          auto_print: autoPrint,
        }),
      });
      showSuccess('Pengaturan Kiosk & Thermal Printer berhasil disimpan');
      onSaveSuccess();
      onClose();
    } catch (err: any) {
      showError(err.message || 'Gagal memperbarui pengaturan Kiosk & Thermal Printer');
    } finally {
      setLoading(false);
    }
  };


  return (
    <div className="fixed inset-0 bg-slate-950/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-slate-900 border border-slate-800 rounded-3xl p-6 max-w-xl w-full space-y-6 shadow-2xl relative">
        <div className="flex justify-between items-center border-b border-slate-800 pb-4">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 rounded-xl bg-blue-600/20 text-blue-400 border border-blue-500/30 flex items-center justify-center">
              <Printer className="w-5 h-5" />
            </div>
            <div>
              <h2 className="text-lg font-extrabold text-white">Pengaturan Kiosk & Thermal Printer</h2>
              <p className="text-xs text-slate-400">Cabang: {branch.name} ({branch.code})</p>
            </div>
          </div>

          <button
            onClick={onClose}
            className="p-2 text-slate-400 hover:text-white rounded-xl hover:bg-slate-800 transition"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Status Enable Kiosk */}
          <div className="flex items-center justify-between p-4 bg-slate-800/60 rounded-2xl border border-slate-700/70">
            <div>
              <div className="text-xs font-bold text-white">Aktifkan Terminal Kiosk</div>
              <div className="text-[11px] text-slate-400">Izinkan publik mengakses URL Kiosk cabang ini (`/kiosk/${branch.code.toLowerCase()}`)</div>

            </div>
            <input
              type="checkbox"
              checked={kioskEnabled}
              onChange={(e) => setKioskEnabled(e.target.checked)}
              className="w-5 h-5 rounded bg-slate-900 border-slate-700 text-blue-600 focus:ring-blue-500"
            />
          </div>

          {/* Mode Kiosk */}
          <div>
            <label className="block text-xs font-bold text-slate-300 mb-1.5">Mode Cetak Kiosk</label>
            <div className="grid grid-cols-3 gap-3">
              {[
                { id: 'DUAL', label: 'Dual Mode', desc: 'Bisa pilih Cetak / Digital' },
                { id: 'PHYSICAL', label: 'Cetak Fisik Only', desc: 'Struk Thermal' },
                { id: 'PAPERLESS', label: 'Paperless Only', desc: 'Scan QR Code HP' },
              ].map((m) => (
                <button
                  key={m.id}
                  type="button"
                  onClick={() => setKioskMode(m.id as any)}
                  className={`p-3 rounded-2xl border text-left transition ${
                    kioskMode === m.id
                      ? 'bg-blue-600/20 border-blue-500 text-white font-bold'
                      : 'bg-slate-800/60 border-slate-700/60 text-slate-400 hover:bg-slate-800'
                  }`}
                >
                  <div className="text-xs">{m.label}</div>
                  <div className="text-[10px] text-slate-400 mt-0.5">{m.desc}</div>
                </button>
              ))}
            </div>
          </div>

          {/* Ukuran Kertas */}
          <div>
            <label className="block text-xs font-bold text-slate-300 mb-1.5">Ukuran Kertas Thermal Printer</label>
            <div className="grid grid-cols-2 gap-3">
              {[
                { id: '58mm', label: '58 mm (Standard Roll)' },
                { id: '80mm', label: '80 mm (Wide POS Roll)' },
              ].map((p) => (
                <button
                  key={p.id}
                  type="button"
                  onClick={() => setPaperSize(p.id as any)}
                  className={`p-3 rounded-2xl border text-center text-xs font-bold transition ${
                    paperSize === p.id
                      ? 'bg-blue-600/20 border-blue-500 text-blue-400'
                      : 'bg-slate-800/60 border-slate-700/60 text-slate-400 hover:bg-slate-800'
                  }`}
                >
                  {p.label}
                </button>
              ))}
            </div>
          </div>

          {/* Header & Footer Custom */}
          <div className="space-y-3">
            <div>
              <label className="block text-xs font-bold text-slate-300 mb-1">Header Struk Custom (Opsional)</label>
              <input
                type="text"
                value={receiptHeader}
                onChange={(e) => setReceiptHeader(e.target.value)}
                placeholder="Contoh: Selamat Datang di Klinik Utama"
                className="w-full bg-slate-800/60 border border-slate-700 rounded-xl px-3.5 py-2 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-blue-500"
              />
            </div>

            <div>
              <label className="block text-xs font-bold text-slate-300 mb-1">Pesan Footer Struk</label>
              <textarea
                rows={2}
                value={receiptFooter}
                onChange={(e) => setReceiptFooter(e.target.value)}
                placeholder="Pesan ucapan di bagian bawah struk..."
                className="w-full bg-slate-800/60 border border-slate-700 rounded-xl px-3.5 py-2 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-blue-500"
              />
            </div>
          </div>

          <div className="pt-3 border-t border-slate-800 flex justify-end space-x-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2.5 bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-bold rounded-xl border border-slate-700 transition"
            >
              Batal
            </button>
            <button
              type="submit"
              disabled={loading}
              className="px-5 py-2.5 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-xl shadow-lg shadow-blue-600/30 transition flex items-center space-x-2"
            >
              <Check className="w-4 h-4" />
              <span>{loading ? 'Menyimpan...' : 'Simpan Pengaturan'}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
