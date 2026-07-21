import { MoneyPage, type MoneyPageProps } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'
import { useLocale } from '../lib/i18n'

/**
 * Bilingual money page. The same component serves the FR route
 * (/cours-echecs-ados-competition) and the EN route (/en/junior-chess-coaching);
 * it reads the locale and renders the matching content.
 */

function MethodeFR() {
  return (
    <DeepDive eyebrow="La méthode" title="Comment on construit la progression">
      <StepList
        steps={[
          {
            t: 'Objectifs partagés',
            d: 'Un point avec le jeune et les parents pour fixer des objectifs Elo réalistes et un calendrier de tournois.',
          },
          {
            t: 'Technique adaptée',
            d: 'Répertoire d’ouvertures adapté à l’âge, entraînement au calcul et finales essentielles.',
          },
          {
            t: 'Compétition',
            d: 'Préparation des tournois, gestion du temps et du stress, débrief des parties jouées.',
          },
          {
            t: 'Bilan régulier',
            d: 'Suivi du classement et ajustement du plan, avec un retour transparent aux parents.',
          },
        ]}
      />
    </DeepDive>
  )
}

function MethodeEN() {
  return (
    <DeepDive eyebrow="The method" title="How we build progress">
      <StepList
        steps={[
          {
            t: 'Shared goals',
            d: 'A meeting with the teenager and parents to set realistic Elo goals and a tournament calendar.',
          },
          {
            t: 'Tailored technique',
            d: 'An age-appropriate opening repertoire, calculation training and essential endgames.',
          },
          {
            t: 'Competition',
            d: 'Tournament preparation, time and stress management, and debriefing the games played.',
          },
          {
            t: 'Regular review',
            d: 'Rating tracking and plan adjustments, with transparent feedback to parents.',
          },
        ]}
      />
    </DeepDive>
  )
}

const FR: MoneyPageProps = {
  path: '/cours-echecs-ados-competition',
  title: 'Cours d’échecs pour ados en compétition',
  metaTitle: 'Cours d’échecs ados compétition — Maître FIDE Genève',
  description:
    'Coaching d’échecs pour adolescents en compétition à Genève : progression Elo, préparation tournoi et suivi individualisé par un Maître FIDE.',
  cluster: 'tournoi',
  eyebrow: 'Jeunes joueurs · Compétition',
  lead: (
    <>
      Pour les adolescents qui jouent déjà en compétition et veulent{' '}
      <strong>franchir un palier</strong> : un coaching exigeant et bienveillant, orienté{' '}
      <strong>progression Elo</strong> et résultats en tournoi.
    </>
  ),
  facts: [
    { label: 'Public', value: 'Ados classés' },
    { label: 'Objectif', value: 'Progression Elo' },
    { label: 'Format', value: 'Individuel / duo' },
    { label: 'Suivi', value: 'Parents informés' },
  ],
  intro: (
    <>
      <p>
        <strong>Le coaching d’échecs pour adolescents en compétition</strong> vise les jeunes
        joueurs déjà classés qui disputent des tournois et cherchent à progresser vite. Le
        travail combine technique (ouvertures, calcul, finales) et compétences de compétiteur
        (gestion du temps, du stress et des émotions).
      </p>
      <p>
        Un point régulier avec les parents permet de suivre les objectifs et l’évolution du
        classement, sans mettre de pression inutile sur le jeune.
      </p>
    </>
  ),
  benefits: [
    { title: 'Progression Elo', body: 'Un plan clair orienté résultats, adapté à l’âge et au niveau du joueur.' },
    { title: 'Mental de compétiteur', body: 'Gérer le stress, la pendule et l’après-défaite pour performer en tournoi.' },
    { title: 'Technique solide', body: 'Répertoire adapté, calcul et finales — les fondations d’une progression durable.' },
    { title: 'Suivi transparent', body: 'Objectifs partagés et retours réguliers aux parents.' },
  ],
  course: {
    name: 'Coaching d’échecs pour ados en compétition',
    description:
      'Coaching individualisé pour adolescents joueurs de compétition : progression Elo, préparation tournoi et mental, par un Maître FIDE à Genève.',
    url: '/cours-echecs-ados-competition',
    price: 120,
    priceUnit: 'la séance individuelle de 60 min',
    courseMode: 'blended',
  },
  related: [
    { to: '/preparation-tournoi-echecs', label: 'Préparer un tournoi précis' },
    { to: '/stages-echecs-geneve', label: 'Stages pendant les vacances' },
    { to: '/tarifs', label: 'Tarifs' },
  ],
  faq: [
    {
      question: 'À quel âge et quel niveau s’adresse ce coaching ?',
      answer:
        'Aux adolescents (environ 11–18 ans) déjà classés et disputant des compétitions. Pour une initiation, ce format n’est pas adapté.',
    },
    {
      question: 'Les parents sont-ils tenus informés ?',
      answer:
        'Oui. Les objectifs sont définis ensemble et un point régulier fait le bilan des progrès et du classement.',
    },
    {
      question: 'Le coaching peut-il se faire en ligne ?',
      answer:
        'Oui, en présentiel à Genève ou en ligne, selon vos contraintes. La méthode et le suivi restent identiques.',
    },
    {
      question: 'À quelle fréquence faut-il des séances ?',
      answer:
        'Un rythme hebdomadaire ou bimensuel donne les meilleurs résultats, complété par du travail personnel guidé entre les séances. On adapte selon les tournois et le calendrier scolaire.',
    },
    {
      question: 'Comment gérez-vous la pression de la compétition ?',
      answer:
        'On travaille explicitement le mental : routines d’avant-partie, gestion de la défaite, concentration. L’objectif est que le jeune prenne du plaisir tout en performant.',
    },
  ],
  children: <MethodeFR />,
}

