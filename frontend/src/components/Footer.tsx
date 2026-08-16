import { Link } from 'vite-react-ssg'
import { Container } from './Container'
import { NewsletterSignup } from './NewsletterSignup'
import { SITE } from '../lib/site'
import { IconKnight, IconMail, IconPhone, IconPin } from './icons'
import { categoryPath } from '../lib/categories'
import { useLocale, t, pathFor, type Locale, type PageKey } from '../lib/i18n'

/**
 * A footer link names a PAGE, never a URL — pathFor() resolves it per locale, so
 * a French URL can never surface in the English footer. `cat` covers the blog
 * category archives, which live outside the PAGES registry.
 */
type FooterLink = { page: PageKey; label: string } | { cat: string; label: string }
type Column = { title: string; links: FooterLink[] }

const linkPath = (l: FooterLink, locale: Locale): string =>
  'page' in l ? pathFor(l.page, locale) : categoryPath(l.cat, locale)

const COLUMNS_FR: Column[] = [
  {
    title: 'Cours',
    links: [
      { page: 'coursAdultes', label: 'Cours adultes à Genève' },
      { page: 'coursEnLigne', label: 'Cours en ligne' },
      { page: 'coursGroupe', label: 'Cours en groupe' },
      { page: 'coursAdos', label: 'Ados en compétition' },
    ],
  },
  {
    title: 'Aller plus loin',
    links: [
      { page: 'preparationTournoi', label: 'Préparation tournoi' },
      { page: 'stages', label: 'Stages à Genève' },
      { page: 'conferences', label: 'Conférences entreprise' },
      { page: 'teamBuilding', label: 'Team building' },
    ],
  },
  {
    title: 'Informations',
    links: [
      { page: 'apropos', label: 'À propos' },
      { page: 'resultats', label: 'Résultats & avis' },
      { page: 'tarifs', label: 'Tarifs' },
      { page: 'blog', label: 'Blog échecs' },
      { cat: 'carnet-de-tournoi', label: 'Carnet de tournoi' },
      { page: 'calendrier', label: 'Calendrier des tournois' },
      { page: 'reserver', label: 'Réserver un cours' },
      { page: 'contact', label: 'Contact' },
      { page: 'confidentialite', label: 'Confidentialité' },
    ],
  },
]

const COLUMNS_EN: Column[] = [
  {
    title: 'Lessons',
    links: [
      { page: 'coursAdultes', label: 'Adult chess lessons in Geneva' },
      { page: 'coursEnLigne', label: 'Online chess lessons' },
      { page: 'coursGroupe', label: 'Group chess lessons' },
      { page: 'coursAdos', label: 'Junior chess coaching' },
    ],
  },
  {
    title: 'Going further',
    links: [
      { page: 'preparationTournoi', label: 'Tournament preparation' },
      { page: 'stages', label: 'Chess camps in Geneva' },
      { page: 'conferences', label: 'Corporate talks' },
      { page: 'teamBuilding', label: 'Team building' },
    ],
  },
  {
    title: 'Information',
    links: [
      { page: 'apropos', label: 'About' },
      { page: 'resultats', label: 'Results & reviews' },
      { page: 'tarifs', label: 'Pricing' },
      { page: 'tactiques', label: 'Weekly tactics' },
      { page: 'blog', label: 'Blog' },
      { page: 'calendrier', label: 'Tournament calendar' },
      { page: 'reserver', label: 'Book a lesson' },
      { page: 'contact', label: 'Contact' },
      { page: 'confidentialite', label: 'Privacy' },
    ],
  },
]

export function Footer() {
  const locale = useLocale()
  const s = t(locale)
  const columns = locale === 'en' ? COLUMNS_EN : COLUMNS_FR
  const tagline: Record<Locale, string> = {
    fr: 'Genève · Arc lémanique · Cours en ligne',
    en: 'Geneva · Lake Geneva region · Online lessons',
  }

  return (
    <footer className="mt-8 border-t-2 border-gold-500/60 bg-ink-950 text-ink-300">
      <div className="border-b border-white/10">
        <Container className="py-12">
          <div className="max-w-md">
            <NewsletterSignup />
          </div>
        </Container>
      </div>

      <Container className="grid gap-12 py-16 md:grid-cols-4">
        {/* NAP block — must stay consistent with the Google Business Profile. */}
        <div>
          <p className="flex items-center gap-2.5 text-lg font-bold text-white">
            <span aria-hidden className="flex h-9 w-9 items-center justify-center rounded-full bg-white/5 text-gold-500 ring-1 ring-inset ring-gold-500/30">
              <IconKnight size={19} />
            </span>
            <span className="font-display tracking-tight">Alexandre Iwanesko</span>
          </p>
          <p className="mt-3 text-sm font-medium text-gold-400">{s.footerRole}</p>
          <address className="mt-5 space-y-2.5 not-italic text-sm leading-relaxed">
            <span className="flex items-start gap-2.5">
              <IconPin size={16} className="mt-0.5 flex-none text-gold-500/80" />
              <span>
                {SITE.address.street}
                <br />
                {SITE.address.postalCode} {SITE.address.locality}, {locale === 'en' ? 'Switzerland' : 'Suisse'}
              </span>
            </span>
            <a href={SITE.contact.phoneHref} className="flex items-center gap-2.5 transition-colors hover:text-white">
              <IconPhone size={16} className="flex-none text-gold-500/80" />
              {SITE.contact.phone}
            </a>
            <a href={`mailto:${SITE.contact.email}`} className="flex items-center gap-2.5 transition-colors hover:text-white">
              <IconMail size={16} className="flex-none text-gold-500/80" />
              {SITE.contact.email}
            </a>
          </address>
          <p className="mt-4 text-xs leading-relaxed text-ink-400">
            {s.footerZones} : {SITE.areaServed.join(', ')}.
          </p>
        </div>

        {columns.map((col) => (
          <nav key={col.title} aria-label={col.title}>
            <p className="flex items-center gap-2 text-xs font-semibold uppercase tracking-[0.16em] text-white">
              <span aria-hidden className="h-px w-4 bg-gold-500" />
              {col.title}
            </p>
            <ul className="mt-5 space-y-2.5 text-sm">
              {col.links.map((l) => (
                <li key={l.label}>
                  <Link to={linkPath(l, locale)} className="text-ink-300 transition-colors hover:text-gold-400">
                    {l.label}
                  </Link>
                </li>
              ))}
            </ul>
          </nav>
        ))}
      </Container>

      <div className="border-t border-white/10">
        <Container className="flex flex-col items-center justify-between gap-2 py-6 text-xs text-ink-400 sm:flex-row">
          <p>© {new Date().getFullYear()} Alexandre Iwanesko — {s.footerRights}</p>
          <p>{tagline[locale]}</p>
        </Container>
      </div>
    </footer>
  )
}
