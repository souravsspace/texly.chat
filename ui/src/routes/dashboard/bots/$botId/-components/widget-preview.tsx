import { Bot, Send, X } from "lucide-react";
import { useEffect, useState } from "react";
import type { WidgetConfig } from "@/api/index.types";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { cn } from "@/lib/utils";

interface WidgetPreviewProps {
  config: WidgetConfig;
  botName: string;
}

export function WidgetPreview({ config, botName }: WidgetPreviewProps) {
  const [isOpen, setIsOpen] = useState(true);
  const [previewMessage, setPreviewMessage] = useState("");

  // Mock messages for preview
  const mockMessages = [
    {
      id: 1,
      role: "assistant",
      content: config.initial_message || "Hi! How can I help you today?",
      timestamp: new Date(),
    },
  ];

  // Reset open state when position changes to show animation
  // biome-ignore lint/correctness/useExhaustiveDependencies: Intentionally trigger on position change
  useEffect(() => {
    setIsOpen(false);
    const timer = setTimeout(() => setIsOpen(true), 300);
    return () => clearTimeout(timer);
  }, [config.position]);

  const positionClasses =
    config.position === "bottom-left" ? "bottom-4 left-4" : "bottom-4 right-4";

  const launcherPositionClasses =
    config.position === "bottom-left" ? "bottom-6 left-6" : "bottom-6 right-6";

  return (
    <Card className="overflow-hidden">
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Live Preview</CardTitle>
          <Badge variant="secondary">Interactive Demo</Badge>
        </div>
      </CardHeader>
      <CardContent>
        {/* Mock Website Background */}
        <div className="relative h-[600px] overflow-hidden rounded-lg border-2 border-dashed bg-muted/30 p-6">
          <div className="space-y-4">
            <div className="h-4 w-3/4 rounded bg-muted" />
            <div className="h-4 w-1/2 rounded bg-muted" />
            <div className="h-4 w-2/3 rounded bg-muted" />
            <div className="mt-8 space-y-2">
              <div className="h-3 w-full rounded bg-muted" />
              <div className="h-3 w-full rounded bg-muted" />
              <div className="h-3 w-5/6 rounded bg-muted" />
            </div>
          </div>

          {/* Widget Launcher Button */}
          {!isOpen && (
            <button
              className={cn(
                "absolute z-50 flex h-14 w-14 items-center justify-center rounded-full shadow-lg transition-all hover:scale-110",
                launcherPositionClasses
              )}
              onClick={() => setIsOpen(true)}
              style={{ backgroundColor: config.theme_color }}
              type="button"
            >
              <Bot className="h-6 w-6 text-white" />
            </button>
          )}

          {/* Widget Window */}
          {isOpen && (
            <div
              className={cn(
                "absolute z-50 flex w-[380px] flex-col overflow-hidden rounded-2xl bg-background shadow-2xl transition-all",
                positionClasses
              )}
              style={{ height: "500px", maxHeight: "85vh" }}
            >
              {/* Widget Header */}
              <div
                className="flex items-center justify-between px-4 py-3"
                style={{ backgroundColor: config.theme_color }}
              >
                <div className="flex items-center gap-3">
                  {config.bot_avatar ? (
                    <img
                      alt={botName}
                      className="h-8 w-8 rounded-full"
                      height={32}
                      src={config.bot_avatar}
                      width={32}
                    />
                  ) : (
                    <div className="flex h-8 w-8 items-center justify-center rounded-full bg-white/20">
                      <Bot className="h-5 w-5 text-white" />
                    </div>
                  )}
                  <div>
                    <h3 className="font-semibold text-sm text-white">
                      {botName}
                    </h3>
                    <p className="text-white/80 text-xs">Online</p>
                  </div>
                </div>
                <button
                  className="rounded-full p-1 text-white/80 transition-colors hover:bg-white/20 hover:text-white"
                  onClick={() => setIsOpen(false)}
                  type="button"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>

              {/* Messages Area */}
              <div className="flex-1 space-y-4 overflow-y-auto p-4">
                {mockMessages.map((message) => (
                  <div className="flex gap-3" key={message.id}>
                    <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
                      <Bot className="h-4 w-4 text-primary" />
                    </div>
                    <div className="flex-1">
                      <div className="rounded-lg bg-muted px-3 py-2">
                        <p className="text-sm">{message.content}</p>
                      </div>
                      <p className="mt-1 text-muted-foreground text-xs">
                        {message.timestamp.toLocaleTimeString([], {
                          hour: "2-digit",
                          minute: "2-digit",
                        })}
                      </p>
                    </div>
                  </div>
                ))}

                {/* Placeholder for user message */}
                {previewMessage && (
                  <div className="flex justify-end">
                    <div className="max-w-[80%]">
                      <div
                        className="rounded-lg px-3 py-2 text-white"
                        style={{ backgroundColor: config.theme_color }}
                      >
                        <p className="text-sm">{previewMessage}</p>
                      </div>
                    </div>
                  </div>
                )}
              </div>

              {/* Input Area */}
              <div className="border-t p-4">
                <div className="flex gap-2">
                  <input
                    className="flex-1 rounded-lg border bg-background px-3 py-2 text-sm outline-none transition-colors focus:border-primary"
                    onChange={(e) => setPreviewMessage(e.target.value)}
                    placeholder="Type a message..."
                    type="text"
                    value={previewMessage}
                  />
                  <button
                    className="flex h-9 w-9 items-center justify-center rounded-lg transition-colors"
                    style={{ backgroundColor: config.theme_color }}
                    type="button"
                  >
                    <Send className="h-4 w-4 text-white" />
                  </button>
                </div>
                <p className="mt-2 text-center text-muted-foreground text-xs">
                  Powered by Texly
                </p>
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
