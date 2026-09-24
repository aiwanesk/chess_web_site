// Package content is the single source of truth for the site's indexable
// routes and blog posts. The Go backend uses it to generate sitemap.xml and
// llms.txt so those stay in sync with the pages shipped by the frontend.
package content

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Page is one indexable route of the site.
//
// Updated is the date the page's copy last changed (YYYY-MM-DD), written by
// hand on purpose. The sitemap used to send time.Now() for every static page,
// so all forty URLs claimed to have changed on every single crawl. A lastmod
// that moves without the content moving is a lastmod Google learns to ignore —
// and the signal is then lost for the articles too, where it is true. Bump the
// line when you rewrite a page: leaving it stale is harmless, a self-updating
// lie is not. Pages built from data (the blog index, the category archives)
// carry no date here — LastMod computes theirs from the articles.
//
// No changefreq, no priority: Google reads neither. They were noise in the
// file and one more column to keep in sync.
type Page struct {
	Path    string // e.g. "/cours-echecs-adultes-geneve"
	Title   string // human label, reused in llms.txt
	Updated string // YYYY-MM-DD, empty for data-driven pages (see LastMod)
	Summary string // one-line description for llms.txt / GEO
}

// StaticPages lists every non-blog indexable route. Keep this aligned with the
// frontend router in frontend/src/routes.tsx.
var StaticPages = []Page{
	{"/", "Accueil — Alexandre Iwanesko, Maître FIDE & coach d'échecs à Genève", "2026-09-24",
		"Coach d'échecs à Genève, Maître FIDE. Cours pour adultes 1200–2200 Elo et ados en compétition, préparation tournoi, en présentiel et en ligne."},
	{"/cours-echecs-adultes-geneve", "Cours d'échecs pour adultes à Genève", "2026-09-24",
		"Cours d'échecs pour adultes (1200–2200 Elo) à Genève avec un Maître FIDE : méthode structurée, plan de progression, présentiel et en ligne."},
	{"/preparation-tournoi-echecs", "Préparation tournoi d'échecs", "2026-09-24",
		"Préparation ciblée aux tournois d'échecs : ouvertures, gestion du temps, préparation adverse et mental de compétition."},
	{"/cours-echecs-en-ligne", "Cours d'échecs en ligne", "2026-09-24",
		"Cours d'échecs particuliers en ligne avec un Maître FIDE, depuis toute la Suisse romande et la France voisine."},
	{"/cours-echecs-groupe-geneve", "Cours d'échecs en groupe à Genève", "2026-09-24",
		"Cours d'échecs en petit groupe à Genève : émulation, tarif réduit, niveau homogène."},
	{"/cours-echecs-ados-competition", "Cours d'échecs pour ados en compétition", "2026-09-24",
		"Coaching pour adolescents joueurs de compétition : progression Elo, préparation tournoi et suivi individualisé."},
	{"/stages-echecs-geneve", "Stages d'échecs à Genève", "2026-09-24",
		"Stages d'échecs intensifs à Genève pendant les vacances scolaires, encadrés par un Maître FIDE."},
	{"/conferences-echecs-entreprise", "Conférences d'échecs en entreprise", "2026-09-24",
		"Conférences et interventions échecs en entreprise : stratégie, prise de décision et gestion du risque."},
	{"/team-building-echecs-geneve", "Team building échecs à Genève", "2026-09-24",
		"Ateliers de team building autour des échecs pour entreprises à Genève et dans l'arc lémanique."},
	{"/a-propos", "À propos d'Alexandre Iwanesko, Maître FIDE", "2026-09-24",
		"Parcours d'Alexandre Iwanesko, Maître FIDE et coach d'échecs à Genève : titre, résultats et méthode d'enseignement."},
	{"/resultats", "Résultats & témoignages", "2026-09-24",
		"Résultats des élèves et témoignages : progressions Elo, performances en tournoi et retours d'expérience."},
	{"/tarifs", "Tarifs des cours d'échecs", "2026-09-24",
		"Tarifs des cours d'échecs à Genève : cours particuliers, en groupe, en ligne et forfaits de préparation tournoi."},
	{"/contact", "Contact", "2026-09-24",
		"Contacter Alexandre Iwanesko pour un cours d'échecs à Genève ou en ligne : premier échange pour définir vos objectifs."},
	{"/reserver", "Réserver un cours d'échecs", "2026-09-24",
		"Réserver un cours d'échecs particulier avec Alexandre Iwanesko, Maître FIDE, en soirée (17h30-20h00). Confirmation immédiate par e-mail."},
	{"/confidentialite", "Politique de confidentialité", "2026-09-24",
		"Traitement des données personnelles sur iwanesko.ch (contact et newsletter), conforme au RGPD et à la LPD suisse."},
	{"/calendrier", "Calendrier des tournois", "2026-09-24",
		"Le calendrier de compétition d'Alexandre Iwanesko, Maître FIDE : tournois à venir et carnets des opens déjà joués."},
	{"/blog", "Blog échecs", "",
		"Articles d'échecs : stratégie, ouvertures, préparation tournoi et progression pour joueurs intermédiaires et avancés."},
	{"/tactiques", "Tactiques de la semaine", "2026-09-24",
		"Résolvez les plus belles tactiques d'échecs de la semaine, sélectionnées par un Maître FIDE. Nouveaux puzzles chaque semaine."},
	{"/blog/categorie/progresser", "Progresser aux échecs — guides", "",
		"Guides d'un Maître FIDE pour progresser aux échecs : ouvertures, finales, tactique, préparation de tournoi et mental."},
	{"/blog/categorie/carnet-de-tournoi", "Carnet de tournoi", "",
		"Le journal de compétition d'Alexandre Iwanesko, Maître FIDE : parties marquantes, décisions sous pression et leçons de tournoi."},

	// English (EN) — grows as the i18n rollout continues.
	{"/en", "Chess coach in Geneva — Alexandre Iwanesko, FIDE Master", "2026-09-24",
		"Chess coach in Geneva, FIDE Master. Lessons for adults (1200–2200 Elo) and competitive teens, tournament prep, in person and online."},
	{"/en/adult-chess-lessons-geneva", "Adult chess lessons in Geneva", "2026-09-24",
		"Chess lessons for adults (1200–2200 Elo) in Geneva with a FIDE Master: structured method, progression plan, in person and online."},
	{"/en/tournament-preparation", "Chess tournament preparation", "2026-09-24",
		"Targeted chess tournament preparation with a FIDE Master: repertoire, opponent prep, time management and competitive mindset."},
	{"/en/online-chess-lessons", "Online chess lessons", "2026-09-24",
		"Private online chess lessons with a FIDE Master, from anywhere in French-speaking Switzerland and neighbouring France."},
	{"/en/group-chess-lessons-geneva", "Group chess lessons in Geneva", "2026-09-24",
		"Small-group chess lessons in Geneva: matched level, healthy competition and a reduced per-person rate."},
	{"/en/junior-chess-coaching", "Junior chess coaching", "2026-09-24",
		"Coaching for competitive teenagers: Elo progression, tournament preparation and individual follow-up."},
	{"/en/chess-camps-geneva", "Chess camps in Geneva", "2026-09-24",
		"Intensive chess camps in Geneva during the school holidays, led by a FIDE Master."},
	{"/en/corporate-chess-talks", "Corporate chess talks", "2026-09-24",
		"Chess talks and keynotes for companies: strategy, decision-making and risk management."},
	{"/en/chess-team-building-geneva", "Chess team building in Geneva", "2026-09-24",
		"Chess team-building workshops for companies in Geneva and the Lake Geneva region."},
	{"/en/about", "About Alexandre Iwanesko, FIDE Master", "2026-09-24",
		"The path of Alexandre Iwanesko, FIDE Master and chess coach in Geneva: title, results and teaching method."},
	{"/en/results", "Results & testimonials", "2026-09-24",
		"Student results and testimonials: Elo progress, tournament performances and feedback."},
	{"/en/tactics", "Tactics of the week", "2026-09-24",
		"Solve the best chess tactics of the week, hand-picked by a FIDE Master. New puzzles every week."},
	{"/en/blog/category/improve", "Improve at chess — guides", "",
		"Guides from a FIDE Master to improve at chess: openings, endgames, tactics, tournament preparation and mindset."},
	{"/en/blog/category/tournament-diary", "Tournament diary", "",
		"The competition diary of Alexandre Iwanesko, FIDE Master: key games, decisions under pressure and lessons from the road."},
	{"/en/calendar", "Tournament calendar", "2026-09-24",
		"The competition calendar of Alexandre Iwanesko, FIDE Master: upcoming tournaments and diaries of the opens already played."},
	{"/en/blog", "Chess blog", "",
		"Guides from a FIDE Master to improve at chess (openings, endgames, tournament prep) and Alexandre Iwanesko's tournament diary."},
	{"/en/pricing", "Chess lesson pricing", "2026-09-24",
		"Chess lesson pricing in Geneva: private lessons, group lessons, online and packages."},
	{"/en/contact", "Contact", "2026-09-24",
		"Contact Alexandre Iwanesko, FIDE Master, for chess lessons in Geneva or online."},
	{"/en/book", "Book a chess lesson", "2026-09-24",
		"Book a private chess lesson with Alexandre Iwanesko, FIDE Master, in the evening (17:30-20:00). Instant e-mail confirmation."},
	{"/en/privacy", "Privacy policy", "2026-09-24",
		"How personal data is handled on iwanesko.ch (contact and newsletter), GDPR and Swiss FADP compliant."},
}

