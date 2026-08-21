import React, { useEffect, useState } from 'react';
import { fetchApi } from '../../api/client';
import { Counter, QueueTicket, Branch, Service } from '../../types';
import { Link } from 'react-router-dom';
import {
  Volume2,
  CheckCircle2,
  SkipForward,
  UserX,
  Play,
  RotateCcw,
  ArrowRightLeft,
  Kanban,
  Power,
  Clock,
  Sparkles,
  AlertTriangle,
  Building2,
  ChevronDown,
  Monitor,
} from 'lucide-react';


import { useToast } from '../../components/Toast';

export const StaffCounterPage: React.FC = () => {
  const { showError, showSuccess, showInfo, showWarning } = useToast();
  const [branches, setBranches] = useState<Branch[]>([]);
  const [counters, setCounters] = useState<Counter[]>([]);
  const [selectedBranchId, setSelectedBranchId] = useState<number>(0);
  const [selectedCounterId, setSelectedCounterId] = useState<number>(0);
  const [activeCounter, setActiveCounter] = useState<Counter | null>(null);
  const [currentTicket, setCurrentTicket] = useState<QueueTicket | null>(null);
  const [waitingTickets, setWaitingTickets] = useState<QueueTicket[]>([]);
  const [services, setServices] = useState<Service[]>([]);
  const [transferModalOpen, setTransferModalOpen] = useState(false);
  const [targetServiceId, setTargetServiceId] = useState<number>(0);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchBranches();
  }, []);

  useEffect(() => {
    if (selectedBranchId > 0) {
      fetchCounters(selectedBranchId);
      fetchServices(selectedBranchId);
    }
  }, [selectedBranchId]);

  const fetchBranches = async () => {
    try {
      const res = await fetchApi<Branch[]>('/branches');
      setBranches(res.data || []);
      if (res.data && res.data.length > 0) {
        setSelectedBranchId(res.data[0].id);
      }
    } catch (err: any) {
      showError(err.message || 'Gagal memuat cabang');
    }
  };

  const fetchCounters = async (bId: number) => {
    try {
      const res = await fetchApi<Counter[]>(`/counters?branch_id=${bId}`);
      const counterList = res.data || [];
      setCounters(counterList);
      if (counterList.length > 0) {
        setSelectedCounterId(counterList[0].id);
        setActiveCounter(counterList[0]);
      } else {
        setSelectedCounterId(0);
        setActiveCounter(null);
      }
    } catch (err: any) {
      showError(err.message || 'Gagal memuat loket');
    }
  };

  const fetchServices = async (bId: number) => {
    try {
      const res = await fetchApi<Service[]>(`/services?branch_id=${bId}`);
      setServices(res.data || []);
    } catch (err: any) {
      showError(err.message || 'Gagal memuat layanan');
    }
  };

  const fetchWaitingQueue = async () => {
    if (!selectedBranchId) return;
    try {
      const res = await fetchApi<QueueTicket[]>(`/tickets/waiting?branch_id=${selectedBranchId}`);
      setWaitingTickets(res.data || []);
    } catch (err: any) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchWaitingQueue();
    const interval = setInterval(fetchWaitingQueue, 5000);
    return () => clearInterval(interval);
  }, [selectedBranchId]);

  const handleOpenCounter = async () => {
    if (!selectedCounterId) {
      showWarning('Silakan buat Loket terlebih dahulu di menu Branches & Counters');
      return;
    }
    try {
      await fetchApi(`/counters/${selectedCounterId}/open`, { method: 'POST' });
      const c = counters.find((ct) => ct.id === selectedCounterId);
      if (c) setActiveCounter({ ...c, status: 'OPEN' });
      showSuccess('Loket berhasil dibuka');
    } catch (err: any) {
      showError(err.message || 'Gagal membuka loket');
    }
  };

  const handleCloseCounter = async () => {
    if (!selectedCounterId) return;
    try {
      await fetchApi(`/counters/${selectedCounterId}/close`, { method: 'POST' });
      const c = counters.find((ct) => ct.id === selectedCounterId);
      if (c) setActiveCounter({ ...c, status: 'CLOSED' });
      setCurrentTicket(null);
      showInfo('Loket berhasil ditutup');
    } catch (err: any) {
      showError(err.message || 'Gagal menutup loket');
    }
  };

  const handleCallNext = async () => {
    if (!selectedCounterId || selectedCounterId === 0) {
      showWarning('Silakan pilih Loket terlebih dahulu atau buat Loket baru di menu Branches & Counters');
      return;
    }
    if (activeCounter?.status !== 'OPEN' && activeCounter?.status !== 'BUSY') {
      showWarning('Buka Loket terlebih dahulu dengan menekan tombol "Open Counter"');
      return;
    }

    setLoading(true);
    try {
      const res = await fetchApi<QueueTicket>(`/tickets/counters/${selectedCounterId}/next`, {
        method: 'POST',
      });
      setCurrentTicket(res.data);
      fetchWaitingQueue();
      showSuccess(`Memanggil Antrean ${res.data.ticket_number}`);
    } catch (err: any) {
      showError(err.message || 'Tidak ada tiket antrean yang menunggu');
    } finally {
      setLoading(false);
    }
  };


  const handleRecall = async () => {
    if (!currentTicket) return;
    try {
      await fetchApi(`/tickets/${currentTicket.id}/recall`, { method: 'POST' });
      showInfo(`Memanggil Ulang Tiket ${currentTicket.ticket_number}`);
    } catch (err: any) {
      showError(err.message || 'Gagal memanggil ulang');
    }
  };

  const handleStartServing = async () => {
    if (!currentTicket) return;
    try {
      const res = await fetchApi<QueueTicket>(`/tickets/${currentTicket.id}/start`, { method: 'POST' });
      setCurrentTicket(res.data);
      showSuccess(`Memulai pelayanan untuk ${res.data.ticket_number}`);
    } catch (err: any) {
      showError(err.message || 'Gagal melayani tiket');
    }
  };

  const handleComplete = async () => {
    if (!currentTicket) return;
    try {
      await fetchApi(`/tickets/${currentTicket.id}/complete`, { method: 'POST' });
      showSuccess(`Tiket ${currentTicket.ticket_number} Selesai Dilayani`);
      setCurrentTicket(null);
      fetchWaitingQueue();
    } catch (err: any) {
      showError(err.message || 'Gagal menyelesaikan pelayanan');
    }
  };

  const handleSkip = async () => {
    if (!currentTicket) return;
    try {
      await fetchApi(`/tickets/${currentTicket.id}/skip`, { method: 'POST' });
      showInfo(`Tiket ${currentTicket.ticket_number} Dilewati`);
      setCurrentTicket(null);
      fetchWaitingQueue();
    } catch (err: any) {
      showError(err.message || 'Gagal melewati tiket');
    }
  };

  const handleNoShow = async () => {
    if (!currentTicket) return;
    try {
      await fetchApi(`/tickets/${currentTicket.id}/no-show`, { method: 'POST' });
      showInfo(`Tiket ${currentTicket.ticket_number} Ditandai Tidak Hadir`);
      setCurrentTicket(null);
      fetchWaitingQueue();
    } catch (err: any) {
      showError(err.message || 'Gagal memproses tiket');
    }
  };

  return (
    <div className="space-y-6">
      {/* Top Station Selector */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-2xl p-5 flex flex-col md:flex-row items-center justify-between gap-4">
        <div className="flex items-center space-x-3">
          <div className="w-10 h-10 rounded-xl bg-blue-600/20 text-blue-400 border border-blue-500/30 flex items-center justify-center">
            <Kanban className="w-5 h-5" />
          </div>
          <div>
            <h1 className="text-xl font-extrabold text-white">Staff Counter Terminal</h1>
            <p className="text-xs text-slate-400">Manage calls, serving, and ticket disposition</p>
          </div>
        </div>

        <div className="flex flex-col sm:flex-row items-center gap-3 w-full md:w-auto">
          {/* Branch Selector Dropdown */}
          <div className="relative w-full sm:w-52">
            <Building2 className="w-4 h-4 text-blue-400 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none z-10" />
            <select
              value={selectedBranchId}
              onChange={(e) => setSelectedBranchId(Number(e.target.value))}
              className="w-full appearance-none bg-slate-800/90 hover:bg-slate-800 border border-slate-700/80 focus:border-blue-500 rounded-xl pl-9 pr-9 py-2.5 text-xs text-white font-bold transition shadow-sm cursor-pointer focus:outline-none focus:ring-2 focus:ring-blue-500/20"
            >
              {branches.map((b) => (
                <option key={b.id} value={b.id} className="bg-slate-900 text-white font-semibold py-1">
                  {b.name} ({b.code})
                </option>
              ))}
            </select>
            <ChevronDown className="w-4 h-4 text-slate-400 absolute right-3.5 top-1/2 -translate-y-1/2 pointer-events-none z-10" />
          </div>

          {/* Counter Selector Dropdown */}
          <div className="relative w-full sm:w-64">
            <Monitor className="w-4 h-4 text-indigo-400 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none z-10" />
            <select
              value={selectedCounterId}
              onChange={(e) => {
                const cid = Number(e.target.value);
                setSelectedCounterId(cid);
                const ct = counters.find((c) => c.id === cid);
                setActiveCounter(ct || null);
              }}
              className="w-full appearance-none bg-slate-800/90 hover:bg-slate-800 border border-slate-700/80 focus:border-indigo-500 rounded-xl pl-9 pr-9 py-2.5 text-xs text-white font-bold transition shadow-sm cursor-pointer focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
            >
              {counters.length === 0 ? (
                <option value={0} className="bg-slate-900 text-amber-400 font-semibold">-- Belum Ada Loket --</option>
              ) : (
                counters.map((c) => (
                  <option key={c.id} value={c.id} className="bg-slate-900 text-white font-semibold py-1">
                    Loket {c.counter_number} — {c.name} ({c.status})
                  </option>
                ))
              )}
            </select>
            <ChevronDown className="w-4 h-4 text-slate-400 absolute right-3.5 top-1/2 -translate-y-1/2 pointer-events-none z-10" />
          </div>


          {activeCounter?.status === 'OPEN' || activeCounter?.status === 'BUSY' ? (
            <button
              onClick={handleCloseCounter}
              className="flex items-center space-x-1.5 px-4 py-2 rounded-xl bg-red-600/20 hover:bg-red-600/30 text-red-400 border border-red-500/30 text-xs font-bold transition"
            >
              <Power className="w-3.5 h-3.5" />
              <span>Close Counter</span>
            </button>
          ) : (
            <button
              onClick={handleOpenCounter}
              className="flex items-center space-x-1.5 px-4 py-2 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-bold shadow-lg shadow-emerald-600/20 transition"
            >
              <Power className="w-3.5 h-3.5" />
              <span>Open Counter</span>
            </button>
          )}
        </div>
      </div>

      {counters.length === 0 && (
        <div className="bg-amber-500/10 border border-amber-500/30 rounded-2xl p-4 flex items-center justify-between text-amber-300 text-xs shadow-lg">
          <div className="flex items-center space-x-2">
            <AlertTriangle className="w-4 h-4 text-amber-400 flex-shrink-0" />
            <span>Cabang ini belum memiliki Loket Pelayanan. Silakan buat Loket di menu <strong>Branches & Counters</strong> agar dapat melakukan panggilan antrean.</span>
          </div>
          <Link to="/branches" className="px-3 py-1.5 bg-amber-500 text-slate-950 font-bold rounded-xl hover:bg-amber-400 transition flex-shrink-0">
            Buat Loket Sekarang
          </Link>
        </div>
      )}


      {/* Main Counter Console */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column: Currently Active Ticket */}
        <div className="lg:col-span-2 space-y-6">
          <div className="bg-slate-900/90 border border-slate-800 rounded-3xl p-8 text-center relative overflow-hidden shadow-2xl">
            <div className="absolute top-4 left-4 text-xs font-semibold text-slate-400 uppercase tracking-widest flex items-center space-x-2">
              <span className="w-2 h-2 rounded-full bg-emerald-400 animate-ping"></span>
              <span>
                Counter {activeCounter?.counter_number || '01'} ({activeCounter?.name || 'General'})
              </span>
            </div>

            {currentTicket ? (
              <div className="space-y-6 py-6">
                <div>
                  <span className="inline-block px-3 py-1 rounded-full text-xs font-bold bg-blue-500/10 text-blue-400 border border-blue-500/20 mb-3">
                    STATUS: {currentTicket.status}
                  </span>
                  <h2 className="text-7xl font-black text-white tracking-tighter drop-shadow-lg">
                    {currentTicket.ticket_number}
                  </h2>
                  <p className="text-sm font-semibold text-slate-300 mt-2">{currentTicket.service_name}</p>
                </div>

                {/* Operations Bar */}
                <div className="flex flex-wrap items-center justify-center gap-3 pt-4 border-t border-slate-800">
                  <button
                    onClick={handleRecall}
                    className="flex items-center space-x-2 px-5 py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-amber-400 border border-slate-700 text-xs font-bold transition"
                  >
                    <Volume2 className="w-4 h-4" />
                    <span>Recall Audio</span>
                  </button>

                  {currentTicket.status === 'CALLED' && (
                    <button
                      onClick={handleStartServing}
                      className="flex items-center space-x-2 px-5 py-2.5 rounded-xl bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold shadow-lg shadow-blue-600/25 transition"
                    >
                      <Play className="w-4 h-4" />
                      <span>Start Serving</span>
                    </button>
                  )}

                  <button
                    onClick={handleComplete}
                    className="flex items-center space-x-2 px-6 py-2.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-bold shadow-lg shadow-emerald-600/25 transition"
                  >
                    <CheckCircle2 className="w-4 h-4" />
                    <span>Complete Ticket</span>
                  </button>

                  <button
                    onClick={handleSkip}
                    className="flex items-center space-x-2 px-4 py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 border border-slate-700 text-xs font-semibold transition"
                  >
                    <SkipForward className="w-4 h-4" />
                    <span>Skip</span>
                  </button>

                  <button
                    onClick={handleNoShow}
                    className="flex items-center space-x-2 px-4 py-2.5 rounded-xl bg-red-600/20 hover:bg-red-600/30 text-red-400 border border-red-500/30 text-xs font-semibold transition"
                  >
                    <UserX className="w-4 h-4" />
                    <span>No Show</span>
                  </button>
                </div>
              </div>
            ) : (
              <div className="py-12 space-y-4">
                <div className="w-16 h-16 rounded-full bg-slate-800 flex items-center justify-center text-slate-500 mx-auto">
                  <Sparkles className="w-8 h-8" />
                </div>
                <h3 className="text-lg font-bold text-white">No Customer Currently Assigned</h3>
                <p className="text-xs text-slate-400 max-w-sm mx-auto">
                  Click <strong className="text-blue-400">CALL NEXT TICKET</strong> to atomically fetch the next customer from the priority queue.
                </p>
                <button
                  onClick={handleCallNext}
                  disabled={loading}
                  className="px-8 py-3.5 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white font-extrabold text-sm rounded-2xl shadow-xl shadow-blue-600/30 transition transform hover:scale-105"
                >
                  {loading ? 'Fetching Next...' : 'CALL NEXT TICKET'}
                </button>
              </div>
            )}
          </div>

          {currentTicket && (
            <div className="text-center">
              <button
                onClick={handleCallNext}
                disabled={loading}
                className="w-full py-4 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white font-extrabold text-sm rounded-2xl shadow-xl shadow-blue-600/25 transition"
              >
                {loading ? 'Calling Next...' : 'NEXT CUSTOMER (CALL NEXT)'}
              </button>
            </div>
          )}
        </div>

        {/* Right Column: Branch Waiting Queue */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-2xl p-5 flex flex-col h-[500px]">
          <div className="flex items-center justify-between mb-4 pb-3 border-b border-slate-800">
            <div>
              <h3 className="text-sm font-bold text-white">Branch Waiting Queue</h3>
              <p className="text-[11px] text-slate-400">{waitingTickets.length} customers in line</p>
            </div>
            <Clock className="w-4 h-4 text-blue-400" />
          </div>

          <div className="flex-1 overflow-y-auto space-y-2 pr-1">
            {waitingTickets.length === 0 ? (
              <div className="h-full flex items-center justify-center text-xs text-slate-500">
                Queue is completely clear!
              </div>
            ) : (
              waitingTickets.map((t) => (
                <div
                  key={t.id}
                  className="bg-slate-800/50 border border-slate-700/60 rounded-xl p-3 flex items-center justify-between hover:bg-slate-800 transition"
                >
                  <div>
                    <div className="text-base font-black text-blue-400">{t.ticket_number}</div>
                    <div className="text-[11px] font-medium text-slate-300">{t.service_name}</div>
                  </div>
                  <div className="text-right">
                    <span
                      className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${
                        t.priority === 'EMERGENCY'
                          ? 'bg-red-500/20 text-red-400'
                          : t.priority === 'PRIORITY'
                          ? 'bg-amber-500/20 text-amber-400'
                          : 'bg-slate-700 text-slate-300'
                      }`}
                    >
                      {t.priority}
                    </span>
                    <div className="text-[10px] text-slate-400 mt-1">Est ~{Math.floor(t.estimated_wait_seconds / 60)}m</div>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
