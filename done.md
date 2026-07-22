# done.md — iwanesko.ch (archive du terminé)

Ce fichier liste ce qui est **fait** (référencé depuis `TODO.md`, qui ne garde que
le reste-à-faire). Une ligne = une chose terminée ; le détail est dans le code et
l'historique git. Dernière mise à jour : 2026-07-22.

## Socle / infra
- **Site + SSG** : toutes les routes pré-rendues (meta + JSON-LD), design premium,
  logo/favicon, maillage interne. Backend Go (chi) sert le front embarqué (`go:embed`).
- **CI** verte : lint/test/build + Lighthouse (gates a11y/BP/CLS). **Release** sur tag `v*` → Docker Hub.
- **Cause racine des échecs de déploiement trouvée + corrigée** : Virtuozzo injecte
  `curl`/`iptables-persistent` via `apt` → image scratch/distroless cassait tout.
  Runtime rebasé sur `debian:bookworm-slim` (+ ca-certificates, `mkdir /data`, root).

## Identité / NAP / SEO
- **Coordonnées réelles** : Alexandre Iwanesko (Maître FIDE), +41 78 783 56 89,
  alexandre@iwanesko.ch, Route de Florissant 2, 1206 Genève, `sameAs` FIDE 682136,
  niveau 1200–2200 Elo, tarifs 120/1000/60 CHF (groupe = 60 CHF/pers/60 min).
- **Pin GPS carte** précis (46.196817, 6.153666) + carte OpenStreetMap sur `/contact`.
- **SEO/GEO** : titres ≤60, meta, canonical, un seul H1, `/sitemap.xml`, `/robots.txt`
  (crawlers IA), `/llms.txt`, image OG 1200×630, Search Console + sitemap soumis.

## Bilingue (rollout EN COMPLET)
- **Toutes les pages** traduites FR/EN (`/en/…`) : money pages, à-propos, résultats,
  tarifs, contact, tactiques, confidentialité — routes EN, `PAGES`/hreflang, sitemap EN, footer EN.
- **Blog bilingue** : catégories locale-aware (clé stable + slug FR/EN), loaders
  (`content/blog` FR, `content/blog/en` EN), PostCard, index, archives, articles EN.

## Contact / e-mail
- **Formulaire de contact** → envoi SMTP Infomaniak vers alexandre@iwanesko.ch (rate-limit /api).
- **Accusé de réception** e-mail au visiteur (dans sa langue, Reply-To = Alexandre).
- Carte contact alignée sur le titre ; doublon « Genève et Genève » corrigé.

## Tactiques
- Fetch parties (chess.com + lichess), détection Stockfish, **anonymisation par inversion**
  (prouvée légale), sortie 100 % anonyme (leak-guard testé).
- **Échiquier façon lichess** : grille 8×8 (fix rangées vides), couleurs, coordonnées,
  coup surligné, état résolu clair, bouton « Voir la solution ».
- **Détection** : 1er coup = échec ou sacrifice (pas de prise gratuite en 1) ; solution
  coupée à sa partie forcée et **finissant sur le coup décisif** (prise/échec/promo/mat) ; cap 7 demi-coups.
- **Cron hebdo** (`tactics-weekly.yml`, lundi 05:00 UTC) : génère → commit → rebuild `:latest`.
- Templates/exemples d'articles supprimés ; blog vide en attendant les vrais articles.

## Stats privées + /admin
- `internal/stats` (SQLite `modernc`, pur Go), `POST /api/tactics/event` (validé, rate-limité, ≤1 KB),
  dashboard **`/admin`** Basic Auth **constant-time**, `ADMIN_TOKEN` obligatoire (vide → 404),
  compteurs agrégés only (pas d'IP/identifiant). Garde-fou : DB inaccessible → stats désactivées, site debout.

## Analytics & réservation (full privacy)
- **Analytics self-hosté** : comptage des pages vues côté serveur (agrégats `(jour, path)`,
  zéro IP/cookie/tiers), affiché dans `/admin` (fréquentation 14 j + top pages).
- **Réservation de cours** : `internal/booking` (SQLite, anti-double-réservation),
  `POST /api/booking` (créneaux 17h30–20h00 par 30 min, calcul du montant via `HOURLY_RATE`,
  e-mail au visiteur + à Alexandre), page **/reserver** (+ **/en/book**) avec calendrier
  et prix en direct, CTA « Réserver » du header pointant dessus, réservations listées dans `/admin`.

## Anti-spam (fait maison, zéro tiers)
- **Jeton HMAC** : `GET /api/form-token` signe un jeton court ; contact/newsletter/réservation
  l'exigent → un bot qui POST direct sur l'API (sans jeton) est **403**. Front : jeton
  préchargé au montage + refresh transparent sur 403.
- **Filtre anti-liens** sur le texte libre du contact (silencieusement ignoré).
- **Rate-limit resserré** partagé sur les 3 formulaires (~6 puis 1 / 20 s) + honeypot existant.

## Divers SEO / sécu
- **security.txt** (`/.well-known/security.txt`, RFC 9116), servi par le backend.
- **hreflang** sur les articles de blog (via `altSlug` du frontmatter) + articles EN dans le sitemap.
- Durcissement anti-scan (404 silencieux sur wp-login.php, .env, etc.). Movetime tactiques 250→400.
- **Event JSON-LD** complété (eventStatus, image, organizer.name, offers.availability/validFrom).

## Newsletter (RGPD + nLPD)
- `internal/newsletter` (même DB) : **double opt-in**, minimisation (email+langue+jetons),
  **désabo 1 clic** (suppression) + `List-Unsubscribe`, page **/confidentialite** (+ `/en/privacy`).
- **Envoi 100 % auto** à la publication (blog + tactiques + `content/events.json`), seed au 1er run.
- **Notif à Alexandre** par e-mail quand un abonné confirme.
