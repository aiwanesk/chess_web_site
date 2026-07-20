import { Container } from './Container'
import { Section, Eyebrow } from './ui'
import type { FaqItem } from '../lib/schema'

/**
 * FAQ block for the bottom of every money page. Questions as headings +
 * concise, self-contained answers up top = ideal for GEO extraction. The
 * matching FAQPage JSON-LD is emitted via the Seo component.
 *
 * Uses native <details>/<summary> so it works with zero JS (progressive
 * enhancement) and stays accessible.
 */
export function Faq({ items, title = 'Questions fréquentes' }: { items: FaqItem[]; title?: string }) {
  return (
    <Section className="bg-cream-100">
      <Container>
        <Eyebrow>FAQ</Eyebrow>
        <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">{title}</h2>
        <div className="mx-auto mt-8 max-w-4xl space-y-3">
          {items.map((item) => (
            <details
              key={item.question}
              className="group rounded-2xl border border-ink-200/80 bg-white px-6 shadow-soft transition-colors open:border-gold-300"
            >
              <summary className="flex cursor-pointer list-none items-center justify-between gap-4 py-5 text-left text-lg font-semibold text-ink-900">
                {item.question}
                <span
                  aria-hidden
                  className="flex h-8 w-8 flex-none items-center justify-center rounded-full border border-gold-300 bg-gold-50 text-gold-700 transition-transform duration-300 group-open:rotate-45"
                >
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
                    <path d="M12 5v14M5 12h14" />
                  </svg>
                </span>
              </summary>
              <p className="max-w-3xl pb-6 leading-relaxed text-ink-600">{item.answer}</p>
            </details>
          ))}
        </div>
      </Container>
    </Section>
  )
}
