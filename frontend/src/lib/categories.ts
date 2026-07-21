/**
 * Blog categories — the user-facing browsable sections (findability), distinct
 * from `cluster` (which links an article to its money page for the SEO silo).
 *
 * Bilingual: each category has a stable `key` (used in article front-matter,
 * language-independent) and a localized view per locale (FR default + EN).
 */
import type { Locale } from './i18n'

export interface CategoryView {
  slug: string // URL segment for this locale
  label: string // full title (archive page H1)
  short: string // nav / badge label
  description: string // meta description of the archive page
  intro: string // archive page hero subtitle
}

export interface Category extends CategoryView {
  key: string // stable id (== FR slug), referenced by article `category` front-matter
  en: CategoryView
}

export const CATEGORIES: Category[] = [
  {
    key: 'progresser',
    slug: 'progresser',
    label: 'Progresser aux échecs',
    short: 'Progresser',
    description:
      'Guides et conseils d’un Maître FIDE pour progresser aux échecs : ouvertures, finales, tactique, préparation de tournoi et mental de compétition.',
    intro:
      'Des méthodes concrètes pour gagner des points Elo, par Alexandre Iwanesko, Maître FIDE.',
    en: {
      slug: 'improve',
      label: 'Improve at chess',
      short: 'Improve',
      description:
        'Guides from a FIDE Master to improve at chess: openings, endgames, tactics, tournament preparation and competitive mindset.',
      intro: 'Concrete methods to gain Elo points, by Alexandre Iwanesko, FIDE Master.',
    },
  },
  {
    key: 'carnet-de-tournoi',
    slug: 'carnet-de-tournoi',
    label: 'Carnet de tournoi',
    short: 'Carnet de tournoi',
    description:
      'Le journal de compétition d’Alexandre Iwanesko, Maître FIDE : parties marquantes, décisions sous pression, coulisses et leçons de mes tournois.',
    intro:
      'Mes tournois de l’intérieur — parties marquantes, choix sous pression et ce que j’en retiens.',
    en: {
      slug: 'tournament-diary',
      label: 'Tournament diary',
      short: 'Tournament diary',
      description:
        'The competition diary of Alexandre Iwanesko, FIDE Master: memorable games, decisions under pressure, behind the scenes and lessons from my tournaments.',
      intro: 'My tournaments from the inside — key games, choices under pressure and what I take away.',
    },
  },
]

export const DEFAULT_CATEGORY = 'progresser'

export const CATEGORY_KEYS = CATEGORIES.map((c) => c.key)

/** The localized view of a category (FR fields live on the object itself). */
export function catView(c: Category, locale: Locale): CategoryView {
  return locale === 'en' ? c.en : c
}

/** Look up a category by its stable key (== the value in article front-matter). */
export function getCategory(key: string): Category | undefined {
  return CATEGORIES.find((c) => c.key === key)
}

/** Look up a category by the slug used in a given locale's URL. */
export function getCategoryBySlug(slug: string, locale: Locale): Category | undefined {
  return CATEGORIES.find((c) => catView(c, locale).slug === slug)
}

/** Root-relative path of a category archive page, per locale. */
export function categoryPath(key: string, locale: Locale): string {
  const c = getCategory(key)
  if (!c) return locale === 'en' ? '/en/blog' : '/blog'
  return locale === 'en' ? `/en/blog/category/${c.en.slug}` : `/blog/categorie/${c.slug}`
}

/** getStaticPaths helpers (one per locale). */
export const categoryStaticPaths = (): string[] => CATEGORIES.map((c) => `/blog/categorie/${c.slug}`)
export const categoryStaticPathsEN = (): string[] => CATEGORIES.map((c) => `/en/blog/category/${c.en.slug}`)
