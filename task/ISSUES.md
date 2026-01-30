# Issue: Docker Build Fails with ResourceExhausted

## Description
The Docker build process fails during the UI build stage with a `ResourceExhausted` error. The specific error message is:
`failed to solve: ResourceExhausted: process "/bin/sh -c bun run build && rm -rf dist && mv .output/public dist" did not complete successfully: cannot allocate memory`

This occurs because the `bun run build` command (which runs `vite build`) consumes more memory than is available in the Docker container environment. The `Dockerfile` currently sets `NODE_OPTIONS="--max-old-space-size=4096"`, which requests 4GB of memory for Node.js. If the Docker VM/container has less than ~4.5GB of available RAM, this can lead to an OOM (Out Of Memory) kill.

## Steps to Fix
1.  **Analyze current memory usage**: Review `Dockerfile` settings (done). Found `ENV NODE_OPTIONS="--max-old-space-size=4096"`.
2.  **Adjust Memory Limit**: Lower the `max-old-space-size` to a safer value (e.g., 2048MB) to leave room for the OS and other processes within the container, or advise increasing the Docker VM memory if 4GB is strictly required (unlikely for this project size).
3.  **Optimize Build Command**: Ensure the build command is efficient.
4.  **Verify Fix**: Run `make docker-build` to confirm the build completes successfully.

## Recommended Fix
Update `Dockerfile` line 6:
```dockerfile
ENV NODE_OPTIONS="--max-old-space-size=2048"
```
Or potentially remove the env var if not strictly needed, allowing Node/Bun to manage memory naturally, but limiting it explicitly prevents it from trying to take more than available. Given the crash, explicit lower limit or increasing VM resources is key.

Start with lowering to 2048MB.
