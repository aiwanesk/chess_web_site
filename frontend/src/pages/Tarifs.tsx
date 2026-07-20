import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section, Eyebrow, CtaLink } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { Faq } from '../components/Faq'
import { IconCheck } from '../components/icons'
import { breadcrumbSchema, faqSchema, type Crumb } from '../lib/schema'

const PLANS = [
  {
    name: 'Séance individuelle',
    price: '120 CHF',
    unit: '/ 60 min',
    features: ['Diagnostic personnalisé', 'Présentiel Genève ou en ligne', 'Exercices entre séances'],
    cta: 'Réserver une séance',
    featured: false,
  },
  {
    name: 'Pack 10 séances',
    price: '1000 CHF',
    unit: '10 × 60 min · soit 100 CHF/séance',
    features: ['Plan de progression complet', 'Suivi Elo & points de contrôle', 'Tarif horaire réduit'],
    cta: 'Choisir le pack',
    featured: true,
  },
  {
    name: 'Cours en groupe',
    price: '60 CHF',
    unit: '/ pers. / 90 min',
    features: ['Petit groupe de niveau', 'Émulation collective', 'À Genève'],
    cta: 'Rejoindre un groupe',
    featured: false,
  },
]

const FAQ = [
  {
    question: 'Quels sont les moyens de paiement acceptés ?',
    answer:
      'Paiement par TWINT, virement bancaire ou espèces. Les forfaits sont réglés à l’avance ; le détail est confirmé lors du premier contact.',
  },
  {
    question: 'Les forfaits ont-ils une durée de validité ?',
    answer:
      'Oui, un forfait de 10 séances est en général à utiliser sur 4 à 6 mois. Les modalités précises sont convenues ensemble.',
  },
  {
    question: 'Proposez-vous une première séance découverte ?',
    answer:
      'Le premier contact permet de définir vos objectifs et le format adapté. Contactez-moi pour en discuter sans engagement.',
  },
]

export function Component() {
  const crumbs: Crumb[] = [
    { name: 'Accueil', path: '/' },
    { name: 'Tarifs', path: '/tarifs' },
  ]

  return (
    <>
      <Seo
        title="Tarifs des cours d’échecs"
        description="Tarifs des cours d’échecs à Genève : cours particuliers, en groupe, en ligne et forfaits de préparation tournoi. Prix clairs en CHF."
        path="/tarifs"
        jsonLd={[breadcrumbSchema(crumbs), faqSchema(FAQ)]}
      />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero
        eyebrow="Tarifs"
        title="Des tarifs clairs, en francs suisses"
        lead="Cours particuliers, forfaits et cours en groupe. Choisissez le format qui correspond à vos objectifs et à votre rythme."
        secondaryCta={{ to: '/contact', label: 'Poser une question' }}
      />

      <Section>
        <Container>
          <Eyebrow>Formules</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">
            Choisissez votre formule
          </h2>
          <div className="mt-10 grid items-start gap-6 lg:grid-cols-3">
            {PLANS.map((p) => (
              <div
                key={p.name}
                className={`relative flex flex-col rounded-3xl p-8 transition-shadow ${
                  p.featured
                    ? 'border-2 border-gold-400 bg-white shadow-lift lg:-mt-4 lg:pt-12'
                    : 'border border-ink-200/80 bg-white shadow-soft hover:shadow-card'
                }`}
              >
                {p.featured ? (
                  <span className="absolute -top-3.5 left-1/2 inline-flex -translate-x-1/2 items-center gap-1.5 rounded-full bg-gold-500 px-4 py-1.5 text-xs font-bold uppercase tracking-wide text-ink-950 shadow-gold">
                    ★ Le plus choisi
                  </span>
                ) : null}
                <h3 className="font-display text-xl font-bold text-ink-900">{p.name}</h3>
                <p className="mt-4 flex items-baseline gap-2">
                  <span className="font-display text-4xl font-extrabold tracking-tight text-ink-900">
                    {p.price}
                  </span>
                  <span className="text-sm text-ink-500">{p.unit}</span>
                </p>
                <div className="my-6 h-px bg-ink-100" />
                <ul className="flex-1 space-y-3 text-sm text-ink-700">
                  {p.features.map((f) => (
                    <li key={f} className="flex items-start gap-3">
                      <span
                        aria-hidden
                        className="mt-px flex h-5 w-5 flex-none items-center justify-center rounded-full bg-gold-100 text-gold-700"
                      >
                        <IconCheck size={13} />
                      </span>
                      {f}
                    </li>
                  ))}
                </ul>
                <div className="mt-8">
                  <CtaLink to="/contact" variant={p.featured ? 'primary' : 'secondary'} className="w-full">
                    {p.cta}
                  </CtaLink>
                </div>
              </div>
            ))}
          </div>
          <p className="mt-8 text-center text-sm text-ink-500">
            Conférences et team building en entreprise sur devis.
          </p>
        </Container>
      </Section>

      <Faq items={FAQ} title="Questions sur les tarifs" />
    </>
  )
}
