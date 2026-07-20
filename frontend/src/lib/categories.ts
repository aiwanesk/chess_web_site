/**
 * Blog categories — the user-facing browsable sections (findability), distinct
 * from `cluster` (which links an article to its money page for the SEO silo).
 *
 * - "progresser"       → guides/conseils (acquisition : chaque article pointe vers une money page)
 * - "carnet-de-tournoi"→ journal de compétition d'Alexandre (récits de première main)
 */
export interface Category {
  slug: string
  label: string // full title (archive page H1)
  short: string // nav / badge label
  description: string // meta description of the archive page
  intro: string // archive page hero subtitle
}

export const CATEGORIES: Category[] = [
  {
    slug: 'progresser',
    label: 'Progresser aux échecs',
    short: 'Progresser',
    description:
      'Guides et conseils d’un Maître FIDE pour progresser aux échecs : ouvertures, finales, tactique, préparation de tournoi et mental de compétition.',
    intro:
      'Des méthodes concrètes pour gagner des points Elo, par Alexandre Iwanesko, Maître FIDE.',
  },
  {
    slug: 'carnet-de-tournoi',
    label: 'Carnet de tournoi',
    short: 'Carnet de tournoi',
    description:
      'Le journal de compétition d’Alexandre Iwanesko, Maître FIDE : parties marquantes, décisions sous pression, coulisses et leçons de mes tournois.',
    intro:
      'Mes tournois de l’intérieur — parties marquantes, choix sous pression et ce que j’en retiens.',
  },
]

export const DEFAULT_CATEGORY = 'progresser'

export const CATEGORY_SLUGS = CATEGORIES.map((c) => c.slug)

export function getCategory(slug: string): Category | undefined {
  return CATEGORIES.find((c) => c.slug === slug)
}

/** Root-relative path of a category archive page. */
export const categoryPath = (slug: string): string => `/blog/categorie/${slug}`

/** Build-time list of category archive paths for the router's getStaticPaths. */
export const categoryStaticPaths = (): string[] => CATEGORY_SLUGS.map(categoryPath)
