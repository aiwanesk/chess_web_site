import { MoneyPage, type MoneyPageProps } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'
import { useLocale } from '../lib/i18n'

const FR: MoneyPageProps = {
  path: '/preparation-tournoi-echecs',
  title: 'Préparation à un tournoi d’échecs',
  metaTitle: 'Préparation tournoi d’échecs — coaching Maître FIDE',
  description:
    'Préparation ciblée aux tournois d’échecs avec un Maître FIDE : répertoire, préparation adverse, gestion du temps et mental de compétition.',
  cluster: 'tournoi',
  eyebrow: 'Compétition',
  lead: (
    <>
      Un open approche&nbsp;? On prépare votre <strong>répertoire</strong>, votre{' '}
      <strong>gestion du temps</strong> et votre <strong>mental de compétition</strong> pour
      arriver le jour&nbsp;J avec un plan clair et de la confiance.
    </>
  ),
  facts: [
    { label: 'Objectif', value: 'Performance en open' },
    { label: 'Format', value: 'Forfait dédié' },
    { label: 'Délai', value: 'Dès 3 semaines avant' },
    { label: 'Suivi', value: 'Analyse post-rondes' },
  ],
  intro: (
    <>
      <p>
        <strong>La préparation à un tournoi d’échecs</strong> consiste à optimiser, sur une
        période courte, les facteurs qui décident d’une performance : répertoire d’ouvertures
        fiable, préparation contre des adversaires connus, endurance sur plusieurs rondes et
        gestion de la pendule.
      </p>
      <p>
        On part de vos parties et du format du tournoi visé pour bâtir un plan réaliste, puis
        on peut assurer un <strong>suivi pendant l’événement</strong> avec analyse entre les
        rondes.
      </p>
    </>
  ),
  benefits: [
    { title: 'Répertoire verrouillé', body: 'Un répertoire que vous maîtrisez vraiment, sans trous théoriques exploitables.' },
    { title: 'Préparation adverse', body: 'Analyse des adversaires probables et lignes ciblées quand l’appariement est connu.' },
    { title: 'Gestion du temps', body: 'Routines pour éviter le zeitnot et garder de la lucidité en fin de partie.' },
    { title: 'Endurance & mental', body: 'Gérer la fatigue sur plusieurs rondes, encaisser une défaite et rebondir.' },
    { title: 'Analyse post-rondes', body: 'Débrief rapide entre les rondes pour corriger et préparer la suivante.' },
    { title: 'Plan jour J', body: 'Une routine d’avant-partie claire pour aborder chaque ronde sereinement.' },
  ],
  course: {
    name: 'Préparation à un tournoi d’échecs',
    description:
      'Coaching de préparation aux tournois d’échecs avec un Maître FIDE : répertoire, préparation adverse, gestion du temps et suivi pendant l’événement.',
    url: '/preparation-tournoi-echecs',
    price: 120,
    priceUnit: 'la séance (forfaits dédiés)',
    courseMode: 'blended',
  },
  related: [
    { to: '/cours-echecs-adultes-geneve', label: 'Progresser toute l’année (cours adultes)' },
    { to: '/cours-echecs-ados-competition', label: 'Ados en compétition' },
    { to: '/tarifs', label: 'Forfaits et tarifs' },
  ],
  faq: [
    {
      question: 'Combien de temps avant le tournoi faut-il commencer ?',
      answer:
        'Idéalement 3 à 6 semaines avant, pour consolider le répertoire et installer les routines. Une préparation express reste utile même à quelques jours de l’événement.',
    },
    {
      question: 'Assurez-vous un suivi pendant le tournoi ?',
      answer:
        'Oui. Selon le forfait, un débrief entre les rondes permet d’ajuster la préparation et d’aborder chaque adversaire avec un plan.',
    },
    {
      question: 'Faut-il un niveau minimum ?',
      answer:
        'La préparation s’adresse aux joueurs classés (environ 1200 Elo et plus) qui disputent des opens ou des championnats.',
    },
    {
      question: 'Que travaille-t-on concrètement dans le répertoire ?',
      answer:
        'On sécurise vos lignes principales avec les deux couleurs, on comble les trous théoriques que vos adversaires pourraient exploiter, et on choisit des systèmes que vous comprenez plutôt que des variantes à mémoriser.',
    },
    {
      question: 'La préparation fonctionne-t-elle pour le rapide et le blitz ?',
      answer:
        'Oui. On adapte le contenu à la cadence : en rapide et en blitz, l’accent porte davantage sur des schémas fiables, l’intuition et les pièges d’ouverture que sur la théorie profonde.',
    },
  ],
  children: (
    <DeepDive eyebrow="Le déroulé" title="Comment se déroule une préparation">
      <StepList
        steps={[
          {
            t: 'Bilan express',
            d: 'Analyse de vos dernières parties de compétition et du format visé (cadence, nombre de rondes, niveau du plateau).',
          },
          {
            t: 'Répertoire & plans',
            d: 'On verrouille les ouvertures et on prépare des plans-types pour vos structures, avec un aide-mémoire à réviser avant chaque ronde.',
          },
          {
            t: 'Routines de compétition',
            d: 'Gestion du temps, routine d’avant-partie, alimentation et récupération entre les rondes pour tenir la distance.',
          },
          {
            t: 'Suivi pendant l’événement',
            d: 'Selon le forfait, débrief entre les rondes pour corriger et, quand l’appariement est connu, préparer l’adversaire suivant.',
          },
        ]}
      />
    </DeepDive>
  ),
}

