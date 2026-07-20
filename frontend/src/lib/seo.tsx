import { Head } from 'vite-react-ssg'
import { SITE, absoluteUrl } from './site'

export interface SeoProps {
  /** ≤ 60 chars ideally. The site name is appended automatically unless isHome. */
  title: string
  /** ≤ 155 chars. */
  description: string
  /** Root-relative path of this page, e.g. "/tarifs". Used for canonical + OG. */
  path: string
  /** Social share image (root-relative or absolute). */
  image?: string
  /** "website" (default) or "article". */
  ogType?: 'website' | 'article'
  /** Set on pages that must not be indexed (none by default). */
  noindex?: boolean
  isHome?: boolean
  /** One or more JSON-LD objects injected as <script type="application/ld+json">. */
  jsonLd?: Array<Record<string, unknown>>
}

/**
 * Emits every per-page head tag into the server-rendered HTML. Because
 * vite-react-ssg pre-renders each route, these tags are present at the first
 * byte with no JS execution — verifiable with `curl`.
 */
export function Seo({
  title,
  description,
  path,
  image = SITE.defaultOgImage,
  ogType = 'website',
  noindex = false,
  isHome = false,
  jsonLd = [],
}: SeoProps) {
  const fullTitle = isHome ? title : `${title} | ${SITE.person.name}`
  const canonical = absoluteUrl(path)
  const ogImage = absoluteUrl(image)

  return (
    <Head>
      <html lang={SITE.lang} />
      <title>{fullTitle}</title>
      <meta name="description" content={description} />
      <link rel="canonical" href={canonical} />
      <meta name="robots" content={noindex ? 'noindex,nofollow' : 'index,follow,max-image-preview:large'} />

      {/* Open Graph */}
      <meta property="og:type" content={ogType} />
      <meta property="og:site_name" content={SITE.shortName} />
      <meta property="og:locale" content={SITE.locale} />
      <meta property="og:title" content={fullTitle} />
      <meta property="og:description" content={description} />
      <meta property="og:url" content={canonical} />
      <meta property="og:image" content={ogImage} />

      {/* Twitter / X */}
      <meta name="twitter:card" content="summary_large_image" />
      <meta name="twitter:title" content={fullTitle} />
      <meta name="twitter:description" content={description} />
      <meta name="twitter:image" content={ogImage} />
      {SITE.twitter ? <meta name="twitter:site" content={SITE.twitter} /> : null}

      {/*
        hreflang — no-op while the site is FR only. When the EN v2 ships, add:
        <link rel="alternate" hrefLang="en" href={absoluteUrl('/en' + path)} />
      */}
      <link rel="alternate" hrefLang="fr" href={canonical} />
      <link rel="alternate" hrefLang="x-default" href={canonical} />

      {jsonLd.map((data, i) => (
        <script key={i} type="application/ld+json">
          {JSON.stringify(data)}
        </script>
      ))}
    </Head>
  )
}
