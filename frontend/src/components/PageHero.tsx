import type { ReactNode } from 'react'
import { Container } from './Container'
import { Eyebrow, CtaLink } from './ui'

/**
 * Reusable hero. The <h1> here is the single H1 of every page. Answer-first
 * copy (the lead) doubles as a citable definition for generative engines.
 */
export function PageHero({
  eyebrow,
  title,
  lead,
  primaryCta = { to: '/contact', label: 'Réserver un premier cours' },
  secondaryCta,
  aside,
}: {
  eyebrow?: string
  title: string
  lead: ReactNode
  primaryCta?: { to: string; label: string }
  secondaryCta?: { to: string; label: string }
  aside?: ReactNode
}) {
  return (
    <div className="relative overflow-hidden border-b border-ink-100 bg-gradient-to-b from-cream-100 to-white">
      <div aria-hidden className="board-texture absolute inset-0 opacity-70" />
      {/* Soft gold glow, top-right. */}
      <div
        aria-hidden
        className="pointer-events-none absolute -right-24 -top-24 h-72 w-72 rounded-full bg-gold-300/20 blur-3xl"
      />
      <Container className="relative grid gap-12 py-16 sm:py-24 lg:grid-cols-[1.35fr_1fr] lg:items-center">
        <div>
          {eyebrow ? <Eyebrow>{eyebrow}</Eyebrow> : null}
          <h1 className="font-display text-[2rem] font-bold leading-[1.08] tracking-tight text-ink-900 sm:text-[2.6rem] lg:text-[3.1rem]">
            {title}
          </h1>
          <div className="mt-6 max-w-[56ch] text-lg leading-relaxed text-ink-600">{lead}</div>
          <div className="mt-9 flex flex-wrap items-center gap-3">
            <CtaLink to={primaryCta.to} variant="primary">
              {primaryCta.label}
            </CtaLink>
            {secondaryCta ? (
              <CtaLink to={secondaryCta.to} variant="ghost">
                {secondaryCta.label}
              </CtaLink>
            ) : null}
          </div>
        </div>
        {aside ? <div className="lg:pl-4">{aside}</div> : null}
      </Container>
    </div>
  )
}

/** Compact "key facts" card — specific, citable numbers the LLMs love. */
export function FactCard({ facts }: { facts: Array<{ label: string; value: string }> }) {
  return (
    <dl className="grid grid-cols-2 gap-px overflow-hidden rounded-2xl border border-ink-200/80 bg-ink-200/70 shadow-card">
      {facts.map((f) => (
        <div key={f.label} className="bg-white p-5 transition-colors hover:bg-cream-50">
          <dt className="text-[0.68rem] font-semibold uppercase tracking-[0.12em] text-ink-500">{f.label}</dt>
          <dd className="mt-1.5 font-display text-lg font-bold text-ink-900">{f.value}</dd>
        </div>
      ))}
    </dl>
  )
}
