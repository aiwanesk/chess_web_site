import { MoneyPage } from '../components/MoneyPage'
import { DeepDive, StepList } from '../components/ui'

export function Component() {
  return (
    <MoneyPage
      path="/conferences-echecs-entreprise"
      title="Conférences d’échecs en entreprise"
      metaTitle="Conférences échecs en entreprise — Genève & lémanique"
      description="Conférences et interventions échecs en entreprise à Genève : stratégie, prise de décision et gestion du risque, par un Maître FIDE."
      eyebrow="Entreprises"
      lead={
        <>
          Une intervention inspirante qui relie les <strong>échecs</strong> aux enjeux de
          l’entreprise : <strong>stratégie</strong>, <strong>prise de décision sous pression</strong>{' '}
          et gestion du risque, animée par un Maître&nbsp;FIDE.
        </>
      }
      facts={[
        { label: 'Format', value: 'Conférence / atelier' },
        { label: 'Durée', value: '45–90 min' },
        { label: 'Public', value: 'Équipes, comités' },
        { label: 'Lieu', value: 'Vos locaux / Genève' },
      ]}
      intro={
        <>
          <p>
            <strong>Une conférence d’échecs en entreprise</strong> utilise le jeu comme métaphore
            concrète des décisions stratégiques : anticipation, évaluation des risques, gestion du
            temps et lucidité sous pression. Le propos est vivant, illustré de positions célèbres
            et adapté à votre secteur.
          </p>
          <p>
            L’intervention peut prendre la forme d’une conférence, d’un atelier interactif ou d’une
            simultanée où le conférencier affronte plusieurs collaborateurs à la fois.
          </p>
        </>
      }
      benefits={[
        { title: 'Message sur mesure', body: 'Un fil conducteur aligné sur vos enjeux : décision, risque, stratégie.' },
        { title: 'Format flexible', body: 'Conférence, atelier participatif ou simultanée, selon l’effet recherché.' },
        { title: 'Interactif', body: 'Des exemples concrets et des moments de participation qui marquent les esprits.' },
        { title: 'Clé en main', body: 'Organisation simple, dans vos locaux à Genève ou dans l’arc lémanique.' },
      ]}
      course={{
        name: 'Conférence d’échecs en entreprise',
        description:
          'Conférence ou atelier échecs en entreprise sur la stratégie et la prise de décision, animé par un Maître FIDE à Genève.',
        url: '/conferences-echecs-entreprise',
        courseMode: 'onsite',
      }}
      related={[
        { to: '/team-building-echecs-geneve', label: 'Team building échecs' },
        { to: '/contact', label: 'Demander un devis' },
      ]}
      faq={[
        {
          question: 'Combien de temps dure une conférence ?',
          answer:
            'De 45 à 90 minutes selon le format choisi, avec possibilité d’ajouter un atelier ou une simultanée.',
        },
        {
          question: 'Faut-il savoir jouer aux échecs pour en profiter ?',
          answer:
            'Non. Le propos est accessible à tous ; les notions échiquéennes sont expliquées et servent de support aux idées de stratégie et de décision.',
        },
        {
          question: 'Intervenez-vous en dehors de Genève ?',
          answer:
            'Oui, dans tout l’arc lémanique et au-delà sur demande. Contactez-moi pour un devis adapté.',
        },
        {
          question: 'Quels thèmes pouvez-vous relier aux échecs ?',
          answer:
            'Prise de décision sous incertitude, anticipation, gestion du risque et du temps, apprentissage de l’erreur, calme sous pression. Le fil est adapté à votre secteur et à votre message.',
        },
        {
          question: 'Combien de participants pour une simultanée ?',
          answer:
            'Une simultanée peut opposer le conférencier à une dizaine, voire une vingtaine de collaborateurs en parallèle — un format spectaculaire et fédérateur.',
        },
      ]}
    >
      <DeepDive eyebrow="Formats" title="Trois manières d’intervenir">
        <StepList
          steps={[
            {
              t: 'Conférence',
              d: 'Un propos structuré (45–60 min) reliant les échecs à vos enjeux de décision et de stratégie, illustré d’exemples célèbres.',
            },
            {
              t: 'Atelier interactif',
              d: 'Les participants jouent et expérimentent les notions : anticipation, plans, gestion du risque, en équipes.',
            },
            {
              t: 'Simultanée',
              d: 'Le conférencier affronte plusieurs collaborateurs à la fois — un moment marquant qui clôt idéalement un séminaire.',
            },
          ]}
        />
      </DeepDive>
    </MoneyPage>
  )
}
