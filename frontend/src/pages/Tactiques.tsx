import { useEffect, useState } from 'react'
import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { PuzzleBoard } from '../components/PuzzleBoard'
import { breadcrumbSchema, type Crumb } from '../lib/schema'
import { useLocale, homePath, t, type Locale } from '../lib/i18n'

interface Puzzle {
  id: string
  fen: string
  sideToMove: 'w' | 'b'
  solution: string[]
  mate?: boolean
  sacrifice?: boolean
}

const STR: Record<Locale, {
  title: string; desc: string; eyebrow: string; heroTitle: string; heroLead: string
  loading: string; empty: string; whiteWins: string; blackWins: string; mate: string; sac: string
  board: { yourMove: string; solved: string; tryAgain: string; retry: string; whiteToPlay: string; blackToPlay: string }
}> = {
  fr: {
    title: 'Tactiques de la semaine', desc: 'Résolvez les plus belles tactiques d’échecs de la semaine, sélectionnées par un Maître FIDE. Nouveaux puzzles chaque semaine.',
    eyebrow: 'Puzzles', heroTitle: 'Les tactiques de la semaine',
    heroLead: 'Une sélection des plus belles combinaisons, à résoudre directement sur l’échiquier. De nouveaux puzzles chaque semaine.',
    loading: 'Chargement des puzzles…', empty: 'Les puzzles de la semaine arrivent bientôt.',
    whiteWins: 'Les Blancs jouent et gagnent', blackWins: 'Les Noirs jouent et gagnent', mate: 'Mat', sac: 'Sacrifice',
    board: { yourMove: 'à vous de jouer', solved: 'Résolu !', tryAgain: 'Essayez encore', retry: 'Recommencer', whiteToPlay: 'Les Blancs jouent', blackToPlay: 'Les Noirs jouent' },
  },
  en: {
    title: 'Tactics of the week', desc: 'Solve the week’s best chess tactics, hand-picked by a FIDE Master. New puzzles every week.',
    eyebrow: 'Puzzles', heroTitle: 'Tactics of the week',
    heroLead: 'A selection of the best combinations to solve right on the board. New puzzles every week.',
    loading: 'Loading puzzles…', empty: 'This week’s puzzles are coming soon.',
    whiteWins: 'White to play and win', blackWins: 'Black to play and win', mate: 'Mate', sac: 'Sacrifice',
    board: { yourMove: 'your move', solved: 'Solved!', tryAgain: 'Try again', retry: 'Restart', whiteToPlay: 'White to play', blackToPlay: 'Black to play' },
  },
}

export function Component() {
  const locale = useLocale()
  const s = STR[locale]
  const path = locale === 'en' ? '/en/tactics' : '/tactiques'
  const [puzzles, setPuzzles] = useState<Puzzle[] | null>(null)

  useEffect(() => {
    let ok = true
    fetch('/api/tactics')
      .then((r) => r.json())
      .then((d) => { if (ok) setPuzzles(Array.isArray(d.puzzles) ? d.puzzles : []) })
      .catch(() => { if (ok) setPuzzles([]) })
    return () => { ok = false }
  }, [])

  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: s.title, path },
  ]

  return (
    <>
      <Seo title={s.title} description={s.desc} path={path} jsonLd={[breadcrumbSchema(crumbs)]} />
      <Breadcrumbs crumbs={crumbs} />
      <PageHero eyebrow={s.eyebrow} title={s.heroTitle} lead={s.heroLead} primaryCta={{ to: locale === 'en' ? '/en/contact' : '/contact', label: locale === 'en' ? 'Book a lesson' : 'Réserver un cours' }} />

      <Section>
        <Container>
          {puzzles === null ? (
            <p className="text-center text-ink-500">{s.loading}</p>
          ) : puzzles.length === 0 ? (
            <p className="text-center text-ink-500">{s.empty}</p>
          ) : (
            <ol className="grid gap-10 sm:grid-cols-2">
              {puzzles.map((p, i) => (
                <li key={p.id} className="rounded-2xl border border-ink-200/80 bg-white p-5 shadow-soft sm:p-6">
                  <div className="mb-3 flex items-center justify-between">
                    <span className="font-display text-lg font-bold text-ink-900">#{i + 1}</span>
                    <span className="flex gap-2">
                      {p.mate ? <span className="rounded-full bg-ink-900 px-2.5 py-0.5 text-xs font-semibold text-gold-300">{s.mate}</span> : null}
                      {p.sacrifice ? <span className="rounded-full bg-gold-100 px-2.5 py-0.5 text-xs font-semibold text-gold-700">{s.sac}</span> : null}
                    </span>
                  </div>
                  <p className="mb-3 text-sm font-medium text-ink-600">{p.sideToMove === 'w' ? s.whiteWins : s.blackWins}</p>
                  <PuzzleBoard fen={p.fen} sideToMove={p.sideToMove} solution={p.solution} labels={s.board} />
                </li>
              ))}
            </ol>
          )}
        </Container>
      </Section>
    </>
  )
}
