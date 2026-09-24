package server

import (
	"net/http"

	"github.com/iwanesko/chess-web-site/backend/internal/content"
)

// gonePosts sont les articles retirés du dépôt dont l'URL reste indexée, et la
// page la plus proche du sujet vers laquelle les envoyer.
//
// Un 404 sur une URL qui a été indexée perd tout d'un coup : le lien entrant,
// la position acquise, et le budget de crawl que Google dépense à revenir la
// relire pour rien. L'archive de la catégorie est la destination honnête — pas
// l'accueil, qui ne répond pas à la question que se posait le visiteur.
//
// Ces deux-là sont le placeholder « plateau 1500 Elo » supprimé le 22.07.2026
// (commit 6e7270d), FR et EN. La Search Console n'a remonté que la version FR,
// mais l'anglaise est morte de la même façon.
var gonePosts = map[string]string{
	"/blog/sortir-du-plateau-1500-elo":       "/blog/categorie/progresser",
	"/en/blog/breaking-the-1500-elo-plateau": "/en/blog/category/improve",
}

// buildRedirects assemble la table des 301, une fois au démarrage : le contenu
// vit dans l'image, il ne change pas entre deux requêtes.
//
// Outre les articles retirés, elle rattrape le **mauvais préfixe de langue**.
// Un slug anglais servi sous /blog/ (ou l'inverse) n'est pas une URL qui a
// existé : c'est une URL qu'un lien interne fautif a fabriquée, et que Google a
// suivie. Le lien est corrigé (`articleUrl` dans i18n.tsx), mais l'URL est déjà
// indexée. Plutôt qu'une ligne en dur par cas, la règle se déduit des fichiers
// présents : un slug qui n'existe que dans une langue redirige vers cette
// langue. Elle rattrapera donc aussi la prochaine, sans rien à écrire.
//
// Un slug présent dans les DEUX langues (les opens espagnols partagent le leur)
// n'entre pas dans la table : les deux URLs sont bonnes.
func buildRedirects(contentDir string) map[string]string {
	out := make(map[string]string, len(gonePosts)+8)
	for from, to := range gonePosts {
		out[from] = to
	}

	slugs := func(dir string) map[string]bool {
		posts, err := content.LoadBlogPosts(dir)
		if err != nil {
			return nil
		}
		set := make(map[string]bool, len(posts))
		for _, p := range posts {
			set[p.Slug] = true
		}
		return set
	}
	fr, en := slugs(contentDir), slugs(contentDir+"/en")

	for slug := range en {
		if !fr[slug] {
			out["/blog/"+slug] = "/en/blog/" + slug
		}
	}
	for slug := range fr {
		if !en[slug] {
			out["/en/blog/"+slug] = "/blog/" + slug
		}
	}
	return out
}

// redirectTo répond 301 vers target en conservant la chaîne de requête, pour
// qu'un lien campagne garde son attribution jusqu'à la page d'arrivée.
func redirectTo(w http.ResponseWriter, r *http.Request, target string) {
	u := *r.URL
	u.Path = target
	http.Redirect(w, r, u.RequestURI(), http.StatusMovedPermanently)
}
