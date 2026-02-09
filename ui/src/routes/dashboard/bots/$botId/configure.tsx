import { createFileRoute } from "@tanstack/react-router";
import { SourceManager } from "./-components/source-manager";

export const Route = createFileRoute("/dashboard/bots/$botId/configure")({
  component: ConfigurePage,
});

function ConfigurePage() {
  const { botId } = Route.useParams();

  return <SourceManager botId={botId} />;
}
