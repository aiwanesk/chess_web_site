#!/usr/bin/env node
/**
 * Generates chess board diagram SVGs using official Lichess cburnett piece set.
 * Pieces are embedded directly from lichess-org/lila (viewBox 0 0 45 45).
 *
 * Usage: node scripts/gen-diagrams.mjs
 */

import { writeFileSync, mkdirSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const root = join(__dirname, '..');
const outDir = join(root, 'public/images/blog/pontevedra-2026');

mkdirSync(outDir, { recursive: true });

const LIGHT = '#f0d9b5';
const DARK = '#b58863';
const LAST_MOVE_LIGHT = '#cdd16a';
const LAST_MOVE_DARK = '#a6a83a';

// ─── Official Lichess cburnett piece SVG inner content (viewBox 0 0 45 45) ────
// Source: https://github.com/lichess-org/lila/tree/master/public/piece/cburnett
const PIECE_SVG = {
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
};

function parseFen(fen) {
  const rows = fen.split(' ')[0].split('/');
  const board = [];
  for (let r = 0; r < 8; r++) {
    const row = [];
    for (const ch of rows[r]) {
      if (ch >= '1' && ch <= '8') {
        for (let i = 0; i < parseInt(ch); i++) row.push(null);
      } else {
        row.push(ch);
      }
    }
    board.push(row);
  }
  return board;
}

function generateBoardSVG(fen, options = {}) {
  const { lastMove, flip = false } = options;
  const size = 360;
  const sq = size / 8; // 45
  const board = parseFen(fen);
  const files = 'abcdefgh';

  const lastFrom = lastMove ? lastMove.slice(0, 2) : null;
  const lastTo = lastMove ? lastMove.slice(2, 4) : null;

  let svg = `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${size} ${size}" width="${size}" height="${size}">
<rect width="${size}" height="${size}" rx="3" fill="${DARK}"/>`;

  // Draw squares
  for (let r = 0; r < 8; r++) {
    for (let c = 0; c < 8; c++) {
      const displayR = flip ? 7 - r : r;
      const displayC = flip ? 7 - c : c;
      const isDark = (displayR + displayC) % 2 === 1;
      const x = c * sq;
      const y = r * sq;
      const sqName = files[displayC] + (8 - displayR);
      const isLastMove = sqName === lastFrom || sqName === lastTo;

      let fill;
      if (isLastMove) {
        fill = isDark ? LAST_MOVE_DARK : LAST_MOVE_LIGHT;
      } else {
        fill = isDark ? DARK : LIGHT;
      }
      svg += `\n<rect x="${x}" y="${y}" width="${sq}" height="${sq}" fill="${fill}"/>`;
    }
  }

  // Coordinates
  for (let i = 0; i < 8; i++) {
    const displayI = flip ? 7 - i : i;
    const rankDark = (displayI) % 2 === 1;
    const rankColor = rankDark ? LIGHT : DARK;
    svg += `\n<text x="2" y="${i * sq + 11}" font-size="9" font-weight="700" font-family="-apple-system,system-ui,sans-serif" fill="${rankColor}">${8 - displayI}</text>`;
    const fileDark = (7 + displayI) % 2 === 0;
    const fileColor = fileDark ? LIGHT : DARK;
    svg += `\n<text x="${i * sq + sq - 2}" y="${size - 2}" font-size="9" font-weight="700" font-family="-apple-system,system-ui,sans-serif" text-anchor="end" fill="${fileColor}">${files[displayI]}</text>`;
  }

  // Draw pieces — direct embed of lichess cburnett SVGs
  for (let r = 0; r < 8; r++) {
    for (let c = 0; c < 8; c++) {
      const displayR = flip ? 7 - r : r;
      const displayC = flip ? 7 - c : c;
      const piece = board[displayR][displayC];
      if (!piece) continue;

      const x = c * sq;
      const y = r * sq;
      const pieceSvg = PIECE_SVG[piece];
      if (pieceSvg) {
        // Pieces are designed for 45x45 viewBox, squares are 45x45 → 1:1 fit
        svg += `\n<g transform="translate(${x},${y})">${pieceSvg}</g>`;
      }
    }
  }

  svg += '\n</svg>';
  return svg;
}

// ──────────────────────────────────────────────────────────────
// POSITIONS — Replace with your actual game positions
// ──────────────────────────────────────────────────────────────

const diagrams = [
  {
    file: 'ronde1-diagramme1.svg',
    fen: 'r1bq1rk1/p3npb1/1pnp2pp/1Pp1p3/2P5/P1NP1NP1/4PPBP/1RBQ1RK1 b - - 0 11',
    lastMove: 'a4b5',
  },
  {
    file: 'ronde1-diagramme2.svg',
    fen: '2bq1rk1/p3npb1/1p1p2pp/1Pp5/2P5/P1pP2P1/4PPBP/1RBQ1RK1 b - - 1 14',
  },
  {
    file: 'ronde1-diagramme3.svg',
    fen: '5k2/p2b1p2/1p4p1/7p/3BP3/P2P2P1/4KP1P/8 w - - 1 31',
  },
  {
    file: 'ronde1-diagramme4.svg',
    fen: '8/p4p2/1p1k1Bp1/3P3p/4PKPP/P4P2/8/7b b - - 0 38',
  },
  {
    file: 'ronde2-diagramme1.svg',
    fen: '1rr5/pp3p1p/4nkp1/8/8/1P4P1/P3RPBP/R6K b - - 1 22',
    flip: true,
  },
  {
    file: 'ronde2-diagramme2.svg',
    fen: '8/5p2/1p1rnkpp/p1r5/P1B4P/1P2R1P1/5PK1/4R3 b - - 2 29',
    flip: true,
  },
  {
    file: 'ronde2-diagramme3.svg',
    fen: '8/5p2/1p1rnk1p/p1r3p1/P1B3PP/1P2R1K1/5P2/4R3 w - - 0 32',
    flip: true,
  },
  {
    file: 'ronde2-diagramme4.svg',
    fen: '8/5p2/1p2nk1p/p1r3p1/P1BrR1PP/1P4K1/5P2/4R3 w - - 2 33',
    flip: true,
  },
  {
    file: 'ronde3-diagramme.svg',
    fen: 'r4rk1/pp2bppp/2n1pn2/3p4/3P4/2N1BN2/PPQ1BPPP/R4RK1 b - - 0 22',
    lastMove: 'd1c2',
  },
  {
    file: 'ronde4-diagramme1.svg',
    fen: 'r1r3k1/1pq1bppp/p2pbn2/4p3/P3P3/1PNQ2PP/2PB1PB1/R3R1K1 b - - 0 15',
  },
  {
    file: 'ronde4-diagramme2.svg',
    fen: 'r1r3k1/1pq2ppp/pb1pbn2/4p3/P3PP2/1P1Q2PP/2PBN1BK/R1R5 b - - 2 19',
  },
  {
    file: 'ronde4-diagramme3.svg',
    fen: 'r2r2k1/1p3pp1/p1q1bn1p/8/P1P1pP2/1P2Q1PP/4N1BK/R1R5 b - - 2 28',
  },
  {
    file: 'ronde4-diagramme4.svg',
    fen: '3r2k1/1p3pp1/p1q2n1p/8/P1P1pP2/1P1r2PB/4NQ1K/R1R5 w - - 3 31',
  },
  {
    file: 'ronde5-diagramme1.svg',
    fen: 'r4rk1/pp2b1pp/1nq2p2/2p1p3/2P1P3/1P2NPP1/1B4KP/R3QR2 b - - 0 19',
  },
  {
    file: 'ronde5-diagramme2.svg',
    fen: 'r2r2k1/1p1q2pp/3b1p2/p1pP4/2P1R3/1P4P1/1B2Q1KP/R7 b - - 2 25',
  },
  {
    file: 'ronde5-diagramme3.svg',
    fen: '4R3/1p1q1kpp/3b1p2/r1pP4/2P5/1P4P1/1B2Q1KP/8 w - - 1 28',
  },
  {
    file: 'ronde5-diagramme4.svg',
    fen: '5k2/1p1q2pp/3bRp2/2pP4/2P5/1P4P1/rBQ4P/6K1 b - - 6 30',
  },
  {
    file: 'ronde5-diagramme5.svg',
    fen: '6k1/6p1/2R4p/2pP4/2P5/6P1/r6P/6K1 b - - 0 38',
  },
  {
    file: 'ronde6-diagramme1.svg',
    fen: 'r1bqk1nr/1p2p1bp/p1n3p1/4P3/3p1B2/5N1P/PP2NPP1/R2QKB1R b KQkq - 1 11',
  },
  {
    file: 'ronde6-diagramme2.svg',
    fen: 'r1b1k1nr/1p2p2p/p5p1/8/3N4/7P/PP3PP1/R3KB1R b KQkq - 0 16',
    flip: true,
  },
  {
    file: 'ronde6-diagramme3.svg',
    fen: '4k3/1p2p2p/p2n2p1/8/2r3P1/7P/PPN1RPK1/8 b - - 6 30',
    flip: true,
  },
  {
    file: 'ronde6-diagramme4.svg',
    fen: '8/1p2k2p/p3p1n1/6p1/6P1/6NP/PP1R1PK1/2r5 w - - 10 39',
    flip: true,
  },
  {
    file: 'ronde6-diagramme5.svg',
    fen: '8/1p2k2p/p3p1n1/6p1/6P1/5K1P/PP1RNP2/4r3 b - - 17 42',
    flip: true,
  },
  {
    file: 'ronde6-diagramme6.svg',
    fen: '8/1p5p/p3pkn1/6p1/6P1/5PKP/PP1RN3/7r w - - 1 46',
    flip: true,
  },
  {
    file: 'ronde7-diagramme1.svg',
    fen: '1rbq1rk1/3nbpp1/2n1p2p/2ppP2P/1p6/3P1NP1/1PP2PB1/R1BQRNK1 w - - 0 15',
  },
  {
    file: 'ronde7-diagramme2.svg',
    fen: '2b1qrk1/4b1p1/2n1p2p/2ppN2P/5B2/2PP2PB/5P2/3QR1K1 b - - 4 25',
  },
  {
    file: 'ronde7-diagramme3.svg',
    fen: '2b3k1/3q2p1/4pr1p/3pR2P/3P4/6PB/4QP2/6K1 b - - 0 31',
  },
  {
    file: 'ronde7-diagramme4.svg',
    fen: 'R7/3b2pk/4pr1p/3p3P/2qP4/2Q3PB/5P2/6K1 w - - 13 38',
  },
  {
    file: 'ronde7-diagramme5.svg',
    fen: '4b3/3r3k/4p2p/6pP/2RP4/6P1/5PB1/6K1 w - - 1 43',
  },
  {
    file: 'ronde7-diagramme6.svg',
    fen: '8/5k2/8/5R1P/3pr1p1/6P1/5P2/5K2 b - - 1 53',
  },
  {
    file: 'ronde7-diagramme7.svg',
    fen: '8/8/4k3/4r2P/3p1Rp1/6P1/5P2/5K2 w - - 4 55',
  },
  {
    file: 'ronde7-diagramme8.svg',
    fen: '8/4r3/7P/5k2/3R2p1/6P1/5P2/5K2 w - - 1 57',
  },
  {
    file: 'ronde8-diagramme1.svg',
    fen: 'r2q1rk1/1p2bppp/p1npb3/2nNp3/2P1P3/2NBB3/PP3PPP/R2Q1RK1 w - - 12 13',
  },
  {
    file: 'ronde8-diagramme2.svg',
    fen: 'r2q1rk1/1p2bppp/p1n1b3/2p1p3/2P1P3/2NBB3/P4PPP/1R1Q1RK1 w - - 0 16',
  },
  {
    file: 'ronde8-diagramme3.svg',
    fen: 'r2q1rk1/1p3ppp/p3b3/2b1p3/1NP1P3/3B4/P4PPP/1R1Q1RK1 b - - 0 18',
  },
  {
    file: 'ronde8-diagramme4.svg',
    fen: '2b2rk1/1p3pp1/p2r4/2bNp2p/2P1P3/3B2PP/P5P1/1R3R1K b - - 0 24',
  },
  {
    file: 'ronde8-diagramme5.svg',
    fen: '2b2r2/1p3pk1/p2r2p1/2bNp2p/2P1P3/3B1RPP/P5P1/5R1K b - - 3 26',
  },
  {
    file: 'ronde8-diagramme6.svg',
    fen: '1r6/5pk1/p3b1p1/1p1Np3/2PbP3/7P/P3B1P1/5R1K b - - 1 33',
  },
  {
    file: 'ronde8-diagramme7.svg',
    fen: '8/5pk1/B3b1p1/3Np3/3bP3/7P/r5P1/5R1K w - - 0 36',
  },
  {
    file: 'ronde9-diagramme1.svg',
    fen: '1r1r2k1/ppq1bpp1/2n1bn1p/P3p3/Q1p1P3/2P2NP1/1P3PBP/R1B1RNK1 w - - 0 16',
  },
  {
    file: 'ronde9-diagramme2.svg',
    fen: '2qr2k1/pp3pp1/4b3/P7/2pP1B2/6Pp/1P3PBP/R3R1K1 w - - 0 26',
  },
  {
    file: 'ronde9-diagramme3.svg',
    fen: '8/1p3ppk/1B2b3/P7/2pR4/5qPp/1P3P1P/5BK1 b - - 4 33',
  },
  {
    file: 'ronde9-diagramme4.svg',
    fen: '8/1p3p1k/1B6/P4bp1/8/4RP1p/4BK1P/2q5 b - - 3 40',
    flip: true,
  },
];

console.log('Generating diagrams with official Lichess cburnett pieces…\n');

for (const d of diagrams) {
  const svg = generateBoardSVG(d.fen, { lastMove: d.lastMove, flip: d.flip });
  writeFileSync(join(outDir, d.file), svg, 'utf-8');
  console.log(`  ✓ ${d.file}`);
}

writeFileSync(join(outDir, 'og-pontevedra.svg'), generateBoardSVG(diagrams[0].fen, { lastMove: diagrams[0].lastMove, flip: diagrams[0].flip }), 'utf-8');
console.log(`  ✓ og-pontevedra.svg`);

console.log(`\nDone! public/images/blog/pontevedra-2026/`);
