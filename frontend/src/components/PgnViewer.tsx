import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { PieceSvg } from './chessPieces'
import { getGame } from '../lib/pgn'
import { toFrenchSan } from '../lib/chess'
import { useLocale, type Locale } from '../lib/i18n'

/**
 * Read-only game replayer for the tournament diaries. The score is compiled at
 * build time (see lib/pgn.ts), so the pre-rendered HTML already contains the
 * starting position and the whole move list — the client only steps through it.
 *
 * Main line only, both colours from either side (flip button), keyboard arrows.
 */

const FILE_LETTERS = 'abcdefgh'
const LIGHT = '#f0d9b5'
const DARK = '#b58863'

const STR: Record<Locale, {
  flip: string; first: string; prev: string; next: string; last: string
  play: string; pause: string; missing: string; moves: string; startPos: string; games: string
}> = {
  fr: {
    flip: 'Retourner l’échiquier',
    first: 'Position de départ',
    prev: 'Coup précédent',
    next: 'Coup suivant',
    last: 'Dernier coup',
    play: 'Lecture automatique',
    pause: 'Pause',
    missing: 'Partie introuvable.',
    moves: 'Liste des coups',
    startPos: 'Début de la partie',
    games: 'Choisir la partie',
  },
  en: {
    flip: 'Flip the board',
    first: 'Starting position',
    prev: 'Previous move',
    next: 'Next move',
    last: 'Last move',
    play: 'Autoplay',
    pause: 'Pause',
    missing: 'Game not found.',
    moves: 'Move list',
    startPos: 'Start of the game',
    games: 'Choose the game',
  },
}

export interface PgnViewerProps {
  /** File names (without .pgn) under /content/games. Several ids = one viewer with a picker. */
  gameIds: string[]
  caption?: string
  /** Anchor target, so an article can link down to the games from its intro. */
  anchorId?: string
}

/** Tab label for the picker: the opponent's surname, prefixed by the round when known. */
function tabLabel(id: string, index: number): string {
  const g = getGame(id)
  if (!g) return `#${index + 1}`
  const opponent = /iwanesko/i.test(g.white) ? g.black : g.white
  const surname = opponent.trim().split(/\s+/).pop() || opponent
  return g.round ? `R${g.round} · ${surname}` : surname
}

