import React, { useEffect, useState } from 'react';
import { fetchApi } from '../../api/client';
import { UsageMeter, Invoice } from '../../types';
import {
  CreditCard,
  Receipt,
  Ticket,
  Clock,
  CheckCircle2,
  AlertTriangle,
  ExternalLink,
  ShieldAlert,
  Zap,
  TrendingUp,
  ListFilter,
  Eye,
  X,
} from 'lucide-react';

import { useToast } from '../../components/Toast';

export const BillingPortalPage: React.FC = () => {
  const { showError, showSuccess, showInfo } = useToast();
  const [usage, setUsage] = useState<UsageMeter | null>(null);

  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [selectedInvoice, setSelectedInvoice] = useState<Invoice | null>(null);
  const [loading, setLoading] = useState(true);
  const [payLoadingId, setPayLoadingId] = useState<number | null>(null);

  useEffect(() => {
    fetchBillingData();
  }, []);

  const fetchBillingData = async () => {
    setLoading(true);
    try {
      const [uRes, iRes] = await Promise.all([
        fetchApi<UsageMeter>('/billing/usage'),
        fetchApi<Invoice[]>('/billing/invoices'),
      ]);
      setUsage(uRes.data);
      setInvoices(iRes.data || []);
    } catch (err) {
      console.error('Failed to load billing data:', err);
    } finally {
      setLoading(false);
    }
  };

  const handlePayInvoice = async (inv: Invoice) => {
    setPayLoadingId(inv.id);
    try {
      const res = await fetchApi<{ payment_number: string; snap_token: string; snap_redirect_url: string }>(
        `/billing/invoices/${inv.id}/pay`,
        { method: 'POST' }
      );

      const redirectUrl = res.data.snap_redirect_url;
      if (redirectUrl) {
        showSuccess('Link pembayaran Midtrans Snap berhasil dibuat');
        window.open(redirectUrl, '_blank');
      } else {
        showInfo('Token pembayaran dibuat. Silakan selesaikan pembayaran.');
      }
    } catch (err: any) {
      showError(err.message || 'Gagal menghasilkan token pembayaran');
    } finally {
      setPayLoadingId(null);
    }

  };

  const currentMonthTickets = usage?.ticket_count || 0;
  const currentEstCost = currentMonthTickets * 500; // Rp 500 / ticket

  return (
    <div className="space-y-6">
      {/* Page Title */}
      <div className="flex items-center space-x-3">
        <div className="w-12 h-12 rounded-2xl bg-gradient-to-tr from-emerald-600 to-teal-600 text-white flex items-center justify-center shadow-lg shadow-emerald-600/20">
          <CreditCard className="w-6 h-6" />
        </div>
        <div>
          <h1 className="text-2xl font-black text-white">Billing & Subscription Portal</h1>
          <p className="text-xs text-slate-400">Postpaid Metered Usage (Rp 500 / tiket) & Midtrans Payment Gateway</p>
        </div>
      </div>

      {/* Top Cards Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Metered Usage Box */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-3xl p-6 shadow-xl relative overflow-hidden">
          <div className="flex justify-between items-center mb-4">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-widest">
              Tiket Bulan Ini ({usage?.billing_period || '2026-08'})
            </span>
            <Ticket className="w-5 h-5 text-emerald-400" />
          </div>
          <div className="text-4xl font-black text-white tracking-tight">{currentMonthTickets.toLocaleString()}</div>
          <p className="text-xs text-slate-400 mt-2 flex items-center space-x-1">
            <TrendingUp className="w-3.5 h-3.5 text-emerald-400" />
            <span>Akumulasi otomatis per pemanggilan tiket</span>
          </p>
        </div>

        {/* Estimated Current Bill */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-3xl p-6 shadow-xl relative overflow-hidden">
          <div className="flex justify-between items-center mb-4">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-widest">Estimasi Tagihan</span>
            <Zap className="w-5 h-5 text-amber-400" />
          </div>
          <div className="text-4xl font-black text-amber-400 tracking-tight">
            Rp {currentEstCost.toLocaleString('id-ID')}
          </div>
          <p className="text-xs text-slate-400 mt-2">Dihitung dari Rp 500 / tiket antrean</p>
        </div>

        {/* Current Plan Card */}
        <div className="bg-gradient-to-br from-slate-900 to-slate-900/90 border border-blue-500/30 rounded-3xl p-6 shadow-xl relative overflow-hidden">
          <div className="flex justify-between items-center mb-4">
            <span className="text-xs font-bold text-blue-400 uppercase tracking-widest">Paket Aktif</span>
            <span className="px-2.5 py-1 bg-blue-500/20 text-blue-400 rounded-full text-[10px] font-black border border-blue-500/30">
              POSTPAID METERED
            </span>
          </div>
          <div className="text-xl font-extrabold text-white">Postpaid Standard</div>
          <p className="text-xs text-slate-400 mt-2">Terpisah: Invoice Line Items & Payment Transactions</p>
        </div>
      </div>

      {/* Invoices List */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-3xl p-6 shadow-xl space-y-4">
        <div className="flex justify-between items-center border-b border-slate-800 pb-4">
          <div className="flex items-center space-x-2">
            <Receipt className="w-5 h-5 text-blue-400" />
            <h2 className="text-lg font-bold text-white">Riwayat Invoice Tagihan</h2>
          </div>
        </div>

        {loading ? (
          <div className="text-center py-8 text-slate-400 text-xs">Memuat data invoice...</div>
        ) : invoices.length === 0 ? (
          <div className="text-center py-12 text-slate-500 space-y-2">
            <Receipt className="w-12 h-12 mx-auto text-slate-700" />
            <p className="text-sm font-bold">Belum Ada Invoice Tagihan</p>
            <p className="text-xs text-slate-500">Invoice tagihan bulanan akan diterbitkan setiap tanggal 1 bulan baru.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs text-slate-300">
              <thead className="bg-slate-800/60 text-slate-400 uppercase text-[10px] font-bold border-b border-slate-800">
                <tr>
                  <th className="py-3 px-4">No. Invoice</th>
                  <th className="py-3 px-4">Periode</th>
                  <th className="py-3 px-4">Rincian Item</th>
                  <th className="py-3 px-4">Total Tagihan</th>
                  <th className="py-3 px-4">Jatuh Tempo</th>
                  <th className="py-3 px-4">Status</th>
                  <th className="py-3 px-4 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-800/60">
                {invoices.map((inv) => (
                  <tr key={inv.id} className="hover:bg-slate-800/30 transition">
                    <td className="py-3.5 px-4 font-mono font-bold text-white">{inv.invoice_number}</td>
                    <td className="py-3.5 px-4 font-semibold">{inv.billing_period}</td>
                    <td className="py-3.5 px-4">
                      <button
                        onClick={() => setSelectedInvoice(inv)}
                        className="text-blue-400 hover:text-blue-300 font-semibold flex items-center space-x-1"
                      >
                        <Eye className="w-3.5 h-3.5" />
                        <span>Lihat Item ({inv.items?.length || 0})</span>
                      </button>
                    </td>
                    <td className="py-3.5 px-4 font-bold text-emerald-400">
                      Rp {inv.total_amount.toLocaleString('id-ID')}
                    </td>
                    <td className="py-3.5 px-4 text-slate-400">
                      {new Date(inv.due_date).toLocaleDateString('id-ID')}
                    </td>
                    <td className="py-3.5 px-4">
                      <span
                        className={`px-2.5 py-1 rounded-full text-[10px] font-extrabold ${
                          inv.status === 'PAID'
                            ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
                            : inv.status === 'OVERDUE'
                            ? 'bg-red-500/20 text-red-400 border border-red-500/30'
                            : 'bg-amber-500/20 text-amber-400 border border-amber-500/30'
                        }`}
                      >
                        {inv.status}
                      </span>
                    </td>
                    <td className="py-3.5 px-4 text-right">
                      {inv.status !== 'PAID' ? (
                        <button
                          onClick={() => handlePayInvoice(inv)}
                          disabled={payLoadingId === inv.id}
                          className="px-3.5 py-1.5 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white font-bold text-[11px] rounded-xl shadow-md transition flex items-center space-x-1.5 ml-auto"
                        >
                          <CreditCard className="w-3.5 h-3.5" />
                          <span>{payLoadingId === inv.id ? 'Memuat Midtrans...' : 'Bayar Sekarang'}</span>
                        </button>
                      ) : (
                        <span className="text-emerald-400 font-bold text-[11px] flex items-center justify-end space-x-1">
                          <CheckCircle2 className="w-3.5 h-3.5" />
                          <span>Lunas</span>
                        </span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Detail Invoice Items Modal */}
      {selectedInvoice && (
        <div className="fixed inset-0 bg-slate-950/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-3xl p-6 max-w-lg w-full space-y-4 shadow-2xl relative">
            <div className="flex justify-between items-center border-b border-slate-800 pb-3">
              <div>
                <h3 className="text-base font-bold text-white">Rincian Invoice Items</h3>
                <p className="text-xs text-slate-400 font-mono">{selectedInvoice.invoice_number}</p>
              </div>
              <button
                onClick={() => setSelectedInvoice(null)}
                className="p-2 text-slate-400 hover:text-white rounded-xl hover:bg-slate-800 transition"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-2">
              {selectedInvoice.items && selectedInvoice.items.length > 0 ? (
                selectedInvoice.items.map((item) => (
                  <div key={item.id} className="bg-slate-800/50 p-3 rounded-2xl border border-slate-700/60 flex justify-between items-center text-xs">
                    <div>
                      <div className="font-bold text-white">{item.description}</div>
                      <div className="text-[11px] text-slate-400">
                        {item.quantity} x Rp {item.unit_price.toLocaleString('id-ID')}
                      </div>
                    </div>
                    <div className="font-mono font-bold text-emerald-400">
                      Rp {item.total_price.toLocaleString('id-ID')}
                    </div>
                  </div>
                ))
              ) : (
                <div className="text-center py-4 text-slate-400 text-xs">Tidak ada rincian item</div>
              )}
            </div>

            <div className="pt-3 border-t border-slate-800 flex justify-between items-center text-xs font-bold">
              <span className="text-slate-400">Total Tagihan:</span>
              <span className="text-emerald-400 text-sm">Rp {selectedInvoice.total_amount.toLocaleString('id-ID')}</span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
