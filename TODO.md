# TODO — iwanesko.ch

Point de reprise. Coché = fait. Dernière mise à jour : session du 2026-07-21.

---

## ✅ État actuel (ce qui marche déjà)

- **Site complet** Go + React SSG : toutes les routes pré-rendues (meta + JSON-LD),
  design premium, blog en 2 catégories (Progresser / Carnet de tournoi), maillage.
- **Coordonnées réelles** : Alexandre Iwanesko (Maître FIDE), +41 78 783 56 89,
  alexandre@iwanesko.ch, adresse Swiss Tax Horizon (Route de Florissant 2, 1206
  Genève), `sameAs` fiche FIDE 682136. Tarifs 120 CHF/h · 1000 CHF pack 10 · 60 CHF/pers.
- **Bilingue (partiel)** : FR + EN (`/en/…`) sur Home, Cours adultes, Tarifs,
  Contact — sélecteur de langue + hreflang OK.
- **Bandeau « site en construction »** (bilingue, à retirer au lancement :
  `<SiteBanner/>` dans `Layout.tsx`).
- **Tactiques** : fetch parties (chess.com + lichess), détection Stockfish
  (jouées/ratées), anonymisation (inversion 100 %, prouvée légale), **échiquier
  interactif** sur `/tactiques` + `GET /api/tactics`.
- **CI** vert : lint/test/build + Lighthouse (gates bloquants a11y/best-practices/CLS).
- **Release Docker** : workflow sur tag `v*` → `alex42000iwa/chess` (Docker Hub).

---

## 🔴 Actions immédiates (toi)

- [ ] **Secrets GitHub** (pour le push Docker Hub) : `DOCKERHUB_USERNAME=alex42000iwa`
      + `DOCKERHUB_TOKEN` (nouveau token — révoquer l'ancien exposé). Puis Actions →
      run « Release (Docker Hub) » → **Re-run failed jobs**.
- [ ] **Déploiement Jelastic** :
  - Déployer l'image `alex42000iwa/chess:latest`.
  - Env : `BASE_URL=https://iwanesko.ch` (+ plus tard `ADMIN_TOKEN`).
  - **Volume persistant `/data`** (local, pas de nœud storage réseau) pour SQLite.
  - Domaine `iwanesko.ch` + TLS.
- [ ] Après mise en ligne : soumettre `sitemap.xml` à **Google Search Console**,
      créer/aligner le **Google Business Profile** (NAP identique au footer).

---

## 🟠 Prochain chantier dev : tactiques (stats + auto)

- [ ] **Stats privées (SQLite)** : `modernc.org/sqlite` (pur Go, compat distroless)
      sur `DB_PATH=/data/tactics.db`. `POST /api/tactics/event` (vues/tentatives/
      résolu) + **tableau de bord `/admin`** protégé par `ADMIN_TOKEN`.
      Câbler `onSolved`/`onAttempt` de `PuzzleBoard` vers l'endpoint.
- [ ] **Ingestion runtime** : `POST /api/admin/tactics` (Bearer `ADMIN_TOKEN`) pour
      que le batch pousse les puzzles **sans redeploy** (modèle A). Le serveur lit
      alors depuis `/data` au lieu du fichier bâti.
- [ ] **Cron hebdo** (GitHub Actions) : `apt-get install stockfish` + secrets
      pseudos → `go run ./cmd/tactics` → `POST /api/admin/tactics`.
- [ ] Affiner les seuils de détection sur un plus gros échantillon (movetime ↑,
      plus de parties) — la qualité s'améliore avec les vraies données.

---

## 📝 Contenu à fournir (toi)

- [ ] **Articles de blog** réels (voir `docs/plan-editorial-blog.md`). Le guide
      « Sortir du plateau 1500 Elo » est amorcé → compléter le bloc `[À COMPLÉTER]`
      (ton vécu). Remplacer l'exemple de carnet par un vrai récit.
- [ ] **Carnet de tournoi** : ~1 récit par tournoi marquant (2-3/mois) —
      dupliquer `content/_templates/carnet.md`.
- [ ] `/resultats` : témoignages **réels et attribuables** + vraies stats.
- [ ] **Portrait** d'Alexandre (Home / À propos).
- [ ] `public/og/` : images Open Graph 1200×630 (PNG).
- [ ] Police **Inter variable** self-hostée (`public/fonts/` + `@font-face`).
- [ ] `/contact` : carte de Genève légère (chargée en différé).
- [ ] Dates réelles des stages (`pages/Stages.tsx`, `eventSchema`).
- [ ] Coordonnées géo précises (actuellement approx. Florissant).

---

## 🌐 Parké (décision : « oublie le rollout EN complet » pour l'instant)

- [ ] Traduire le reste en EN : money pages (tournoi, en-ligne, groupe, ados,
      stages, conférences, team building), pages info (à-propos, résultats),
      **articles de blog** (structure blog i18n), et compléter nav/footer EN.
      → Pattern posé : ajouter dans `PAGES` (`src/lib/i18n.tsx`) + route EN.

---

## 🚀 Post-lancement (plus tard)

- [ ] **Formulaire contact → email/CRM** : `CONTACT_WEBHOOK_URL` (Resend/Postmark/
      Formspree) + accusé de réception + rate-limit / anti-spam.
- [ ] **Prise de RDV** : Cal.com (self-host) ou lien externe sur `/contact`.
- [ ] **Analytics privacy-first** : Plausible/Umami self-host (cohérent CSP).
- [ ] Workflow optionnel « push `main` → build+push `:latest` » pour fluidifier
      les redeploys d'articles (au lieu d'un tag à chaque fois).
