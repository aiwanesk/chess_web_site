import { useLocale, t } from '../lib/i18n'
import { toggleTheme } from '../lib/theme'

/**
 * Bouton jour / nuit.
 *
 * Aucun état React, et c'est délibéré : sur un site pré-rendu, un bouton dont
 * l'icône dépend d'un `useState` se rend en clair côté serveur puis saute au
 * bon état à l'hydratation. Ici les deux icônes sont dans le DOM et c'est le
 * CSS, via `html[data-theme]`, qui montre la bonne — donc pas de saut, et
 * l'icône est déjà juste au premier pixel peint.
 */
export function ThemeToggle() {
  const s = t(useLocale())
  return (
    <button
      type="button"
      onClick={() => toggleTheme()}
      className="rounded-md border border-ink-200 p-1.5 text-ink-600 transition-colors hover:border-gold-400 hover:text-ink-950"
      aria-label={s.themeToggle}
      title={s.themeToggle}
    >
      {/* Lune = « passer en nuit », visible tant qu'on est en clair. */}
      <svg
        className="theme-icon-moon"
        width="17"
        height="17"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
        aria-hidden="true"
      >
        <path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z" />
      </svg>
      {/* Soleil = « repasser en clair », visible en mode nuit. */}
      <svg
        className="theme-icon-sun"
        width="17"
        height="17"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
        aria-hidden="true"
      >
        <circle cx="12" cy="12" r="4" />
        <path d="M12 2v2m0 16v2M2 12h2m16 0h2M4.9 4.9l1.4 1.4m11.4 11.4 1.4 1.4M19.1 4.9l-1.4 1.4M6.3 17.7l-1.4 1.4" />
      </svg>
    </button>
  )
}
