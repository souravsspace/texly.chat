import { useEffect, useState } from "react";

export type Theme = "dark" | "light" | "system";

/**
 * Hook to detect and manage dark mode within the widget
 * This works independently of the host page's theme
 */
export function useDarkMode() {
  const [theme, setTheme] = useState<Theme>(() => {
    // Try to get saved preference from localStorage
    try {
      const saved = localStorage.getItem("texly-widget-theme");
      if (saved === "dark" || saved === "light" || saved === "system") {
        return saved as Theme;
      }
    } catch (e) {
      // localStorage might not be available
      console.warn("[useDarkMode] Failed to access localStorage:", e);
    }
    return "system";
  });

  const [isDark, setIsDark] = useState(() => {
    // Initialize with system preference
    try {
      if (typeof window !== 'undefined' && window.matchMedia) {
        return window.matchMedia("(prefers-color-scheme: dark)").matches;
      }
    } catch (e) {
      console.warn("[useDarkMode] Failed to check system preference:", e);
    }
    return false;
  });

  useEffect(() => {
    console.log("[useDarkMode] Effect running, theme:", theme);
    
    const determineTheme = () => {
      try {
        if (theme === "system") {
          // Use system preference
          const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
          setIsDark(mediaQuery.matches);
        } else {
          setIsDark(theme === "dark");
        }
      } catch (e) {
        console.error("[useDarkMode] Error determining theme:", e);
      }
    };

    determineTheme();

    // Listen for system theme changes if using system preference
    if (theme === "system") {
      try {
        const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
        const handler = (e: MediaQueryListEvent) => setIsDark(e.matches);
        
        // Modern browsers
        if (mediaQuery.addEventListener) {
          mediaQuery.addEventListener("change", handler);
          return () => mediaQuery.removeEventListener("change", handler);
        } else {
          // Legacy browsers
          mediaQuery.addListener(handler);
          return () => mediaQuery.removeListener(handler);
        }
      } catch (e) {
        console.error("[useDarkMode] Error setting up media query listener:", e);
      }
    }
  }, [theme]);

  const setThemeWithStorage = (newTheme: Theme) => {
    try {
      localStorage.setItem("texly-widget-theme", newTheme);
    } catch (e) {
      console.warn("[useDarkMode] Failed to save theme preference:", e);
    }
    setTheme(newTheme);
  };

  return { theme, setTheme: setThemeWithStorage, isDark };
}
