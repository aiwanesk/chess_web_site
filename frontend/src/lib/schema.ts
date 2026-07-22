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
    founder: { '@id': PERSON_ID },
    employee: { '@id': PERSON_ID },
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

export interface CourseInput {
  name: string
  description: string
  url: string
  price?: number // CHF, per unit
  priceUnit?: string // e.g. 'la séance (60 min)'
  courseMode?: 'onsite' | 'online' | 'blended'
}

/** A Course + Offer with a CHF price. Emitted on each money page. */
export function courseSchema(c: CourseInput): JsonLd {
  const modeMap = { onsite: 'Onsite', online: 'Online', blended: 'Blended' } as const
  const offers = c.price
    ? {
        offers: {
          '@type': 'Offer',
          price: c.price,
          priceCurrency: 'CHF',
          availability: 'https://schema.org/InStock',
          url: absoluteUrl(c.url),
          ...(c.priceUnit ? { description: `Tarif ${c.priceUnit}` } : {}),
        },
      }
    : {}
  return {
    '@context': 'https://schema.org',
    '@type': 'Course',
    name: c.name,
    description: c.description,
    url: absoluteUrl(c.url),
    inLanguage: 'fr',
    provider: { '@id': BUSINESS_ID },
    hasCourseInstance: {
      '@type': 'CourseInstance',
      courseMode: modeMap[c.courseMode ?? 'blended'],
      courseWorkload: 'PT1H',
      instructor: { '@id': PERSON_ID },
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
    organizer: { '@type': 'Organization', '@id': BUSINESS_ID, name: SITE.person.name, url: SITE.url },
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
}

export function articleSchema(a: ArticleInput): JsonLd {
  return {
    '@context': 'https://schema.org',
    '@type': 'BlogPosting',
    headline: a.title,
    description: a.description,
    url: absoluteUrl(a.url),
    mainEntityOfPage: absoluteUrl(a.url),
    datePublished: a.datePublished,
    dateModified: a.dateModified ?? a.datePublished,
    inLanguage: 'fr',
    author: { '@id': PERSON_ID },
    publisher: { '@id': BUSINESS_ID },
    ...(a.image ? { image: absoluteUrl(a.image) } : {}),
  }
}
