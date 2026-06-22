# syntax=docker/dockerfile:1

# --- Stage 1: build the React SPA into frontend/dist ---
FROM node:22-slim AS frontend
RUN corepack enable
# CI=true makes pnpm non-interactive: it skips the pre-run deps-status check that
# otherwise tries to purge node_modules and aborts on "no TTY" inside the build.
ENV CI=true
WORKDIR /app/frontend
# Install deps first for better layer caching.
COPY frontend/package.json frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm build

# --- Stage 2: build the Go server with the SPA embedded ---
FROM golang:1.26.4 AS backend
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# The embed.go pattern requires frontend/dist to exist at compile time.
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

# --- Stage 3: minimal runtime image ---
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=backend /out/server /server
EXPOSE 8090
ENTRYPOINT ["/server"]
