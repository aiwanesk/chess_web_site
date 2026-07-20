import type { ReactNode } from 'react'
import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { Container } from './Container'
import { Section, Eyebrow } from './ui'
import { Breadcrumbs } from './Breadcrumbs'
import { PageHero, FactCard } from './PageHero'
import { Faq } from './Faq'
import { postsByCluster } from '../lib/postMeta'
import {
  IconTarget,
  IconRoute,
  IconBoard,
  IconClock,
  IconTrophy,
  IconGraduation,
  IconArrowRight,
} from './icons'
import {
  breadcrumbSchema,
  courseSchema,
  faqSchema,
  type CourseInput,
  type FaqItem,
  type Crumb,
} from '../lib/schema'

/** Decorative icons cycled across the benefit cards (chess / coaching theme). */
const BENEFIT_ICONS = [IconTarget, IconRoute, IconBoard, IconClock, IconTrophy, IconGraduation]

export interface RelatedLink {
  to: string
  label: string
}

export interface MoneyPageProps {
  /** SEO */
  title: string // <title> + H1
  metaTitle?: string // override <title> if it should differ from H1
  description: string
  path: string
  /** Hero */
  eyebrow: string
  lead: ReactNode
  facts: Array<{ label: string; value: string }>
  /** Answer-first intro paragraph(s) — great for GEO extraction. */
  intro: ReactNode
  /** "What you get" bullet points. */
  benefitsTitle?: string
  benefits: Array<{ title: string; body: string }>
  /** Structured data */
  course: CourseInput
  faq: FaqItem[]
  /** Internal linking (silo). */
  related?: RelatedLink[]
  /** Content silo key: surfaces matching blog articles (see postMeta.ts). */
  cluster?: string
  /** Extra JSON-LD objects (e.g. Event for stages/simultanées). */
  extraJsonLd?: Array<Record<string, unknown>>
  /** Optional extra content injected before the FAQ. */
  children?: ReactNode
}

/**
 * Template shared by every "money page": semantic HTML, single H1, answer-first
 * intro, benefits, internal-link silo, FAQ, and the full JSON-LD stack
 * (BreadcrumbList + Course/Offer + FAQPage).
 */
export function MoneyPage(props: MoneyPageProps) {
  const crumbs: Crumb[] = [
    { name: 'Accueil', path: '/' },
    { name: props.title, path: props.path },
  ]

  const articles = props.cluster ? postsByCluster(props.cluster) : []

  const jsonLd = [
    breadcrumbSchema(crumbs),
    courseSchema(props.course),
    faqSchema(props.faq),
    ...(props.extraJsonLd ?? []),
  ]

  return (
    <>
      <Seo
        title={props.metaTitle ?? props.title}
        description={props.description}
        path={props.path}
        jsonLd={jsonLd}
      />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero
        eyebrow={props.eyebrow}
        title={props.title}
        lead={props.lead}
        secondaryCta={{ to: '/tarifs', label: 'Voir les tarifs' }}
        aside={<FactCard facts={props.facts} />}
      />

      {/* Answer-first intro */}
      <Section>
        <Container>
          <div className="prose">{props.intro}</div>
        </Container>
      </Section>

      {/* Benefits */}
      <Section className="border-y border-ink-100 bg-cream-100">
        <Container>
          <Eyebrow>La méthode</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">
            {props.benefitsTitle ?? 'Ce que vous obtenez'}
          </h2>
          <div className="mt-10 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {props.benefits.map((b, i) => {
              const Icon = BENEFIT_ICONS[i % BENEFIT_ICONS.length]!
              return (
                <div
                  key={b.title}
                  className="hover-lift rounded-2xl border border-ink-200/80 bg-white p-7 shadow-soft hover:border-gold-300 hover:shadow-card"
                >
                  <span
                    aria-hidden
                    className="flex h-11 w-11 items-center justify-center rounded-xl border border-gold-200 bg-gold-50 text-gold-700"
                  >
                    <Icon size={22} />
                  </span>
                  <h3 className="mt-5 text-lg font-semibold text-ink-900">{b.title}</h3>
                  <p className="mt-2 leading-relaxed text-ink-600">{b.body}</p>
                </div>
              )
            })}
          </div>
        </Container>
      </Section>

      {props.children}

      {/* Internal-link silo */}
      {props.related && props.related.length > 0 ? (
        <Section>
          <Container>
            <Eyebrow>Explorer</Eyebrow>
            <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">Pour aller plus loin</h2>
            <ul className="mt-8 grid gap-4 sm:grid-cols-2">
              {props.related.map((r) => (
                <li key={r.to}>
                  <Link
                    to={r.to}
                    className="group flex items-center justify-between gap-4 rounded-2xl border border-ink-200/80 bg-white px-6 py-5 font-medium text-ink-800 shadow-soft transition-colors hover:border-gold-300 hover:text-ink-950"
                  >
                    {r.label}
                    <span
                      aria-hidden
                      className="flex h-9 w-9 flex-none items-center justify-center rounded-full bg-cream-100 text-gold-700 transition-[background-color,transform] group-hover:translate-x-0.5 group-hover:bg-gold-100"
                    >
                      <IconArrowRight size={18} />
                    </span>
                  </Link>
                </li>
              ))}
            </ul>
          </Container>
        </Section>
      ) : null}

      {/* Related articles (blog → money page silo, reciprocal link). */}
      {articles.length > 0 ? (
        <Section className="border-y border-ink-100 bg-cream-100">
          <Container>
            <Eyebrow>À lire</Eyebrow>
            <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">Articles liés</h2>
            <ul className="mt-8 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {articles.map((a) => (
                <li key={a.slug}>
                  <Link
                    to={`/blog/${a.slug}`}
                    className="hover-lift flex h-full flex-col rounded-2xl border border-ink-200/80 bg-white p-6 shadow-soft transition-colors hover:border-gold-300 hover:shadow-card"
                  >
                    <h3 className="font-semibold text-ink-900">{a.title}</h3>
                    <p className="mt-2 flex-1 text-sm leading-relaxed text-ink-600">{a.description}</p>
                    <span className="mt-4 inline-flex items-center gap-1.5 text-sm font-semibold text-gold-700">
                      Lire l’article
                      <IconArrowRight size={15} />
                    </span>
                  </Link>
                </li>
              ))}
            </ul>
          </Container>
        </Section>
      ) : null}

      <Faq items={props.faq} />

      {/* Closing CTA */}
      <Section className="board-texture-dark bg-ink-950 text-center text-white">
        <Container>
          <span aria-hidden className="mx-auto mb-6 block h-0.5 w-12 rounded-full bg-gold-500" />
          <h2 className="font-display text-3xl font-bold sm:text-4xl">Prêt à progresser ?</h2>
          <p className="mx-auto mt-4 max-w-xl leading-relaxed text-ink-300">
            Un premier échange pour définir vos objectifs et le format qui vous convient.
          </p>
          <div className="mt-9">
            <Link
              to="/contact"
              className="inline-flex rounded-full bg-gold-500 px-7 py-3 font-semibold text-ink-950 shadow-gold transition-[background-color,transform] duration-200 hover:bg-gold-400 active:translate-y-px"
            >
              Réserver un premier cours
            </Link>
          </div>
        </Container>
      </Section>
    </>
  )
}
