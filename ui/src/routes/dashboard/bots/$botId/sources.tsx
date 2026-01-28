import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { AlertCircle, ExternalLink, Loader2, Plus, Trash2 } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { api } from "@/api";
import type { Source } from "@/api/index.types";
import type { SourceStatus } from "@/api/index.types.manual";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
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

export const Route = createFileRoute("/dashboard/bots/$botId/sources")({
  component: SourcesPage,
});

function SourcesPage() {
  const { botId } = Route.useParams();
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false);
  const [url, setUrl] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const queryClient = useQueryClient();

  // Fetch sources with auto-refresh
  const {
    data: sources,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["sources", botId],
    queryFn: () => api.sources.list(botId),
    refetchInterval: 5000, // Poll every 5 seconds for status updates
  });

  // Create source mutation
  const createSourceMutation = useMutation({
    mutationFn: (url: string) => api.sources.create(botId, url),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sources", botId] });
      setIsAddDialogOpen(false);
      setUrl("");
      toast.success(
        "Source added successfully! Processing will begin shortly."
      );
    },
    onError: (error: Error) => {
      toast.error(`Failed to add source: ${error.message}`);
    },
    onSettled: () => {
      setIsSubmitting(false);
    },
  });

  // Delete source mutation
  const deleteSourceMutation = useMutation({
    mutationFn: (sourceId: string) => api.sources.delete(botId, sourceId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sources", botId] });
      toast.success("Source deleted successfully");
    },
    onError: (error: Error) => {
      toast.error(`Failed to delete source: ${error.message}`);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!url.trim()) {
      toast.error("Please enter a valid URL");
      return;
    }
    setIsSubmitting(true);
    createSourceMutation.mutate(url);
  };

  const handleDelete = (sourceId: string) => {
      deleteSourceMutation.mutate(sourceId);
  };

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-destructive">
              <AlertCircle className="h-5 w-5" />
              Error
            </CardTitle>
            <CardDescription>
              Failed to load sources: {error.message}
            </CardDescription>
          </CardHeader>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto max-w-6xl p-6">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="font-bold text-3xl">Data Sources</h1>
          <p className="mt-1 text-muted-foreground">
            Manage training data sources for your chatbot
          </p>
        </div>

        <Dialog onOpenChange={setIsAddDialogOpen} open={isAddDialogOpen}>
          <DialogTrigger>
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              Add Source
            </Button>
          </DialogTrigger>
          <DialogContent>
            <form onSubmit={handleSubmit}>
              <DialogHeader>
                <DialogTitle>Add Data Source</DialogTitle>
                <DialogDescription>
                  Enter a URL to scrape and add to your bot's knowledge base.
                </DialogDescription>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="grid gap-2">
                  <Label htmlFor="url">URL</Label>
                  <Input
                    onChange={(e) => setUrl(e.target.value)}
                    placeholder="https://example.com"
                    required
                    type="url"
                    value={url}
                  />
                </div>
              </div>
              <DialogFooter>
                <Button
                  disabled={isSubmitting}
                  onClick={() => setIsAddDialogOpen(false)}
                  type="button"
                  variant="outline"
                >
                  Cancel
                </Button>
                <Button disabled={isSubmitting} type="submit">
                  {isSubmitting && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  Add Source
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {!sources || sources.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <p className="mb-4 text-muted-foreground">No sources added yet</p>
            <Button onClick={() => setIsAddDialogOpen(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Add Your First Source
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4">
          {sources.map((source) => (
            <SourceCard
              key={source.id}
              onDelete={() => handleDelete(source.id)}
              source={source}
            />
          ))}
        </div>
      )}
    </div>
  );
}

interface SourceCardProps {
  source: Source;
  onDelete: () => void;
}

function SourceCard({ source, onDelete }: SourceCardProps) {
  const statusConfig: Record<
    SourceStatus,
    {
      label: string;
      variant: "default" | "secondary" | "destructive" | "outline";
    }
  > = {
    pending: { label: "Pending", variant: "secondary" },
    processing: { label: "Processing", variant: "default" },
    completed: { label: "Completed", variant: "outline" },
    failed: { label: "Failed", variant: "destructive" },
  };

  const config = statusConfig[source.status as SourceStatus];

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex-1 space-y-1">
            <div className="flex items-center gap-2">
              <CardTitle className="text-lg">{source.url}</CardTitle>
              <a
                className="inline-flex items-center text-muted-foreground transition-colors hover:text-foreground"
                href={source.url}
                rel="noopener noreferrer"
                target="_blank"
              >
                <ExternalLink className="h-4 w-4" />
              </a>
            </div>
            <CardDescription>
              Added {new Date(source.created_at).toLocaleDateString()}
            </CardDescription>
          </div>

          <div className="flex items-center gap-2">
            <Badge variant={config.variant}>{config.label}</Badge>
            <Button
              className="text-destructive hover:bg-destructive/10 hover:text-destructive"
              onClick={onDelete}
              size="icon"
              variant="ghost"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </CardHeader>

      {source.error_message && (
        <CardContent>
          <div className="flex items-start gap-2 rounded-md bg-destructive/10 p-3 text-destructive">
            <AlertCircle className="mt-0.5 h-4 w-4 flex-shrink-0" />
            <p className="text-sm">{source.error_message}</p>
          </div>
        </CardContent>
      )}

      {source.processed_at && (
        <CardContent>
          <p className="text-muted-foreground text-sm">
            Processed {new Date(source.processed_at).toLocaleString()}
          </p>
        </CardContent>
      )}
    </Card>
  );
}
