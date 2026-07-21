import { useParams } from 'react-router-dom'
import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { PuzzleBoard } from '../components/PuzzleBoard'
import { getWeek, formatWeek } from '../lib/tactics'
import { recordTacticsEvent } from '../lib/tacticsEvents'
import { breadcrumbSchema, type Crumb } from '../lib/schema'
import { useLocale, homePath, t, type Locale } from '../lib/i18n'

const STR: Record<Locale, {
  eyebrow: string; titlePrefix: string; lead: string; intro: string
  whiteWins: string; blackWins: string; mate: string; sac: string; notFound: string; back: string
  board: { yourMove: string; solved: string; tryAgain: string; retry: string; whiteToPlay: string; blackToPlay: string }
}> = {
  fr: {
    eyebrow: 'Tactiques de la semaine', titlePrefix: 'Tactiques du',
    lead: 'Les plus belles combinaisons de la semaine — à résoudre directement sur l’échiquier.',
    intro: 'Chaque position vient d’une vraie partie. Cliquez la pièce puis sa case d’arrivée pour jouer la solution.',
    whiteWins: 'Les Blancs jouent et gagnent', blackWins: 'Les Noirs jouent et gagnent',
    mate: 'Mat', sac: 'Sacrifice', notFound: 'Cette série n’existe pas.', back: 'Toutes les semaines',
    board: { yourMove: 'à vous de jouer', solved: 'Résolu !', tryAgain: 'Essayez encore', retry: 'Recommencer', whiteToPlay: 'Les Blancs jouent', blackToPlay: 'Les Noirs jouent' },
  },
  en: {
    eyebrow: 'Tactics of the week', titlePrefix: 'Tactics of',
    lead: 'The best combinations of the week — solve them right on the board.',
    intro: 'Each position is from a real game. Click the piece, then its destination square, to play the solution.',
    whiteWins: 'White to play and win', blackWins: 'Black to play and win',
    mate: 'Mate', sac: 'Sacrifice', notFound: 'This set does not exist.', back: 'All weeks',
    board: { yourMove: 'your move', solved: 'Solved!', tryAgain: 'Try again', retry: 'Restart', whiteToPlay: 'White to play', blackToPlay: 'Black to play' },
  },
}

export function Component() {
  const locale = useLocale()
  const s = STR[locale]
  const { date = '' } = useParams()
  const week = getWeek(date)
  const indexPath = locale === 'en' ? '/en/tactics' : '/tactiques'

  if (!week) {
    return (
      <Section>
        <Container>
          <Seo title={s.eyebrow} description={s.lead} path={`/tactiques/${date}`} noindex />
          <h1 className="font-display text-2xl font-bold text-ink-900">{s.notFound}</h1>
          <p className="mt-3"><Link to={indexPath} className="text-gold-700 underline">{s.back}</Link></p>
        </Container>
      </Section>
    )
  }

  const label = formatWeek(week.slug, locale)
  const path = `${indexPath}/${week.slug}`
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: s.eyebrow, path: indexPath },
    { name: label, path },
  ]

  return (
    <>
      <Seo
        title={`${s.titlePrefix} ${label}`}
        description={`${s.lead} — ${label}.`}
        path={path}
        jsonLd={[breadcrumbSchema(crumbs)]}
      />
      <Breadcrumbs crumbs={crumbs} />
      <PageHero eyebrow={s.eyebrow} title={`${s.titlePrefix} ${label}`} lead={s.lead} primaryCta={{ to: locale === 'en' ? '/en/contact' : '/contact', label: locale === 'en' ? 'Book a lesson' : 'Réserver un cours' }} />

      <Section>
        <Container>
          <p className="mb-8 max-w-2xl leading-relaxed text-ink-600">{s.intro}</p>
          <ol className="grid gap-10 sm:grid-cols-2">
            {week.puzzles.map((p, i) => (
              <li key={p.id} className="rounded-2xl border border-ink-200/80 bg-white p-5 shadow-soft sm:p-6">
                <div className="mb-3 flex items-center justify-between">
                  <span className="font-display text-lg font-bold text-ink-900">#{i + 1}</span>
                  <span className="flex gap-2">
                    {p.mate ? <span className="rounded-full bg-ink-900 px-2.5 py-0.5 text-xs font-semibold text-gold-300">{s.mate}</span> : null}
                    {p.sacrifice ? <span className="rounded-full bg-gold-100 px-2.5 py-0.5 text-xs font-semibold text-gold-700">{s.sac}</span> : null}
                  </span>
                </div>
                <p className="mb-3 text-sm font-medium text-ink-600">{p.sideToMove === 'w' ? s.whiteWins : s.blackWins}</p>
                <PuzzleBoard
                  fen={p.fen}
                  sideToMove={p.sideToMove}
                  solution={p.solution}
                  labels={s.board}
                  onView={() => recordTacticsEvent(week.slug, p.id, 'view')}
                  onAttempt={(correct) => { if (correct) recordTacticsEvent(week.slug, p.id, 'attempt') }}
                  onSolved={() => recordTacticsEvent(week.slug, p.id, 'solved')}
                />
              </li>
            ))}
          </ol>
        </Container>
      </Section>
    </>
  )
}
