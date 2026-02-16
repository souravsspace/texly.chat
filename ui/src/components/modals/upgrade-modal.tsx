import { Check } from "lucide-react";
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
  const handleUpgrade = async () => {
    try {
      const response = await fetch("/api/billing/checkout", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      const data = await response.json();
      if (data.url) {
        window.location.href = data.url;
      }
    } catch (error) {
      console.error("Error creating checkout session:", error);
    }
  };

  return (
    <Dialog onOpenChange={onOpenChange} open={open}>
      {trigger && <DialogTrigger>{trigger}</DialogTrigger>}
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Upgrade to Pro</DialogTitle>
          <DialogDescription>
            Unlock the full potential of Texly with our Pro plan.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="flex items-center space-x-2">
            <Check className="h-4 w-4 text-green-500" />
            <span>Unlimited Messages (Pay-as-you-go)</span>
          </div>
          <div className="flex items-center space-x-2">
            <Check className="h-4 w-4 text-green-500" />
            <span>5 Included Bots</span>
          </div>
          <div className="flex items-center space-x-2">
            <Check className="h-4 w-4 text-green-500" />
            <span>Increased Limits & Storage</span>
          </div>
          <div className="mt-4 text-center">
            <div className="font-bold text-3xl">
              $20
              <span className="font-normal text-muted-foreground text-sm">
                /mo
              </span>
            </div>
            <p className="text-muted-foreground text-xs">
              Includes $20 in monthly credits
            </p>
          </div>
        </div>
        <DialogFooter>
          <Button className="w-full" onClick={handleUpgrade} type="submit">
            Upgrade Now
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
