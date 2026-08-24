/**
 * Minimal chess rules engine — just enough to replay a PGN main line and know
 * which square each move came from. No engine, no evaluation: legal-move
 * generation exists only to resolve SAN (`Rfd8`, `exd5`, `O-O`) into a
 * from → to pair, which is what the board needs in order to highlight it.
 *
 * Squares are indexed 0 (a1) … 63 (h8): index = file + 8 × rank.
 * Pieces follow FEN casing — uppercase is White.
 */

const FILES = 'abcdefgh'

export type Square = number
export type Board = (string | null)[]

export interface Position {
  board: Board
  turn: 'w' | 'b'
  castling: string // subset of "KQkq"
  ep: Square | null // en-passant target square
}

export interface Move {
  from: Square
  to: Square
  promotion?: string // uppercase piece letter
}

export const START_FEN = 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1'

export const squareName = (sq: Square): string => FILES[sq & 7]! + (1 + (sq >> 3))
export const squareIndex = (name: string): Square =>
  FILES.indexOf(name[0]!) + (Number(name[1]) - 1) * 8

const isWhite = (piece: string): boolean => piece === piece.toUpperCase()
const fileOf = (sq: Square): number => sq & 7
const rankOf = (sq: Square): number => sq >> 3

export function parseFen(fen: string): Position {
  const [placement, turn, castling, ep] = fen.trim().split(/\s+/)
  const board: Board = new Array(64).fill(null)
  placement!.split('/').forEach((row, i) => {
    const rank = 7 - i
    let file = 0
    for (const ch of row) {
      if (ch >= '1' && ch <= '8') file += Number(ch)
      else board[rank * 8 + file++] = ch
    }
  })
  return {
    board,
    turn: turn === 'b' ? 'b' : 'w',
    castling: !castling || castling === '-' ? '' : castling,
    ep: !ep || ep === '-' ? null : squareIndex(ep),
  }
}

// --- Move generation --------------------------------------------------------

const KNIGHT_STEPS = [
  [1, 2], [2, 1], [2, -1], [1, -2], [-1, -2], [-2, -1], [-2, 1], [-1, 2],
] as const
const BISHOP_RAYS = [[1, 1], [1, -1], [-1, 1], [-1, -1]] as const
const ROOK_RAYS = [[1, 0], [-1, 0], [0, 1], [0, -1]] as const
const KING_STEPS = [...BISHOP_RAYS, ...ROOK_RAYS] as const

/** Squares reachable from `from` along `dirs` (one step only when not sliding). */
function* rays(from: Square, dirs: readonly (readonly [number, number])[], board: Board, sliding: boolean) {
  for (const [df, dr] of dirs) {
    let f = fileOf(from) + df
    let r = rankOf(from) + dr
    while (f >= 0 && f < 8 && r >= 0 && r < 8) {
      const sq = r * 8 + f
      yield sq
      if (!sliding || board[sq]) break
      f += df
      r += dr
    }
  }
}

/** Is `sq` attacked by the side `byWhite`? (pawn pushes do not attack) */
function isAttacked(board: Board, sq: Square, byWhite: boolean): boolean {
  const f = fileOf(sq)
  const r = rankOf(sq)

  // Pawns capture diagonally forward, so look backwards from `sq`.
  const pawnRank = r + (byWhite ? -1 : 1)
  if (pawnRank >= 0 && pawnRank < 8) {
    for (const df of [-1, 1]) {
      const pf = f + df
      if (pf < 0 || pf > 7) continue
      if (board[pawnRank * 8 + pf] === (byWhite ? 'P' : 'p')) return true
    }
  }
  for (const [df, dr] of KNIGHT_STEPS) {
    const nf = f + df
    const nr = r + dr
    if (nf < 0 || nf > 7 || nr < 0 || nr > 7) continue
    const p = board[nr * 8 + nf]
    if (p && p.toLowerCase() === 'n' && isWhite(p) === byWhite) return true
  }
  for (const [df, dr] of KING_STEPS) {
    const nf = f + df
    const nr = r + dr
    if (nf < 0 || nf > 7 || nr < 0 || nr > 7) continue
    const p = board[nr * 8 + nf]
    if (p && p.toLowerCase() === 'k' && isWhite(p) === byWhite) return true
  }
  for (const [dirs, types] of [[BISHOP_RAYS, 'bq'], [ROOK_RAYS, 'rq']] as const) {
    for (const target of rays(sq, dirs, board, true)) {
      const p = board[target]
      if (!p) continue
      if (isWhite(p) === byWhite && types.includes(p.toLowerCase())) return true
      break // the ray is blocked either way
    }
  }
  return false
}

