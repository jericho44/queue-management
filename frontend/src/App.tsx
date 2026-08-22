import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { DashboardLayout } from './layouts/DashboardLayout';
import { LoginPage } from './features/auth/LoginPage';
import { RegisterOrgPage } from './features/auth/RegisterOrgPage';
import { DashboardPage } from './features/dashboard/DashboardPage';
import { StaffCounterPage } from './features/counter/StaffCounterPage';
import { ReceptionPage } from './features/reception/ReceptionPage';
import { PublicDisplayPage } from './features/display/PublicDisplayPage';
import { PublicTicketPage } from './features/customer/PublicTicketPage';
import { KioskPage } from './features/kiosk/KioskPage';
import { BranchesPage } from './features/branches/BranchesPage';
import { ReportsPage } from './features/reports/ReportsPage';
import { SuperAdminPage } from './features/superadmin/SuperAdminPage';
import { BillingPortalPage } from './features/billing/BillingPortalPage';
import { SuperadminBillingPage } from './features/billing/SuperadminBillingPage';

import { ToastProvider } from './components/Toast';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

export const App: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <ToastProvider>
        <BrowserRouter>
          <Routes>
            {/* Public Authentication & Customer Pages */}
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register-org" element={<RegisterOrgPage />} />
            <Route path="/ticket/:publicToken" element={<PublicTicketPage />} />
            <Route path="/kiosk/:branchId" element={<KioskPage />} />
            <Route path="/display" element={<PublicDisplayPage />} />
            <Route path="/display/:branchIdentifier" element={<PublicDisplayPage />} />


            {/* Protected Dashboard Routes */}
            <Route path="/" element={<DashboardLayout />}>
              <Route index element={<DashboardPage />} />
              <Route path="superadmin" element={<SuperAdminPage />} />
              <Route path="superadmin/billing" element={<SuperadminBillingPage />} />
              <Route path="counter" element={<StaffCounterPage />} />
              <Route path="reception" element={<ReceptionPage />} />
              <Route path="branches" element={<BranchesPage />} />
              <Route path="billing" element={<BillingPortalPage />} />
              <Route path="reports" element={<ReportsPage />} />
            </Route>

            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </ToastProvider>
    </QueryClientProvider>
  );
};

