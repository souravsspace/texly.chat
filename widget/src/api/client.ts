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
    console.log("[Widget API] Sending message:", message, "to session:", sessionId);
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

    console.log("[Widget API] Response status:", response.status);

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

    console.log("[Widget API] Starting to read SSE stream...");

    const decoder = new TextDecoder();
    let buffer = "";
    let completedSuccessfully = false;

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
            console.log("[Widget API] Received SSE data:", data);
            try {
              const parsed: ChatTokenResponse = JSON.parse(data);

              if (parsed.type === "token" && parsed.content) {
                console.log("[Widget API] Token received:", parsed.content);
                onToken(parsed.content);
              } else if (parsed.type === "done") {
                console.log("[Widget API] Stream completed successfully");
                completedSuccessfully = true;
                onComplete();
                return;
              } else if (parsed.type === "error" && parsed.error) {
                console.error("[Widget API] Stream error:", parsed.error);
                onError(parsed.error);
                return;
              }
            } catch (e) {
              console.error("Failed to parse SSE message:", e);
            }
          }
        }
      }

      // If we exit the loop without receiving a "done" message, still call onComplete
      if (!completedSuccessfully) {
        console.log("[Widget API] Stream ended without 'done' message, completing anyway");
        onComplete();
      }
    } catch (error) {
      onError(error instanceof Error ? error.message : "Unknown error");
    }
  }
}
