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
	// TacticsDir holds the weekly anonymised puzzle files (AAAA-Sxx.json) served
	// at /api/tactics. Empty disables the endpoint.
	TacticsDir string

	// SMTP settings for sending contact-form emails. If SMTPHost/SMTPUser are
	// empty, email sending is skipped (the lead is still logged/forwarded).
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	MailFrom string // sender address (defaults to SMTPUser)
	MailTo   string // recipient (defaults to alexandre@iwanesko.ch)

	// DBPath is the SQLite file for private stats. Empty disables stats.
	// The newsletter re-uses the same file (separate tables).
	DBPath string
	// AdminToken protects the /admin dashboard + stats. Empty disables /admin.
	AdminToken string
	// EventsFile is an optional JSON file of stages/events used for newsletter
	// announcements (see internal/server/announce.go).
	EventsFile string
}

// LoadConfig reads configuration from the environment, applying safe defaults.
func LoadConfig() Config {
	return Config{
		Addr:              env("ADDR", ":8080"),
		BaseURL:           strings.TrimRight(env("BASE_URL", "https://iwanesko.ch"), "/"),
		DevStaticDir:      os.Getenv("DEV_STATIC_DIR"),
		ContentDir:        env("CONTENT_DIR", "../content/blog"),
		ContactWebhookURL: os.Getenv("CONTACT_WEBHOOK_URL"),
		TacticsDir:        env("TACTICS_DIR", "../content/tactiques"),
		SMTPHost:          os.Getenv("SMTP_HOST"),
		SMTPPort:          env("SMTP_PORT", "587"),
		SMTPUser:          os.Getenv("SMTP_USER"),
		SMTPPass:          os.Getenv("SMTP_PASS"),
		MailFrom:          env("MAIL_FROM", os.Getenv("SMTP_USER")),
		MailTo:            env("MAIL_TO", "alexandre@iwanesko.ch"),
		DBPath:            os.Getenv("DB_PATH"),
		AdminToken:        os.Getenv("ADMIN_TOKEN"),
		EventsFile:        env("EVENTS_FILE", "../content/events.json"),
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
