import { create } from 'zustand';
import { User, Organization } from '../types';

interface AuthState {
  user: User | null;
  organization: Organization | null;
  token: string | null;
  setAuth: (user: User, token: string, organization?: Organization) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: JSON.parse(localStorage.getItem('user') || 'null'),
  organization: JSON.parse(localStorage.getItem('organization') || 'null'),
  token: localStorage.getItem('access_token'),

  setAuth: (user, token, organization) => {
    localStorage.setItem('access_token', token);
    localStorage.setItem('user', JSON.stringify(user));
    if (organization) {
      localStorage.setItem('organization', JSON.stringify(organization));
    }
    set({ user, token, organization: organization || null });
  },

  logout: () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('user');
    localStorage.removeItem('organization');
    set({ user: null, token: null, organization: null });
  },
}));
