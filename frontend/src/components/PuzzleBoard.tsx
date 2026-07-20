import { useMemo, useState } from 'react'

/**
 * Lightweight interactive chess puzzle. No client-side chess engine: the
 * solution line is known (trusted, from Stockfish), so we only compare the
 * player's move to the expected one and animate the position. Fully anonymised
 * data (mirrored FEN + solution) — see backend/internal/tactics.
 */
export interface PuzzleBoardProps {
  fen: string
  sideToMove: 'w' | 'b' // solver's colour
  solution: string[] // UCI moves, solver plays even indices
  onSolved?: () => void
  onAttempt?: (correct: boolean) => void
  labels: { yourMove: string; solved: string; tryAgain: string; retry: string; whiteToPlay: string; blackToPlay: string }
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

export function PuzzleBoard({ fen, sideToMove, solution, onSolved, onAttempt, labels }: PuzzleBoardProps) {
  const [applied, setApplied] = useState(0)
  const [selected, setSelected] = useState<string | null>(null)
  const [wrong, setWrong] = useState(false)

  const pieces = useMemo(() => applyMoves(fen, solution.slice(0, applied)), [fen, solution, applied])
  const solved = applied >= solution.length
  const solverTurn = applied % 2 === 0 && !solved

  const ranks = sideToMove === 'w' ? [8, 7, 6, 5, 4, 3, 2, 1] : [1, 2, 3, 4, 5, 6, 7, 8]
  const files = sideToMove === 'w' ? FILES : [...FILES].reverse()
  const isSolver = (piece: string) => (sideToMove === 'w' ? piece === piece.toUpperCase() : piece === piece.toLowerCase())

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
    const expected = solution[applied]!
    const candidate = selected + sq
    if (expected.slice(0, 4) === candidate) {
      onAttempt?.(true)
      setSelected(null)
      setWrong(false)
      const next = applied + 1
      setApplied(next)
      if (next >= solution.length) onSolved?.()
      else setTimeout(() => setApplied(next + 1), 350) // opponent's forced reply
    } else if (pieces[sq] && isSolver(pieces[sq]!)) {
      setSelected(sq) // reselect own piece
    } else {
      onAttempt?.(false)
      setWrong(true)
      setSelected(null)
      setTimeout(() => setWrong(false), 600)
    }
  }

  return (
    <div>
      <div
        className={`mx-auto grid aspect-square w-full max-w-[26rem] grid-cols-8 overflow-hidden rounded-xl ring-1 ring-ink-900/10 ${
          wrong ? 'ring-2 ring-red-500' : ''
        }`}
        role="group"
        aria-label="Échiquier du puzzle"
      >
        {ranks.map((rank) =>
          files.map((file) => {
            const sq = file + rank
            const piece = pieces[sq]
            const dark = (file.charCodeAt(0) - 97 + rank) % 2 === 0
            const isSel = selected === sq
            return (
              <button
                key={sq}
                type="button"
                onClick={() => clickSquare(sq)}
                disabled={!solverTurn}
                aria-label={sq + (piece ? ` ${piece}` : '')}
                className={`relative flex items-center justify-center ${dark ? 'bg-[#b58863]' : 'bg-[#f0d9b5]'} ${
                  isSel ? 'outline outline-4 -outline-offset-4 outline-gold-500' : ''
                } ${solverTurn ? 'cursor-pointer' : 'cursor-default'}`}
              >
                {piece ? (
                  <span
                    className="select-none text-[7vw] leading-none sm:text-[2rem]"
                    style={{
                      color: piece === piece.toUpperCase() ? '#fafafa' : '#1a1a1a',
                      textShadow: piece === piece.toUpperCase() ? '0 0 1px #000, 0 1px 2px rgba(0,0,0,.35)' : '0 1px 1px rgba(255,255,255,.25)',
                    }}
                  >
                    {GLYPH[piece.toLowerCase()]}
                  </span>
                ) : null}
              </button>
            )
          }),
        )}
      </div>

      <div className="mt-4 flex items-center justify-between gap-3">
        <p className={`text-sm font-semibold ${wrong ? 'text-red-600' : solved ? 'text-green-700' : 'text-ink-700'}`} role="status">
          {solved ? `✓ ${labels.solved}` : wrong ? labels.tryAgain : `${sideToMove === 'w' ? labels.whiteToPlay : labels.blackToPlay} — ${labels.yourMove}`}
        </p>
        <button
          type="button"
          onClick={() => {
            setApplied(0)
            setSelected(null)
            setWrong(false)
          }}
          className="rounded-md border border-ink-200 px-3 py-1 text-xs font-medium text-ink-600 transition-colors hover:border-gold-400 hover:text-ink-900"
        >
          {labels.retry}
        </button>
      </div>
    </div>
  )
}
