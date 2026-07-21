package server

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/content"
)

// --- sitemap.xml -----------------------------------------------------------

type urlEntry struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float64 `xml:"priority,omitempty"`
}

type urlSet struct {
	XMLName xml.Name   `xml:"urlset"`
	Xmlns   string     `xml:"xmlns,attr"`
	URLs    []urlEntry `xml:"url"`
}

func (s *Server) handleSitemap(w http.ResponseWriter, _ *http.Request) {
	today := time.Now().UTC().Format("2006-01-02")
	set := urlSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}

	for _, p := range content.StaticPages {
		set.URLs = append(set.URLs, urlEntry{
			Loc:        s.abs(p.Path),
			LastMod:    today,
			ChangeFreq: p.Changefreq,
			Priority:   p.Priority,
		})
	}

	if posts, err := content.LoadBlogPosts(s.cfg.ContentDir); err == nil {
		for _, post := range posts {
			last := today
			if !post.Updated.IsZero() {
				last = post.Updated.Format("2006-01-02")
			}
			set.URLs = append(set.URLs, urlEntry{
				Loc:        s.abs("/blog/" + post.Slug),
				LastMod:    last,
				ChangeFreq: "yearly",
				Priority:   0.6,
			})
		}
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_, _ = w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	_ = enc.Encode(set)
}

// --- robots.txt ------------------------------------------------------------

// aiCrawlers are explicitly welcomed so generative engines can cite the site.
var aiCrawlers = []string{"GPTBot", "OAI-SearchBot", "ChatGPT-User", "PerplexityBot", "ClaudeBot", "Claude-Web", "Google-Extended", "CCBot", "Applebot-Extended"}

func (s *Server) handleRobots(w http.ResponseWriter, _ *http.Request) {
	var b strings.Builder
	b.WriteString("# Search + AI crawlers are welcome.\n")
	b.WriteString("User-agent: *\n")
	b.WriteString("Allow: /\n\n")
	for _, ua := range aiCrawlers {
		fmt.Fprintf(&b, "User-agent: %s\nAllow: /\n\n", ua)
	}
	b.WriteString("Disallow: /api/\n\n")
	fmt.Fprintf(&b, "Sitemap: %s\n", s.abs("/sitemap.xml"))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = w.Write([]byte(b.String()))
}

// --- llms.txt (GEO) --------------------------------------------------------

// handleLLMs implements the emerging /llms.txt convention: a concise, machine-
// friendly map of the site for generative engines.
func (s *Server) handleLLMs(w http.ResponseWriter, _ *http.Request) {
	var b strings.Builder
	b.WriteString("# Alexandre Iwanesko — Coach d'échecs à Genève (Maître FIDE)\n\n")
	b.WriteString("> Alexandre Iwanesko est Maître FIDE et coach d'échecs à Genève. Il enseigne aux ")
	b.WriteString("adultes (1200–2200 Elo) et aux adolescents en compétition, en présentiel à ")
	b.WriteString("Genève et en ligne dans toute la Suisse romande. Spécialités : progression ")
	b.WriteString("structurée, préparation de tournoi et coaching individualisé. Public visé : ")
	b.WriteString("joueurs intermédiaires à avancés (pas de cours débutant).\n\n")

	b.WriteString("## Faits clés\n")
	b.WriteString("- Titre : Maître FIDE (FIDE Master)\n")
	b.WriteString("- Zone : Genève, Vaud, arc lémanique, France voisine\n")
	b.WriteString("- Formats : cours particulier, petit groupe, en ligne, stages, entreprise\n")
	b.WriteString("- Langue : français (anglais possible)\n\n")

	b.WriteString("## Pages\n")
	for _, p := range content.StaticPages {
		fmt.Fprintf(&b, "- [%s](%s): %s\n", p.Title, s.abs(p.Path), p.Summary)
	}

	if posts, err := content.LoadBlogPosts(s.cfg.ContentDir); err == nil && len(posts) > 0 {
		b.WriteString("\n## Articles\n")
		for _, post := range posts {
			fmt.Fprintf(&b, "- [%s](%s): %s\n", post.Title, s.abs("/blog/"+post.Slug), post.Description)
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_, _ = w.Write([]byte(b.String()))
}

// abs builds an absolute URL from a root-relative path.
func (s *Server) abs(p string) string {
	if p == "/" {
		return s.cfg.BaseURL + "/"
	}
	return s.cfg.BaseURL + p
}
