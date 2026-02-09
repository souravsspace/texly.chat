import { useForm } from "@tanstack/react-form";
import { Globe, MapPin, MessageSquare, Palette } from "lucide-react";
import type { WidgetConfig } from "@/api/index.types";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Field, FieldGroup } from "@/components/ui/fields";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";

interface WidgetConfigFormProps {
  initialConfig: WidgetConfig;
  initialOrigins: string[];
  onSubmit: (config: WidgetConfig, allowedOrigins: string[]) => void;
  onConfigChange?: (config: WidgetConfig) => void;
  isSubmitting: boolean;
}

export function WidgetConfigForm({
  initialConfig,
  initialOrigins,
  onSubmit,
  onConfigChange,
  isSubmitting,
}: WidgetConfigFormProps) {
  const form = useForm({
    defaultValues: {
      themeColor: initialConfig.theme_color,
      initialMessage: initialConfig.initial_message,
      position: initialConfig.position,
      botAvatar: initialConfig.bot_avatar,
      allowedOrigins: initialOrigins.join("\n"),
    },
    onSubmit: ({ value }) => {
      const widgetConfig: WidgetConfig = {
        theme_color: value.themeColor,
        initial_message: value.initialMessage,
        position: value.position,
        bot_avatar: value.botAvatar,
      };

      const originsArray = value.allowedOrigins
        .split("\n")
        .map((origin) => origin.trim())
        .filter((origin) => origin.length > 0);

      onSubmit(widgetConfig, originsArray);
    },
  });

  // Helper to update preview when any field changes
  const updatePreview = (values: {
    themeColor: string;
    initialMessage: string;
    position: string;
    botAvatar: string;
  }) => {
    if (onConfigChange) {
      const widgetConfig: WidgetConfig = {
        theme_color: values.themeColor,
        initial_message: values.initialMessage,
        position: values.position,
        bot_avatar: values.botAvatar,
      };
      onConfigChange(widgetConfig);
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Widget Settings</CardTitle>
      </CardHeader>
      <CardContent>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
        >
          <FieldGroup className="space-y-6">
            {/* Theme Color */}
            <form.Field name="themeColor">
              {(field) => (
                <Field>
                  <Label htmlFor={field.name}>
                    <Palette className="mr-2 inline h-4 w-4" />
                    Theme Color
                  </Label>
                  <div className="flex gap-3">
                    <Input
                      className="h-10 w-20 cursor-pointer"
                      id={field.name}
                      onChange={(e) => {
                        field.handleChange(e.target.value);
                        updatePreview({
                          themeColor: e.target.value,
                          initialMessage: form.state.values.initialMessage,
                          position: form.state.values.position,
                          botAvatar: form.state.values.botAvatar,
                        });
                      }}
                      type="color"
                      value={field.state.value}
                    />
                    <Input
                      className="flex-1"
                      onChange={(e) => {
                        field.handleChange(e.target.value);
                        updatePreview({
                          themeColor: e.target.value,
                          initialMessage: form.state.values.initialMessage,
                          position: form.state.values.position,
                          botAvatar: form.state.values.botAvatar,
                        });
                      }}
                      placeholder="#6366f1"
                      type="text"
                      value={field.state.value}
                    />
                  </div>
                  <p className="mt-1 text-muted-foreground text-sm">
                    Primary color for the widget interface
                  </p>
                </Field>
              )}
            </form.Field>

            {/* Initial Message */}
            <form.Field name="initialMessage">
              {(field) => (
                <Field>
                  <Label htmlFor={field.name}>
                    <MessageSquare className="mr-2 inline h-4 w-4" />
                    Welcome Message
                  </Label>
                  <Textarea
                    id={field.name}
                    onChange={(e) => {
                      field.handleChange(e.target.value);
                      updatePreview({
                        themeColor: form.state.values.themeColor,
                        initialMessage: e.target.value,
                        position: form.state.values.position,
                        botAvatar: form.state.values.botAvatar,
                      });
                    }}
                    placeholder="Hi! How can I help you today?"
                    rows={3}
                    value={field.state.value}
                  />
                  <p className="mt-1 text-muted-foreground text-sm">
                    First message shown when the widget opens
                  </p>
                </Field>
              )}
            </form.Field>

            {/* Position */}
            <form.Field name="position">
              {(field) => (
                <Field>
                  <Label htmlFor={field.name}>
                    <MapPin className="mr-2 inline h-4 w-4" />
                    Widget Position
                  </Label>
                  <Select
                    onValueChange={(value) => {
                      field.handleChange(value);
                      updatePreview({
                        themeColor: form.state.values.themeColor,
                        initialMessage: form.state.values.initialMessage,
                        position: value,
                        botAvatar: form.state.values.botAvatar,
                      });
                    }}
                    value={field.state.value}
                  >
                    <SelectTrigger id={field.name}>
                      <SelectValue placeholder="Select position" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="bottom-right">Bottom Right</SelectItem>
                      <SelectItem value="bottom-left">Bottom Left</SelectItem>
                    </SelectContent>
                  </Select>
                  <p className="mt-1 text-muted-foreground text-sm">
                    Where the widget appears on the page
                  </p>
                </Field>
              )}
            </form.Field>

            {/* Bot Avatar URL */}
            <form.Field name="botAvatar">
              {(field) => (
                <Field>
                  <Label htmlFor={field.name}>Bot Avatar URL (Optional)</Label>
                  <Input
                    id={field.name}
                    onChange={(e) => {
                      field.handleChange(e.target.value);
                      updatePreview({
                        themeColor: form.state.values.themeColor,
                        initialMessage: form.state.values.initialMessage,
                        position: form.state.values.position,
                        botAvatar: e.target.value,
                      });
                    }}
                    placeholder="https://example.com/avatar.png"
                    type="url"
                    value={field.state.value}
                  />
                  <p className="mt-1 text-muted-foreground text-sm">
                    URL to an image for the bot's avatar
                  </p>
                </Field>
              )}
            </form.Field>

            {/* Allowed Origins */}
            <form.Field name="allowedOrigins">
              {(field) => (
                <Field>
                  <Label htmlFor={field.name}>
                    <Globe className="mr-2 inline h-4 w-4" />
                    Allowed Origins
                  </Label>
                  <Textarea
                    className="font-mono text-sm"
                    id={field.name}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="https://example.com&#10;https://www.example.com&#10;http://localhost:3000"
                    rows={5}
                    value={field.state.value}
                  />
                  <p className="mt-1 text-muted-foreground text-sm">
                    One domain per line. Leave empty to allow all origins (not
                    recommended for production)
                  </p>
                </Field>
              )}
            </form.Field>

            {/* Submit Button */}
            <Button className="w-full" disabled={isSubmitting} type="submit">
              {isSubmitting ? "Saving..." : "Save Configuration"}
            </Button>
          </FieldGroup>
        </form>
      </CardContent>
    </Card>
  );
}
