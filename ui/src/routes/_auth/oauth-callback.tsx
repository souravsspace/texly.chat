import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect, useRef } from "react";
import { z } from "zod";
import { Icons } from "@/components/icons";
import { useAuthStore } from "@/stores/auth";

const oauthCallbackSearchSchema = z.object({
  error: z.string().optional(),
});

export const Route = createFileRoute("/_auth/oauth-callback")({
  validateSearch: oauthCallbackSearchSchema,
  component: OAuthCallback,
});

function OAuthCallback() {
  const navigate = useNavigate();
  const { error } = Route.useSearch();
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
            useAuthStore.getState().login(token, user);
            navigate({ to: "/dashboard" });
          } else {
            console.error("Failed to fetch user info");
            navigate({ to: "/login" });
          }
        } catch (e) {
          console.error("Error processing callback", e);
          navigate({ to: "/login" });
        }
      } else {
        navigate({ to: "/login" });
      }
    };

    handleCallback();
  }, [navigate, error]);

  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="flex flex-col items-center gap-4">
        <Icons.logo className="h-12 w-12 animate-pulse" />
        <p className="text-muted-foreground">Authenticating...</p>
      </div>
    </div>
  );
}
