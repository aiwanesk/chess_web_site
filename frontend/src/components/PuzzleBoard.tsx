import { useEffect, useMemo, useRef, useState, type PointerEvent as ReactPointerEvent } from 'react'

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

// Official Lichess cburnett piece SVG inner content (viewBox 0 0 45 45)
// Source: https://github.com/lichess-org/lila/tree/master/public/piece/cburnett
const PIECE_SVG: Record<string, string> = {
  K: `<g fill="none" fill-rule="evenodd" stroke="#000" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><path stroke-linejoin="miter" d="M22.5 11.63V6M20 8h5"/><path fill="#fff" stroke-linecap="butt" stroke-linejoin="miter" d="M22.5 25s4.5-7.5 3-10.5c0 0-1-2.5-3-2.5s-3 2.5-3 2.5c-1.5 3 3 10.5 3 10.5"/><path fill="#fff" d="M11.5 37c5.5 3.5 15.5 3.5 21 0v-7s9-4.5 6-10.5c-4-6.5-13.5-3.5-16 4V27v-3.5c-3.5-7.5-13-10.5-16-4-3 6 5 10 5 10z"/><path d="M11.5 30c5.5-3 15.5-3 21 0m-21 3.5c5.5-3 15.5-3 21 0m-21 3.5c5.5-3 15.5-3 21 0"/></g>`,
  Q: `<g fill="#fff" fill-rule="evenodd" stroke="#000" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><path d="M8 12a2 2 0 1 1-4 0 2 2 0 1 1 4 0m16.5-4.5a2 2 0 1 1-4 0 2 2 0 1 1 4 0M41 12a2 2 0 1 1-4 0 2 2 0 1 1 4 0M16 8.5a2 2 0 1 1-4 0 2 2 0 1 1 4 0M33 9a2 2 0 1 1-4 0 2 2 0 1 1 4 0"/><path stroke-linecap="butt" d="M9 26c8.5-1.5 21-1.5 27 0l2-12-7 11V11l-5.5 13.5-3-15-3 15-5.5-14V25L7 14z"/><path stroke-linecap="butt" d="M9 26c0 2 1.5 2 2.5 4 1 1.5 1 1 .5 3.5-1.5 1-1.5 2.5-1.5 2.5-1.5 1.5.5 2.5.5 2.5 6.5 1 16.5 1 23 0 0 0 1.5-1 0-2.5 0 0 .5-1.5-1-2.5-.5-2.5-.5-2 .5-3.5 1-2 2.5-2 2.5-4-8.5-1.5-18.5-1.5-27 0z"/><path fill="none" d="M11.5 30c3.5-1 18.5-1 22 0M12 33.5c6-1 15-1 21 0"/></g>`,
  R: `<g fill="#fff" fill-rule="evenodd" stroke="#000" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><path stroke-linecap="butt" d="M9 39h27v-3H9zm3-3v-4h21v4zm-1-22V9h4v2h5V9h5v2h5V9h4v5"/><path d="m34 14-3 3H14l-3-3"/><path stroke-linecap="butt" stroke-linejoin="miter" d="M31 17v12.5H14V17"/><path d="m31 29.5 1.5 2.5h-20l1.5-2.5"/><path fill="none" stroke-linejoin="miter" d="M11 14h23"/></g>`,
  B: `<g fill="none" fill-rule="evenodd" stroke="#000" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><g fill="#fff" stroke-linecap="butt"><path d="M9 36c3.39-.97 10.11.43 13.5-2 3.39 2.43 10.11 1.03 13.5 2 0 0 1.65.54 3 2-.68.97-1.65.99-3 .5-3.39-.97-10.11.46-13.5-1-3.39 1.46-10.11.03-13.5 1-1.35.49-2.32.47-3-.5 1.35-1.94 3-2 3-2z"/><path d="M15 32c2.5 2.5 12.5 2.5 15 0 .5-1.5 0-2 0-2 0-2.5-2.5-4-2.5-4 5.5-1.5 6-11.5-5-15.5-11 4-10.5 14-5 15.5 0 0-2.5 1.5-2.5 4 0 0-.5.5 0 2z"/><path d="M25 8a2.5 2.5 0 1 1-5 0 2.5 2.5 0 1 1 5 0z"/></g><path stroke-linejoin="miter" d="M17.5 26h10M15 30h15m-7.5-14.5v5M20 18h5"/></g>`,
  N: `<g fill="none" fill-rule="evenodd" stroke="#000" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><path fill="#fff" d="M22 10c10.5 1 16.5 8 16 29H15c0-9 10-6.5 8-21"/><path fill="#fff" d="M24 18c.38 2.91-5.55 7.37-8 9-3 2-2.82 4.34-5 4-1.042-.94 1.41-3.04 0-3-1 0 .19 1.23-1 2-1 0-4.003 1-4-4 0-2 6-12 6-12s1.89-1.9 2-3.5c-.73-.994-.5-2-.5-3 1-1 3 2.5 3 2.5h2s.78-1.992 2.5-3c1 0 1 3 1 3"/><path fill="#000" d="M9.5 25.5a.5.5 0 1 1-1 0 .5.5 0 1 1 1 0m5.433-9.75a.5 1.5 30 1 1-.866-.5.5 1.5 30 1 1 .866.5"/></g>`,
  P: `<path fill="#fff" stroke="#000" stroke-linecap="round" stroke-width="1.5" d="M22.5 9c-2.21 0-4 1.79-4 4 0 .89.29 1.71.78 2.38C17.33 16.5 16 18.59 16 21c0 2.03.94 3.84 2.41 5.03-3 1.06-7.41 5.55-7.41 13.47h23c0-7.92-4.41-12.41-7.41-13.47 1.47-1.19 2.41-3 2.41-5.03 0-2.41-1.33-4.5-3.28-5.62.49-.67.78-1.49.78-2.38 0-2.21-1.79-4-4-4z"/>`,
  k: `<g fill="none" fill-rule="evenodd" stroke="#000" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><path stroke-linejoin="miter" d="M22.5 11.6V6"/><path fill="#000" stroke-linecap="butt" stroke-linejoin="miter" d="M22.5 25s4.5-7.5 3-10.5c0 0-1-2.5-3-2.5s-3 2.5-3 2.5c-1.5 3 3 10.5 3 10.5"/><path fill="#000" d="M11.5 37a22.3 22.3 0 0 0 21 0v-7s9-4.5 6-10.5c-4-6.5-13.5-3.5-16 4V27v-3.5c-3.5-7.5-13-10.5-16-4-3 6 5 10 5 10z"/><path stroke-linejoin="miter" d="M20 8h5"/><path stroke="#ececec" d="M32 29.5s8.5-4 6-9.7C34.1 14 25 18 22.5 24.6v2.1-2.1C20 18 9.9 14 7 19.9c-2.5 5.6 4.8 9 4.8 9"/><path stroke="#ececec" d="M11.5 30c5.5-3 15.5-3 21 0m-21 3.5c5.5-3 15.5-3 21 0m-21 3.5c5.5-3 15.5-3 21 0"/></g>`,
  q: `<g fill-rule="evenodd" stroke="#000" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><g stroke="none"><circle cx="6" cy="12" r="2.75"/><circle cx="14" cy="9" r="2.75"/><circle cx="22.5" cy="8" r="2.75"/><circle cx="31" cy="9" r="2.75"/><circle cx="39" cy="12" r="2.75"/></g><path stroke-linecap="butt" d="M9 26c8.5-1.5 21-1.5 27 0l2.5-12.5L31 25l-.3-14.1-5.2 13.6-3-14.5-3 14.5-5.2-13.6L14 25 6.5 13.5z"/><path stroke-linecap="butt" d="M9 26c0 2 1.5 2 2.5 4 1 1.5 1 1 .5 3.5-1.5 1-1.5 2.5-1.5 2.5-1.5 1.5.5 2.5.5 2.5 6.5 1 16.5 1 23 0 0 0 1.5-1 0-2.5 0 0 .5-1.5-1-2.5-.5-2.5-.5-2 .5-3.5 1-2 2.5-2 2.5-4-8.5-1.5-18.5-1.5-27 0z"/><path fill="none" stroke-linecap="butt" d="M11 38.5a35 35 1 0 0 23 0"/><path fill="none" stroke="#ececec" d="M11 29a35 35 1 0 1 23 0m-21.5 2.5h20m-21 3a35 35 1 0 0 22 0m-23 3a35 35 1 0 0 24 0"/></g>`,
  r: `<g fill-rule="evenodd" stroke="#000" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><path stroke-linecap="butt" d="M9 39h27v-3H9zm3.5-7 1.5-2.5h17l1.5 2.5zm-.5 4v-4h21v4z"/><path stroke-linecap="butt" stroke-linejoin="miter" d="M14 29.5v-13h17v13z"/><path stroke-linecap="butt" d="M14 16.5 11 14h23l-3 2.5zM11 14V9h4v2h5V9h5v2h5V9h4v5z"/><path fill="none" stroke="#ececec" stroke-linejoin="miter" stroke-width="1" d="M12 35.5h21m-20-4h19m-18-2h17m-17-13h17M11 14h23"/></g>`,
  b: `<g fill="none" fill-rule="evenodd" stroke="#000" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><g fill="#000" stroke-linecap="butt"><path d="M9 36c3.4-1 10.1.4 13.5-2 3.4 2.4 10.1 1 13.5 2 0 0 1.6.5 3 2-.7 1-1.6 1-3 .5-3.4-1-10.1.5-13.5-1-3.4 1.5-10.1 0-13.5 1-1.4.5-2.3.5-3-.5 1.4-2 3-2 3-2z"/><path d="M15 32c2.5 2.5 12.5 2.5 15 0 .5-1.5 0-2 0-2 0-2.5-2.5-4-2.5-4 5.5-1.5 6-11.5-5-15.5-11 4-10.5 14-5 15.5 0 0-2.5 1.5-2.5 4 0 0-.5.5 0 2z"/><path d="M25 8a2.5 2.5 0 1 1-5 0 2.5 2.5 0 1 1 5 0z"/></g><path stroke="#ececec" stroke-linejoin="miter" d="M17.5 26h10M15 30h15m-7.5-14.5v5M20 18h5"/></g>`,
  n: `<g fill="none" fill-rule="evenodd" stroke="#000" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><path fill="#000" d="M22 10c10.5 1 16.5 8 16 29H15c0-9 10-6.5 8-21"/><path fill="#000" d="M24 18c.38 2.91-5.55 7.37-8 9-3 2-2.82 4.34-5 4-1.04-.94 1.41-3.04 0-3-1 0 .19 1.23-1 2-1 0-4 1-4-4 0-2 6-12 6-12s1.89-1.9 2-3.5c-.73-1-.5-2-.5-3 1-1 3 2.5 3 2.5h2s.78-2 2.5-3c1 0 1 3 1 3"/><path fill="#ececec" stroke="#ececec" d="M9.5 25.5a.5.5 0 1 1-1 0 .5.5 0 1 1 1 0m5.43-9.75a.5 1.5 30 1 1-.86-.5.5 1.5 30 1 1 .86.5"/><path fill="#ececec" stroke="none" d="m24.55 10.4-.45 1.45.5.15c3.15 1 5.65 2.49 7.9 6.75S35.75 29.06 35.25 39l-.05.5h2.25l.05-.5c.5-10.06-.88-16.85-3.25-21.34s-5.79-6.64-9.19-7.16z"/></g>`,
  p: `<path stroke="#000" stroke-linecap="round" stroke-width="1.5" d="M22.5 9a4 4 0 0 0-3.22 6.38 6.48 6.48 0 0 0-.87 10.65c-3 1.06-7.41 5.55-7.41 13.47h23c0-7.92-4.41-12.41-7.41-13.47a6.46 6.46 0 0 0-.87-10.65A4.01 4.01 0 0 0 22.5 9z"/>`,
}

