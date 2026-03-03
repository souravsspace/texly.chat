import { useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { CheckCircle } from "lucide-react";
import { useEffect } from "react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useAuth } from "@/providers/auth";

export const Route = createFileRoute("/dashboard/billing/")({
  component: BillingPage,
  validateSearch: (search: Record<string, unknown>) => {
    return {
      success: search.success === "true" || search.success === true,
      customer_session_token:
        typeof search.customer_session_token === "string"
          ? search.customer_session_token
          : undefined,
      type: typeof search.type === "string" ? search.type : undefined,
    };
  },
});

function BillingPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { success, customer_session_token, type } = Route.useSearch();
  const { initializeAuth } = useAuth();

  // Refresh user data when returning from successful checkout
  useEffect(() => {
    if (success) {
      // Wait a bit for webhook to process, then refresh user data and invalidate all caches
      const timer = setTimeout(() => {
        // Invalidate all cached queries to force refetch
        queryClient.invalidateQueries({ queryKey: ["billing-usage"] });
        queryClient.invalidateQueries({ queryKey: ["bots"] });
        queryClient.invalidateQueries();

        // Refresh auth store with fresh user data
        initializeAuth();
      }, 2000); // 2 second delay to allow webhook processing
      return () => clearTimeout(timer);
    }
  }, [success, initializeAuth, queryClient]);

  // Redirect to subscription page after 3 seconds if success
  useEffect(() => {
    if (success) {
      const timer = setTimeout(() => {
        navigate({ to: "/dashboard/billing/subscription" });
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [success, navigate]);

  if (success) {
    return (
      <div className="max-w-4xl space-y-6">
        <div>
          <h2 className="font-bold text-3xl tracking-tight">
            Payment Successful
          </h2>
          <p className="text-muted-foreground">
            Your payment has been processed successfully.
          </p>
        </div>

        <Alert className="border-green-500 bg-green-50 dark:bg-green-950">
          <CheckCircle className="h-4 w-4 text-green-600 dark:text-green-400" />
          <AlertTitle className="text-green-600 dark:text-green-400">
            Success!
          </AlertTitle>
          <AlertDescription className="text-green-700 dark:text-green-300">
            {type === "usage"
              ? "Your usage charge has been processed successfully."
              : "Your subscription has been activated successfully."}
            {customer_session_token && (
              <div className="mt-2 text-xs">
                Session: {customer_session_token.substring(0, 20)}...
              </div>
            )}
          </AlertDescription>
        </Alert>

        <Card>
          <CardHeader>
            <CardTitle>What's Next?</CardTitle>
            <CardDescription>
              You can now access all Pro features.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-muted-foreground text-sm">
              Redirecting you to the subscription management page in 3
              seconds...
            </p>

            <div className="flex gap-4">
              <Button
                onClick={() =>
                  navigate({ to: "/dashboard/billing/subscription" })
                }
              >
                Manage Subscription
              </Button>
              <Button
                onClick={() => navigate({ to: "/dashboard" })}
                variant="outline"
              >
                Go to Dashboard
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Default billing overview page
  return (
    <div className="max-w-4xl space-y-6">
      <div>
        <h2 className="font-bold text-3xl tracking-tight">Billing</h2>
        <p className="text-muted-foreground">
          Manage your subscription and billing information.
        </p>
      </div>

      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Subscription</CardTitle>
            <CardDescription>
              View and manage your subscription plan.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button
              onClick={() =>
                navigate({ to: "/dashboard/billing/subscription" })
              }
            >
              View Subscription
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Usage</CardTitle>
            <CardDescription>
              Track your usage and associated costs.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button
              onClick={() => navigate({ to: "/dashboard/billing/usage" })}
            >
              View Usage
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
