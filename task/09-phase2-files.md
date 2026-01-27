# Phase 2.2: File Upload Training

## Goal
Allow users to upload PDF, DOCX, and TXT files to train their chatbots.

## Backend Tasks

### Step 1: File Parser Service
**File**: `internal/services/parser/parser.go`
- Libs: `ledongthuc/pdf`, `docconv` (for docx).
- Method: `ParseFile(path, mimeType) -> (text, error)`.

### Step 2: Upload Handler
**File**: `internal/handlers/source/file.go`
- POST `/api/bots/:id/sources/upload`
- Multipart Form Data (`file`).
- Limit: 10MB.
- Logic:
    1. Save file to disk (`./uploads` or S3).
    2. Create `Source` record (Type: "file").
    3. Enqueue `ProcessFileTask`.

## Frontend Tasks

### Step 1: File Upload Component
**File**: `ui/src/components/file-upload.tsx`
- Drag & Drop zone.
- Progress bar during upload.
