/**
 * Single source of truth for brand / NAP / entity data. Reused by SEO meta,
 * JSON-LD builders, header and footer so the E-E-A-T signals stay consistent
 * everywhere (a key GEO/local-SEO requirement).
 *
 * ⚠️ Replace the placeholder phone/address/sameAs values with the real ones,
 * kept identical to the Google Business Profile.
 */
export const SITE = {
  name: 'Alexandre Iwanesko — Coach d’échecs à Genève',
  shortName: 'Alexandre Iwanesko échecs',
  url: 'https://iwanesko.ch',
  locale: 'fr_CH',
  lang: 'fr',
  defaultOgImage: '/og/default.png',
  // Logo carré, pour les endroits qui attendent un logo et pas une carte de
  // partage : publisher.logo en JSON-LD réclame une image (Google demande au
  // moins 112 px de haut), pas un visuel 1200×630 qui s'y affiche rogné ou
  // refusé. Régénéré par scripts/gen-brand-assets.mjs depuis assets/logo-source.png.
  logo: { url: '/logo-512.png', width: 512, height: 512 },
  // No X/Twitter account for now, so no `twitter:site` tag: pointing it at a
  // handle nobody owns is a broken attribution (and an open door to someone
  // else claiming it). Add the handle back here and restore the tag in seo.tsx
  // if an account is ever created. The other twitter:* tags stay — they drive
  // the link preview and need no account.

  person: {
    name: 'Alexandre Iwanesko',
    jobTitle: 'Coach d’échecs, Maître FIDE',
    honorific: 'Maître FIDE',
    description:
      'Maître FIDE et coach d’échecs à Genève, spécialisé dans la progression des adultes (1200–2200 Elo) et des adolescents en compétition.',
    // Désambiguïsation de l'entité pour le graphe de connaissances : plus il y a
    // de profils officiels concordants, plus un moteur peut relier « Alexandre
    // Iwanesko » à une seule personne réelle. Uniquement des fiches vérifiables.
    //
    // TODO(alex) — trois pistes manquent, et je ne peux pas les deviner :
    //   1. chess-results.com : il n'y a pas d'URL de joueur stable, seulement des
    //      pages de tournoi. Envoie le lien de ta fiche s'il en existe une.
    //   2. Lichess et/ou Chess.com : ils étaient volontairement exclus jusqu'ici
    //      (ce ne sont pas des fédérations). À rétablir seulement si tu veux
    //      assumer ces comptes publiquement — donne les pseudos exacts.
    //   3. Fiche Google Business : l'URL n'existera qu'une fois la fiche créée
    //      (voir le point SEO local de l'audit).
    sameAs: [
      'https://ratings.fide.com/profile/682136',
      'https://www.echecs.asso.fr/FicheJoueur.aspx?Id=335623',
      'https://www.365chess.com/players/Alexandre_Iwanesko',
      'https://ch.linkedin.com/in/alexandre-iwanesko-720593144',
    ],
    // TODO(alex) — Person.image : il n'y a aucun portrait dans le dépôt. Dépose
    // une photo dans public/ (paysage ou carré, ≥ 800 px, toi reconnaissable) et
    // renseigne son chemin ici ; personSchema l'émettra automatiquement, et
    // /a-propos l'affichera. Un logo ne peut pas servir d'image de personne.
    image: undefined as string | undefined,
  },

  // NAP — must match the Google Business Profile exactly.
  contact: {
    email: 'alexandre@iwanesko.ch',
    phone: '+41 78 783 56 89',
    phoneHref: 'tel:+41787835689',
  },
  // Adresse pro (domiciliation Swiss Tax Horizon). NPA à confirmer.
  address: {
    street: 'Route de Florissant 2',
    locality: 'Genève',
    region: 'GE',
    postalCode: '1206',
    country: 'CH',
    // Route de Florissant 2, 1206 Genève (coordonnées précises).
    geo: { lat: 46.196817, lng: 6.153666 },
  },
  areaServed: ['Genève', 'Vaud', 'Arc lémanique', 'France voisine'],
  // Une fourchette, pas un code de devise : « CHF » seul n'apprend rien à
  // personne. Les bornes viennent de Tarifs.tsx — 70 CHF le cours en groupe,
  // 120 CHF la séance individuelle. À tenir à jour avec la page tarifs.
  priceRange: 'CHF 70–120',
} as const

export const absoluteUrl = (path: string): string => {
  if (path.startsWith('http')) return path
  return SITE.url + (path.startsWith('/') ? path : `/${path}`)
}
