import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section, Eyebrow, CtaLink } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { Faq } from '../components/Faq'
import { IconCheck } from '../components/icons'
import { breadcrumbSchema, faqSchema, type Crumb, type FaqItem } from '../lib/schema'
import { useLocale, homePath, t, PAGES, type Locale } from '../lib/i18n'

interface Plan { name: string; price: string; unit: string; features: string[]; cta: string; featured: boolean }
interface TData {
  path: string; title: string; description: string; eyebrow: string; heroTitle: string; heroLead: string
  ask: string; formules: string; chooseTitle: string; enterprise: string; faqTitle: string; featured: string
  plans: Plan[]; faq: FaqItem[]
}

const DATA: Record<Locale, TData> = {
  fr: {
    path: '/tarifs', title: 'Tarifs des cours d’échecs',
    description: 'Tarifs des cours d’échecs à Genève : cours particuliers, en groupe, en ligne et forfaits. Prix clairs en CHF.',
    eyebrow: 'Tarifs', heroTitle: 'Des tarifs clairs, en francs suisses',
    heroLead: 'Cours particuliers, forfaits et cours en groupe. Choisissez le format qui correspond à vos objectifs et à votre rythme.',
    ask: 'Poser une question', formules: 'Formules', chooseTitle: 'Choisissez votre formule',
    enterprise: 'Conférences et team building en entreprise sur devis.',
    faqTitle: 'Questions sur les tarifs', featured: '★ Le plus choisi',
    plans: [
      { name: 'Séance individuelle', price: '120 CHF', unit: '/ 60 min', features: ['Diagnostic personnalisé', 'Présentiel Genève ou en ligne', 'Exercices entre séances'], cta: 'Réserver une séance', featured: false },
      { name: 'Pack 10 séances', price: '1000 CHF', unit: '10 × 60 min · soit 100 CHF/séance', features: ['Plan de progression complet', 'Suivi Elo & points de contrôle', 'Tarif horaire réduit'], cta: 'Choisir le pack', featured: true },
      { name: 'Cours en groupe', price: '60 CHF', unit: '/ pers. / 90 min', features: ['Petit groupe de niveau', 'Émulation collective', 'À Genève'], cta: 'Rejoindre un groupe', featured: false },
    ],
    faq: [
      { question: 'Quels sont les moyens de paiement acceptés ?', answer: 'Paiement par TWINT, virement bancaire ou espèces. Les forfaits sont réglés à l’avance ; le détail est confirmé lors du premier contact.' },
      { question: 'Les forfaits ont-ils une durée de validité ?', answer: 'Oui, un forfait de 10 séances est en général à utiliser sur 4 à 6 mois. Les modalités précises sont convenues ensemble.' },
      { question: 'Proposez-vous une première séance découverte ?', answer: 'Le premier contact permet de définir vos objectifs et le format adapté. Contactez-moi pour en discuter sans engagement.' },
    ],
  },
  en: {
    path: '/en/pricing', title: 'Chess lesson pricing',
    description: 'Chess lesson pricing in Geneva: private lessons, group lessons, online and packages. Clear prices in CHF.',
    eyebrow: 'Pricing', heroTitle: 'Clear pricing, in Swiss francs',
    heroLead: 'Private lessons, packages and group lessons. Choose the format that matches your goals and your pace.',
    ask: 'Ask a question', formules: 'Plans', chooseTitle: 'Choose your plan',
    enterprise: 'Corporate talks and team building on request.',
    faqTitle: 'Pricing questions', featured: '★ Most chosen',
    plans: [
      { name: 'Single session', price: '120 CHF', unit: '/ 60 min', features: ['Personalised diagnosis', 'In person in Geneva or online', 'Exercises between sessions'], cta: 'Book a session', featured: false },
      { name: '10-session package', price: '1000 CHF', unit: '10 × 60 min · i.e. 100 CHF/session', features: ['Full progression plan', 'Elo tracking & checkpoints', 'Reduced hourly rate'], cta: 'Choose the package', featured: true },
      { name: 'Group lesson', price: '60 CHF', unit: '/ person / 90 min', features: ['Small level-matched group', 'Collective momentum', 'In Geneva'], cta: 'Join a group', featured: false },
    ],
    faq: [
      { question: 'What payment methods are accepted?', answer: 'Payment by TWINT, bank transfer or cash. Packages are paid in advance; details are confirmed at first contact.' },
      { question: 'Do packages have an expiry?', answer: 'Yes, a 10-session package is usually to be used within 4 to 6 months. Exact terms are agreed together.' },
      { question: 'Do you offer a first discovery session?', answer: 'The first contact defines your goals and the right format. Get in touch to discuss with no commitment.' },
    ],
  },
}

export function Component() {
  const locale = useLocale()
  const d = DATA[locale]
  const contactPath = PAGES.contact[locale]
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: d.eyebrow, path: d.path },
  ]

  return (
    <>
      <Seo title={d.title} description={d.description} path={d.path} jsonLd={[breadcrumbSchema(crumbs), faqSchema(d.faq)]} />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero eyebrow={d.eyebrow} title={d.heroTitle} lead={d.heroLead} secondaryCta={{ to: contactPath, label: d.ask }} />

      <Section>
        <Container>
          <Eyebrow>{d.formules}</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">{d.chooseTitle}</h2>
          <div className="mt-10 grid items-start gap-6 lg:grid-cols-3">
            {d.plans.map((p) => (
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
                    {d.featured}
                  </span>
                ) : null}
                <h3 className="font-display text-xl font-bold text-ink-900">{p.name}</h3>
                <p className="mt-4 flex items-baseline gap-2">
                  <span className="font-display text-4xl font-extrabold tracking-tight text-ink-900">{p.price}</span>
                  <span className="text-sm text-ink-500">{p.unit}</span>
                </p>
                <div className="my-6 h-px bg-ink-100" />
                <ul className="flex-1 space-y-3 text-sm text-ink-700">
                  {p.features.map((f) => (
                    <li key={f} className="flex items-start gap-3">
                      <span aria-hidden className="mt-px flex h-5 w-5 flex-none items-center justify-center rounded-full bg-gold-100 text-gold-700">
                        <IconCheck size={13} />
                      </span>
                      {f}
                    </li>
                  ))}
                </ul>
                <div className="mt-8">
                  <CtaLink to={contactPath} variant={p.featured ? 'primary' : 'secondary'} className="w-full">
                    {p.cta}
                  </CtaLink>
                </div>
              </div>
            ))}
          </div>
          <p className="mt-8 text-center text-sm text-ink-500">{d.enterprise}</p>
        </Container>
      </Section>

      <Faq items={d.faq} title={d.faqTitle} />
    </>
  )
}
