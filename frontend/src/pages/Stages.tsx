import { MoneyPage, type MoneyPageProps } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'
import { eventSchema } from '../lib/schema'
import { useLocale } from '../lib/i18n'

function ProgrammeFR() {
  return (
    <DeepDive eyebrow="Le programme" title="Une journée type de stage">
      <StepList
        steps={[
          {
            t: 'Matin — Thème stratégique',
            d: 'Cours interactif sur une notion clé (structures de pions, attaque du roi, finales), illustré par des parties de maîtres.',
          },
          {
            t: 'Midi — Résolution d’exercices',
            d: 'Ateliers de calcul et de tactique par niveau, en individuel et en équipe.',
          },
          {
            t: 'Après-midi — Parties & analyse',
            d: 'Parties d’entraînement avec cadence, puis analyse collective à chaud des positions marquantes.',
          },
          {
            t: 'Fin de stage — Bilan',
            d: 'Un retour individuel avec des axes de travail personnalisés pour la suite.',
          },
        ]}
      />
    </DeepDive>
  )
}

function ProgrammeEN() {
  return (
    <DeepDive eyebrow="The programme" title="A typical camp day">
      <StepList
        steps={[
          {
            t: 'Morning — Strategic theme',
            d: 'Interactive class on a key concept (pawn structures, king attack, endgames), illustrated with master games.',
          },
          {
            t: 'Midday — Solving exercises',
            d: 'Calculation and tactics workshops by level, individually and in teams.',
          },
          {
            t: 'Afternoon — Games & analysis',
            d: 'Timed training games, followed by a collective on-the-spot analysis of the key positions.',
          },
          {
            t: 'End of camp — Review',
            d: 'An individual debrief with personalised areas of work for the future.',
          },
        ]}
      />
    </DeepDive>
  )
}

const FR: MoneyPageProps = {
  path: '/stages-echecs-geneve',
  title: 'Stages d’échecs à Genève',
  metaTitle: 'Stages d’échecs à Genève — vacances, Maître FIDE',
  description:
    'Stages d’échecs intensifs à Genève pendant les vacances scolaires, encadrés par un Maître FIDE. Progression rapide sur quelques jours.',
  eyebrow: 'Vacances scolaires · Genève',
  lead: (
    <>
      Des <strong>stages intensifs</strong> pendant les vacances pour progresser vite :
      plusieurs jours de travail structuré, d’analyse et de parties, encadrés par un{' '}
      <strong>Maître&nbsp;FIDE</strong>.
    </>
  ),
  facts: [
    { label: 'Durée', value: '2 à 5 jours' },
    { label: 'Rythme', value: 'Intensif' },
    { label: 'Niveau', value: 'Groupes par niveau' },
    { label: 'Lieu', value: 'Genève' },
  ],
  intro: (
    <>
      <p>
        <strong>Un stage d’échecs</strong> est une session intensive de quelques jours qui
        concentre un maximum de travail utile : thèmes stratégiques, exercices de calcul,
        finales et parties commentées. Le format vacances permet de progresser nettement en
        peu de temps.
      </p>
      <p>
        Les stages sont organisés par niveau pour garantir un rythme adapté à chaque
        participant.
      </p>
    </>
  ),
  benefits: [
    { title: 'Immersion', body: 'Plusieurs jours dédiés aux échecs pour ancrer durablement les progrès.' },
    { title: 'Programme dense', body: 'Stratégie, calcul, finales et parties commentées en alternance.' },
    { title: 'Par niveau', body: 'Des groupes homogènes pour un rythme adapté à chacun.' },
    { title: 'Bilan personnalisé', body: 'Un retour individuel avec des axes de travail pour la suite.' },
  ],
  course: {
    name: 'Stage d’échecs à Genève',
    description:
      'Stage d’échecs intensif à Genève pendant les vacances scolaires, encadré par un Maître FIDE, organisé par niveau.',
    url: '/stages-echecs-geneve',
    price: 240,
    priceUnit: 'le stage de plusieurs jours',
    courseMode: 'onsite',
  },
  extraJsonLd: [
    eventSchema({
      name: 'Stage d’échecs intensif — Genève',
      description:
        'Stage d’échecs de plusieurs jours à Genève pendant les vacances scolaires, encadré par un Maître FIDE.',
      url: '/stages-echecs-geneve',
      // Placeholder dates — mettre à jour à chaque session programmée.
      startDate: '2026-10-19',
      endDate: '2026-10-21',
      price: 240,
    }),
  ],
  related: [
    { to: '/cours-echecs-ados-competition', label: 'Coaching ados en compétition' },
    { to: '/cours-echecs-groupe-geneve', label: 'Cours en groupe à l’année' },
    { to: '/tarifs', label: 'Tarifs' },
  ],
  faq: [
    {
      question: 'Quand ont lieu les stages ?',
      answer:
        'Pendant les vacances scolaires genevoises. Les dates des prochaines sessions sont annoncées sur cette page et par e-mail ; contactez-moi pour être informé.',
    },
    {
      question: 'Les stages sont-ils adaptés aux adultes comme aux jeunes ?',
      answer:
        'Oui, des groupes distincts sont formés par niveau et par public afin que chacun travaille au bon rythme.',
    },
    {
      question: 'Quel niveau faut-il avoir ?',
      answer:
        'Les stages s’adressent à des joueurs qui connaissent déjà les bases et jouent régulièrement. Précisez votre niveau à l’inscription.',
    },
    {
      question: 'Le repas et le matériel sont-ils inclus ?',
      answer:
        'Le matériel d’échecs est fourni. Les modalités pratiques (horaires, repas, pauses) sont précisées dans le programme de chaque session.',
    },
  ],
  children: <ProgrammeFR />,
}

