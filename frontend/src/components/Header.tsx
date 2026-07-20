import { useState } from 'react'
import { Link, NavLink } from 'vite-react-ssg'
import { Container } from './Container'
import { CtaLink } from './ui'
import { IconKnight } from './icons'

/** Primary nav. Kept small: links to money pages + key info pages. */
const NAV = [
  { to: '/cours-echecs-adultes-geneve', label: 'Cours adultes' },
  { to: '/preparation-tournoi-echecs', label: 'Préparation tournoi' },
  { to: '/cours-echecs-en-ligne', label: 'En ligne' },
  { to: '/tarifs', label: 'Tarifs' },
  { to: '/resultats', label: 'Résultats' },
  { to: '/blog', label: 'Blog' },
]

function Brand() {
  return (
    <Link
      to="/"
      className="flex items-center gap-2.5 font-bold text-ink-900"
      aria-label="Alexandre Iwanesko — accueil"
    >
      <span
        aria-hidden
        className="flex h-9 w-9 items-center justify-center rounded-full bg-ink-900 text-gold-500 ring-1 ring-inset ring-gold-500/30"
      >
        <IconKnight size={19} />
      </span>
      <span className="leading-none">
        <span className="font-display text-[1.05rem] tracking-tight">Alexandre&nbsp;Iwanesko</span>
        <span className="mt-0.5 block text-[0.62rem] font-semibold uppercase tracking-[0.18em] text-gold-700">
          Maître FIDE · Coach d’échecs
        </span>
      </span>
    </Link>
  )
}

export function Header() {
  const [open, setOpen] = useState(false)

  return (
    <header className="sticky top-0 z-40 border-b border-ink-100 bg-white/85 backdrop-blur-md supports-[backdrop-filter]:bg-white/75">
      <Container className="flex h-[4.5rem] items-center justify-between gap-4">
        <Brand />

        <nav aria-label="Navigation principale" className="hidden items-center gap-7 lg:flex">
          {NAV.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              className={({ isActive }) =>
                `relative text-sm font-medium transition-colors after:absolute after:-bottom-1.5 after:left-0 after:h-px after:bg-gold-500 after:transition-all after:duration-300 hover:text-ink-950 ${
                  isActive
                    ? 'text-ink-950 after:w-full'
                    : 'text-ink-600 after:w-0 hover:after:w-full'
                }`
              }
            >
              {item.label}
            </NavLink>
          ))}
          <CtaLink to="/contact" variant="primary" className="!px-5 !py-2.5">
            Réserver un cours
          </CtaLink>
        </nav>

        <button
          type="button"
          className="-mr-1 rounded-lg p-2 text-ink-800 transition-colors hover:bg-ink-50 lg:hidden"
          aria-expanded={open}
          aria-controls="mobile-nav"
          aria-label={open ? 'Fermer le menu' : 'Ouvrir le menu'}
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
      </Container>

      {open ? (
        <nav id="mobile-nav" aria-label="Navigation mobile" className="border-t border-ink-100 bg-white lg:hidden">
          <Container className="flex flex-col gap-0.5 py-3">
            {NAV.map((item) => (
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
                {item.label}
              </NavLink>
            ))}
            <CtaLink to="/contact" variant="primary" className="mt-3 w-full" >
              Réserver un cours
            </CtaLink>
          </Container>
        </nav>
      ) : null}
    </header>
  )
}
