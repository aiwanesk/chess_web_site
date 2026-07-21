import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section, Eyebrow, CtaLink } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { IconQuote } from '../components/icons'
import { breadcrumbSchema, aggregateRatingSchema, type Crumb, type ReviewInput } from '../lib/schema'
import { useLocale, homePath, pathFor, t, type Locale } from '../lib/i18n'

// ⚠️ Placeholder testimonials — replace with real, attributable reviews only.
const REVIEWS: Record<Locale, ReviewInput[]> = {
  fr: [
    { author: 'Marc D.', rating: 5, body: 'En trois mois j’ai gagné 120 points Elo. La méthode est claire et le suivi vraiment personnalisé.' },
    { author: 'Sophie L.', rating: 5, body: 'Enfin un coach qui part de mes propres parties. Mes ouvertures ne sont plus un point faible.' },
    { author: 'Thomas R.', rating: 5, body: 'Préparation d’open impeccable : j’ai fait ma meilleure performance en tournoi.' },
    { author: 'Nadia B.', rating: 4, body: 'Cours en ligne très efficaces, l’échiquier partagé rend l’analyse limpide.' },
  ],
  en: [
    { author: 'Marc D.', rating: 5, body: 'In three months I gained 120 Elo points. The method is clear and the coaching is genuinely tailored to me.' },
    { author: 'Sophie L.', rating: 5, body: 'Finally a coach who works from my own games. My openings are no longer a weakness.' },
    { author: 'Thomas R.', rating: 5, body: 'Flawless open preparation: I posted my best-ever tournament performance.' },
    { author: 'Nadia B.', rating: 4, body: 'Very effective online lessons — the shared board makes analysis crystal clear.' },
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
  statsNote: string
  reviewsEyebrow: string
  reviewsTitle: string
  ratingLabel: (rating: number) => string
  ctaTitle: string
  ctaLead: string
  ctaButton: string
}> = {
  fr: {
    path: '/resultats',
    title: 'Résultats & témoignages',
    description:
      'Résultats des élèves et témoignages : progressions Elo, performances en tournoi et retours d’expérience sur les cours d’échecs à Genève.',
    eyebrow: 'Résultats',
    heroTitle: 'Des progrès concrets, mesurés en points Elo',
    heroLead:
      'Ce que les élèves obtiennent : progression au classement, meilleures performances en tournoi et un jeu plus clair.',
    secondaryCtaLabel: 'Découvrir la méthode',
    secondaryCtaTo: '/a-propos',
    stats: [
      { k: '+100', l: 'points Elo en moyenne sur un trimestre de travail assidu' },
      { k: '10/10', l: 'élèves recommandent le coaching à un autre joueur' },
      { k: '100 %', l: 'des plans de progression construits sur vos propres parties' },
    ],
    statsNote: 'Chiffres indicatifs, à remplacer par vos statistiques réelles avant mise en ligne.',
    reviewsEyebrow: 'Témoignages',
    reviewsTitle: 'Ils ont progressé',
    ratingLabel: (rating) => `Note : ${rating} sur 5`,
    ctaTitle: 'Prêt à obtenir les mêmes résultats ?',
    ctaLead: 'Réservez un premier cours pour définir vos objectifs et lancer votre progression.',
    ctaButton: 'Obtenir les mêmes résultats',
  },
  en: {
    path: '/en/results',
    title: 'Results & testimonials',
    description:
      'Student results and testimonials: Elo gains, tournament performances and first-hand feedback on chess lessons in Geneva.',
    eyebrow: 'Results',
    heroTitle: 'Real progress, measured in Elo points',
    heroLead:
      'What students achieve: a higher rating, stronger tournament performances and clearer play.',
    secondaryCtaLabel: 'Discover the method',
    secondaryCtaTo: '/a-propos',
    stats: [
      { k: '+100', l: 'Elo points on average over a quarter of dedicated work' },
      { k: '10/10', l: 'students would recommend the coaching to another player' },
      { k: '100 %', l: 'of progress plans built on your own games' },
    ],
    statsNote: 'Indicative figures, to be replaced with your real statistics before going live.',
    reviewsEyebrow: 'Testimonials',
    reviewsTitle: 'They improved',
    ratingLabel: (rating) => `Rating: ${rating} out of 5`,
    ctaTitle: 'Ready to get the same results?',
    ctaLead: 'Book a first lesson to define your goals and kick-start your progress.',
    ctaButton: 'Get the same results',
  },
}

export function Component() {
  const locale = useLocale()
  const c = COPY[locale]
  const reviews = REVIEWS[locale]
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: c.title, path: c.path },
  ]

  return (
    <>
      <Seo
        title={c.title}
        description={c.description}
        path={c.path}
        jsonLd={[breadcrumbSchema(crumbs), aggregateRatingSchema(reviews)]}
      />
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
          <p className="mt-4 text-center text-xs text-ink-400">
            {c.statsNote}
          </p>
        </Container>
      </Section>

      <Section className="border-y border-ink-100 bg-cream-100">
        <Container>
          <Eyebrow>{c.reviewsEyebrow}</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">{c.reviewsTitle}</h2>
          <div className="mt-10 grid gap-6 sm:grid-cols-2">
            {reviews.map((r) => (
              <figure
                key={r.author}
                className="relative flex flex-col rounded-2xl border border-ink-200/80 bg-white p-7 shadow-soft"
              >
                <IconQuote size={38} className="text-gold-200" />
                <div className="mt-1 flex items-center gap-0.5 text-gold-500">
                  <span className="sr-only">{c.ratingLabel(r.rating)}</span>
                  {Array.from({ length: 5 }).map((_, i) => (
                    <svg
                      key={i}
                      width="18"
                      height="18"
                      viewBox="0 0 24 24"
                      fill={i < r.rating ? 'currentColor' : 'none'}
                      stroke="currentColor"
                      strokeWidth="1.5"
                      aria-hidden="true"
                      className={i < r.rating ? '' : 'text-ink-300'}
                    >
                      <path d="M12 3.5l2.6 5.3 5.9.9-4.3 4.1 1 5.8L12 17l-5.2 2.7 1-5.8-4.3-4.1 5.9-.9L12 3.5Z" />
                    </svg>
                  ))}
                </div>
                <blockquote className="mt-4 flex-1 leading-relaxed text-ink-700">“{r.body}”</blockquote>
                <figcaption className="mt-5 flex items-center gap-3 border-t border-ink-100 pt-4 text-sm font-semibold text-ink-900">
                  <span
                    aria-hidden
                    className="flex h-9 w-9 items-center justify-center rounded-full bg-cream-200 font-display text-gold-700"
                  >
                    {r.author.charAt(0)}
                  </span>
                  {r.author}
                </figcaption>
              </figure>
            ))}
          </div>
        </Container>
      </Section>

      <Section>
        <Container>
          <div className="board-texture-dark rounded-3xl bg-ink-950 px-6 py-16 text-center shadow-lift">
            <span aria-hidden className="mx-auto mb-6 block h-0.5 w-12 rounded-full bg-gold-500" />
            <h2 className="font-display text-3xl font-bold text-white sm:text-4xl">
              {c.ctaTitle}
            </h2>
            <p className="mx-auto mt-4 max-w-xl leading-relaxed text-ink-300">
              {c.ctaLead}
            </p>
            <div className="mt-9">
              <CtaLink to={pathFor('contact', locale)} variant="primary">
                {c.ctaButton}
              </CtaLink>
            </div>
          </div>
        </Container>
      </Section>
    </>
  )
}
