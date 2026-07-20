import { MoneyPage } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'
import { eventSchema } from '../lib/schema'

export function Component() {
  return (
    <MoneyPage
      path="/stages-echecs-geneve"
      title="Stages d’échecs à Genève"
      metaTitle="Stages d’échecs à Genève — vacances, Maître FIDE"
      description="Stages d’échecs intensifs à Genève pendant les vacances scolaires, encadrés par un Maître FIDE. Progression rapide sur quelques jours."
      eyebrow="Vacances scolaires · Genève"
      lead={
        <>
          Des <strong>stages intensifs</strong> pendant les vacances pour progresser vite :
          plusieurs jours de travail structuré, d’analyse et de parties, encadrés par un{' '}
          <strong>Maître&nbsp;FIDE</strong>.
        </>
      }
      facts={[
        { label: 'Durée', value: '2 à 5 jours' },
        { label: 'Rythme', value: 'Intensif' },
        { label: 'Niveau', value: 'Groupes par niveau' },
        { label: 'Lieu', value: 'Genève' },
      ]}
      intro={
        <>
          <p>
            <strong>Un stage d’échecs</strong> est une session intensive de quelques jours qui
            concentre un maximum de travail utile : thèmes stratégiques, exercices de calcul,
            finales et parties commentées. Le format vacances permet de progresser nettement en
            peu de temps.
          </p>
          <p>
            Les stages sont organisés par niveau pour garantir un rythme adapté à chaque
            participant.
          </p>
        </>
      }
      benefits={[
        { title: 'Immersion', body: 'Plusieurs jours dédiés aux échecs pour ancrer durablement les progrès.' },
        { title: 'Programme dense', body: 'Stratégie, calcul, finales et parties commentées en alternance.' },
        { title: 'Par niveau', body: 'Des groupes homogènes pour un rythme adapté à chacun.' },
        { title: 'Bilan personnalisé', body: 'Un retour individuel avec des axes de travail pour la suite.' },
      ]}
      course={{
        name: 'Stage d’échecs à Genève',
        description:
          'Stage d’échecs intensif à Genève pendant les vacances scolaires, encadré par un Maître FIDE, organisé par niveau.',
        url: '/stages-echecs-geneve',
        price: 240,
        priceUnit: 'le stage de plusieurs jours',
        courseMode: 'onsite',
      }}
      extraJsonLd={[
        eventSchema({
          name: 'Stage d’échecs intensif — Genève',
          description:
            'Stage d’échecs de plusieurs jours à Genève pendant les vacances scolaires, encadré par un Maître FIDE.',
          url: '/stages-echecs-geneve',
          // Placeholder dates — mettre à jour à chaque session programmée.
          startDate: '2026-10-19',
          endDate: '2026-10-21',
          price: 240,
        }),
      ]}
      related={[
        { to: '/cours-echecs-ados-competition', label: 'Coaching ados en compétition' },
        { to: '/cours-echecs-groupe-geneve', label: 'Cours en groupe à l’année' },
        { to: '/tarifs', label: 'Tarifs' },
      ]}
      faq={[
        {
          question: 'Quand ont lieu les stages ?',
          answer:
            'Pendant les vacances scolaires genevoises. Les dates des prochaines sessions sont annoncées sur cette page et par e-mail ; contactez-moi pour être informé.',
        },
        {
          question: 'Les stages sont-ils adaptés aux adultes comme aux jeunes ?',
          answer:
            'Oui, des groupes distincts sont formés par niveau et par public afin que chacun travaille au bon rythme.',
        },
        {
          question: 'Quel niveau faut-il avoir ?',
          answer:
            'Les stages s’adressent à des joueurs qui connaissent déjà les bases et jouent régulièrement. Précisez votre niveau à l’inscription.',
        },
        {
          question: 'Le repas et le matériel sont-ils inclus ?',
          answer:
            'Le matériel d’échecs est fourni. Les modalités pratiques (horaires, repas, pauses) sont précisées dans le programme de chaque session.',
        },
      ]}
    >
      <DeepDive eyebrow="Le programme" title="Une journée type de stage">
        <StepList
          steps={[
            {
              t: 'Matin — Thème stratégique',
              d: 'Cours interactif sur une notion clé (structures de pions, attaque du roi, finales), illustré par des parties de maîtres.',
            },
            {
              t: 'Midi — Résolution d’exercices',
              d: 'Ateliers de calcul et de tactique par niveau, en individuel et en équipe.',
            },
            {
              t: 'Après-midi — Parties & analyse',
              d: 'Parties d’entraînement avec cadence, puis analyse collective à chaud des positions marquantes.',
            },
            {
              t: 'Fin de stage — Bilan',
              d: 'Un retour individuel avec des axes de travail personnalisés pour la suite.',
            },
          ]}
        />
      </DeepDive>
    </MoneyPage>
  )
}
