# TODO — iwanesko.ch

Suivi des travaux restants. Coché = fait.

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
- [ ] Articles de blog : **à réécrire** (les 4 actuels sont provisoires).
- [ ] `/resultats` : témoignages **réels et attribuables** uniquement + vraies stats.
- [ ] `public/og/` : images Open Graph 1200×630 par type de page (PNG).
- [ ] Police **Inter variable** self-hostée dans `public/fonts/` + réactiver
      `preload` (index.html) et `@font-face` (styles.css).
- [ ] `/contact` : carte de Genève légère (chargée en différé).
- [ ] Dates réelles des stages (`pages/Stages.tsx`, `eventSchema`).

## En cours (itération design + contenu)
- Voir l'historique git. Money pages étoffées, design premium, blog + maillage.
