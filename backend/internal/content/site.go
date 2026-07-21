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
type Page struct {
	Path       string  // e.g. "/cours-echecs-adultes-geneve"
	Title      string  // human label, reused in llms.txt
	Changefreq string  // sitemap changefreq hint
	Priority   float64 // sitemap priority, 0..1
	Summary    string  // one-line description for llms.txt / GEO
}

// StaticPages lists every non-blog indexable route. Keep this aligned with the
// frontend router in frontend/src/routes.tsx.
var StaticPages = []Page{
	{"/", "Accueil — Alexandre Iwanesko, Maître FIDE & coach d'échecs à Genève", "weekly", 1.0,
		"Coach d'échecs à Genève, Maître FIDE. Cours pour adultes 1200–2200 Elo et ados en compétition, préparation tournoi, en présentiel et en ligne."},
	{"/cours-echecs-adultes-geneve", "Cours d'échecs pour adultes à Genève", "monthly", 0.9,
		"Cours d'échecs pour adultes (1200–2200 Elo) à Genève avec un Maître FIDE : méthode structurée, plan de progression, présentiel et en ligne."},
	{"/preparation-tournoi-echecs", "Préparation tournoi d'échecs", "monthly", 0.9,
		"Préparation ciblée aux tournois d'échecs : ouvertures, gestion du temps, préparation adverse et mental de compétition."},
	{"/cours-echecs-en-ligne", "Cours d'échecs en ligne", "monthly", 0.8,
		"Cours d'échecs particuliers en ligne avec un Maître FIDE, depuis toute la Suisse romande et la France voisine."},
	{"/cours-echecs-groupe-geneve", "Cours d'échecs en groupe à Genève", "monthly", 0.8,
		"Cours d'échecs en petit groupe à Genève : émulation, tarif réduit, niveau homogène."},
	{"/cours-echecs-ados-competition", "Cours d'échecs pour ados en compétition", "monthly", 0.8,
		"Coaching pour adolescents joueurs de compétition : progression Elo, préparation tournoi et suivi individualisé."},
	{"/stages-echecs-geneve", "Stages d'échecs à Genève", "monthly", 0.7,
		"Stages d'échecs intensifs à Genève pendant les vacances scolaires, encadrés par un Maître FIDE."},
	{"/conferences-echecs-entreprise", "Conférences d'échecs en entreprise", "monthly", 0.6,
		"Conférences et interventions échecs en entreprise : stratégie, prise de décision et gestion du risque."},
	{"/team-building-echecs-geneve", "Team building échecs à Genève", "monthly", 0.6,
		"Ateliers de team building autour des échecs pour entreprises à Genève et dans l'arc lémanique."},
	{"/a-propos", "À propos d'Alexandre Iwanesko, Maître FIDE", "yearly", 0.7,
		"Parcours d'Alexandre Iwanesko, Maître FIDE et coach d'échecs à Genève : titre, résultats et méthode d'enseignement."},
	{"/resultats", "Résultats & témoignages", "monthly", 0.7,
		"Résultats des élèves et témoignages : progressions Elo, performances en tournoi et retours d'expérience."},
	{"/tarifs", "Tarifs des cours d'échecs", "monthly", 0.8,
		"Tarifs des cours d'échecs à Genève : cours particuliers, en groupe, en ligne et forfaits de préparation tournoi."},
	{"/contact", "Contact", "yearly", 0.6,
		"Contacter Alexandre Iwanesko pour un cours d'échecs à Genève ou en ligne : premier échange pour définir vos objectifs."},
	{"/confidentialite", "Politique de confidentialité", "yearly", 0.3,
		"Traitement des données personnelles sur iwanesko.ch (contact et newsletter), conforme au RGPD et à la LPD suisse."},
	{"/blog", "Blog échecs", "weekly", 0.7,
		"Articles d'échecs : stratégie, ouvertures, préparation tournoi et progression pour joueurs intermédiaires et avancés."},
	{"/tactiques", "Tactiques de la semaine", "weekly", 0.7,
		"Résolvez les plus belles tactiques d'échecs de la semaine, sélectionnées par un Maître FIDE. Nouveaux puzzles chaque semaine."},
	{"/blog/categorie/progresser", "Progresser aux échecs — guides", "weekly", 0.7,
		"Guides d'un Maître FIDE pour progresser aux échecs : ouvertures, finales, tactique, préparation de tournoi et mental."},
	{"/blog/categorie/carnet-de-tournoi", "Carnet de tournoi", "weekly", 0.6,
		"Le journal de compétition d'Alexandre Iwanesko, Maître FIDE : parties marquantes, décisions sous pression et leçons de tournoi."},

	// English (EN) — grows as the i18n rollout continues.
	{"/en", "Chess coach in Geneva — Alexandre Iwanesko, FIDE Master", "weekly", 0.9,
		"Chess coach in Geneva, FIDE Master. Lessons for adults (1200–2200 Elo) and competitive teens, tournament prep, in person and online."},
	{"/en/adult-chess-lessons-geneva", "Adult chess lessons in Geneva", "monthly", 0.8,
		"Chess lessons for adults (1200–2200 Elo) in Geneva with a FIDE Master: structured method, progression plan, in person and online."},
	{"/en/tournament-preparation", "Chess tournament preparation", "monthly", 0.8,
		"Targeted chess tournament preparation with a FIDE Master: repertoire, opponent prep, time management and competitive mindset."},
	{"/en/online-chess-lessons", "Online chess lessons", "monthly", 0.7,
		"Private online chess lessons with a FIDE Master, from anywhere in French-speaking Switzerland and neighbouring France."},
	{"/en/group-chess-lessons-geneva", "Group chess lessons in Geneva", "monthly", 0.7,
		"Small-group chess lessons in Geneva: matched level, healthy competition and a reduced per-person rate."},
	{"/en/junior-chess-coaching", "Junior chess coaching", "monthly", 0.7,
		"Coaching for competitive teenagers: Elo progression, tournament preparation and individual follow-up."},
	{"/en/chess-camps-geneva", "Chess camps in Geneva", "monthly", 0.6,
		"Intensive chess camps in Geneva during the school holidays, led by a FIDE Master."},
	{"/en/corporate-chess-talks", "Corporate chess talks", "monthly", 0.5,
		"Chess talks and keynotes for companies: strategy, decision-making and risk management."},
	{"/en/chess-team-building-geneva", "Chess team building in Geneva", "monthly", 0.5,
		"Chess team-building workshops for companies in Geneva and the Lake Geneva region."},
	{"/en/about", "About Alexandre Iwanesko, FIDE Master", "yearly", 0.6,
		"The path of Alexandre Iwanesko, FIDE Master and chess coach in Geneva: title, results and teaching method."},
	{"/en/results", "Results & testimonials", "monthly", 0.6,
		"Student results and testimonials: Elo progress, tournament performances and feedback."},
	{"/en/tactics", "Tactics of the week", "weekly", 0.6,
		"Solve the best chess tactics of the week, hand-picked by a FIDE Master. New puzzles every week."},
	{"/en/pricing", "Chess lesson pricing", "monthly", 0.7,
		"Chess lesson pricing in Geneva: private lessons, group lessons, online and packages."},
	{"/en/contact", "Contact", "yearly", 0.6,
		"Contact Alexandre Iwanesko, FIDE Master, for chess lessons in Geneva or online."},
	{"/en/privacy", "Privacy policy", "yearly", 0.3,
		"How personal data is handled on iwanesko.ch (contact and newsletter), GDPR and Swiss FADP compliant."},
}

// BlogPost is the metadata parsed from a Markdown file's front matter.
type BlogPost struct {
	Slug        string
	Title       string
	Description string
	Date        time.Time
	Updated     time.Time
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
		posts = append(posts, BlogPost{
			Slug:        slug,
			Title:       fm["title"],
			Description: fm["description"],
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
