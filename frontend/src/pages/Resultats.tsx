import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section, Eyebrow, CtaLink } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { IconQuote } from '../components/icons'
import { breadcrumbSchema, aggregateRatingSchema, type Crumb, type ReviewInput } from '../lib/schema'

// ⚠️ Placeholder testimonials — replace with real, attributable reviews only.
const REVIEWS: ReviewInput[] = [
  { author: 'Marc D.', rating: 5, body: 'En trois mois j’ai gagné 120 points Elo. La méthode est claire et le suivi vraiment personnalisé.' },
  { author: 'Sophie L.', rating: 5, body: 'Enfin un coach qui part de mes propres parties. Mes ouvertures ne sont plus un point faible.' },
  { author: 'Thomas R.', rating: 5, body: 'Préparation d’open impeccable : j’ai fait ma meilleure performance en tournoi.' },
  { author: 'Nadia B.', rating: 4, body: 'Cours en ligne très efficaces, l’échiquier partagé rend l’analyse limpide.' },
]

export function Component() {
  const crumbs: Crumb[] = [
    { name: 'Accueil', path: '/' },
    { name: 'Résultats', path: '/resultats' },
  ]

  return (
    <>
      <Seo
        title="Résultats & témoignages"
        description="Résultats des élèves et témoignages : progressions Elo, performances en tournoi et retours d’expérience sur les cours d’échecs à Genève."
        path="/resultats"
        jsonLd={[breadcrumbSchema(crumbs), aggregateRatingSchema(REVIEWS)]}
      />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero
        eyebrow="Résultats"
        title="Des progrès concrets, mesurés en points Elo"
        lead="Ce que les élèves obtiennent : progression au classement, meilleures performances en tournoi et un jeu plus clair."
        secondaryCta={{ to: '/a-propos', label: 'Découvrir la méthode' }}
      />

      <Section>
        <Container>
          <div className="overflow-hidden rounded-3xl border border-ink-800 bg-ink-950 shadow-lift">
            <div className="board-texture-dark grid gap-px bg-white/5 sm:grid-cols-3">
              {[
                { k: '+100', l: 'points Elo en moyenne sur un trimestre de travail assidu' },
                { k: '10/10', l: 'élèves recommandent le coaching à un autre joueur' },
                { k: '100 %', l: 'des plans de progression construits sur vos propres parties' },
              ].map((s) => (
                <div key={s.k} className="bg-ink-950 p-8 text-center">
                  <p className="font-display text-5xl font-extrabold tracking-tight text-gold-400">{s.k}</p>
                  <p className="mx-auto mt-3 max-w-[22ch] text-sm leading-relaxed text-ink-300">{s.l}</p>
                </div>
              ))}
            </div>
          </div>
          <p className="mt-4 text-center text-xs text-ink-400">
            Chiffres indicatifs, à remplacer par vos statistiques réelles avant mise en ligne.
          </p>
        </Container>
      </Section>

      <Section className="border-y border-ink-100 bg-cream-100">
        <Container>
          <Eyebrow>Témoignages</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">Ils ont progressé</h2>
          <div className="mt-10 grid gap-6 sm:grid-cols-2">
            {REVIEWS.map((r) => (
              <figure
                key={r.author}
                className="relative flex flex-col rounded-2xl border border-ink-200/80 bg-white p-7 shadow-soft"
              >
                <IconQuote size={38} className="text-gold-200" />
                <div className="mt-1 flex items-center gap-0.5 text-gold-500">
                  <span className="sr-only">Note : {r.rating} sur 5</span>
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
              Prêt à obtenir les mêmes résultats ?
            </h2>
            <p className="mx-auto mt-4 max-w-xl leading-relaxed text-ink-300">
              Réservez un premier cours pour définir vos objectifs et lancer votre progression.
            </p>
            <div className="mt-9">
              <CtaLink to="/contact" variant="primary">
                Obtenir les mêmes résultats
              </CtaLink>
            </div>
          </div>
        </Container>
      </Section>
    </>
  )
}
