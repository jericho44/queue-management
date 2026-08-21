import React from 'react';
import { NavLink } from 'react-router-dom';
import {
  LayoutDashboard,
  Kanban,
  TicketPlus,
  Monitor,
  Building,
  Settings,
  FileText,
  Users,
  ShieldAlert,
  CreditCard,
  DollarSign,
} from 'lucide-react';
import { useAuthStore } from '../stores/authStore';

export const Sidebar: React.FC = () => {
  const { user } = useAuthStore();
  const role = user?.role || 'STAFF';

  const links = [
    { to: '/', label: 'Overview', icon: LayoutDashboard, roles: ['SUPER_ADMIN', 'OWNER', 'MANAGER', 'STAFF', 'RECEPTIONIST'] },
    { to: '/superadmin', label: 'Super Admin Console', icon: ShieldAlert, roles: ['SUPER_ADMIN'] },
    { to: '/superadmin/billing', label: 'SaaS Revenue & Billing', icon: DollarSign, roles: ['SUPER_ADMIN'] },
    { to: '/counter', label: 'Counter Station', icon: Kanban, roles: ['SUPER_ADMIN', 'OWNER', 'MANAGER', 'STAFF'] },
    { to: '/reception', label: 'Ticket Reception', icon: TicketPlus, roles: ['SUPER_ADMIN', 'OWNER', 'MANAGER', 'RECEPTIONIST'] },
    { to: '/branches', label: 'Branches & Counters', icon: Building, roles: ['SUPER_ADMIN', 'OWNER', 'MANAGER'] },
    { to: '/billing', label: 'Billing & Invoices', icon: CreditCard, roles: ['SUPER_ADMIN', 'OWNER', 'MANAGER'] },
    { to: '/reports', label: 'Analytics & Audit', icon: FileText, roles: ['SUPER_ADMIN', 'OWNER', 'MANAGER'] },
  ];


  return (
    <aside className="w-64 bg-slate-900/60 border-r border-slate-800 p-4 flex flex-col justify-between hidden md:flex">
      <div className="space-y-1">
        <div className="px-3 py-2 text-[11px] font-bold tracking-wider text-slate-500 uppercase">
          Navigation Menu
        </div>
        {links.map((link) => {
          if (!link.roles.includes(role)) return null;
          const Icon = link.icon;
          return (
            <NavLink
              key={link.to}
              to={link.to}
              end={link.to === '/'}
              className={({ isActive }) =>
                `flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-xs font-medium transition ${
                  isActive
                    ? 'bg-blue-600/15 text-blue-400 border border-blue-500/20 font-semibold shadow-sm'
                    : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'
                }`
              }
            >
              <Icon className="w-4 h-4" />
              <span>{link.label}</span>
            </NavLink>
          );
        })}
      </div>

      <div className="pt-4 border-t border-slate-800">
        <div className="bg-slate-800/40 rounded-xl p-3 border border-slate-700/50">
          <div className="flex items-center space-x-2 text-xs font-semibold text-slate-300">
            <ShieldAlert className="w-4 h-4 text-emerald-400" />
            <span>Isolation Active</span>
          </div>
          <p className="text-[10px] text-slate-400 mt-1">
            Tenant context is verified on every request.
          </p>
        </div>
      </div>
    </aside>
  );
};
