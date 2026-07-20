# TODO — iwanesko.ch

Suivi des travaux restants. Coché = fait.

## Tactiques de la semaine (batch auto) — EN COURS

- [x] Confidentialité : pseudos dans `.env` (gitignoré) / GitHub Secrets — **jamais
      commités** (repo public). Vérifié : 0 pseudo dans les fichiers suivis.
- [x] `backend/internal/tactics` : config (env), **inversion 100 % FEN/coups**
      (`MirrorFEN`/`MirrorUCIMove` + tests), modèle `Puzzle` anonyme (test anti-fuite),
      clients chess.com + lichess (`FetchWeek`). `cmd/tactics` runnable.
- [x] Fetch prouvé : 239 parties/semaine récupérées (bullet inclus).
- [ ] **Détection Stockfish** (UCI) : combinaisons jouées vs ratées via seuils
      swing d'éval + écart coup-unique ; scoring « beauté » ; top 10 ; écrire
      `content/tactiques/<année>-S<semaine>.json`.
- [x] **Rendu interactif** : `PuzzleBoard` (échiquier jouable, sans lib d'échecs
      côté client) + page `/tactiques` qui charge `/api/tactics`. Backend sert le
      dernier fichier de puzzles (`TACTICS_DIR`). Nav + sitemap OK.
- [ ] **Stats privées (SQLite)** : `modernc.org/sqlite` sur volume `/data`
      (`DB_PATH`), endpoint `POST /api/tactics/event` + tableau de bord `/admin`
      protégé par `ADMIN_TOKEN`. (Dockerfile : `VOLUME /data` déjà déclaré.)
- [ ] **Ingestion runtime** : `POST /api/admin/tactics` (Bearer ADMIN_TOKEN) pour
      que le batch pousse les puzzles sans redeploy (modèle A).
- [ ] **Cron hebdo** (GitHub Actions) : batch Stockfish → push vers `/api/admin/tactics`.

## Traduction anglaise (i18n) — EN COURS

- [x] Infra i18n : FR par défaut (racine) + **EN sous `/en/` avec slugs anglais**.
      Contexte de locale, dictionnaire UI, registre de routes FR↔EN, **sélecteur
      de langue**, **hreflang** bidirectionnel, chrome bilingue (header/footer/seo/form).
- [x] Pages EN livrées : **Home, Cours adultes (phare), Tarifs, Contact** +
      entrées sitemap EN. Vérifié : `lang=en`, hreflang FR↔EN, contenu anglais, 1 H1.
- [ ] **Rollout EN restant** : money pages (tournoi, en-ligne, groupe, ados,
      stages, conférences, team building), pages info (à-propos, résultats),
      **articles de blog** (versions EN), et compléter la nav/footer EN.
      → Ajouter chaque page traduite dans `PAGES` (src/lib/i18n.tsx) + une route EN.
- [ ] Blog i18n : structure EN pour les articles + catégories EN.

## Différé (à traiter plus tard, décision utilisateur)

### Intégrations
- [ ] **Formulaire de contact → email/CRM** : brancher `CONTACT_WEBHOOK_URL`
      (relais mail type Resend/Postmark/Formspree, ou webhook CRM). Actuellement
      les leads sont validés + loggés côté serveur, et forwardés si l'URL est
      définie. Ajouter : accusé de réception e-mail, anti-spam renforcé
      (rate-limit par IP), éventuel reCAPTCHA/Turnstile *self-host-friendly*.
- [ ] **Prise de RDV / booking** : calendrier (Cal.com self-host, ou lien
      externe) intégré à `/contact`.
- [ ] **Analytics privacy-first** : Plausible/Umami self-host (pas de scripts
      tiers bloquants, cohérent avec la CSP).

### Déploiement
- [ ] Pipeline de déploiement réel (build Docker → registry → host).
- [ ] Domaine `iwanesko.ch` + TLS (reverse proxy / plateforme).
- [ ] Renseigner `BASE_URL` de prod, vérifier `sitemap.xml`/canonical.
- [ ] Soumettre le sitemap à Google Search Console + Bing Webmaster.
- [ ] Créer/aligner le **Google Business Profile** (NAP identique au footer).

## Placeholders à remplacer avant mise en ligne
- [x] `frontend/src/lib/site.ts` : nom (Alexandre Iwanesko), téléphone
      (+41 78 783 56 89), `sameAs` (fiche FIDE 682136 ; lichess/chess.com
      volontairement exclus).
- [x] Adresse pro : Swiss Tax Horizon, Route de Florissant 2, Genève
      (footer + contact + schema LocalBusiness).
- [x] NPA **1206** et e-mail **contact@iwanesko.ch** confirmés.
- [ ] Coordonnées géo : approximatives (Florissant) — affiner si besoin pour le pin carte.
- [x] `/tarifs` : prix réels en CHF (120 CHF/h · 1000 CHF le pack 10 · 60 CHF/pers groupe).
- [x] Blog restructuré en 2 catégories : **Progresser** (guides/acquisition) et
      **Carnet de tournoi** (journal de compet'). Templates dans `content/_templates/`.
- [ ] **Écrire les vrais articles** en suivant `docs/plan-editorial-blog.md`.
      Le guide « Sortir du plateau 1500 Elo » est amorcé : compléter l'anecdote vécue
      (bloc `[À COMPLÉTER]`). Remplacer l'exemple de carnet par un vrai récit.
- [ ] Carnet de tournoi : tu joues 2-3 tournois/mois → viser ~1 récit par tournoi
      marquant (dupliquer `content/_templates/carnet.md`).
- [ ] `/resultats` : témoignages **réels et attribuables** uniquement + vraies stats.
- [ ] `public/og/` : images Open Graph 1200×630 par type de page (PNG).
- [ ] Police **Inter variable** self-hostée dans `public/fonts/` + réactiver
      `preload` (index.html) et `@font-face` (styles.css).
- [ ] `/contact` : carte de Genève légère (chargée en différé).
- [ ] Dates réelles des stages (`pages/Stages.tsx`, `eventSchema`).

## En cours (itération design + contenu)
- Voir l'historique git. Money pages étoffées, design premium, blog + maillage.
