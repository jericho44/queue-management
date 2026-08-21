import React, { useEffect, useState } from 'react';
import { fetchApi } from '../../api/client';
import { Invoice, SuperadminBillingStats } from '../../types';
import {
  DollarSign,
  TrendingUp,
  Receipt,
  AlertTriangle,
  CheckCircle2,
  Clock,
  ShieldCheck,
  Building,
  RefreshCw,
} from 'lucide-react';

export const SuperadminBillingPage: React.FC = () => {
  const [stats, setStats] = useState<SuperadminBillingStats | null>(null);
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchSuperadminBilling();
  }, []);

  const fetchSuperadminBilling = async () => {
    setLoading(true);
    try {
      const [sRes, iRes] = await Promise.all([
        fetchApi<SuperadminBillingStats>('/superadmin/billing/stats'),
        fetchApi<Invoice[]>('/superadmin/billing/invoices'),
      ]);
      setStats(sRes.data);
      setInvoices(iRes.data || []);
    } catch (err) {
      console.error('Failed to load superadmin billing data:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Title */}
      <div className="flex items-center space-x-3">
        <div className="w-12 h-12 rounded-2xl bg-gradient-to-tr from-blue-600 to-indigo-600 text-white flex items-center justify-center shadow-lg shadow-blue-600/25">
          <DollarSign className="w-6 h-6" />
        </div>
        <div>
          <h1 className="text-2xl font-black text-white">Superadmin Billing & SaaS Revenue Control</h1>
          <p className="text-xs text-slate-400">Postpaid Metered Revenue, Invoice Tracking, & Midtrans Payment Logs</p>
        </div>
      </div>

      {/* Metric Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="bg-slate-900/90 border border-slate-800 rounded-3xl p-6 shadow-xl">
          <div className="flex justify-between items-center mb-3">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-widest">Total Pendapatan</span>
            <DollarSign className="w-5 h-5 text-emerald-400" />
          </div>
          <div className="text-3xl font-black text-emerald-400 tracking-tight">
            Rp {(stats?.total_revenue || 0).toLocaleString('id-ID')}
          </div>
          <p className="text-xs text-slate-400 mt-2 font-semibold">{stats?.paid_invoices_count || 0} Invoice Lunas</p>
        </div>

        <div className="bg-slate-900/90 border border-slate-800 rounded-3xl p-6 shadow-xl">
          <div className="flex justify-between items-center mb-3">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-widest">Pending Revenue</span>
            <Clock className="w-5 h-5 text-amber-400" />
          </div>
          <div className="text-3xl font-black text-amber-400 tracking-tight">
            Rp {(stats?.pending_revenue || 0).toLocaleString('id-ID')}
          </div>
          <p className="text-xs text-slate-400 mt-2 font-semibold">
            {stats?.unpaid_invoices_count || 0} Unpaid • {stats?.overdue_invoices_count || 0} Overdue
          </p>
        </div>

        <div className="bg-slate-900/90 border border-slate-800 rounded-3xl p-6 shadow-xl">
          <div className="flex justify-between items-center mb-3">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-widest">Total Tiket Sistem</span>
            <TrendingUp className="w-5 h-5 text-blue-400" />
          </div>
          <div className="text-3xl font-black text-white tracking-tight">
            {(stats?.current_month_tickets || 0).toLocaleString()}
          </div>
          <p className="text-xs text-slate-400 mt-2">Bulan Berjalan ({new Date().toISOString().slice(0, 7)})</p>
        </div>

        <div className="bg-slate-900/90 border border-slate-800 rounded-3xl p-6 shadow-xl">
          <div className="flex justify-between items-center mb-3">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-widest">Tarif Default</span>
            <ShieldCheck className="w-5 h-5 text-indigo-400" />
          </div>
          <div className="text-3xl font-black text-indigo-400 tracking-tight">Rp 500</div>
          <p className="text-xs text-slate-400 mt-2">Per Tiket Antrean (Postpaid)</p>
        </div>
      </div>

      {/* Invoices List for Superadmin */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-3xl p-6 shadow-xl space-y-4">
        <div className="flex justify-between items-center border-b border-slate-800 pb-4">
          <div className="flex items-center space-x-2">
            <Receipt className="w-5 h-5 text-blue-400" />
            <h2 className="text-lg font-bold text-white">Semua Invoice Tagihan Tenant ({invoices.length})</h2>
          </div>

          <button
            onClick={fetchSuperadminBilling}
            className="p-2 text-slate-400 hover:text-white rounded-xl hover:bg-slate-800 transition"
          >
            <RefreshCw className="w-4 h-4" />
          </button>
        </div>

        {loading ? (
          <div className="text-center py-8 text-slate-400 text-xs">Memuat data invoice...</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs text-slate-300">
              <thead className="bg-slate-800/60 text-slate-400 uppercase text-[10px] font-bold border-b border-slate-800">
                <tr>
                  <th className="py-3 px-4">No. Invoice</th>
                  <th className="py-3 px-4">Nama Organisasi</th>
                  <th className="py-3 px-4">Periode</th>
                  <th className="py-3 px-4">Jumlah Tiket</th>
                  <th className="py-3 px-4">Total Tagihan</th>
                  <th className="py-3 px-4">Status</th>
                  <th className="py-3 px-4">Metode Bayar</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-800/60">
                {invoices.map((inv) => (
                  <tr key={inv.id} className="hover:bg-slate-800/30 transition">
                    <td className="py-3.5 px-4 font-mono font-bold text-white">{inv.invoice_number}</td>
                    <td className="py-3.5 px-4 font-semibold text-blue-400">{inv.org_name || `Tenant #${inv.organization_id}`}</td>
                    <td className="py-3.5 px-4 font-semibold">{inv.billing_period}</td>
                    <td className="py-3.5 px-4">
                      {(inv.items?.find((i) => i.item_type === 'TICKET_USAGE')?.quantity || 0).toLocaleString()} Tiket
                    </td>
                    <td className="py-3.5 px-4 font-bold text-emerald-400">
                      Rp {inv.total_amount.toLocaleString('id-ID')}
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
                    <td className="py-3.5 px-4 text-slate-400">
                      {inv.latest_payment?.payment_method ? inv.latest_payment.payment_method.toUpperCase() : '-'}
                    </td>

                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};
