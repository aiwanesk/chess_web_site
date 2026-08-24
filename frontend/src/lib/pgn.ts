/**
 * PGN → replayable game. Scores live as plain `.pgn` files in /content/games so
 * a new game is one paste away: drop the file, reference it from an article with
 * `[[pgn:file-name]]`.
 *
 * Only the MAIN LINE is kept. Sidelines, comments and NAGs are stripped: the
 * viewer under an article is there to follow the story, not to browse a tree.
 */
import { applyMove, parseFen, sanToMove, START_FEN, type Board, type Square } from './chess'

const files = import.meta.glob('../../../content/games/*.pgn', {
  query: '?raw',
  import: 'default',
  eager: true,
}) as Record<string, string>

export interface Ply {
  san: string // as written in the PGN (English notation)
  from: Square
  to: Square
  board: Board // position AFTER the move
  number: number // move number (1, 1, 2, 2, …)
  color: 'w' | 'b'
}

export interface Game {
  id: string
  headers: Record<string, string>
  white: string
  black: string
  whiteElo?: string
  blackElo?: string
  event?: string
  round?: string
  date?: string
  result: string
  /** Side shown at the bottom on first render. */
  orientation: 'w' | 'b'
  startBoard: Board
  plies: Ply[]
}

const HEADER_RE = /\[\s*(\w+)\s+"([^"]*)"\s*\]/g
const RESULTS = new Set(['1-0', '0-1', '1/2-1/2', '*'])

/** Drop comments, then sidelines (they can nest), then NAGs and move numbers. */
function mainLineTokens(movetext: string): string[] {
  let text = movetext.replace(/\{[^}]*\}/g, ' ').replace(/;[^\n]*/g, ' ')

  let out = ''
  let depth = 0
  for (const ch of text) {
    if (ch === '(') depth++
    else if (ch === ')') depth = Math.max(0, depth - 1)
    else if (depth === 0) out += ch
  }
  text = out.replace(/\$\d+/g, ' ')

  return text
    .split(/\s+/)
    .map((tok) => tok.replace(/^\d+\.(\.\.)?/, '')) // "12.e4" and "12...e4" glued forms
    .filter((tok) => tok && !RESULTS.has(tok) && !/^\d+\.*$/.test(tok) && tok !== '...')
}

function buildGame(id: string, raw: string): Game {
  const text = raw.replace(/\r\n/g, '\n')
  const headers: Record<string, string> = {}
  for (const m of text.matchAll(HEADER_RE)) headers[m[1]!] = m[2]!

  const movetext = text.slice(text.lastIndexOf(']') + 1)
  const startFen = headers.FEN || START_FEN
  let pos = parseFen(startFen)
  const startBoard = pos.board

  const plies: Ply[] = []
  let number = Number(startFen.split(/\s+/)[5]) || 1
  for (const san of mainLineTokens(movetext)) {
    const move = sanToMove(pos, san)
    if (!move) {
      // An unreadable score would silently truncate; make it loud in dev.
      if (import.meta.env.DEV) console.warn(`[pgn:${id}] illegal move "${san}" after ${plies.length} plies`)
      break
    }
    const color = pos.turn
    pos = applyMove(pos, move)
    plies.push({ san, from: move.from, to: move.to, board: pos.board, number, color })
    if (color === 'b') number++
  }

  const result = headers.Result || '*'
  const orientation: 'w' | 'b' =
    headers.Orientation === 'black' || headers.Orientation === 'b'
      ? 'b'
      : headers.Orientation === 'white' || headers.Orientation === 'w'
        ? 'w'
        : /iwanesko/i.test(headers.Black ?? '')
          ? 'b'
          : 'w'

  return {
    id,
    headers,
    white: headers.White ?? '?',
    black: headers.Black ?? '?',
    whiteElo: headers.WhiteElo,
    blackElo: headers.BlackElo,
    event: headers.Event,
    round: headers.Round,
    date: headers.Date && headers.Date !== '????.??.??' ? headers.Date : undefined,
    result,
    orientation,
    startBoard,
    plies,
  }
}

const games = new Map<string, Game>()
for (const [path, raw] of Object.entries(files)) {
  const id = path.split('/').pop()!.replace(/\.pgn$/, '')
  games.set(id, buildGame(id, raw))
}

export const getGame = (id: string): Game | undefined => games.get(id)
