import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { formatDistanceToNow } from "date-fns";
import {
  AlertCircle,
  CheckCircle2,
  Clock,
  FileText,
  Link as LinkIcon,
  Loader2,
  Plus,
  Trash2,
} from "lucide-react";
import { toast } from "sonner";
import { api } from "@/api";
import type { Source } from "@/api/index.types";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { SOURCE_STATUS, SOURCE_TYPE } from "@/lib/constants";

export function SourceList({
  botId,
  onAddSource,
}: {
  botId: string;
  onAddSource: () => void;
}) {
  const queryClient = useQueryClient();

  const { data: sources = [], isLoading } = useQuery({
    queryKey: ["sources", botId],
    queryFn: () => api.sources.list(botId),
    refetchInterval: (query) => {
      // Poll every 2 seconds if any source is processing
      const hasProcessing = query.state.data?.some(
        (s: Source) => s.status === SOURCE_STATUS.PROCESSING
      );
      return hasProcessing ? 2000 : false;
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (sourceId: string) => api.sources.delete(botId, sourceId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sources", botId] });
      toast.success("Source deleted successfully");
    },
    onError: (error: Error) => {
      toast.error(`Failed to delete source: ${error.message}`);
    },
  });

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (sources.length === 0) {
    return (
      <Card className="border-dashed">
        <CardContent className="flex h-64 flex-col items-center justify-center space-y-4">
          <div className="flex flex-col items-center space-y-3 text-center">
            <div className="rounded-full bg-muted p-4">
              <FileText className="h-8 w-8 text-muted-foreground" />
            </div>
            <div className="space-y-1">
              <h3 className="font-semibold text-lg">No sources yet</h3>
              <p className="text-muted-foreground text-sm">
                Add your first data source to get started
              </p>
            </div>
          </div>
          <Button onClick={onAddSource} size="lg">
            <Plus className="mr-2 h-4 w-4" />
            Add Source
          </Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-muted-foreground text-sm">
          {sources.length} {sources.length === 1 ? "source" : "sources"}
        </p>
        <Button onClick={onAddSource} size="sm">
          <Plus className="mr-2 h-4 w-4" />
          Add Source
        </Button>
      </div>

      <div className="overflow-hidden rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Type</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Progress</TableHead>
              <TableHead>Added</TableHead>
              <TableHead className="w-[100px]">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {sources.map((source) => (
              <TableRow
                className="transition-colors hover:bg-muted/50"
                key={source.id}
              >
                <TableCell>
                  <SourceTypeIcon type={source.source_type} />
                </TableCell>
                <TableCell className="font-medium">
                  <SourceName source={source} />
                </TableCell>
                <TableCell>
                  <SourceStatusBadge status={source.status} />
                </TableCell>
                <TableCell>
                  <SourceProgress source={source} />
                </TableCell>
                <TableCell className="text-muted-foreground text-sm">
                  {formatDistanceToNow(new Date(source.created_at), {
                    addSuffix: true,
                  })}
                </TableCell>
                <TableCell>
                  <AlertDialog>
                    <AlertDialogTrigger>
                      <Button
                        disabled={deleteMutation.isPending}
                        size="sm"
                        variant="ghost"
                      >
                        <Trash2 className="h-4 w-4 text-destructive" />
                      </Button>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>Delete source?</AlertDialogTitle>
                        <AlertDialogDescription>
                          This will permanently delete this data source and all
                          associated content. This action cannot be undone.
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction
                          className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                          onClick={() => deleteMutation.mutate(source.id)}
                        >
                          Delete
                        </AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}

function SourceTypeIcon({ type }: { type: string }) {
  switch (type) {
    case SOURCE_TYPE.URL:
      return (
        <div className="flex items-center gap-2">
          <div className="rounded-md bg-chart-1/10 p-1.5">
            <LinkIcon className="h-4 w-4 text-chart-1" />
          </div>
          <span className="text-sm">URL</span>
        </div>
      );
    case SOURCE_TYPE.FILE:
      return (
        <div className="flex items-center gap-2">
          <div className="rounded-md bg-chart-2/10 p-1.5">
            <FileText className="h-4 w-4 text-chart-2" />
          </div>
          <span className="text-sm">File</span>
        </div>
      );
    case SOURCE_TYPE.TEXT:
      return (
        <div className="flex items-center gap-2">
          <div className="rounded-md bg-chart-3/10 p-1.5">
            <FileText className="h-4 w-4 text-chart-3" />
          </div>
          <span className="text-sm">Text</span>
        </div>
      );
    default:
      return <span className="text-sm">{type}</span>;
  }
}

function SourceName({ source }: { source: Source }) {
  if (source.source_type === SOURCE_TYPE.URL) {
    return (
      <a
        className="hover:underline"
        href={source.url}
        rel="noopener noreferrer"
        target="_blank"
      >
        {source.url}
      </a>
    );
  }

  if (source.original_filename) {
    return <span>{source.original_filename}</span>;
  }

  return <span className="text-muted-foreground">Unnamed source</span>;
}

function SourceStatusBadge({ status }: { status: string }) {
  switch (status) {
    case SOURCE_STATUS.PENDING:
      return (
        <Badge className="gap-1" variant="secondary">
          <Clock className="h-3 w-3" />
          Pending
        </Badge>
      );
    case SOURCE_STATUS.PROCESSING:
      return (
        <Badge className="gap-1" variant="secondary">
          <Loader2 className="h-3 w-3 animate-spin" />
          Processing
        </Badge>
      );
    case SOURCE_STATUS.COMPLETED:
      return (
        <Badge
          className="gap-1 border-chart-2 bg-chart-2/10 text-chart-2"
          variant="outline"
        >
          <CheckCircle2 className="h-3 w-3" />
          Completed
        </Badge>
      );
    case SOURCE_STATUS.FAILED:
      return (
        <Badge className="gap-1" variant="destructive">
          <AlertCircle className="h-3 w-3" />
          Failed
        </Badge>
      );
    default:
      return <Badge variant="outline">{status}</Badge>;
  }
}

function SourceProgress({ source }: { source: Source }) {
  if (source.status === SOURCE_STATUS.COMPLETED) {
    return <span className="text-muted-foreground text-sm">100%</span>;
  }

  if (source.status === SOURCE_STATUS.FAILED) {
    return (
      <span className="text-destructive text-sm" title={source.error_message}>
        Error
      </span>
    );
  }

  if (source.status === SOURCE_STATUS.PROCESSING) {
    return (
      <div className="flex items-center gap-2">
        <Progress className="w-24" value={source.processing_progress} />
        <span className="text-muted-foreground text-sm">
          {source.processing_progress}%
        </span>
      </div>
    );
  }

  return <span className="text-muted-foreground text-sm">-</span>;
}
