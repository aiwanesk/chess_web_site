import { MoneyPage, type MoneyPageProps } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'
import { useLocale } from '../lib/i18n'

function FormatsFR() {
  return (
    <DeepDive eyebrow="Formats" title="Trois manières d’intervenir">
      <StepList
        steps={[
          {
            t: 'Conférence',
            d: 'Un propos structuré (45–60 min) reliant les échecs à vos enjeux de décision et de stratégie, illustré d’exemples célèbres.',
          },
          {
            t: 'Atelier interactif',
            d: 'Les participants jouent et expérimentent les notions : anticipation, plans, gestion du risque, en équipes.',
          },
          {
            t: 'Simultanée',
            d: 'Le conférencier affronte plusieurs collaborateurs à la fois — un moment marquant qui clôt idéalement un séminaire.',
          },
        ]}
      />
    </DeepDive>
  )
}

function FormatsEN() {
  return (
    <DeepDive eyebrow="Formats" title="Three ways to deliver">
      <StepList
        steps={[
          {
            t: 'Keynote',
            d: 'A structured talk (45–60 min) connecting chess to your decision-making and strategy challenges, illustrated with famous examples.',
          },
          {
            t: 'Interactive workshop',
            d: 'Participants play and experience the concepts first-hand: anticipation, plans and risk management, in teams.',
          },
          {
            t: 'Simultaneous exhibition',
            d: 'The speaker takes on several colleagues at once — a memorable moment that makes the ideal close to a seminar.',
          },
        ]}
      />
    </DeepDive>
  )
}

const FR: MoneyPageProps = {
  path: '/conferences-echecs-entreprise',
  title: 'Conférences d’échecs en entreprise',
  metaTitle: 'Conférences échecs en entreprise — Genève & lémanique',
  description:
    'Conférences et interventions échecs en entreprise à Genève : stratégie, prise de décision et gestion du risque, par un Maître FIDE.',
  eyebrow: 'Entreprises',
  lead: (
    <>
      Une intervention inspirante qui relie les <strong>échecs</strong> aux enjeux de
      l’entreprise : <strong>stratégie</strong>, <strong>prise de décision sous pression</strong>{' '}
      et gestion du risque, animée par un Maître&nbsp;FIDE.
    </>
  ),
  facts: [
    { label: 'Format', value: 'Conférence / atelier' },
    { label: 'Durée', value: '45–90 min' },
    { label: 'Public', value: 'Équipes, comités' },
    { label: 'Lieu', value: 'Vos locaux / Genève' },
  ],
  intro: (
    <>
      <p>
        <strong>Une conférence d’échecs en entreprise</strong> utilise le jeu comme métaphore
        concrète des décisions stratégiques : anticipation, évaluation des risques, gestion du
        temps et lucidité sous pression. Le propos est vivant, illustré de positions célèbres
        et adapté à votre secteur.
      </p>
      <p>
        L’intervention peut prendre la forme d’une conférence, d’un atelier interactif ou d’une
        simultanée où le conférencier affronte plusieurs collaborateurs à la fois.
      </p>
    </>
  ),
  benefits: [
    { title: 'Message sur mesure', body: 'Un fil conducteur aligné sur vos enjeux : décision, risque, stratégie.' },
    { title: 'Format flexible', body: 'Conférence, atelier participatif ou simultanée, selon l’effet recherché.' },
    { title: 'Interactif', body: 'Des exemples concrets et des moments de participation qui marquent les esprits.' },
    { title: 'Clé en main', body: 'Organisation simple, dans vos locaux à Genève ou dans l’arc lémanique.' },
  ],
  course: {
    name: 'Conférence d’échecs en entreprise',
    description:
      'Conférence ou atelier échecs en entreprise sur la stratégie et la prise de décision, animé par un Maître FIDE à Genève.',
    url: '/conferences-echecs-entreprise',
    courseMode: 'onsite',
  },
  related: [
    { to: '/team-building-echecs-geneve', label: 'Team building échecs' },
    { to: '/contact', label: 'Demander un devis' },
  ],
  faq: [
    {
      question: 'Combien de temps dure une conférence ?',
      answer:
        'De 45 à 90 minutes selon le format choisi, avec possibilité d’ajouter un atelier ou une simultanée.',
    },
    {
      question: 'Faut-il savoir jouer aux échecs pour en profiter ?',
      answer:
        'Non. Le propos est accessible à tous ; les notions échiquéennes sont expliquées et servent de support aux idées de stratégie et de décision.',
    },
    {
      question: 'Intervenez-vous en dehors de Genève ?',
      answer:
        'Oui, dans tout l’arc lémanique et au-delà sur demande. Contactez-moi pour un devis adapté.',
    },
    {
      question: 'Quels thèmes pouvez-vous relier aux échecs ?',
      answer:
        'Prise de décision sous incertitude, anticipation, gestion du risque et du temps, apprentissage de l’erreur, calme sous pression. Le fil est adapté à votre secteur et à votre message.',
    },
    {
      question: 'Combien de participants pour une simultanée ?',
      answer:
        'Une simultanée peut opposer le conférencier à une dizaine, voire une vingtaine de collaborateurs en parallèle — un format spectaculaire et fédérateur.',
    },
  ],
  children: <FormatsFR />,
}

