# chunky

## Reliable, resumable uploads for big files — because losing progress at 99% fucking sucks.

We've all been there. You're waiting for a massive file to upload, staring at the progress bar inching forward like a 3AM microwave. And just when it's about to finish, it fails. Maybe your internet cut out. Maybe your laptop died. Or maybe you fat-fingered `Ctrl-C` and terminated the upload. Whatever it was, you’re back to square one.

That’s where `chunky` comes in.

**`chunky` is a Go library for resumable uploads of arbitrarily large files. Now, if a upload fails halfway through, it picks up where it left off so you don’t have to start over.**

---

## Features

- Resumable uploads, even after crashes or disconnects
- File locking and hashing for integrity
- Efficient streaming/chunking for large files
- Upload progress persistence via local storage
- Retry logic with exponential backoff
- Duplicate-safe: avoids conflicts across users uploading identical files

---

## Protocol Overview

### 1. Initialization

The client is initialized with:
- A **server URL**
- A **local file path**

The file is locked for reading to prevent accidental modification during upload.  
If it's large, it’s read as a stream.

A **SHA-256 hash** is calculated to uniquely identify the file version.

---

### 2. Initiate Upload Session

The client sends a `POST /initiateUploadSession` request with:
- File hash
- Total bytes

The server:
- Generates a unique session ID
- Associates this session ID with the file hash, byte size, and a TTL (default: 5 days)
- Starts a job to automatically clean up expired sessions
- Returns the session ID and TTL

The client:
- Stores the session metadata (session ID ↔ file) locally (e.g., via SQLite)

---

### 3. Uploading Chunks

The file is divided into chunks and uploaded using a `PATCH /uploadChunk` request containing:
- **Session ID**
- **File hash**
- **Chunk start & end byte**
- **Chunk data**

> 🔐 **Why both session ID *and* file hash?**  
> - File hash ensures versioning and detects file changes  
> - Session ID avoids user conflicts (e.g., two people uploading the same file)

---

### 4. Server Handling of Chunks

The server:
- Verifies session ID (else returns `Invalid session ID`)
- Validates file hash (else returns `Invalid file hash`)
- Checks if the chunk is already stored (if yes, responds with 2xx)
- Stores new chunks and responds with `200 OK`

---

### 5. Chunk Response Handling

The client:
- Proceeds on any `2xx` response
- On network error: retries with exponential backoff (max 5 attempts), then pauses
- On `invalid session id` or `invalid file hash`: terminates upload.  
  The user must restart from **step 2** to continue.

---

### 6. Resuming Uploads

If the upload is paused or interrupted:
- The client starts again from **step 1**, but calls `GET /uploadStatus`
- The server responds with **missing byte ranges**

Note: The server may report chunks as missing if they haven't been stored yet (even if just received). Reuploading overlapping chunks is safe and expected.

---

### 7. Finalizing Upload

Once all chunks are uploaded:
- The client sends another `GET /uploadStatus`
- If the server reports status as `complete`, the client sends `POST /endUpload`
- The server deletes all metadata associated with the session

If the session is already gone, it returns a `session doesn't exist` error (safe to ignore if the upload is truly complete).

---

## Edge Cases Handled

- **Out-of-Order Chunk Uploads**: server tracks uploaded chunks
- **Double uploads**: okay, safely ignored
- **Paused uploads**: resume without re-uploading everything
- **Modified file in between upload pause and resume**: caught via file hash mismatch
- **Concurrent users uploading same file**: handled via per-user session IDs
- **Automatic session cleanup**: via TTL and background cleanup job



## Edge Cases To Handle
- **Upload Session TTL Expiration Mid-Upload**
- **Storage Space Full on Server**
- **Multiple Clients Trying to Resume the Same Session**

---

## Contributing

Want to improve `chunky`? PRs welcome.  
Found a bug or edge case? Open an issue and describe what broke!

---

## License

MIT © 2025  
Because uploads shouldn't be painful.

