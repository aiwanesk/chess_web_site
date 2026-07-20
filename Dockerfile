# syntax=docker/dockerfile:1

# ---------- Stage 1: build the frontend (React → static HTML via SSG) ----------
FROM node:22-alpine AS frontend
WORKDIR /app/frontend
# Install deps first (better layer caching).
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
# Content lives at the repo root; copy what the build reads.
COPY frontend/ ./
COPY content/ /app/content/
COPY public/ /app/public/
RUN npm run build
# → produces /app/frontend/dist

# ---------- Stage 2: build the Go binary with the frontend embedded ----------
FROM golang:1.26-alpine AS backend
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
# Replace the committed placeholder with the real build so //go:embed picks it up.
RUN rm -rf ./dist
COPY --from=frontend /app/frontend/dist ./dist
# Static, stripped binary.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /server .

# ---------- Stage 3: minimal runtime ----------
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=backend /server /app/server
# Blog Markdown is read at runtime for sitemap.xml / llms.txt generation.
COPY --from=frontend /app/content /app/content
# Persistent volume for runtime state (SQLite: stats + published puzzles). On
# Jelastic, mount /data as a persistent volume so it survives redeploys.
VOLUME /data
ENV ADDR=":8080" \
    CONTENT_DIR="/app/content/blog" \
    TACTICS_DIR="/app/content/tactiques" \
    DB_PATH="/data/tactics.db"
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
