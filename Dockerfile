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

# ---------- Stage 3: runtime ----------
# Debian base (NOT scratch/distroless): the Virtuozzo/Jelastic platform injects
# its integration packages (curl for health-checks, iptables-persistent for the
# network layer) via apt on every deploy. A shell-less, apt-less image makes that
# step fail → the whole deployment errors out. A Debian base has apt + a shell,
# so deploys succeed and the Web SSH works. The Go binary runs identically.
FROM debian:bookworm-slim
WORKDIR /app
# ca-certificates → TLS for outbound SMTP (Infomaniak) and HTTPS.
# curl + iptables-persistent → expected by the Virtuozzo integration layer.
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        ca-certificates curl iptables-persistent \
    && rm -rf /var/lib/apt/lists/*
COPY --from=backend /server /app/server
# Blog Markdown is read at runtime for sitemap.xml / llms.txt generation.
COPY --from=frontend /app/content /app/content
# Create /data so SQLite can write even before a persistent volume is mounted
# over it (prevents "unable to open database file" on first boot).
RUN mkdir -p /data
# On Jelastic, mount /data as a persistent volume so it survives redeploys.
VOLUME /data
ENV ADDR=":8080" \
    CONTENT_DIR="/app/content/blog" \
    TACTICS_DIR="/app/content/tactiques" \
    DB_PATH="/data/tactics.db"
EXPOSE 8080
# Runs as root: the platform's deploy-time apt injection needs it (and it lets
# SQLite write to a freshly-mounted /data without uid-permission friction).
ENTRYPOINT ["/app/server"]
