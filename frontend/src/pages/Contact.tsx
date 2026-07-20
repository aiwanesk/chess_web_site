import { Seo } from '../lib/seo'
import { SITE } from '../lib/site'
import { Container } from '../components/Container'
import { Section, Eyebrow } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { ContactForm } from '../components/ContactForm'
import { IconMail, IconPhone, IconPin } from '../components/icons'
import { localBusinessSchema, breadcrumbSchema, type Crumb } from '../lib/schema'

export function Component() {
  const crumbs: Crumb[] = [
    { name: 'Accueil', path: '/' },
    { name: 'Contact', path: '/contact' },
  ]

  return (
    <>
      <Seo
        title="Contact"
        description="Contactez Alexandre Iwanesko, Maître FIDE, pour un cours d’échecs à Genève ou en ligne. Premier échange pour définir vos objectifs."
        path="/contact"
        jsonLd={[breadcrumbSchema(crumbs), localBusinessSchema()]}
      />
      <Breadcrumbs crumbs={crumbs} />

      <Section>
        <Container className="grid gap-12 lg:grid-cols-[1.2fr_1fr]">
          <div>
            <Eyebrow>Contact</Eyebrow>
            <h1 className="font-display text-[2rem] font-bold tracking-tight text-ink-900 sm:text-[2.5rem]">
              Parlons de vos objectifs
            </h1>
            <p className="mt-5 max-w-xl text-lg leading-relaxed text-ink-600">
              Décrivez votre niveau et ce que vous souhaitez travailler. Je vous réponds
              rapidement pour définir le format le plus adapté — présentiel à Genève ou en ligne.
            </p>
            <div className="mt-9 rounded-3xl border border-ink-200/80 bg-white p-7 shadow-card sm:p-8">
              <ContactForm />
            </div>
          </div>

          <aside className="h-fit rounded-3xl border border-ink-200/80 bg-cream-100 p-7 shadow-soft">
            <h2 className="font-display text-xl font-bold text-ink-900">Coordonnées</h2>
            <address className="mt-6 space-y-5 not-italic text-ink-700">
              <p className="flex items-start gap-3">
                <span aria-hidden className="mt-0.5 flex h-9 w-9 flex-none items-center justify-center rounded-full bg-white text-gold-700 shadow-sm">
                  <IconMail size={18} />
                </span>
                <span>
                  <span className="block text-xs font-semibold uppercase tracking-wide text-ink-500">E-mail</span>
                  <a href={`mailto:${SITE.contact.email}`} className="font-medium text-gold-700 underline underline-offset-2">
                    {SITE.contact.email}
                  </a>
                </span>
              </p>
              <p className="flex items-start gap-3">
                <span aria-hidden className="mt-0.5 flex h-9 w-9 flex-none items-center justify-center rounded-full bg-white text-gold-700 shadow-sm">
                  <IconPhone size={18} />
                </span>
                <span>
                  <span className="block text-xs font-semibold uppercase tracking-wide text-ink-500">Téléphone</span>
                  <a href={SITE.contact.phoneHref} className="font-medium text-gold-700 underline underline-offset-2">
                    {SITE.contact.phone}
                  </a>
                </span>
              </p>
              <p className="flex items-start gap-3">
                <span aria-hidden className="mt-0.5 flex h-9 w-9 flex-none items-center justify-center rounded-full bg-white text-gold-700 shadow-sm">
                  <IconPin size={18} />
                </span>
                <span>
                  <span className="block text-xs font-semibold uppercase tracking-wide text-ink-500">Adresse</span>
                  {SITE.address.street}
                  <br />
                  {SITE.address.postalCode} {SITE.address.locality}
                </span>
              </p>
              <p className="flex items-start gap-3">
                <span aria-hidden className="mt-0.5 flex h-9 w-9 flex-none items-center justify-center rounded-full bg-white text-gold-700 shadow-sm">
                  <IconPin size={18} />
                </span>
                <span>
                  <span className="block text-xs font-semibold uppercase tracking-wide text-ink-500">Zone</span>
                  {SITE.address.locality} et {SITE.areaServed.join(', ')}.
                </span>
              </p>
            </address>

            {/*
              Carte : chargée à la demande pour ne pas pénaliser le LCP.
              Remplacer par une carte statique légère ou un embed cliquable
              plutôt qu'une iframe Google lourde au chargement initial.
            */}
            <div className="mt-7 flex aspect-video items-center justify-center rounded-2xl border border-dashed border-ink-300 bg-white text-center text-sm text-ink-500">
              Carte de Genève (chargement différé) — à intégrer
            </div>
          </aside>
        </Container>
      </Section>
    </>
  )
}
