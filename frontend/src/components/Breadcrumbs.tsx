import { Link } from 'vite-react-ssg'
import { Container } from './Container'
import type { Crumb } from '../lib/schema'

/**
 * Visible breadcrumb trail. The matching BreadcrumbList JSON-LD is emitted by
 * each page via the Seo component (keep the two crumb arrays identical).
 */
export function Breadcrumbs({ crumbs }: { crumbs: Crumb[] }) {
  return (
    <nav aria-label="Fil d’Ariane" className="border-b border-ink-100 bg-cream-50">
      <Container>
        <ol className="flex flex-wrap items-center gap-x-2 gap-y-1 py-3.5 text-sm text-ink-500">
          {crumbs.map((c, i) => {
            const last = i === crumbs.length - 1
            return (
              <li key={c.path} className="flex items-center gap-2">
                {last ? (
                  <span aria-current="page" className="font-medium text-ink-700">
                    {c.name}
                  </span>
                ) : (
                  <Link to={c.path} className="transition-colors hover:text-gold-700">
                    {c.name}
                  </Link>
                )}
                {!last && <span aria-hidden className="text-ink-300">/</span>}
              </li>
            )
          })}
        </ol>
      </Container>
    </nav>
  )
}
