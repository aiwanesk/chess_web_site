import type { ReactNode } from 'react'
import { Seo } from '../lib/seo'
import { SITE } from '../lib/site'
import { Container } from '../components/Container'
import { Section, Eyebrow, CtaLink } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { ChessMotif } from '../components/ChessMotif'
import { personSchema, breadcrumbSchema, type Crumb } from '../lib/schema'
import { useLocale, homePath, t, type Locale } from '../lib/i18n'

type Field = { label: string; value: string }

const COPY: Record<Locale, {
  path: string; title: string; description: string; eyebrow: string; crumb: string
  heroTitle: string; lead: ReactNode; primaryCta: string; secondaryCta: string
  prose: ReactNode; retrouvezEyebrow: string; retrouvezP: string; bottomCta: string
  fields: Field[]
}> = {
  fr: {
    path: '/a-propos',
    title: 'À propos d’Alexandre Iwanesko, Maître FIDE',
    description:
      'Parcours d’Alexandre Iwanesko, Maître FIDE et coach d’échecs à Genève : titre, méthode d’enseignement et philosophie de progression.',
    eyebrow: 'À propos',
    crumb: 'À propos',
    heroTitle: 'Alexandre Iwanesko, Maître FIDE & coach d’échecs à Genève',
    lead: (
      <>
        {SITE.person.description} Mon objectif&nbsp;: rendre votre progression{' '}
        <strong>concrète et mesurable</strong>.
      </>
    ),
    primaryCta: 'Me contacter',
    secondaryCta: 'Voir les résultats',
    prose: (
      <>
        <h2>Un Maître FIDE au service de votre jeu</h2>
        <p>
          <strong>Un Maître FIDE</strong> est un joueur ayant obtenu un titre international
          décerné par la Fédération internationale des échecs (FIDE), reconnaissant un niveau
          Elo élevé et stable. Ce titre est le socle d’une expertise que je mets au service de
          votre progression.
        </p>
        <p>
          J’enseigne aux <strong>adultes de 1200 à 2200 Elo</strong> et aux{' '}
          <strong>adolescents en compétition</strong>. Mon approche est structurée : on part de
          vos parties réelles, on identifie les leviers de progrès les plus rentables, et on
          construit un plan de travail sur plusieurs semaines.
        </p>
        <h2>Ma méthode</h2>
        <p>
          Pas de recettes toutes faites : un diagnostic honnête, des priorités claires, et un
          travail régulier entre les séances. Les échecs se progressent par la compréhension,
          pas par la mémorisation.
        </p>
      </>
    ),
    retrouvezEyebrow: 'Retrouvez-moi',
    retrouvezP:
      'Profil FIDE, parties commentées et jeu en ligne — les liens officiels (à compléter) renforcent la vérifiabilité de mon expertise.',
    bottomCta: 'Réserver un premier cours',
    fields: [
      { label: 'Titre', value: 'Maître FIDE' },
      { label: 'Public', value: '1200–2200 Elo' },
      { label: 'Lieu', value: 'Genève / en ligne' },
      { label: 'Langues', value: 'FR · EN' },
    ],
  },
  en: {
    path: '/en/about',
    title: 'About Alexandre Iwanesko, FIDE Master',
    description:
      'The background of Alexandre Iwanesko, FIDE Master and chess coach in Geneva: title, teaching method and philosophy of progress.',
    eyebrow: 'About',
    crumb: 'About',
    heroTitle: 'Alexandre Iwanesko, FIDE Master & chess coach in Geneva',
    lead: (
      <>
        FIDE Master and chess coach in Geneva, specialising in the progress of adults (1200–2200
        Elo) and competitive teenagers. My goal: to make your progress{' '}
        <strong>concrete and measurable</strong>.
      </>
    ),
    primaryCta: 'Get in touch',
    secondaryCta: 'See the results',
    prose: (
      <>
        <h2>A FIDE Master at the service of your game</h2>
        <p>
          <strong>A FIDE Master</strong> is a player who has earned an international title awarded
          by the International Chess Federation (FIDE), recognising a high and stable Elo level.
          This title is the foundation of an expertise I put at the service of your progress.
        </p>
        <p>
          I teach <strong>adults from 1200 to 2200 Elo</strong> and{' '}
          <strong>competitive teenagers</strong>. My approach is structured: we start from your
          real games, identify the highest-value levers for improvement, and build a work plan
          spanning several weeks.
        </p>
        <h2>My method</h2>
        <p>
          No ready-made recipes: an honest diagnosis, clear priorities, and regular work between
          sessions. In chess you improve through understanding, not memorisation.
        </p>
      </>
    ),
    retrouvezEyebrow: 'Find me',
    retrouvezP:
      'FIDE profile, annotated games and online play — the official links (to be completed) reinforce the verifiability of my expertise.',
    bottomCta: 'Book a first lesson',
    fields: [
      { label: 'Title', value: 'FIDE Master' },
      { label: 'Audience', value: '1200–2200 Elo' },
      { label: 'Location', value: 'Geneva / online' },
      { label: 'Languages', value: 'FR · EN' },
    ],
  },
}

export function Component() {
  const locale = useLocale()
  const c = COPY[locale]
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: c.crumb, path: c.path },
  ]

  return (
    <>
      <Seo
        title={c.title}
        description={c.description}
        path={c.path}
        jsonLd={[personSchema(), breadcrumbSchema(crumbs)]}
      />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero
        eyebrow={c.eyebrow}
        title={c.heroTitle}
        lead={c.lead}
        primaryCta={{ to: '/contact', label: c.primaryCta }}
        secondaryCta={{ to: '/resultats', label: c.secondaryCta }}
      />

      <Section>
        <Container className="grid gap-12 lg:grid-cols-[1.5fr_1fr] lg:items-start">
          <div>
          <div className="prose">
            {c.prose}
          </div>

          <div className="mt-10 rounded-2xl border border-ink-200/80 bg-cream-100 p-7 shadow-soft">
            <Eyebrow>{c.retrouvezEyebrow}</Eyebrow>
            <p className="text-ink-600">
              {c.retrouvezP}
            </p>
            <ul className="mt-5 flex flex-wrap gap-3 text-sm font-medium">
              {SITE.person.sameAs.map((url) => (
                <li key={url}>
                  <a
                    href={url}
                    rel="me noopener"
                    target="_blank"
                    className="inline-flex rounded-full border border-ink-300 bg-white px-4 py-2 text-ink-700 transition-colors hover:border-gold-400 hover:text-gold-700"
                  >
                    {new URL(url).hostname.replace('www.', '')}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          <div className="mt-10">
            <CtaLink to="/contact" variant="primary">
              {c.bottomCta}
            </CtaLink>
          </div>
          </div>

          {/* Visual aside — brand motif + credentials at a glance. */}
          <aside className="lg:sticky lg:top-24">
            <div className="overflow-hidden rounded-3xl border border-ink-200/80 bg-white p-6 shadow-card">
              <ChessMotif className="mx-auto w-full max-w-[16rem] text-ink-900" />
              <dl className="mt-6 grid grid-cols-2 gap-x-4 gap-y-4">
                {c.fields.map((f) => (
                  <div key={f.label}>
                    <dt className="text-[0.68rem] font-semibold uppercase tracking-[0.12em] text-ink-500">
                      {f.label}
                    </dt>
                    <dd className="mt-1 font-display font-bold text-ink-900">{f.value}</dd>
                  </div>
                ))}
              </dl>
            </div>
          </aside>
        </Container>
      </Section>
    </>
  )
}
