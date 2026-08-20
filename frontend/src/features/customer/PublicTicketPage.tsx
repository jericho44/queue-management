import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { QueueTicket } from '../../types';
import { fetchApi } from '../../api/client';
import { Clock, Users, ArrowUpRight, CheckCircle2, Volume2, ShieldCheck } from 'lucide-react';

export const PublicTicketPage: React.FC = () => {
  const { publicToken } = useParams<{ publicToken: string }>();
  const [ticket, setTicket] = useState<QueueTicket | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (publicToken) {
      fetchTicketStatus();
      setupWS(publicToken);
    }
  }, [publicToken]);

  const fetchTicketStatus = async () => {
    try {
      const res = await fetchApi<QueueTicket>(`/public/tickets/${publicToken}`);
      setTicket(res.data);
    } catch (err: any) {
      setError(err.message || 'Ticket not found');
    } finally {
      setLoading(false);
    }
  };

  const setupWS = (token: string) => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws?channel=ticket:${token}`;

    const ws = new WebSocket(wsUrl);
    ws.onmessage = () => {
      fetchTicketStatus();
    };
    return () => ws.close();
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-slate-950 flex items-center justify-center p-4">
        <div className="text-xs text-slate-400 font-semibold animate-pulse">Loading live ticket tracking...</div>
      </div>
    );
  }

  if (error || !ticket) {
    return (
      <div className="min-h-screen bg-slate-950 flex items-center justify-center p-4">
        <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 text-center max-w-sm">
          <div className="text-red-400 text-sm font-bold mb-2">Ticket Not Found</div>
          <p className="text-xs text-slate-400">{error || 'Invalid or expired ticket token.'}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-950 text-white flex flex-col justify-between p-4 max-w-md mx-auto">
      {/* Header */}
      <div className="flex items-center justify-between py-4 border-b border-slate-800">
        <div className="flex items-center space-x-2">
          <div className="w-8 h-8 rounded-xl bg-blue-600 flex items-center justify-center font-black text-sm">Q</div>
          <span className="font-extrabold text-sm tracking-tight">QFlow Live Tracker</span>
        </div>
        <div className="flex items-center space-x-1 text-[11px] text-emerald-400 bg-emerald-500/10 px-2.5 py-1 rounded-full border border-emerald-500/20 font-semibold">
          <span className="w-2 h-2 rounded-full bg-emerald-400 animate-ping"></span>
          <span>Live WebSocket</span>
        </div>
      </div>

      {/* Main Ticket Card */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-3xl p-8 text-center shadow-2xl space-y-6 my-auto">
        <div>
          <span className="text-[11px] font-bold tracking-widest text-slate-400 uppercase">YOUR TICKET NUMBER</span>
          <h1 className="text-6xl font-black text-blue-400 tracking-tighter mt-2 drop-shadow-md">
            {ticket.ticket_number}
          </h1>
          <div className="text-xs font-semibold text-slate-300 mt-2">{ticket.service_name}</div>
        </div>

        {/* Status Callout */}
        <div className="bg-slate-950 rounded-2xl p-5 border border-slate-800/80 space-y-4">
          <div className="flex items-center justify-between text-xs">
            <span className="text-slate-400 font-medium">Ticket Status:</span>
            <span
              className={`px-3 py-1 rounded-full font-black text-xs ${
                ticket.status === 'CALLED'
                  ? 'bg-amber-500/20 text-amber-300 animate-pulse border border-amber-500/40'
                  : ticket.status === 'SERVING'
                  ? 'bg-blue-500/20 text-blue-300 border border-blue-500/40'
                  : ticket.status === 'COMPLETED'
                  ? 'bg-emerald-500/20 text-emerald-400'
                  : 'bg-slate-800 text-slate-300'
              }`}
            >
              {ticket.status}
            </span>
          </div>

          {ticket.status === 'CALLED' && (
            <div className="p-4 bg-gradient-to-r from-amber-500/20 to-amber-600/20 border border-amber-500/30 rounded-xl text-amber-300 text-xs font-bold animate-bounce">
              PROCEED TO COUNTER {ticket.counter_number || '01'} NOW!
            </div>
          )}

          {ticket.status === 'WAITING' && (
            <div className="grid grid-cols-2 gap-3 pt-2 border-t border-slate-800/80">
              <div>
                <div className="text-[10px] text-slate-400 uppercase font-semibold">People Ahead</div>
                <div className="text-2xl font-black text-amber-400">{ticket.people_ahead || 0}</div>
              </div>
              <div>
                <div className="text-[10px] text-slate-400 uppercase font-semibold">Est. Wait</div>
                <div className="text-2xl font-black text-white">~{Math.floor(ticket.estimated_wait_seconds / 60)}m</div>
              </div>
            </div>
          )}
        </div>

        <p className="text-[11px] text-slate-400">
          This page automatically refreshes in real-time when your ticket is called.
        </p>
      </div>

      {/* Footer */}
      <div className="py-4 text-center text-[10px] text-slate-500 border-t border-slate-800">
        QFlow Multi-Tenant Queue Platform — End-to-End Encryption
      </div>
    </div>
  );
};
