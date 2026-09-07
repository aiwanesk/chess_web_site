import { MoneyPage, type MoneyPageProps } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'
import { useLocale } from '../lib/i18n'

function DeroulementFR() {
  return (
    <DeepDive eyebrow="Le déroulé" title="Comment se passe un cours en ligne">
      <StepList
        steps={[
          {
            t: 'Connexion',
            d: 'Un lien de visioconférence et un échiquier partagé Lichess : rien à installer, on démarre en deux clics.',
          },
          {
            t: 'Travail sur vos parties',
            d: 'On analyse vos parties récentes à l’écran, avec flèches et variantes, exactement comme sur un échiquier physique.',
          },
          {
            t: 'Exercices ciblés',
            d: 'Tactique, finales ou plans selon le thème du jour, avec un accès direct aux positions travaillées.',
          },
          {
            t: 'Rediffusion',
            d: 'La séance peut être enregistrée pour que vous la revoyiez et ancriez les idées entre deux cours.',
          },
        ]}
      />
    </DeepDive>
  )
}

function DeroulementEN() {
  return (
    <DeepDive eyebrow="The process" title="How an online lesson works">
      <StepList
        steps={[
          {
            t: 'Connecting',
            d: 'A video-call link and a shared Lichess board: nothing to install, we start in two clicks.',
          },
          {
            t: 'Working on your games',
            d: 'We analyse your recent games on screen, with arrows and variations, exactly as on a physical board.',
          },
          {
            t: 'Targeted exercises',
            d: 'Tactics, endgames or plans depending on the day’s theme, with direct access to the positions we work on.',
          },
          {
            t: 'Replay',
            d: 'The session can be recorded so you can watch it again and lock in the ideas between lessons.',
          },
        ]}
      />
    </DeepDive>
  )
}

const FR: MoneyPageProps = {
  path: '/cours-echecs-en-ligne',
  title: 'Cours d’échecs en ligne',
  metaTitle: 'Cours d’échecs en ligne — Maître FIDE (Suisse romande)',
  description:
    'Cours d’échecs particuliers en ligne avec un Maître FIDE, partout en Suisse romande et en France voisine. Même méthode qu’en présentiel.',
  eyebrow: 'À distance',
  lead: (
    <>
      Où que vous soyez en Suisse romande ou en France voisine, suivez des cours particuliers
      avec un <strong>Maître&nbsp;FIDE</strong> par visioconférence et échiquier partagé —{' '}
      <strong>même méthode, même suivi</strong> qu’en présentiel.
    </>
  ),
  facts: [
    { label: 'Format', value: 'Visio + échiquier' },
    { label: 'Zone', value: 'Suisse romande / FR' },
    { label: 'Séance', value: '60 minutes' },
    { label: 'Souplesse', value: 'Créneaux flexibles' },
  ],
  intro: (
    <>
      <p>
        <strong>Un cours d’échecs en ligne</strong> se déroule par visioconférence avec un
        échiquier interactif partagé : le coach voit vos analyses en temps réel, annote les
        positions et enregistre la séance si besoin. La qualité pédagogique est identique à
        celle d’un cours en présentiel.
      </p>
      <p>
        C’est la solution idéale si vous habitez Lausanne, Nyon, Fribourg, Annecy ou ailleurs,
        ou simplement pour gagner du temps de trajet.
      </p>
    </>
  ),
  benefits: [
    { title: 'Échiquier partagé', body: 'Analyse en direct des positions, avec annotations et flèches visibles à l’écran.' },
    { title: 'Séances enregistrées', body: 'Revoyez la séance à tête reposée pour ancrer les idées travaillées.' },
    { title: 'Créneaux flexibles', body: 'Des horaires adaptés à votre emploi du temps, en soirée ou le week-end.' },
    { title: 'Même méthode', body: 'Diagnostic, plan de progression et exercices, exactement comme en présentiel.' },
  ],
  course: {
    name: 'Cours d’échecs en ligne',
    description:
      'Cours particuliers d’échecs en ligne avec un Maître FIDE, par visioconférence et échiquier partagé, pour toute la Suisse romande.',
    url: '/cours-echecs-en-ligne',
    price: 120,
    priceUnit: 'la séance en ligne de 60 min',
    courseMode: 'online',
  },
  related: [
    { to: '/cours-echecs-adultes-geneve', label: 'Cours pour adultes (présentiel à Genève)' },
    { to: '/preparation-tournoi-echecs', label: 'Préparer un tournoi' },
    { to: '/tarifs', label: 'Tarifs' },
  ],
  faq: [
    {
      question: 'Quel matériel faut-il pour un cours en ligne ?',
      answer:
        'Un ordinateur avec webcam et une connexion correcte suffisent. On utilise un échiquier en ligne partagé (type Lichess) et une visioconférence.',
    },
    {
      question: 'La qualité est-elle vraiment équivalente au présentiel ?',
      answer:
        'Oui. L’échiquier partagé permet d’analyser ensemble en temps réel ; beaucoup d’élèves préfèrent même le format en ligne pour sa souplesse et les séances enregistrées.',
    },
    {
      question: 'Depuis quelles régions peut-on suivre les cours ?',
      answer:
        'Depuis toute la Suisse romande (Genève, Vaud, Fribourg, Valais, Neuchâtel) et la France voisine. Le fuseau horaire suisse est utilisé pour les créneaux.',
    },
    {
      question: 'Le cours en ligne convient-il à la préparation de tournoi ?',
      answer:
        'Tout à fait. Le format en ligne est même idéal pour préparer des adversaires et réviser un répertoire, grâce au partage d’écran et aux bases de données.',
    },
    {
      question: 'Peut-on mélanger présentiel et en ligne ?',
      answer:
        'Oui, beaucoup d’élèves alternent : présentiel à Genève quand c’est possible, en ligne le reste du temps. Le suivi reste continu d’un format à l’autre.',
    },
  ],
  children: <DeroulementFR />,
}

