import React from "react";
import { createRoot } from "react-dom/client";
import { App } from "./App";
import "./styles/index.css";

// Self-executing function to initialize widget
(() => {
  console.log("[Texly Widget] Initializing...");
  
  // Find the script tag with data-bot-id
  const scriptTag = document.currentScript as HTMLScriptElement;
  const botId = scriptTag?.getAttribute("data-bot-id");

  console.log("[Texly Widget] Bot ID:", botId);
  console.log("[Texly Widget] Script tag:", scriptTag);

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
    console.log("[Texly Widget] DOM ready, initializing...");
    
    // Create container
    const container = document.createElement("div");
    container.id = `texly-widget-${botId}`;
    document.body.appendChild(container);

    console.log("[Texly Widget] Container appended to body:", container);

    // Attach Shadow DOM for style isolation
    const shadowRoot = container.attachShadow({ mode: "open" });
    console.log("[Texly Widget] Shadow DOM created:", shadowRoot);

    // Inject CSS
    const cssUrl = getCssUrl();
    console.log("[Texly Widget] CSS URL:", cssUrl);
    
    if (cssUrl) {
      const link = document.createElement("link");
      link.rel = "stylesheet";
      link.href = cssUrl;
      link.onload = () => console.log("[Texly Widget] CSS loaded successfully");
      link.onerror = () => console.error("[Texly Widget] CSS failed to load");
      shadowRoot.appendChild(link);
    } else {
      console.warn("[Texly Widget] No CSS URL found");
    }

    // Create React root container
    const reactRoot = document.createElement("div");
    reactRoot.id = "react-root";
    shadowRoot.appendChild(reactRoot);

    console.log("[Texly Widget] React root container created:", reactRoot);

    // Render React app
    const root = createRoot(reactRoot);
    console.log("[Texly Widget] Rendering React app...");
    
    root.render(
      <React.StrictMode>
        <App botId={botId} />
      </React.StrictMode>,
    );
    
    console.log("[Texly Widget] React app rendered");
  };

  // Initialize when DOM is ready
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();
