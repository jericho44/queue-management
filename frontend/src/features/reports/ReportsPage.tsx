import React, { useEffect, useState } from 'react';
import { fetchApi } from '../../api/client';
import { AuditLog } from '../../types';
import { FileText, Shield, Search, RefreshCw } from 'lucide-react';

export const ReportsPage: React.FC = () => {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchAuditLogs();
  }, []);

  const fetchAuditLogs = async () => {
    setLoading(true);
    try {
      const res = await fetchApi<AuditLog[]>('/audit-logs');
      setLogs(res.data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black text-white">Audit & Event Logs</h1>
          <p className="text-xs text-slate-400 mt-1">Immutable security and state machine audit trail</p>
        </div>
        <button
          onClick={fetchAuditLogs}
          className="flex items-center space-x-2 px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold border border-slate-700 transition"
        >
          <RefreshCw className="w-3.5 h-3.5" />
          <span>Refresh</span>
        </button>
      </div>

      <div className="bg-slate-900/90 border border-slate-800 rounded-2xl p-5 shadow-xl">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="text-slate-400 border-b border-slate-800 uppercase font-semibold text-[10px]">
              <tr>
                <th className="pb-3">Timestamp</th>
                <th className="pb-3">User</th>
                <th className="pb-3">Action</th>
                <th className="pb-3">Entity</th>
                <th className="pb-3">IP Address</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/60">
              {logs.length === 0 ? (
                <tr>
                  <td colSpan={5} className="py-8 text-center text-slate-500">
                    No audit log records recorded yet.
                  </td>
                </tr>
              ) : (
                logs.map((l) => (
                  <tr key={l.id} className="hover:bg-slate-800/40 transition">
                    <td className="py-3 text-slate-400 font-mono text-[11px]">
                      {new Date(l.created_at).toLocaleString()}
                    </td>
                    <td className="py-3 font-semibold text-slate-200">{l.user_name || 'System'}</td>
                    <td className="py-3 font-bold text-blue-400">{l.action}</td>
                    <td className="py-3 text-slate-300">
                      {l.entity_type} #{l.entity_id || '-'}
                    </td>
                    <td className="py-3 text-slate-400 font-mono text-[11px]">{l.ip_address || '127.0.0.1'}</td>
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
