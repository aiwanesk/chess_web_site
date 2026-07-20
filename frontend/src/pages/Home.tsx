import type { ReactNode } from 'react'
import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section, Eyebrow, CtaLink } from '../components/ui'
import { PageHero, FactCard } from '../components/PageHero'
import { ChessMotif } from '../components/ChessMotif'
import { IconGraduation, IconTrophy, IconMonitor, IconUsers, IconSpark, IconBoard, IconCheck, IconArrowRight } from '../components/icons'
import { personSchema, localBusinessSchema, faqSchema, breadcrumbSchema, type FaqItem } from '../lib/schema'
import { useLocale, homePath, t, PAGES, type Locale } from '../lib/i18n'

type Offer = { fr: string; en: string; key?: string; icon: typeof IconGraduation; desc: { fr: string; en: string }; title: { fr: string; en: string } }

const OFFERS: Offer[] = [
  { fr: '/cours-echecs-adultes-geneve', en: PAGES.coursAdultes.en, key: 'coursAdultes', icon: IconGraduation, title: { fr: 'Cours pour adultes', en: 'Adult lessons' }, desc: { fr: 'Progression structurée pour joueurs 1200–1900 Elo.', en: 'Structured progress for 1200–1900 Elo players.' } },
  { fr: '/preparation-tournoi-echecs', en: '/preparation-tournoi-echecs', icon: IconTrophy, title: { fr: 'Préparation tournoi', en: 'Tournament prep' }, desc: { fr: 'Répertoire, gestion du temps et mental de compétition.', en: 'Repertoire, time management and competitive mindset.' } },
  { fr: '/cours-echecs-en-ligne', en: '/cours-echecs-en-ligne', icon: IconMonitor, title: { fr: 'Cours en ligne', en: 'Online lessons' }, desc: { fr: 'Même méthode à distance, partout en Suisse romande.', en: 'Same method remotely, anywhere in French-speaking Switzerland.' } },
  { fr: '/cours-echecs-groupe-geneve', en: '/cours-echecs-groupe-geneve', icon: IconUsers, title: { fr: 'Cours en groupe', en: 'Group lessons' }, desc: { fr: 'Petits groupes de niveau homogène, tarif partagé.', en: 'Small level-matched groups, shared rate.' } },
  { fr: '/cours-echecs-ados-competition', en: '/cours-echecs-ados-competition', icon: IconSpark, title: { fr: 'Ados en compétition', en: 'Competitive teens' }, desc: { fr: 'Coaching orienté progression Elo et résultats.', en: 'Coaching focused on Elo progress and results.' } },
  { fr: '/stages-echecs-geneve', en: '/stages-echecs-geneve', icon: IconBoard, title: { fr: 'Stages à Genève', en: 'Camps in Geneva' }, desc: { fr: 'Sessions intensives pendant les vacances scolaires.', en: 'Intensive sessions during school holidays.' } },
]

const FAQ: Record<Locale, FaqItem[]> = {
  fr: [
    { question: 'À qui s’adressent les cours d’Alexandre Iwanesko ?', answer: 'Aux adultes de 1200 à 1900 Elo et aux adolescents en compétition. Ce sont des cours de perfectionnement : les grands débutants ne sont pas le public visé.' },
    { question: 'Qui est Alexandre Iwanesko ?', answer: 'Alexandre Iwanesko est Maître FIDE et coach d’échecs à Genève. Il enseigne aux joueurs intermédiaires et avancés, en présentiel et en ligne.' },
    { question: 'Les cours ont-ils lieu à Genève ou en ligne ?', answer: 'Les deux : en présentiel à Genève et en ligne pour toute la Suisse romande et la France voisine.' },
  ],
  en: [
    { question: 'Who are Alexandre Iwanesko’s lessons for?', answer: 'For adults from 1200 to 1900 Elo and competitive teenagers. These are improvement lessons: complete beginners are not the target audience.' },
    { question: 'Who is Alexandre Iwanesko?', answer: 'Alexandre Iwanesko is a FIDE Master and chess coach in Geneva. He coaches intermediate and advanced players, in person and online.' },
    { question: 'Are lessons in Geneva or online?', answer: 'Both: in person in Geneva and online across French-speaking Switzerland and neighbouring France.' },
  ],
}

