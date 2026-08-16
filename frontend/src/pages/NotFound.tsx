import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section, CtaLink } from '../components/ui'
import { useLocale, homePath, pathFor } from '../lib/i18n'

export function Component() {
  // The copy is French-only, but the links must still follow the locale: a
  // visitor who mistypes an /en/ URL should not be sent back into the FR tree.
  const locale = useLocale()
  return (
    <>
      <Seo title="Page introuvable (404)" description="La page demandée n’existe pas ou a été déplacée." path="/404" noindex />
      <Section>
        <Container className="text-center">
          <p className="font-display text-7xl font-extrabold text-ink-200 sm:text-8xl">404</p>
          <h1 className="mt-4 font-display text-2xl font-bold text-ink-900 sm:text-3xl">Cette page n’existe pas</h1>
          <p className="mx-auto mt-3 max-w-md leading-relaxed text-ink-600">
            Le lien est peut-être erroné ou la page a été déplacée. Reprenons depuis le début.
          </p>
          <div className="mt-8 flex flex-wrap justify-center gap-3">
            <CtaLink to={homePath(locale)} variant="primary">
              Retour à l’accueil
            </CtaLink>
            <CtaLink to={pathFor('coursAdultes', locale)} variant="secondary">
              Voir les cours
            </CtaLink>
          </div>
        </Container>
      </Section>
    </>
  )
}
