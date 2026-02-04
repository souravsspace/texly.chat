export interface WidgetConfig {
  themeColor: string;
  initialMessage: string;
  position: "bottom-right" | "bottom-left" | "top-right" | "top-left";
  botAvatar?: string;
}

export interface BotConfig {
  id: string;
  name: string;
  widget_config: WidgetConfig;
}

export interface Message {
  role: "user" | "assistant";
  content: string;
  timestamp: Date;
}

export interface Session {
  sessionId: string;
  botId: string;
  expiresAt: string;
}

export interface ChatTokenResponse {
  type: "token" | "done" | "error";
  content?: string;
  error?: string;
}

export interface CreateSessionRequest {
  bot_id: string;
}

export interface ChatRequest {
  message: string;
}
