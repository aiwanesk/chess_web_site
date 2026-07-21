# TODO — iwanesko.ch

Point de reprise. Coché = fait. Dernière mise à jour : 2026-07-22.

> 🚀 **`v0.2.4` = dernier tag, à déployer.** Il embarque TOUT : fix déploiement
> (base Debian, voir ci-dessous), stats + `/admin`, newsletter RGPD, accusé de
> réception contact, **rollout EN complet + blog bilingue**, tarif groupe (60 min),
> et le **pin GPS précis** de la carte. → **Action (toi) : mettre à jour le conteneur
> Jelastic sur `v0.2.4`** (les tags v0.2.0→v0.2.3 sont superseded).
>
> ✅ **CAUSE RACINE des échecs de déploiement** : Virtuozzo/Jelastic injecte
> `curl` + `iptables-persistent` via `apt` à chaque deploy → une image
> scratch/distroless (sans apt ni shell) fait **échouer tout le déploiement** (et
> casse le Web SSH). **Fix** : runtime rebasé sur `debian:bookworm-slim`
> (+ ca-certificates, `mkdir /data`, root). Binaire Go pur inchangé.
> Après déploiement OK : remettre le mount NFS `/data` proprement (persistance).

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
- **Bilingue COMPLET** : FR + EN (`/en/…`) sur **toutes** les pages + **blog bilingue**
  — sélecteur de langue + hreflang FR↔EN. (Committé, arrive en ligne avec `v0.2.4`.)
- **Tactiques** : fetch parties (chess.com + lichess), détection Stockfish, anonymisation
  (inversion 100 % prouvée légale), **échiquier interactif** `/tactiques` + `GET /api/tactics`.
- **Bandeau « site en construction »** (à retirer au lancement : `<SiteBanner/>` dans `Layout.tsx`).
- **CI** vert (lint/test/build + Lighthouse, gates a11y/BP/CLS). **Release** sur tag `v*` → Docker Hub.

### 🛠️ Committé (dans `v0.2.4`), en ligne dès le redéploiement Jelastic
- **Stats privées + `/admin`** : `internal/stats`, `POST /api/tactics/event`,
  dashboard Basic Auth (token constant-time, rate-limité anti-brute-force).
- **Newsletter RGPD/nLPD** : `internal/newsletter`, double opt-in, désabo 1 clic,
  envoi **auto** à la publication (blog + tactiques + `content/events.json`),
  formulaire dans le footer, page **/confidentialite** (+ `/en/privacy`).
- **DB SQLite** créée automatiquement au boot si `DB_PATH` défini (garde-fou
  non-bloquant : une DB inaccessible désactive stats/newsletter sans casser le site).
- **Accusé de réception** e-mail au visiteur après le formulaire de contact.
- **Rollout EN complet** (toutes pages + blog bilingue) + **pin GPS carte** précis.

---

## 🔴 Actions immédiates (toi)

- [x] **Adresse** ✅ : « Route de Florissant 2, 1206 Genève » confirmée exacte dans le
      code (`src/lib/site.ts`) — le « faux » vu en ligne = l'ancienne image `v0.1.0`.
      Pin GPS carte corrigé (46.196817, 6.153666).
- [x] **PUSH + TAG `v0.2.4`** ✅ (dernier tag ; supersede v0.2.0→v0.2.3). Build Docker
      Hub via release.yml (suivre onglet **Actions**).
- [ ] **⭐ METTRE À JOUR LE CONTENEUR JELASTIC sur `v0.2.4`** une fois le build terminé.
      C'est ce qui met en ligne : fix déploiement Debian, stats + `/admin`, newsletter,
      accusé de réception, **EN complet + blog bilingue**, tarif groupe, adresse + pin.
- [ ] **Jelastic — volume + variables** (pour stats + newsletter) :
      volume `/data` en lecture/écriture **déjà en place** ✅ ; ajouter les variables
      `DB_PATH=/data/tactics.db` et `ADMIN_TOKEN=<token aléatoire>` (`ADMIN_TOKEN`
      déjà posé). La DB se crée seule au 1er boot dans `/data` (volume persistant,
      survit aux redeploys ; local = mieux que NFS pour SQLite).
- [ ] **Google Business Profile** : créer/aligner (NAP **identique** au footer :
      Route de Florissant 2, 1206 Genève) → pour le *local pack* « coach échecs Genève ».
- [ ] **Search Console** : Inspection d'URL → **Demander l'indexation** des pages clés
      (home, cours adultes) pour accélérer.
- [ ] **Secrets GitHub pour le cron tactiques** (Settings → Secrets → Actions) :
      `CHESSCOM_USERNAME`, `LICHESS_USERNAME` (+ `LICHESS_TOKEN` facultatif). Sans
      eux, le cron hebdo échoue. → Puis Actions → « Tactics (weekly) » → **Run workflow**
      pour tester tout de suite.
