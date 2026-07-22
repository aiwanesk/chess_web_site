import { Link } from 'vite-react-ssg'
import { Container } from './Container'
import { NewsletterSignup } from './NewsletterSignup'
import { SITE } from '../lib/site'
import { IconKnight, IconMail, IconPhone, IconPin } from './icons'
import { useLocale, t, type Locale } from '../lib/i18n'

type Column = { title: string; links: { to: string; label: string }[] }

const COLUMNS_FR: Column[] = [
  {
    title: 'Cours',
    links: [
      { to: '/cours-echecs-adultes-geneve', label: 'Cours adultes à Genève' },
      { to: '/cours-echecs-en-ligne', label: 'Cours en ligne' },
      { to: '/cours-echecs-groupe-geneve', label: 'Cours en groupe' },
      { to: '/cours-echecs-ados-competition', label: 'Ados en compétition' },
    ],
  },
  {
    title: 'Aller plus loin',
    links: [
      { to: '/preparation-tournoi-echecs', label: 'Préparation tournoi' },
      { to: '/stages-echecs-geneve', label: 'Stages à Genève' },
      { to: '/conferences-echecs-entreprise', label: 'Conférences entreprise' },
      { to: '/team-building-echecs-geneve', label: 'Team building' },
    ],
  },
  {
    title: 'Informations',
    links: [
      { to: '/a-propos', label: 'À propos' },
      { to: '/resultats', label: 'Résultats & avis' },
      { to: '/tarifs', label: 'Tarifs' },
      { to: '/blog', label: 'Blog échecs' },
      { to: '/blog/categorie/carnet-de-tournoi', label: 'Carnet de tournoi' },
      { to: '/reserver', label: 'Réserver un cours' },
      { to: '/contact', label: 'Contact' },
      { to: '/confidentialite', label: 'Confidentialité' },
    ],
  },
]

const COLUMNS_EN: Column[] = [
  {
    title: 'Lessons',
    links: [
      { to: '/en/adult-chess-lessons-geneva', label: 'Adult chess lessons in Geneva' },
      { to: '/en/online-chess-lessons', label: 'Online chess lessons' },
      { to: '/en/group-chess-lessons-geneva', label: 'Group chess lessons' },
      { to: '/en/junior-chess-coaching', label: 'Junior chess coaching' },
    ],
  },
  {
    title: 'Going further',
    links: [
      { to: '/en/tournament-preparation', label: 'Tournament preparation' },
      { to: '/en/chess-camps-geneva', label: 'Chess camps in Geneva' },
      { to: '/en/corporate-chess-talks', label: 'Corporate talks' },
      { to: '/en/chess-team-building-geneva', label: 'Team building' },
    ],
  },
  {
    title: 'Information',
    links: [
      { to: '/en/about', label: 'About' },
      { to: '/en/results', label: 'Results & reviews' },
      { to: '/en/pricing', label: 'Pricing' },
      { to: '/en/tactics', label: 'Weekly tactics' },
      { to: '/en/blog', label: 'Blog' },
      { to: '/en/book', label: 'Book a lesson' },
      { to: '/en/contact', label: 'Contact' },
      { to: '/en/privacy', label: 'Privacy' },
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
                <li key={l.to}>
                  <Link to={l.to} className="text-ink-300 transition-colors hover:text-gold-400">
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
