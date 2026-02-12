import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect, useRef } from "react";
import { z } from "zod";
import { Icons } from "@/components/icons";
import { useAuth } from "@/providers/auth";

const oauthCallbackSearchSchema = z.object({
  error: z.string().optional(),
});

export const Route = createFileRoute("/_auth/auth/callback")({
  validateSearch: oauthCallbackSearchSchema,
  component: OAuthCallback,
});

function OAuthCallback() {
  const navigate = useNavigate();
  const { error } = Route.useSearch();
  const { login } = useAuth();
  const processedRef = useRef(false);

  useEffect(() => {
    if (processedRef.current) return;

    const handleCallback = async () => {
      processedRef.current = true;
      if (error) {
        console.error("OAuth error:", error);
        navigate({ to: "/login" });
        return;
      }

      // Parse token from URL fragment
      const hash = window.location.hash.substring(1); // remove #
      const params = new URLSearchParams(hash);
      const token = params.get("token");

      if (token) {
        try {
          const response = await fetch("/api/users/me", {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });

          if (response.ok) {
            const user = await response.json();
            login(token, user);
            navigate({ to: "/dashboard" });
          } else {
            console.error("Failed to fetch user info", await response.text());
            navigate({ to: "/login" });
          }
        } catch (e) {
          console.error("Error processing callback", e);
          navigate({ to: "/login" });
        }
      } else {
        console.error("No token found in URL fragment");
        navigate({ to: "/login" });
      }
    };

    handleCallback();
  }, [navigate, error, login]);

  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="flex flex-col items-center gap-4">
        <Icons.logo className="h-12 w-12 animate-pulse" />
        <p className="text-muted-foreground">Authenticating...</p>
      </div>
    </div>
  );
}
