import React, { useEffect, useState } from 'react';
import { fetchApi } from '../../api/client';
import { Branch, Service, Counter } from '../../types';
import { Building, Plus, Layers, Kanban, CheckCircle, ShieldAlert, Printer, Settings2 } from 'lucide-react';
import { KioskSettingsModal } from './KioskSettingsModal';
import { useToast } from '../../components/Toast';

export const BranchesPage: React.FC = () => {
  const { showError, showSuccess } = useToast();
  const [branches, setBranches] = useState<Branch[]>([]);
  const [services, setServices] = useState<Service[]>([]);
  const [counters, setCounters] = useState<Counter[]>([]);
  const [selectedBranchId, setSelectedBranchId] = useState<number>(0);
  const [showKioskModal, setShowKioskModal] = useState(false);

  // New Branch Form State
  const [newBranchName, setNewBranchName] = useState('');
  const [newBranchCode, setNewBranchCode] = useState('');

  // New Service Form State
  const [newServiceName, setNewServiceName] = useState('');
  const [newServicePrefix, setNewServicePrefix] = useState('A');

  // New Counter Form State
  const [newCounterNum, setNewCounterNum] = useState('');
  const [newCounterName, setNewCounterName] = useState('');

  const selectedBranch = branches.find((b) => b.id === selectedBranchId);

  useEffect(() => {
    fetchBranches();
  }, []);

  useEffect(() => {
    if (selectedBranchId > 0) {
      fetchBranchData(selectedBranchId);
    }
  }, [selectedBranchId]);

  const fetchBranches = async () => {
    try {
      const res = await fetchApi<Branch[]>('/branches');
      setBranches(res.data || []);
      if (res.data && res.data.length > 0 && selectedBranchId === 0) {
        setSelectedBranchId(res.data[0].id);
      }
    } catch (err: any) {
      showError(err.message || 'Gagal memuat daftar cabang');
    }
  };

  const fetchBranchData = async (branchId: number) => {
    try {
      const [sRes, cRes] = await Promise.all([
        fetchApi<Service[]>(`/services?branch_id=${branchId}`),
        fetchApi<Counter[]>(`/counters?branch_id=${branchId}`),
      ]);
      setServices(sRes.data || []);
      setCounters(cRes.data || []);
    } catch (err: any) {
      showError(err.message || 'Gagal memuat detail cabang');
    }
  };

  const handleCreateBranch = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await fetchApi('/branches', {
        method: 'POST',
        body: JSON.stringify({ name: newBranchName, code: newBranchCode }),
      });
      showSuccess(`Cabang "${newBranchName}" berhasil dibuat`);
      setNewBranchName('');
      setNewBranchCode('');
      fetchBranches();
    } catch (err: any) {
      showError(err.message || 'Gagal membuat cabang baru');
    }
  };

  const handleCreateService = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedBranchId) return;
    try {
      await fetchApi('/services', {
        method: 'POST',
        body: JSON.stringify({
          branch_id: selectedBranchId,
          name: newServiceName,
          code: newServiceName.toUpperCase().replace(/\s+/g, '_'),
          prefix: newServicePrefix,
        }),
      });
      showSuccess(`Layanan "${newServiceName}" berhasil ditambahkan`);
      setNewServiceName('');
      fetchBranchData(selectedBranchId);
    } catch (err: any) {
      showError(err.message || 'Gagal menambahkan layanan');
    }
  };

  const handleCreateCounter = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedBranchId) return;
    try {
      const serviceIDs = services.map((s) => s.id);
      await fetchApi('/counters', {
        method: 'POST',
        body: JSON.stringify({
          branch_id: selectedBranchId,
          counter_number: newCounterNum,
          name: newCounterName,
          service_ids: serviceIDs,
        }),
      });
      showSuccess(`Loket "${newCounterName}" berhasil dibuat`);
      setNewCounterNum('');
      setNewCounterName('');
      fetchBranchData(selectedBranchId);
    } catch (err: any) {
      showError(err.message || 'Gagal membuat loket');
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-black text-white">Branch & Counter Management</h1>
          <p className="text-xs text-slate-400 mt-1">Configure multi-tenant locations, services, and counter assignments</p>
        </div>

        {selectedBranch && (
          <button
            onClick={() => setShowKioskModal(true)}
            className="px-4 py-2.5 bg-slate-800 hover:bg-slate-700 text-blue-400 hover:text-blue-300 border border-slate-700 rounded-xl text-xs font-bold transition flex items-center space-x-2 shadow-lg"
          >
            <Printer className="w-4 h-4" />
            <span>Kiosk & Printer Settings</span>
          </button>
        )}
      </div>

      {showKioskModal && selectedBranch && (
        <KioskSettingsModal
          branch={selectedBranch}
          onClose={() => setShowKioskModal(false)}
          onSaveSuccess={fetchBranches}
        />
      )}


      {/* Branch Selector */}
      <div className="flex items-center space-x-3 overflow-x-auto pb-2">
        {branches.map((b) => (
          <button
            key={b.id}
            onClick={() => setSelectedBranchId(b.id)}
            className={`px-4 py-2.5 rounded-xl border text-xs font-bold transition flex items-center space-x-2 whitespace-nowrap ${
              selectedBranchId === b.id
                ? 'bg-blue-600 border-blue-500 text-white shadow-lg shadow-blue-600/20'
                : 'bg-slate-900 border-slate-800 text-slate-400 hover:text-slate-200'
            }`}
          >
            <Building className="w-4 h-4" />
            <span>{b.name} ({b.code})</span>
          </button>
        ))}
      </div>

      {/* Grid Layout for Configuration */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Create Branch */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-2xl p-5 shadow-xl space-y-4">
          <div className="flex items-center space-x-2 text-sm font-bold text-white border-b border-slate-800 pb-3">
            <Building className="w-4 h-4 text-blue-400" />
            <span>Add New Branch</span>
          </div>

          <form onSubmit={handleCreateBranch} className="space-y-3">
            <div>
              <label className="block text-[11px] font-semibold text-slate-300 mb-1">Branch Name</label>
              <input
                type="text"
                value={newBranchName}
                onChange={(e) => setNewBranchName(e.target.value)}
                required
                placeholder="Downtown Main Branch"
                className="w-full bg-slate-800 border border-slate-700 rounded-xl px-3 py-2 text-xs text-white focus:outline-none focus:border-blue-500"
              />
            </div>
            <div>
              <label className="block text-[11px] font-semibold text-slate-300 mb-1">Branch Code</label>
              <input
                type="text"
                value={newBranchCode}
                onChange={(e) => setNewBranchCode(e.target.value.toUpperCase())}
                required
                placeholder="MAIN"
                className="w-full bg-slate-800 border border-slate-700 rounded-xl px-3 py-2 text-xs text-white focus:outline-none focus:border-blue-500"
              />
            </div>
            <button
              type="submit"
              className="w-full py-2.5 bg-blue-600 hover:bg-blue-500 text-white font-bold text-xs rounded-xl shadow-md transition"
            >
              Save Branch
            </button>
          </form>
        </div>

        {/* Services List & Create */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-2xl p-5 shadow-xl space-y-4">
          <div className="flex items-center space-x-2 text-sm font-bold text-white border-b border-slate-800 pb-3">
            <Layers className="w-4 h-4 text-indigo-400" />
            <span>Services ({services.length})</span>
          </div>

          <form onSubmit={handleCreateService} className="space-y-3">
            <div>
              <input
                type="text"
                value={newServiceName}
                onChange={(e) => setNewServiceName(e.target.value)}
                required
                placeholder="New Service Name (e.g. Pharmacy)"
                className="w-full bg-slate-800 border border-slate-700 rounded-xl px-3 py-2 text-xs text-white focus:outline-none focus:border-blue-500"
              />
            </div>
            <div className="flex items-center space-x-2">
              <select
                value={newServicePrefix}
                onChange={(e) => setNewServicePrefix(e.target.value)}
                className="bg-slate-800 border border-slate-700 rounded-xl px-3 py-2 text-xs text-white font-bold"
              >
                <option value="A">Prefix A</option>
                <option value="B">Prefix B</option>
                <option value="C">Prefix C</option>
                <option value="P">Prefix P</option>
              </select>
              <button
                type="submit"
                className="flex-1 py-2 bg-indigo-600 hover:bg-indigo-500 text-white font-bold text-xs rounded-xl transition"
              >
                Add Service
              </button>
            </div>
          </form>

          <div className="space-y-2 max-h-48 overflow-y-auto pt-2">
            {services.map((s) => (
              <div key={s.id} className="bg-slate-800/40 p-2.5 rounded-xl border border-slate-700/60 flex items-center justify-between text-xs">
                <div>
                  <span className="font-bold text-white">{s.name}</span>
                  <span className="text-[10px] text-slate-400 block">Avg ~{Math.floor(s.avg_service_time_sec / 60)} min</span>
                </div>
                <span className="font-extrabold text-indigo-400 bg-indigo-500/10 px-2 py-0.5 rounded-md border border-indigo-500/20">
                  {s.prefix}
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Counters List & Create */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-2xl p-5 shadow-xl space-y-4">
          <div className="flex items-center space-x-2 text-sm font-bold text-white border-b border-slate-800 pb-3">
            <Kanban className="w-4 h-4 text-emerald-400" />
            <span>Counters ({counters.length})</span>
          </div>

          <form onSubmit={handleCreateCounter} className="space-y-3">
            <div className="grid grid-cols-2 gap-2">
              <input
                type="text"
                value={newCounterNum}
                onChange={(e) => setNewCounterNum(e.target.value)}
                required
                placeholder="Num (e.g. 01)"
                className="bg-slate-800 border border-slate-700 rounded-xl px-3 py-2 text-xs text-white focus:outline-none focus:border-blue-500"
              />
              <input
                type="text"
                value={newCounterName}
                onChange={(e) => setNewCounterName(e.target.value)}
                required
                placeholder="Counter Name"
                className="bg-slate-800 border border-slate-700 rounded-xl px-3 py-2 text-xs text-white focus:outline-none focus:border-blue-500"
              />
            </div>
            <button
              type="submit"
              className="w-full py-2 bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs rounded-xl transition"
            >
              Add Counter
            </button>
          </form>

          <div className="space-y-2 max-h-48 overflow-y-auto pt-2">
            {counters.map((c) => (
              <div key={c.id} className="bg-slate-800/40 p-2.5 rounded-xl border border-slate-700/60 flex items-center justify-between text-xs">
                <div>
                  <span className="font-bold text-white">Counter {c.counter_number}</span>
                  <span className="text-[10px] text-slate-400 block">{c.name}</span>
                </div>
                <span
                  className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${
                    c.status === 'OPEN'
                      ? 'bg-emerald-500/20 text-emerald-400'
                      : c.status === 'BUSY'
                      ? 'bg-amber-500/20 text-amber-400'
                      : 'bg-slate-700 text-slate-400'
                  }`}
                >
                  {c.status}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};
