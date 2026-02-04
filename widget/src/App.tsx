import React, { useEffect, useState } from "react";
import { WidgetAPI } from "./api/client";
import { ChatWindow } from "./components/ChatWindow";
import { Launcher } from "./components/Launcher";
import type { BotConfig, Message, Session } from "./types";
import { SessionManager } from "./utils/session";

interface AppProps {
  botId: string;
}

export const App: React.FC<AppProps> = ({ botId }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [config, setConfig] = useState<BotConfig | null>(null);
  const [session, setSession] = useState<Session | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const api = new WidgetAPI(botId);
  const sessionManager = new SessionManager(botId);

  // Load configuration on mount
  useEffect(() => {
    const loadConfig = async () => {
      try {
        const botConfig = await api.getConfig();
        setConfig(botConfig);

        // Add initial message
        if (botConfig.widget_config.initialMessage) {
          setMessages([
            {
              role: "assistant",
              content: botConfig.widget_config.initialMessage,
              timestamp: new Date(),
            },
          ]);
        }
      } catch (err) {
        console.error("Failed to load widget config:", err);
        setError("Failed to load chat widget");
      }
    };

    loadConfig();
  }, [botId]);

  // Load or create session
  useEffect(() => {
    const initSession = async () => {
      // Try to get existing session
      let existingSession = sessionManager.getSession();

      if (existingSession) {
        setSession(existingSession);
        return;
      }

      // Create new session
      try {
        const newSession = await api.createSession();
        sessionManager.setSession(newSession);
        setSession(newSession);
      } catch (err) {
        console.error("Failed to create session:", err);
        setError("Failed to start chat session");
      }
    };

    if (config && !session) {
      initSession();
    }
  }, [config]);

  const handleSendMessage = async (content: string) => {
    if (!session) {
      setError("No active session");
      return;
    }

    // Add user message
    // Add user message immediately
    const userMessage: Message = {
      role: "user",
      content,
      timestamp: new Date(),
    };

    // Create placeholder assistant message
    const assistantMessage: Message = {
      role: "assistant",
      content: "",
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage, assistantMessage]);
    setIsLoading(true);

    let assistantContent = "";

    try {
      await api.streamMessage(
        session.sessionId,
        content,
        // onToken
        (token) => {
          assistantContent += token;
          setMessages((prev) => {
            const newMessages = [...prev];
            const lastIndex = newMessages.length - 1;
            
            // Should always be updating the last message (the placeholder)
            if (lastIndex >= 0) {
              newMessages[lastIndex] = {
                ...newMessages[lastIndex],
                content: assistantContent,
              };
            }

            return newMessages;
          });
        },
        // onComplete
        () => {
          setIsLoading(false);
        },
        // onError
        (err) => {
          console.error("Chat error:", err);
          if (err === "SESSION_EXPIRED") {
             sessionManager.clearSession();
             setSession(null);
             setError("Session expired. Please refresh.");
          } else {
             setError("Failed to get response");
          }
          setIsLoading(false);
        },
      );
    } catch (err) {
      console.error("Failed to send message:", err);
       if (err instanceof Error && err.message === "SESSION_EXPIRED") {
         sessionManager.clearSession();
         setSession(null);
         setError("Session expired. Please refresh.");
       } else {
          setError("Failed to send message");
       }
      setIsLoading(false);
    }
  };

  if (error) {
    return (
      <div className="fixed bottom-5 right-5 z-[9999] bg-red-500 text-white px-4 py-2 rounded-lg shadow-lg">
        {error}
      </div>
    );
  }

  if (!config) {
    return null; // Loading...
  }

  const { widget_config } = config;

  return (
    <>
      <Launcher
        isOpen={isOpen}
        onClick={() => setIsOpen(!isOpen)}
        themeColor={widget_config.themeColor}
        position={widget_config.position}
      />
      <ChatWindow
        isOpen={isOpen}
        messages={messages}
        botName={config.name}
        botAvatar={widget_config.botAvatar}
        themeColor={widget_config.themeColor}
        position={widget_config.position}
        isLoading={isLoading}
        onSendMessage={handleSendMessage}
        onClose={() => setIsOpen(false)}
      />
    </>
  );
};
