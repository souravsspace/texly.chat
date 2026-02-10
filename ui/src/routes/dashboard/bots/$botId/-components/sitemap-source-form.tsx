import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Globe, Loader2 } from "lucide-react";
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

export function SitemapSourceForm({
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
    mutationFn: (url: string) => api.sources.createSitemap(botId, url),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["sources", botId] });
      toast.success(
        data.message || `Created ${data.created_count} sources from sitemap`
      );
      setUrl("");
      onSuccess();
    },
    onError: (error: Error) => {
      toast.error(`Failed to crawl sitemap: ${error.message}`);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!url.trim()) return;

    // Basic URL validation
    try {
      const parsedUrl = new URL(url);
      // Ensure URL has http or https scheme
      if (parsedUrl.protocol !== "http:" && parsedUrl.protocol !== "https:") {
        toast.error("URL must start with http:// or https://");
        return;
      }
      mutation.mutate(url);
    } catch {
      toast.error("Please enter a valid URL starting with http:// or https://");
    }
  };

  return (
    <Card className="border-2">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <div className="rounded-md bg-primary/10 p-2">
            <Globe className="h-5 w-5 text-primary" />
          </div>
          Crawl Entire Website
        </CardTitle>
        <CardDescription>
          Automatically discover and index all pages from a website's sitemap
        </CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit}>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="sitemap-url">Website URL or Sitemap URL</Label>
            <div className="relative">
              <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                <Globe className="h-4 w-4 text-muted-foreground" />
              </div>
              <Input
                className="pl-9"
                disabled={mutation.isPending}
                id="sitemap-url"
                onChange={(e) => setUrl(e.target.value)}
                placeholder="https://example.com or https://example.com/sitemap.xml"
                required
                type="url"
                value={url}
              />
            </div>
            <p className="text-muted-foreground text-xs">
              Enter a website URL or a direct link to sitemap.xml. We'll
              automatically discover the sitemap.
            </p>
          </div>

          {mutation.isError && (
            <Alert variant="destructive">
              <AlertDescription>
                {mutation.error?.message || "Failed to crawl sitemap"}
              </AlertDescription>
            </Alert>
          )}

          <Alert>
            <AlertDescription className="text-xs">
              <strong>Note:</strong> This will crawl and index multiple pages
              from the website. Processing may take several minutes depending on
              the number of pages (max 1000 URLs). Each page will be scraped,
              chunked, and embedded for semantic search.
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
            {mutation.isPending ? "Discovering URLs..." : "Crawl Website"}
          </Button>
        </CardFooter>
      </form>
    </Card>
  );
}
