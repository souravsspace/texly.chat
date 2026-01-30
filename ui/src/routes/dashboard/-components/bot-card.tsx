import { Link } from "@tanstack/react-router";
import { formatDistanceToNow } from "date-fns";
import { MessageSquare, Settings, Trash2 } from "lucide-react";
import type { Bot } from "@/api/index.types";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

interface BotCardProps {
  bot: Bot;
  onDelete: (id: string) => void;
  isDeleting?: boolean;
}

export function BotCard({ bot, onDelete, isDeleting }: BotCardProps) {
  return (
    <Card className="group relative overflow-hidden border-primary/50 transition-all duration-300 hover:shadow-lg">
      {/* Subtle gradient accent */}
      <div className="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-primary/50 via-primary to-primary/50 opacity-0 transition-opacity duration-300 group-hover:opacity-100" />

      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-2">
          <CardTitle className="text-xl">{bot.name}</CardTitle>
          <Badge className="text-xs" variant="secondary">
            Active
          </Badge>
        </div>
      </CardHeader>

      <CardContent className="pb-4">
        <p className="line-clamp-3 text-muted-foreground text-sm leading-relaxed">
          {bot.system_prompt || "No system prompt configured."}
        </p>
      </CardContent>

      <Separator />

      <CardFooter className="flex flex-col gap-3 pt-4">
        <div className="flex w-full items-center justify-between text-muted-foreground text-xs">
          <span>
            Created{" "}
            {formatDistanceToNow(new Date(bot.created_at), { addSuffix: true })}
          </span>

          <Button
            disabled={isDeleting}
            onClick={() => onDelete(bot.id)}
            size="sm"
            variant="outline"
          >
            <Trash2 className="h-4 w-4 text-destructive" />
          </Button>
        </div>

        <div className="flex w-full gap-2">
          <Link
            className="flex-1"
            params={{ botId: bot.id }}
            to="/dashboard/bots/$botId/configure"
          >
            <Button
              className="w-full transition-colors hover:bg-accent"
              size="sm"
              variant="outline"
            >
              <Settings className="mr-2 h-4 w-4" />
              Configure
            </Button>
          </Link>

          <Link
            className="flex-1"
            params={{ botId: bot.id }}
            to="/dashboard/bots/$botId/chat"
          >
            <Button className="w-full" size="sm" variant="default">
              <MessageSquare className="mr-2 h-4 w-4" />
              Chat
            </Button>
          </Link>
        </div>
      </CardFooter>
    </Card>
  );
}