// BlogPost is the metadata parsed from a Markdown file's front matter.
type BlogPost struct {
	Slug        string
	Title       string
	Description string
	Category    string
	Date        time.Time
	Updated     time.Time
}

// DefaultCategory mirrors DEFAULT_CATEGORY in frontend/src/lib/categories.ts.
const DefaultCategory = "progresser"

// blogCategory mirrors one entry of frontend/src/lib/categories.ts. Only the
// stable key and the two URL slugs are needed here — the sitemap has no use for
// the labels.
type blogCategory struct{ Key, SlugFR, SlugEN string }

var blogCategories = []blogCategory{
	{"progresser", "progresser", "improve"},
	{"carnet-de-tournoi", "carnet-de-tournoi", "tournament-diary"},
}

// EmptyCategoryPaths returns the archive URLs that currently hold no article,
// in both locales. They must stay out of the sitemap: an archive that promises
// guides and lists none is a thin page. The rule is generic — publish an article
// in a category and its URL returns on the next build, with nothing to toggle.
// The frontend applies the matching noindex in pages/BlogCategory.tsx.
func EmptyCategoryPaths(frDir, enDir string) map[string]bool {
	count := func(dir string) map[string]int {
		out := map[string]int{}
		posts, err := LoadBlogPosts(dir)
		if err != nil {
			return out
		}
		for _, p := range posts {
			out[p.Category]++
		}
		return out
	}
	fr, en := count(frDir), count(enDir)

	empty := map[string]bool{}
	for _, c := range blogCategories {
		if fr[c.Key] == 0 {
			empty["/blog/categorie/"+c.SlugFR] = true
		}
		if en[c.Key] == 0 {
			empty["/en/blog/category/"+c.SlugEN] = true
		}
	}
	return empty
}

