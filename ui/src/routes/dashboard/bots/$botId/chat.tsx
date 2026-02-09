"use client";

import { createFileRoute } from "@tanstack/react-router";
import { ArrowUpIcon, BotIcon, Loader2Icon } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { api } from "@/api";
import {
  Conversation,
  ConversationContent,
  ConversationEmptyState,
  ConversationScrollButton,
} from "@/components/ai-elements/conversation";
import {
  Message,
  MessageContent,
  MessageResponse,
} from "@/components/ai-elements/message";
import { Suggestion, Suggestions } from "@/components/ai-elements/suggestion";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";

export const Route = createFileRoute("/dashboard/bots/$botId/chat")({
  component: ChatPage,
});

interface ChatMessage {
  id: string;
  role: "user" | "assistant";
  content: string;
}

const SUGGESTIONS = [
  "What information do you have?",
  "Summarize the main topics",
  "What are the key points?",
  "Explain in detail",
];

function ChatPage() {
  const { botId } = Route.useParams();
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputValue, setInputValue] = useState("");
  const [isStreaming, setIsStreaming] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputValue.trim() || isStreaming) return;

    const userMessage: ChatMessage = {
      id: `user-${Date.now()}`,
      role: "user",
      content: inputValue.trim(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInputValue("");
    setIsStreaming(true);

    const assistantMessageId = `assistant-${Date.now()}`;
    let accumulatedContent = "";

    try {
      // Add empty assistant message that will be updated
      setMessages((prev) => [
        ...prev,
        {
          id: assistantMessageId,
          role: "assistant",
          content: "",
        },
      ]);

      // Stream tokens from the API
      for await (const token of api.chat.stream(botId, userMessage.content)) {
        if (token.type === "token" && token.content) {
          accumulatedContent += token.content;

          // Update the assistant message with accumulated content
          setMessages((prev) =>
            prev.map((msg) =>
              msg.id === assistantMessageId
                ? { ...msg, content: accumulatedContent }
                : msg
            )
          );
        } else if (token.type === "error") {
          throw new Error(token.error || "An error occurred during streaming");
        }
      }
    } catch (error) {
      console.error("Chat error:", error);
      toast.error(
        error instanceof Error ? error.message : "Failed to send message"
      );

      // Remove the failed assistant message
      setMessages((prev) =>
        prev.filter((msg) => msg.id !== assistantMessageId)
      );
    } finally {
      setIsStreaming(false);
    }
  };

  const handleSuggestionClick = (suggestion: string) => {
    setInputValue(suggestion);
    textareaRef.current?.focus();
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  // Auto-resize textarea
  // biome-ignore lint/correctness/useExhaustiveDependencies: We want to resize when inputValue changes
  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = "auto";
      textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`;
    }
  }, [inputValue]);

  return (
    <Card className="flex flex-col overflow-hidden">
      {/* Messages */}
      <Conversation className="min-h-[500px] flex-1">
        <ConversationContent>
          {messages.length === 0 ? (
            <ConversationEmptyState
              description="Ask me anything about your knowledge base"
              icon={<BotIcon className="size-8" />}
              title="Start a conversation"
            >
              <div className="mt-6 w-full max-w-2xl">
                <Suggestions>
                  {SUGGESTIONS.map((suggestion) => (
                    <Suggestion
                      key={suggestion}
                      onClick={handleSuggestionClick}
                      suggestion={suggestion}
                    />
                  ))}
                </Suggestions>
              </div>
            </ConversationEmptyState>
          ) : (
            <>
              {messages.map((message) => (
                <Message from={message.role} key={message.id}>
                  <MessageContent>
                    {message.role === "assistant" ? (
                      <MessageResponse>{message.content}</MessageResponse>
                    ) : (
                      message.content
                    )}
                  </MessageContent>
                </Message>
              ))}
              {isStreaming && (
                <div className="flex items-center gap-2 text-muted-foreground">
                  <Loader2Icon className="size-4 animate-spin" />
                  <span className="text-sm">Thinking...</span>
                </div>
              )}
            </>
          )}
        </ConversationContent>
        <ConversationScrollButton />
      </Conversation>

      {/* Input */}
      <div className="shrink-0 border-t p-4">
        <form className="mx-auto max-w-3xl" onSubmit={handleSubmit}>
          <div className="relative flex items-end gap-2">
            <Textarea
              className="max-h-[200px] min-h-[60px] resize-none pr-12"
              disabled={isStreaming}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Type your message..."
              ref={textareaRef}
              rows={1}
              value={inputValue}
            />
            <Button
              className="absolute right-2 bottom-2"
              disabled={!inputValue.trim() || isStreaming}
              size="icon"
              type="submit"
            >
              {isStreaming ? (
                <Loader2Icon className="size-4 animate-spin" />
              ) : (
                <ArrowUpIcon className="size-4" />
              )}
            </Button>
          </div>
        </form>
      </div>
    </Card>
  );
}
