import { useMutation } from "@tanstack/react-query";
import { Check, Loader2, Sparkles } from "lucide-react";
import { api } from "@/api";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

interface UpgradeModalProps {
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  trigger?: React.ReactNode;
}

export function UpgradeModal({
  open,
  onOpenChange,
  trigger,
}: UpgradeModalProps) {
  const { mutate: handleUpgrade, isPending } = useMutation({
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

  const features = [
    "Unlimited Messages (Pay-as-you-go)",
    "5 Included Bots",
    "Increased Limits & Storage",
    "$20 in Monthly Credits Included",
  ];

  return (
    <Dialog onOpenChange={onOpenChange} open={open}>
      {trigger && <DialogTrigger>{trigger}</DialogTrigger>}
      <DialogContent className="overflow-hidden border-border bg-card p-0 text-card-foreground shadow-lg sm:max-w-[425px]">
        <div className="space-y-4 p-6 pt-8 text-center">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
            <Sparkles className="h-6 w-6 text-primary" />
          </div>
          <DialogHeader>
            <DialogTitle className="font-bold text-2xl tracking-tight">
              Unlock Pro Features
            </DialogTitle>
            <DialogDescription className="mx-auto max-w-xs text-muted-foreground text-sm">
              Upgrade to unleash the full power of your AI assistants.
            </DialogDescription>
          </DialogHeader>
        </div>

        <div className="border-border/50 border-y bg-muted/30 px-6 py-4">
          <ul className="space-y-3">
            {features.map((feature, index) => (
              <li
                className="flex items-center font-medium text-foreground/90 text-sm"
                // biome-ignore lint/suspicious/noArrayIndexKey: index is fine here
                key={index}
              >
                <Check className="mr-3 h-4 w-4 shrink-0 text-primary" />
                {feature}
              </li>
            ))}
          </ul>
        </div>

        <div className="space-y-4 p-6">
          <div className="flex items-baseline justify-center gap-1">
            <span className="font-bold text-3xl text-foreground">$20</span>
            <span className="font-medium text-muted-foreground text-sm">
              /month
            </span>
          </div>

          <DialogFooter className="sm:justify-center">
            <Button
              className="w-full font-semibold text-base shadow-sm transition-all hover:shadow"
              disabled={isPending}
              onClick={() => handleUpgrade()}
              size="lg"
              type="submit"
            >
              {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Upgrade to Pro
            </Button>
          </DialogFooter>
          <p className="text-center text-muted-foreground/60 text-xs">
            Cancel anytime. Secure checkout.
          </p>
        </div>
      </DialogContent>
    </Dialog>
  );
}
