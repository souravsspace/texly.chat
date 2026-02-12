import { useForm } from "@tanstack/react-form-start";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { z } from "zod";
import { api } from "#api";
import { Icons } from "@/components/icons";
import { Button, buttonVariants } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";
import { useAuth } from "@/providers/auth";

export const Route = createFileRoute("/_auth/signup")({
  component: Signup,
});

const signupSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters"),
  email: z.email("Invalid email address"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

function Signup() {
  const navigate = useNavigate();
  const { login } = useAuth();
  const apiBaseUrl = import.meta.env.DEV ? "http://localhost:8080" : "";

  const [globalError, setGlobalError] = useState("");

  const form = useForm({
    defaultValues: {
      name: "",
      email: "",
      password: "",
    },
    validators: {
      onSubmit: signupSchema,
    },
    onSubmit: async ({ value }) => {
      setGlobalError("");
      try {
        const response = await api.auth.signup(
          value.email,
          value.password,
          value.name
        );
        login(response.token, response.user);
        navigate({ to: "/dashboard" });
      } catch (err) {
        setGlobalError(err instanceof Error ? err.message : "Signup failed");
      }
    },
  });

  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="relative w-full max-w-md overflow-hidden border bg-gradient-to-b from-muted/50 to-card px-8 py-8 shadow-lg/5 dark:from-transparent dark:shadow-xl">
        <div
          className="absolute inset-0 -top-px -left-px z-0"
          style={{
            backgroundImage: `
        linear-gradient(to right, color-mix(in srgb, var(--card-foreground) 8%, transparent) 1px, transparent 1px),
        linear-gradient(to bottom, color-mix(in srgb, var(--card-foreground) 8%, transparent) 1px, transparent 1px)
      `,
            backgroundSize: "20px 20px",
            backgroundPosition: "0 0, 0 0",
            maskImage: `
        repeating-linear-gradient(
              to right,
              black 0px,
              black 3px,
              transparent 3px,
              transparent 8px
            ),
            repeating-linear-gradient(
              to bottom,
              black 0px,
              black 3px,
              transparent 3px,
              transparent 8px
            ),
            radial-gradient(ellipse 70% 50% at 50% 0%, #000 60%, transparent 100%)
      `,
            WebkitMaskImage: `
 repeating-linear-gradient(
              to right,
              black 0px,
              black 3px,
              transparent 3px,
              transparent 8px
            ),
            repeating-linear-gradient(
              to bottom,
              black 0px,
              black 3px,
              transparent 3px,
              transparent 8px
            ),
            radial-gradient(ellipse 70% 50% at 50% 0%, #000 60%, transparent 100%)
      `,
            maskComposite: "intersect",
            WebkitMaskComposite: "source-in",
          }}
        />

        <div className="relative isolate flex flex-col items-center">
          <Icons.logo className="h-9 w-9" />
          <p className="mt-4 font-semibold text-xl tracking-tight">
            Sign up for Texly AI
          </p>

          <a
            className={cn(
              buttonVariants({ variant: "outline", size: "lg" }),
              "mt-8 w-full gap-3"
            )}
            href={`${apiBaseUrl}/api/auth/google`}
          >
            <Icons.google />
            Continue with Google
          </a>

          <div className="my-7 flex w-full items-center justify-center overflow-hidden">
            <Separator className="flex-1" />
            <span className="px-2 text-muted-foreground text-sm">OR</span>
            <Separator className="flex-1" />
          </div>

          <form
            className="w-full space-y-4"
            onSubmit={(e) => {
              e.preventDefault();
              e.stopPropagation();
              form.handleSubmit();
            }}
          >
            {globalError && (
              <div className="mb-4 rounded-lg border border-destructive bg-destructive/15 px-4 py-3 text-destructive text-sm">
                {globalError}
              </div>
            )}

            <form.Field name="name">
              {(field) => (
                <div className="space-y-2">
                  <Label htmlFor={field.name}>Name</Label>
                  <Input
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="John Doe"
                    type="text"
                    value={field.state.value}
                  />
                  {field.state.meta.errors ? (
                    <p className="text-destructive text-sm">
                      {field.state.meta.errors
                        .map((err) => err?.message || String(err))
                        .join(", ")}
                    </p>
                  ) : null}
                </div>
              )}
            </form.Field>

            <form.Field name="email">
              {(field) => (
                <div className="space-y-2">
                  <Label htmlFor={field.name}>Email</Label>
                  <Input
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="you@example.com"
                    type="email"
                    value={field.state.value}
                  />
                  {field.state.meta.errors ? (
                    <p className="text-destructive text-sm">
                      {field.state.meta.errors
                        .map((err) => err?.message || String(err))
                        .join(", ")}
                    </p>
                  ) : null}
                </div>
              )}
            </form.Field>

            <form.Field name="password">
              {(field) => (
                <div className="space-y-2">
                  <Label htmlFor={field.name}>Password</Label>
                  <Input
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="••••••••"
                    type="password"
                    value={field.state.value}
                  />
                  {field.state.meta.errors ? (
                    <p className="text-destructive text-sm">
                      {field.state.meta.errors
                        .map((err) => err?.message || String(err))
                        .join(", ")}
                    </p>
                  ) : null}
                </div>
              )}
            </form.Field>

            <form.Subscribe
              selector={(state) => [state.canSubmit, state.isSubmitting]}
            >
              {([canSubmit, isSubmitting]) => (
                <Button
                  className="mt-4 w-full"
                  disabled={!canSubmit}
                  size="lg"
                  type="submit"
                >
                  {isSubmitting ? "Creating account..." : "Continue with Email"}
                </Button>
              )}
            </form.Subscribe>
          </form>

          <p className="mt-5 text-center text-sm">
            Already have an account?
            <Link className="ml-1 text-muted-foreground underline" to="/login">
              Log in
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
