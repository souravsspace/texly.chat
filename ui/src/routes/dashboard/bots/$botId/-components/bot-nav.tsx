import { Link } from "@tanstack/react-router";
import { Code, Database, MessageSquare } from "lucide-react";
import { cn } from "@/lib/utils";

interface BotNavProps {
  botId: string;
  currentPath: "configure" | "widget" | "chat";
}

export function BotNav({ botId, currentPath }: BotNavProps) {
  const navItems = [
    {
      path: "configure" as const,
      label: "Configure",
      icon: Database,
      href: `/dashboard/bots/${botId}/configure`,
    },
    {
      path: "widget" as const,
      label: "Widget",
      icon: Code,
      href: `/dashboard/bots/${botId}/widget`,
    },
    {
      path: "chat" as const,
      label: "Chat",
      icon: MessageSquare,
      href: `/dashboard/bots/${botId}/chat`,
    },
  ];

  return (
    <div className="border-b">
      <nav aria-label="Bot navigation" className="flex gap-1">
        {navItems.map((item) => {
          const Icon = item.icon;
          const isActive = currentPath === item.path;

          return (
            <Link
              className={cn(
                "flex items-center gap-2 border-b-2 px-4 py-3 font-medium text-sm transition-colors",
                isActive
                  ? "border-primary text-primary"
                  : "border-transparent text-muted-foreground hover:text-foreground"
              )}
              key={item.path}
              to={item.href}
            >
              <Icon className="h-4 w-4" />
              {item.label}
            </Link>
          );
        })}
      </nav>
    </div>
  );
}
