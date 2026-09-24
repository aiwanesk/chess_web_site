/**
 * schema.org / JSON-LD builders. Each returns a plain object that the Seo
 * component serialises into a <script type="application/ld+json"> tag in the
 * server-rendered HTML — exactly what search engines and generative engines
 * parse to attribute expertise (E-E-A-T).
 */
import { SITE, absoluteUrl } from './site'

type JsonLd = Record<string, unknown>

const PERSON_ID = `${SITE.url}/#person`
const BUSINESS_ID = `${SITE.url}/#business`

/**
 * Références vers les deux entités du site, AVEC leur nom.
 *
 * Les nœuds complets (Person, ProfessionalService) ne sont émis que sur
 * l'accueil et /a-propos. Ailleurs, un `{ '@id': … }` seul est une référence
 * pendante : un analyseur qui lit la page isolément voit un fournisseur, un
 * intervenant ou un auteur sans nom. Le nom voyage donc avec la référence, et
 * l'`@id` reste pour que tout se réconcilie à l'échelle du site.
 *
 * articleSchema appliquait déjà cette règle pour lui seul ; les pages argent, le
 * calendrier et /contact émettaient encore quarante références nues.
 * scripts/check-seo.mjs échoue maintenant dessus.
 */
const personRef = () => ({ '@type': 'Person', '@id': PERSON_ID, name: SITE.person.name, url: SITE.url })
const businessRef = () => ({
  '@type': 'ProfessionalService',
  '@id': BUSINESS_ID,
  name: SITE.name,
  url: SITE.url,
})

/** The coach as a stable Person entity. Emitted on home + /a-propos. */
export function personSchema(): JsonLd {
  return {
    '@context': 'https://schema.org',
    '@type': 'Person',
    '@id': PERSON_ID,
    name: SITE.person.name,
    jobTitle: SITE.person.jobTitle,
    description: SITE.person.description,
    url: SITE.url,
    knowsAbout: ['Échecs', 'Stratégie', 'Préparation de tournoi', 'Ouvertures', 'Finales'],
    hasCredential: {
      '@type': 'EducationalOccupationalCredential',
      credentialCategory: 'FIDE title',
      name: 'Maître FIDE (FIDE Master)',
    },
    sameAs: SITE.person.sameAs,
    ...(SITE.person.image ? { image: absoluteUrl(SITE.person.image) } : {}),
  }
}

/** The coaching practice as a LocalBusiness for local SEO. */
export function localBusinessSchema(): JsonLd {
  return {
    '@context': 'https://schema.org',
    '@type': 'ProfessionalService',
    '@id': BUSINESS_ID,
    name: SITE.name,
    url: SITE.url,
    image: absoluteUrl(SITE.defaultOgImage),
    email: SITE.contact.email,
    telephone: SITE.contact.phone,
    priceRange: SITE.priceRange,
    founder: personRef(),
    employee: personRef(),
    address: {
      '@type': 'PostalAddress',
      streetAddress: SITE.address.street,
      addressLocality: SITE.address.locality,
      addressRegion: SITE.address.region,
      postalCode: SITE.address.postalCode,
      addressCountry: SITE.address.country,
    },
    geo: {
      '@type': 'GeoCoordinates',
      latitude: SITE.address.geo.lat,
      longitude: SITE.address.geo.lng,
    },
    areaServed: SITE.areaServed.map((name) => ({ '@type': 'Place', name })),
    availableLanguage: ['fr', 'en'],
    knowsLanguage: ['fr', 'en'],
  }
}

export interface OfferingInput {
  name: string
  description: string
  url: string
  price?: number // CHF, per unit
  priceUnit?: string // e.g. 'la séance (60 min)'
  courseMode?: 'onsite' | 'online' | 'blended'
  /**
   * Durée d'une session, en durée ISO 8601 (`PT1H`). Elle valait `PT1H` en dur
   * pour tout le monde : le stage « 2 à 5 jours » se déclarait donc comme une
   * heure de cours, et son tarif forfaitaire se lisait comme un tarif horaire.
   * Facultative — sans elle, rien n'est affirmé sur la durée.
   */
  courseWorkload?: string
  /**
   * `course` (défaut) pour ce qui s'enseigne à un élève sur la durée, `service`
   * pour une prestation ponctuelle vendue à une entreprise. La conférence et le
   * team building sortaient en `Course` : ni inscription, ni programme, ni
   * élève — juste une intervention facturée. Un `Course` qui n'enseigne rien
   * est une donnée structurée fausse, et Google le teste sur ces signaux-là.
   */
  kind?: 'course' | 'service'
}

