import { MoneyPage } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'

export function Component() {
  return (
    <MoneyPage
      path="/team-building-echecs-geneve"
      title="Team building échecs à Genève"
      metaTitle="Team building échecs à Genève — ateliers d’équipe"
      description="Ateliers de team building autour des échecs pour entreprises à Genève et dans l’arc lémanique : coopération, stratégie et fun, par un Maître FIDE."
      eyebrow="Cohésion d’équipe · Genève"
      lead={
        <>
          Un <strong>team building original</strong> autour des échecs : tournois par équipes,
          échecs coopératifs et défis contre un Maître&nbsp;FIDE — pour souder vos équipes en
          s’amusant intelligemment.
        </>
      }
      facts={[
        { label: 'Format', value: 'Atelier d’équipe' },
        { label: 'Groupe', value: 'Jusqu’à 30+' },
        { label: 'Durée', value: '1 à 3 h' },
        { label: 'Lieu', value: 'Vos locaux / Genève' },
      ]}
      intro={
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
      }
      benefits={[
        { title: 'Sans prérequis', body: 'Tout le monde participe, quel que soit le niveau d’échecs.' },
        { title: 'Coopération', body: 'Des formats en équipe qui font travailler la communication et la décision collective.' },
        { title: 'Ludique', body: 'Variantes originales, défis et simultanée pour une ambiance conviviale.' },
        { title: 'Modulable', body: 'De petites équipes à plus de trente personnes, selon votre événement.' },
      ]}
      course={{
        name: 'Team building échecs à Genève',
        description:
          'Atelier de team building autour des échecs pour entreprises à Genève et dans l’arc lémanique, animé par un Maître FIDE.',
        url: '/team-building-echecs-geneve',
        courseMode: 'onsite',
      }}
      related={[
        { to: '/conferences-echecs-entreprise', label: 'Conférence d’échecs en entreprise' },
        { to: '/contact', label: 'Demander un devis' },
      ]}
      faq={[
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
      ]}
    >
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
    </MoneyPage>
  )
}
