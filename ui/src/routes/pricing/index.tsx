import { useMutation } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { Check, X } from "lucide-react";
import { api } from "@/api";
import { Badge } from "@/components/ui/badge";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { useAuth } from "@/providers/auth";

export const Route = createFileRoute("/pricing/")({
  component: PricingPage,
});

function PricingPage() {
  const { user } = useAuth();
  const isAuthenticated = !!user;

  const mutation = useMutation({
    mutationFn: () => api.billing.checkout(),
    onSuccess: (data) => {
      if (data.url) {
        window.location.href = data.url;
      }
    },
    onError: (error) => {
      console.error("Error creating checkout session:", error);
    },
  });

  const handleUpgrade = () => {
    if (!isAuthenticated) return;
    mutation.mutate();
  };

  const renderProButton = () => {
    if (!isAuthenticated) {
      return (
        <Link
          className={cn(buttonVariants({ variant: "default" }), "w-full")}
          to="/signup"
        >
          Start Free Trial
        </Link>
      );
    }
    if (user?.tier === "pro") {
      return (
        <Button className="w-full" disabled variant="outline">
          Current Plan
        </Button>
      );
    }
    return (
      <Button
        className="w-full"
        disabled={mutation.isPending}
        onClick={handleUpgrade}
      >
        {mutation.isPending ? "Loading..." : "Upgrade to Pro"}
      </Button>
    );
  };

  return (
    <div className="container mx-auto py-20">
      <div className="mb-16 text-center">
        <h1 className="mb-4 font-extrabold text-4xl tracking-tight lg:text-5xl">
          Simple, Transparent Pricing
        </h1>
        <p className="text-muted-foreground text-xl">
          Choose the plan that's right for you. Change or cancel anytime.
        </p>
      </div>

      <div className="mx-auto grid max-w-6xl grid-cols-1 gap-8 md:grid-cols-3">
        {/* Free Tier */}
        <Card className="flex flex-col">
          <CardHeader>
            <CardTitle className="text-2xl">Free</CardTitle>
            <CardDescription>Perfect for trying out Texly.</CardDescription>
          </CardHeader>
          <CardContent className="flex-1">
            <div className="mb-6 font-bold text-3xl">
              $0
              <span className="font-normal text-lg text-muted-foreground">
                /mo
              </span>
            </div>
            <ul className="space-y-3 text-sm">
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> 1 Bot
              </li>
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> 100
                Messages/mo
              </li>
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> 5 Sources per
                Bot
              </li>
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> 10MB Storage
              </li>
              <li className="flex items-center text-muted-foreground">
                <X className="mr-2 h-4 w-4" /> No Credits Included
              </li>
            </ul>
          </CardContent>
          <CardFooter>
            <Button
              className="w-full"
              disabled={isAuthenticated && user?.tier === "free"}
              variant="outline"
            >
              {isAuthenticated && user?.tier === "free"
                ? "Current Plan"
                : "Get Started"}
            </Button>
          </CardFooter>
        </Card>

        {/* Pro Tier */}
        <Card className="relative flex flex-col overflow-hidden border-primary">
          <div className="absolute top-0 right-0 p-3">
            <Badge className="bg-primary text-primary-foreground">
              Popular
            </Badge>
          </div>
          <CardHeader>
            <CardTitle className="text-2xl">Pro</CardTitle>
            <CardDescription>For power users and creators.</CardDescription>
          </CardHeader>
          <CardContent className="flex-1">
            <div className="mb-6 font-bold text-3xl">
              $20
              <span className="font-normal text-lg text-muted-foreground">
                /mo
              </span>
            </div>
            <p className="mb-4 text-muted-foreground text-sm">
              Includes <strong>$20</strong> in monthly usage credits!
            </p>
            <ul className="space-y-3 text-sm">
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> 5 Included
                Bots
              </li>
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> Pay-as-you-go
                Messages
              </li>
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> 50 Sources per
                Bot
              </li>
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> 1GB Storage
              </li>
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> Priority
                Support
              </li>
            </ul>
          </CardContent>
          <CardFooter>{renderProButton()}</CardFooter>
        </Card>

        {/* Enterprise Tier */}
        <Card className="flex flex-col">
          <CardHeader>
            <CardTitle className="text-2xl">Enterprise</CardTitle>
            <CardDescription>Custom solutions for large teams.</CardDescription>
          </CardHeader>
          <CardContent className="flex-1">
            <div className="mb-6 font-bold text-3xl">Custom</div>
            <ul className="space-y-3 text-sm">
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> Unlimited Bots
              </li>
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> Custom Usage
                Rates
              </li>
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> Dedicated
                Support
              </li>
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> SSO & Audit
                Logs
              </li>
              <li className="flex items-center">
                <Check className="mr-2 h-4 w-4 text-green-500" /> SLA Guarantee
              </li>
            </ul>
          </CardContent>
          <CardFooter>
            <a
              className={cn(buttonVariants({ variant: "outline" }), "w-full")}
              href="mailto:sales@texly.chat"
            >
              Contact Sales
            </a>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
