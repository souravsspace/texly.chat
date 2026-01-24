import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import * as React from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { api, type Post } from "@/lib/api";
import { useAuth } from "@/lib/auth";

export const Route = createFileRoute("/dashboard")({
  component: Dashboard,
});

function Dashboard() {
  const navigate = useNavigate();
  const { user, token, logout, loading: authLoading } = useAuth();
  const [posts, setPosts] = React.useState<Post[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState("");
  const [showForm, setShowForm] = React.useState(false);
  const [title, setTitle] = React.useState("");
  const [content, setContent] = React.useState("");
  const [editingPost, setEditingPost] = React.useState<Post | null>(null);

  React.useEffect(() => {
    if (!(authLoading || token)) {
      navigate({ to: "/login" });
    }
  }, [authLoading, token, navigate]);

  const loadPosts = React.useCallback(async () => {
    setLoading(true);
    try {
      const data = await api.posts.list();
      setPosts(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load posts");
    } finally {
      setLoading(false);
    }
  }, []);

  React.useEffect(() => {
    if (user) {
      loadPosts();
    }
  }, [user, loadPosts]);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    try {
      if (editingPost) {
        await api.posts.update(editingPost.id, title, content);
      } else {
        await api.posts.create(title, content);
      }

      setTitle("");
      setContent("");
      setShowForm(false);
      setEditingPost(null);
      await loadPosts();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save post");
    }
  }

  async function handleDelete(id: string) {
    if (!confirm("Are you sure you want to delete this post?")) return;

    try {
      await api.posts.delete(id);
      await loadPosts();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete post");
    }
  }

  function startEdit(post: Post) {
    setEditingPost(post);
    setTitle(post.title);
    setContent(post.content);
    setShowForm(true);
  }

  function cancelEdit() {
    setEditingPost(null);
    setTitle("");
    setContent("");
    setShowForm(false);
  }

  if (authLoading || (!user && loading)) {
    return (
      <div className="flex h-screen items-center justify-center">
        Loading...
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
      <header className="mb-8 flex items-center justify-between border-border border-b-2 pb-6">
        <div>
          <h1 className="mb-2 font-bold text-4xl text-foreground">Dashboard</h1>
          {user && (
            <p className="text-muted-foreground">Welcome, {user.name}!</p>
          )}
        </div>
        <div className="flex gap-3">
          <Button asChild variant="secondary">
            <Link to="/">Home</Link>
          </Button>
          <Button onClick={() => logout()} variant="secondary">
            Logout
          </Button>
        </div>
      </header>

      {error && (
        <div className="mb-6 rounded-lg border border-destructive bg-destructive/15 px-4 py-3 text-destructive text-sm">
          {error}
        </div>
      )}

      <div className="mb-6 flex items-center justify-between">
        <h2 className="font-bold text-2xl text-foreground">Your Posts</h2>
        <Button
          onClick={() => {
            if (showForm) cancelEdit();
            else setShowForm(true);
          }}
        >
          {showForm ? "Cancel" : "New Post"}
        </Button>
      </div>

      {showForm && (
        <Card className="mb-8">
          <CardContent className="pt-6">
            <form className="space-y-4" onSubmit={handleSubmit}>
              <Input
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Post title"
                required
                type="text"
                value={title}
              />
              <Textarea
                onChange={(e) => setContent(e.target.value)}
                placeholder="Post content"
                required
                rows={4}
                value={content}
              />
              <div className="flex gap-3">
                <Button type="submit">
                  {editingPost ? "Update" : "Create"} Post
                </Button>
                {editingPost && (
                  <Button
                    onClick={cancelEdit}
                    type="button"
                    variant="secondary"
                  >
                    Cancel
                  </Button>
                )}
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {loading ? (
        <p className="py-12 text-center text-muted-foreground">
          Loading posts...
        </p>
      ) : posts.length === 0 ? (
        <p className="py-12 text-center text-muted-foreground">
          No posts yet. Create your first post!
        </p>
      ) : (
        <div className="space-y-4">
          {posts.map((post) => (
            <Card className="transition-shadow hover:shadow-md" key={post.id}>
              <CardHeader>
                <CardTitle>{post.title}</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="whitespace-pre-wrap text-muted-foreground">
                  {post.content}
                </p>
              </CardContent>
              <CardFooter className="flex justify-between border-t pt-4">
                <small className="text-muted-foreground">
                  {new Date(post.created_at).toLocaleDateString()}
                </small>
                {post.user_id === user?.id && (
                  <div className="flex gap-4">
                    <Button
                      className="h-auto p-0 text-primary hover:bg-transparent hover:text-primary/80"
                      onClick={() => startEdit(post)}
                      variant="ghost"
                    >
                      Edit
                    </Button>
                    <Button
                      className="h-auto p-0 text-destructive hover:bg-transparent hover:text-destructive/80"
                      onClick={() => handleDelete(post.id)}
                      variant="ghost"
                    >
                      Delete
                    </Button>
                  </div>
                )}
              </CardFooter>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
