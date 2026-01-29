import Cookies from "js-cookie";
import type {
  AuthResponse,
  Bot,
  ChatTokenResponse,
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

  chat = {
    async *stream(
      botId: string,
      message: string
    ): AsyncGenerator<ChatTokenResponse, void, unknown> {
      const token = Cookies.get("auth_token");
      const headers: Record<string, string> = {
        "Content-Type": "application/json",
      };

      if (token) {
        headers.Authorization = `Bearer ${token}`;
      }

      const response = await fetch(`/api/bots/${botId}/chat`, {
        method: "POST",
        headers,
        body: JSON.stringify({ message }),
      });

      if (!response.ok) {
        const error = await response
          .json()
          .catch(() => ({ message: "Request failed" }));
        throw new Error(error.message || `HTTP ${response.status}`);
      }

      const reader = response.body?.getReader();
      if (!reader) {
        throw new Error("No response body");
      }

      const decoder = new TextDecoder();
      let buffer = "";

      try {
        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          buffer += decoder.decode(value, { stream: true });
          const lines = buffer.split("\n");
          buffer = lines.pop() || "";

          for (const line of lines) {
            const event = parseSSELine(line);
            if (event) {
              if (event === "DONE") return;
              yield event;
            }
          }
        }
      } finally {
        reader.releaseLock();
      }
    },
  };
}

function parseSSELine(line: string): ChatTokenResponse | "DONE" | null {
  if (!line.startsWith("data: ")) return null;

  const data = line.slice(6);
  if (data === "[DONE]") return "DONE";

  try {
    return JSON.parse(data) as ChatTokenResponse;
  } catch (e) {
    console.error("Failed to parse SSE data:", data, e);
    return null;
  }
}

export const api = new ApiClient();