// LoadBlogPosts reads the front matter of every Markdown file in dir and
// returns the posts sorted newest first. Missing dir yields an empty slice.
func LoadBlogPosts(dir string) ([]BlogPost, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var posts []BlogPost
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		fm := parseFrontMatter(string(raw))
		slug := strings.TrimSuffix(e.Name(), ".md")
		if v := fm["slug"]; v != "" {
			slug = v
		}
		date, _ := time.Parse("2006-01-02", fm["date"])
		updated := date
		if v := fm["updated"]; v != "" {
			if t, err := time.Parse("2006-01-02", v); err == nil {
				updated = t
			}
		}
		category := fm["category"]
		if category == "" {
			category = DefaultCategory
		}
		posts = append(posts, BlogPost{
			Slug:        slug,
			Title:       fm["title"],
			Description: fm["description"],
			Category:    category,
			Date:        date,
			Updated:     updated,
		})
	}
	sort.Slice(posts, func(i, j int) bool { return posts[i].Date.After(posts[j].Date) })
	return posts, nil
}

// parseFrontMatter extracts a minimal YAML front-matter block (key: "value").
// It intentionally supports only flat string keys — enough for sitemap/llms.txt.
func parseFrontMatter(raw string) map[string]string {
	out := map[string]string{}
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	if !strings.HasPrefix(raw, "---\n") {
		return out
	}
	end := strings.Index(raw[4:], "\n---")
	if end < 0 {
		return out
	}
	block := raw[4 : 4+end]
	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimSpace(line)
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		v = strings.TrimSpace(v)
		v = strings.Trim(v, `"'`)
		out[strings.TrimSpace(k)] = v
	}
	return out
}