/**
 * L'offre d'une page argent : un `Course` pour un cours, un `Service` pour une
 * prestation en entreprise, plus l'`Offer` quand un prix public existe (les
 * pages entreprise sont sur devis : pas de prix, donc pas d'offre inventée).
 */
export function offeringSchema(c: OfferingInput): JsonLd {
  const modeMap = { onsite: 'Onsite', online: 'Online', blended: 'Blended' } as const
  const url = absoluteUrl(c.url)
  // La locale se lit dans l'URL : les pages EN se déclaraient toutes en `fr`,
  // exactement le bug déjà corrigé sur les articles.
  const inLanguage = c.url.startsWith('/en') ? 'en' : 'fr'
  const offers = c.price
    ? {
        offers: {
          '@type': 'Offer',
          price: c.price,
          priceCurrency: 'CHF',
          availability: 'https://schema.org/InStock',
          url,
          ...(c.priceUnit ? { description: `Tarif ${c.priceUnit}` } : {}),
        },
      }
    : {}

  if (c.kind === 'service') {
    return {
      '@context': 'https://schema.org',
      '@type': 'Service',
      name: c.name,
      serviceType: c.name,
      description: c.description,
      url,
      provider: businessRef(),
      areaServed: SITE.areaServed.map((name) => ({ '@type': 'Place', name })),
      availableChannel: { '@type': 'ServiceChannel', serviceUrl: url },
      ...offers,
    }
  }

  return {
    '@context': 'https://schema.org',
    '@type': 'Course',
    name: c.name,
    description: c.description,
    url,
    inLanguage,
    provider: businessRef(),
    hasCourseInstance: {
      '@type': 'CourseInstance',
      courseMode: modeMap[c.courseMode ?? 'blended'],
      ...(c.courseWorkload ? { courseWorkload: c.courseWorkload } : {}),
      instructor: personRef(),
    },
    ...offers,
  }
}

export interface EventInput {
  name: string
  description: string
  url: string
  startDate: string // ISO
  endDate?: string
  price?: number
  image?: string // defaults to the site OG image
  validFrom?: string // offer valid-from (ISO); defaults to Jan 1 of the event year
}

export function eventSchema(e: EventInput): JsonLd {
  const image = absoluteUrl(e.image ?? SITE.defaultOgImage)
  const validFrom = e.validFrom ?? `${e.startDate.slice(0, 4)}-01-01`
  return {
    '@context': 'https://schema.org',
    '@type': 'Event',
    name: e.name,
    description: e.description,
    image: [image],
    url: absoluteUrl(e.url),
    startDate: e.startDate,
    ...(e.endDate ? { endDate: e.endDate } : {}),
    eventStatus: 'https://schema.org/EventScheduled',
    eventAttendanceMode: 'https://schema.org/OfflineEventAttendanceMode',
    location: {
      '@type': 'Place',
      name: 'Genève',
      address: { '@type': 'PostalAddress', addressLocality: 'Genève', addressCountry: 'CH' },
    },
    organizer: businessRef(),
    ...(e.price != null
      ? {
          offers: {
            '@type': 'Offer',
            price: e.price,
            priceCurrency: 'CHF',
            availability: 'https://schema.org/InStock',
            validFrom,
            url: absoluteUrl(e.url),
          },
        }
      : {}),
  }
}

export interface TournamentEventInput {
  name: string
  startDate: string // YYYY-MM-DD
  endDate: string
  location?: string // free text, e.g. "Badalona, Espagne"
  url: string // where the visitor reads about it (calendar page, or its diary)
  description?: string
}

/**
 * A tournament Alexandre is competing in. `performer` points at the Person
 * entity, so the graph states that the FIDE Master is playing it — first-hand
 * competitive experience, which is exactly the E-E-A-T signal a coaching site
 * wants. Emit for upcoming/ongoing events only: a past event carries no value
 * as a rich result.
 */
export function sportsEventSchema(e: TournamentEventInput): JsonLd {
  return {
    '@context': 'https://schema.org',
    '@type': 'SportsEvent',
    name: e.name,
    ...(e.description ? { description: e.description } : {}),
    startDate: e.startDate,
    endDate: e.endDate,
    url: absoluteUrl(e.url),
    eventStatus: 'https://schema.org/EventScheduled',
    eventAttendanceMode: 'https://schema.org/OfflineEventAttendanceMode',
    sport: 'Chess',
    location: e.location
      ? { '@type': 'Place', name: e.location, address: e.location }
      : { '@type': 'Place', name: 'Europe' },
    performer: personRef(),
  }
}

