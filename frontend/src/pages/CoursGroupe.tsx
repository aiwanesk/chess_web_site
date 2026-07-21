import { MoneyPage, type MoneyPageProps } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'
import { useLocale } from '../lib/i18n'

/**
 * Money page, bilingual. The same component serves the FR route
 * (/cours-echecs-groupe-geneve) and the EN route (/en/group-chess-lessons-geneva);
 * it reads the locale and renders the matching content.
 */

function DeroulementFR() {
  return (
    <DeepDive eyebrow="Le déroulé" title="Comment se déroule une séance de groupe">
      <StepList
        steps={[
          {
            t: 'Thème du jour',
            d: 'Une notion stratégique ou tactique commune, illustrée par des positions marquantes et adaptée au niveau du groupe.',
          },
          {
            t: 'Exercices collectifs',
            d: 'On cherche ensemble : chacun propose ses idées, on confronte les plans, on corrige en direct.',
          },
          {
            t: 'Mini-tournoi commenté',
            d: 'Des parties entre participants, analysées à chaud pour appliquer immédiatement la notion travaillée.',
          },
          {
            t: 'À retenir',
            d: 'Une synthèse et un exercice à faire d’ici la prochaine séance pour ancrer les acquis.',
          },
        ]}
      />
    </DeepDive>
  )
}

function DeroulementEN() {
  return (
    <DeepDive eyebrow="The process" title="How a group session unfolds">
      <StepList
        steps={[
          {
            t: 'Theme of the day',
            d: 'A shared strategic or tactical concept, illustrated with striking positions and adapted to the group’s level.',
          },
          {
            t: 'Group exercises',
            d: 'We work it out together: everyone shares their ideas, we compare plans and correct them live.',
          },
          {
            t: 'Commented mini-tournament',
            d: 'Games between participants, analysed on the spot to apply the concept just covered right away.',
          },
          {
            t: 'Key takeaways',
            d: 'A summary and an exercise to do before the next session to lock in what you’ve learned.',
          },
        ]}
      />
    </DeepDive>
  )
}

const FR: MoneyPageProps = {
  path: '/cours-echecs-groupe-geneve',
  title: 'Cours d’échecs en groupe à Genève',
  metaTitle: 'Cours d’échecs en groupe à Genève — Maître FIDE',
  description:
    'Cours d’échecs en petit groupe à Genève : niveau homogène, émulation et tarif réduit, encadrés par un Maître FIDE.',
  eyebrow: 'Petit groupe · Genève',
  lead: (
    <>
      Apprenez à plusieurs, dans un <strong>petit groupe de niveau homogène</strong> (3 à 6
      joueurs) : l’émulation d’un club, la rigueur d’un coaching, et un{' '}
      <strong>tarif par personne réduit</strong>.
    </>
  ),
  facts: [
    { label: 'Taille', value: '3 à 6 joueurs' },
    { label: 'Niveau', value: 'Groupes homogènes' },
    { label: 'Séance', value: '60 minutes' },
    { label: 'Lieu', value: 'Genève' },
  ],
  intro: (
    <>
      <p>
        <strong>Un cours d’échecs en groupe</strong> réunit quelques joueurs de niveau proche
        autour d’un même programme. On y travaille des thèmes communs (ouvertures, plans,
        finales) avec des exercices collectifs, des mini-tournois commentés et des analyses
        partagées.
      </p>
      <p>
        Le format est plus économique que le cours particulier et développe une saine
        émulation, tout en gardant un suivi de qualité grâce aux effectifs réduits.
      </p>
    </>
  ),
  benefits: [
    { title: 'Émulation', body: 'Progresser en confrontant ses idées à des joueurs de niveau proche.' },
    { title: 'Tarif partagé', body: 'Un coût par personne réduit par rapport au cours individuel.' },
    { title: 'Thèmes structurés', body: 'Un programme progressif commun, avec exercices collectifs et corrigés.' },
    { title: 'Mini-tournois', body: 'Des parties commentées en direct pour appliquer immédiatement les notions.' },
  ],
  course: {
    name: 'Cours d’échecs en groupe à Genève',
    description:
      'Cours d’échecs en petit groupe (3 à 6 joueurs) de niveau homogène à Genève, encadrés par un Maître FIDE.',
    url: '/cours-echecs-groupe-geneve',
    price: 60,
    priceUnit: 'par personne et par séance de 60 min',
    courseMode: 'onsite',
  },
  related: [
    { to: '/cours-echecs-adultes-geneve', label: 'Préférer un cours individuel' },
    { to: '/stages-echecs-geneve', label: 'Stages intensifs' },
    { to: '/tarifs', label: 'Tarifs' },
  ],
  faq: [
    {
      question: 'Comment les groupes sont-ils constitués ?',
      answer:
        'Par niveau, pour que chacun progresse au bon rythme. Un court échange initial permet de vous placer dans le groupe adapté à votre Elo.',
    },
    {
      question: 'Combien de participants par groupe ?',
      answer:
        'De 3 à 6 joueurs, afin de garder un suivi individualisé tout en profitant de l’émulation collective.',
    },
    {
      question: 'Peut-on rejoindre un groupe en cours d’année ?',
      answer:
        'Oui, selon les places disponibles et l’adéquation de niveau. Contactez-moi pour connaître les groupes ouverts.',
    },
    {
      question: 'Le cours en groupe remplace-t-il le cours particulier ?',
      answer:
        'Ce sont deux formats complémentaires. Le groupe est idéal pour l’émulation et le coût ; le cours particulier permet un suivi 100 % personnalisé. Certains élèves combinent les deux.',
    },
    {
      question: 'Faut-il venir avec son matériel ?',
      answer:
        'Non, le matériel est fourni sur place. Il suffit d’apporter de quoi prendre des notes si vous le souhaitez.',
    },
  ],
  children: <DeroulementFR />,
}

