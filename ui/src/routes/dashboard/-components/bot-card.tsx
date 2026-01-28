import type { Bot } from "@/api/index.types";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

interface BotCardProps {
  bot: Bot;
  onDelete: (id: string) => void;
  isDeleting?: boolean;
}

export function BotCard({ bot, onDelete, isDeleting }: BotCardProps) {
  return (
    <Card className="transition-shadow hover:shadow-md">
      <CardHeader>
        <CardTitle>{bot.name}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="line-clamp-3 text-muted-foreground text-sm">
          {bot.system_prompt || "No system prompt configured."}
        </p>
      </CardContent>
      <CardFooter className="flex justify-between border-t pt-4">
        <small className="text-muted-foreground">
          {new Date(bot.created_at).toLocaleDateString()}
        </small>
        <div className="flex gap-2">
          <Button
            className="h-auto p-2 text-destructive hover:bg-transparent hover:text-destructive/80"
            disabled={isDeleting}
            onClick={() => onDelete(bot.id)}
            size="sm"
            variant="ghost"
          >
            {isDeleting ? "Deleting..." : "Delete"}
          </Button>
        </div>
      </CardFooter>
    </Card>
  );
}
