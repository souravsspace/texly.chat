import Cookies from "js-cookie";
import type {
  AuthResponse,
  Bot,
  CreateBotRequest,
  CreateSourceRequest,
  Source,
  UpdateBotRequest,
  User,
} from "@/api/index.types";

class ApiClient {
  private readonly BASE_URL = "/api";

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const token = Cookies.get("auth_token");

    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...(options.headers as Record<string, string>),
    };

    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }

    const response = await fetch(`${this.BASE_URL}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: "Request failed" }));
      throw new Error(error.message || `HTTP ${response.status}`);
    }

    if (response.status === 204) {
      return {} as T;
    }

    return response.json();
  }

  auth = {
    signup: (email: string, password: string, name: string) =>
      this.request<AuthResponse>("/auth/signup", {
        method: "POST",
        body: JSON.stringify({ email, password, name }),
      }),

    login: (email: string, password: string) =>
      this.request<AuthResponse>("/auth/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      }),
  };

  users = {
    getMe: () => this.request<User>("/users/me"),
  };

  bots = {
    list: () => this.request<Bot[]>("/bots"),

    get: (id: string) => this.request<Bot>(`/bots/${id}`),

    create: (name: string, system_prompt?: string) => {
      const payload: CreateBotRequest = {
        name,
        system_prompt: system_prompt || "",
      };
      return this.request<Bot>("/bots", {
        method: "POST",
        body: JSON.stringify(payload),
      });
    },

    update: (id: string, name: string, system_prompt?: string) => {
      const payload: UpdateBotRequest = {
        name,
        system_prompt: system_prompt || "",
      };
      return this.request<Bot>(`/bots/${id}`, {
        method: "PUT",
        body: JSON.stringify(payload),
      });
    },

    delete: (id: string) =>
      this.request<void>(`/bots/${id}`, {
        method: "DELETE",
      }),
  };

  sources = {
    list: (botId: string) => this.request<Source[]>(`/bots/${botId}/sources`),

    get: (botId: string, sourceId: string) =>
      this.request<Source>(`/bots/${botId}/sources/${sourceId}`),

    create: (botId: string, url: string) => {
      const payload: CreateSourceRequest = { url };
      return this.request<Source>(`/bots/${botId}/sources`, {
        method: "POST",
        body: JSON.stringify(payload),
      });
    },

    delete: (botId: string, sourceId: string) =>
      this.request<void>(`/bots/${botId}/sources/${sourceId}`, {
        method: "DELETE",
      }),
  };
}

export const api = new ApiClient();
