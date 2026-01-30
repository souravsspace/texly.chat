import { useMutation, useQueryClient } from "@tanstack/react-query";
import { FileUp, Loader2, X } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { api } from "@/api";
import {
  Dropzone,
  DropzoneContent,
  DropzoneEmptyState,
} from "@/components/kibo-ui/dropzone";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import {
  MAX_FILE_SIZE_BYTES,
  SUPPORTED_FILE_MIME_TYPES,
} from "@/lib/constants";

export function FileUploadForm({
  botId,
  onSuccess,
  onCancel,
}: {
  botId: string;
  onSuccess: () => void;
  onCancel: () => void;
}) {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploadProgress, setUploadProgress] = useState(0);
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (file: File) =>
      api.sources.uploadFile(botId, file, (progress) => {
        setUploadProgress(progress);
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sources", botId] });
      toast.success("File uploaded successfully");
      setSelectedFile(null);
      setUploadProgress(0);
      onSuccess();
    },
    onError: (error: Error) => {
      toast.error(`Failed to upload file: ${error.message}`);
      setUploadProgress(0);
    },
  });

  const handleFileSelect = (files: File[]) => {
    if (files.length > 0) {
      const file = files[0];

      // Validate file size
      if (file.size > MAX_FILE_SIZE_BYTES) {
        toast.error(
          `File size exceeds ${MAX_FILE_SIZE_BYTES / 1024 / 1024}MB limit`
        );
        return;
      }

      setSelectedFile(file);
    }
  };

  const handleUpload = () => {
    if (!selectedFile) return;
    mutation.mutate(selectedFile);
  };

  const handleRemove = () => {
    setSelectedFile(null);
    setUploadProgress(0);
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${Math.round((bytes / k ** i) * 100) / 100} ${sizes[i]}`;
  };

  return (
    <Card className="border-2">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <div className="rounded-md bg-primary/10 p-2">
            <FileUp className="h-5 w-5 text-primary" />
          </div>
          Upload File
        </CardTitle>
        <CardDescription>
          Upload documents (PDF, Excel, CSV, Text) to use as knowledge for your
          chatbot
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {selectedFile ? (
          <div className="space-y-4">
            <Card className="border-2 border-dashed">
              <CardContent className="p-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="rounded-md bg-chart-2/10 p-2">
                      <FileUp className="h-6 w-6 text-chart-2" />
                    </div>
                    <div>
                      <p className="font-medium text-sm">{selectedFile.name}</p>
                      <p className="text-muted-foreground text-xs">
                        {formatFileSize(selectedFile.size)}
                      </p>
                    </div>
                  </div>
                  {!mutation.isPending && (
                    <Button
                      className="h-8 w-8 p-0"
                      onClick={handleRemove}
                      size="sm"
                      variant="ghost"
                    >
                      <X className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              </CardContent>
            </Card>

            {mutation.isPending && (
              <div className="space-y-2">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted-foreground">Uploading...</span>
                  <span className="font-medium">{uploadProgress}%</span>
                </div>
                <Progress className="h-2" value={uploadProgress} />
              </div>
            )}
          </div>
        ) : (
          <Dropzone
            accept={SUPPORTED_FILE_MIME_TYPES.reduce(
              // biome-ignore lint/performance/noAccumulatingSpread: ikr
              (acc, type) => ({ ...acc, [type]: [] }),
              {}
            )}
            disabled={mutation.isPending}
            maxSize={MAX_FILE_SIZE_BYTES}
            onDrop={handleFileSelect}
          >
            <DropzoneEmptyState />
            <DropzoneContent />
          </Dropzone>
        )}

        {mutation.isError && (
          <Alert variant="destructive">
            <AlertDescription>
              {mutation.error?.message || "Failed to upload file"}
            </AlertDescription>
          </Alert>
        )}

        <Alert>
          <AlertDescription className="text-xs">
            <strong>Supported formats:</strong> PDF, Excel (.xlsx, .xls), CSV,
            Text (.txt, .md).
            <br />
            <strong>Maximum file size:</strong> 100MB.
          </AlertDescription>
        </Alert>
      </CardContent>
      <CardFooter className="flex justify-between border-t pt-4">
        <Button
          disabled={mutation.isPending}
          onClick={onCancel}
          type="button"
          variant="outline"
        >
          Cancel
        </Button>
        <Button
          disabled={!selectedFile || mutation.isPending}
          onClick={handleUpload}
        >
          {mutation.isPending && (
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
          )}
          Upload File
        </Button>
      </CardFooter>
    </Card>
  );
}
