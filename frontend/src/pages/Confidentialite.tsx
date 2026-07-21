import { Seo } from '../lib/seo'
import { SITE } from '../lib/site'
import { Container } from '../components/Container'
import { Section } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { breadcrumbSchema, type Crumb } from '../lib/schema'
import { useLocale, homePath, pathFor, t, type Locale } from '../lib/i18n'

type Block = { h: string; p: string[] }

const COPY: Record<Locale, {
  title: string; description: string; eyebrow: string; lead: string; updated: string; blocks: Block[]
}> = {
  fr: {
    title: 'Politique de confidentialité',
    description:
      'Comment les données personnelles sont traitées sur iwanesko.ch (formulaire de contact et newsletter). Conforme au RGPD et à la LPD suisse.',
    eyebrow: 'Confidentialité',
    lead: 'Cette page explique quelles données sont collectées sur iwanesko.ch, pourquoi, et comment exercer vos droits. Traitement conforme au RGPD (UE) et à la LPD suisse (nLPD).',
    updated: 'Dernière mise à jour : juillet 2026.',
    blocks: [
      {
        h: 'Responsable du traitement',
        p: [
          `Alexandre Iwanesko — ${SITE.address.street}, ${SITE.address.postalCode} ${SITE.address.locality}, Suisse.`,
          `Contact : ${SITE.contact.email}.`,
        ],
      },
      {
        h: 'Newsletter — données et finalité',
        p: [
          'Si vous vous inscrivez à la newsletter, nous conservons uniquement votre adresse e-mail et votre langue préférée. Aucune autre donnée (nom, adresse IP, traceur) n’est collectée à cette fin.',
          'Finalité : vous informer des nouveautés du site (nouveaux articles, tactiques hebdomadaires, dates de stages).',
          'Base légale : votre consentement explicite, recueilli via une case à cocher non pré-cochée, puis confirmé par un lien envoyé à votre adresse (double opt-in). Tant que vous n’avez pas cliqué ce lien, aucun autre e-mail ne vous est envoyé.',
        ],
      },
      {
        h: 'Formulaire de contact',
        p: [
          'Les informations que vous saisissez dans le formulaire de contact (nom, e-mail, niveau, message) servent uniquement à traiter votre demande. Elles ne sont ni revendues, ni utilisées à des fins publicitaires.',
        ],
      },
      {
        h: 'Durée de conservation',
        p: [
          'Newsletter : jusqu’à votre désinscription. Lorsque vous vous désabonnez, votre adresse est supprimée définitivement de la liste.',
          'Demandes de contact : conservées le temps nécessaire au suivi de votre demande.',
        ],
      },
      {
        h: 'Vos droits',
        p: [
          'Vous disposez d’un droit d’accès, de rectification, d’effacement et d’opposition, ainsi que du droit de retirer votre consentement à tout moment.',
          `Pour la newsletter, le retrait est immédiat : chaque e-mail contient un lien de désinscription en un clic. Vous pouvez aussi écrire à ${SITE.contact.email}.`,
        ],
      },
      {
        h: 'Hébergement et sous-traitants',
        p: [
          'Le site et la base d’abonnés sont hébergés sur une infrastructure en Europe. L’envoi des e-mails passe par un prestateur de messagerie (Infomaniak, Suisse). Aucune donnée n’est transmise à des fins commerciales à des tiers.',
        ],
      },
    ],
  },
  en: {
    title: 'Privacy policy',
    description:
      'How personal data is handled on iwanesko.ch (contact form and newsletter). GDPR and Swiss FADP compliant.',
    eyebrow: 'Privacy',
    lead: 'This page explains what data is collected on iwanesko.ch, why, and how to exercise your rights. Processing complies with the GDPR (EU) and the Swiss FADP (nFADP).',
    updated: 'Last updated: July 2026.',
    blocks: [
      {
        h: 'Data controller',
        p: [
          `Alexandre Iwanesko — ${SITE.address.street}, ${SITE.address.postalCode} ${SITE.address.locality}, Switzerland.`,
          `Contact: ${SITE.contact.email}.`,
        ],
      },
      {
        h: 'Newsletter — data and purpose',
        p: [
          'If you subscribe to the newsletter, we store only your e-mail address and preferred language. No other data (name, IP address, tracker) is collected for this purpose.',
          'Purpose: to let you know about new content (new articles, weekly tactics, upcoming events).',
          'Legal basis: your explicit consent, collected via an unchecked checkbox and then confirmed through a link sent to your address (double opt-in). Until you click that link, no other e-mail is sent to you.',
        ],
      },
      {
        h: 'Contact form',
        p: [
          'The information you enter in the contact form (name, e-mail, level, message) is used only to handle your request. It is never sold or used for advertising.',
        ],
      },
      {
        h: 'Retention',
        p: [
          'Newsletter: until you unsubscribe. When you unsubscribe, your address is permanently deleted from the list.',
          'Contact requests: kept for as long as needed to follow up on your request.',
        ],
      },
      {
        h: 'Your rights',
        p: [
          'You have the right to access, rectify, erase and object, as well as the right to withdraw your consent at any time.',
          `For the newsletter, withdrawal is immediate: every e-mail contains a one-click unsubscribe link. You can also write to ${SITE.contact.email}.`,
        ],
      },
      {
        h: 'Hosting and processors',
        p: [
          'The site and the subscriber database are hosted on infrastructure in Europe. E-mails are sent through a mail provider (Infomaniak, Switzerland). No data is shared with third parties for commercial purposes.',
        ],
      },
    ],
  },
}

export function Component() {
  const locale = useLocale()
  const c = COPY[locale]
  const path = pathFor('confidentialite', locale)
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: c.title, path },
  ]

  return (
    <>
      <Seo title={c.title} description={c.description} path={path} jsonLd={[breadcrumbSchema(crumbs)]} />
      <Breadcrumbs crumbs={crumbs} />
      <PageHero eyebrow={c.eyebrow} title={c.title} lead={c.lead} />

      <Section>
        <Container>
          <div className="mx-auto max-w-2xl space-y-8">
            {c.blocks.map((b) => (
              <section key={b.h}>
                <h2 className="font-display text-xl font-bold text-ink-900">{b.h}</h2>
                {b.p.map((para, i) => (
                  <p key={i} className="mt-3 leading-relaxed text-ink-600">
                    {para}
                  </p>
                ))}
              </section>
            ))}
            <p className="border-t border-ink-200 pt-6 text-sm text-ink-500">{c.updated}</p>
          </div>
        </Container>
      </Section>
    </>
  )
}
