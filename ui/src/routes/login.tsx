import { useForm } from "@tanstack/react-form";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { zodValidator } from "@tanstack/zod-form-adapter";
import * as React from "react";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Field,
  FieldContent,
  FieldError,
  FieldLabel,
} from "@/components/ui/fields";
import { Input } from "@/components/ui/input";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth";

export const Route = createFileRoute("/login")({
  component: Login,
});

const loginSchema = z.object({
  email: z.string().email("Invalid email address"),
  password: z.string().min(1, "Password is required"),
});

function Login() {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [globalError, setGlobalError] = React.useState("");

  const form = useForm({
    defaultValues: {
      email: "",
      password: "",
    },
    validatorAdapter: zodValidator(),
    validators: {
      onChange: loginSchema,
    },
    onSubmit: async ({ value }) => {
      setGlobalError("");
      try {
        const response = await api.auth.login(value.email, value.password);
        login(response.token, response.user);
        navigate({ to: "/dashboard" });
      } catch (err) {
        setGlobalError(err instanceof Error ? err.message : "Login failed");
      }
    },
  });

  return (
    <div className="flex min-h-screen items-center justify-center bg-muted px-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="mb-2 text-center font-bold text-3xl text-foreground">
            Login
          </CardTitle>
        </CardHeader>
        <CardContent>
          {globalError && (
            <div className="mb-4 rounded-lg border border-destructive bg-destructive/15 px-4 py-3 text-destructive text-sm">
              {globalError}
            </div>
          )}

          <form
            className="space-y-4"
            onSubmit={(e) => {
              e.preventDefault();
              e.stopPropagation();
              form.handleSubmit();
            }}
          >
            <form.Field
              children={(field) => (
                <Field>
                  <FieldLabel htmlFor={field.name}>Email</FieldLabel>
                  <FieldContent>
                    <Input
                      id={field.name}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      placeholder="you@example.com"
                      type="email"
                      value={field.state.value}
                    />
                  </FieldContent>
                  <FieldError
                    errors={field.state.meta.errors.map((err) => ({
                      message: err?.message || String(err),
                    }))}
                  />
                </Field>
              )}
              name="email"
            />

            <form.Field
              children={(field) => (
                <Field>
                  <FieldLabel htmlFor={field.name}>Password</FieldLabel>
                  <FieldContent>
                    <Input
                      id={field.name}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      placeholder="••••••••"
                      type="password"
                      value={field.state.value}
                    />
                  </FieldContent>
                  <FieldError
                    errors={field.state.meta.errors.map((err) => ({
                      message: err?.message || String(err),
                    }))}
                  />
                </Field>
              )}
              name="password"
            />

            <form.Subscribe
              children={([canSubmit, isSubmitting]) => (
                <Button className="w-full" disabled={!canSubmit} type="submit">
                  {isSubmitting ? "Logging in..." : "Login"}
                </Button>
              )}
              selector={(state) => [state.canSubmit, state.isSubmitting]}
            />
          </form>
        </CardContent>
        <CardFooter className="justify-center">
          <p className="text-center text-muted-foreground text-sm">
            Don't have an account?{" "}
            <Link
              className="font-semibold text-primary hover:text-primary/80"
              to="/signup"
            >
              Sign up
            </Link>
          </p>
        </CardFooter>
      </Card>
    </div>
  );
}