const EN: MoneyPageProps = {
  path: '/en/chess-camps-geneva',
  title: 'Chess camps in Geneva',
  metaTitle: 'Chess camps in Geneva — school holidays, FIDE Master',
  description:
    'Intensive chess camps in Geneva during the school holidays, led by a FIDE Master. Fast progress over a few days.',
  eyebrow: 'School holidays · Geneva',
  lead: (
    <>
      <strong>Intensive camps</strong> during the holidays to progress fast:
      several days of structured work, analysis and games, led by a{' '}
      <strong>FIDE&nbsp;Master</strong>.
    </>
  ),
  facts: [
    { label: 'Duration', value: '2 to 5 days' },
    { label: 'Pace', value: 'Intensive' },
    { label: 'Level', value: 'Groups by level' },
    { label: 'Location', value: 'Geneva' },
  ],
  intro: (
    <>
      <p>
        <strong>A chess camp</strong> is an intensive session of a few days that packs in
        as much useful work as possible: strategic themes, calculation exercises,
        endgames and annotated games. The holiday format makes it possible to progress
        markedly in a short time.
      </p>
      <p>
        Camps are organised by level to ensure a pace suited to every participant.
      </p>
    </>
  ),
  benefits: [
    { title: 'Immersion', body: 'Several days devoted to chess to anchor progress for the long term.' },
    { title: 'Packed programme', body: 'Strategy, calculation, endgames and annotated games in alternation.' },
    { title: 'By level', body: 'Homogeneous groups for a pace suited to everyone.' },
    { title: 'Personalised review', body: 'Individual feedback with areas of work for what comes next.' },
  ],
  course: {
    name: 'Chess camp in Geneva',
    description:
      'Intensive chess camp in Geneva during the school holidays, led by a FIDE Master, organised by level.',
    url: '/en/chess-camps-geneva',
    price: 240,
    priceUnit: 'per multi-day camp',
    courseMode: 'onsite',
  },
  extraJsonLd: [
    eventSchema({
      name: 'Intensive chess camp — Geneva',
      description:
        'Multi-day chess camp in Geneva during the school holidays, led by a FIDE Master.',
      url: '/en/chess-camps-geneva',
      // Placeholder dates — update for each scheduled session.
      startDate: '2026-10-19',
      endDate: '2026-10-21',
      price: 240,
    }),
  ],
  related: [
    { to: '/en/junior-chess-coaching', label: 'Junior competition coaching' },
    { to: '/en/group-chess-lessons-geneva', label: 'Year-round group lessons' },
    { to: '/en/pricing', label: 'Pricing' },
  ],
  faq: [
    {
      question: 'When do the camps take place?',
      answer:
        'During the Geneva school holidays. The dates of upcoming sessions are announced on this page and by email; contact me to be kept informed.',
    },
    {
      question: 'Are the camps suitable for adults as well as young players?',
      answer:
        'Yes, separate groups are formed by level and by audience so everyone works at the right pace.',
    },
    {
      question: 'What level do I need?',
      answer:
        'The camps are for players who already know the basics and play regularly. Please indicate your level when registering.',
    },
    {
      question: 'Are meals and equipment included?',
      answer:
        'Chess equipment is provided. The practical details (schedule, meals, breaks) are set out in the programme for each session.',
    },
  ],
  children: <ProgrammeEN />,
}

export function Component() {
  const locale = useLocale()
  return <MoneyPage {...(locale === 'en' ? EN : FR)} />
}
