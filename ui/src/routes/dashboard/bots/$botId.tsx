import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet, useLocation } from "@tanstack/react-router";
import { Bot, ChevronRight, Home } from "lucide-react";
import { api } from "@/api";
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
import { BotNav } from "./$botId/-components/bot-nav";

export const Route = createFileRoute("/dashboard/bots/$botId")({
  component: BotLayout,
});

function BotLayout() {
  const { botId } = Route.useParams();
  const location = useLocation();

  const { data: bot, isLoading } = useQuery({
    queryKey: ["bot", botId],
    queryFn: () => api.bots.get(botId),
  });

  if (isLoading) {
    return (
      <div className="container mx-auto max-w-6xl space-y-6 py-8">
        <Skeleton className="h-6 w-64" />
        <div className="space-y-4">
          <Skeleton className="h-10 w-full" />
          <Skeleton className="h-64 w-full" />
        </div>
      </div>
    );
  }

  if (!bot) {
    return (
      <div className="container mx-auto max-w-6xl py-8">
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

  const currentPath = getCurrentPath(location.pathname);
  const description = getPageDescription(currentPath);

  return (
    <div className="container mx-auto max-w-6xl space-y-6 py-8">
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
            <p className="text-muted-foreground text-sm">{description}</p>
          </div>
        </div>
      </div>

      {/* Navigation Tabs */}
      <BotNav botId={botId} currentPath={currentPath} />

      {/* Child Routes */}
      <Outlet />
    </div>
  );
}

function getCurrentPath(pathname: string): "configure" | "chat" | "widget" {
  if (pathname.includes("/chat")) return "chat";
  if (pathname.includes("/widget")) return "widget";
  return "configure";
}

function getPageDescription(
  currentPath: "configure" | "chat" | "widget"
): string {
  switch (currentPath) {
    case "chat":
      return "Test your chatbot and see how it responds";
    case "widget":
      return "Customize and embed your chatbot widget";
    default:
      return "Configure data sources for your chatbot";
  }
}
