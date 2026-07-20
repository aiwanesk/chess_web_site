import { MoneyPage } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'

export function Component() {
  return (
    <MoneyPage
      path="/cours-echecs-ados-competition"
      title="Cours d’échecs pour ados en compétition"
      metaTitle="Cours d’échecs ados compétition — Maître FIDE Genève"
      description="Coaching d’échecs pour adolescents en compétition à Genève : progression Elo, préparation tournoi et suivi individualisé par un Maître FIDE."
      cluster="tournoi"
      eyebrow="Jeunes joueurs · Compétition"
      lead={
        <>
          Pour les adolescents qui jouent déjà en compétition et veulent{' '}
          <strong>franchir un palier</strong> : un coaching exigeant et bienveillant, orienté{' '}
          <strong>progression Elo</strong> et résultats en tournoi.
        </>
      }
      facts={[
        { label: 'Public', value: 'Ados classés' },
        { label: 'Objectif', value: 'Progression Elo' },
        { label: 'Format', value: 'Individuel / duo' },
        { label: 'Suivi', value: 'Parents informés' },
      ]}
      intro={
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
      }
      benefits={[
        { title: 'Progression Elo', body: 'Un plan clair orienté résultats, adapté à l’âge et au niveau du joueur.' },
        { title: 'Mental de compétiteur', body: 'Gérer le stress, la pendule et l’après-défaite pour performer en tournoi.' },
        { title: 'Technique solide', body: 'Répertoire adapté, calcul et finales — les fondations d’une progression durable.' },
        { title: 'Suivi transparent', body: 'Objectifs partagés et retours réguliers aux parents.' },
      ]}
      course={{
        name: 'Coaching d’échecs pour ados en compétition',
        description:
          'Coaching individualisé pour adolescents joueurs de compétition : progression Elo, préparation tournoi et mental, par un Maître FIDE à Genève.',
        url: '/cours-echecs-ados-competition',
        price: 120,
        priceUnit: 'la séance individuelle de 60 min',
        courseMode: 'blended',
      }}
      related={[
        { to: '/preparation-tournoi-echecs', label: 'Préparer un tournoi précis' },
        { to: '/stages-echecs-geneve', label: 'Stages pendant les vacances' },
        { to: '/tarifs', label: 'Tarifs' },
      ]}
      faq={[
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
      ]}
    >
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
    </MoneyPage>
  )
}
