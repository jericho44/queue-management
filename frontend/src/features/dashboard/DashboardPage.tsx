import React, { useEffect, useState } from 'react';
import { fetchApi } from '../../api/client';
import { DashboardStats, QueueTicket, Branch } from '../../types';
import {
  Users,
  CheckCircle,
  Clock,
  Kanban,
  Activity,
  ArrowUpRight,
  Monitor,
  TicketPlus,
  TrendingUp,
} from 'lucide-react';
import { Link } from 'react-router-dom';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts';

export const DashboardPage: React.FC = () => {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [tickets, setTickets] = useState<QueueTicket[]>([]);
  const [branches, setBranches] = useState<Branch[]>([]);
  const [selectedBranchId, setSelectedBranchId] = useState<number>(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchBranches();
  }, []);

  useEffect(() => {
    loadDashboardData();
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

  const loadDashboardData = async () => {
    setLoading(true);
    try {
      const branchParam = selectedBranchId ? `?branch_id=${selectedBranchId}` : '';
      const statsRes = await fetchApi<DashboardStats>(`/reports/dashboard${branchParam}`);
      setStats(statsRes.data);

      if (selectedBranchId > 0) {
        const ticketRes = await fetchApi<QueueTicket[]>(`/tickets?branch_id=${selectedBranchId}`);
        setTickets(ticketRes.data || []);
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const formatSec = (sec: number) => {
    const m = Math.floor(sec / 60);
    const s = sec % 60;
    return `${m}m ${s}s`;
  };

  const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6', '#ec4899'];

  return (
    <div className="space-y-6">
      {/* Top Bar */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-black tracking-tight text-white">Queue Overview</h1>
          <p className="text-xs text-slate-400 mt-1">Realtime multi-branch operational telemetry</p>
        </div>

        <div className="flex items-center space-x-3">
          {branches.length > 0 && (
            <select
              value={selectedBranchId}
              onChange={(e) => setSelectedBranchId(Number(e.target.value))}
              className="bg-slate-900 border border-slate-800 rounded-xl px-3 py-2 text-xs text-slate-200 font-semibold focus:outline-none focus:border-blue-500"
            >
              {branches.map((b) => (
                <option key={b.id} value={b.id}>
                  {b.name} ({b.code})
                </option>
              ))}
            </select>
          )}

          <Link
            to="/counter"
            className="flex items-center space-x-2 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold px-4 py-2 rounded-xl transition shadow-lg shadow-blue-600/20"
          >
            <Kanban className="w-4 h-4" />
            <span>Staff Counter UI</span>
          </Link>
          <Link
            to="/reception"
            className="flex items-center space-x-2 bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-bold px-4 py-2 rounded-xl border border-slate-700 transition"
          >
            <TicketPlus className="w-4 h-4" />
            <span>Issue Ticket</span>
          </Link>
        </div>
      </div>

      {/* Metrics Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-slate-900/80 border border-slate-800/80 rounded-2xl p-5 relative overflow-hidden">
          <div className="flex items-center justify-between text-slate-400 mb-3">
            <span className="text-xs font-semibold">Total Issued Today</span>
            <Activity className="w-4 h-4 text-blue-400" />
          </div>
          <div className="text-3xl font-black text-white">{stats?.total_tickets_today || 0}</div>
          <div className="text-[11px] text-blue-400 font-medium mt-2 flex items-center">
            <TrendingUp className="w-3 h-3 mr-1" /> Realtime sequence count
          </div>
        </div>

        <div className="bg-slate-900/80 border border-slate-800/80 rounded-2xl p-5 relative overflow-hidden">
          <div className="flex items-center justify-between text-slate-400 mb-3">
            <span className="text-xs font-semibold">Currently Waiting</span>
            <Clock className="w-4 h-4 text-amber-400" />
          </div>
          <div className="text-3xl font-black text-amber-400">{stats?.waiting_count || 0}</div>
          <div className="text-[11px] text-slate-400 font-medium mt-2">
            Avg wait ~{formatSec(stats?.avg_wait_time_sec || 0)}
          </div>
        </div>

        <div className="bg-slate-900/80 border border-slate-800/80 rounded-2xl p-5 relative overflow-hidden">
          <div className="flex items-center justify-between text-slate-400 mb-3">
            <span className="text-xs font-semibold">Completed Today</span>
            <CheckCircle className="w-4 h-4 text-emerald-400" />
          </div>
          <div className="text-3xl font-black text-emerald-400">{stats?.completed_today || 0}</div>
          <div className="text-[11px] text-slate-400 font-medium mt-2">
            Avg service ~{formatSec(stats?.avg_service_time_sec || 0)}
          </div>
        </div>

        <div className="bg-slate-900/80 border border-slate-800/80 rounded-2xl p-5 relative overflow-hidden">
          <div className="flex items-center justify-between text-slate-400 mb-3">
            <span className="text-xs font-semibold">Active Counters</span>
            <Kanban className="w-4 h-4 text-indigo-400" />
          </div>
          <div className="text-3xl font-black text-indigo-400">{stats?.active_counters || 0}</div>
          <div className="text-[11px] text-slate-400 font-medium mt-2">
            {stats?.serving_count || 0} currently serving
          </div>
        </div>
      </div>

      {/* Analytics Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 bg-slate-900/80 border border-slate-800/80 rounded-2xl p-5">
          <h3 className="text-sm font-bold text-white mb-4">Tickets Issued Per Hour (Today)</h3>
          <div className="h-64">
            {stats?.hourly_distribution && stats.hourly_distribution.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={stats.hourly_distribution}>
                  <XAxis dataKey="hour" stroke="#64748b" fontSize={11} tickFormatter={(h) => `${h}:00`} />
                  <YAxis stroke="#64748b" fontSize={11} />
                  <Tooltip
                    contentStyle={{ backgroundColor: '#0f172a', borderColor: '#334155', borderRadius: '12px' }}
                    labelStyle={{ color: '#94a3b8' }}
                  />
                  <Bar dataKey="count" fill="#3b82f6" radius={[6, 6, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="h-full flex items-center justify-center text-xs text-slate-500">
                No hourly ticket data available for today yet.
              </div>
            )}
          </div>
        </div>

        <div className="bg-slate-900/80 border border-slate-800/80 rounded-2xl p-5">
          <h3 className="text-sm font-bold text-white mb-4">Service Volume Breakdown</h3>
          <div className="h-64 flex items-center justify-center">
            {stats?.service_distribution && stats.service_distribution.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={stats.service_distribution}
                    dataKey="count"
                    nameKey="service_name"
                    cx="50%"
                    cy="50%"
                    outerRadius={80}
                    label={({ service_name }) => service_name}
                  >
                    {stats.service_distribution.map((_, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip contentStyle={{ backgroundColor: '#0f172a', borderColor: '#334155', borderRadius: '12px' }} />
                </PieChart>
              </ResponsiveContainer>
            ) : (
              <div className="text-xs text-slate-500">No service distribution data</div>
            )}
          </div>
        </div>
      </div>

      {/* Live Ticket Stream Table */}
      <div className="bg-slate-900/80 border border-slate-800/80 rounded-2xl p-5">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-sm font-bold text-white">Live Queue Monitor (Branch)</h3>
          <span className="text-xs text-slate-400">Auto-refreshing</span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="text-slate-400 border-b border-slate-800 uppercase font-semibold text-[10px]">
              <tr>
                <th className="pb-3">Ticket #</th>
                <th className="pb-3">Service</th>
                <th className="pb-3">Counter</th>
                <th className="pb-3">Priority</th>
                <th className="pb-3">Status</th>
                <th className="pb-3">Est. Wait</th>
                <th className="pb-3 text-right">Track Link</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/60">
              {tickets.length === 0 ? (
                <tr>
                  <td colSpan={7} className="py-8 text-center text-slate-500">
                    No active tickets issued for this branch today.
                  </td>
                </tr>
              ) : (
                tickets.slice(0, 10).map((t) => (
                  <tr key={t.id} className="hover:bg-slate-800/30 transition">
                    <td className="py-3.5 font-extrabold text-blue-400 text-sm">{t.ticket_number}</td>
                    <td className="py-3.5 font-medium text-slate-200">{t.service_name}</td>
                    <td className="py-3.5 text-slate-300">{t.counter_number ? `Counter ${t.counter_number}` : '-'}</td>
                    <td className="py-3.5">
                      <span
                        className={`px-2 py-0.5 rounded-full text-[10px] font-bold ${
                          t.priority === 'EMERGENCY'
                            ? 'bg-red-500/20 text-red-400 border border-red-500/30'
                            : t.priority === 'PRIORITY'
                            ? 'bg-amber-500/20 text-amber-400 border border-amber-500/30'
                            : 'bg-slate-800 text-slate-400'
                        }`}
                      >
                        {t.priority}
                      </span>
                    </td>
                    <td className="py-3.5">
                      <span
                        className={`px-2.5 py-0.5 rounded-full text-[10px] font-bold ${
                          t.status === 'CALLED'
                            ? 'bg-amber-500/20 text-amber-300 animate-pulse'
                            : t.status === 'SERVING'
                            ? 'bg-blue-500/20 text-blue-300'
                            : t.status === 'COMPLETED'
                            ? 'bg-emerald-500/20 text-emerald-400'
                            : 'bg-slate-800 text-slate-300'
                        }`}
                      >
                        {t.status}
                      </span>
                    </td>
                    <td className="py-3.5 text-slate-400">{formatSec(t.estimated_wait_seconds)}</td>
                    <td className="py-3.5 text-right">
                      <Link
                        to={`/ticket/${t.public_token}`}
                        target="_blank"
                        className="inline-flex items-center space-x-1 text-blue-400 hover:text-blue-300 font-semibold"
                      >
                        <span>Public Link</span>
                        <ArrowUpRight className="w-3 h-3" />
                      </Link>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
