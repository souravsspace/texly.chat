import { createFileRoute } from "@tanstack/react-router";
import { ExternalLink } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export const Route = createFileRoute("/dashboard/billing/subscription")({
  component: SubscriptionPage,
});

import { useMutation } from "@tanstack/react-query";
import { api } from "@/api";

function SubscriptionPage() {
  const mutation = useMutation({
    mutationFn: () => api.billing.portal(),
    onSuccess: (data) => {
      if (data.url) {
        window.location.href = data.url;
      }
    },
    onError: (error) => {
      console.error("Failed to get portal URL:", error);
    },
  });

  const handleManageSubscription = () => {
    mutation.mutate();
  };

  return (
    <div className="max-w-4xl space-y-6">
      <div>
        <h2 className="font-bold text-3xl tracking-tight">Subscription</h2>
        <p className="text-muted-foreground">
          Manage your subscription and payment methods.
        </p>
      </div>

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
    </div>
  );
}
