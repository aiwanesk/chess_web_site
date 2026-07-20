import { MoneyPage } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'

export function Component() {
  return (
    <MoneyPage
      path="/preparation-tournoi-echecs"
      title="Préparation à un tournoi d’échecs"
      metaTitle="Préparation tournoi d’échecs — coaching Maître FIDE"
      description="Préparation ciblée aux tournois d’échecs avec un Maître FIDE : répertoire, préparation adverse, gestion du temps et mental de compétition."
      cluster="tournoi"
      eyebrow="Compétition"
      lead={
        <>
          Un open approche&nbsp;? On prépare votre <strong>répertoire</strong>, votre{' '}
          <strong>gestion du temps</strong> et votre <strong>mental de compétition</strong> pour
          arriver le jour&nbsp;J avec un plan clair et de la confiance.
        </>
      }
      facts={[
        { label: 'Objectif', value: 'Performance en open' },
        { label: 'Format', value: 'Forfait dédié' },
        { label: 'Délai', value: 'Dès 3 semaines avant' },
        { label: 'Suivi', value: 'Analyse post-rondes' },
      ]}
      intro={
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
      }
      benefits={[
        { title: 'Répertoire verrouillé', body: 'Un répertoire que vous maîtrisez vraiment, sans trous théoriques exploitables.' },
        { title: 'Préparation adverse', body: 'Analyse des adversaires probables et lignes ciblées quand l’appariement est connu.' },
        { title: 'Gestion du temps', body: 'Routines pour éviter le zeitnot et garder de la lucidité en fin de partie.' },
        { title: 'Endurance & mental', body: 'Gérer la fatigue sur plusieurs rondes, encaisser une défaite et rebondir.' },
        { title: 'Analyse post-rondes', body: 'Débrief rapide entre les rondes pour corriger et préparer la suivante.' },
        { title: 'Plan jour J', body: 'Une routine d’avant-partie claire pour aborder chaque ronde sereinement.' },
      ]}
      course={{
        name: 'Préparation à un tournoi d’échecs',
        description:
          'Coaching de préparation aux tournois d’échecs avec un Maître FIDE : répertoire, préparation adverse, gestion du temps et suivi pendant l’événement.',
        url: '/preparation-tournoi-echecs',
        price: 120,
        priceUnit: 'la séance (forfaits dédiés)',
        courseMode: 'blended',
      }}
      related={[
        { to: '/cours-echecs-adultes-geneve', label: 'Progresser toute l’année (cours adultes)' },
        { to: '/cours-echecs-ados-competition', label: 'Ados en compétition' },
        { to: '/tarifs', label: 'Forfaits et tarifs' },
      ]}
      faq={[
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
      ]}
    >
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
    </MoneyPage>
  )
}
