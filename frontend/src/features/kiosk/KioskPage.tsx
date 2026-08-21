import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { fetchApi } from '../../api/client';
import { Branch, Service, QueueTicket } from '../../types';
import { ThermalReceipt } from '../../components/ThermalReceipt';
import { printDirectWebUSB } from '../../utils/escpos';
import {
  Ticket,
  Printer,
  QrCode,
  Sparkles,
  Clock,
  CheckCircle2,
  Maximize,
  RefreshCw,
  Users,
  Smartphone,
  ChevronRight,
  Zap,
} from 'lucide-react';

import { useToast } from '../../components/Toast';

export const KioskPage: React.FC = () => {
  const { showError, showSuccess } = useToast();
  const { branchId: branchIdentifier } = useParams<{ branchId: string }>();

  const [branch, setBranch] = useState<Branch | null>(null);
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);
  const [issuingServiceId, setIssuingServiceId] = useState<number | null>(null);

  // Issued ticket modal & state
  const [issuedTicket, setIssuedTicket] = useState<QueueTicket | null>(null);
  const [modalStep, setModalStep] = useState<'CHOICE' | 'PHYSICAL_PRINTING' | 'PAPERLESS_QR' | null>(null);
  const [countdown, setCountdown] = useState(10);

  const [timeStr, setTimeStr] = useState(new Date().toLocaleTimeString());

  useEffect(() => {
    const timer = setInterval(() => setTimeStr(new Date().toLocaleTimeString()), 1000);
    return () => clearInterval(timer);
  }, []);

  useEffect(() => {
    if (branchIdentifier) {
      loadKioskData(branchIdentifier);
    }
  }, [branchIdentifier]);

  // Auto reset countdown timer when ticket modal is open
  useEffect(() => {
    let interval: any = null;
    if (modalStep && countdown > 0) {
      interval = setInterval(() => {
        setCountdown((prev) => prev - 1);
      }, 1000);
    } else if (modalStep && countdown === 0) {
      resetKiosk();
    }
    return () => {
      if (interval) clearInterval(interval);
    };
  }, [modalStep, countdown]);

  const loadKioskData = async (bIdentifier: string) => {
    setLoading(true);
    try {
      const [branchRes, servicesRes] = await Promise.all([
        fetchApi<Branch>(`/public/branches/${bIdentifier}`),
        fetchApi<Service[]>(`/public/services?branch_id=${bIdentifier}`),
      ]);
      setBranch(branchRes.data);
      setServices(servicesRes.data || []);
    } catch (err: any) {
      console.error(err);
      showError('Gagal memuat data cabang kiosk');
    } finally {
      setLoading(false);
    }
  };

  const resetKiosk = () => {
    setIssuedTicket(null);
    setModalStep(null);
    setCountdown(10);
    setIssuingServiceId(null);
  };

  const handleSelectService = async (service: Service) => {
    if (!branchIdentifier || issuingServiceId) return;

    setIssuingServiceId(service.id);
    try {
      const res = await fetchApi<QueueTicket>('/public/tickets', {
        method: 'POST',
        body: JSON.stringify({
          branch_id: branch?.id || service.branch_id,
          service_id: service.id,
          priority: 'NORMAL',
        }),
      });


      const ticket = res.data;
      setIssuedTicket(ticket);
      showSuccess(`Nomor Antrean ${ticket.ticket_number} Berhasil Dibuat!`, 'Antrean Terbit');

      const mode = branch?.kiosk_mode || 'DUAL';

      if (mode === 'PHYSICAL') {
        executePhysicalPrint(ticket);
      } else if (mode === 'PAPERLESS') {
        setModalStep('PAPERLESS_QR');
      } else {
        // DUAL MODE
        setModalStep('CHOICE');
      }
      setCountdown(12);
    } catch (err: any) {
      showError(err.message || 'Gagal membuat tiket antrean');
    } finally {
      setIssuingServiceId(null);
    }
  };



  const executePhysicalPrint = async (ticket: QueueTicket) => {
    setModalStep('PHYSICAL_PRINTING');

    // Try WebUSB first
    if (branch) {
      const success = await printDirectWebUSB({
        branchName: branch.name,
        serviceName: ticket.service_name || '',
        ticketNumber: ticket.ticket_number,
        estimatedWaitMinutes: Math.max(1, Math.floor(ticket.estimated_wait_seconds / 60)),
        peopleAhead: ticket.people_ahead || 0,
        publicToken: ticket.public_token,
        dateStr: new Date().toLocaleTimeString(),
        headerText: branch.receipt_header,
        footerText: branch.receipt_footer,
        paperSize: branch.paper_size,
      });

      if (!success) {
        // Fallback to browser print
        setTimeout(() => {
          window.print();
        }, 300);
      }
    }
  };

  const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen().catch(() => {});
    } else {
      document.exitFullscreen().catch(() => {});
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-slate-950 flex items-center justify-center text-white">
        <div className="text-center space-y-4">
          <RefreshCw className="w-12 h-12 text-blue-500 animate-spin mx-auto" />
          <p className="text-lg font-bold">Memuat Terminal Kiosk Antrean...</p>
        </div>
      </div>
    );
  }

  if (!branch) {
    return (
      <div className="min-h-screen bg-slate-950 flex items-center justify-center text-white p-6">
        <div className="bg-slate-900 border border-slate-800 rounded-3xl p-8 max-w-md text-center">
          <h2 className="text-xl font-bold text-red-400 mb-2">Cabang Tidak Ditemukan</h2>
          <p className="text-slate-400 text-sm">Silakan periksa ID cabang pada URL terminal Kiosk Anda.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-950 text-white flex flex-col justify-between p-6 select-none relative overflow-hidden">
      {/* Hidden print element */}
      {issuedTicket && branch && <ThermalReceipt ticket={issuedTicket} branch={branch} />}

      {/* Top Header */}
      <header className="flex justify-between items-center bg-slate-900/80 border border-slate-800 rounded-2xl px-6 py-4 backdrop-blur-xl shadow-2xl">
        <div className="flex items-center space-x-4">
          <div className="w-14 h-14 rounded-2xl bg-gradient-to-tr from-blue-600 to-indigo-600 flex items-center justify-center shadow-lg shadow-blue-500/25">
            <Sparkles className="w-7 h-7 text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-black tracking-tight text-white">{branch.name}</h1>
            <p className="text-xs font-semibold text-blue-400">Terminal Ambil Tiket Antrean Mandiri</p>
          </div>
        </div>

        <div className="flex items-center space-x-6">
          <div className="flex items-center space-x-2 bg-slate-800/80 border border-slate-700/60 rounded-xl px-4 py-2 text-slate-200">
            <Clock className="w-5 h-5 text-amber-400" />
            <span className="font-mono text-base font-bold">{timeStr}</span>
          </div>

          <button
            onClick={toggleFullscreen}
            className="p-3 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl border border-slate-700 transition"
            title="Toggle Fullscreen"
          >
            <Maximize className="w-5 h-5" />
          </button>
        </div>
      </header>

      {/* Main Content: Touch Service Tiles */}
      <main className="my-auto py-8 max-w-6xl mx-auto w-full space-y-6">
        <div className="text-center space-y-2">
          <h2 className="text-3xl font-black tracking-tight text-white">Silakan Pilih Jenis Layanan Anda</h2>
          <p className="text-sm font-medium text-slate-400">Sentuh salah satu kotak layanan di bawah ini untuk mengambil nomor antrean</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {services.map((service) => {
            const isIssuing = issuingServiceId === service.id;
            return (
              <button
                key={service.id}
                disabled={Boolean(issuingServiceId)}
                onClick={() => handleSelectService(service)}
                className={`group relative bg-gradient-to-b from-slate-900 to-slate-900/90 border-2 rounded-3xl p-6 text-left shadow-2xl transition-all duration-300 transform active:scale-95 flex flex-col justify-between min-h-[220px] ${
                  isIssuing
                    ? 'border-amber-500 bg-amber-500/10'
                    : 'border-slate-800 hover:border-blue-500/80 hover:shadow-blue-500/20 hover:-translate-y-1'
                }`}
              >
                {/* Top Badge */}
                <div className="flex justify-between items-center">
                  <div className="w-12 h-12 rounded-2xl bg-blue-600/20 border border-blue-500/30 text-blue-400 font-black text-xl flex items-center justify-center group-hover:scale-110 transition">
                    {service.prefix}
                  </div>
                  <span className="px-3 py-1 bg-slate-800 rounded-full text-[11px] font-bold text-slate-300 border border-slate-700">
                    Est. ~{Math.floor(service.avg_service_time_sec / 60)} min
                  </span>
                </div>

                {/* Service Details */}
                <div className="space-y-1 my-4">
                  <h3 className="text-xl font-black text-white group-hover:text-blue-300 transition">
                    {service.name}
                  </h3>
                  <p className="text-xs text-slate-400 font-mono">Kode Servis: {service.code}</p>
                </div>

                {/* Bottom Call to Action */}
                <div className="flex items-center justify-between border-t border-slate-800/80 pt-3">
                  <span className="text-xs font-bold text-blue-400 group-hover:underline flex items-center space-x-1">
                    <span>{isIssuing ? 'Memproses Tiket...' : 'Ambil Tiket'}</span>
                    {!isIssuing && <ChevronRight className="w-4 h-4" />}
                  </span>
                  {isIssuing && <RefreshCw className="w-5 h-5 text-amber-400 animate-spin" />}
                </div>
              </button>
            );
          })}
        </div>
      </main>

      {/* Footer Instructions */}
      <footer className="bg-slate-900/60 border border-slate-800 rounded-2xl px-6 py-3 text-center text-xs text-slate-400 backdrop-blur-md flex justify-between items-center">
        <span>Powered by <strong>Multi-Tenant QMS SaaS</strong></span>
        <span className="flex items-center space-x-2 text-emerald-400 font-semibold">
          <Zap className="w-4 h-4" />
          <span>Sistem Antrean Concurrency-Safe Aktif</span>
        </span>
      </footer>

      {/* Ticket Modal (Choice / Physical / Paperless) */}
      {modalStep && issuedTicket && (
        <div className="fixed inset-0 bg-slate-950/90 backdrop-blur-xl z-50 flex items-center justify-center p-6 animate-in fade-in duration-200">
          <div className="bg-slate-900 border-2 border-blue-500/40 rounded-3xl p-8 max-w-lg w-full text-center space-y-6 shadow-2xl relative">
            {/* Auto reset badge */}
            <div className="absolute top-4 right-4 bg-slate-800 text-amber-400 px-3.5 py-1 rounded-full text-xs font-bold border border-slate-700 flex items-center space-x-1">
              <Clock className="w-3.5 h-3.5" />
              <span>Reset otomatis ({countdown}s)</span>
            </div>

            {/* DUAL MODE CHOICE */}
            {modalStep === 'CHOICE' && (
              <div className="space-y-6 pt-2">
                <div className="w-16 h-16 rounded-3xl bg-blue-600/20 text-blue-400 border border-blue-500/40 flex items-center justify-center mx-auto shadow-lg shadow-blue-500/20">
                  <Ticket className="w-8 h-8" />
                </div>

                <div>
                  <h3 className="text-2xl font-black text-white">Nomor Antrean Anda</h3>
                  <div className="text-6xl font-black text-amber-400 my-2 tracking-tight">
                    {issuedTicket.ticket_number}
                  </div>
                  <p className="text-xs text-slate-300">Pilih metode tiket yang Anda inginkan:</p>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 pt-2">
                  <button
                    onClick={() => executePhysicalPrint(issuedTicket)}
                    className="p-5 bg-gradient-to-b from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white rounded-2xl font-bold flex flex-col items-center space-y-2 shadow-lg transition active:scale-95"
                  >
                    <Printer className="w-8 h-8" />
                    <span className="text-sm">Cetak Struk Fisik</span>
                    <span className="text-[10px] text-blue-200 font-normal">Kertas Thermal</span>
                  </button>

                  <button
                    onClick={() => setModalStep('PAPERLESS_QR')}
                    className="p-5 bg-slate-800 hover:bg-slate-700 border border-emerald-500/40 text-emerald-400 rounded-2xl font-bold flex flex-col items-center space-y-2 shadow-lg transition active:scale-95"
                  >
                    <Smartphone className="w-8 h-8 text-emerald-400" />
                    <span className="text-sm">Tiket Digital</span>
                    <span className="text-[10px] text-emerald-300 font-normal">Scan Paperless QR</span>
                  </button>
                </div>
              </div>
            )}

            {/* PHYSICAL PRINTING STEP */}
            {modalStep === 'PHYSICAL_PRINTING' && (
              <div className="space-y-6 py-4">
                <div className="w-16 h-16 rounded-3xl bg-emerald-500/20 text-emerald-400 border border-emerald-500/40 flex items-center justify-center mx-auto animate-bounce">
                  <Printer className="w-8 h-8" />
                </div>
                <div>
                  <h3 className="text-2xl font-black text-white">Mencetak Struk Antrean...</h3>
                  <div className="text-6xl font-black text-amber-400 my-2 tracking-tight">
                    {issuedTicket.ticket_number}
                  </div>
                  <p className="text-sm text-slate-300">Harap ambil struk kertas fisik Anda dari mesin printer.</p>
                </div>
                <button
                  onClick={resetKiosk}
                  className="w-full py-3 bg-slate-800 hover:bg-slate-700 text-white font-bold rounded-xl border border-slate-700 text-xs transition"
                >
                  Selesai / Layar Awal
                </button>
              </div>
            )}

            {/* PAPERLESS QR STEP */}
            {modalStep === 'PAPERLESS_QR' && (
              <div className="space-y-4 py-2">
                <div className="text-xs font-bold text-emerald-400 uppercase tracking-widest">
                  🌱 Mode Paperless (Hemat Kertas)
                </div>

                <div>
                  <div className="text-5xl font-black text-white tracking-tight">{issuedTicket.ticket_number}</div>
                  <p className="text-xs text-slate-400 mt-1">Scan QR Code di bawah untuk melacak antrean di HP Anda</p>
                </div>

                <div className="bg-white p-4 rounded-2xl w-48 h-48 mx-auto shadow-2xl flex items-center justify-center">
                  <img
                    src={`https://api.qrserver.com/v1/create-qr-code/?size=180x180&data=${encodeURIComponent(
                      `${window.location.origin}/ticket/${issuedTicket.public_token}`
                    )}`}
                    alt="Paperless QR Code"
                    className="w-full h-full"
                  />
                </div>

                <div className="bg-slate-950 p-3 rounded-xl border border-slate-800 text-xs text-slate-300">
                  Est. Waktu Tunggu: <strong>~{Math.max(1, Math.floor(issuedTicket.estimated_wait_seconds / 60))} Menit</strong>
                  <div className="text-[11px] text-slate-400">Sisa Antrean: {issuedTicket.people_ahead || 0} Orang</div>
                </div>

                <button
                  onClick={resetKiosk}
                  className="w-full py-3.5 bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-500 hover:to-teal-500 text-white font-black rounded-xl text-xs shadow-lg transition"
                >
                  Sudah Dican / Tutup
                </button>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};