const EN: MoneyPageProps = {
  path: '/en/tournament-preparation',
  title: 'Chess tournament preparation',
  metaTitle: 'Chess tournament preparation — FIDE Master coaching',
  description:
    'Focused chess tournament preparation with a FIDE Master: repertoire, opponent prep, time management and competitive mindset.',
  eyebrow: 'Competition',
  lead: (
    <>
      An open on the horizon? We sharpen your <strong>repertoire</strong>, your{' '}
      <strong>time management</strong> and your <strong>competitive mindset</strong> so you
      arrive on the day with a clear plan and real confidence.
    </>
  ),
  facts: [
    { label: 'Goal', value: 'Open performance' },
    { label: 'Format', value: 'Dedicated package' },
    { label: 'Lead time', value: 'From 3 weeks out' },
    { label: 'Follow-up', value: 'Post-round analysis' },
  ],
  intro: (
    <>
      <p>
        <strong>Chess tournament preparation</strong> is about fine-tuning, over a short
        window, the factors that decide a result: a reliable opening repertoire, targeted prep
        against known opponents, stamina across multiple rounds and clock management.
      </p>
      <p>
        We start from your own games and the format you are aiming for to build a realistic
        plan, and can then provide <strong>support during the event</strong> with analysis
        between rounds.
      </p>
    </>
  ),
  benefits: [
    { title: 'A watertight repertoire', body: 'A repertoire you genuinely own, with no theoretical gaps for opponents to exploit.' },
    { title: 'Opponent preparation', body: 'Analysis of likely opponents and targeted lines once the pairing is known.' },
    { title: 'Time management', body: 'Routines to avoid time trouble and stay clear-headed deep into the game.' },
    { title: 'Stamina & mindset', body: 'Managing fatigue across rounds, taking a loss on the chin and bouncing back.' },
    { title: 'Post-round analysis', body: 'A quick debrief between rounds to fix mistakes and prepare the next one.' },
    { title: 'Match-day plan', body: 'A clear pre-game routine so you approach every round with a calm head.' },
  ],
  course: {
    name: 'Chess tournament preparation',
    description:
      'Tournament preparation coaching with a FIDE Master: repertoire, opponent prep, time management and support during the event.',
    url: '/en/tournament-preparation',
    price: 120,
    priceUnit: 'per session (dedicated packages)',
    courseMode: 'blended',
  },
  related: [
    { to: '/en/adult-chess-lessons-geneva', label: 'Progress all year round (adult lessons)' },
    { to: '/en/junior-chess-coaching', label: 'Juniors in competition' },
    { to: '/en/pricing', label: 'Packages & pricing' },
  ],
  faq: [
    {
      question: 'How far ahead of the tournament should we start?',
      answer:
        'Ideally 3 to 6 weeks out, to consolidate the repertoire and lock in the routines. Even an express prep in the final few days can make a real difference.',
    },
    {
      question: 'Do you provide support during the tournament?',
      answer:
        'Yes. Depending on the package, a debrief between rounds lets us adjust the preparation and approach each opponent with a plan.',
    },
    {
      question: 'Is there a minimum level?',
      answer:
        'Preparation is for rated players (roughly 1200 Elo and up) competing in opens or championships.',
    },
    {
      question: 'What exactly do we work on in the repertoire?',
      answer:
        'We shore up your main lines with both colours, close the theoretical gaps your opponents could exploit, and choose systems you understand rather than variations to memorise.',
    },
    {
      question: 'Does this work for rapid and blitz?',
      answer:
        'Yes. We tailor the content to the time control: in rapid and blitz the focus shifts toward reliable patterns, intuition and opening traps rather than deep theory.',
    },
  ],
  children: (
    <DeepDive eyebrow="The process" title="How a preparation unfolds">
      <StepList
        steps={[
          {
            t: 'Quick assessment',
            d: 'A review of your recent competitive games and the target format (time control, number of rounds, strength of the field).',
          },
          {
            t: 'Repertoire & plans',
            d: 'We lock in the openings and prepare model plans for your structures, with a cheat sheet to revise before each round.',
          },
          {
            t: 'Competition routines',
            d: 'Time management, a pre-game routine, and nutrition and recovery between rounds to go the distance.',
          },
          {
            t: 'Support during the event',
            d: 'Depending on the package, a debrief between rounds to fix mistakes and, once the pairing is known, prepare the next opponent.',
          },
        ]}
      />
    </DeepDive>
  ),
}

export function Component() {
  const locale = useLocale()
  return <MoneyPage {...(locale === 'en' ? EN : FR)} />
}