export interface FaqItem {
  question: string
  answer: string
}

export function faqSchema(items: FaqItem[]): JsonLd {
  return {
    '@context': 'https://schema.org',
    '@type': 'FAQPage',
    mainEntity: items.map((it) => ({
      '@type': 'Question',
      name: it.question,
      acceptedAnswer: { '@type': 'Answer', text: it.answer },
    })),
  }
}

export interface Crumb {
  name: string
  path: string
}

export function breadcrumbSchema(crumbs: Crumb[]): JsonLd {
  return {
    '@context': 'https://schema.org',
    '@type': 'BreadcrumbList',
    itemListElement: crumbs.map((c, i) => ({
      '@type': 'ListItem',
      position: i + 1,
      name: c.name,
      item: absoluteUrl(c.path),
    })),
  }
}

export interface ReviewInput {
  author: string
  body: string
  rating: number
}

/**
 * ⚠️ Only ever call this with reviews a student actually wrote and rated.
 * Ratings nobody gave are fabricated structured data: a manual-action risk with
 * Google, and misleading advertising under Swiss LCD art. 3 / the EU Omnibus
 * directive. Currently unused — /resultats presents coached cases instead, with
 * no stars. Wire it back in the day real written reviews come in.
 */
export function aggregateRatingSchema(reviews: ReviewInput[]): JsonLd {
  const value = reviews.reduce((s, r) => s + r.rating, 0) / reviews.length
  return {
    '@context': 'https://schema.org',
    '@type': 'ProfessionalService',
    '@id': BUSINESS_ID,
    name: SITE.name,
    aggregateRating: {
      '@type': 'AggregateRating',
      ratingValue: value.toFixed(1),
      reviewCount: reviews.length,
      bestRating: 5,
    },
    review: reviews.map((r) => ({
      '@type': 'Review',
      author: { '@type': 'Person', name: r.author },
      reviewRating: { '@type': 'Rating', ratingValue: r.rating, bestRating: 5 },
      reviewBody: r.body,
    })),
  }
}

export interface ArticleInput {
  title: string
  description: string
  url: string
  datePublished: string
  dateModified?: string
  image?: string
  /** Page language. Sans elle, un article EN se declarait francais. */
  locale?: 'fr' | 'en'
  /** Rubrique lisible, ex. « Progresser ». */
  section?: string
  /** Temps de lecture affiche sur la page, en minutes. */
  readingMinutes?: number
}

/**
 * Un article de blog.
 *
 * Deux choses valent d'etre expliquees, parce qu'elles ne se voient pas :
 *  - une page d'article ne porte AUCUN noeud Person ni Organization (ils ne
 *    sont emis que sur l'accueil et /a-propos). Un `@id` seul y serait donc une
 *    reference pendante : un analyseur qui lit cette page isolement verrait un
 *    auteur sans nom. Le nom voyage avec la reference, l'`@id` reste pour que
 *    les deux se reconcilient a l'echelle du site.
 *  - `image` n'est jamais omise : la carte de l'article s'il en a une, celle du
 *    site sinon — comme og:image, qui a toujours eu ce repli.
 */
export function articleSchema(a: ArticleInput): JsonLd {
  const url = absoluteUrl(a.url)
  return {
    '@context': 'https://schema.org',
    '@type': 'BlogPosting',
    headline: a.title,
    description: a.description,
    url,
    mainEntityOfPage: { '@type': 'WebPage', '@id': url },
    datePublished: a.datePublished,
    dateModified: a.dateModified ?? a.datePublished,
    inLanguage: a.locale === 'en' ? 'en' : 'fr',
    author: personRef(),
    publisher: {
      ...businessRef(),
      logo: {
        '@type': 'ImageObject',
        url: absoluteUrl(SITE.logo.url),
        width: SITE.logo.width,
        height: SITE.logo.height,
      },
    },
    image: absoluteUrl(a.image ?? SITE.defaultOgImage),
    ...(a.section ? { articleSection: a.section } : {}),
    ...(a.readingMinutes ? { timeRequired: `PT${a.readingMinutes}M` } : {}),
  }
}
