import { createFileRoute, Link } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/providers/auth";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  const { user, logout } = useAuth();

  return (
    <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
      <header className="mb-12 text-center">
        <h1 className="mb-4 font-bold text-5xl text-foreground">
          Welcome to My SaaS
        </h1>
        <p className="text-muted-foreground text-xl">
          A modern full-stack application built with Go and TanStack Start
        </p>
      </header>

      {user ? (
        <div className="mb-16 flex justify-center gap-4">
          <Button className="h-auto px-8 py-3 text-lg">
            <Link to="/dashboard">Go to Dashboard</Link>
          </Button>
          <Button
            className="h-auto px-8 py-3 text-lg"
            onClick={() => logout()}
            variant="secondary"
          >
            Logout
          </Button>
        </div>
      ) : (
        <div className="mb-16 flex justify-center gap-4">
          <Button className="h-auto px-8 py-3 text-lg">
            <Link to="/login">Login</Link>
          </Button>
          <Button className="h-auto px-8 py-3 text-lg" variant="secondary">
            <Link to="/signup">Sign Up</Link>
          </Button>
        </div>
      )}

      <section className="mt-16">
        <h2 className="mb-8 text-center font-bold text-3xl">Features</h2>
        <div className="grid grid-cols-1 gap-8 md:grid-cols-3">
          <div className="rounded-lg bg-secondary p-8 text-center">
            <h3 className="mb-2 font-semibold text-2xl">🚀 Fast</h3>
            <p className="text-muted-foreground">
              Built with performance in mind
            </p>
          </div>
          <div className="rounded-lg bg-secondary p-8 text-center">
            <h3 className="mb-2 font-semibold text-2xl">🔒 Secure</h3>
            <p className="text-muted-foreground">
              JWT authentication & bcrypt hashing
            </p>
          </div>
          <div className="rounded-lg bg-secondary p-8 text-center">
            <h3 className="mb-2 font-semibold text-2xl">📦 Simple</h3>
            <p className="text-muted-foreground">Single binary deployment</p>
          </div>
        </div>
      </section>
    </div>
  );
}
