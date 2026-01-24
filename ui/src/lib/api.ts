import type {
  AuthResponse,
  CreatePostRequest,
  Post,
  UpdatePostRequest,
  User,
} from "#api-types";

const API_BASE = "/api";

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token =
    typeof window !== "undefined" ? localStorage.getItem("token") : null;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE}${endpoint}`, {
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

export const api = {
  auth: {
    signup: (email: string, password: string, name: string) =>
      request<AuthResponse>("/auth/signup", {
        method: "POST",
        body: JSON.stringify({ email, password, name }),
      }),

    login: (email: string, password: string) =>
      request<AuthResponse>("/auth/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      }),
  },

  users: {
    getMe: () => request<User>("/users/me"),
  },

  posts: {
    list: () => request<Post[]>("/posts"),

    get: (id: string) => request<Post>(`/posts/${id}`),

    create: (title: string, content: string) => {
      const payload: CreatePostRequest = { title, content };
      return request<Post>("/posts", {
        method: "POST",
        body: JSON.stringify(payload),
      });
    },

    update: (id: string, title: string, content: string) => {
      const payload: UpdatePostRequest = { title, content };
      return request<Post>(`/posts/${id}`, {
        method: "PUT",
        body: JSON.stringify(payload),
      });
    },

    delete: (id: string) =>
      request<void>(`/posts/${id}`, {
        method: "DELETE",
      }),
  },
};
