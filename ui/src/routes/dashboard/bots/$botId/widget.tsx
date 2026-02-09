import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { api } from "@/api";
import type { UpdateBotRequest, WidgetConfig } from "@/api/index.types";
import { EmbedCodeDisplay } from "./-components/embed-code-display";
import { WidgetConfigForm } from "./-components/widget-config-form";
import { WidgetPreview } from "./-components/widget-preview";

export const Route = createFileRoute("/dashboard/bots/$botId/widget")({
  component: WidgetPage,
});

function WidgetPage() {
  const { botId } = Route.useParams();
  const queryClient = useQueryClient();

  // State for live preview updates
  const [previewConfig, setPreviewConfig] = useState<WidgetConfig>({
    theme_color: "#6366f1",
    initial_message: "Hi! How can I help you today?",
    position: "bottom-right",
    bot_avatar: "",
  });

  const { data: bot } = useQuery({
    queryKey: ["bot", botId],
    queryFn: () => api.bots.get(botId),
  });

  const updateBotMutation = useMutation({
    mutationFn: (data: UpdateBotRequest) => api.bots.update(botId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bot", botId] });
      toast.success("Widget configuration updated successfully");
    },
    onError: (error: Error) => {
      toast.error(`Failed to update widget: ${error.message}`);
    },
  });

  // Parse and set initial config when bot loads
  // biome-ignore lint/correctness/useExhaustiveDependencies: We only want to update when widget_config changes
  useEffect(() => {
    if (!bot) return;

    try {
      if (bot.widget_config) {
        const parsed = JSON.parse(bot.widget_config) as WidgetConfig;
        setPreviewConfig(parsed);
      }
    } catch (error) {
      console.error("Failed to parse widget config:", error);
    }
  }, [bot?.widget_config]);

  const handleConfigSubmit = (
    widgetConfig: WidgetConfig,
    allowedOrigins: string[]
  ) => {
    if (!bot) return;

    const updateData: UpdateBotRequest = {
      name: bot.name,
      system_prompt: bot.system_prompt,
      widget_config: widgetConfig,
      allowed_origins: allowedOrigins,
    };

    updateBotMutation.mutate(updateData);
  };

  if (!bot) {
    return null;
  }

  // Parse widget config and allowed origins from JSON strings
  let initialWidgetConfig: WidgetConfig = {
    theme_color: "#6366f1",
    initial_message: "Hi! How can I help you today?",
    position: "bottom-right",
    bot_avatar: "",
  };

  let allowedOrigins: string[] = [];

  try {
    if (bot.widget_config) {
      initialWidgetConfig = JSON.parse(bot.widget_config) as WidgetConfig;
    }
  } catch (error) {
    console.error("Failed to parse widget config:", error);
  }

  try {
    if (bot.allowed_origins) {
      allowedOrigins = JSON.parse(bot.allowed_origins) as string[];
    }
  } catch (error) {
    console.error("Failed to parse allowed origins:", error);
  }

  return (
    <div className="mx-auto max-w-6xl">
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Left Column - Configuration Form */}
        <WidgetConfigForm
          initialConfig={initialWidgetConfig}
          initialOrigins={allowedOrigins}
          isSubmitting={updateBotMutation.isPending}
          onConfigChange={setPreviewConfig}
          onSubmit={handleConfigSubmit}
        />

        {/* Right Column - Live Preview */}
        <div className="space-y-6">
          <div className="lg:sticky lg:top-6 lg:h-fit">
            <WidgetPreview botName={bot.name} config={previewConfig} />
          </div>
          {/* Embed Code Section */}
          <EmbedCodeDisplay botId={botId} />
        </div>
      </div>
    </div>
  );
}