function PieceSvg({ piece, className, style }: { piece: string; className?: string; style?: React.CSSProperties }) {
  const svg = PIECE_SVG[piece]
  if (!svg) return null
  return (
    <svg
      viewBox="0 0 45 45"
      className={className}
      style={style}
      aria-hidden
      dangerouslySetInnerHTML={{ __html: svg }}
    />
  )
}

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

  // Pointer-based drag (lichess-style): a piece that follows the cursor, not the
  // native HTML5 drag ghost (which drags the whole coloured square).
  const boardRef = useRef<HTMLDivElement>(null)
  const [drag, setDrag] = useState<{ from: string; piece: string; x: number; y: number; size: number } | null>(null)
  const movedRef = useRef(false)
  const wasSelectedRef = useRef(false)

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

  // Which square is under a client coordinate (given the current orientation)?
  function squareAt(clientX: number, clientY: number): string | null {
    const el = boardRef.current
    if (!el) return null
    const r = el.getBoundingClientRect()
    const col = Math.floor(((clientX - r.left) / r.width) * 8)
    const row = Math.floor(((clientY - r.top) / r.height) * 8)
    if (col < 0 || col > 7 || row < 0 || row > 7) return null
    return files[col]! + ranks[row]!
  }

  function onSquarePointerDown(sq: string, e: ReactPointerEvent) {
    if (!solverTurn) return
    if (pieces[sq] && isSolver(pieces[sq]!)) {
      e.preventDefault()
      movedRef.current = false
      wasSelectedRef.current = selected === sq
      setSelected(sq)
      const r = boardRef.current?.getBoundingClientRect()
      setDrag({
        from: sq,
        piece: pieces[sq]!,
        x: e.clientX,
        y: e.clientY,
        size: r ? r.width / 8 : 40,
      })
      boardRef.current?.setPointerCapture(e.pointerId)
    }
    // taps on empty/enemy squares are resolved on pointer-up (click-to-move)
  }

  function onBoardPointerMove(e: ReactPointerEvent) {
    if (!drag) return
    movedRef.current = true
    setDrag((d) => (d ? { ...d, x: e.clientX, y: e.clientY } : d))
  }

  function onBoardPointerUp(e: ReactPointerEvent) {
    const target = squareAt(e.clientX, e.clientY)
    if (drag) {
      boardRef.current?.releasePointerCapture(e.pointerId)
      const from = drag.from
      setDrag(null)
      if (movedRef.current && target && target !== from) {
        dropOn(target, from) // a real drag → play it
      } else if (!movedRef.current && wasSelectedRef.current) {
        setSelected(null) // tapped the already-selected piece → deselect
      }
      return
    }
    // pointer started on an empty/enemy square → click-to-move to it
    if (selected && target && target !== selected && !(pieces[target] && isSolver(pieces[target]!))) {
      attemptMove(selected, target)
    }
  }

  const ringClass = finishedByUser
    ? 'ring-2 ring-green-500'
    : wrong
      ? 'ring-2 ring-red-500'
      : 'ring-1 ring-ink-900/10'

  return (
    <div>
      {/* Floating piece that follows the cursor while dragging (position:fixed so
          it escapes the board's overflow-hidden). */}
      {drag ? (
        <PieceSvg
          piece={drag.piece}
          className="pointer-events-none fixed z-50 select-none"
          style={{
            left: drag.x,
            top: drag.y,
            transform: 'translate(-50%, -50%)',
            width: drag.size * 0.9,
            height: drag.size * 0.9,
          }}
        />
      ) : null}

      <div
        ref={boardRef}
        onPointerMove={onBoardPointerMove}
        onPointerUp={onBoardPointerUp}
        className={`relative mx-auto grid aspect-square w-full max-w-[26rem] select-none touch-none grid-cols-8 grid-rows-8 overflow-hidden rounded-lg shadow-md ${ringClass}`}
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
                onPointerDown={(e) => onSquarePointerDown(sq, e)}
                onClick={(e) => {
                  if (e.detail === 0) clickSquare(sq) // keyboard activation only (mouse/touch go through pointer events)
                }}
                disabled={!solverTurn}
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
                  <PieceSvg
                    piece={piece}
                    className="pointer-events-none relative w-[90%] h-[90%] select-none"
                    style={{
                      opacity: drag && drag.from === sq ? 0.25 : 1,
                    }}
                  />
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
