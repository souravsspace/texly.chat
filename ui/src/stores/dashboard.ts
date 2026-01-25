import { create } from "zustand";
import { api } from "@/api";
import type { Post } from "@/api/index.types";

interface DashboardState {
  posts: Post[];
  loading: boolean;
  error: string;
  showForm: boolean;
  title: string;
  content: string;
  editingPost: Post | null;

  // Actions
  setShowForm: (show: boolean) => void;
  setTitle: (title: string) => void;
  setContent: (content: string) => void;
  startEdit: (post: Post) => void;
  cancelEdit: () => void;

  // Async Actions
  fetchPosts: () => Promise<void>;
  savePost: () => Promise<void>;
  deletePost: (id: string) => Promise<void>;
}

export const useDashboardStore = create<DashboardState>((set, get) => ({
  posts: [],
  loading: true,
  error: "",
  showForm: false,
  title: "",
  content: "",
  editingPost: null,

  setShowForm: (show) => set({ showForm: show }),
  setTitle: (title) => set({ title }),
  setContent: (content) => set({ content }),

  startEdit: (post) =>
    set({
      editingPost: post,
      title: post.title,
      content: post.content,
      showForm: true,
    }),

  cancelEdit: () =>
    set({
      editingPost: null,
      title: "",
      content: "",
      showForm: false,
    }),

  fetchPosts: async () => {
    set({ loading: true, error: "" });
    try {
      const posts = await api.posts.list();
      set({ posts });
    } catch (err) {
      set({
        error: err instanceof Error ? err.message : "Failed to load posts",
      });
    } finally {
      set({ loading: false });
    }
  },

  savePost: async () => {
    const { editingPost, title, content } = get();
    set({ error: "" });

    try {
      if (editingPost) {
        await api.posts.update(editingPost.id, title, content);
      } else {
        await api.posts.create(title, content);
      }

      // Reset form and reload
      get().cancelEdit();
      await get().fetchPosts();
    } catch (err) {
      set({
        error: err instanceof Error ? err.message : "Failed to save post",
      });
    }
  },

  deletePost: async (id) => {
    set({ error: "" });
    try {
      await api.posts.delete(id);
      await get().fetchPosts();
    } catch (err) {
      set({
        error: err instanceof Error ? err.message : "Failed to delete post",
      });
    }
  },
}));