const EN: MoneyPageProps = {
  path: '/en/junior-chess-coaching',
  title: 'Chess coaching for competitive teenagers',
  metaTitle: 'Junior chess coaching — FIDE Master, Geneva',
  description:
    'Chess coaching for competitive teenagers in Geneva: Elo progress, tournament preparation and individual follow-up by a FIDE Master.',
  eyebrow: 'Young players · Competition',
  lead: (
    <>
      For teenagers who already play competitively and want to{' '}
      <strong>reach the next level</strong>: demanding yet supportive coaching, focused on{' '}
      <strong>Elo progress</strong> and tournament results.
    </>
  ),
  facts: [
    { label: 'Audience', value: 'Rated teens' },
    { label: 'Goal', value: 'Elo progress' },
    { label: 'Format', value: 'One-to-one / pairs' },
    { label: 'Follow-up', value: 'Parents kept informed' },
  ],
  intro: (
    <>
      <p>
        <strong>Chess coaching for competitive teenagers</strong> is aimed at young players who
        are already rated, play in tournaments and want to progress quickly. The work combines
        technique (openings, calculation, endgames) with competitor skills (managing time,
        stress and emotions).
      </p>
      <p>
        Regular check-ins with parents make it possible to track goals and rating progress,
        without putting unnecessary pressure on the young player.
      </p>
    </>
  ),
  benefits: [
    { title: 'Elo progress', body: 'A clear, results-oriented plan tailored to the player’s age and level.' },
    { title: 'A competitor’s mindset', body: 'Managing stress, the clock and the aftermath of a loss to perform in tournaments.' },
    { title: 'Solid technique', body: 'A suitable repertoire, calculation and endgames — the foundations of lasting progress.' },
    { title: 'Transparent follow-up', body: 'Shared goals and regular feedback to parents.' },
  ],
  course: {
    name: 'Chess coaching for competitive teenagers',
    description:
      'Individual coaching for teenage competition players: Elo progress, tournament preparation and mental game, by a FIDE Master in Geneva.',
    url: '/en/junior-chess-coaching',
    price: 120,
    priceUnit: 'per one-to-one 60-min session',
    courseMode: 'blended',
  },
  related: [
    { to: '/en/tournament-preparation', label: 'Prepare for a specific tournament' },
    { to: '/en/chess-camps-geneva', label: 'Holiday chess camps' },
    { to: '/en/pricing', label: 'Pricing' },
  ],
  faq: [
    {
      question: 'What age and level is this coaching for?',
      answer:
        'For teenagers (roughly 11–18) who are already rated and playing in competitions. It is not suited to a first introduction to the game.',
    },
    {
      question: 'Are parents kept informed?',
      answer:
        'Yes. Goals are set together and regular check-ins review progress and rating.',
    },
    {
      question: 'Can the coaching take place online?',
      answer:
        'Yes, in person in Geneva or online, depending on your constraints. The method and follow-up remain the same.',
    },
    {
      question: 'How often should sessions take place?',
      answer:
        'A weekly or fortnightly rhythm gives the best results, supplemented by guided personal work between sessions. We adapt to tournaments and the school calendar.',
    },
    {
      question: 'How do you handle the pressure of competition?',
      answer:
        'We work on the mental game explicitly: pre-game routines, coping with defeat, concentration. The goal is for the young player to enjoy the game while performing.',
    },
  ],
  children: <MethodeEN />,
}

export function Component() {
  const locale = useLocale()
  return <MoneyPage {...(locale === 'en' ? EN : FR)} />
}
