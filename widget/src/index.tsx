import React from "react";
import { createRoot } from "react-dom/client";
import { App } from "./App";
import "./styles/index.css";

// Self-executing function to initialize widget
(() => {
  // Find the script tag with data-bot-id
  const scriptTag = document.currentScript as HTMLScriptElement;
  const botId = scriptTag?.getAttribute("data-bot-id");

  if (!botId) {
    console.error("Texly Widget: data-bot-id attribute is required");
    return;
  }

  // Determine CSS URL based on script location
  const getCssUrl = () => {
    // In development, use local CSS path or hardcoded local server
    if (process.env.NODE_ENV === "development") {
      return "http://localhost:8080/widget/texly-widget.css";
    }

    const scriptTag = document.querySelector(
      'script[src*="texly-widget.js"]',
    ) as HTMLScriptElement;

    if (scriptTag && scriptTag.src) {
      // Replace .js with .css
      return scriptTag.src.replace(".js", ".css");
    }

    return null;
  };

  // Wait for DOM to be ready
  const init = () => {
    // Create container
    const container = document.createElement("div");
    container.id = `texly-widget-${botId}`;
    document.body.appendChild(container);

    // Attach Shadow DOM for style isolation
    const shadowRoot = container.attachShadow({ mode: "open" });

    // Inject CSS
    const cssUrl = getCssUrl();
    if (cssUrl) {
      const link = document.createElement("link");
      link.rel = "stylesheet";
      link.href = cssUrl;
      shadowRoot.appendChild(link);
    }

    // Create React root container
    const reactRoot = document.createElement("div");
    reactRoot.id = "react-root";
    shadowRoot.appendChild(reactRoot);

    // Render React app
    const root = createRoot(reactRoot);
    root.render(
      <React.StrictMode>
        <App botId={botId} />
      </React.StrictMode>,
    );
  };

  // Initialize when DOM is ready
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();
