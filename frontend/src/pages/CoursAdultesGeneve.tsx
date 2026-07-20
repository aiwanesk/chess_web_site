import { MoneyPage } from '../components/MoneyPage'
import { Container } from '../components/Container'
import { Section, Eyebrow } from '../components/ui'

/**
 * Reference money page — the 100 %-optimised template to duplicate:
 * single H1, answer-first intro, specific/citable facts, benefits, a
 * program breakdown, internal-link silo, FAQ, and full JSON-LD
 * (BreadcrumbList + Course/Offer + FAQPage) via <MoneyPage>.
 */
export function Component() {
  return (
    <MoneyPage
      path="/cours-echecs-adultes-geneve"
      title="Cours d’échecs pour adultes à Genève"
      metaTitle="Cours d’échecs adultes à Genève — Maître FIDE"
      description="Cours d’échecs pour adultes (1200–1900 Elo) à Genève avec un Maître FIDE. Méthode structurée, plan de progression, présentiel et en ligne."
      cluster="adultes"
      eyebrow="Cours pour adultes · Genève"
      lead={
        <>
          Vous jouez déjà et vous stagnez&nbsp;? Ces cours d’échecs pour adultes à Genève,
          encadrés par un <strong>Maître&nbsp;FIDE</strong>, s’adressent aux joueurs de{' '}
          <strong>1200 à 1900 Elo</strong> qui veulent une progression réelle et mesurable —
          pas une remise à niveau pour débutants.
        </>
      }
      facts={[
        { label: 'Public', value: '1200–1900 Elo' },
        { label: 'Format', value: 'Individuel / en ligne' },
        { label: 'Séance', value: '60 minutes' },
        { label: 'Langue', value: 'FR (EN possible)' },
      ]}
      intro={
        <>
          <p>
            <strong>Un cours d’échecs pour adultes</strong> est un accompagnement individualisé
            visant à corriger les faiblesses récurrentes d’un joueur et à structurer sa
            progression. À Genève, Alexandre&nbsp;Iwanesko, Maître&nbsp;FIDE, enseigne aux adultes
            de niveau intermédiaire à avancé en présentiel et en ligne.
          </p>
          <p>
            La démarche commence par un <strong>diagnostic</strong> de vos parties récentes :
            on identifie les schémas d’erreur (calcul, plans, gestion du temps, ouvertures
            mal comprises), puis on construit un <strong>plan de travail sur 8 à 12 semaines</strong>{' '}
            avec des objectifs Elo concrets et des exercices ciblés entre les séances.
          </p>
        </>
      }
      benefits={[
        {
          title: 'Diagnostic de vos parties',
          body: 'Analyse de vos parties classées pour cibler les 2–3 leviers qui débloquent le plus de points Elo.',
        },
        {
          title: 'Répertoire d’ouvertures adapté',
          body: 'Un répertoire cohérent avec votre style, que vous comprenez au lieu de le mémoriser.',
        },
        {
          title: 'Calcul & prise de décision',
          body: 'Méthode de calcul concret, gestion des candidats et du temps à la pendule.',
        },
        {
          title: 'Finales qui rapportent',
          body: 'Les finales essentielles pour convertir vos avantages et sauver les positions inférieures.',
        },
        {
          title: 'Travail entre les séances',
          body: 'Exercices et parties commentées à réaliser seul, corrigés à la séance suivante.',
        },
        {
          title: 'Suivi de progression',
          body: 'Objectifs Elo, points de contrôle réguliers et ajustement du plan selon vos résultats.',
        },
      ]}
      course={{
        name: 'Cours d’échecs pour adultes à Genève',
        description:
          'Cours particuliers d’échecs pour adultes 1200–1900 Elo à Genève, avec un Maître FIDE : diagnostic, plan de progression, ouvertures, calcul et finales.',
        url: '/cours-echecs-adultes-geneve',
        price: 120,
        priceUnit: 'la séance individuelle de 60 min',
        courseMode: 'blended',
      }}
      related={[
        { to: '/preparation-tournoi-echecs', label: 'Préparer un tournoi spécifique' },
        { to: '/cours-echecs-en-ligne', label: 'Suivre les cours en ligne' },
        { to: '/cours-echecs-groupe-geneve', label: 'Apprendre en petit groupe' },
        { to: '/tarifs', label: 'Consulter les tarifs et forfaits' },
      ]}
      faq={[
        {
          question: 'À partir de quel niveau ces cours sont-ils adaptés ?',
          answer:
            'Ils s’adressent aux joueurs adultes d’environ 1200 à 1900 Elo. Si vous connaissez les règles et jouez déjà en ligne ou en club, vous êtes au bon endroit. Les grands débutants ne sont pas le public visé.',
        },
        {
          question: 'Combien coûte un cours d’échecs pour adultes à Genève ?',
          answer:
            'La séance individuelle de 60 minutes est à 120 CHF, et un pack de 10 séances est à 1000 CHF (soit 100 CHF la séance). Le détail figure sur la page Tarifs.',
        },
        {
          question: 'Les cours ont-ils lieu en présentiel ou en ligne ?',
          answer:
            'Les deux. En présentiel à Genève, ou en ligne par visioconférence avec un échiquier partagé — même méthode, même qualité de suivi, où que vous soyez en Suisse romande.',
        },
        {
          question: 'En combien de temps peut-on gagner des points Elo ?',
          answer:
            'La plupart des élèves assidus constatent une progression sur 2 à 3 mois de travail régulier. Le rythme dépend de votre point de départ et du temps consacré aux exercices entre les séances.',
        },
        {
          question: 'Qu’est-ce qu’un Maître FIDE ?',
          answer:
            'Un Maître FIDE (FIDE Master) est un titre international décerné par la Fédération internationale des échecs, attribué aux joueurs ayant atteint un niveau Elo élevé et stable. C’est une garantie d’expertise pour l’enseignement.',
        },
      ]}
    >
      {/* Program breakdown — extra citable structure. */}
      <Section>
        <Container>
          <Eyebrow>Le déroulé</Eyebrow>
          <h2 className="text-2xl font-bold text-ink-900 sm:text-3xl">
            Comment se passe un accompagnement type
          </h2>
          <ol className="mt-8 space-y-6">
            {[
              {
                t: 'Séance 1 — Diagnostic',
                d: 'Revue de vos parties et de vos résultats, définition d’objectifs Elo et du plan sur 8–12 semaines.',
              },
              {
                t: 'Séances 2 à N — Travail ciblé',
                d: 'Ouvertures, calcul, milieu de jeu et finales selon vos faiblesses prioritaires, avec exercices entre les séances.',
              },
              {
                t: 'Points de contrôle',
                d: 'Toutes les 4 séances, on mesure les progrès et on réajuste le plan en fonction de vos parties récentes.',
              },
            ].map((s, i) => (
              <li key={s.t} className="flex gap-4">
                <span className="flex h-9 w-9 flex-none items-center justify-center rounded-full bg-gold-500 font-bold text-ink-950">
                  {i + 1}
                </span>
                <div>
                  <h3 className="text-lg font-semibold text-ink-900">{s.t}</h3>
                  <p className="mt-1 text-ink-600">{s.d}</p>
                </div>
              </li>
            ))}
          </ol>
        </Container>
      </Section>
    </MoneyPage>
  )
}
