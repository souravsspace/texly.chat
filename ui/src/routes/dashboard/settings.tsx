import { createFileRoute } from "@tanstack/react-router";
import { Icons } from "@/components/icons";
import { Header } from "@/components/layout/header";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useAuth } from "@/providers/auth";

export const Route = createFileRoute("/dashboard/settings")({
  component: SettingsPage,
});

function SettingsPage() {
  const { user } = useAuth();

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
      .slice(0, 2);
  };

  const isGoogleConnected =
    user?.auth_provider === "google" || !!user?.google_id;

  return (
    <div className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
      <Header
        description="Manage your account settings and preferences."
        title="Settings"
      />

      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>Profile</CardTitle>
            <CardDescription>
              Your personal information and account details.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="flex items-center gap-6">
              <Avatar className="h-20 w-20">
                <AvatarImage alt={user?.name} src={user?.avatar} />
                <AvatarFallback className="text-lg">
                  {user?.name ? getInitials(user.name) : "??"}
                </AvatarFallback>
              </Avatar>
              <div className="space-y-1">
                <h3 className="font-medium text-lg">{user?.name}</h3>
                <p className="text-muted-foreground">{user?.email}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Connected Accounts</CardTitle>
            <CardDescription>
              Manage your linked social accounts for easier login.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between rounded-lg border p-4">
              <div className="flex items-center gap-4">
                <div className="flex h-10 w-10 items-center justify-center rounded-full bg-muted">
                  <Icons.google className="h-5 w-5" />
                </div>
                <div>
                  <p className="font-medium">Google</p>
                  <p className="text-muted-foreground text-sm">
                    {isGoogleConnected ? "Connected" : "Not connected"}
                  </p>
                </div>
              </div>
              {isGoogleConnected ? (
                <Badge
                  className="border-green-500 text-green-500"
                  variant="outline"
                >
                  Connected
                </Badge>
              ) : (
                <Badge variant="outline">Not connected</Badge>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
