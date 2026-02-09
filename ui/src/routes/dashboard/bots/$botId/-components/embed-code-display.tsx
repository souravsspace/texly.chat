import { Check, Code2, Copy } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface EmbedCodeDisplayProps {
  botId: string;
}

export function EmbedCodeDisplay({ botId }: EmbedCodeDisplayProps) {
  const [copied, setCopied] = useState(false);

  // Generate embed code
  const embedCode = `<script src="${window.location.origin}/widget.js" data-bot-id="${botId}"></script>`;

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(embedCode);
      setCopied(true);
      toast.success("Embed code copied to clipboard");
      setTimeout(() => setCopied(false), 2000);
    } catch {
      toast.error("Failed to copy to clipboard");
    }
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center gap-2">
          <Code2 className="h-5 w-5 text-primary" />
          <CardTitle>Embed Code</CardTitle>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <p className="text-muted-foreground text-sm">
          Add this code to your website before the closing{" "}
          <code className="rounded bg-muted px-1 py-0.5 font-mono text-xs">
            {"</body>"}
          </code>{" "}
          tag to display your chatbot widget.
        </p>

        {/* Code Block */}
        <div className="relative">
          <pre className="overflow-x-auto border bg-muted p-4">
            <code className="font-mono text-sm">{embedCode}</code>
          </pre>

          {/* Copy Button */}
          <Button
            className="absolute top-2 right-2"
            onClick={handleCopy}
            size="sm"
            variant="ghost"
          >
            {copied ? (
              <>
                <Check className="mr-2 h-4 w-4" />
                Copied
              </>
            ) : (
              <>
                <Copy className="mr-2 h-4 w-4" />
                Copy
              </>
            )}
          </Button>
        </div>

        {/* Instructions */}
        <div className="space-y-2 rounded-lg border bg-muted/50 p-4">
          <h4 className="font-medium text-sm">Installation Steps</h4>
          <ol className="ml-4 list-decimal space-y-1 text-muted-foreground text-sm">
            <li>Copy the embed code above</li>
            <li>
              Paste it in your website's HTML before the{" "}
              <code className="rounded bg-muted px-1 py-0.5 font-mono text-xs">
                {"</body>"}
              </code>{" "}
              tag
            </li>
            <li>Save and refresh your website</li>
            <li>The widget will appear in the configured position</li>
          </ol>
        </div>

        {/* Security Note */}
        <div className="rounded-lg border border-orange-200 bg-orange-50 p-3 dark:border-orange-900 dark:bg-orange-950/30">
          <p className="text-orange-900 text-sm dark:text-orange-200">
            <strong>Security:</strong> Make sure to add your website domain to
            the "Allowed Origins" field above to enable the widget on your site.
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
