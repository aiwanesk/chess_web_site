import { MoneyPage } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'

export function Component() {
  return (
    <MoneyPage
      path="/cours-echecs-en-ligne"
      title="Cours d’échecs en ligne"
      metaTitle="Cours d’échecs en ligne — Maître FIDE (Suisse romande)"
      description="Cours d’échecs particuliers en ligne avec un Maître FIDE, partout en Suisse romande et en France voisine. Même méthode qu’en présentiel."
      eyebrow="À distance"
      lead={
        <>
          Où que vous soyez en Suisse romande ou en France voisine, suivez des cours particuliers
          avec un <strong>Maître&nbsp;FIDE</strong> par visioconférence et échiquier partagé —{' '}
          <strong>même méthode, même suivi</strong> qu’en présentiel.
        </>
      }
      facts={[
        { label: 'Format', value: 'Visio + échiquier' },
        { label: 'Zone', value: 'Suisse romande / FR' },
        { label: 'Séance', value: '60 minutes' },
        { label: 'Souplesse', value: 'Créneaux flexibles' },
      ]}
      intro={
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
      }
      benefits={[
        { title: 'Échiquier partagé', body: 'Analyse en direct des positions, avec annotations et flèches visibles à l’écran.' },
        { title: 'Séances enregistrées', body: 'Revoyez la séance à tête reposée pour ancrer les idées travaillées.' },
        { title: 'Créneaux flexibles', body: 'Des horaires adaptés à votre emploi du temps, en soirée ou le week-end.' },
        { title: 'Même méthode', body: 'Diagnostic, plan de progression et exercices, exactement comme en présentiel.' },
      ]}
      course={{
        name: 'Cours d’échecs en ligne',
        description:
          'Cours particuliers d’échecs en ligne avec un Maître FIDE, par visioconférence et échiquier partagé, pour toute la Suisse romande.',
        url: '/cours-echecs-en-ligne',
        price: 120,
        priceUnit: 'la séance en ligne de 60 min',
        courseMode: 'online',
      }}
      related={[
        { to: '/cours-echecs-adultes-geneve', label: 'Cours pour adultes (présentiel à Genève)' },
        { to: '/preparation-tournoi-echecs', label: 'Préparer un tournoi' },
        { to: '/tarifs', label: 'Tarifs' },
      ]}
      faq={[
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
      ]}
    >
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
    </MoneyPage>
  )
}