const T: Record<Locale, {
  metaTitle: string; metaDesc: string; heroEyebrow: string; heroTitle: string
  formats: string; formatsTitle: string; whyEyebrow: string; whyTitle: string; whyLead: ReactNode; whyCta: string
  bullets: string[]; faqTitle: string; ctaTitle: string; ctaLead: string; ctaBtn: string
  facts: { label: string; value: string }[]
}> = {
  fr: {
    metaTitle: 'Coach d’échecs à Genève — Alexandre Iwanesko, Maître FIDE',
    metaDesc: 'Coach d’échecs à Genève, Maître FIDE. Cours pour adultes (1200–1900 Elo) et ados en compétition, préparation tournoi, en présentiel et en ligne.',
    heroEyebrow: 'Maître FIDE · Genève & en ligne', heroTitle: 'Progressez aux échecs avec un Maître FIDE à Genève',
    formats: 'Formats', formatsTitle: 'Trouvez le cours qui vous correspond',
    whyEyebrow: 'Pourquoi un Maître FIDE', whyTitle: 'Une expertise de titre, une pédagogie de coach',
    whyLead: (<>Le titre de <strong className="text-ink-900">Maître FIDE</strong> garantit un niveau de jeu élevé. Mais progresser, c’est surtout une méthode : diagnostic précis, plan de travail personnalisé et exercices ciblés entre les séances.</>),
    whyCta: 'En savoir plus sur mon parcours',
    bullets: ['Diagnostic de vos parties classées pour cibler les vrais leviers.', 'Plan de progression sur 8 à 12 semaines avec objectifs Elo.', 'Présentiel à Genève ou en ligne, au même niveau d’exigence.', 'Public sérieux : joueurs qui veulent progresser, pas s’initier.'],
    faqTitle: 'Questions fréquentes', ctaTitle: 'Envie de franchir un palier ?',
    ctaLead: 'Réservez un premier cours à Genève ou en ligne pour définir vos objectifs.', ctaBtn: 'Réserver un premier cours',
    facts: [{ label: 'Titre', value: 'Maître FIDE' }, { label: 'Public', value: '1200–1900 Elo' }, { label: 'Lieu', value: 'Genève / en ligne' }, { label: 'Langues', value: 'FR · EN' }],
  },
  en: {
    metaTitle: 'Chess coach in Geneva — Alexandre Iwanesko, FIDE Master',
    metaDesc: 'Chess coach in Geneva, FIDE Master. Lessons for adults (1200–1900 Elo) and competitive teens, tournament prep, in person and online.',
    heroEyebrow: 'FIDE Master · Geneva & online', heroTitle: 'Improve your chess with a FIDE Master in Geneva',
    formats: 'Formats', formatsTitle: 'Find the lesson that fits you',
    whyEyebrow: 'Why a FIDE Master', whyTitle: 'A titled expertise, a coach’s teaching',
    whyLead: (<>The <strong className="text-ink-900">FIDE Master</strong> title guarantees a high level of play. But improving is mostly a method: precise diagnosis, a personalised work plan and targeted exercises between sessions.</>),
    whyCta: 'More about my background',
    bullets: ['Diagnosis of your rated games to target the real levers.', 'An 8–12 week progression plan with Elo goals.', 'In person in Geneva or online, at the same standard.', 'A serious audience: players who want to improve, not to start out.'],
    faqTitle: 'Frequently asked questions', ctaTitle: 'Ready to reach the next level?',
    ctaLead: 'Book a first lesson in Geneva or online to define your goals.', ctaBtn: 'Book a first lesson',
    facts: [{ label: 'Title', value: 'FIDE Master' }, { label: 'Level', value: '1200–1900 Elo' }, { label: 'Where', value: 'Geneva / online' }, { label: 'Languages', value: 'EN · FR' }],
  },
}

