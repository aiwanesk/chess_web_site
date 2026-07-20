import { MoneyPage } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'

export function Component() {
  return (
    <MoneyPage
      path="/cours-echecs-groupe-geneve"
      title="Cours d’échecs en groupe à Genève"
      metaTitle="Cours d’échecs en groupe à Genève — Maître FIDE"
      description="Cours d’échecs en petit groupe à Genève : niveau homogène, émulation et tarif réduit, encadrés par un Maître FIDE."
      eyebrow="Petit groupe · Genève"
      lead={
        <>
          Apprenez à plusieurs, dans un <strong>petit groupe de niveau homogène</strong> (3 à 6
          joueurs) : l’émulation d’un club, la rigueur d’un coaching, et un{' '}
          <strong>tarif par personne réduit</strong>.
        </>
      }
      facts={[
        { label: 'Taille', value: '3 à 6 joueurs' },
        { label: 'Niveau', value: 'Groupes homogènes' },
        { label: 'Séance', value: '90 minutes' },
        { label: 'Lieu', value: 'Genève' },
      ]}
      intro={
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
      }
      benefits={[
        { title: 'Émulation', body: 'Progresser en confrontant ses idées à des joueurs de niveau proche.' },
        { title: 'Tarif partagé', body: 'Un coût par personne réduit par rapport au cours individuel.' },
        { title: 'Thèmes structurés', body: 'Un programme progressif commun, avec exercices collectifs et corrigés.' },
        { title: 'Mini-tournois', body: 'Des parties commentées en direct pour appliquer immédiatement les notions.' },
      ]}
      course={{
        name: 'Cours d’échecs en groupe à Genève',
        description:
          'Cours d’échecs en petit groupe (3 à 6 joueurs) de niveau homogène à Genève, encadrés par un Maître FIDE.',
        url: '/cours-echecs-groupe-geneve',
        price: 60,
        priceUnit: 'par personne et par séance de 90 min',
        courseMode: 'onsite',
      }}
      related={[
        { to: '/cours-echecs-adultes-geneve', label: 'Préférer un cours individuel' },
        { to: '/stages-echecs-geneve', label: 'Stages intensifs' },
        { to: '/tarifs', label: 'Tarifs' },
      ]}
      faq={[
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
      ]}
    >
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
    </MoneyPage>
  )
}
