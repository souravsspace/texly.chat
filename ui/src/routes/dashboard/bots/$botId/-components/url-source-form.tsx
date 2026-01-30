import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Link as LinkIcon, Loader2 } from "lucide-react";
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

export function UrlSourceForm({
  botId,
  onSuccess,
  onCancel,
}: {
  botId: string;
  onSuccess: () => void;
  onCancel: () => void;
}) {
  const [url, setUrl] = useState("");
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (url: string) => api.sources.create(botId, url),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sources", botId] });
      toast.success("URL source added successfully");
      setUrl("");
      onSuccess();
    },
    onError: (error: Error) => {
      toast.error(`Failed to add source: ${error.message}`);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!url.trim()) return;

    // Basic URL validation
    try {
      new URL(url);
      mutation.mutate(url);
    } catch {
      toast.error("Please enter a valid URL");
    }
  };

  return (
    <Card className="border-2">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <div className="rounded-md bg-primary/10 p-2">
            <LinkIcon className="h-5 w-5 text-primary" />
          </div>
          Add URL Source
        </CardTitle>
        <CardDescription>
          Scrape content from a web page to use as knowledge for your chatbot
        </CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit}>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="url">Website URL</Label>
            <div className="relative">
              <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                <LinkIcon className="h-4 w-4 text-muted-foreground" />
              </div>
              <Input
                className="pl-9"
                disabled={mutation.isPending}
                onChange={(e) => setUrl(e.target.value)}
                placeholder="https://example.com/docs"
                required
                type="url"
                value={url}
              />
            </div>
          </div>

          {mutation.isError && (
            <Alert variant="destructive">
              <AlertDescription>
                {mutation.error?.message || "Failed to add source"}
              </AlertDescription>
            </Alert>
          )}

          <Alert>
            <AlertDescription className="text-xs">
              The content will be scraped, chunked, and embedded for semantic
              search. This may take a few moments depending on the page size.
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
          <Button disabled={mutation.isPending || !url.trim()} type="submit">
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
