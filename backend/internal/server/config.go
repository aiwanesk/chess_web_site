package server

import (
	"os"
	"strings"
)

// Config holds the runtime configuration, all sourced from environment
// variables so no secret is ever hard-coded. See .env.example.
type Config struct {
	Addr string // listen address, e.g. ":8080"
	// BaseURL is the public origin (no trailing slash), used to build
	// absolute URLs in sitemap.xml, canonical tags and llms.txt.
	BaseURL string
	// DevStaticDir, when set, serves the frontend from the filesystem
	// (e.g. ../frontend/dist) instead of the binary-embedded copy. Handy in
	// development so you don't have to rebuild the Go binary on every front change.
	DevStaticDir string
	// ContentDir is where blog Markdown lives, read for sitemap generation.
	ContentDir string
	// ContactWebhookURL, when set, forwards contact submissions (e.g. to an
	// email relay or CRM). Empty means submissions are only logged.
	ContactWebhookURL string
}

// LoadConfig reads configuration from the environment, applying safe defaults.
func LoadConfig() Config {
	return Config{
		Addr:              env("ADDR", ":8080"),
		BaseURL:           strings.TrimRight(env("BASE_URL", "https://iwanesko.ch"), "/"),
		DevStaticDir:      os.Getenv("DEV_STATIC_DIR"),
		ContentDir:        env("CONTENT_DIR", "../content/blog"),
		ContactWebhookURL: os.Getenv("CONTACT_WEBHOOK_URL"),
	}
}

// StaticMode reports whether static assets are served from the embedded FS
// ("embedded") or from disk ("dev:<dir>").
func (c Config) StaticMode() string {
	if c.DevStaticDir != "" {
		return "dev:" + c.DevStaticDir
	}
	return "embedded"
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