// BlogLastMod returns the lastmod date of the pages that are built FROM the
// articles: each locale's blog index and each category archive. Their content
// changes the day an article lands and never on a redeploy, so the date of the
// newest article they list is their date. A path with no article at all is
// absent from the map and ends up with no lastmod — better than a made-up one.
func BlogLastMod(frDir, enDir string) map[string]string {
	out := map[string]string{}

	fill := func(dir, index string, catPath func(blogCategory) string) {
		posts, err := LoadBlogPosts(dir)
		if err != nil {
			return
		}
		newest := func(keep func(BlogPost) bool) string {
			var best time.Time
			for _, p := range posts {
				if keep(p) && p.Updated.After(best) {
					best = p.Updated
				}
			}
			if best.IsZero() {
				return ""
			}
			return best.Format("2006-01-02")
		}
		if d := newest(func(BlogPost) bool { return true }); d != "" {
			out[index] = d
		}
		for _, c := range blogCategories {
			cat := c
			if d := newest(func(p BlogPost) bool { return p.Category == cat.Key }); d != "" {
				out[catPath(cat)] = d
			}
		}
	}

	fill(frDir, "/blog", func(c blogCategory) string { return "/blog/categorie/" + c.SlugFR })
	fill(enDir, "/en/blog", func(c blogCategory) string { return "/en/blog/category/" + c.SlugEN })
	return out
}

// TacticsWeek is one weekly puzzle page, named after the Monday it covers.
type TacticsWeek struct {
	Slug string // "14-09-26", the URL segment
	Date string // "2026-09-14", the same day in sitemap form
}

// TacticsWeeks lists the weekly tactics pages found in dir, newest first.
//
// Each week is a real, indexable page in both locales (/tactiques/<slug> and
// /en/tactics/<slug>, pre-rendered by the frontend and linked from the index),
// but the sitemap only ever declared the index itself: a growing series of
// pages that had to be discovered by crawling. They are the one part of the
// site that publishes on a weekly rhythm — exactly what a sitemap is for.
//
// The file name is JJ-MM-AA, which does NOT sort chronologically as text
// ("31-08-26" would come after "14-09-26"), so the date is rebuilt before
// sorting. A name that doesn't parse is skipped rather than guessed at.
func TacticsWeeks(dir string) []TacticsWeek {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var weeks []TacticsWeek
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		slug := strings.TrimSuffix(e.Name(), ".json")
		d, err := time.Parse("02-01-06", slug)
		if err != nil {
			continue
		}
		weeks = append(weeks, TacticsWeek{Slug: slug, Date: d.Format("2006-01-02")})
	}
	sort.Slice(weeks, func(i, j int) bool { return weeks[i].Date > weeks[j].Date })
	return weeks
}
