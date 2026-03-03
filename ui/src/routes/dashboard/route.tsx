import { createFileRoute, Outlet } from "@tanstack/react-router";
import { Header } from "@/components/layout/header";

export const Route = createFileRoute("/dashboard")({
  component: DashboardLayout,
});

function DashboardLayout() {
  return (
    <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <Header />
      <Outlet />
    </div>
  );
}
