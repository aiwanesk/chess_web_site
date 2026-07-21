import { MoneyPage, type MoneyPageProps } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'
import { useLocale } from '../lib/i18n'

function DeroulementFR() {
  return (
    <DeepDive eyebrow="Le déroulé" title="Une animation en trois temps">
      <StepList
        steps={[
          {
            t: 'Échauffement',
            d: 'Les règles utiles en quelques minutes et un premier défi ludique pour mettre tout le monde à l’aise, quel que soit le niveau.',
          },
          {
            t: 'Défis en équipe',
            d: 'Échecs par équipes, variantes coopératives et quiz tactiques : la communication et la stratégie collective au cœur du jeu.',
          },
          {
            t: 'Simultanée conviviale',
            d: 'Les volontaires affrontent le Maître FIDE en simultanée — un final fédérateur et mémorable.',
          },
        ]}
      />
    </DeepDive>
  )
}

function DeroulementEN() {
  return (
    <DeepDive eyebrow="The format" title="A workshop in three acts">
      <StepList
        steps={[
          {
            t: 'Warm-up',
            d: 'The key rules in a few minutes and a first playful challenge to put everyone at ease, whatever their level.',
          },
          {
            t: 'Team challenges',
            d: 'Team chess, cooperative variants and tactical quizzes: communication and collective strategy at the heart of the game.',
          },
          {
            t: 'Friendly simul',
            d: 'Volunteers take on the FIDE Master in a simultaneous display — a memorable, team-uniting finale.',
          },
        ]}
      />
    </DeepDive>
  )
}

const FR: MoneyPageProps = {
  path: '/team-building-echecs-geneve',
  title: 'Team building échecs à Genève',
  metaTitle: 'Team building échecs à Genève — ateliers d’équipe',
  description:
    'Ateliers de team building autour des échecs pour entreprises à Genève et dans l’arc lémanique : coopération, stratégie et fun, par un Maître FIDE.',
  eyebrow: 'Cohésion d’équipe · Genève',
  lead: (
    <>
      Un <strong>team building original</strong> autour des échecs : tournois par équipes,
      échecs coopératifs et défis contre un Maître&nbsp;FIDE — pour souder vos équipes en
      s’amusant intelligemment.
    </>
  ),
  facts: [
    { label: 'Format', value: 'Atelier d’équipe' },
    { label: 'Groupe', value: 'Jusqu’à 30+' },
    { label: 'Durée', value: '1 à 3 h' },
    { label: 'Lieu', value: 'Vos locaux / Genève' },
  ],
  intro: (
    <>
      <p>
        <strong>Un team building échecs</strong> transforme le jeu en activité collective :
        échecs en équipe, variantes coopératives, quiz tactiques et simultanée conviviale.
        Aucun niveau requis — l’objectif est la <strong>coopération</strong>, la communication
        et le plaisir partagé.
      </p>
      <p>
        L’atelier s’adapte à la taille de votre groupe et à vos objectifs, dans vos locaux à
        Genève ou sur un lieu de séminaire.
      </p>
    </>
  ),
  benefits: [
    { title: 'Sans prérequis', body: 'Tout le monde participe, quel que soit le niveau d’échecs.' },
    { title: 'Coopération', body: 'Des formats en équipe qui font travailler la communication et la décision collective.' },
    { title: 'Ludique', body: 'Variantes originales, défis et simultanée pour une ambiance conviviale.' },
    { title: 'Modulable', body: 'De petites équipes à plus de trente personnes, selon votre événement.' },
  ],
  course: {
    name: 'Team building échecs à Genève',
    description:
      'Atelier de team building autour des échecs pour entreprises à Genève et dans l’arc lémanique, animé par un Maître FIDE.',
    url: '/team-building-echecs-geneve',
    courseMode: 'onsite',
  },
  related: [
    { to: '/conferences-echecs-entreprise', label: 'Conférence d’échecs en entreprise' },
    { to: '/contact', label: 'Demander un devis' },
  ],
  faq: [
    {
      question: 'Combien de personnes peuvent participer ?',
      answer:
        'De petits groupes jusqu’à une trentaine de personnes ou plus, avec des formats adaptés à l’effectif.',
    },
    {
      question: 'Faut-il savoir jouer aux échecs ?',
      answer:
        'Non, aucun niveau n’est requis. Les règles utiles sont expliquées et les activités sont pensées pour que chacun s’amuse.',
    },
    {
      question: 'Où se déroule l’atelier ?',
      answer:
        'Dans vos locaux à Genève, sur votre lieu de séminaire dans l’arc lémanique, ou tout autre lieu adapté.',
    },
    {
      question: 'Combien de temps dure un atelier ?',
      answer:
        'De 1 à 3 heures selon le format et l’effectif, avec la possibilité de l’intégrer à une demi-journée de séminaire.',
    },
    {
      question: 'Quel est l’intérêt des échecs pour la cohésion ?',
      answer:
        'Les formats coopératifs font travailler la communication, la répartition des rôles et la décision collective — des compétences directement transposables en équipe.',
    },
  ],
  children: <DeroulementFR />,
}

