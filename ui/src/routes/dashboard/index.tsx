import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useEffect } from "react";
import { api } from "@/api";
import { ModeToggle } from "@/components/mode-toggle";
import { Button } from "@/components/ui/button";
import { useCreateBotDialog } from "@/hooks/use-create-bot-dialog";
import { useConfirm } from "@/providers/alert-dialog";
import { useAuth } from "@/providers/auth";
import { BotList } from "@/routes/dashboard/_components/bot-list";
import { CreateBotDialog } from "@/routes/dashboard/_components/create-bot-dialog";

export const Route = createFileRoute("/dashboard/")({
  component: Dashboard,
});

function Dashboard() {
  const navigate = useNavigate();
  const { user, token, logout, loading: authLoading } = useAuth();
  const queryClient = useQueryClient();
  const createBotDialog = useCreateBotDialog();
  const confirm = useConfirm();

  useEffect(() => {
    if (!(authLoading || token)) {
      navigate({ to: "/login" });
    }
  }, [authLoading, token, navigate]);

  const {
    data: bots,
    isLoading: botsLoading,
    error: botsError,
  } = useQuery({
    queryKey: ["bots"],
    queryFn: () => api.bots.list(),
    enabled: !!user,
  });

  const deleteBotMutation = useMutation({
    mutationFn: (id: string) => api.bots.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bots"] });
    },
  });

  const handleDelete = async (id: string) => {
    const ok = await confirm.confirm({
      title: "Delete Bot",
      description:
        "Are you sure you want to delete this bot? This action cannot be undone.",
      confirmText: "Delete",
      variant: "destructive",
    });
    if (ok) {
      deleteBotMutation.mutate(id);
    }
  };

  if (authLoading || (!user && botsLoading)) {
    return (
      <div className="flex h-screen items-center justify-center">
        Loading...
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
      <header className="mb-8 flex items-center justify-between border-border border-b-2 pb-6">
        <div>
          <h1 className="mb-2 font-bold text-4xl text-foreground">Dashboard</h1>
          {user && (
            <p className="text-muted-foreground">Welcome, {user.name}!</p>
          )}
        </div>
        <div className="flex gap-3">
          <ModeToggle />
          <Button variant="secondary">
            <Link to="/">Home</Link>
          </Button>
          <Button onClick={() => logout()} variant="secondary">
            Logout
          </Button>
        </div>
      </header>

      {botsError && (
        <div className="mb-6 rounded-lg border border-destructive bg-destructive/15 px-4 py-3 text-destructive text-sm">
          Failed to load bots
        </div>
      )}

      <div className="mb-6 flex items-center justify-between">
        <h2 className="font-bold text-2xl text-foreground">Your Chatbots</h2>
        <Button onClick={createBotDialog.onOpen}>New Chatbot</Button>
        <CreateBotDialog
          onOpenChange={(open) =>
            open ? createBotDialog.onOpen() : createBotDialog.onClose()
          }
          open={createBotDialog.isOpen}
        />
      </div>

      <BotList
        bots={bots}
        deletingId={deleteBotMutation.variables}
        isLoading={botsLoading}
        onDelete={handleDelete}
      />
    </div>
  );
}
