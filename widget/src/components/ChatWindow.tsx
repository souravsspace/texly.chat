import React, { useEffect, useRef, useState } from "react";
import type { Message } from "../types";

interface ChatWindowProps {
  isOpen: boolean;
  messages: Message[];
  botName: string;
  botAvatar?: string;
  themeColor: string;
  position: "bottom-right" | "bottom-left" | "top-right" | "top-left";
  isLoading: boolean;
  onSendMessage: (message: string) => void;
  onClose: () => void;
}

export const ChatWindow: React.FC<ChatWindowProps> = ({
  isOpen,
  messages,
  botName,
  botAvatar,
  themeColor,
  position,
  isLoading,
  onSendMessage,
  onClose,
}) => {
  const [input, setInput] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (input.trim() && !isLoading) {
      onSendMessage(input.trim());
      setInput("");
    }
  };

  const getPositionClasses = () => {
    switch (position) {
      case "bottom-right":
        return "bottom-24 right-5";
      case "bottom-left":
        return "bottom-24 left-5";
      case "top-right":
        return "top-24 right-5";
      case "top-left":
        return "top-24 left-5";
      default:
        return "bottom-24 right-5";
    }
  };

  if (!isOpen) return null;

  return (
    <div
      className={`fixed ${getPositionClasses()} z-9998 w-96 h-[600px] bg-white rounded-lg shadow-2xl flex flex-col animate-slide-up`}
      style={{ maxHeight: "calc(100vh - 120px)" }}
    >
      {/* Header */}
      <div
        className="flex items-center justify-between p-4 rounded-t-lg text-white"
        style={{ backgroundColor: themeColor }}
      >
        <div className="flex items-center gap-3">
          {botAvatar ? (
            <img
              src={botAvatar}
              alt={botName}
              className="w-10 h-10 rounded-full"
            />
          ) : (
            <div className="w-10 h-10 rounded-full bg-white/20 flex items-center justify-center">
              <span className="text-lg font-semibold">
                {botName.charAt(0).toUpperCase()}
              </span>
            </div>
          )}
          <div>
            <h3 className="font-semibold text-sm">{botName}</h3>
            <p className="text-xs opacity-90">Online</p>
          </div>
        </div>
        <button
          onClick={onClose}
          className="p-1 hover:bg-white/20 rounded transition-colors"
          aria-label="Close chat"
          type="button"
        >
          <svg
            className="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <title>Close</title>
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4 widget-scrollbar">
        {messages.map((msg, idx) => (
          <div
            key={idx}
            className={`flex ${msg.role === "user" ? "justify-end" : "justify-start"}`}
          >
            <div
              className={`max-w-[80%] rounded-lg px-4 py-2 ${
                msg.role === "user"
                  ? "bg-gray-100 text-gray-900"
                  : "text-white"
              }`}
              style={
                msg.role === "assistant"
                  ? { backgroundColor: themeColor }
                  : undefined
              }
            >
              <p className="text-sm whitespace-pre-wrap">{msg.content}</p>
              <p className="text-xs opacity-70 mt-1">
                {new Date(msg.timestamp).toLocaleTimeString([], {
                  hour: "2-digit",
                  minute: "2-digit",
                })}
              </p>
            </div>
          </div>
        ))}
        {isLoading && (
          <div className="flex justify-start">
            <div
              className="rounded-lg px-4 py-2 text-white"
              style={{ backgroundColor: themeColor }}
            >
              <div className="flex gap-1">
                <div
                  className="w-2 h-2 bg-white rounded-full animate-bounce"
                  style={{ animationDelay: "0ms" }}
                />
                <div
                  className="w-2 h-2 bg-white rounded-full animate-bounce"
                  style={{ animationDelay: "150ms" }}
                />
                <div
                  className="w-2 h-2 bg-white rounded-full animate-bounce"
                  style={{ animationDelay: "300ms" }}
                />
              </div>
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <form onSubmit={handleSubmit} className="p-4 border-t">
        <div className="flex gap-2">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Type a message..."
            disabled={isLoading}
            className="flex-1 px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 disabled:opacity-50 disabled:cursor-not-allowed text-sm"
          />
          <button
            type="submit"
            disabled={isLoading || !input.trim()}
            className="px-4 py-2 rounded-lg text-white font-medium transition-all hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed"
            style={{ backgroundColor: themeColor }}
          >
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <title>Send</title>
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"
              />
            </svg>
          </button>
        </div>
      </form>
    </div>
  );
};