const EN: MoneyPageProps = {
  path: '/en/corporate-chess-talks',
  title: 'Corporate chess keynotes & talks',
  metaTitle: 'Corporate chess talks — Geneva & Lake Geneva region',
  description:
    'Corporate chess keynotes and talks in Geneva: strategy, decision-making and risk management, delivered by a FIDE Master.',
  eyebrow: 'For companies',
  lead: (
    <>
      An inspiring keynote that connects <strong>chess</strong> to your business challenges:{' '}
      <strong>strategy</strong>, <strong>decision-making under pressure</strong> and risk
      management, delivered by a FIDE&nbsp;Master.
    </>
  ),
  facts: [
    { label: 'Format', value: 'Keynote / workshop' },
    { label: 'Duration', value: '45–90 min' },
    { label: 'Audience', value: 'Teams, boards' },
    { label: 'Location', value: 'Your offices / Geneva' },
  ],
  intro: (
    <>
      <p>
        <strong>A corporate chess talk</strong> uses the game as a concrete metaphor for strategic
        decisions: anticipation, risk assessment, time management and clear thinking under
        pressure. The talk is lively, illustrated with famous positions and tailored to your
        industry.
      </p>
      <p>
        It can take the form of a keynote, an interactive workshop or a simultaneous exhibition in
        which the speaker takes on several colleagues at once.
      </p>
    </>
  ),
  benefits: [
    { title: 'Tailored message', body: 'A narrative aligned with your priorities: decision-making, risk and strategy.' },
    { title: 'Flexible format', body: 'Keynote, hands-on workshop or simultaneous exhibition, depending on the effect you want.' },
    { title: 'Interactive', body: 'Concrete examples and moments of participation that leave a lasting impression.' },
    { title: 'Turnkey', body: 'Simple to organise, at your offices in Geneva or across the Lake Geneva region.' },
  ],
  course: {
    name: 'Corporate chess talk',
    description:
      'A corporate chess keynote or workshop on strategy and decision-making, delivered by a FIDE Master in Geneva.',
    url: '/en/corporate-chess-talks',
    courseMode: 'onsite',
  },
  related: [
    { to: '/en/chess-team-building-geneva', label: 'Chess team building' },
    { to: '/en/contact', label: 'Request a quote' },
  ],
  faq: [
    {
      question: 'How long does a talk last?',
      answer:
        '45 to 90 minutes depending on the format chosen, with the option to add a workshop or a simultaneous exhibition.',
    },
    {
      question: 'Do you need to know how to play chess to enjoy it?',
      answer:
        'No. The talk is accessible to everyone; the chess concepts are explained and used to support the ideas around strategy and decision-making.',
    },
    {
      question: 'Do you travel outside Geneva?',
      answer:
        'Yes, throughout the Lake Geneva region and beyond on request. Get in touch for a tailored quote.',
    },
    {
      question: 'Which themes can you connect to chess?',
      answer:
        'Decision-making under uncertainty, anticipation, risk and time management, learning from mistakes, staying calm under pressure. The narrative is tailored to your industry and your message.',
    },
    {
      question: 'How many participants for a simultaneous exhibition?',
      answer:
        'A simultaneous exhibition can pit the speaker against ten or even twenty colleagues at the same time — a spectacular and unifying format.',
    },
  ],
  children: <FormatsEN />,
}

export function Component() {
  const locale = useLocale()
  return <MoneyPage {...(locale === 'en' ? EN : FR)} />
}
