import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { QueueTicket, Branch } from '../../types';
import { fetchApi } from '../../api/client';
import { Volume2, Monitor, Building2, Clock, Sparkles } from 'lucide-react';

export const PublicDisplayPage: React.FC = () => {
  const { branchIdentifier } = useParams<{ branchIdentifier?: string }>();
  const [activeBranch, setActiveBranch] = useState<Branch | null>(null);
  const [calledTickets, setCalledTickets] = useState<QueueTicket[]>([]);
  const [waitingTickets, setWaitingTickets] = useState<QueueTicket[]>([]);
  const [latestCall, setLatestCall] = useState<QueueTicket | null>(null);
  const [currentTime, setCurrentTime] = useState<string>('');

  const targetIdentifier = branchIdentifier || 'sdr';

  useEffect(() => {
    const clockTimer = setInterval(() => {
      setCurrentTime(new Date().toLocaleTimeString());
    }, 1000);
    return () => clearInterval(clockTimer);
  }, []);

  useEffect(() => {
    if (!targetIdentifier) return;

    loadBranchData(targetIdentifier);
    loadDisplayTickets(targetIdentifier);

    const wsCleanup = setupWebSocket(targetIdentifier);
    const pollTimer = setInterval(() => {
      loadDisplayTickets(targetIdentifier);
    }, 3000);

    return () => {
      if (wsCleanup) wsCleanup();
      clearInterval(pollTimer);
    };
  }, [targetIdentifier]);


  const loadBranchData = async (bIdentifier: string) => {
    try {
      const res = await fetchApi<Branch>(`/public/branches/${bIdentifier}`);
      setActiveBranch(res.data);
    } catch (err) {
      console.error(err);
    }
  };

  const loadDisplayTickets = async (bIdentifier: string) => {
    try {
      const calledRes = await fetchApi<QueueTicket[]>(`/public/display?branch_id=${bIdentifier}&status=CALLED`);
      const waitingRes = await fetchApi<QueueTicket[]>(`/public/display?branch_id=${bIdentifier}&status=WAITING`);
      const calledList = calledRes.data || [];
      const waitingList = waitingRes.data || [];

      setCalledTickets(calledList);
      setWaitingTickets(waitingList);
      if (calledList.length > 0) {
        setLatestCall(calledList[0]);
      } else {
        setLatestCall(null);
      }
    } catch (err) {
      console.error(err);
    }
  };

  const playChimeSound = () => {
    try {
      const AudioCtx = window.AudioContext || (window as any).webkitAudioContext;
      if (!AudioCtx) return;
      const ctx = new AudioCtx();
      const osc = ctx.createOscillator();
      const gain = ctx.createGain();

      osc.type = 'sine';
      osc.frequency.setValueAtTime(880, ctx.currentTime); // A5 note
      osc.frequency.exponentialRampToValueAtTime(440, ctx.currentTime + 0.5); // A4 note

      gain.gain.setValueAtTime(0.3, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.8);

      osc.connect(gain);
      gain.connect(ctx.destination);

      osc.start();
      osc.stop(ctx.currentTime + 0.8);
    } catch (e) {
      console.log('Audio chime auto-play prevented', e);
    }
  };

  const setupWebSocket = (bIdentifier: string) => {
    const orgId = activeBranch?.organization_id || 1;
    const branchId = activeBranch?.id || 1;
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws?channel=org:${orgId}:branch:${branchId}`;

    const ws = new WebSocket(wsUrl);

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (msg.event === 'QUEUE_TICKET_CALLED' || msg.event === 'QUEUE_TICKET_RECALLED') {
          const ticket: QueueTicket = msg.data;
          setLatestCall(ticket);
          setCalledTickets((prev) => [ticket, ...prev.filter((t) => t.id !== ticket.id)].slice(0, 4));
          playChimeSound();
        }
        loadDisplayTickets(bIdentifier);
      } catch (e) {
        console.error('WS error', e);
      }
    };

    return () => ws.close();
  };



  return (
    <div className="min-h-screen bg-slate-950 text-white flex flex-col p-6 font-sans select-none overflow-hidden">
      {/* Top Banner */}
      <header className="h-20 bg-slate-900/90 border border-slate-800 rounded-3xl px-8 flex items-center justify-between mb-6 shadow-2xl backdrop-blur-xl">
        <div className="flex items-center space-x-4">
          <div className="w-12 h-12 rounded-2xl bg-gradient-to-tr from-blue-600 to-indigo-500 flex items-center justify-center font-black text-2xl shadow-xl shadow-blue-500/25">
            Q
          </div>
          <div>
            <h1 className="text-2xl font-black tracking-tight text-white">{activeBranch?.name || 'PUBLIC QUEUE DISPLAY'}</h1>
            <p className="text-xs text-blue-400 font-semibold tracking-wider uppercase">Live Realtime Feed {activeBranch?.code ? `• ${activeBranch.code}` : ''}</p>
          </div>
        </div>

        <div className="flex items-center space-x-6">


          <div className="text-right">
            <div className="text-2xl font-black text-white font-mono tracking-tight">{currentTime}</div>
            <div className="text-[10px] text-slate-400 uppercase font-bold tracking-widest">Local Time</div>
          </div>
        </div>
      </header>

      {/* Main Grid */}
      <div className="flex-1 grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column: Big Now Serving Display */}
        <div className="lg:col-span-2 bg-slate-900/90 border border-slate-800 rounded-3xl p-10 flex flex-col items-center justify-center text-center shadow-2xl relative overflow-hidden">
          <div className="absolute top-6 left-8 text-sm font-extrabold text-blue-400 uppercase tracking-widest flex items-center space-x-2">
            <span className="w-3 h-3 rounded-full bg-blue-500 animate-ping"></span>
            <span>NOW SERVING</span>
          </div>

          {latestCall ? (
            <div className="space-y-6 my-auto animate-pulse-ring">
              <div className="text-[120px] lg:text-[160px] font-black text-white leading-none tracking-tighter drop-shadow-[0_10px_35px_rgba(59,130,246,0.3)]">
                {latestCall.ticket_number}
              </div>

              <div className="bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-3xl px-12 py-6 inline-block shadow-2xl shadow-blue-600/30">
                <div className="text-2xl font-bold uppercase tracking-wider text-blue-100">PROCEED TO</div>
                <div className="text-5xl font-black tracking-tight mt-1">COUNTER {latestCall.counter_number || '01'}</div>
              </div>

              <div className="text-xl font-bold text-slate-300 pt-4">Service: {latestCall.service_name}</div>
            </div>
          ) : (
            <div className="my-auto space-y-4">
              <Sparkles className="w-16 h-16 text-slate-600 mx-auto" />
              <h2 className="text-3xl font-black text-slate-400">Waiting for Next Call...</h2>
              <p className="text-sm text-slate-500">Please relax in the waiting lounge.</p>
            </div>
          )}
        </div>

        {/* Right Column: Waiting & Recent Calls List */}
        <div className="space-y-6 flex flex-col">
          {/* Active Called Board */}
          <div className="bg-slate-900/90 border border-slate-800 rounded-3xl p-6 shadow-xl flex-1 flex flex-col">
            <div className="text-xs font-extrabold text-slate-400 uppercase tracking-widest mb-4 pb-3 border-b border-slate-800">
              Active Counter Assignments
            </div>

            <div className="space-y-3 flex-1 overflow-y-auto">
              {calledTickets.length === 0 ? (
                <div className="text-xs text-slate-500 text-center py-6">No counters currently calling</div>
              ) : (
                calledTickets.map((t) => (
                  <div
                    key={t.id}
                    className="bg-slate-800/60 border border-slate-700/80 rounded-2xl p-4 flex items-center justify-between shadow-md"
                  >
                    <div>
                      <div className="text-3xl font-black text-blue-400">{t.ticket_number}</div>
                      <div className="text-xs font-medium text-slate-300">{t.service_name}</div>
                    </div>
                    <div className="text-right">
                      <div className="text-sm font-black text-emerald-400">COUNTER {t.counter_number}</div>
                      <div className="text-[10px] text-slate-400 font-semibold uppercase">{t.status}</div>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>

          {/* Waiting Count Box */}
          <div className="bg-slate-900/90 border border-slate-800 rounded-3xl p-6 shadow-xl flex items-center justify-between">
            <div>
              <div className="text-xs font-bold text-slate-400 uppercase tracking-wider">Total Waiting In Line</div>
              <div className="text-4xl font-black text-amber-400 mt-1">{waitingTickets.length}</div>
            </div>
            <Clock className="w-10 h-10 text-amber-400/30" />
          </div>
        </div>
      </div>
    </div>
  );
};
