import type { ReactNode } from 'react'
import { Link } from 'vite-react-ssg'

type Variant = 'primary' | 'secondary' | 'ghost'

const base =
  'group/btn inline-flex items-center justify-center gap-2 rounded-full font-semibold tracking-[0.01em] transition-[background-color,color,box-shadow,transform] duration-200 focus-visible:outline-none active:translate-y-px'
const sizes = 'px-6 py-3 text-[0.95rem]'
const variants: Record<Variant, string> = {
  primary: 'bg-gold-500 text-ink-950 shadow-soft hover:bg-gold-400 hover:shadow-gold',
  secondary: 'bg-ink-900 text-white shadow-soft hover:bg-ink-800',
  ghost:
    'text-ink-800 underline decoration-gold-400 decoration-1 underline-offset-4 hover:text-ink-950 hover:decoration-gold-600',
}

/** Internal-link CTA. Use descriptive anchor text (never "cliquez ici"). */
export function CtaLink({
  to,
  children,
  variant = 'primary',
  className = '',
}: {
  to: string
  children: ReactNode
  variant?: Variant
  className?: string
}) {
  return (
    <Link to={to} className={`${base} ${sizes} ${variants[variant]} ${className}`}>
      {children}
    </Link>
  )
}

export function Section({
  children,
  className = '',
  id,
}: {
  children: ReactNode
  className?: string
  id?: string
}) {
  return (
    <section id={id} className={`py-16 sm:py-24 ${className}`}>
      {children}
    </section>
  )
}

export function Eyebrow({ children, tone = 'light' }: { children: ReactNode; tone?: 'light' | 'dark' }) {
  return (
    <p
      className={`mb-4 flex items-center gap-3 text-xs font-semibold uppercase tracking-[0.2em] ${
        tone === 'dark' ? 'text-gold-400' : 'text-gold-700'
      }`}
    >
      <span aria-hidden className="h-px w-7 bg-gold-500" />
      {children}
    </p>
  )
}

export function Lead({ children }: { children: ReactNode }) {
  return <p className="max-w-[60ch] text-lg leading-relaxed text-ink-600">{children}</p>
}

/** A titled deep-dive section — used on money pages to add citable depth. */
export function DeepDive({
  eyebrow,
  title,
  children,
  className = '',
}: {
  eyebrow?: string
  title: string
  children: ReactNode
  className?: string
}) {
  return (
    <Section className={className}>
      <div className="mx-auto w-full max-w-6xl px-5 sm:px-8">
        {eyebrow ? <Eyebrow>{eyebrow}</Eyebrow> : null}
        <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">{title}</h2>
        <div className="mt-8">{children}</div>
      </div>
    </Section>
  )
}

/** Ordered list of labelled steps, e.g. a program/method breakdown. */
export function StepList({ steps }: { steps: Array<{ t: string; d: string }> }) {
  return (
    <ol className="mt-2 space-y-1">
      {steps.map((s, i) => (
        <li key={s.t} className="relative flex gap-5 pb-8 last:pb-0">
          {i < steps.length - 1 ? (
            <span aria-hidden className="absolute left-[1.35rem] top-11 bottom-2 w-px bg-ink-200" />
          ) : null}
          <span className="relative z-10 flex h-11 w-11 flex-none items-center justify-center rounded-full border border-gold-300 bg-gold-50 font-display text-lg font-bold text-gold-700">
            {i + 1}
          </span>
          <div className="pt-1.5">
            <h3 className="text-lg font-semibold text-ink-900">{s.t}</h3>
            <p className="mt-1.5 leading-relaxed text-ink-600">{s.d}</p>
          </div>
        </li>
      ))}
    </ol>
  )
}
