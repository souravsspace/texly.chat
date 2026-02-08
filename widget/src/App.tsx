import React, { useEffect, useMemo, useState } from "react";
import { WidgetAPI } from "./api/client";
import { ChatWindow } from "./components/ChatWindow";
import { Launcher } from "./components/Launcher";
import type { BotConfig, Message, Session } from "./types";
import { SessionManager } from "./utils/session";
import { useDarkMode } from "./hooks/useDarkMode";

interface AppProps {
	botId: string;
}

export const App: React.FC<AppProps> = ({ botId }) => {
	console.log("[Texly Widget App] Rendering with botId:", botId);

	const { isDark } = useDarkMode();
	console.log("[Texly Widget App] Dark mode status:", isDark);
	
	const [isOpen, setIsOpen] = useState(false);
	const [config, setConfig] = useState<BotConfig | null>(null);
	const [session, setSession] = useState<Session | null>(null);
	const [messages, setMessages] = useState<Message[]>([]);
	const [isLoading, setIsLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const api = useMemo(() => new WidgetAPI(botId), [botId]);
	const sessionManager = useMemo(() => new SessionManager(botId), [botId]);

	// Load configuration on mount
	useEffect(() => {
		console.log("[Texly Widget App] Loading bot config...");

		const loadConfig = async () => {
			try {
				const botConfig = await api.getConfig();
				console.log("[Texly Widget App] Bot config loaded:", botConfig);
				setConfig(botConfig);

				// Add initial message
				if (botConfig.widget_config.initial_message) {
					setMessages([
						{
							role: "assistant",
							content: botConfig.widget_config.initial_message,
							timestamp: new Date(),
						},
					]);
				}
			} catch (err) {
				console.error("[Texly Widget App] Failed to load widget config:", err);
				setError("Failed to load chat widget");
			}
		};

		loadConfig();
	}, [api]);

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
	}, [config, session, api, sessionManager]);

	const handleSendMessage = async (content: string) => {
		console.log("[Widget App] handleSendMessage called with:", content);
		console.log("[Widget App] Current session:", session);
		console.log("[Widget App] Current messages count:", messages.length);

		if (!session) {
			setError("No active session");
			return;
		}

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

		console.log("[Widget App] Adding messages to state");
		setMessages((prev) => {
			const updated = [...prev, userMessage, assistantMessage];
			console.log("[Widget App] Messages updated, new count:", updated.length);
			return updated;
		});
		setIsLoading(true);

		// Use a ref to track accumulated content to avoid closure issues
		const contentRef = { current: "" };

		try {
			await api.streamMessage(
				session.sessionId,
				content,
				// onToken
				(token) => {
					console.log("[Widget App] Token received in callback:", token);
					contentRef.current += token;
					console.log("[Widget App] Accumulated content:", contentRef.current);

					// Update the last message with the accumulated content
					setMessages((prev) => {
						const newMessages = [...prev];
						const lastIndex = newMessages.length - 1;

						console.log(
							"[Widget App] Updating message at index:",
							lastIndex,
							"Role:",
							newMessages[lastIndex]?.role,
						);

						if (lastIndex >= 0 && newMessages[lastIndex].role === "assistant") {
							newMessages[lastIndex] = {
								...newMessages[lastIndex],
								content: contentRef.current,
							};
							console.log(
								"[Widget App] Updated assistant message content to:",
								contentRef.current,
							);
						}

						return newMessages;
					});
				},
				// onComplete
				() => {
					console.log(
						"[Widget App] Stream completed, final content length:",
						contentRef.current.length,
					);
					setIsLoading(false);
				},
				// onError
				(err) => {
					console.error("[Widget App] Chat error:", err);
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
			console.error("[Widget App] Failed to send message:", err);
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
		console.error("[Texly Widget App] Error state:", error);
		return (
			<div className="fixed bottom-5 right-5 z-[9999] bg-red-500 text-white px-4 py-2 rounded-lg shadow-lg">
				{error}
			</div>
		);
	}

	if (!config) {
		console.log("[Texly Widget App] Waiting for config...");
		return null; // Loading...
	}

	console.log("[Texly Widget App] Rendering launcher and chat window");

	const { widget_config } = config;

	return (
		<>
			<Launcher
				isOpen={isOpen}
				onClick={() => setIsOpen(!isOpen)}
				themeColor={widget_config.theme_color}
				position={widget_config.position}
			/>
			<ChatWindow
				isOpen={isOpen}
				messages={messages}
				botName={config.name}
				botAvatar={widget_config.bot_avatar}
				themeColor={widget_config.theme_color}
				position={widget_config.position}
				isLoading={isLoading}
				onSendMessage={handleSendMessage}
				onClose={() => setIsOpen(false)}
				isDark={isDark}
			/>
		</>
	);
};
