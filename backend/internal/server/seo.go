package server

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/iwanesko/chess-web-site/backend/internal/content"
)

// --- sitemap.xml -----------------------------------------------------------

// urlEntry carries a loc and, when it is true, a lastmod. Nothing else:
// changefreq and priority are hints Google has said for years it ignores, and
// they were two more values to keep honest for no return.
type urlEntry struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type urlSet struct {
	XMLName xml.Name   `xml:"urlset"`
	Xmlns   string     `xml:"xmlns,attr"`
	URLs    []urlEntry `xml:"url"`
}

// handleSitemap lists every indexable URL with the date it really changed.
//
// It used to stamp time.Now() on all of them, recomputed at each request: the
// whole file looked rewritten every time a crawler read it. That is the one
// way to lose the lastmod signal entirely — Google drops a lastmod it catches
// lying, and the articles, where the date IS true, lose it along with the rest.
// Now a static page states its own date (content.Page.Updated), a blog index or
// a category archive inherits the date of its newest article, and an article
// carries its `updated` front matter, or failing that its publication date.
func (s *Server) handleSitemap(w http.ResponseWriter, _ *http.Request) {
	enDir := s.cfg.ContentDir + "/en"
	set := urlSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}

	// Category archives with no article are excluded — see EmptyCategoryPaths.
	empty := content.EmptyCategoryPaths(s.cfg.ContentDir, enDir)
	fromBlog := content.BlogLastMod(s.cfg.ContentDir, enDir)

	// Les deux index de tactiques datent de leur semaine la plus récente : c'est
	// le jour où la page a changé, et il tombe tout seul chaque lundi.
	weeks := content.TacticsWeeks(s.cfg.TacticsDir)
	if len(weeks) > 0 {
		fromBlog["/tactiques"] = weeks[0].Date
		fromBlog["/en/tactics"] = weeks[0].Date
	}

	for _, p := range content.StaticPages {
		if empty[p.Path] {
			continue
		}
		last := p.Updated
		if d, ok := fromBlog[p.Path]; ok {
			last = d
		}
		set.URLs = append(set.URLs, urlEntry{Loc: s.abs(p.Path), LastMod: last})
	}

	addPosts := func(dir, prefix string) {
		posts, err := content.LoadBlogPosts(dir)
		if err != nil {
			return
		}
		for _, post := range posts {
			last := ""
			if !post.Updated.IsZero() {
				last = post.Updated.Format("2006-01-02")
			}
			set.URLs = append(set.URLs, urlEntry{Loc: s.abs(prefix + post.Slug), LastMod: last})
		}
	}
	// Une page par semaine, dans les deux langues. Elles existent depuis le
	// début, pré-rendues et liées depuis l'index, mais étaient en noindex : des
	// positions tirées de vraies parties, que personne d'autre ne publie, et que
	// rien ne déclarait. Le noindex est levé (chaque série porte désormais un
	// texte calculé sur ses propres positions), elles entrent donc au sitemap.
	for _, wk := range weeks {
		set.URLs = append(set.URLs,
			urlEntry{Loc: s.abs("/tactiques/" + wk.Slug), LastMod: wk.Date},
			urlEntry{Loc: s.abs("/en/tactics/" + wk.Slug), LastMod: wk.Date},
		)
	}

	addPosts(s.cfg.ContentDir, "/blog/")          // FR
	addPosts(s.cfg.ContentDir+"/en", "/en/blog/") // EN

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

// --- security.txt (RFC 9116) -----------------------------------------------

func (s *Server) handleSecurityTxt(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	fmt.Fprintf(w, "Contact: mailto:%s\nExpires: 2027-12-31T23:59:59.000Z\nPreferred-Languages: fr, en\nCanonical: %s\n",
		s.cfg.MailTo, s.abs("/.well-known/security.txt"))
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
	// Une source externe vérifiable vaut mieux qu'une affirmation : classements et
	// tournois joués s'y contrôlent, et la fiche reste à jour toute seule.
	b.WriteString("- Joueur en activité — fiche FIDE officielle (ID 682136) : https://ratings.fide.com/profile/682136\n")
	b.WriteString("- Zone : Genève, Vaud, arc lémanique, France voisine\n")
	b.WriteString("- Formats : cours particulier, petit groupe, en ligne, stages, entreprise\n")
	b.WriteString("- Langue : français (anglais possible)\n")
	b.WriteString("- Le blog publie ses propres carnets de tournoi, chiffres et parties à l'appui.\n\n")

	// La convention llms.txt prévoit une section « Optional » : ce qu'un moteur peut
	// sauter s'il manque de contexte. Pages légales et listes de catégories y vont,
	// pour que le haut du fichier ne contienne que l'essentiel.
	var fr, en, optional []content.Page
	for _, p := range content.StaticPages {
		switch {
		case isSecondary(p.Path):
			optional = append(optional, p)
		case strings.HasPrefix(p.Path, "/en"):
			en = append(en, p)
		default:
			fr = append(fr, p)
		}
	}

	writePages := func(heading string, pages []content.Page) {
		if len(pages) == 0 {
			return
		}
		fmt.Fprintf(&b, "## %s\n", heading)
		for _, p := range pages {
			fmt.Fprintf(&b, "- [%s](%s): %s\n", p.Title, s.abs(p.Path), p.Summary)
		}
		b.WriteString("\n")
	}
	writePages("Pages principales (français)", fr)
	writePages("Main pages (English)", en)

	// Les articles des DEUX langues — le pendant de ce que fait déjà le sitemap.
	writePosts := func(heading, dir, prefix string) {
		posts, err := content.LoadBlogPosts(dir)
		if err != nil || len(posts) == 0 {
			return
		}
		fmt.Fprintf(&b, "## %s\n", heading)
		for _, post := range posts {
			fmt.Fprintf(&b, "- [%s](%s): %s\n", post.Title, s.abs(prefix+post.Slug), post.Description)
		}
		b.WriteString("\n")
	}
	writePosts("Articles (français)", s.cfg.ContentDir, "/blog/")
	writePosts("Articles (English)", s.cfg.ContentDir+"/en", "/en/blog/")

	writePages("Optional", optional)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_, _ = w.Write([]byte(b.String()))
}

// isSecondary marque les pages qu'un moteur peut ignorer sans rien perdre de
// l'essentiel : mentions légales et pages de catégorie, qui ne font que lister
// des articles déjà présents plus haut.
func isSecondary(path string) bool {
	for _, frag := range []string{"/confidentialite", "/en/privacy", "/blog/categorie/", "/en/blog/category/"} {
		if strings.Contains(path, frag) {
			return true
		}
	}
	return false
}

// abs builds an absolute URL from a root-relative path.
func (s *Server) abs(p string) string {
	if p == "/" {
		return s.cfg.BaseURL + "/"
	}
	return s.cfg.BaseURL + p
}
