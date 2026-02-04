import React from "react";

interface LauncherProps {
  isOpen: boolean;
  onClick: () => void;
  themeColor: string;
  position: "bottom-right" | "bottom-left" | "top-right" | "top-left";
}

export const Launcher: React.FC<LauncherProps> = ({
  isOpen,
  onClick,
  themeColor,
  position,
}) => {
  const getPositionClasses = () => {
    switch (position) {
      case "bottom-right":
        return "bottom-5 right-5";
      case "bottom-left":
        return "bottom-5 left-5";
      case "top-right":
        return "top-5 right-5";
      case "top-left":
        return "top-5 left-5";
      default:
        return "bottom-5 right-5";
    }
  };

  return (
    <button
      type="button"
      onClick={onClick}
      className={`fixed ${getPositionClasses()} z-9999 w-14 h-14 rounded-full shadow-lg transition-all duration-300 hover:scale-110 focus:outline-none focus:ring-4 focus:ring-opacity-50 animate-fade-in`}
      style={{
        backgroundColor: themeColor,
      }}
      aria-label={isOpen ? "Close chat" : "Open chat"}
    >
      {isOpen ? (
        <svg
          className="w-6 h-6 mx-auto text-white"
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
      ) : (
        <svg
          className="w-6 h-6 mx-auto text-white"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <title>Open</title>
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"
          />
        </svg>
      )}
    </button>
  );
};
