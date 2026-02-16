import {
  queryOptions,
  useMutation,
  useSuspenseQuery,
} from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { format } from "date-fns";
import { CreditCard, Loader2 } from "lucide-react";
import { api } from "@/api";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { cn } from "@/lib/utils";

const usageQueryOptions = queryOptions({
  queryKey: ["billing-usage"],
  queryFn: () => api.billing.usage(),
});

export const Route = createFileRoute("/dashboard/billing/usage")({
  component: BillingUsagePage,
  loader: ({ context }) => {
    return context.queryClient.ensureQueryData(usageQueryOptions);
  },
});

function BillingUsagePage() {
  const { data } = useSuspenseQuery(usageQueryOptions);

  const portalMutation = useMutation({
    mutationFn: () => api.billing.portal(),
    onSuccess: (data) => {
      if (data.url) window.location.href = data.url;
    },
    onError: (err) => console.error("Portal error:", err),
  });

  const handlePortal = () => {
    portalMutation.mutate();
  };

  if (!data)
    return (
      <div className="p-12 text-center">
        <Loader2 className="inline-block animate-spin" />
      </div>
    );

  return (
    <div className="max-w-4xl space-y-6">
      <div>
        <h2 className="font-bold text-3xl tracking-tight">Billing & Usage</h2>
        <p className="text-muted-foreground">
          Manage your plan and track your usage.
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="font-medium text-sm">Current Plan</CardTitle>
            <CreditCard className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="font-bold text-2xl capitalize">{data.tier}</div>
            <p className="text-muted-foreground text-xs">
              {data.billing_cycle_end
                ? `Renews ${format(new Date(data.billing_cycle_end), "MMM d, yyyy")}`
                : "No active subscription"}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="font-medium text-sm">
              Credits Balance
            </CardTitle>
            <div className="font-bold text-2xl">
              ${data.credits_balance.toFixed(2)}
            </div>
          </CardHeader>
          <CardContent>
            <div className="mb-2 text-muted-foreground text-xs">
              of ${data.credits_allocated.toFixed(2)} monthly
            </div>
            <Progress
              className="h-2"
              value={
                (data.credits_balance / Math.max(data.credits_allocated, 1)) *
                100
              }
            />
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="font-medium text-sm">Current Usage</CardTitle>
            <div className="font-bold text-2xl">
              ${data.current_period_usage.toFixed(2)}
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground text-xs">
              Total value consumed this cycle
            </p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Subscription Management</CardTitle>
          <CardDescription>
            Update your payment method, download invoices, or change your plan.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-x-4">
          {data.tier === "free" ? (
            <Link
              className={cn(buttonVariants(), "no-underline")}
              to="/pricing"
            >
              Upgrade to Pro
            </Link>
          ) : (
            <Button
              disabled={portalMutation.isPending}
              onClick={handlePortal}
              variant="outline"
            >
              {portalMutation.isPending ? "Loading..." : "Manage Subscription"}
            </Button>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
