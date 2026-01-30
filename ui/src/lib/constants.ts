// Source type constants
export const SOURCE_TYPE = {
  URL: "url",
  FILE: "file",
  TEXT: "text",
} as const;

export type SourceTypeValue = (typeof SOURCE_TYPE)[keyof typeof SOURCE_TYPE];

// Source status constants
export const SOURCE_STATUS = {
  PENDING: "pending",
  PROCESSING: "processing",
  COMPLETED: "completed",
  FAILED: "failed",
} as const;

export type SourceStatusValue =
  (typeof SOURCE_STATUS)[keyof typeof SOURCE_STATUS];

// Supported file types
export const SUPPORTED_FILE_TYPES = [
  ".txt",
  ".md",
  ".pdf",
  ".xlsx",
  ".xls",
  ".csv",
] as const;

export const SUPPORTED_FILE_MIME_TYPES = [
  "text/plain",
  "text/markdown",
  "application/pdf",
  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
  "application/vnd.ms-excel",
  "text/csv",
] as const;

// File upload limits
export const MAX_FILE_SIZE_MB = 100;
export const MAX_FILE_SIZE_BYTES = MAX_FILE_SIZE_MB * 1024 * 1024;

// Text source limits
export const MAX_TEXT_SIZE_MB = 10;
export const MAX_TEXT_SIZE_BYTES = MAX_TEXT_SIZE_MB * 1024 * 1024;
