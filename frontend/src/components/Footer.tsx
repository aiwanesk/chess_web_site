import { Link } from 'vite-react-ssg'
import { Container } from './Container'
import { SITE } from '../lib/site'
import { IconKnight, IconMail, IconPhone, IconPin } from './icons'

const COLUMNS = [
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
      { to: '/contact', label: 'Contact' },
    ],
  },
]

export function Footer() {
  return (
    <footer className="mt-8 border-t-2 border-gold-500/60 bg-ink-950 text-ink-300">
      <Container className="grid gap-12 py-16 md:grid-cols-4">
        {/* NAP block — must stay consistent with the Google Business Profile. */}
        <div>
          <p className="flex items-center gap-2.5 text-lg font-bold text-white">
            <span
              aria-hidden
              className="flex h-9 w-9 items-center justify-center rounded-full bg-white/5 text-gold-500 ring-1 ring-inset ring-gold-500/30"
            >
              <IconKnight size={19} />
            </span>
            <span className="font-display tracking-tight">Alexandre Iwanesko</span>
          </p>
          <p className="mt-3 text-sm font-medium text-gold-400">Maître FIDE · Coach d’échecs</p>
          <address className="mt-5 space-y-2.5 not-italic text-sm leading-relaxed">
            <span className="flex items-start gap-2.5">
              <IconPin size={16} className="mt-0.5 flex-none text-gold-500/80" />
              <span>
                {SITE.address.street}
                <br />
                {SITE.address.postalCode} {SITE.address.locality}, Suisse
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
            Zones : {SITE.areaServed.join(', ')}.
          </p>
        </div>

        {COLUMNS.map((col) => (
          <nav key={col.title} aria-label={col.title}>
            <p className="flex items-center gap-2 text-xs font-semibold uppercase tracking-[0.16em] text-white">
              <span aria-hidden className="h-px w-4 bg-gold-500" />
              {col.title}
            </p>
            <ul className="mt-5 space-y-2.5 text-sm">
              {col.links.map((l) => (
                <li key={l.to}>
                  <Link
                    to={l.to}
                    className="text-ink-300 transition-colors hover:text-gold-400"
                  >
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
          <p>© {new Date().getFullYear()} Alexandre Iwanesko — Tous droits réservés.</p>
          <p>Genève · Arc lémanique · Cours en ligne</p>
        </Container>
      </div>
    </footer>
  )
}
