import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section, Eyebrow, CtaLink } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { IconGraduation, IconTrophy, IconArrowRight } from '../components/icons'
import { breadcrumbSchema, type Crumb } from '../lib/schema'
import { TOURNAMENTS } from '../lib/tournaments'
import { useLocale, homePath, pathFor, t, type Locale } from '../lib/i18n'

/**
 * Student results.
 *
 * These are cases Alexandre describes himself — what was worked on and what it
 * produced — NOT reviews written by the students. So the page carries no star
 * ratings and emits no AggregateRating/Review JSON-LD: declaring ratings nobody
 * actually gave would be fabricated structured data (a manual-action risk with
 * Google, and misleading advertising under Swiss LCD art. 3 / the EU Omnibus
 * directive). If a student later sends a written review, add the quote here and
 * only then restore the rating schema.
 */
interface StudentCase {
  name: string
  focus: string // what we worked on
  outcome: string // what it produced
  metric?: string // the headline figure, when there is one
}

const CASES: Record<Locale, StudentCase[]> = {
  fr: [
    {
      name: 'Flavien S.',
      focus: 'Compréhension du jeu — lire les positions au lieu de calculer à l’aveugle.',
      outcome:
        'Deux mois de travail sur la compréhension positionnelle, et le classement a suivi : passé de 2000 à 2100 Elo FIDE.',
      metric: '+100 Elo FIDE en 2 mois',
    },
    {
      name: 'Lucas M.',
      focus: 'Construction d’un répertoire aux Noirs autour de la Benoni.',
      outcome:
        'Les plans typiques et les structures avant les variantes : il sait maintenant ce qu’il cherche dans la position, pas seulement les quinze premiers coups.',
      metric: 'Répertoire Benoni',
    },
    {
      name: 'Mathilde B.',
      focus: 'Apprentissage de l’Est-Indienne, une défense ambitieuse et exigeante.',
      outcome:
        'Attaque à l’aile roi, tension centrale, choix du bon moment pour rompre — une ouverture qui apprend le milieu de jeu autant que la théorie.',
      metric: 'Répertoire Est-Indienne',
    },
  ],
  en: [
    {
      name: 'Flavien S.',
      focus: 'Understanding the game — reading positions instead of calculating blind.',
      outcome:
        'Two months of work on positional understanding, and the rating followed: from 2000 to 2100 FIDE.',
      metric: '+100 FIDE Elo in 2 months',
    },
    {
      name: 'Lucas M.',
      focus: 'Building a Black repertoire around the Benoni.',
      outcome:
        'Typical plans and structures before variations: he now knows what he is looking for in the position, not just the first fifteen moves.',
      metric: 'Benoni repertoire',
    },
    {
      name: 'Mathilde B.',
      focus: 'Learning the King’s Indian, an ambitious and demanding defence.',
      outcome:
        'Kingside attack, central tension, picking the right moment to break — an opening that teaches the middlegame as much as the theory.',
      metric: 'King’s Indian repertoire',
    },
  ],
}

const COPY: Record<Locale, {
  path: string
  title: string
  description: string
  eyebrow: string
  heroTitle: string
  heroLead: string
  secondaryCtaLabel: string
  secondaryCtaTo: string
  stats: { k: string; l: string }[]
  casesEyebrow: string
  casesTitle: string
  casesLead: string
  focusLabel: string
  calendarLink: string
  ctaTitle: string
  ctaLead: string
  ctaButton: string
}> = {
  fr: {
    path: '/resultats',
    title: 'Résultats des élèves',
    description:
      'Ce que les élèves obtiennent avec un Maître FIDE : progression Elo, répertoire d’ouvertures construit et compréhension du jeu. Trois parcours détaillés.',
    eyebrow: 'Résultats',
    heroTitle: 'Des progrès concrets, mesurés en points Elo',
    heroLead:
      'Trois élèves, trois chantiers différents : la compréhension du jeu, un répertoire aux Noirs, une défense ambitieuse. Voici ce qu’on a travaillé et ce que ça a donné.',
    secondaryCtaLabel: 'Découvrir la méthode',
    secondaryCtaTo: '/a-propos',
    stats: [
      { k: '+100', l: 'points Elo FIDE pris par Flavien en deux mois (2000 → 2100)' },
      { k: '100 %', l: 'des plans de progression construits sur vos propres parties' },
      { k: String(TOURNAMENTS.length), l: 'tournois que je joue moi-même cette saison — je suis toujours en compétition' },
    ],
    casesEyebrow: 'Parcours d’élèves',
    casesTitle: 'Ce qu’on a travaillé, ce que ça a donné',
    casesLead:
      'Chaque plan part des parties de l’élève. Ces trois-là n’avaient pas le même point faible — donc pas le même programme.',
    focusLabel: 'Le chantier',
    calendarLink: 'Voir mon calendrier de tournois',
    ctaTitle: 'Prêt à obtenir les mêmes résultats ?',
    ctaLead: 'Réservez un premier cours pour définir vos objectifs et lancer votre progression.',
    ctaButton: 'Obtenir les mêmes résultats',
  },
  en: {
    path: '/en/results',
    title: 'Student results',
    description:
      'What students achieve with a FIDE Master: Elo progress, a built opening repertoire and real understanding of the game. Three detailed cases.',
    eyebrow: 'Results',
    heroTitle: 'Real progress, measured in Elo points',
    heroLead:
      'Three students, three different jobs: understanding the game, a Black repertoire, an ambitious defence. Here is what we worked on and what it produced.',
    secondaryCtaLabel: 'Discover the method',
    secondaryCtaTo: '/en/about',
    stats: [
      { k: '+100', l: 'FIDE Elo points gained by Flavien in two months (2000 → 2100)' },
      { k: '100 %', l: 'of progress plans built on your own games' },
      { k: String(TOURNAMENTS.length), l: 'tournaments I play myself this season — I am still competing' },
    ],
    casesEyebrow: 'Student cases',
    casesTitle: 'What we worked on, and what it produced',
    casesLead:
      'Every plan starts from the student’s own games. These three did not share the same weakness — so they did not get the same programme.',
    focusLabel: 'The job',
    calendarLink: 'See my tournament calendar',
    ctaTitle: 'Ready to get the same results?',
    ctaLead: 'Book a first lesson to define your goals and kick-start your progress.',
    ctaButton: 'Get the same results',
  },
}

