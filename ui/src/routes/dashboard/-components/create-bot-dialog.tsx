import { useForm } from "@tanstack/react-form-start";
import { useQueryClient } from "@tanstack/react-query";
import { AlertCircleIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { z } from "zod";
import { api } from "@/api";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
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
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Spinner } from "@/components/ui/spinner";
import { Textarea } from "@/components/ui/textarea";

interface CreateBotDialogProps {
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}

const botSchema = z.object({
  name: z
    .string()
    .min(1, "Name is required")
    .max(50, "Name must be less than 50 characters"),
  systemPrompt: z.string(),
});

export function CreateBotDialog({
  open: controlledOpen,
  onOpenChange,
}: CreateBotDialogProps) {
  const [internalOpen, setInternalOpen] = useState(false);
  const open = controlledOpen !== undefined ? controlledOpen : internalOpen;
  const queryClient = useQueryClient();
  const [submitError, setSubmitError] = useState<string | null>(null);

  function handleOpenChange(newOpen: boolean) {
    if (controlledOpen === undefined) {
      setInternalOpen(newOpen);
    }
    onOpenChange?.(newOpen);
  }

  const form = useForm({
    defaultValues: {
      name: "",
      systemPrompt: "",
    },
    validators: {
      onChange: botSchema,
    },
    onSubmit: async ({ value }) => {
      setSubmitError(null);
      try {
        await api.bots.create(value.name, value.systemPrompt || "");
        await queryClient.invalidateQueries({ queryKey: ["bots"] });
        toast.success("Bot created successfully");
        handleOpenChange(false);
      } catch (error) {
        console.error("Failed to create bot:", error);
        toast.error("Failed to create bot. Please try again.");
        if (error instanceof Error) {
          setSubmitError(error.message);
        } else {
          setSubmitError("An unexpected error occurred.");
        }
      }
    },
  });

  useEffect(() => {
    if (!open) {
      const timer = setTimeout(() => {
        form.reset();
        setSubmitError(null);
      }, 300);
      return () => clearTimeout(timer);
    }
  }, [open, form]);

  return (
    <Dialog onOpenChange={handleOpenChange} open={open}>
      {controlledOpen === undefined && (
        <DialogTrigger>
          <Button>Create Bot</Button>
        </DialogTrigger>
      )}
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create Chatbot</DialogTitle>
          <DialogDescription>
            Give your new chatbot a name and an optional personality.
          </DialogDescription>
        </DialogHeader>

        <form
          className="grid gap-4 py-4"
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
        >
          {submitError && (
            <Alert variant="destructive">
              <AlertCircleIcon className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>{submitError}</AlertDescription>
            </Alert>
          )}

          <form.Field name="name">
            {(field) => (
              <div className="grid grid-cols-4 items-center gap-4">
                <Label className="text-right" htmlFor={field.name}>
                  Name
                </Label>
                <div className="col-span-3">
                  <Input
                    className={
                      field.state.meta.errors.length
                        ? "border-destructive focus-visible:ring-destructive"
                        : ""
                    }
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    value={field.state.value}
                  />
                  {field.state.meta.errors.length > 0 && (
                    <p className="mt-1 text-destructive text-sm">
                      {field.state.meta.errors.join(", ")}
                    </p>
                  )}
                </div>
              </div>
            )}
          </form.Field>

          <form.Field name="systemPrompt">
            {(field) => (
              <div className="grid grid-cols-4 items-center gap-4">
                <Label className="text-right" htmlFor={field.name}>
                  Prompt
                </Label>
                <div className="col-span-3">
                  <Textarea
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="You are a helpful assistant..."
                    value={field.state.value}
                  />
                  {field.state.meta.errors.length > 0 && (
                    <p className="mt-1 text-destructive text-sm">
                      {field.state.meta.errors.join(", ")}
                    </p>
                  )}
                </div>
              </div>
            )}
          </form.Field>

          <DialogFooter>
            <form.Subscribe
              selector={(state) => [state.canSubmit, state.isSubmitting]}
            >
              {([canSubmit, isSubmitting]) => (
                <Button disabled={!canSubmit || isSubmitting} type="submit">
                  {isSubmitting && (
                    <Spinner className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  {isSubmitting ? "Creating..." : "Create Bot"}
                </Button>
              )}
            </form.Subscribe>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
