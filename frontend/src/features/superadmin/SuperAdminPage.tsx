import React, { useState } from 'react';
import { 
  Building2, 
  ShieldCheck, 
  Users, 
  Crown, 
  Sliders, 
  FileText, 
  Plus, 
  CheckCircle, 
  AlertTriangle,
  Search,
  Check,
  X
} from 'lucide-react';
import { useAuthStore } from '../../stores/authStore';

interface TenantOrg {
  id: number;
  name: string;
  code: string;
  plan: 'FREE' | 'PRO' | 'ENTERPRISE';
  branches: number;
  maxBranches: number;
  counters: number;
  maxCounters: number;
  status: 'ACTIVE' | 'SUSPENDED';
  createdAt: string;
}

interface PlatformUser {
  id: number;
  name: string;
  email: string;
  role: string;
  orgName: string;
  status: string;
}

export const SuperAdminPage: React.FC = () => {
  const { user } = useAuthStore();
  const [activeTab, setActiveTab] = useState<'tenants' | 'users' | 'master' | 'audit'>('tenants');
  const [searchQuery, setSearchQuery] = useState('');

  // Initial Mock State for Master Data Management
  const [organizations, setOrganizations] = useState<TenantOrg[]>([
    {
      id: 1,
      name: 'Demo Healthcare Org',
      code: 'DEMO',
      plan: 'ENTERPRISE',
      branches: 3,
      maxBranches: 10,
      counters: 12,
      maxCounters: 50,
      status: 'ACTIVE',
      createdAt: '2026-01-15',
    },
    {
      id: 2,
      name: 'BCA Main Branch Queue',
      code: 'BCA01',
      plan: 'PRO',
      branches: 5,
      maxBranches: 5,
      counters: 20,
      maxCounters: 25,
      status: 'ACTIVE',
      createdAt: '2026-02-01',
    },
    {
      id: 3,
      name: 'Klinik Medika Pratama',
      code: 'MEDIKA',
      plan: 'FREE',
      branches: 1,
      maxBranches: 1,
      counters: 2,
      maxCounters: 2,
      status: 'ACTIVE',
      createdAt: '2026-02-10',
    },
  ]);

  const [platformUsers, setPlatformUsers] = useState<PlatformUser[]>([
    { id: 1, name: 'Super Admin Master', email: 'superadmin@system.com', role: 'SUPER_ADMIN', orgName: 'System Global', status: 'ACTIVE' },
    { id: 2, name: 'Dr. Sarah Connor', email: 'owner@healthcare.com', role: 'OWNER', orgName: 'Demo Healthcare Org', status: 'ACTIVE' },
    { id: 3, name: 'Alex Mercer', email: 'manager@healthcare.com', role: 'MANAGER', orgName: 'Demo Healthcare Org', status: 'ACTIVE' },
    { id: 4, name: 'John Doe', email: 'staff@healthcare.com', role: 'STAFF', orgName: 'Demo Healthcare Org', status: 'ACTIVE' },
    { id: 5, name: 'Emily Watson', email: 'receptionist@healthcare.com', role: 'RECEPTIONIST', orgName: 'Demo Healthcare Org', status: 'ACTIVE' },
  ]);

  const [isAddOrgModalOpen, setIsAddOrgModalOpen] = useState(false);
  const [newOrgName, setNewOrgName] = useState('');
  const [newOrgCode, setNewOrgCode] = useState('');
  const [newOrgPlan, setNewOrgPlan] = useState<'FREE' | 'PRO' | 'ENTERPRISE'>('PRO');

  const handleAddOrg = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newOrgName || !newOrgCode) return;

    const maxB = newOrgPlan === 'ENTERPRISE' ? 20 : newOrgPlan === 'PRO' ? 5 : 1;
    const maxC = newOrgPlan === 'ENTERPRISE' ? 100 : newOrgPlan === 'PRO' ? 25 : 2;

    const newOrg: TenantOrg = {
      id: Date.now(),
      name: newOrgName,
      code: newOrgCode.toUpperCase(),
      plan: newOrgPlan,
      branches: 1,
      maxBranches: maxB,
      counters: 1,
      maxCounters: maxC,
      status: 'ACTIVE',
      createdAt: new Date().toISOString().split('T')[0],
    };

    setOrganizations([newOrg, ...organizations]);
    setNewOrgName('');
    setNewOrgCode('');
    setIsAddOrgModalOpen(false);
  };

  const toggleOrgStatus = (id: number) => {
    setOrganizations(organizations.map(org => {
      if (org.id === id) {
        return {
          ...org,
          status: org.status === 'ACTIVE' ? 'SUSPENDED' : 'ACTIVE'
        };
      }
      return org;
    }));
  };

  const filteredOrgs = organizations.filter(org => 
    org.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
    org.code.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="space-y-6">
      {/* Super Admin Banner */}
      <div className="bg-gradient-to-r from-purple-900 via-indigo-900 to-slate-900 p-6 rounded-2xl text-white shadow-xl flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <div className="p-3 bg-purple-500/20 backdrop-blur-md rounded-xl border border-purple-400/30">
            <ShieldCheck className="w-10 h-10 text-purple-400" />
          </div>
          <div>
            <div className="flex items-center space-x-2">
              <h1 className="text-2xl font-extrabold tracking-tight">Super Admin Master Control</h1>
              <span className="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-purple-500/30 text-purple-200 border border-purple-400/30">
                System Administrator
              </span>
            </div>
            <p className="text-purple-200 text-sm mt-1">
              Global Platform Tenant Management, Master Configurations, and Audit Oversight
            </p>
          </div>
        </div>

        <button
          onClick={() => setIsAddOrgModalOpen(true)}
          className="flex items-center space-x-2 bg-purple-600 hover:bg-purple-500 text-white font-medium px-4 py-2.5 rounded-xl shadow-lg transition-all transform hover:-translate-y-0.5"
        >
          <Plus className="w-5 h-5" />
          <span>New Tenant Org</span>
        </button>
      </div>

      {/* Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-5">
        <div className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-slate-500">Total Organizations</p>
            <h3 className="text-2xl font-extrabold text-slate-800 mt-1">{organizations.length}</h3>
            <p className="text-xs text-emerald-600 mt-1 flex items-center">
              <CheckCircle className="w-3.5 h-3.5 mr-1" /> All system tenants
            </p>
          </div>
          <div className="w-12 h-12 bg-blue-50 text-blue-600 rounded-xl flex items-center justify-center">
            <Building2 className="w-6 h-6" />
          </div>
        </div>

        <div className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-slate-500">Enterprise Plans</p>
            <h3 className="text-2xl font-extrabold text-purple-700 mt-1">
              {organizations.filter(o => o.plan === 'ENTERPRISE').length}
            </h3>
            <p className="text-xs text-purple-600 mt-1 flex items-center">
              <Crown className="w-3.5 h-3.5 mr-1" /> Unlimited capacity
            </p>
          </div>
          <div className="w-12 h-12 bg-purple-50 text-purple-600 rounded-xl flex items-center justify-center">
            <Crown className="w-6 h-6" />
          </div>
        </div>

        <div className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-slate-500">Platform Users</p>
            <h3 className="text-2xl font-extrabold text-slate-800 mt-1">{platformUsers.length}</h3>
            <p className="text-xs text-slate-500 mt-1">Across all roles</p>
          </div>
          <div className="w-12 h-12 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center">
            <Users className="w-6 h-6" />
          </div>
        </div>

        <div className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-slate-500">System Status</p>
            <h3 className="text-2xl font-extrabold text-emerald-600 mt-1">HEALTHY</h3>
            <p className="text-xs text-emerald-600 mt-1">Multi-tenant isolates active</p>
          </div>
          <div className="w-12 h-12 bg-emerald-50 text-emerald-600 rounded-xl flex items-center justify-center">
            <ShieldCheck className="w-6 h-6" />
          </div>
        </div>
      </div>

      {/* Tabs Navigation */}
      <div className="border-b border-slate-200 flex space-x-6">
        <button
          onClick={() => setActiveTab('tenants')}
          className={`pb-3 text-sm font-semibold flex items-center space-x-2 transition-colors border-b-2 ${
            activeTab === 'tenants' 
              ? 'border-purple-600 text-purple-600' 
              : 'border-transparent text-slate-500 hover:text-slate-800'
          }`}
        >
          <Building2 className="w-4 h-4" />
          <span>Tenant Organizations</span>
        </button>

        <button
          onClick={() => setActiveTab('users')}
          className={`pb-3 text-sm font-semibold flex items-center space-x-2 transition-colors border-b-2 ${
            activeTab === 'users' 
              ? 'border-purple-600 text-purple-600' 
              : 'border-transparent text-slate-500 hover:text-slate-800'
          }`}
        >
          <Users className="w-4 h-4" />
          <span>User Master</span>
        </button>

        <button
          onClick={() => setActiveTab('master')}
          className={`pb-3 text-sm font-semibold flex items-center space-x-2 transition-colors border-b-2 ${
            activeTab === 'master' 
              ? 'border-purple-600 text-purple-600' 
              : 'border-transparent text-slate-500 hover:text-slate-800'
          }`}
        >
          <Sliders className="w-4 h-4" />
          <span>Global Configurations</span>
        </button>

        <button
          onClick={() => setActiveTab('audit')}
          className={`pb-3 text-sm font-semibold flex items-center space-x-2 transition-colors border-b-2 ${
            activeTab === 'audit' 
              ? 'border-purple-600 text-purple-600' 
              : 'border-transparent text-slate-500 hover:text-slate-800'
          }`}
        >
          <FileText className="w-4 h-4" />
          <span>System Audit Log</span>
        </button>
      </div>

      {/* TAB 1: Tenants List */}
      {activeTab === 'tenants' && (
        <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-bold text-slate-800">Master Organizations & Subscriptions</h2>
            <div className="relative w-64">
              <Search className="w-4 h-4 absolute left-3 top-3 text-slate-400" />
              <input
                type="text"
                placeholder="Search org name or code..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-purple-500"
              />
            </div>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-600">
              <thead className="bg-slate-50 text-slate-500 font-semibold border-b border-slate-200 uppercase text-xs">
                <tr>
                  <th className="py-3 px-4">Organization Name</th>
                  <th className="py-3 px-4">Code</th>
                  <th className="py-3 px-4">Subscription Plan</th>
                  <th className="py-3 px-4">Branch Limit</th>
                  <th className="py-3 px-4">Counter Limit</th>
                  <th className="py-3 px-4">Status</th>
                  <th className="py-3 px-4 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {filteredOrgs.map((org) => (
                  <tr key={org.id} className="hover:bg-slate-50/50">
                    <td className="py-4 px-4 font-semibold text-slate-900">{org.name}</td>
                    <td className="py-4 px-4 font-mono text-slate-600">{org.code}</td>
                    <td className="py-4 px-4">
                      <span className={`px-2.5 py-1 rounded-full text-xs font-bold ${
                        org.plan === 'ENTERPRISE' ? 'bg-purple-100 text-purple-700' :
                        org.plan === 'PRO' ? 'bg-blue-100 text-blue-700' : 'bg-slate-100 text-slate-700'
                      }`}>
                        {org.plan}
                      </span>
                    </td>
                    <td className="py-4 px-4 font-medium">{org.branches} / {org.maxBranches}</td>
                    <td className="py-4 px-4 font-medium">{org.counters} / {org.maxCounters}</td>
                    <td className="py-4 px-4">
                      <span className={`px-2.5 py-1 rounded-full text-xs font-semibold ${
                        org.status === 'ACTIVE' ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'
                      }`}>
                        {org.status}
                      </span>
                    </td>
                    <td className="py-4 px-4 text-right">
                      <button
                        onClick={() => toggleOrgStatus(org.id)}
                        className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors ${
                          org.status === 'ACTIVE'
                            ? 'bg-amber-50 text-amber-700 hover:bg-amber-100'
                            : 'bg-emerald-50 text-emerald-700 hover:bg-emerald-100'
                        }`}
                      >
                        {org.status === 'ACTIVE' ? 'Suspend' : 'Activate'}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* TAB 2: Users Master */}
      {activeTab === 'users' && (
        <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 space-y-4">
          <h2 className="text-lg font-bold text-slate-800">Global Users Administration</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-600">
              <thead className="bg-slate-50 text-slate-500 font-semibold border-b border-slate-200 uppercase text-xs">
                <tr>
                  <th className="py-3 px-4">Full Name</th>
                  <th className="py-3 px-4">Email</th>
                  <th className="py-3 px-4">Role</th>
                  <th className="py-3 px-4">Tenant Org</th>
                  <th className="py-3 px-4">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {platformUsers.map((u) => (
                  <tr key={u.id} className="hover:bg-slate-50/50">
                    <td className="py-4 px-4 font-semibold text-slate-900">{u.name}</td>
                    <td className="py-4 px-4 font-mono text-slate-600">{u.email}</td>
                    <td className="py-4 px-4">
                      <span className={`px-2.5 py-1 rounded-full text-xs font-bold ${
                        u.role === 'SUPER_ADMIN' ? 'bg-purple-100 text-purple-700 border border-purple-300' :
                        u.role === 'OWNER' ? 'bg-indigo-100 text-indigo-700' : 'bg-slate-100 text-slate-700'
                      }`}>
                        {u.role}
                      </span>
                    </td>
                    <td className="py-4 px-4">{u.orgName}</td>
                    <td className="py-4 px-4">
                      <span className="px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-100 text-emerald-700">
                        {u.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* TAB 3: Global Configurations */}
      {activeTab === 'master' && (
        <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 space-y-6">
          <h2 className="text-lg font-bold text-slate-800">Master Platform Defaults & Limits</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="p-5 border border-slate-200 rounded-xl space-y-3">
              <h3 className="font-semibold text-slate-800 text-sm">Free Plan Limits</h3>
              <div className="space-y-2 text-xs text-slate-600">
                <div className="flex justify-between"><span>Max Branches:</span><span className="font-bold">1</span></div>
                <div className="flex justify-between"><span>Max Counters:</span><span className="font-bold">2</span></div>
                <div className="flex justify-between"><span>Max Monthly Tickets:</span><span className="font-bold">1,000</span></div>
              </div>
            </div>

            <div className="p-5 border border-purple-200 bg-purple-50/50 rounded-xl space-y-3">
              <h3 className="font-semibold text-purple-900 text-sm">Enterprise Plan Defaults</h3>
              <div className="space-y-2 text-xs text-purple-800">
                <div className="flex justify-between"><span>Max Branches:</span><span className="font-bold">20+</span></div>
                <div className="flex justify-between"><span>Max Counters:</span><span className="font-bold">100</span></div>
                <div className="flex justify-between"><span>Max Monthly Tickets:</span><span className="font-bold">100,000</span></div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* TAB 4: Audit Logs */}
      {activeTab === 'audit' && (
        <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 space-y-4">
          <h2 className="text-lg font-bold text-slate-800">System Platform Audit Trail</h2>
          <p className="text-sm text-slate-500">Live security & configuration changes across all organization tenants.</p>
          <div className="border border-slate-200 rounded-xl p-4 bg-slate-50 font-mono text-xs space-y-2 text-slate-700">
            <p>[2026-08-19 14:54:09 UTC] USER 99 (SUPER_ADMIN) initialized platform system controls.</p>
            <p>[2026-08-19 14:58:33 UTC] DB Config updated: host port mapped to 5436.</p>
            <p>[2026-08-19 15:10:00 UTC] Organization ID 1 (Demo Healthcare Org) subscription updated to ENTERPRISE.</p>
          </div>
        </div>
      )}

      {/* Modal Add Org */}
      {isAddOrgModalOpen && (
        <div className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-5 animate-in fade-in zoom-in duration-200">
            <div className="flex items-center justify-between border-b border-slate-100 pb-3">
              <h3 className="text-lg font-bold text-slate-800 flex items-center space-x-2">
                <Building2 className="w-5 h-5 text-purple-600" />
                <span>Create New Tenant Org</span>
              </h3>
              <button 
                onClick={() => setIsAddOrgModalOpen(false)}
                className="text-slate-400 hover:text-slate-600 rounded-lg p-1"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleAddOrg} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-700 uppercase mb-1">Organization Name</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. RS Siloam Utama"
                  value={newOrgName}
                  onChange={(e) => setNewOrgName(e.target.value)}
                  className="w-full px-3.5 py-2 border border-slate-200 rounded-xl text-sm focus:ring-2 focus:ring-purple-500 focus:outline-none"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 uppercase mb-1">Organization Code</label>
                <input
                  type="text"
                  required
                  placeholder="e.g. SILOAM"
                  value={newOrgCode}
                  onChange={(e) => setNewOrgCode(e.target.value)}
                  className="w-full px-3.5 py-2 border border-slate-200 rounded-xl text-sm focus:ring-2 focus:ring-purple-500 focus:outline-none uppercase"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 uppercase mb-1">Subscription Plan</label>
                <select
                  value={newOrgPlan}
                  onChange={(e) => setNewOrgPlan(e.target.value as any)}
                  className="w-full px-3.5 py-2 border border-slate-200 rounded-xl text-sm focus:ring-2 focus:ring-purple-500 focus:outline-none"
                >
                  <option value="FREE">FREE (1 Branch, 2 Counters)</option>
                  <option value="PRO">PRO (5 Branches, 25 Counters)</option>
                  <option value="ENTERPRISE">ENTERPRISE (20 Branches, 100 Counters)</option>
                </select>
              </div>

              <div className="flex justify-end space-x-3 pt-3">
                <button
                  type="button"
                  onClick={() => setIsAddOrgModalOpen(false)}
                  className="px-4 py-2 border border-slate-200 rounded-xl text-sm font-medium text-slate-600 hover:bg-slate-50"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-purple-600 hover:bg-purple-500 text-white rounded-xl text-sm font-medium shadow-md transition-colors"
                >
                  Create Organization
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
