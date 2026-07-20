import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { SITE } from '../lib/site'
import { Container } from '../components/Container'
import { Section, Eyebrow, CtaLink } from '../components/ui'
import { PageHero, FactCard } from '../components/PageHero'
import { ChessMotif } from '../components/ChessMotif'
import {
  IconGraduation,
  IconTrophy,
  IconMonitor,
  IconUsers,
  IconSpark,
  IconBoard,
  IconCheck,
  IconArrowRight,
} from '../components/icons'
import { personSchema, localBusinessSchema, faqSchema, breadcrumbSchema } from '../lib/schema'

const OFFERS = [
  { to: '/cours-echecs-adultes-geneve', title: 'Cours pour adultes', desc: 'Progression structurée pour joueurs 1200–1900 Elo.', Icon: IconGraduation },
  { to: '/preparation-tournoi-echecs', title: 'Préparation tournoi', desc: 'Répertoire, gestion du temps et mental de compétition.', Icon: IconTrophy },
  { to: '/cours-echecs-en-ligne', title: 'Cours en ligne', desc: 'Même méthode à distance, partout en Suisse romande.', Icon: IconMonitor },
  { to: '/cours-echecs-groupe-geneve', title: 'Cours en groupe', desc: 'Petits groupes de niveau homogène, tarif partagé.', Icon: IconUsers },
  { to: '/cours-echecs-ados-competition', title: 'Ados en compétition', desc: 'Coaching orienté progression Elo et résultats.', Icon: IconSpark },
  { to: '/stages-echecs-geneve', title: 'Stages à Genève', desc: 'Sessions intensives pendant les vacances scolaires.', Icon: IconBoard },
]

const HOME_FAQ = [
  {
    question: 'À qui s’adressent les cours d’Alexandre Iwanesko ?',
    answer:
      'Aux adultes de 1200 à 1900 Elo et aux adolescents en compétition. Ce sont des cours de perfectionnement : les grands débutants ne sont pas le public visé.',
  },
  {
    question: 'Qui est Alexandre Iwanesko ?',
    answer:
      'Alexandre Iwanesko est Maître FIDE et coach d’échecs à Genève. Il enseigne aux joueurs intermédiaires et avancés, en présentiel et en ligne.',
  },
  {
    question: 'Les cours ont-ils lieu à Genève ou en ligne ?',
    answer:
      'Les deux : en présentiel à Genève et en ligne pour toute la Suisse romande et la France voisine, avec la même méthode et le même suivi.',
  },
]

