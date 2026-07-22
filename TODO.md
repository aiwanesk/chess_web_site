# TODO — iwanesko.ch

Ce fichier ne garde que **le reste à faire**. Tout ce qui est **terminé** est archivé
dans **[`done.md`](./done.md)** (une ligne par élément, pour référence). Dernière mise
à jour : 2026-07-22.

> 🚀 **Dernier tag = `v0.2.9`** (à déployer sur Jelastic ; remplace v0.2.8 —
> corrige le drag-and-drop : seule la pièce flotte, plus la case). Ajoute depuis v0.2.7 :
> durcissement anti-scan, **détection tactiques revue** (sac juste, pas de prise de
> dame en 1, ≥ 3 coups, fin sur le coup décisif, deux camps analysés), **drag-and-drop**
> sur l'échiquier, **nav EN complétée** (Tournament prep + Blog), **Event JSON-LD**
> complété (alertes Search Console), et les **puzzles régénérés** (8, conformes).

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

- [ ] *(Option future)* **Ingestion runtime** `POST /api/admin/tactics` (Bearer `ADMIN_TOKEN`)
      pour pousser les puzzles **sans redeploy** — non fait, le modèle cron+rebuild suffit.
- [x] ~~Anti-spam renforcé~~ ✅ **fait maison** (jeton HMAC `/api/form-token` + filtre
      anti-liens + rate-limit resserré, en plus du honeypot) → voir `done.md`. Turnstile plus nécessaire.
- [x] ~~Prise de RDV~~ ✅ **fait maison** (`/reserver`, calcul du prix, e-mails) → voir `done.md`.
- [x] ~~Analytics privacy-first~~ ✅ **fait maison** (server-side, dans `/admin`) → voir `done.md`.
- [x] ~~hreflang articles blog~~ ✅ / ~~movetime tactiques~~ ✅ (400) → voir `done.md`.

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
- **Réservation `/reserver`** : nécessite `DB_PATH` + SMTP (sinon la résa est stockée mais
  les e-mails ne partent pas). Tarif via `HOURLY_RATE` (défaut 120). Les résas + l'analytics
  (fréquentation, top pages) apparaissent dans **`/admin`**. **À tester avant lancement** (envoie de vrais mails).
- Détail de tout le terminé → **[`done.md`](./done.md)**.
