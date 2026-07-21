# TODO — iwanesko.ch

Ce fichier ne garde que **le reste à faire**. Tout ce qui est **terminé** est archivé
dans **[`done.md`](./done.md)** (une ligne par élément, pour référence). Dernière mise
à jour : 2026-07-22.

> 🚀 **Dernier tag = `v0.2.7`** (à déployer sur Jelastic). Embarque tout le terminé
> (voir `done.md`). ⚠️ **Régénération des tactiques en cours** (local, nouvelles règles) —
> commit **en attente de ton feu vert**, pas encore taggé.

---

## 🔴 Actions immédiates (toi)

- [ ] **⭐ DÉPLOYER le conteneur Jelastic sur le dernier `:latest`/tag** (`v0.2.7`).
      Met en ligne : fix Debian, stats + `/admin`, newsletter, accusé de réception,
      EN complet + blog bilingue, tarif groupe, adresse + pin, échiquier + détection tactiques.
- [ ] **Auto-republish Jelastic — OPTION A (recommandée)** : activer **Automatic Updates**
      sur le conteneur Chess (dashboard → réglages du nœud → Automatic Updates), intervalle
      quotidien/horaire, tag **`:latest`**. Jelastic compare le digest de `:latest` et
      **redéploie tout seul** quand il change (après le push du cron), **volumes `/data`
      préservés**. → l'article tactiques du lundi part en ligne sans intervention.
      (Option B, immédiat via API Jelastic depuis GitHub Actions : dispo si besoin, nécessite
      un token API en secret GitHub — non fait, l'option A suffit pour du hebdo.)
- [ ] **Secrets GitHub pour le cron tactiques** (Settings → Secrets → Actions) :
      `CHESSCOM_USERNAME`, `LICHESS_USERNAME` (+ `LICHESS_TOKEN` facultatif). Sans eux,
      le cron du lundi échoue. (La régénération manuelle locale, elle, tourne déjà.)
- [ ] **Google Business Profile** : créer/aligner (NAP **identique** au footer :
      Route de Florissant 2, 1206 Genève) → *local pack* « coach échecs Genève ».
- [ ] **Search Console** : Inspection d'URL → **Demander l'indexation** des pages clés.
- [ ] **Retirer le bandeau « site en construction »** au lancement (`<SiteBanner/>` dans `Layout.tsx`).

---

## 🟠 Dev — reste à faire

- [ ] Affiner les seuils de détection tactiques sur un plus gros échantillon (movetime ↑).
- [ ] *(Mineur)* hreflang réciproque sur les articles de blog individuels (l'`altSlug`
      est déjà dans le frontmatter, à brancher sur `<Seo>`).
- [ ] *(Option future)* **Ingestion runtime** `POST /api/admin/tactics` (Bearer `ADMIN_TOKEN`)
      pour pousser les puzzles **sans redeploy** — non fait, le modèle cron+rebuild suffit.
- [ ] **Prise de RDV** : Cal.com (self-host) ou lien externe sur `/contact`.
- [ ] **Analytics privacy-first** : Plausible/Umami self-host (cohérent CSP).
- [ ] Anti-spam renforcé optionnel (Turnstile) sur contact/newsletter.

---

## 📝 Contenu à fournir (toi)

- [ ] **Articles de blog** réels (voir `docs/plan-editorial-blog.md`). Le blog est **vide**
      (les exemples/placeholders ont été supprimés) → déposer les vrais `.md` dans
      `content/blog/` (FR) et `content/blog/en/` (EN, `lang: "en"` + `altSlug`).
- [ ] **Carnet de tournoi** : ~1 récit par tournoi marquant (gabarit `content/_templates/carnet.md`).
- [ ] `/resultats` : **témoignages réels et attribuables** + vraies stats (placeholders
      actuellement). → **Flavius** (élève) dès réception de son témoignage.
- [ ] **Portrait** d'Alexandre (Home / À propos).
- [ ] Police **Inter variable** self-hostée (`public/fonts/` + `@font-face`).
- [ ] **Dates réelles des stages** (`pages/Stages.tsx`, `eventSchema`).

---

## 🧩 Notes prod (rappel)
- L'image v0.2.x définit `DB_PATH=/data/tactics.db` par défaut + `mkdir /data` + tourne
  en root → la DB SQLite se crée seule au boot sur le volume `/data` (persistant).
- Newsletter : nécessite `DB_PATH` **et** SMTP configuré. Pour annoncer un stage → entrée
  dans `content/events.json` puis redéploie (l'annonce part auto au boot).
- Détail de tout le terminé → **[`done.md`](./done.md)**.