export function Component() {
  const locale = useLocale()
  const c = T[locale]
  const aboutPath = locale === 'en' ? '/a-propos' : '/a-propos' // à-propos EN à venir
  const contactPath = PAGES.contact[locale]
  const jsonLd = [personSchema(), localBusinessSchema(), breadcrumbSchema([{ name: t(locale).breadcrumbHome, path: homePath(locale) }]), faqSchema(FAQ[locale])]

  return (
    <>
      <Seo isHome title={c.metaTitle} description={c.metaDesc} path={homePath(locale)} jsonLd={jsonLd} />

      <PageHero
        eyebrow={c.heroEyebrow}
        title={c.heroTitle}
        lead={locale === 'en'
          ? (<>Chess coaching for <strong>adults (1200–1900 Elo)</strong> and <strong>competitive teens</strong>. A clear method, a progression plan and real follow-up — in person in Geneva or online.</>)
          : (<>Coaching d’échecs pour <strong>adultes (1200–1900 Elo)</strong> et <strong>ados en compétition</strong>. Une méthode claire, un plan de progression et un vrai suivi — en présentiel à Genève ou en ligne.</>)}
        primaryCta={{ to: contactPath, label: c.ctaBtn }}
        secondaryCta={{ to: locale === 'en' ? PAGES.coursAdultes.en : '/cours-echecs-adultes-geneve', label: locale === 'en' ? 'Discover adult lessons' : 'Découvrir les cours adultes' }}
        aside={<div className="relative"><ChessMotif className="mx-auto w-56 text-ink-900 drop-shadow-xl sm:w-64" /><div className="mt-6"><FactCard facts={c.facts} /></div></div>}
      />

      <Section>
        <Container>
          <Eyebrow>{c.formats}</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">{c.formatsTitle}</h2>
          <div className="mt-10 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {OFFERS.map((o) => (
              <Link key={o.fr} to={locale === 'en' ? o.en : o.fr} className="hover-lift group flex flex-col rounded-2xl border border-ink-200/80 bg-white p-7 shadow-soft transition-colors hover:border-gold-300 hover:shadow-card">
                <span aria-hidden className="flex h-12 w-12 items-center justify-center rounded-xl border border-gold-200 bg-gold-50 text-gold-700"><o.icon size={24} /></span>
                <h3 className="mt-5 flex items-center justify-between text-lg font-semibold text-ink-900">
                  {o.title[locale]}
                  <span aria-hidden className="text-gold-600 transition-transform group-hover:translate-x-1"><IconArrowRight size={18} /></span>
                </h3>
                <p className="mt-2 leading-relaxed text-ink-600">{o.desc[locale]}</p>
              </Link>
            ))}
          </div>
        </Container>
      </Section>

      <Section className="border-y border-ink-100 bg-cream-100">
        <Container className="grid gap-12 lg:grid-cols-2 lg:items-center">
          <div>
            <Eyebrow>{c.whyEyebrow}</Eyebrow>
            <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">{c.whyTitle}</h2>
            <p className="mt-5 text-lg leading-relaxed text-ink-600">{c.whyLead}</p>
            <div className="mt-8"><CtaLink to={aboutPath} variant="secondary">{c.whyCta}</CtaLink></div>
          </div>
          <ul className="space-y-3">
            {c.bullets.map((b) => (
              <li key={b} className="flex items-start gap-4 rounded-2xl border border-ink-200/80 bg-white p-5 shadow-soft">
                <span aria-hidden className="mt-0.5 flex h-7 w-7 flex-none items-center justify-center rounded-full bg-gold-100 text-gold-700"><IconCheck size={16} /></span>
                <span className="leading-relaxed text-ink-700">{b}</span>
              </li>
            ))}
          </ul>
        </Container>
      </Section>

      <Section>
        <Container>
          <Eyebrow>FAQ</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">{c.faqTitle}</h2>
          <div className="mx-auto mt-8 max-w-4xl space-y-3">
            {FAQ[locale].map((f) => (
              <details key={f.question} className="group rounded-2xl border border-ink-200/80 bg-white px-6 py-4 shadow-soft">
                <summary className="flex cursor-pointer list-none items-center justify-between gap-4 text-lg font-semibold text-ink-900">
                  {f.question}
                  <span aria-hidden className="flex h-7 w-7 flex-none items-center justify-center rounded-full bg-gold-100 text-gold-700 transition-transform group-open:rotate-45">+</span>
                </summary>
                <p className="mt-3 leading-relaxed text-ink-600">{f.answer}</p>
              </details>
            ))}
          </div>
        </Container>
      </Section>

      <Section className="board-texture-dark bg-ink-950 text-center text-white">
        <Container>
          <span aria-hidden className="mx-auto mb-6 block h-0.5 w-12 rounded-full bg-gold-500" />
          <h2 className="font-display text-3xl font-bold sm:text-4xl">{c.ctaTitle}</h2>
          <p className="mx-auto mt-4 max-w-xl leading-relaxed text-ink-300">{c.ctaLead}</p>
          <div className="mt-9"><CtaLink to={contactPath} variant="primary">{c.ctaBtn}</CtaLink></div>
        </Container>
      </Section>
    </>
  )
}