const kingSquare = (board: Board, white: boolean): Square => board.indexOf(white ? 'K' : 'k')

/** Pseudo-legal moves for the side to move — king safety is filtered later. */
function pseudoMoves(pos: Position): Move[] {
  const { board, turn } = pos
  const white = turn === 'w'
  const moves: Move[] = []
  const push = (from: Square, to: Square) => {
    const target = board[to]
    if (target && isWhite(target) === white) return
    moves.push({ from, to })
  }

  for (let from = 0; from < 64; from++) {
    const piece = board[from]
    if (!piece || isWhite(piece) !== white) continue
    const type = piece.toLowerCase()

    if (type === 'p') {
      const dir = white ? 1 : -1
      const startRank = white ? 1 : 6
      const lastRank = white ? 7 : 0
      const f = fileOf(from)
      const r = rankOf(from)
      if (r + dir >= 0 && r + dir < 8) {
        const one = (r + dir) * 8 + f
        if (!board[one]) {
          if (rankOf(one) === lastRank) for (const p of 'QRBN') moves.push({ from, to: one, promotion: p })
          else moves.push({ from, to: one })
          const two = (r + 2 * dir) * 8 + f
          if (r === startRank && !board[two]) moves.push({ from, to: two })
        }
        for (const df of [-1, 1]) {
          const cf = f + df
          if (cf < 0 || cf > 7) continue
          const to = (r + dir) * 8 + cf
          const target = board[to]
          if (!((target && isWhite(target) !== white) || to === pos.ep)) continue
          if (rankOf(to) === lastRank) for (const p of 'QRBN') moves.push({ from, to, promotion: p })
          else moves.push({ from, to })
        }
      }
    } else if (type === 'n') {
      for (const [df, dr] of KNIGHT_STEPS) {
        const nf = fileOf(from) + df
        const nr = rankOf(from) + dr
        if (nf >= 0 && nf < 8 && nr >= 0 && nr < 8) push(from, nr * 8 + nf)
      }
    } else if (type === 'b') {
      for (const sq of rays(from, BISHOP_RAYS, board, true)) push(from, sq)
    } else if (type === 'r') {
      for (const sq of rays(from, ROOK_RAYS, board, true)) push(from, sq)
    } else if (type === 'q') {
      for (const sq of rays(from, KING_STEPS, board, true)) push(from, sq)
    } else if (type === 'k') {
      for (const sq of rays(from, KING_STEPS, board, false)) push(from, sq)
      // Castling: rights present, path clear, king never crossing an attacked square.
      const rank = white ? 0 : 7
      const rights = white ? 'KQ' : 'kq'
      if (from === rank * 8 + 4 && !isAttacked(board, from, !white)) {
        if (pos.castling.includes(rights[0]!) && !board[from + 1] && !board[from + 2] && !isAttacked(board, from + 1, !white)) {
          moves.push({ from, to: from + 2 })
        }
        if (
          pos.castling.includes(rights[1]!) &&
          !board[from - 1] && !board[from - 2] && !board[from - 3] &&
          !isAttacked(board, from - 1, !white)
        ) {
          moves.push({ from, to: from - 2 })
        }
      }
    }
  }
  return moves
}

