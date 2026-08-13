import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { TournamentCalendar } from '../components/TournamentCalendar'
import { breadcrumbSchema, sportsEventSchema, type Crumb } from '../lib/schema'
import { categoryPath } from '../lib/categories'
import { TOURNAMENTS } from '../lib/tournaments'
import { locationFor, nameFor, statusOf } from '../lib/calendar'
import { useLocale, homePath, pathFor, t, type Locale } from '../lib/i18n'

const COPY: Record<Locale, {
  title: string
  description: string
  eyebrow: string
  heroTitle: string
  heroLead: string
  ctaLabel: string
  secondaryLabel: string
}> = {
  fr: {
    title: 'Calendrier des tournois',
    description:
      'Le calendrier de compétition d’Alexandre Iwanesko, Maître FIDE : les tournois à venir et les carnets de tournoi des opens déjà joués.',
    eyebrow: 'Calendrier',
    heroTitle: 'Où je joue cette saison',
    heroLead:
      'Les opens que je prépare, ceux que je viens de jouer. Un tournoi passé mène directement à son carnet — la compétition racontée de l’intérieur, ronde par ronde.',
    ctaLabel: 'Préparer votre tournoi',
    secondaryLabel: 'Lire les carnets de tournoi',
  },
  en: {
    title: 'Tournament calendar',
    description:
      'The competition calendar of Alexandre Iwanesko, FIDE Master: upcoming tournaments and the tournament diaries of the opens already played.',
    eyebrow: 'Calendar',
    heroTitle: 'Where I’m playing this season',
    heroLead:
      'The opens I’m preparing for, and the ones I’ve just played. Past tournaments link straight to their diary — competition told from the inside, round by round.',
    ctaLabel: 'Prepare your tournament',
    secondaryLabel: 'Read the tournament diaries',
  },
}

export function Component() {
  const locale = useLocale()
  const c = COPY[locale]
  const path = pathFor('calendrier', locale)
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: c.title, path },
  ]

  // Declare the tournaments still ahead as SportsEvents. Uses the build date
  // rather than the visitor's: JSON-LD is read from the served HTML, so what
  // matters is what was true when the page was generated.
  const eventsJsonLd = TOURNAMENTS.filter((e) => statusOf(e, __BUILD_DATE__) !== 'past').map((e) =>
    sportsEventSchema({
      name: nameFor(e, locale),
      startDate: e.start,
      endDate: e.end,
      location: locationFor(e, locale),
      url: path,
    }),
  )

  return (
    <>
      <Seo
        title={c.title}
        description={c.description}
        path={path}
        jsonLd={[breadcrumbSchema(crumbs), ...eventsJsonLd]}
      />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero
        eyebrow={c.eyebrow}
        title={c.heroTitle}
        lead={c.heroLead}
        primaryCta={{ to: pathFor('preparationTournoi', locale), label: c.ctaLabel }}
        secondaryCta={{ to: categoryPath('carnet-de-tournoi', locale), label: c.secondaryLabel }}
      />

      <Section>
        <Container>
          <TournamentCalendar />
        </Container>
      </Section>
    </>
  )
}
