import { MoneyPage, type MoneyPageProps } from '../components/MoneyPage'
import { Container } from '../components/Container'
import { Section, Eyebrow } from '../components/ui'
import { useLocale } from '../lib/i18n'

/**
 * Flagship money page, bilingual. The same component serves the FR route
 * (/cours-echecs-adultes-geneve) and the EN route (/en/adult-chess-lessons-geneva);
 * it reads the locale and renders the matching content. Reference to duplicate.
 */

function DeroulementFR() {
  return (
    <Section>
      <Container>
        <Eyebrow>Le déroulé</Eyebrow>
        <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">
          Comment se passe un accompagnement type
        </h2>
        <ol className="mt-8 space-y-6">
          {[
            { t: 'Séance 1 — Diagnostic', d: 'Revue de vos parties et de vos résultats, définition d’objectifs Elo et du plan sur 8–12 semaines.' },
            { t: 'Séances 2 à N — Travail ciblé', d: 'Ouvertures, calcul, milieu de jeu et finales selon vos faiblesses prioritaires, avec exercices entre les séances.' },
            { t: 'Points de contrôle', d: 'Toutes les 4 séances, on mesure les progrès et on réajuste le plan en fonction de vos parties récentes.' },
          ].map((s, i) => (
            <li key={s.t} className="flex gap-4">
              <span className="flex h-9 w-9 flex-none items-center justify-center rounded-full bg-gold-500 font-bold text-ink-950">{i + 1}</span>
              <div>
                <h3 className="text-lg font-semibold text-ink-900">{s.t}</h3>
                <p className="mt-1 text-ink-600">{s.d}</p>
              </div>
            </li>
          ))}
        </ol>
      </Container>
    </Section>
  )
}

function DeroulementEN() {
  return (
    <Section>
      <Container>
        <Eyebrow>The process</Eyebrow>
        <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">How a typical programme works</h2>
        <ol className="mt-8 space-y-6">
          {[
            { t: 'Session 1 — Diagnosis', d: 'Review of your games and results, Elo goals and an 8–12 week plan.' },
            { t: 'Sessions 2 to N — Targeted work', d: 'Openings, calculation, middlegame and endgames based on your priority weaknesses, with exercises between sessions.' },
            { t: 'Checkpoints', d: 'Every 4 sessions we measure progress and adjust the plan based on your recent games.' },
          ].map((s, i) => (
            <li key={s.t} className="flex gap-4">
              <span className="flex h-9 w-9 flex-none items-center justify-center rounded-full bg-gold-500 font-bold text-ink-950">{i + 1}</span>
              <div>
                <h3 className="text-lg font-semibold text-ink-900">{s.t}</h3>
                <p className="mt-1 text-ink-600">{s.d}</p>
              </div>
            </li>
          ))}
        </ol>
      </Container>
    </Section>
  )
}