export function Component() {
  const jsonLd = [
    personSchema(),
    localBusinessSchema(),
    breadcrumbSchema([{ name: 'Accueil', path: '/' }]),
    faqSchema(HOME_FAQ),
  ]

  return (
    <>
      <Seo
        isHome
        title="Coach d’échecs à Genève — Alexandre Iwanesko, Maître FIDE"
        description="Coach d’échecs à Genève, Maître FIDE. Cours pour adultes (1200–1900 Elo) et ados en compétition, préparation tournoi, en présentiel et en ligne."
        path="/"
        jsonLd={jsonLd}
      />

      <PageHero
        eyebrow="Maître FIDE · Genève & en ligne"
        title="Progressez aux échecs avec un Maître FIDE à Genève"
        lead={
          <>
            Coaching d’échecs pour <strong>adultes (1200–1900 Elo)</strong> et{' '}
            <strong>ados en compétition</strong>. Une méthode claire, un plan de progression et un
            vrai suivi — en présentiel à Genève ou en ligne.
          </>
        }
        secondaryCta={{ to: '/cours-echecs-adultes-geneve', label: 'Découvrir les cours adultes' }}
        aside={
          <div className="relative">
            <ChessMotif className="mx-auto w-56 text-ink-900 drop-shadow-xl sm:w-64" />
            <div className="mt-6">
              <FactCard
                facts={[
                  { label: 'Titre', value: 'Maître FIDE' },
                  { label: 'Public', value: '1200–1900 Elo' },
                  { label: 'Lieu', value: 'Genève / en ligne' },
                  { label: 'Langues', value: 'FR · EN' },
                ]}
              />
            </div>
          </div>
        }
      />

      {/* Offers hub — links to every money page (silo top). */}
      <Section>
        <Container>
          <Eyebrow>Formats</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">
            Trouvez le cours qui vous correspond
          </h2>
          <div className="mt-10 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {OFFERS.map((o) => (
              <Link
                key={o.to}
                to={o.to}
                className="hover-lift group flex flex-col rounded-2xl border border-ink-200/80 bg-white p-7 shadow-soft transition-colors hover:border-gold-300 hover:shadow-card"
              >
                <span
                  aria-hidden
                  className="flex h-12 w-12 items-center justify-center rounded-xl border border-gold-200 bg-gold-50 text-gold-700"
                >
                  <o.Icon size={24} />
                </span>
                <h3 className="mt-5 flex items-center justify-between text-lg font-semibold text-ink-900">
                  {o.title}
                  <span
                    aria-hidden
                    className="text-gold-600 transition-transform group-hover:translate-x-1"
                  >
                    <IconArrowRight size={18} />
                  </span>
                </h3>
                <p className="mt-2 leading-relaxed text-ink-600">{o.desc}</p>
              </Link>
            ))}
          </div>
        </Container>
      </Section>

      {/* Why me */}
      <Section className="border-y border-ink-100 bg-cream-100">
        <Container className="grid gap-12 lg:grid-cols-2 lg:items-center">
          <div>
            <Eyebrow>Pourquoi un Maître FIDE</Eyebrow>
            <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">
              Une expertise de titre, une pédagogie de coach
            </h2>
            <p className="mt-5 text-lg leading-relaxed text-ink-600">
              Le titre de <strong className="text-ink-900">Maître FIDE</strong> garantit un niveau de
              jeu élevé. Mais progresser, c’est surtout une méthode : diagnostic précis, plan de
              travail personnalisé et exercices ciblés entre les séances.
            </p>
            <div className="mt-8">
              <CtaLink to="/a-propos" variant="secondary">
                En savoir plus sur mon parcours
              </CtaLink>
            </div>
          </div>
          <ul className="space-y-3">
            {[
              'Diagnostic de vos parties classées pour cibler les vrais leviers.',
              'Plan de progression sur 8 à 12 semaines avec objectifs Elo.',
              'Présentiel à Genève ou en ligne, au même niveau d’exigence.',
              'Public sérieux : joueurs qui veulent progresser, pas s’initier.',
            ].map((t) => (
              <li
                key={t}
                className="flex items-start gap-4 rounded-2xl border border-ink-200/80 bg-white p-5 shadow-soft"
              >
                <span
                  aria-hidden
                  className="mt-0.5 flex h-7 w-7 flex-none items-center justify-center rounded-full bg-gold-100 text-gold-700"
                >
                  <IconCheck size={16} />
                </span>
                <span className="leading-relaxed text-ink-700">{t}</span>
              </li>
            ))}
          </ul>
        </Container>
      </Section>

      {/* FAQ */}
      <Section>
        <Container>
          <Eyebrow>FAQ</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">Questions fréquentes</h2>
          <div className="mx-auto mt-8 max-w-4xl space-y-3">
            {HOME_FAQ.map((f) => (
              <details
                key={f.question}
                className="group rounded-2xl border border-ink-200/80 bg-white px-6 shadow-soft transition-colors open:border-gold-300"
              >
                <summary className="flex cursor-pointer list-none items-center justify-between gap-4 py-5 text-lg font-semibold text-ink-900">
                  {f.question}
                  <span
                    aria-hidden
                    className="flex h-8 w-8 flex-none items-center justify-center rounded-full border border-gold-300 bg-gold-50 text-gold-700 transition-transform duration-300 group-open:rotate-45"
                  >
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                      <path d="M12 5v14M5 12h14" />
                    </svg>
                  </span>
                </summary>
                <p className="max-w-3xl pb-6 leading-relaxed text-ink-600">{f.answer}</p>
              </details>
            ))}
          </div>
        </Container>
      </Section>

      {/* CTA band */}
      <Section className="board-texture-dark bg-ink-950 text-center text-white">
        <Container>
          <span aria-hidden className="mx-auto mb-6 block h-0.5 w-12 rounded-full bg-gold-500" />
          <h2 className="font-display text-3xl font-bold sm:text-4xl">Envie de franchir un palier ?</h2>
          <p className="mx-auto mt-4 max-w-xl leading-relaxed text-ink-300">
            Réservez un premier cours à {SITE.address.locality} ou en ligne pour définir vos
            objectifs.
          </p>
          <div className="mt-9">
            <CtaLink to="/contact" variant="primary">
              Réserver un premier cours
            </CtaLink>
          </div>
        </Container>
      </Section>
    </>
  )
}
