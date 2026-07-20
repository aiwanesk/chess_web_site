import { Seo } from '../lib/seo'
import { SITE } from '../lib/site'
import { Container } from '../components/Container'
import { Section, Eyebrow, CtaLink } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { ChessMotif } from '../components/ChessMotif'
import { personSchema, breadcrumbSchema, type Crumb } from '../lib/schema'

export function Component() {
  const crumbs: Crumb[] = [
    { name: 'Accueil', path: '/' },
    { name: 'À propos', path: '/a-propos' },
  ]

  return (
    <>
      <Seo
        title="À propos d’Alexandre Iwanesko, Maître FIDE"
        description="Parcours d’Alexandre Iwanesko, Maître FIDE et coach d’échecs à Genève : titre, méthode d’enseignement et philosophie de progression."
        path="/a-propos"
        jsonLd={[personSchema(), breadcrumbSchema(crumbs)]}
      />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero
        eyebrow="À propos"
        title="Alexandre Iwanesko, Maître FIDE & coach d’échecs à Genève"
        lead={
          <>
            {SITE.person.description} Mon objectif&nbsp;: rendre votre progression{' '}
            <strong>concrète et mesurable</strong>.
          </>
        }
        primaryCta={{ to: '/contact', label: 'Me contacter' }}
        secondaryCta={{ to: '/resultats', label: 'Voir les résultats' }}
      />

      <Section>
        <Container className="grid gap-12 lg:grid-cols-[1.5fr_1fr] lg:items-start">
          <div>
          <div className="prose">
            <h2>Un Maître FIDE au service de votre jeu</h2>
            <p>
              <strong>Un Maître FIDE</strong> est un joueur ayant obtenu un titre international
              décerné par la Fédération internationale des échecs (FIDE), reconnaissant un niveau
              Elo élevé et stable. Ce titre est le socle d’une expertise que je mets au service de
              votre progression.
            </p>
            <p>
              J’enseigne aux <strong>adultes de 1200 à 1900 Elo</strong> et aux{' '}
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
          </div>

          <div className="mt-10 rounded-2xl border border-ink-200/80 bg-cream-100 p-7 shadow-soft">
            <Eyebrow>Retrouvez-moi</Eyebrow>
            <p className="text-ink-600">
              Profil FIDE, parties commentées et jeu en ligne — les liens officiels (à compléter)
              renforcent la vérifiabilité de mon expertise.
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
              Réserver un premier cours
            </CtaLink>
          </div>
          </div>

          {/* Visual aside — brand motif + credentials at a glance. */}
          <aside className="lg:sticky lg:top-24">
            <div className="overflow-hidden rounded-3xl border border-ink-200/80 bg-white p-6 shadow-card">
              <ChessMotif className="mx-auto w-full max-w-[16rem] text-ink-900" />
              <dl className="mt-6 grid grid-cols-2 gap-x-4 gap-y-4">
                {[
                  { label: 'Titre', value: 'Maître FIDE' },
                  { label: 'Public', value: '1200–1900 Elo' },
                  { label: 'Lieu', value: 'Genève / en ligne' },
                  { label: 'Langues', value: 'FR · EN' },
                ].map((f) => (
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