const FR: MoneyPageProps = {
  path: '/cours-echecs-adultes-geneve',
  title: 'Cours d’échecs pour adultes à Genève',
  metaTitle: 'Cours d’échecs adultes à Genève — Maître FIDE',
  description:
    'Cours d’échecs pour adultes (1200–2200 Elo) à Genève avec un Maître FIDE. Méthode structurée, plan de progression, présentiel et en ligne.',
  cluster: 'adultes',
  eyebrow: 'Cours pour adultes · Genève',
  lead: (
    <>
      Vous jouez déjà et vous stagnez&nbsp;? Ces cours d’échecs pour adultes à Genève, encadrés par un{' '}
      <strong>Maître&nbsp;FIDE</strong>, s’adressent aux joueurs de <strong>1200 à 2200 Elo</strong> qui veulent une
      progression réelle et mesurable — pas une remise à niveau pour débutants.
    </>
  ),
  facts: [
    { label: 'Public', value: '1200–2200 Elo' },
    { label: 'Format', value: 'Individuel / en ligne' },
    { label: 'Séance', value: '60 minutes' },
    { label: 'Langue', value: 'FR (EN possible)' },
  ],
  intro: (
    <>
      <p>
        <strong>Un cours d’échecs pour adultes</strong> est un accompagnement individualisé visant à corriger les
        faiblesses récurrentes d’un joueur et à structurer sa progression. À Genève, Alexandre&nbsp;Iwanesko,
        Maître&nbsp;FIDE, enseigne aux adultes de niveau intermédiaire à avancé en présentiel et en ligne.
      </p>
      <p>
        La démarche commence par un <strong>diagnostic</strong> de vos parties récentes : on identifie les schémas
        d’erreur (calcul, plans, gestion du temps, ouvertures mal comprises), puis on construit un{' '}
        <strong>plan de travail sur 8 à 12 semaines</strong> avec des objectifs Elo concrets.
      </p>
    </>
  ),
  benefits: [
    { title: 'Diagnostic de vos parties', body: 'Analyse de vos parties classées pour cibler les 2–3 leviers qui débloquent le plus de points Elo.' },
    { title: 'Répertoire d’ouvertures adapté', body: 'Un répertoire cohérent avec votre style, que vous comprenez au lieu de le mémoriser.' },
    { title: 'Calcul & prise de décision', body: 'Méthode de calcul concret, gestion des candidats et du temps à la pendule.' },
    { title: 'Finales qui rapportent', body: 'Les finales essentielles pour convertir vos avantages et sauver les positions inférieures.' },
    { title: 'Travail entre les séances', body: 'Exercices et parties commentées à réaliser seul, corrigés à la séance suivante.' },
    { title: 'Suivi de progression', body: 'Objectifs Elo, points de contrôle réguliers et ajustement du plan selon vos résultats.' },
  ],
  course: {
    name: 'Cours d’échecs pour adultes à Genève',
    description:
      'Cours particuliers d’échecs pour adultes 1200–2200 Elo à Genève, avec un Maître FIDE : diagnostic, plan de progression, ouvertures, calcul et finales.',
    url: '/cours-echecs-adultes-geneve',
    price: 120,
    priceUnit: 'la séance individuelle de 60 min',
    courseMode: 'blended',
  },
  related: [
    { to: '/preparation-tournoi-echecs', label: 'Préparer un tournoi spécifique' },
    { to: '/cours-echecs-en-ligne', label: 'Suivre les cours en ligne' },
    { to: '/cours-echecs-groupe-geneve', label: 'Apprendre en petit groupe' },
    { to: '/tarifs', label: 'Consulter les tarifs et forfaits' },
  ],
  faq: [
    { question: 'À partir de quel niveau ces cours sont-ils adaptés ?', answer: 'Ils s’adressent aux joueurs adultes d’environ 1200 à 2200 Elo. Si vous connaissez les règles et jouez déjà en ligne ou en club, vous êtes au bon endroit. Les grands débutants ne sont pas le public visé.' },
    { question: 'Combien coûte un cours d’échecs pour adultes à Genève ?', answer: 'La séance individuelle de 60 minutes est à 120 CHF, et un pack de 10 séances est à 1000 CHF (soit 100 CHF la séance). Le détail figure sur la page Tarifs.' },
    { question: 'Les cours ont-ils lieu en présentiel ou en ligne ?', answer: 'Les deux. En présentiel à Genève, ou en ligne par visioconférence avec un échiquier partagé — même méthode, même qualité de suivi.' },
    { question: 'En combien de temps peut-on gagner des points Elo ?', answer: 'La plupart des élèves assidus constatent une progression sur 2 à 3 mois de travail régulier. Le rythme dépend de votre point de départ et du temps consacré aux exercices.' },
    { question: 'Qu’est-ce qu’un Maître FIDE ?', answer: 'Un Maître FIDE (FIDE Master) est un titre international décerné par la Fédération internationale des échecs, attribué aux joueurs ayant atteint un niveau Elo élevé et stable.' },
  ],
  children: <DeroulementFR />,
}

