import type { Bot } from "@/api/index.types";
import { BotCard } from "./bot-card";

interface BotListProps {
  bots: Bot[] | undefined;
  isLoading: boolean;
  onDelete: (id: string) => void;
  deletingId?: string | null;
}

export function BotList({
  bots,
  isLoading,
  onDelete,
  deletingId,
}: BotListProps) {
  if (isLoading) {
    return (
      <p className="py-12 text-center text-muted-foreground">Loading bots...</p>
    );
  }

  if (!bots || bots.length === 0) {
    return (
      <p className="py-12 text-center text-muted-foreground">
        No bots yet. Create your first chatbot!
      </p>
    );
  }

  return (
    <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {bots.map((bot) => (
        <BotCard
          bot={bot}
          isDeleting={deletingId === bot.id}
          key={bot.id}
          onDelete={onDelete}
        />
      ))}
    </div>
  );
}
