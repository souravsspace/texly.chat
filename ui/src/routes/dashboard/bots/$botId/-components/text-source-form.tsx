import { useMutation, useQueryClient } from "@tanstack/react-query";
import { FileText, Loader2 } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { api } from "@/api";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { MAX_TEXT_SIZE_BYTES } from "@/lib/constants";
import { cn } from "@/lib/utils";

export function TextSourceForm({
  botId,
  onSuccess,
  onCancel,
}: {
  botId: string;
  onSuccess: () => void;
  onCancel: () => void;
}) {
  const [text, setText] = useState("");
  const [name, setName] = useState("");
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: ({ text, name }: { text: string; name?: string }) =>
      api.sources.createText(botId, text, name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sources", botId] });
      toast.success("Text source added successfully");
      setText("");
      setName("");
      onSuccess();
    },
    onError: (error: Error) => {
      toast.error(`Failed to add text source: ${error.message}`);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!text.trim()) return;

    // Validate text size
    const textSize = new Blob([text]).size;
    if (textSize > MAX_TEXT_SIZE_BYTES) {
      toast.error(
        `Text size exceeds ${MAX_TEXT_SIZE_BYTES / 1024 / 1024}MB limit`
      );
      return;
    }

    mutation.mutate({ text, name: name.trim() || undefined });
  };

  const textSize = new Blob([text]).size;
  const textSizeMB = (textSize / 1024 / 1024).toFixed(2);
  const maxSizeMB = (MAX_TEXT_SIZE_BYTES / 1024 / 1024).toFixed(0);
  const isOverLimit = textSize > MAX_TEXT_SIZE_BYTES;
  const usagePercentage = (textSize / MAX_TEXT_SIZE_BYTES) * 100;

  return (
    <Card className="border-2">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <div className="rounded-md bg-primary/10 p-2">
            <FileText className="h-5 w-5 text-primary" />
          </div>
          Add Text Source
        </CardTitle>
        <CardDescription>
          Paste or type text content to use as knowledge for your chatbot
        </CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit}>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Name (optional)</Label>
            <Input
              disabled={mutation.isPending}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g., Product Documentation"
              type="text"
              value={name}
            />
          </div>

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label htmlFor="text">Text Content</Label>
              <span
                className={cn(
                  "font-medium text-xs",
                  isOverLimit && "text-destructive",
                  !isOverLimit && usagePercentage > 80 && "text-orange-500",
                  !isOverLimit &&
                    usagePercentage <= 80 &&
                    "text-muted-foreground"
                )}
              >
                {textSizeMB} / {maxSizeMB} MB
              </span>
            </div>
            <Textarea
              className="min-h-[300px] font-mono text-sm"
              disabled={mutation.isPending}
              onChange={(e) => setText(e.target.value)}
              placeholder="Paste your text content here..."
              required
              value={text}
            />
          </div>

          {mutation.isError && (
            <Alert variant="destructive">
              <AlertDescription>
                {mutation.error?.message || "Failed to add text source"}
              </AlertDescription>
            </Alert>
          )}

          {isOverLimit && (
            <Alert variant="destructive">
              <AlertDescription>
                Text size exceeds the maximum limit of {maxSizeMB}MB. Please
                reduce the content or split it into multiple sources.
              </AlertDescription>
            </Alert>
          )}

          <Alert>
            <AlertDescription className="text-xs">
              The text will be chunked and embedded for semantic search. This
              may take a few moments depending on the text length.
            </AlertDescription>
          </Alert>
        </CardContent>
        <CardFooter className="flex justify-between border-t pt-4">
          <Button
            disabled={mutation.isPending}
            onClick={onCancel}
            type="button"
            variant="outline"
          >
            Cancel
          </Button>
          <Button
            disabled={
              mutation.isPending ||
              !text.trim() ||
              textSize > MAX_TEXT_SIZE_BYTES
            }
            type="submit"
          >
            {mutation.isPending && (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            )}
            Add Source
          </Button>
        </CardFooter>
      </form>
    </Card>
  );
}
