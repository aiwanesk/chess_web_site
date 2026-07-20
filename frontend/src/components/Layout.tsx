import { Outlet, ScrollRestoration, useLocation } from 'react-router-dom'
import { Header } from './Header'
import { Footer } from './Footer'
import { SiteBanner } from './SiteBanner'
import { LocaleProvider, localeFromPath, t } from '../lib/i18n'

/** Root layout wrapping every route: detects the locale from the URL, provides
 * it to the tree, then renders skip link, header, main, footer. */
export function Layout() {
  const locale = localeFromPath(useLocation().pathname)
  const strings = t(locale)
  return (
    <LocaleProvider locale={locale}>
      <a
        href="#contenu"
        className="sr-only focus:not-sr-only focus:absolute focus:left-4 focus:top-4 focus:z-50 focus:rounded-md focus:bg-ink-900 focus:px-4 focus:py-2 focus:text-white"
      >
        {strings.skipToContent}
      </a>
      <SiteBanner />
      <Header />
      <main id="contenu">
        <Outlet />
      </main>
      <Footer />
      <ScrollRestoration />
    </LocaleProvider>
  )
}
