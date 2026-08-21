import React from 'react';
import { QueueTicket, Branch } from '../types';

interface ThermalReceiptProps {
  ticket: QueueTicket;
  branch: Branch;
}

export const ThermalReceipt: React.FC<ThermalReceiptProps> = ({ ticket, branch }) => {
  const paperWidthClass = branch.paper_size === '80mm' ? 'w-[78mm]' : 'w-[56mm]';
  const estWaitMin = Math.max(1, Math.floor((ticket.estimated_wait_seconds || 0) / 60));

  return (
    <div className="hidden print:block print:fixed print:top-0 print:left-0 print:bg-white print:text-black">
      <style>{`
        @media print {
          @page {
            margin: 0;
            size: ${branch.paper_size === '80mm' ? '80mm auto' : '58mm auto'};
          }
          body {
            background: white !important;
            color: black !important;
            font-family: 'Courier New', Courier, monospace;
            margin: 0;
            padding: 4px;
          }
        }
      `}</style>
      <div className={`${paperWidthClass} mx-auto text-center p-2 text-black space-y-2`}>
        {/* Optional Header */}
        {branch.receipt_header && (
          <div className="text-[10px] font-bold uppercase tracking-tight border-b border-black pb-1">
            {branch.receipt_header}
          </div>
        )}

        {/* Branch Name */}
        <div className="text-sm font-black uppercase leading-tight">{branch.name}</div>
        <div className="text-[9px] uppercase">
          {typeof branch.address === 'string' ? (branch.address || branch.code) : (branch.address as any)?.String || branch.code}
        </div>



        <div className="border-t border-b border-dashed border-black py-1">
          <div className="text-[10px] font-bold uppercase">LAYANAN: {ticket.service_name}</div>
        </div>

        {/* Big Ticket Number */}
        <div className="py-2">
          <div className="text-[10px] uppercase font-semibold">Nomor Antrean</div>
          <div className="text-4xl font-black tracking-tighter leading-none">{ticket.ticket_number}</div>
        </div>

        {/* Stats */}
        <div className="border-t border-dashed border-black pt-1 text-[10px] space-y-0.5 text-left px-1">
          <div className="flex justify-between">
            <span>Sisa Antrean:</span>
            <span className="font-bold">{ticket.people_ahead || 0} Orang</span>
          </div>
          <div className="flex justify-between">
            <span>Est. Waktu Tunggu:</span>
            <span className="font-bold">~{estWaitMin} Menit</span>
          </div>
          <div className="flex justify-between">
            <span>Waktu Cetak:</span>
            <span>{new Date(ticket.created_at || Date.now()).toLocaleTimeString()}</span>
          </div>
        </div>

        {/* QR Code / Tracking Link */}
        <div className="border-t border-dashed border-black pt-2 flex flex-col items-center">
          <img
            src={`https://api.qrserver.com/v1/create-qr-code/?size=110x110&data=${encodeURIComponent(
              `${window.location.origin}/ticket/${ticket.public_token}`
            )}`}
            alt="Ticket Tracking QR Code"
            className="w-24 h-24 mx-auto"
          />
          <div className="text-[8px] font-mono mt-1">Scan untuk Live Tracking Smartphone</div>
        </div>

        {/* Footer */}
        <div className="border-t border-black pt-2 text-[9px] leading-tight">
          {branch.receipt_footer || 'Terima kasih atas kunjungan Anda.'}
        </div>
      </div>
    </div>
  );
};
