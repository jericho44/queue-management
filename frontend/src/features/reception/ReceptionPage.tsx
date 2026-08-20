import React, { useEffect, useState } from 'react';
import { fetchApi } from '../../api/client';
import { Branch, Service, QueueTicket } from '../../types';
import { TicketPlus, Printer, QrCode, CheckCircle2, User, Phone, Sparkles } from 'lucide-react';
import { Link } from 'react-router-dom';

export const ReceptionPage: React.FC = () => {
  const [branches, setBranches] = useState<Branch[]>([]);
  const [services, setServices] = useState<Service[]>([]);
  const [selectedBranchId, setSelectedBranchId] = useState<number>(0);
  const [selectedServiceId, setSelectedServiceId] = useState<number>(0);
  const [priority, setPriority] = useState<'NORMAL' | 'PRIORITY' | 'EMERGENCY'>('NORMAL');
  const [customerName, setCustomerName] = useState('');
  const [customerPhone, setCustomerPhone] = useState('');
  const [issuedTicket, setIssuedTicket] = useState<QueueTicket | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchBranches();
  }, []);

  useEffect(() => {
    if (selectedBranchId > 0) {
      fetchServices(selectedBranchId);
    }
  }, [selectedBranchId]);

  const fetchBranches = async () => {
    try {
      const res = await fetchApi<Branch[]>('/branches');
      if (res.data && res.data.length > 0) {
        setBranches(res.data);
        setSelectedBranchId(res.data[0].id);
      }
    } catch (err) {
      console.error(err);
    }
  };

  const fetchServices = async (bId: number) => {
    try {
      const res = await fetchApi<Service[]>(`/services?branch_id=${bId}`);
      setServices(res.data || []);
      if (res.data && res.data.length > 0) {
        setSelectedServiceId(res.data[0].id);
      }
    } catch (err) {
      console.error(err);
    }
  };

  const handleIssueTicket = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedBranchId || !selectedServiceId) return;

    setLoading(true);
    try {
      const res = await fetchApi<QueueTicket>('/tickets', {
        method: 'POST',
        body: JSON.stringify({
          branch_id: selectedBranchId,
          service_id: selectedServiceId,
          priority: priority,
          name: customerName,
          phone: customerPhone,
        }),
      });

      setIssuedTicket(res.data);
      setCustomerName('');
      setCustomerPhone('');
    } catch (err: any) {
      alert(err.message || 'Failed to issue ticket');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div className="flex items-center space-x-3">
        <div className="w-10 h-10 rounded-xl bg-blue-600/20 text-blue-400 border border-blue-500/30 flex items-center justify-center">
          <TicketPlus className="w-5 h-5" />
        </div>
        <div>
          <h1 className="text-2xl font-black text-white">Ticket Issuance Terminal</h1>
          <p className="text-xs text-slate-400">Issue customer queue tickets and digital receipts</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Left Column: Issue Ticket Form */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-2xl p-6 shadow-xl">
          <h2 className="text-sm font-bold text-white mb-4">Select Service & Details</h2>

          <form onSubmit={handleIssueTicket} className="space-y-4">
            <div>
              <label className="block text-xs font-semibold text-slate-300 mb-1">Target Branch</label>
              <select
                value={selectedBranchId}
                onChange={(e) => setSelectedBranchId(Number(e.target.value))}
                className="w-full bg-slate-800 border border-slate-700 rounded-xl px-3.5 py-2.5 text-xs text-white focus:outline-none focus:border-blue-500"
              >
                {branches.map((b) => (
                  <option key={b.id} value={b.id}>
                    {b.name} ({b.code})
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-xs font-semibold text-slate-300 mb-1">Service Type</label>
              <div className="grid grid-cols-1 gap-2">
                {services.map((s) => (
                  <button
                    key={s.id}
                    type="button"
                    onClick={() => setSelectedServiceId(s.id)}
                    className={`p-3 rounded-xl border text-left flex items-center justify-between transition ${
                      selectedServiceId === s.id
                        ? 'bg-blue-600/15 border-blue-500 text-white font-semibold'
                        : 'bg-slate-800/60 border-slate-700/70 text-slate-300 hover:bg-slate-800'
                    }`}
                  >
                    <div>
                      <div className="text-xs font-bold">{s.name}</div>
                      <div className="text-[10px] text-slate-400">Prefix [{s.prefix}] • ~{Math.floor(s.avg_service_time_sec / 60)} min</div>
                    </div>
                    <span className="text-xs font-extrabold text-blue-400">Prefix {s.prefix}</span>
                  </button>
                ))}
              </div>
            </div>

            <div>
              <label className="block text-xs font-semibold text-slate-300 mb-1">Queue Priority</label>
              <div className="grid grid-cols-3 gap-2">
                {(['NORMAL', 'PRIORITY', 'EMERGENCY'] as const).map((p) => (
                  <button
                    key={p}
                    type="button"
                    onClick={() => setPriority(p)}
                    className={`py-2 rounded-xl text-xs font-bold border transition ${
                      priority === p
                        ? p === 'EMERGENCY'
                          ? 'bg-red-500/20 border-red-500 text-red-400'
                          : p === 'PRIORITY'
                          ? 'bg-amber-500/20 border-amber-500 text-amber-400'
                          : 'bg-blue-600/20 border-blue-500 text-blue-400'
                        : 'bg-slate-800/60 border-slate-700/70 text-slate-400'
                    }`}
                  >
                    {p}
                  </button>
                ))}
              </div>
            </div>

            <div className="pt-2 border-t border-slate-800 space-y-3">
              <div>
                <label className="block text-xs font-semibold text-slate-300 mb-1">Customer Name (Optional)</label>
                <input
                  type="text"
                  value={customerName}
                  onChange={(e) => setCustomerName(e.target.value)}
                  placeholder="John Doe"
                  className="w-full bg-slate-800/60 border border-slate-700 rounded-xl px-3.5 py-2 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-blue-500"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 mb-1">Phone Number (Optional)</label>
                <input
                  type="text"
                  value={customerPhone}
                  onChange={(e) => setCustomerPhone(e.target.value)}
                  placeholder="+62 812-3456-7890"
                  className="w-full bg-slate-800/60 border border-slate-700 rounded-xl px-3.5 py-2 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-blue-500"
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full py-3.5 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white font-extrabold text-xs rounded-xl shadow-lg shadow-blue-600/25 transition mt-4"
            >
              {loading ? 'Generating Ticket...' : 'PRINT & ISSUE TICKET'}
            </button>
          </form>
        </div>

        {/* Right Column: Issued Ticket Preview */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-2xl p-6 flex flex-col items-center justify-center text-center shadow-xl">
          {issuedTicket ? (
            <div className="w-full max-w-sm bg-slate-950 border border-slate-800 rounded-2xl p-6 space-y-6 shadow-2xl relative">
              <div className="border-b border-slate-800 pb-4">
                <div className="text-xs font-bold text-blue-400 uppercase tracking-widest">TICKET RECEIPT</div>
                <div className="text-xs text-slate-400 mt-1">{issuedTicket.service_name}</div>
              </div>

              <div>
                <div className="text-6xl font-black text-white tracking-tighter">{issuedTicket.ticket_number}</div>
                <div className="text-xs font-semibold text-amber-400 mt-2">
                  People Ahead: {issuedTicket.people_ahead || 0}
                </div>
                <div className="text-xs text-slate-400">
                  Est. Wait: ~{Math.floor(issuedTicket.estimated_wait_seconds / 60)} minutes
                </div>
              </div>

              <div className="bg-slate-900 rounded-xl p-3 border border-slate-800 text-[11px] text-slate-300">
                Tracking Link:{' '}
                <Link
                  to={`/ticket/${issuedTicket.public_token}`}
                  target="_blank"
                  className="text-blue-400 font-bold underline"
                >
                  /ticket/{issuedTicket.public_token.slice(0, 8)}...
                </Link>
              </div>

              <button
                onClick={() => window.print()}
                className="w-full py-2.5 bg-slate-800 hover:bg-slate-700 text-white font-bold text-xs rounded-xl border border-slate-700 flex items-center justify-center space-x-2 transition"
              >
                <Printer className="w-4 h-4" />
                <span>Print Ticket Slip</span>
              </button>
            </div>
          ) : (
            <div className="space-y-3">
              <div className="w-16 h-16 rounded-2xl bg-slate-800 flex items-center justify-center text-slate-500 mx-auto">
                <TicketPlus className="w-8 h-8" />
              </div>
              <h3 className="text-base font-bold text-white">Ticket Preview Terminal</h3>
              <p className="text-xs text-slate-400 max-w-xs mx-auto">
                Select service options on the left to issue a new concurrency-safe ticket.
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
