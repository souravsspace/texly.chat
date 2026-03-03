import { useMutation, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { ExternalLink } from "lucide-react";
import { api } from "@/api";
import { billingUsageQueryOptions } from "@/api/queries";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/dashboard/billing/subscription")({
  component: SubscriptionPage,
  loader: ({ context }) => {
    return context.queryClient.ensureQueryData(billingUsageQueryOptions);
  },
});

function SubscriptionPage() {
  const { data: billingData } = useSuspenseQuery(billingUsageQueryOptions);

  const mutation = useMutation({
    mutationFn: () => api.billing.portal(),
  });

  const handleManageSubscription = async () => {
    try {
      const data = await mutation.mutateAsync();
      if (data.url) {
        window.open(data.url, "_blank", "noopener,noreferrer");
      }
    } catch {
      // Error is already handled by mutation state
    }
  };

  const isFreeTier = billingData.tier === "free";

  return (
    <div className="max-w-4xl space-y-6">
      <div>
        <h2 className="font-bold text-3xl tracking-tight">Subscription</h2>
        <p className="text-muted-foreground">
          Manage your subscription and payment methods.
        </p>
      </div>

      {isFreeTier ? (
        <Card>
          <CardHeader>
            <CardTitle>Free Plan</CardTitle>
            <CardDescription>
              You are currently on the free tier.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-muted-foreground text-sm">
              Upgrade to Pro to unlock additional features:
            </p>
            <ul className="ml-2 list-inside list-disc space-y-1 text-muted-foreground text-sm">
              <li>Increased usage limits</li>
              <li>Priority support</li>
              <li>Advanced analytics</li>
              <li>Custom branding</li>
            </ul>

            <div className="pt-4">
              <Link
                className={cn(buttonVariants(), "no-underline")}
                to="/pricing"
              >
                Upgrade to Pro
              </Link>
            </div>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <CardHeader>
            <CardTitle>Customer Portal</CardTitle>
            <CardDescription>
              Access our secure customer portal to manage your subscription
              details.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-muted-foreground text-sm">
              In the customer portal, you can:
            </p>
            <ul className="ml-2 list-inside list-disc space-y-1 text-muted-foreground text-sm">
              <li>Update your payment method</li>
              <li>Download past invoices</li>
              <li>Upgrade or downgrade your plan</li>
              <li>Cancel your subscription</li>
            </ul>

            <div className="pt-4">
              <Button
                disabled={mutation.isPending}
                onClick={handleManageSubscription}
              >
                {mutation.isPending ? "Loading..." : "Manage Subscription"}{" "}
                <ExternalLink className="ml-2 h-4 w-4" />
              </Button>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
