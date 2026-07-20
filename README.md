# iwanesko.ch — site vitrine & acquisition (Go + React SSG)

Site pour **Alex Iwanesko, Maître FIDE & coach d’échecs à Genève**, optimisé pour
le **SEO technique**, le **GEO** (citabilité par les moteurs génératifs) et les
**Core Web Vitals**.

Architecture : le frontend React est **pré-rendu en HTML statique au build**
(SSG via `vite-react-ssg`) ; le binaire **Go** embarque ce HTML et le sert avec
compression Brotli/gzip et cache immuable, tout en exposant une API (formulaire
de contact) et en générant `sitemap.xml`, `robots.txt` et `llms.txt`.

> **Le HTML servi au premier octet contient tout le texte, les balises meta et le
> JSON-LD, sans exécution JS.** Vérifiable au `curl` (voir plus bas).

---

## Structure

```
.
├── backend/            Serveur Go (chi)
│   ├── main.go         Entrée + graceful shutdown
│   ├── embed.go        //go:embed dist  (frontend embarqué)
│   ├── dist/           Placeholder tracké ; remplacé par le build réel
│   └── internal/
│       ├── server/     Router, middleware, static, seo (sitemap/robots/llms), API contact
│       └── content/    Source de vérité des routes + lecture des articles (Markdown)
├── frontend/           React 18 + TS + Vite + Tailwind v4 + vite-react-ssg
│   └── src/
│       ├── routes.tsx  Une entrée par URL indexable (lazy → 1 chunk/route)
│       ├── lib/        seo.tsx (meta), schema.ts (JSON-LD), site.ts (NAP/entité), content.ts (blog)
│       ├── components/ Layout, Header, Footer, Breadcrumbs, Faq, MoneyPage (template)…
│       └── pages/      Une page par route
├── content/blog/       Articles en Markdown (front-matter + corps)
├── public/             Assets statiques (favicon…) copiés dans dist/
├── scripts/            build-local.(ps1|sh) : build front + embed + binaire Go
├── Dockerfile          Multi-stage : build front → embarqué dans le binaire Go
├── lighthouserc.json   Budget de perf (casse le build si régression)
└── .github/workflows/  CI : vet/test Go + typecheck/lint/build/Lighthouse front
```

---

## Prérequis

- **Go ≥ 1.26**
- **Node ≥ 20** (testé avec 22/26)

## Développement

Deux terminaux, avec rechargement à chaud du front sans rebuild du binaire Go :

```bash
# 1) Frontend en watch (SSG dev server)
cd frontend
npm install
npm run dev            # http://localhost:5173

# 2) (optionnel) Backend servant le build front + l’API
cd frontend && npm run build          # génère frontend/dist
cd ../backend
DEV_STATIC_DIR=../frontend/dist BASE_URL=http://localhost:8080 go run .
# → http://localhost:8080  (sert le HTML pré-rendu + /api + sitemap/robots/llms)
```

Variables d’environnement : voir [`.env.example`](.env.example).

## Build de production (binaire unique avec front embarqué)

```bash
# Windows
./scripts/build-local.ps1
# Linux/macOS
./scripts/build-local.sh

# puis
BASE_URL=https://iwanesko.ch ./server        # (server.exe sous Windows)
```

Ou via Docker (multi-stage, image `distroless` non-root) :

```bash
docker build -t iwanesko-web .
docker run -p 8080:8080 -e BASE_URL=https://iwanesko.ch iwanesko-web
```

## Tests & qualité

```bash
cd backend  && go test ./... -race        # serveur, sitemap, robots, llms, contact
cd frontend && npm run typecheck && npm run lint
```

## Vérifier le rendu serveur (SEO/GEO)

Le contenu doit être présent dans la réponse **brute**, sans JS :

```bash
curl -s http://localhost:8080/cours-echecs-adultes-geneve | grep -c "application/ld+json"   # → 3
curl -s http://localhost:8080/cours-echecs-adultes-geneve | grep -o "<h1[^>]*>"             # → 1 seul
curl -s http://localhost:8080/sitemap.xml
curl -s http://localhost:8080/robots.txt        # GPTBot, PerplexityBot, ClaudeBot… autorisés
curl -s http://localhost:8080/llms.txt          # résumé structuré pour les LLM
```

---

## Ajouter du contenu

### Un article de blog
Créez `content/blog/mon-article.md` :

```markdown
---
title: "Titre de l’article (≤ 60 car. idéalement)"
description: "Meta description ≤ 155 caractères."
author: "Alex Iwanesko"
date: "2026-07-20"
updated: "2026-07-21"        # optionnel
cluster: "adultes"            # silo de maillage interne (optionnel)
clusterPath: "/cours-echecs-adultes-geneve"   # money page liée (optionnel)
---

Corps en **Markdown**. Compilé en HTML au build ; la page est pré-rendue,
ajoutée au sitemap et à llms.txt automatiquement.
```

Rien d’autre à faire : le slug vient du nom de fichier, la route `/blog/:slug`
est pré-rendue via `getStaticPaths`.

### Une nouvelle page (ex. une money page)
1. Créez `frontend/src/pages/MaPage.tsx` exportant `function Component()`.
   Le plus simple : réutilisez le template `MoneyPage` (voir
   `src/pages/CoursAdultesGeneve.tsx`, la **page de référence 100 % optimisée**).
2. Ajoutez la route dans `frontend/src/routes.tsx`.
3. Ajoutez l’entrée correspondante dans `backend/internal/content/site.go`
   (`StaticPages`) pour le sitemap et `llms.txt`.

> Gardez `routes.tsx` et `site.go` alignés : ce sont les deux sources de vérité
> des URLs indexables.

---

## SEO / GEO — ce qui est déjà en place

- **SSG** : chaque route indexable = un HTML complet (texte + meta + JSON-LD) au 1er octet.
- **Par page** : `<title>` unique, meta description, canonical, Open Graph, Twitter Card, `robots`, `hreflang` (FR + prêt pour EN).
- **JSON-LD** : `Person`, `ProfessionalService`/`LocalBusiness`, `Course`+`Offer` (CHF), `Event` (stages), `Review`+`AggregateRating`, `FAQPage`, `BreadcrumbList`.
- **GEO** : `/llms.txt`, crawlers IA autorisés dans `robots.txt`, titres en questions, réponses autoportantes, FAQ en bas des money pages, faits chiffrés citables, entité `Person` cohérente avec `sameAs`.
- **Maillage interne en silos** : home → money pages, money pages entre elles, articles → money page de leur cluster.
- **Perf** : CSS/JS code-splittés par route, HTML pré-rendu, cache `immutable` sur `/assets/*`, Brotli/gzip, CSP stricte, images à dimensionner (anti-CLS).
- **A11y** : HTML sémantique, skip-link, focus visibles, `aria`, navigation clavier, `prefers-reduced-motion`.

## À personnaliser avant mise en ligne (placeholders)

- `frontend/src/lib/site.ts` : **NAP** (téléphone, adresse), `sameAs` (FIDE,
  lichess/chess.com, LinkedIn), coordonnées géo — à aligner sur le Google
  Business Profile.
- **Tarifs** et **témoignages** (`/tarifs`, `/resultats`) : remplacer les
  valeurs indicatives par les vraies.
- **Polices** : self-host `inter-variable.woff2` dans `public/fonts/` puis
  réactiver le `preload` (voir `frontend/index.html`) et le `@font-face`
  (`src/styles.css`). Un stack système est utilisé en attendant (0 requête).
- **Images OG** : ajouter `public/og/default.png` (1200×630).
- **Carte contact** : intégrer une carte légère (chargée en différé) sur `/contact`.