const EN: MoneyPageProps = {
  path: '/en/adult-chess-lessons-geneva',
  title: 'Adult chess lessons in Geneva',
  metaTitle: 'Adult chess lessons in Geneva — FIDE Master',
  description:
    'Chess lessons for adults (1200–2200 Elo) in Geneva with a FIDE Master. Structured method, progression plan, in person and online.',
  eyebrow: 'Adult lessons · Geneva',
  lead: (
    <>
      Already playing but stuck at a plateau? These adult chess lessons in Geneva, taught by a{' '}
      <strong>FIDE&nbsp;Master</strong>, are for players from <strong>1200 to 2200 Elo</strong> who want real,
      measurable progress — not a beginner refresher.
    </>
  ),
  facts: [
    { label: 'Level', value: '1200–2200 Elo' },
    { label: 'Format', value: 'One-to-one / online' },
    { label: 'Session', value: '60 minutes' },
    { label: 'Language', value: 'EN / FR' },
  ],
  intro: (
    <>
      <p>
        <strong>An adult chess lesson</strong> is a one-to-one programme designed to fix a player’s recurring
        weaknesses and structure their progress. In Geneva, Alexandre&nbsp;Iwanesko, a FIDE&nbsp;Master, coaches
        intermediate-to-advanced adults in person and online.
      </p>
      <p>
        It starts with a <strong>diagnosis</strong> of your recent games: we identify your error patterns
        (calculation, plans, time management, poorly understood openings), then build an{' '}
        <strong>8–12 week plan</strong> with concrete Elo goals.
      </p>
    </>
  ),
  benefits: [
    { title: 'Diagnosis of your games', body: 'Analysis of your rated games to target the 2–3 levers that unlock the most Elo.' },
    { title: 'A repertoire that fits you', body: 'A coherent opening repertoire suited to your style — one you understand rather than memorise.' },
    { title: 'Calculation & decisions', body: 'A concrete calculation method, candidate moves and clock management.' },
    { title: 'Endgames that pay off', body: 'The essential endgames to convert your advantages and save worse positions.' },
    { title: 'Work between sessions', body: 'Exercises and annotated games to do on your own, reviewed at the next session.' },
    { title: 'Progress tracking', body: 'Elo goals, regular checkpoints and plan adjustments based on your results.' },
  ],
  course: {
    name: 'Adult chess lessons in Geneva',
    description:
      'Private chess lessons for adults 1200–2200 Elo in Geneva with a FIDE Master: diagnosis, progression plan, openings, calculation and endgames.',
    url: '/en/adult-chess-lessons-geneva',
    price: 120,
    priceUnit: 'per one-to-one 60-min session',
    courseMode: 'blended',
  },
  related: [{ to: '/en/pricing', label: 'See pricing & packages' }],
  faq: [
    { question: 'What level are these lessons for?', answer: 'For adult players roughly 1200 to 2200 Elo. If you know the rules and already play online or in a club, you’re in the right place. Complete beginners are not the target audience.' },
    { question: 'How much does an adult chess lesson in Geneva cost?', answer: 'A one-to-one 60-minute session is 120 CHF, and a 10-session package is 1000 CHF (i.e. 100 CHF per session). Full details are on the Pricing page.' },
    { question: 'Are lessons in person or online?', answer: 'Both. In person in Geneva, or online via video call with a shared board — same method, same quality of follow-up.' },
    { question: 'How fast can I gain Elo?', answer: 'Most committed students see progress within 2 to 3 months of regular work. The pace depends on your starting point and the time spent on exercises.' },
    { question: 'What is a FIDE Master?', answer: 'A FIDE Master is an international title awarded by the International Chess Federation to players who have reached a high, stable Elo level.' },
  ],
  children: <DeroulementEN />,
}

export function Component() {
  const locale = useLocale()
  return <MoneyPage {...(locale === 'en' ? EN : FR)} />
}
