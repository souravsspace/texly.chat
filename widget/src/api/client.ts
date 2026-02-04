import type {
  BotConfig,
  ChatRequest,
  ChatTokenResponse,
  CreateSessionRequest,
  Session,
} from "../types";

// Determine API base URL from the script tag
const getApiBase = () => {
  // In development, use the current origin or localhost:8080
  if (process.env.NODE_ENV === "development") {
    return "http://localhost:8080";
  }

  // Find the script tag that loaded the widget
  const scriptTag = document.querySelector(
    'script[src*="texly-widget.js"]',
  ) as HTMLScriptElement;

  if (scriptTag && scriptTag.src) {
    try {
      const url = new URL(scriptTag.src);
      // Remove /widget/texly-widget.js to get the base
      return url.origin;
    } catch (e) {
      console.error("Failed to parse script URL", e);
    }
  }

  // Fallback to current origin (for same-domain embedding)
  return window.location.origin;
};

const API_BASE = getApiBase();

export class WidgetAPI {
  private botId: string;

  constructor(botId: string) {
    this.botId = botId;
  }

  /**
   * Fetch widget configuration for the bot
   */
  async getConfig(): Promise<BotConfig> {
    const response = await fetch(
      `${API_BASE}/api/public/bots/${this.botId}/config`,
    );

    if (!response.ok) {
      throw new Error("Failed to fetch widget configuration");
    }

    return response.json();
  }

  /**
   * Create a new anonymous session
   */
  async createSession(): Promise<Session> {
    const body: CreateSessionRequest = {
      bot_id: this.botId,
    };

    const response = await fetch(`${API_BASE}/api/public/chats`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    });

    if (!response.ok) {
      throw new Error("Failed to create session");
    }

    const data = await response.json();
    return {
      sessionId: data.session_id,
      botId: data.bot_id,
      expiresAt: data.expires_at,
    };
  }

  /**
   * Send a message and stream the response
   */
  async streamMessage(
    sessionId: string,
    message: string,
    onToken: (token: string) => void,
    onComplete: () => void,
    onError: (error: string) => void,
  ): Promise<void> {
    const body: ChatRequest = { message };

    const response = await fetch(
      `${API_BASE}/api/public/chats/${sessionId}/messages`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
      },
    );

    if (!response.ok) {
      if (response.status === 404 || response.status === 401) {
        throw new Error("SESSION_EXPIRED");
      }
      throw new Error(`Failed to send message: ${response.statusText}`);
    }

    const reader = response.body?.getReader();
    if (!reader) {
      throw new Error("No response body");
    }

    const decoder = new TextDecoder();
    let buffer = "";

    try {
      while (true) {
        const { done, value } = await reader.read();

        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split("\n");
        buffer = lines.pop() || "";

        for (const line of lines) {
          if (line.startsWith("data: ")) {
            const data = line.substring(6);
            try {
              const parsed: ChatTokenResponse = JSON.parse(data);

              if (parsed.type === "token" && parsed.content) {
                onToken(parsed.content);
              } else if (parsed.type === "done") {
                onComplete();
                return;
              } else if (parsed.type === "error" && parsed.error) {
                onError(parsed.error);
                return;
              }
            } catch (e) {
              console.error("Failed to parse SSE message:", e);
            }
          }
        }
      }
    } catch (error) {
      onError(error instanceof Error ? error.message : "Unknown error");
    }
  }
}
