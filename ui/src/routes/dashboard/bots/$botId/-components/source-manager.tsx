import { FileText, FileUp, Link as LinkIcon, List } from "lucide-react";
import { useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { FileUploadForm } from "./file-upload-form";
import { SourceList } from "./source-list";
import { TextSourceForm } from "./text-source-form";
import { UrlSourceForm } from "./url-source-form";

export function SourceManager({ botId }: { botId: string }) {
  const [activeTab, setActiveTab] = useState("list");

  return (
    <Card>
      <CardContent className="p-6">
        <Tabs className="w-full" onValueChange={setActiveTab} value={activeTab}>
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger className="gap-2" value="list">
              <List className="h-4 w-4" />
              Sources
            </TabsTrigger>
            <TabsTrigger className="gap-2" value="url">
              <LinkIcon className="h-4 w-4" />
              Add URL
            </TabsTrigger>
            <TabsTrigger className="gap-2" value="file">
              <FileUp className="h-4 w-4" />
              Upload File
            </TabsTrigger>
            <TabsTrigger className="gap-2" value="text">
              <FileText className="h-4 w-4" />
              Add Text
            </TabsTrigger>
          </TabsList>

          <TabsContent className="mt-6 space-y-4" value="list">
            <SourceList botId={botId} onAddSource={() => setActiveTab("url")} />
          </TabsContent>

          <TabsContent className="mt-6 space-y-4" value="url">
            <UrlSourceForm
              botId={botId}
              onCancel={() => setActiveTab("list")}
              onSuccess={() => setActiveTab("list")}
            />
          </TabsContent>

          <TabsContent className="mt-6 space-y-4" value="file">
            <FileUploadForm
              botId={botId}
              onCancel={() => setActiveTab("list")}
              onSuccess={() => setActiveTab("list")}
            />
          </TabsContent>

          <TabsContent className="mt-6 space-y-4" value="text">
            <TextSourceForm
              botId={botId}
              onCancel={() => setActiveTab("list")}
              onSuccess={() => setActiveTab("list")}
            />
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  );
}
