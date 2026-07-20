import { createContext, useContext, type ReactNode } from 'react'

/**
 * Internationalisation. French is the default at the root ("/"); English lives
 * under "/en/" with its OWN English slugs (best for EN SEO). The PAGES registry
 * maps a page key to its FR and EN paths and drives: routing, the language
 * switcher, hreflang alternates and the sitemap parity.
 */
export type Locale = 'fr' | 'en'
export const DEFAULT_LOCALE: Locale = 'fr'
export const LOCALES: Locale[] = ['fr', 'en']

export function localeFromPath(pathname: string): Locale {
  return pathname === '/en' || pathname.startsWith('/en/') ? 'en' : 'fr'
}

const LocaleCtx = createContext<Locale>(DEFAULT_LOCALE)

export function LocaleProvider({ locale, children }: { locale: Locale; children: ReactNode }) {
  return <LocaleCtx.Provider value={locale}>{children}</LocaleCtx.Provider>
}

export const useLocale = (): Locale => useContext(LocaleCtx)

// --- Page path registry (FR ↔ EN) -----------------------------------------
// Add an entry here as each page gains its English translation.

export interface PagePaths {
  fr: string
  en: string
}

export const PAGES = {
  home: { fr: '/', en: '/en' },
  coursAdultes: { fr: '/cours-echecs-adultes-geneve', en: '/en/adult-chess-lessons-geneva' },
  tarifs: { fr: '/tarifs', en: '/en/pricing' },
  contact: { fr: '/contact', en: '/en/contact' },
} satisfies Record<string, PagePaths>

export type PageKey = keyof typeof PAGES

export function pathFor(key: PageKey, locale: Locale): string {
  return PAGES[key][locale]
}

/** The page key whose FR or EN path matches this pathname. */
export function keyForPath(pathname: string): PageKey | undefined {
  for (const k of Object.keys(PAGES) as PageKey[]) {
    const p = PAGES[k]
    if (p.fr === pathname || p.en === pathname) return k
  }
  return undefined
}

/** The counterpart URL in the other locale, if this page is translated. */
export function altPath(pathname: string): string | undefined {
  const key = keyForPath(pathname)
  if (!key) return undefined
  const other: Locale = localeFromPath(pathname) === 'fr' ? 'en' : 'fr'
  return PAGES[key][other]
}

// --- Shared UI strings ------------------------------------------------------

type Dict = {
  nav: { adultes: string; tournoi: string; enligne: string; tarifs: string; resultats: string; blog: string; tactiques: string }
  reserve: string
  reserveFirst: string
  viewPricing: string
  openMenu: string
  closeMenu: string
  skipToContent: string
  breadcrumbHome: string
  fideTitle: string
  footerRole: string
  footerZones: string
  footerRights: string
  faqTitle: string
  langLabel: string
  banner: string
  money: {
    methodEyebrow: string
    benefitsTitle: string
    goingFurther: string
    readEyebrow: string
    relatedTitle: string
    closingTitle: string
    closingSubtitle: string
    readArticle: string
  }
}

export const UI: Record<Locale, Dict> = {
  fr: {
    nav: { adultes: 'Cours adultes', tournoi: 'Préparation tournoi', enligne: 'En ligne', tarifs: 'Tarifs', resultats: 'Résultats', blog: 'Blog', tactiques: 'Tactiques' },
    reserve: 'Réserver un cours',
    reserveFirst: 'Réserver un premier cours',
    viewPricing: 'Voir les tarifs',
    openMenu: 'Ouvrir le menu',
    closeMenu: 'Fermer le menu',
    skipToContent: 'Aller au contenu',
    breadcrumbHome: 'Accueil',
    fideTitle: 'Maître FIDE',
    footerRole: 'Maître FIDE · Coach d’échecs',
    footerZones: 'Zones',
    footerRights: 'Tous droits réservés.',
    faqTitle: 'Questions fréquentes',
    langLabel: 'EN',
    banner: 'Site en construction — certaines pages arrivent bientôt.',
    money: {
      methodEyebrow: 'La méthode',
      benefitsTitle: 'Ce que vous obtenez',
      goingFurther: 'Pour aller plus loin',
      readEyebrow: 'À lire',
      relatedTitle: 'Articles liés',
      closingTitle: 'Prêt à progresser ?',
      closingSubtitle: 'Un premier échange pour définir vos objectifs et le format qui vous convient.',
      readArticle: 'Lire l’article',
    },
  },
  en: {
    nav: { adultes: 'Adult lessons', tournoi: 'Tournament prep', enligne: 'Online', tarifs: 'Pricing', resultats: 'Results', blog: 'Blog', tactiques: 'Tactics' },
    reserve: 'Book a lesson',
    reserveFirst: 'Book a first lesson',
    viewPricing: 'View pricing',
    openMenu: 'Open menu',
    closeMenu: 'Close menu',
    skipToContent: 'Skip to content',
    breadcrumbHome: 'Home',
    fideTitle: 'FIDE Master',
    footerRole: 'FIDE Master · Chess coach',
    footerZones: 'Areas',
    footerRights: 'All rights reserved.',
    faqTitle: 'Frequently asked questions',
    langLabel: 'FR',
    banner: 'Site under construction — some pages are coming soon.',
    money: {
      methodEyebrow: 'The method',
      benefitsTitle: 'What you get',
      goingFurther: 'Going further',
      readEyebrow: 'Read',
      relatedTitle: 'Related articles',
      closingTitle: 'Ready to improve?',
      closingSubtitle: 'A first conversation to define your goals and the format that suits you.',
      readArticle: 'Read the article',
    },
  },
}

export const t = (locale: Locale): Dict => UI[locale]

/** Locale-aware home path. */
export const homePath = (locale: Locale): string => (locale === 'en' ? '/en' : '/')
