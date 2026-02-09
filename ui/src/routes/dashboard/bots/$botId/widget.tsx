import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { Bot, ChevronRight, Home } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { api } from "@/api";
import type { UpdateBotRequest, WidgetConfig } from "@/api/index.types";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { BotNav } from "./-components/bot-nav";
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

  const { data: bot, isLoading } = useQuery({
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

  if (isLoading) {
    return (
      <div className="container mx-auto max-w-7xl space-y-6 py-8">
        <Skeleton className="h-6 w-64" />
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <Skeleton className="h-[600px] w-full" />
          <Skeleton className="h-[600px] w-full" />
        </div>
      </div>
    );
  }

  if (!bot) {
    return (
      <div className="container mx-auto max-w-7xl py-8">
        <Card>
          <CardContent className="flex h-[50vh] items-center justify-center">
            <div className="space-y-2 text-center">
              <Bot className="mx-auto h-12 w-12 text-muted-foreground" />
              <p className="text-lg text-muted-foreground">Bot not found</p>
            </div>
          </CardContent>
        </Card>
      </div>
    );
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
    <div className="container mx-auto max-w-7xl space-y-6 py-8">
      {/* Breadcrumb Navigation */}
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink href="/dashboard">
              <Home className="h-4 w-4" />
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator>
            <ChevronRight className="h-4 w-4" />
          </BreadcrumbSeparator>
          <BreadcrumbItem>
            <BreadcrumbLink href="/dashboard">Bots</BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator>
            <ChevronRight className="h-4 w-4" />
          </BreadcrumbSeparator>
          <BreadcrumbItem>
            <BreadcrumbPage>{bot.name}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      {/* Header Section */}
      <div className="space-y-3">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
            <Bot className="h-6 w-6 text-primary" />
          </div>
          <div>
            <h1 className="font-bold text-3xl tracking-tight">{bot.name}</h1>
            <p className="text-muted-foreground text-sm">
              Customize and embed your chatbot widget
            </p>
          </div>
        </div>
      </div>

      {/* Navigation Tabs */}
      <BotNav botId={botId} currentPath="widget" />

      {/* Main Content - Two Column Layout */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Left Column - Configuration Form */}
        <div className="space-y-6">
          <WidgetConfigForm
            initialConfig={initialWidgetConfig}
            initialOrigins={allowedOrigins}
            isSubmitting={updateBotMutation.isPending}
            onConfigChange={setPreviewConfig}
            onSubmit={handleConfigSubmit}
          />

          {/* Embed Code Section */}
          <EmbedCodeDisplay botId={botId} />
        </div>

        {/* Right Column - Live Preview */}
        <div className="lg:sticky lg:top-6 lg:h-fit">
          <WidgetPreview botName={bot.name} config={previewConfig} />
        </div>
      </div>
    </div>
  );
}
