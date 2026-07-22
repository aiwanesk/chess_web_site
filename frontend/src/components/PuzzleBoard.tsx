import { useEffect, useMemo, useRef, useState } from 'react'

/**
 * Lightweight interactive chess puzzle, styled like a lichess board. No
 * client-side engine: the solution line is known (trusted, from Stockfish), so
 * we only compare the player's move to the expected one and animate the
 * position. Fully anonymised data (mirrored FEN + solution).
 */
export interface PuzzleBoardProps {
  fen: string
  sideToMove: 'w' | 'b' // solver's colour
  solution: string[] // UCI moves, solver plays even indices
  onSolved?: () => void
  onAttempt?: (correct: boolean) => void
  onView?: () => void
  labels: {
    yourMove: string
    solved: string
    tryAgain: string
    retry: string
    whiteToPlay: string
    blackToPlay: string
    showSolution: string
    solutionShown: string
  }
}

type Pieces = Record<string, string>
const FILES = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h']
const GLYPH: Record<string, string> = { p: '♟', n: '♞', b: '♝', r: '♜', q: '♛', k: '♚' }

function parseFen(fen: string): Pieces {
  const rows = fen.split(' ')[0]!.split('/')
  const pieces: Pieces = {}
  rows.forEach((row, r) => {
    const rank = 8 - r
    let file = 0
    for (const ch of row) {
      if (/\d/.test(ch)) file += Number(ch)
      else {
        pieces[FILES[file]! + rank] = ch
        file++
      }
    }
  })
  return pieces
}

// applyMove mutates a copy for display only (position is trusted, not validated).
function applyMove(prev: Pieces, uci: string): Pieces {
  const p = { ...prev }
  const from = uci.slice(0, 2)
  const to = uci.slice(2, 4)
  const promo = uci[4]
  const piece = p[from]
  if (!piece) return p
  delete p[from]

  const isPawn = piece.toLowerCase() === 'p'
  // En passant: pawn moves diagonally onto an empty square → remove passed pawn.
  if (isPawn && from[0] !== to[0] && !prev[to]) {
    delete p[to[0]! + from[1]]
  }
  // Castling: king moves two files → move the rook too.
  if (piece.toLowerCase() === 'k' && Math.abs(from.charCodeAt(0) - to.charCodeAt(0)) === 2) {
    const rank = from[1]!
    if (to[0] === 'g') {
      p['f' + rank] = p['h' + rank]!
      delete p['h' + rank]
    } else if (to[0] === 'c') {
      p['d' + rank] = p['a' + rank]!
      delete p['a' + rank]
    }
  }
  p[to] = promo ? (piece === piece.toUpperCase() ? promo.toUpperCase() : promo.toLowerCase()) : piece
  return p
}

function applyMoves(fen: string, moves: string[]): Pieces {
  return moves.reduce((acc, m) => applyMove(acc, m), parseFen(fen))
}