export function applyMove(pos: Position, move: Move): Position {
  const board = pos.board.slice()
  const piece = board[move.from]!
  const white = isWhite(piece)
  const type = piece.toLowerCase()

  board[move.from] = null

  // En passant: a pawn stepping diagonally onto an empty square eats behind it.
  if (type === 'p' && fileOf(move.from) !== fileOf(move.to) && !pos.board[move.to]) {
    board[rankOf(move.from) * 8 + fileOf(move.to)] = null
  }
  // Castling drags the rook along.
  if (type === 'k' && Math.abs(fileOf(move.from) - fileOf(move.to)) === 2) {
    const rank = rankOf(move.from) * 8
    if (fileOf(move.to) === 6) {
      board[rank + 5] = board[rank + 7]!
      board[rank + 7] = null
    } else {
      board[rank + 3] = board[rank]!
      board[rank] = null
    }
  }
  board[move.to] = move.promotion ? (white ? move.promotion : move.promotion.toLowerCase()) : piece

  let castling = pos.castling
  if (type === 'k') castling = castling.replace(white ? /[KQ]/g : /[kq]/g, '')
  // A rook that leaves — or gets taken on — a corner kills that right.
  const CORNER: Record<number, string> = { 0: 'Q', 7: 'K', 56: 'q', 63: 'k' }
  for (const sq of [move.from, move.to]) {
    const right = CORNER[sq]
    if (right) castling = castling.replace(right, '')
  }

  const doublePush = type === 'p' && Math.abs(rankOf(move.to) - rankOf(move.from)) === 2
  return {
    board,
    turn: white ? 'b' : 'w',
    castling,
    ep: doublePush ? (rankOf(move.from) + (white ? 1 : -1)) * 8 + fileOf(move.from) : null,
  }
}

export function legalMoves(pos: Position): Move[] {
  const white = pos.turn === 'w'
  return pseudoMoves(pos).filter((m) => {
    const next = applyMove(pos, m)
    return !isAttacked(next.board, kingSquare(next.board, white), !white)
  })
}

// --- SAN ---------------------------------------------------------------------

const SAN_RE = /^([KQRBN])?([a-h])?([1-8])?x?([a-h][1-8])(?:=?([QRBN]))?$/

/**
 * Resolve one SAN token against a position. Returns null when the move does not
 * exist there — a typo in the PGN, or a score the engine cannot follow.
 */
export function sanToMove(pos: Position, san: string): Move | null {
  const clean = san.replace(/[+#!?]+$/, '')
  const legal = legalMoves(pos)

  if (/^(O-O|0-0)(-O|-0)?$/.test(clean)) {
    const rank = pos.turn === 'w' ? 0 : 7
    const target = rank * 8 + (clean.length > 3 ? 2 : 6)
    return legal.find((m) => pos.board[m.from]?.toLowerCase() === 'k' && m.to === target) ?? null
  }

  const parsed = SAN_RE.exec(clean)
  if (!parsed) return null
  const [, pieceLetter, fromFile, fromRank, dest, promotion] = parsed
  const type = (pieceLetter ?? 'P').toLowerCase()
  const to = squareIndex(dest!)

  const candidates = legal.filter((mv) => {
    if (mv.to !== to) return false
    if (pos.board[mv.from]!.toLowerCase() !== type) return false
    if (fromFile && FILES[fileOf(mv.from)] !== fromFile) return false
    if (fromRank && rankOf(mv.from) + 1 !== Number(fromRank)) return false
    return promotion ? mv.promotion === promotion : !mv.promotion
  })
  return candidates[0] ?? null
}

// --- Notation ----------------------------------------------------------------

const FR_PIECES: Record<string, string> = { K: 'R', Q: 'D', R: 'T', B: 'F', N: 'C' }

/** English SAN → French notation (Cf3, Fg2, Td8…). File letters stay lowercase. */
export function toFrenchSan(san: string): string {
  return san
    .replace(/^([KQRBN])/, (_, p: string) => FR_PIECES[p]!)
    .replace(/=([QRBN])/, (_, p: string) => '=' + FR_PIECES[p]!)
}
