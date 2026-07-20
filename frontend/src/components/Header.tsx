import { useState } from 'react'
import { Link, NavLink } from 'vite-react-ssg'
import { useLocation } from 'react-router-dom'
import { Container } from './Container'
import { CtaLink } from './ui'
import { IconKnight } from './icons'
import { useLocale, t, homePath, altPath, PAGES, type Locale } from '../lib/i18n'

type NavItem = { to: string; key: keyof ReturnType<typeof t>['nav'] }

// FR links to every page; EN links only to pages already translated (grows
// as the i18n rollout continues). Labels come from the locale dictionary.
const NAV_FR: NavItem[] = [
  { to: '/cours-echecs-adultes-geneve', key: 'adultes' },
  { to: '/preparation-tournoi-echecs', key: 'tournoi' },
  { to: '/cours-echecs-en-ligne', key: 'enligne' },
  { to: '/tarifs', key: 'tarifs' },
  { to: '/resultats', key: 'resultats' },
  { to: '/blog', key: 'blog' },
]
const NAV_EN: NavItem[] = [
  { to: '/en/adult-chess-lessons-geneva', key: 'adultes' },
  { to: '/en/pricing', key: 'tarifs' },
]

function Brand({ locale, label }: { locale: Locale; label: string }) {
  return (
    <Link to={homePath(locale)} className="flex items-center gap-2.5 font-bold text-ink-900" aria-label={label}>
      <span aria-hidden className="flex h-9 w-9 items-center justify-center rounded-full bg-ink-900 text-gold-500 ring-1 ring-inset ring-gold-500/30">
        <IconKnight size={19} />
      </span>
      <span className="leading-none">
        <span className="font-display text-[1.05rem] tracking-tight">Alexandre&nbsp;Iwanesko</span>
        <span className="mt-0.5 block text-[0.62rem] font-semibold uppercase tracking-[0.18em] text-gold-700">
          {t(locale).footerRole}
        </span>
      </span>
    </Link>
  )
}

function LangSwitch({ locale, label }: { locale: Locale; label: string }) {
  const here = useLocation().pathname
  const other: Locale = locale === 'fr' ? 'en' : 'fr'
  const to = altPath(here) ?? homePath(other)
  return (
    <Link
      to={to}
      hrefLang={other}
      className="rounded-md border border-ink-200 px-2.5 py-1 text-xs font-bold uppercase tracking-wide text-ink-600 transition-colors hover:border-gold-400 hover:text-ink-950"
      aria-label={other === 'en' ? 'Switch to English' : 'Passer en français'}
    >
      {label}
    </Link>
  )
}

export function Header() {
  const [open, setOpen] = useState(false)
  const locale = useLocale()
  const s = t(locale)
  const nav = locale === 'en' ? NAV_EN : NAV_FR
  const contactPath = PAGES.contact[locale]

  return (
    <header className="sticky top-0 z-40 border-b border-ink-100 bg-white/85 backdrop-blur-md supports-[backdrop-filter]:bg-white/75">
      <Container className="flex h-[4.5rem] items-center justify-between gap-4">
        <Brand locale={locale} label={`Alexandre Iwanesko — ${s.breadcrumbHome}`} />

        <nav aria-label="Navigation" className="hidden items-center gap-7 lg:flex">
          {nav.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              className={({ isActive }) =>
                `relative text-sm font-medium transition-colors after:absolute after:-bottom-1.5 after:left-0 after:h-px after:bg-gold-500 after:transition-all after:duration-300 hover:text-ink-950 ${
                  isActive ? 'text-ink-950 after:w-full' : 'text-ink-600 after:w-0 hover:after:w-full'
                }`
              }
            >
              {s.nav[item.key]}
            </NavLink>
          ))}
          <LangSwitch locale={locale} label={s.langLabel} />
          <CtaLink to={contactPath} variant="primary" className="!px-5 !py-2.5">
            {s.reserve}
          </CtaLink>
        </nav>

        <div className="flex items-center gap-2 lg:hidden">
          <LangSwitch locale={locale} label={s.langLabel} />
          <button
            type="button"
            className="-mr-1 rounded-lg p-2 text-ink-800 transition-colors hover:bg-ink-50"
            aria-expanded={open}
            aria-controls="mobile-nav"
            aria-label={open ? s.closeMenu : s.openMenu}
            onClick={() => setOpen((v) => !v)}
          >
            {open ? (
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true">
                <path d="M6 6l12 12M18 6L6 18" />
              </svg>
            ) : (
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true">
                <path d="M4 7h16M4 12h16M4 17h16" />
              </svg>
            )}
          </button>
        </div>
      </Container>

      {open ? (
        <nav id="mobile-nav" aria-label="Navigation" className="border-t border-ink-100 bg-white lg:hidden">
          <Container className="flex flex-col gap-0.5 py-3">
            {nav.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                className={({ isActive }) =>
                  `rounded-lg px-3 py-2.5 text-[0.95rem] font-medium transition-colors ${
                    isActive ? 'bg-cream-100 text-ink-950' : 'text-ink-700 hover:bg-ink-50'
                  }`
                }
                onClick={() => setOpen(false)}
              >
                {s.nav[item.key]}
              </NavLink>
            ))}
            <CtaLink to={contactPath} variant="primary" className="mt-3 w-full">
              {s.reserve}
            </CtaLink>
          </Container>
        </nav>
      ) : null}
    </header>
  )
}