const EN: MoneyPageProps = {
  path: '/en/group-chess-lessons-geneva',
  title: 'Group chess lessons in Geneva',
  metaTitle: 'Group chess lessons in Geneva — FIDE Master',
  description:
    'Small-group chess lessons in Geneva: matched level, healthy competition and a reduced rate, taught by a FIDE Master.',
  eyebrow: 'Small group · Geneva',
  lead: (
    <>
      Learn together, in a <strong>small group of matched level</strong> (3 to 6 players): the
      buzz of a club, the rigour of coaching, and a{' '}
      <strong>reduced per-person rate</strong>.
    </>
  ),
  facts: [
    { label: 'Size', value: '3 to 6 players' },
    { label: 'Level', value: 'Matched groups' },
    { label: 'Session', value: '60 minutes' },
    { label: 'Location', value: 'Geneva' },
  ],
  intro: (
    <>
      <p>
        <strong>A group chess lesson</strong> brings together a few players of similar level
        around the same programme. We work on common themes (openings, plans, endgames) with
        group exercises, commented mini-tournaments and shared analysis.
      </p>
      <p>
        The format is more affordable than private lessons and fosters healthy competition,
        while keeping quality follow-up thanks to the small numbers.
      </p>
    </>
  ),
  benefits: [
    { title: 'Healthy competition', body: 'Improve by testing your ideas against players of a similar level.' },
    { title: 'Shared cost', body: 'A lower per-person cost than one-to-one lessons.' },
    { title: 'Structured themes', body: 'A shared, progressive programme with group exercises and corrections.' },
    { title: 'Mini-tournaments', body: 'Games commented live to apply the concepts straight away.' },
  ],
  course: {
    name: 'Group chess lessons in Geneva',
    description:
      'Small-group chess lessons (3 to 6 players) of matched level in Geneva, taught by a FIDE Master.',
    url: '/en/group-chess-lessons-geneva',
    price: 60,
    priceUnit: 'per person per 60-min session',
    courseMode: 'onsite',
  },
  related: [
    { to: '/en/adult-chess-lessons-geneva', label: 'Prefer one-to-one lessons' },
    { to: '/en/chess-camps-geneva', label: 'Intensive camps' },
    { to: '/en/pricing', label: 'Pricing' },
  ],
  faq: [
    {
      question: 'How are the groups formed?',
      answer:
        'By level, so everyone progresses at the right pace. A short initial chat lets me place you in the group that matches your Elo.',
    },
    {
      question: 'How many participants per group?',
      answer:
        'From 3 to 6 players, to keep individual follow-up while enjoying the benefits of group learning.',
    },
    {
      question: 'Can you join a group during the year?',
      answer:
        'Yes, depending on availability and level fit. Get in touch to find out which groups are open.',
    },
    {
      question: 'Does the group lesson replace private lessons?',
      answer:
        'They are two complementary formats. The group is ideal for competition and cost; private lessons give fully personalised follow-up. Some students combine both.',
    },
    {
      question: 'Do I need to bring my own equipment?',
      answer:
        'No, equipment is provided on site. Just bring something to take notes with if you like.',
    },
  ],
  children: <DeroulementEN />,
}

export function Component() {
  const locale = useLocale()
  return <MoneyPage {...(locale === 'en' ? EN : FR)} />
}
