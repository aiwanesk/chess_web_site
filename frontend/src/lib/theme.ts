/**
 * Thème clair / sombre.
 *
 * L'état vit sur `<html data-theme>`, posé par le script inline de index.html
 * AVANT le premier rendu — c'est ce qui évite le flash blanc au chargement sur
 * un site pré-rendu. Rien n'est stocké dans React : le composant se contente de
 * retourner l'attribut, et le CSS fait le reste (voir styles.css).
 */
export type Theme = 'light' | 'dark'

export const THEME_KEY = 'theme'

/** Couleur de la barre d'adresse mobile, accordée au fond de page. */
const BAR = { light: '#ffffff', dark: '#0d131f' } as const

export function currentTheme(): Theme {
  return document.documentElement.dataset.theme === 'dark' ? 'dark' : 'light'
}

export function applyTheme(theme: Theme): void {
  document.documentElement.dataset.theme = theme
  document.querySelector('meta[name="theme-color"]')?.setAttribute('content', BAR[theme])
}

/** Bascule et mémorise le choix. Le localStorage peut lever (navigation privée,
 * cookies bloqués) : le thème s'applique quand même, il ne survit juste pas. */
export function toggleTheme(): Theme {
  const next: Theme = currentTheme() === 'dark' ? 'light' : 'dark'
  applyTheme(next)
  try {
    localStorage.setItem(THEME_KEY, next)
  } catch {
    /* pas de persistance possible — tant pis, la page reste dans le bon thème */
  }
  return next
}