- [ ] **Jelastic auto-update** sur nouveau `:latest` (ou webhook) → pour que l'article
      hebdo parte en ligne sans intervention. Sinon, redéploiement manuel après le cron.

---

## 🟠 Prochain chantier dev : tactiques (stats + auto)

- [x] **Stats privées (SQLite)** ✅ : `internal/stats` (`modernc.org/sqlite`, pur Go,
      compat distroless). `POST /api/tactics/event` (vues/tentatives/résolu, validé +
      rate-limité + corps ≤ 1KB) et **tableau de bord `/admin`**. `onView`/`onAttempt`/
      `onSolved` câblés dans `PuzzleBoard` → `lib/tacticsEvents.ts` (fire-and-forget).
      🔒 Sécurité : `ADMIN_TOKEN` obligatoire (vide → `/admin` **404, désactivé**),
      Basic Auth **constant-time** (`subtle.ConstantTimeCompare`), `noindex`/`no-store` ;
      compteurs agrégés uniquement (pas d'IP, pas d'identifiant). Tests inclus.
      ⚠️ **En prod** : définir `DB_PATH=/data/tactics.db` sur le volume persistant
      Jelastic `/data` + `ADMIN_TOKEN` (secret). Vide en local → stats désactivées.
- [ ] *(Option future)* **Ingestion runtime** `POST /api/admin/tactics` (Bearer
      `ADMIN_TOKEN`) pour pousser les puzzles **sans redeploy** — non fait, le cron
      commite + rebuild l'image (modèle actuel suffisant).
- [x] **Cron hebdo** (`.github/workflows/tactics-weekly.yml`, lundi 05:00 UTC) :
      Stockfish + go run cmd/tactics → commit `content/tactiques/JJ-MM-AA.json` →
      rebuild image. ⚠️ Nécessite les secrets pseudos (voir Actions immédiates).
- [ ] Affiner les seuils de détection sur un plus gros échantillon (movetime ↑).

---

## 🟢 Newsletter (RGPD + nLPD) ✅

- [x] **Liste d'abonnés conforme, sans service payant** : `internal/newsletter`
      (même DB SQLite que les stats). Consentement **explicite + double opt-in**
      (case non pré-cochée → mail de confirmation → statut `confirmed`).
      Minimisation : on stocke **uniquement** e-mail + langue + statut + 2 jetons
      + horodatages (pas d'IP, pas de nom, jamais loggé). Désinscription 1 clic
      (suppression = droit à l'effacement) + en-têtes `List-Unsubscribe`.
- [x] **Envoi 100% automatique à la publication** : au démarrage, le serveur
      scanne blog + tactiques + `content/events.json` et notifie les abonnés
      confirmés pour tout contenu **nouveau**. 1er run = *seed* (marque l'existant
      comme déjà notifié → aucun envoi en masse du back-catalogue). Throttle SMTP.
- [x] **Formulaire** dans le footer + page **/confidentialite** (FR) et
      **/en/privacy** (EN), liées depuis le footer + la case de consentement.
      ⚠️ **En prod** : nécessite `DB_PATH` **et** le SMTP configuré (sinon
      inscription stockée mais mail de confirmation non envoyé). Pour annoncer un
      stage : ajouter une entrée dans `content/events.json` (`{id,title,url,...}`)
      puis redéployer.

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

## 🌐 Rollout EN COMPLET ✅ (fait — attend `v0.2.3`)

- [x] **Toutes les pages traduites** : money pages (tournoi, en-ligne, groupe, ados,
      stages, conférences, team building) + à-propos + résultats + tarifs + contact +
      tactiques + confidentialité. Routes EN, `PAGES`/hreflang FR↔EN, sitemap EN, footer EN.
- [x] **Blog bilingue** : catégories locale-aware (clé stable + slug/label FR/EN),
      loaders (`content/blog` FR, `content/blog/en` EN), PostCard, index, archives,
      articles. Routes `/en/blog`, `/en/blog/category/:cat`, `/en/blog/:slug`.
      Article « Breaking the 1500 Elo plateau » traduit + gabarit carnet EN.
- [ ] *(Mineur)* hreflang réciproque sur les articles de blog individuels (l'`altSlug`
      est déjà dans le frontmatter, il reste à le brancher sur `<Seo>`).

---

## 🚀 Post-lancement (plus tard)

- [x] Formulaire contact → **e-mail** (SMTP Infomaniak) opérationnel.
- [ ] Accusé de réception e-mail au visiteur + éventuel anti-spam renforcé (Turnstile).
- [ ] **Prise de RDV** : Cal.com (self-host) ou lien externe sur `/contact`.
- [ ] **Analytics privacy-first** : Plausible/Umami self-host (cohérent CSP).
- [ ] Workflow « push `main` → build+push `:latest` » pour fluidifier les redeploys d'articles.
- [ ] Retirer le bandeau « site en construction » au lancement.
