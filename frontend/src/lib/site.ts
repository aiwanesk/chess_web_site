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
  twitter: '@iwanesko', // update or remove if no account

  person: {
    name: 'Alexandre Iwanesko',
    jobTitle: 'Coach d’échecs, Maître FIDE',
    honorific: 'Maître FIDE',
    description:
      'Maître FIDE et coach d’échecs à Genève, spécialisé dans la progression des adultes (1200–1900 Elo) et des adolescents en compétition.',
    // Entity disambiguation for GEO / knowledge graph — official FIDE profile only.
    // (lichess / chess.com intentionally not listed.) Add LinkedIn if desired.
    sameAs: [
      'https://ratings.fide.com/profile/682136',
    ],
  },

  // NAP — must match the Google Business Profile exactly.
  contact: {
    email: 'contact@iwanesko.ch',
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
    // Approx. Route de Florissant / Genève.
    geo: { lat: 46.1986, lng: 6.1668 },
  },
  areaServed: ['Genève', 'Vaud', 'Arc lémanique', 'France voisine'],
  priceRange: 'CHF',
} as const

export const absoluteUrl = (path: string): string => {
  if (path.startsWith('http')) return path
  return SITE.url + (path.startsWith('/') ? path : `/${path}`)
}
