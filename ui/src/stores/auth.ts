import Cookies from "js-cookie";
import { create } from "zustand";
import { api } from "@/api";
import type { User } from "@/api/index.types";

interface AuthState {
  user: User | null;
  token: string | null;
  loading: boolean;
  login: (token: string, user: User) => void;
  logout: () => void;
  setUser: (user: User) => void;
  initializeAuth: () => Promise<void>;
}

const TOKEN_COOKIE_NAME = "auth_token";
const COOKIE_OPTIONS = {
  expires: 7, // 7 days
  sameSite: "strict" as const,
  secure: process.env.NODE_ENV === "production",
};

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  loading: true,

  login: (token: string, user: User) => {
    Cookies.set(TOKEN_COOKIE_NAME, token, COOKIE_OPTIONS);
    set({ token, user });
  },

  logout: () => {
    Cookies.remove(TOKEN_COOKIE_NAME);
    set({ token: null, user: null });
  },

  setUser: (user: User) => {
    set({ user });
  },

  initializeAuth: async () => {
    const savedToken = Cookies.get(TOKEN_COOKIE_NAME);

    if (savedToken) {
      set({ token: savedToken });

      try {
        const user = await api.users.getMe();
        set({ user, loading: false });
      } catch (_error) {
        Cookies.remove(TOKEN_COOKIE_NAME);
        set({ token: null, user: null, loading: false });
      }
    } else {
      set({ loading: false });
    }
  },
}));
