import { useLocale, t } from '../lib/i18n'

/**
 * Temporary "under construction" strip shown above the header. Remove this
 * component (and its use in Layout) when the site is ready to launch.
 */
export function SiteBanner() {
  const s = t(useLocale())
  return (
    <div role="status" className="bg-ink-950 text-center text-sm text-gold-300">
      <p className="mx-auto flex items-center justify-center gap-2 px-4 py-2">
        <span aria-hidden>🚧</span>
        <span className="font-medium">{s.banner}</span>
      </p>
    </div>
  )
}