export function PgnViewer({ gameIds, caption, anchorId }: PgnViewerProps) {
  const locale = useLocale()
  const s = STR[locale]

  // `ply` counts moves played: 0 = starting position, 1 = after White's first.
  const [active, setActive] = useState(0)
  const [ply, setPly] = useState(0)
  const [flipped, setFlipped] = useState(false)
  const [playing, setPlaying] = useState(false)
  const listRef = useRef<HTMLDivElement>(null)

  const game = getGame(gameIds[active] ?? gameIds[0] ?? '')

  // Switching games restarts from the initial position, on the new player's side.
  const select = useCallback((i: number) => {
    setActive(i)
    setPly(0)
    setFlipped(false)
    setPlaying(false)
  }, [])

  const total = game?.plies.length ?? 0
  const go = useCallback((n: number) => setPly(Math.max(0, Math.min(total, n))), [total])

  // Autoplay stops on its own at the end of the game.
  useEffect(() => {
    if (!playing) return
    if (ply >= total) {
      setPlaying(false)
      return
    }
    const id = window.setTimeout(() => setPly((p) => p + 1), 900)
    return () => clearTimeout(id)
  }, [playing, ply, total])

  // Keep the current move visible in the (scrollable) list — and ONLY move that
  // list. scrollIntoView() scrolls every scrollable ancestor, the page included:
  // on mount it dragged the reader straight down to the board at the foot of a
  // diary, which also made the #parties link look broken (you were already there).
  useEffect(() => {
    const list = listRef.current
    const el = list?.querySelector<HTMLElement>('[data-current="true"]')
    if (!list || !el) return
    const box = list.getBoundingClientRect()
    const move = el.getBoundingClientRect()
    if (move.top < box.top) list.scrollTop -= box.top - move.top
    else if (move.bottom > box.bottom) list.scrollTop += move.bottom - box.bottom
  }, [ply])

  const board = useMemo(
    () => (ply === 0 ? game?.startBoard : game?.plies[ply - 1]?.board),
    [game, ply],
  )

  if (!game || !board) {
    return <p className="mx-auto max-w-[68ch] text-sm text-ink-500">{s.missing}</p>
  }

  const current = ply > 0 ? game.plies[ply - 1]! : null
  const orientationWhite = game.orientation === 'w' ? !flipped : flipped
  const rankOrder = orientationWhite ? [7, 6, 5, 4, 3, 2, 1, 0] : [0, 1, 2, 3, 4, 5, 6, 7]
  const fileOrder = orientationWhite ? [0, 1, 2, 3, 4, 5, 6, 7] : [7, 6, 5, 4, 3, 2, 1, 0]
  const label = (san: string) => (locale === 'fr' ? toFrenchSan(san) : san)

  // Player strip: whoever is at the bottom of the board is shown underneath.
  const topPlayer = orientationWhite ? { name: game.black, elo: game.blackElo } : { name: game.white, elo: game.whiteElo }
  const bottomPlayer = orientationWhite ? { name: game.white, elo: game.whiteElo } : { name: game.black, elo: game.blackElo }

  function onKeyDown(e: React.KeyboardEvent) {
    const keys: Record<string, () => void> = {
      ArrowLeft: () => go(ply - 1),
      ArrowRight: () => go(ply + 1),
      Home: () => go(0),
      End: () => go(total),
    }
    const action = keys[e.key]
    if (!action) return
    e.preventDefault()
    setPlaying(false)
    action()
  }

  const NameRow = ({ player }: { player: { name: string; elo?: string } }) => (
    <p className="truncate text-sm font-semibold text-ink-800">
      {player.name}
      {player.elo ? <span className="ml-1.5 font-normal text-ink-500">{player.elo}</span> : null}
    </p>
  )

  return (
    <figure id={anchorId} className="mx-auto my-10 max-w-3xl scroll-mt-24">
      <div
        tabIndex={0}
        onKeyDown={onKeyDown}
        role="group"
        aria-label={`${game.white} – ${game.black}`}
        className="rounded-2xl border border-ink-200 bg-paper p-4 shadow-sm focus:outline-none focus-visible:ring-2 focus-visible:ring-gold-400 sm:p-5"
      >
        {gameIds.length > 1 ? (
          <div role="tablist" aria-label={s.games} className="mb-3 flex flex-wrap gap-1.5">
            {gameIds.map((id, i) => (
              <button
                key={id}
                type="button"
                role="tab"
                aria-selected={i === active}
                onClick={() => select(i)}
                className={`rounded-full px-3 py-1 text-xs font-semibold transition-colors ${
                  i === active
                    ? 'bg-slab-900 text-white'
                    : 'bg-ink-100 text-ink-600 hover:bg-ink-200 hover:text-ink-900'
                }`}
              >
                {tabLabel(id, i)}
              </button>
            ))}
          </div>
        ) : null}

        <header className="mb-4 flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1">
          <p className="text-sm font-semibold text-ink-900">
            {game.white} <span className="font-normal text-ink-500">–</span> {game.black}
          </p>
          <p className="text-xs text-ink-500">
            {[game.event, game.round ? `${locale === 'fr' ? 'ronde' : 'round'} ${game.round}` : null, game.result]
              .filter(Boolean)
              .join(' · ')}
          </p>
        </header>

        <div className="flex flex-col gap-4 sm:flex-row">
          <div className="sm:w-[min(22rem,55%)] sm:shrink-0">
            <NameRow player={topPlayer} />
            <div
              className="mt-1.5 grid aspect-square w-full select-none grid-cols-8 grid-rows-8 overflow-hidden rounded-lg ring-1 ring-ink-900/10"
              aria-hidden
            >
              {rankOrder.map((rank, r) =>
                fileOrder.map((file, f) => {
                  const sq = rank * 8 + file
                  const piece = board[sq]
                  const dark = (file + rank) % 2 === 0
                  const base = dark ? DARK : LIGHT
                  const coord = dark ? LIGHT : DARK
                  const isLast = current && (sq === current.from || sq === current.to)
                  return (
                    <div key={sq} className="relative flex items-center justify-center" style={{ backgroundColor: base }}>
                      {isLast ? (
                        <span className="absolute inset-0" style={{ backgroundColor: 'rgba(155,199,0,0.41)' }} />
                      ) : null}
                      {f === 0 ? (
                        <span
                          className="absolute left-[3px] top-[2px] text-[9px] font-bold leading-none sm:text-[10px]"
                          style={{ color: coord }}
                        >
                          {rank + 1}
                        </span>
                      ) : null}
                      {r === 7 ? (
                        <span
                          className="absolute bottom-[1px] right-[3px] text-[9px] font-bold leading-none sm:text-[10px]"
                          style={{ color: coord }}
                        >
                          {FILE_LETTERS[file]}
                        </span>
                      ) : null}
                      {piece ? <PieceSvg piece={piece} className="relative h-[92%] w-[92%]" /> : null}
                    </div>
                  )
                }),
              )}
            </div>
            <div className="mt-1.5 flex items-baseline justify-between gap-3">
              <NameRow player={bottomPlayer} />
              <span className="shrink-0 text-xs font-semibold tabular-nums text-ink-500">
                {ply}/{total}
              </span>
            </div>
          </div>

          <div className="flex min-w-0 flex-1 flex-col">
            <div
              ref={listRef}
              className="order-2 mt-3 max-h-40 overflow-y-auto rounded-lg bg-ink-50 p-2 text-sm leading-7 sm:order-1 sm:mt-0 sm:max-h-[19rem]"
              aria-label={s.moves}
            >
              <button
                type="button"
                onClick={() => go(0)}
                data-current={ply === 0}
                className={`mr-1 rounded px-1.5 py-0.5 text-xs font-medium transition-colors ${
                  ply === 0 ? 'bg-slab-900 text-white' : 'text-ink-500 hover:bg-ink-200'
                }`}
              >
                {s.startPos}
              </button>
              {game.plies.map((p, i) => (
                <span key={i}>
                  {p.color === 'w' ? (
                    <span className="ml-1 text-xs font-semibold text-ink-400 tabular-nums">{p.number}.</span>
                  ) : null}
                  <button
                    type="button"
                    onClick={() => {
                      setPlaying(false)
                      go(i + 1)
                    }}
                    data-current={ply === i + 1}
                    className={`ml-0.5 rounded px-1 py-0.5 font-medium transition-colors ${
                      ply === i + 1 ? 'bg-gold-500 text-on-gold' : 'text-ink-700 hover:bg-ink-200'
                    }`}
                  >
                    {label(p.san)}
                  </button>
                </span>
              ))}
              <span className="ml-2 text-xs font-semibold text-ink-500">{game.result}</span>
            </div>

            <div className="order-1 flex flex-wrap items-center gap-1.5 sm:order-2 sm:mt-3">
              <NavButton label={s.first} onClick={() => { setPlaying(false); go(0) }} disabled={ply === 0}>‹‹</NavButton>
              <NavButton label={s.prev} onClick={() => { setPlaying(false); go(ply - 1) }} disabled={ply === 0}>‹</NavButton>
              <NavButton
                label={playing ? s.pause : s.play}
                onClick={() => setPlaying((p) => !p)}
                disabled={ply >= total && !playing}
              >
                {playing ? '❚❚' : '▶'}
              </NavButton>
              <NavButton label={s.next} onClick={() => { setPlaying(false); go(ply + 1) }} disabled={ply >= total}>›</NavButton>
              <NavButton label={s.last} onClick={() => { setPlaying(false); go(total) }} disabled={ply >= total}>››</NavButton>
              <button
                type="button"
                onClick={() => setFlipped((v) => !v)}
                title={s.flip}
                className="ml-auto rounded-md border border-ink-200 px-2.5 py-1 text-xs font-medium text-ink-600 transition-colors hover:border-gold-400 hover:text-ink-900"
              >
                ⇅ <span className="sr-only sm:not-sr-only">{s.flip}</span>
              </button>
            </div>
          </div>
        </div>
      </div>
      {caption ? <figcaption className="mt-2 text-center text-sm text-ink-500">{caption}</figcaption> : null}
    </figure>
  )
}

function NavButton({
  label,
  onClick,
  disabled,
  children,
}: {
  label: string
  onClick: () => void
  disabled?: boolean
  children: React.ReactNode
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      title={label}
      aria-label={label}
      className="rounded-md border border-ink-200 px-2.5 py-1 text-sm font-medium text-ink-700 transition-colors hover:border-gold-400 hover:text-ink-950 disabled:cursor-default disabled:opacity-35 disabled:hover:border-ink-200"
    >
      {children}
    </button>
  )
}