const EN: MoneyPageProps = {
  path: '/en/online-chess-lessons',
  title: 'Online chess lessons',
  metaTitle: 'Online chess lessons — FIDE Master (French-speaking Switzerland)',
  description:
    'Private online chess lessons with a FIDE Master, anywhere in French-speaking Switzerland and neighbouring France. The same method as in person.',
  eyebrow: 'Remote',
  lead: (
    <>
      Wherever you are in French-speaking Switzerland or neighbouring France, take private
      lessons with a <strong>FIDE&nbsp;Master</strong> over video call and a shared board —{' '}
      <strong>the same method, the same follow-up</strong> as in person.
    </>
  ),
  facts: [
    { label: 'Format', value: 'Video + board' },
    { label: 'Area', value: 'Romandy / France' },
    { label: 'Session', value: '60 minutes' },
    { label: 'Flexibility', value: 'Flexible slots' },
  ],
  intro: (
    <>
      <p>
        <strong>An online chess lesson</strong> takes place over video call with a shared
        interactive board: the coach sees your analysis in real time, annotates the positions
        and records the session if needed. The teaching quality is identical to an in-person
        lesson.
      </p>
      <p>
        It’s the ideal solution if you live in Lausanne, Nyon, Fribourg, Annecy or elsewhere, or
        simply to save yourself the travel time.
      </p>
    </>
  ),
  benefits: [
    { title: 'Shared board', body: 'Live analysis of positions, with annotations and arrows visible on screen.' },
    { title: 'Recorded sessions', body: 'Watch the session again at your own pace to lock in the ideas you worked on.' },
    { title: 'Flexible slots', body: 'Times that fit your schedule, in the evening or at the weekend.' },
    { title: 'The same method', body: 'Diagnosis, progression plan and exercises, exactly as in person.' },
  ],
  course: {
    name: 'Online chess lessons',
    description:
      'Private online chess lessons with a FIDE Master, over video call and a shared board, for the whole of French-speaking Switzerland.',
    url: '/en/online-chess-lessons',
    price: 120,
    priceUnit: 'per online 60-min session',
    courseMode: 'online',
  },
  related: [
    { to: '/en/adult-chess-lessons-geneva', label: 'Adult lessons (in person in Geneva)' },
    { to: '/en/tournament-preparation', label: 'Prepare for a tournament' },
    { to: '/en/pricing', label: 'Pricing' },
  ],
  faq: [
    {
      question: 'What equipment do I need for an online lesson?',
      answer:
        'A computer with a webcam and a decent connection are enough. We use a shared online board (such as Lichess) and a video call.',
    },
    {
      question: 'Is the quality really equivalent to in person?',
      answer:
        'Yes. The shared board lets us analyse together in real time; many students even prefer the online format for its flexibility and the recorded sessions.',
    },
    {
      question: 'Which regions can I take the lessons from?',
      answer:
        'From anywhere in French-speaking Switzerland (Geneva, Vaud, Fribourg, Valais, Neuchâtel) and neighbouring France. Swiss time is used for scheduling.',
    },
    {
      question: 'Is online coaching suitable for tournament preparation?',
      answer:
        'Absolutely. The online format is actually ideal for preparing against specific opponents and reviewing a repertoire, thanks to screen sharing and databases.',
    },
    {
      question: 'Can I mix in-person and online lessons?',
      answer:
        'Yes, many students alternate: in person in Geneva when possible, online the rest of the time. The follow-up stays continuous from one format to the other.',
    },
  ],
  children: <DeroulementEN />,
}

export function Component() {
  const locale = useLocale()
  return <MoneyPage {...(locale === 'en' ? EN : FR)} />
}