export function Component() {
  const locale = useLocale()
  const c = COPY[locale]
  const cases = CASES[locale]
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: c.title, path: c.path },
  ]

  return (
    <>
      <Seo title={c.title} description={c.description} path={c.path} jsonLd={[breadcrumbSchema(crumbs)]} />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero
        eyebrow={c.eyebrow}
        title={c.heroTitle}
        lead={c.heroLead}
        secondaryCta={{ to: c.secondaryCtaTo, label: c.secondaryCtaLabel }}
      />

      <Section>
        <Container>
          <div className="overflow-hidden rounded-3xl border border-ink-800 bg-ink-950 shadow-lift">
            <div className="board-texture-dark grid gap-px bg-white/5 sm:grid-cols-3">
              {c.stats.map((s) => (
                <div key={s.k} className="bg-ink-950 p-8 text-center">
                  <p className="font-display text-5xl font-extrabold tracking-tight text-gold-400">{s.k}</p>
                  <p className="mx-auto mt-3 max-w-[22ch] text-sm leading-relaxed text-ink-300">{s.l}</p>
                </div>
              ))}
            </div>
          </div>
          <p className="mt-4 text-center text-sm">
            <Link
              to={pathFor('calendrier', locale)}
              className="inline-flex items-center gap-1.5 font-semibold text-gold-700 underline decoration-gold-400 decoration-1 underline-offset-4 hover:text-ink-950"
            >
              <IconTrophy size={15} />
              {c.calendarLink}
            </Link>
          </p>
        </Container>
      </Section>

      <Section className="border-y border-ink-100 bg-cream-100">
        <Container>
          <Eyebrow>{c.casesEyebrow}</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">{c.casesTitle}</h2>
          <p className="mt-4 max-w-[60ch] leading-relaxed text-ink-600">{c.casesLead}</p>

          <div className="mt-10 grid gap-6 md:grid-cols-3">
            {cases.map((s) => (
              <article
                key={s.name}
                className="flex flex-col rounded-2xl border border-ink-200/80 bg-white p-7 shadow-soft"
              >
                <span
                  aria-hidden
                  className="flex h-11 w-11 items-center justify-center rounded-xl border border-gold-200 bg-gold-50 text-gold-700"
                >
                  <IconGraduation size={22} />
                </span>
                <h3 className="mt-5 font-display text-xl font-bold text-ink-900">{s.name}</h3>
                {s.metric ? (
                  <p className="mt-2 inline-flex self-start rounded-full bg-ink-900 px-3 py-1 text-xs font-bold text-gold-400">
                    {s.metric}
                  </p>
                ) : null}
                <p className="mt-5 text-[0.68rem] font-semibold uppercase tracking-[0.14em] text-ink-500">
                  {c.focusLabel}
                </p>
                <p className="mt-1.5 leading-relaxed text-ink-700">{s.focus}</p>
                <p className="mt-4 flex-1 leading-relaxed text-ink-600">{s.outcome}</p>
              </article>
            ))}
          </div>
        </Container>
      </Section>

      <Section>
        <Container>
          <div className="board-texture-dark rounded-3xl bg-ink-950 px-6 py-16 text-center shadow-lift">
            <span aria-hidden className="mx-auto mb-6 block h-0.5 w-12 rounded-full bg-gold-500" />
            <h2 className="font-display text-3xl font-bold text-white sm:text-4xl">{c.ctaTitle}</h2>
            <p className="mx-auto mt-4 max-w-xl leading-relaxed text-ink-300">{c.ctaLead}</p>
            <div className="mt-9">
              <CtaLink to={pathFor('contact', locale)} variant="primary">
                {c.ctaButton}
                <span aria-hidden className="transition-transform group-hover/btn:translate-x-1">
                  <IconArrowRight size={16} />
                </span>
              </CtaLink>
            </div>
          </div>
        </Container>
      </Section>
    </>
  )
}
