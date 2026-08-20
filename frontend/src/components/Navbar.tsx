import React from 'react';
import { useAuthStore } from '../stores/authStore';
import { LogOut, User as UserIcon, Building2, Monitor, Ticket, Layers } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';

export const Navbar: React.FC = () => {
  const { user, organization, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <header className="h-16 bg-slate-900/80 backdrop-blur border-b border-slate-800 px-6 flex items-center justify-between sticky top-0 z-40">
      <div className="flex items-center space-x-3">
        <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-500 flex items-center justify-center text-white font-bold text-xl shadow-lg shadow-blue-500/20">
          Q
        </div>
        <div>
          <span className="font-extrabold text-lg text-white tracking-tight">QFlow</span>
          <span className="ml-2 text-xs font-semibold px-2 py-0.5 rounded-full bg-blue-500/10 text-blue-400 border border-blue-500/20">
            SaaS Multi-Tenant
          </span>
        </div>
      </div>

      <div className="flex items-center space-x-4">
        {organization && (
          <div className="hidden sm:flex items-center space-x-2 text-xs bg-slate-800/80 px-3 py-1.5 rounded-lg border border-slate-700/60 text-slate-300">
            <Building2 className="w-3.5 h-3.5 text-blue-400" />
            <span className="font-medium text-white">{organization.name}</span>
            <span className="text-slate-500">({organization.code})</span>
          </div>
        )}

        <Link
          to="/display"
          target="_blank"
          className="flex items-center space-x-1.5 text-xs bg-indigo-600/20 hover:bg-indigo-600/30 text-indigo-300 border border-indigo-500/30 px-3 py-1.5 rounded-lg transition"
        >
          <Monitor className="w-3.5 h-3.5 text-indigo-400" />
          <span className="font-semibold">Public Display</span>
        </Link>

        {user && (
          <div className="flex items-center space-x-3 pl-3 border-l border-slate-800">
            <div className="text-right hidden md:block">
              <div className="text-xs font-bold text-slate-100">{user.full_name}</div>
              <div className="text-[10px] text-blue-400 font-medium">{user.role}</div>
            </div>
            <button
              onClick={handleLogout}
              className="p-2 rounded-lg bg-slate-800 hover:bg-red-500/10 text-slate-400 hover:text-red-400 border border-slate-700/60 hover:border-red-500/30 transition"
              title="Sign Out"
            >
              <LogOut className="w-4 h-4" />
            </button>
          </div>
        )}
      </div>
    </header>
  );
};