export function PuzzleBoard({ fen, sideToMove, solution, onSolved, onAttempt, onView, labels }: PuzzleBoardProps) {
  const [applied, setApplied] = useState(0)
  const [selected, setSelected] = useState<string | null>(null)
  const [wrong, setWrong] = useState(false)
  const [revealed, setRevealed] = useState(false)
  const timers = useRef<number[]>([])

  // Count one "view" the first time the board is mounted in the browser.
  const viewed = useRef(false)
  useEffect(() => {
    if (viewed.current) return
    viewed.current = true
    onView?.()
  }, [onView])

  // Clear any pending reveal animations on unmount.
  useEffect(() => () => timers.current.forEach(clearTimeout), [])

  const pieces = useMemo(() => applyMoves(fen, solution.slice(0, applied)), [fen, solution, applied])
  const solved = applied >= solution.length
  const solverTurn = applied % 2 === 0 && !solved && !revealed
  const finishedByUser = solved && !revealed

  // Highlight the squares of the last move played (lichess-style).
  const last = applied > 0 ? solution[applied - 1]! : null
  const lastFrom = last ? last.slice(0, 2) : null
  const lastTo = last ? last.slice(2, 4) : null

  const ranks = sideToMove === 'w' ? [8, 7, 6, 5, 4, 3, 2, 1] : [1, 2, 3, 4, 5, 6, 7, 8]
  const files = sideToMove === 'w' ? FILES : [...FILES].reverse()
  const isSolver = (piece: string) => (sideToMove === 'w' ? piece === piece.toUpperCase() : piece === piece.toLowerCase())

  function reset() {
    timers.current.forEach(clearTimeout)
    timers.current = []
    setApplied(0)
    setSelected(null)
    setWrong(false)
    setRevealed(false)
  }

  function revealSolution() {
    setSelected(null)
    setWrong(false)
    setRevealed(true)
    let step = applied
    const advance = () => {
      step += 1
      setApplied(step)
      if (step < solution.length) timers.current.push(window.setTimeout(advance, 650))
    }
    if (step < solution.length) timers.current.push(window.setTimeout(advance, 350))
  }

  // attemptMove compares a from→to move to the expected solution move and
  // advances (or flashes wrong). Shared by click-to-move and drag-and-drop.
  function attemptMove(from: string, to: string) {
    const expected = solution[applied]!
    if (expected.slice(0, 4) === from + to) {
      onAttempt?.(true)
      setSelected(null)
      setWrong(false)
      const next = applied + 1
      setApplied(next)
      if (next >= solution.length) onSolved?.()
      else timers.current.push(window.setTimeout(() => setApplied(next + 1), 350)) // opponent's forced reply
    } else {
      onAttempt?.(false)
      setWrong(true)
      setSelected(null)
      timers.current.push(window.setTimeout(() => setWrong(false), 600))
    }
  }

  function clickSquare(sq: string) {
    if (!solverTurn) return
    if (!selected) {
      if (pieces[sq] && isSolver(pieces[sq]!)) setSelected(sq)
      return
    }
    if (sq === selected) {
      setSelected(null)
      return
    }
    if (pieces[sq] && isSolver(pieces[sq]!)) {
      setSelected(sq) // reselect own piece
      return
    }
    attemptMove(selected, sq)
  }

  function dropOn(target: string, from: string) {
    if (!solverTurn || !from || from === target) return
    if (pieces[target] && isSolver(pieces[target]!)) {
      setSelected(null) // dropped back onto an own piece → cancel
      return
    }
    attemptMove(from, target)
  }

  const ringClass = finishedByUser
    ? 'ring-2 ring-green-500'
    : wrong
      ? 'ring-2 ring-red-500'
      : 'ring-1 ring-ink-900/10'

  return (
    <div>
      <div
        className={`relative mx-auto grid aspect-square w-full max-w-[26rem] select-none grid-cols-8 grid-rows-8 overflow-hidden rounded-lg shadow-md ${ringClass}`}
        role="group"
        aria-label="Échiquier du puzzle"
      >
        {ranks.map((rank, r) =>
          files.map((file, f) => {
            const sq = file + rank
            const piece = pieces[sq]
            const dark = (file.charCodeAt(0) - 97 + rank) % 2 === 1
            const isSel = selected === sq
            const isLast = sq === lastFrom || sq === lastTo
            const light = '#f0d9b5'
            const darkSq = '#b58863'
            const base = dark ? darkSq : light
            const coordColor = dark ? light : darkSq
            return (
              <button
                key={sq}
                type="button"
                onClick={() => clickSquare(sq)}
                disabled={!solverTurn}
                onDragOver={(e) => {
                  if (solverTurn) e.preventDefault() // allow drop
                }}
                onDrop={(e) => {
                  e.preventDefault()
                  dropOn(sq, e.dataTransfer.getData('text/plain'))
                }}
                aria-label={sq + (piece ? ` ${piece}` : '')}
                className={`relative flex items-center justify-center ${
                  solverTurn && piece && isSolver(piece) ? 'cursor-grab' : solverTurn ? 'cursor-pointer' : 'cursor-default'
                }`}
                style={{ backgroundColor: base }}
              >
                {/* last-move tint */}
                {isLast ? <span aria-hidden className="absolute inset-0" style={{ backgroundColor: 'rgba(155,199,0,0.41)' }} /> : null}
                {/* selection ring */}
                {isSel ? <span aria-hidden className="absolute inset-0 outline outline-[3px] -outline-offset-[3px] outline-gold-500" /> : null}
                {/* coordinates (files on bottom row, ranks on left column) */}
                {f === 0 ? (
                  <span aria-hidden className="absolute left-[3px] top-[2px] text-[9px] font-bold leading-none sm:text-[11px]" style={{ color: coordColor }}>
                    {rank}
                  </span>
                ) : null}
                {r === ranks.length - 1 ? (
                  <span aria-hidden className="absolute bottom-[1px] right-[3px] text-[9px] font-bold leading-none sm:text-[11px]" style={{ color: coordColor }}>
                    {file}
                  </span>
                ) : null}
                {piece ? (
                  <span
                    draggable={solverTurn && isSolver(piece)}
                    onDragStart={(e) => {
                      if (!(solverTurn && isSolver(piece))) return
                      setSelected(sq)
                      e.dataTransfer.setData('text/plain', sq)
                      e.dataTransfer.effectAllowed = 'move'
                    }}
                    className="relative select-none text-[9vw] leading-none sm:text-[2.25rem]"
                    style={{
                      color: piece === piece.toUpperCase() ? '#f8f8f8' : '#3a3a38',
                      textShadow:
                        piece === piece.toUpperCase()
                          ? '0 0 1px #000, 0 0 1px #000, 0 1px 2px rgba(0,0,0,.4)'
                          : '0 0 1px rgba(255,255,255,.35), 0 1px 1px rgba(0,0,0,.25)',
                    }}
                  >
                    {GLYPH[piece.toLowerCase()]}
                  </span>
                ) : null}
              </button>
            )
          }),
        )}

        {/* Clear "finished" overlay badge. */}
        {solved ? (
          <div aria-hidden className="pointer-events-none absolute inset-0 flex items-start justify-end p-2">
            <span
              className={`rounded-full px-3 py-1 text-sm font-bold shadow ${
                finishedByUser ? 'bg-green-600 text-white' : 'bg-ink-900/90 text-gold-300'
              }`}
            >
              {finishedByUser ? `✓ ${labels.solved}` : labels.solutionShown}
            </span>
          </div>
        ) : null}
      </div>

      <div className="mt-4 flex flex-wrap items-center justify-between gap-3">
        <p
          className={`text-sm font-semibold ${finishedByUser ? 'text-green-700' : wrong ? 'text-red-600' : 'text-ink-700'}`}
          role="status"
        >
          {finishedByUser
            ? `✓ ${labels.solved}`
            : solved
              ? labels.solutionShown
              : wrong
                ? labels.tryAgain
                : `${sideToMove === 'w' ? labels.whiteToPlay : labels.blackToPlay} — ${labels.yourMove}`}
        </p>
        <div className="flex items-center gap-2">
          {!solved && !revealed ? (
            <button
              type="button"
              onClick={revealSolution}
              className="rounded-md border border-ink-200 px-3 py-1 text-xs font-medium text-ink-600 transition-colors hover:border-gold-400 hover:text-ink-900"
            >
              {labels.showSolution}
            </button>
          ) : null}
          <button
            type="button"
            onClick={reset}
            className="rounded-md border border-ink-200 px-3 py-1 text-xs font-medium text-ink-600 transition-colors hover:border-gold-400 hover:text-ink-900"
          >
            {labels.retry}
          </button>
        </div>
      </div>
    </div>
  )
}
