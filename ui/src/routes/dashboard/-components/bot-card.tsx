import { Link } from "@tanstack/react-router";
import { formatDistanceToNow } from "date-fns";
import {
  Calendar,
  MessageSquare,
  MoreHorizontal,
  Settings,
  Trash2,
} from "lucide-react";
import type { Bot } from "@/api/index.types";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

interface BotCardProps {
  bot: Bot;
  onDelete: (id: string) => void;
  isDeleting?: boolean;
}

export function BotCard({ bot, onDelete, isDeleting }: BotCardProps) {
  return (
    <Card className="group flex h-full flex-col overflow-hidden border-border/60 bg-card transition-all hover:border-primary/50 hover:shadow-md">
      <CardHeader className="flex flex-row items-start justify-between space-y-0 p-4 pb-2">
        <div className="flex flex-col gap-1.5">
          <div className="flex items-center gap-2">
            <h3 className="line-clamp-1 font-semibold text-lg tracking-tight">
              {bot.name}
            </h3>
            <Badge
              className="px-1.5 py-0 font-normal text-[10px]"
              variant="secondary"
            >
              Active
            </Badge>
          </div>
          <div className="flex items-center gap-1 text-muted-foreground text-xs">
            <Calendar className="h-3 w-3" />
            <span>
              {formatDistanceToNow(new Date(bot.created_at), {
                addSuffix: true,
              })}
            </span>
          </div>
        </div>

        <DropdownMenu>
          <DropdownMenuTrigger>
            <Button
              className="h-8 w-8 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100 data-[state=open]:opacity-100"
              size="icon"
              variant="ghost"
            >
              <MoreHorizontal className="h-4 w-4" />
              <span className="sr-only">Open menu</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-[160px]">
            <Link
              params={{ botId: bot.id }}
              to="/dashboard/bots/$botId/configure"
            >
              <DropdownMenuItem>
                <Settings className="mr-2 h-4 w-4" />
                Settings
              </DropdownMenuItem>
            </Link>
            <DropdownMenuItem
              className="text-destructive focus:text-destructive"
              disabled={isDeleting}
              onClick={() => onDelete(bot.id)}
            >
              <Trash2 className="mr-2 h-4 w-4" />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </CardHeader>

      <CardContent className="flex-1 p-4 pt-2">
        <div className="bg-muted/50 p-3 text-muted-foreground text-sm">
          <div className="mb-2 flex items-center gap-1.5 font-medium text-foreground/80 text-xs">
            <MessageSquare className="h-3.5 w-3.5" />
            System Prompt
          </div>
          <p className="line-clamp-3 text-xs leading-relaxed opacity-90">
            {bot.system_prompt || (
              <span className="italic opacity-70">
                No system prompt configured.
              </span>
            )}
          </p>
        </div>
      </CardContent>

      <CardFooter className="bg-muted/20 p-4 pt-0">
        <Link
          className="w-full"
          params={{ botId: bot.id }}
          to="/dashboard/bots/$botId/configure"
        >
          <Button className="w-full gap-2" size="sm" variant="outline">
            <Settings className="h-4 w-4" />
            Manage Chatbot
          </Button>
        </Link>
      </CardFooter>
    </Card>
  );
}
