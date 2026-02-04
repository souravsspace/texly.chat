import type { Session } from "../types";

const SESSION_KEY_PREFIX = "texly-session-";

export class SessionManager {
  private botId: string;

  constructor(botId: string) {
    this.botId = botId;
  }

  private getKey(): string {
    return `${SESSION_KEY_PREFIX}${this.botId}`;
  }

  /**
   * Get stored session if it exists and hasn't expired
   */
  getSession(): Session | null {
    try {
      const stored = localStorage.getItem(this.getKey());
      if (!stored) return null;

      const session: Session = JSON.parse(stored);

      // Check if expired
      const expiresAt = new Date(session.expiresAt);
      if (expiresAt < new Date()) {
        this.clearSession();
        return null;
      }

      return session;
    } catch (error) {
      console.error("Failed to get session:", error);
      return null;
    }
  }

  /**
   * Store session in localStorage
   */
  setSession(session: Session): void {
    try {
      localStorage.setItem(this.getKey(), JSON.stringify(session));
    } catch (error) {
      console.error("Failed to save session:", error);
    }
  }

  /**
   * Clear session from localStorage
   */
  clearSession(): void {
    try {
      localStorage.removeItem(this.getKey());
    } catch (error) {
      console.error("Failed to clear session:", error);
    }
  }
}