const EN: MoneyPageProps = {
  path: '/en/chess-team-building-geneva',
  title: 'Chess team building in Geneva',
  metaTitle: 'Chess team building in Geneva — corporate workshops',
  description:
    'Chess team-building workshops for companies in Geneva and the Lake Geneva region: cooperation, strategy and fun, led by a FIDE Master.',
  eyebrow: 'Team cohesion · Geneva',
  lead: (
    <>
      An <strong>original team-building activity</strong> built around chess: team tournaments,
      cooperative chess and challenges against a FIDE&nbsp;Master — to bring your teams
      together while having smart fun.
    </>
  ),
  facts: [
    { label: 'Format', value: 'Team workshop' },
    { label: 'Group', value: 'Up to 30+' },
    { label: 'Duration', value: '1 to 3 h' },
    { label: 'Location', value: 'Your offices / Geneva' },
  ],
  intro: (
    <>
      <p>
        <strong>A chess team-building session</strong> turns the game into a collective
        activity: team chess, cooperative variants, tactical quizzes and a friendly simul.
        No experience required — the goal is <strong>cooperation</strong>, communication and
        shared enjoyment.
      </p>
      <p>
        The workshop adapts to the size of your group and your objectives, at your offices in
        Geneva or at a seminar venue.
      </p>
    </>
  ),
  benefits: [
    { title: 'No prerequisites', body: 'Everyone takes part, whatever their level of chess.' },
    { title: 'Cooperation', body: 'Team formats that develop communication and collective decision-making.' },
    { title: 'Playful', body: 'Original variants, challenges and a simul for a friendly atmosphere.' },
    { title: 'Scalable', body: 'From small teams to more than thirty people, depending on your event.' },
  ],
  course: {
    name: 'Chess team building in Geneva',
    description:
      'A chess team-building workshop for companies in Geneva and the Lake Geneva region, led by a FIDE Master.',
    url: '/en/chess-team-building-geneva',
    courseMode: 'onsite',
  },
  related: [
    { to: '/en/corporate-chess-talks', label: 'Corporate chess talk' },
    { to: '/en/contact', label: 'Request a quote' },
  ],
  faq: [
    {
      question: 'How many people can take part?',
      answer:
        'From small groups to around thirty people or more, with formats tailored to the number of participants.',
    },
    {
      question: 'Do you need to know how to play chess?',
      answer:
        'No, no experience is required. The useful rules are explained and the activities are designed so everyone has fun.',
    },
    {
      question: 'Where does the workshop take place?',
      answer:
        'At your offices in Geneva, at your seminar venue in the Lake Geneva region, or any other suitable location.',
    },
    {
      question: 'How long does a workshop last?',
      answer:
        'From 1 to 3 hours depending on the format and group size, with the option to include it in a half-day seminar.',
    },
    {
      question: 'Why is chess good for team cohesion?',
      answer:
        'Cooperative formats develop communication, role sharing and collective decision-making — skills that transfer directly to the workplace.',
    },
  ],
  children: <DeroulementEN />,
}

export function Component() {
  const locale = useLocale()
  return <MoneyPage {...(locale === 'en' ? EN : FR)} />
}
