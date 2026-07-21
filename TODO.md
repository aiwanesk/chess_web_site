# TODO — iwanesko.ch

Point de reprise. Coché = fait. Dernière mise à jour : 2026-07-21 (soir).

---

## ✅ État actuel (fait / en ligne)

- **Site en ligne** sur https://iwanesko.ch (Jelastic, image Docker Hub
  `alex42000iwa/chess`). SSG : toutes les routes pré-rendues (meta + JSON-LD),
  design premium, logo/favicon, blog 2 catégories (Progresser / Carnet), maillage.
- **Coordonnées réelles** : Alexandre Iwanesko (Maître FIDE), +41 78 783 56 89,
  alexandre@iwanesko.ch, adresse Swiss Tax Horizon (Route de Florissant 2, 1206
  Genève), `sameAs` FIDE 682136. **Niveau 1200–2200 Elo**. Tarifs 120/1000/60 CHF.
- **Formulaire de contact** : lien dans la nav + **envoi e-mail SMTP** (Infomaniak)
  vers alexandre@iwanesko.ch → **fonctionne** ✅. Rate-limit per-IP sur `/api`.
- **Carte OpenStreetMap** (Route de Florissant) sur `/contact` (gratuite, lazy).
- **SEO/GEO** : titres (≤60) + meta + canonical + **un seul H1** + hreflang FR↔EN ;
  `/sitemap.xml`, `/robots.txt` (crawlers IA), `/llms.txt` en ligne ; **image OG**
  1200×630 générée ; **Search Console + sitemap soumis** ✅.
- **Bilingue (partiel)** : FR + EN (`/en/…`) sur Home, Cours adultes, Tarifs,
  Contact, Tactics — sélecteur de langue + hreflang.
- **Tactiques** : fetch parties (chess.com + lichess), détection Stockfish, anonymisation
  (inversion 100 % prouvée légale), **échiquier interactif** `/tactiques` + `GET /api/tactics`.
- **Bandeau « site en construction »** (à retirer au lancement : `<SiteBanner/>` dans `Layout.tsx`).
- **CI** vert (lint/test/build + Lighthouse, gates a11y/BP/CLS). **Release** sur tag `v*` → Docker Hub.

---

## 🔴 Actions immédiates (toi)

- [ ] **REDÉPLOYER** : le site live tourne une image plus ancienne. Les derniers
      correctifs ne sont **pas encore en ligne** (image OG — actuellement **404** en
      live → pas d'aperçu au partage social —, nav qui ne wrap plus, carte contact,
      titres ≤60, 10/10, message CTA Résultats).
      → `git tag v0.1.1 && git push origin v0.1.1` puis mettre à jour le conteneur Jelastic.
- [ ] **Google Business Profile** : créer/aligner (NAP **identique** au footer :
      Route de Florissant 2, 1206 Genève) → pour le *local pack* « coach échecs Genève ».
- [ ] **Search Console** : Inspection d'URL → **Demander l'indexation** des pages clés
      (home, cours adultes) pour accélérer.

---

## 🟠 Prochain chantier dev : tactiques (stats + auto)

- [ ] **Stats privées (SQLite)** : `modernc.org/sqlite` (pur Go, compat distroless)
      sur `DB_PATH=/data/tactics.db` (volume Jelastic `/data`). `POST /api/tactics/event`
      (vues/tentatives/résolu) + **tableau de bord `/admin`** protégé par `ADMIN_TOKEN`.
      Câbler `onSolved`/`onAttempt` de `PuzzleBoard`.
- [ ] **Ingestion runtime** : `POST /api/admin/tactics` (Bearer `ADMIN_TOKEN`) → le
      batch pousse les puzzles **sans redeploy** (modèle A) ; le serveur lit depuis `/data`.
      🔒 **Sécurité (hyper important)** : `ADMIN_TOKEN` obligatoire (vide → endpoint
      **désactivé**) ; comparaison **constant-time** (`subtle.ConstantTimeCompare`) ;
      corps limité + validation stricte ; rate-limit (déjà en place) ; HTTPS only ;
      zéro secret dans les logs. `/admin` derrière le même token ; `POST /api/tactics/event`
      public mais validé + rate-limité.
- [ ] **Cron hebdo** (GitHub Actions, `cron: '0 5 * * 1'` ≈ lundi 6-7h Genève) :
      `apt-get install stockfish` + secrets pseudos → `go run ./cmd/tactics` →
      `POST /api/admin/tactics`. (Texte de la page déjà générique.)
- [ ] Affiner les seuils de détection sur un plus gros échantillon (movetime ↑).

---

## 📝 Contenu à fournir (toi)

- [ ] **Articles de blog** réels (voir `docs/plan-editorial-blog.md`). Le guide
      « Sortir du plateau 1500 Elo » est amorcé → compléter le `[À COMPLÉTER]` (vécu).
      Remplacer l'exemple de carnet par un vrai récit.
- [ ] **Carnet de tournoi** : ~1 récit par tournoi marquant (dupliquer `content/_templates/carnet.md`).
- [ ] `/resultats` : **témoignages réels et attribuables** + vraies stats (placeholders
      actuellement). → **Flavius** (élève) dès réception de son témoignage.
- [ ] **Portrait** d'Alexandre (Home / À propos).
- [ ] Police **Inter variable** self-hostée (`public/fonts/` + `@font-face`).
- [ ] Dates réelles des stages (`pages/Stages.tsx`, `eventSchema`).
- [ ] Coordonnées géo précises pour le pin carte (actuellement approx. Florissant, `site.ts`).

---

## 🌐 Parké (décision : « oublie le rollout EN complet » pour l'instant)

- [ ] Traduire le reste en EN : money pages (tournoi, en-ligne, groupe, ados, stages,
      conférences, team building), pages info (à-propos, résultats), **articles de blog**
      (structure blog i18n), compléter nav/footer EN.
      → Pattern posé : ajouter dans `PAGES` (`src/lib/i18n.tsx`) + route EN. (Déjà fait :
      Home, Cours adultes, Tarifs, Contact, Tactics.)

---

## 🚀 Post-lancement (plus tard)

- [x] Formulaire contact → **e-mail** (SMTP Infomaniak) opérationnel.
- [ ] Accusé de réception e-mail au visiteur + éventuel anti-spam renforcé (Turnstile).
- [ ] **Prise de RDV** : Cal.com (self-host) ou lien externe sur `/contact`.
- [ ] **Analytics privacy-first** : Plausible/Umami self-host (cohérent CSP).
- [ ] Workflow « push `main` → build+push `:latest` » pour fluidifier les redeploys d'articles.
- [ ] Retirer le bandeau « site en construction » au lancement.
