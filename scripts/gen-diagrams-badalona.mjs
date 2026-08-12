#!/usr/bin/env node
/**
 * Generates Badalona 2026 board diagram SVGs (Lichess cburnett piece set).
 * Reuses the same renderer as gen-diagrams.mjs but writes to badalona-2026/.
 *
 * Usage: node scripts/gen-diagrams-badalona.mjs
 */

import { writeFileSync, mkdirSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const root = join(__dirname, '..');
const outDir = join(root, 'public/images/blog/badalona-2026');

mkdirSync(outDir, { recursive: true });

const LIGHT = '#f0d9b5';
const DARK = '#b58863';
const LAST_MOVE_LIGHT = '#cdd16a';
const LAST_MOVE_DARK = '#a6a83a';

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
  const sq = size / 8;
  const board = parseFen(fen);
  const files = 'abcdefgh';

  const lastFrom = lastMove ? lastMove.slice(0, 2) : null;
  const lastTo = lastMove ? lastMove.slice(2, 4) : null;

  let svg = `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${size} ${size}" width="${size}" height="${size}">
<rect width="${size}" height="${size}" rx="3" fill="${DARK}"/>`;

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
      if (isLastMove) fill = isDark ? LAST_MOVE_DARK : LAST_MOVE_LIGHT;
      else fill = isDark ? DARK : LIGHT;
      svg += `\n<rect x="${x}" y="${y}" width="${sq}" height="${sq}" fill="${fill}"/>`;
    }
  }

  for (let i = 0; i < 8; i++) {
    const displayI = flip ? 7 - i : i;
    const rankDark = displayI % 2 === 1;
    const rankColor = rankDark ? LIGHT : DARK;
    svg += `\n<text x="2" y="${i * sq + 11}" font-size="9" font-weight="700" font-family="-apple-system,system-ui,sans-serif" fill="${rankColor}">${8 - displayI}</text>`;
    const fileDark = (7 + displayI) % 2 === 0;
    const fileColor = fileDark ? LIGHT : DARK;
    svg += `\n<text x="${i * sq + sq - 2}" y="${size - 2}" font-size="9" font-weight="700" font-family="-apple-system,system-ui,sans-serif" text-anchor="end" fill="${fileColor}">${files[displayI]}</text>`;
  }

  for (let r = 0; r < 8; r++) {
    for (let c = 0; c < 8; c++) {
      const displayR = flip ? 7 - r : r;
      const displayC = flip ? 7 - c : c;
      const piece = board[displayR][displayC];
      if (!piece) continue;
      const x = c * sq;
      const y = r * sq;
      const pieceSvg = PIECE_SVG[piece];
      if (pieceSvg) svg += `\n<g transform="translate(${x},${y})">${pieceSvg}</g>`;
    }
  }

  svg += '\n</svg>';
  return svg;
}

const diagrams = [
  {
    // Ronde 1 — setup ambitieux : paire de fous rendue (Bxf3) pour le centre,
    // cavalier planté en d5. Position après 13...Cxd5, trait aux Blancs.
    file: 'ronde1-diagramme1.svg',
    fen: 'r2q1rk1/1p3pbp/p1n3p1/2pnp3/8/2PP1BPP/PPN2P2/R1BQR1K1 w - - 0 14',
  },
  {
    // Ronde 1 — juste avant 19...f4! : les Blancs ont poussé h4-h5-hxg6 sur le
    // roi noir. Trait aux Noirs, le pion blanc vient d'atterrir en g6.
    file: 'ronde1-diagramme2.svg',
    fen: '3r1r1k/1p1qn1bp/p1n3P1/2p1pp2/P7/1QPPNBP1/1P3P2/R1B1R1K1 b - - 0 19',
    lastMove: 'h5g6',
  },
  {
    // Ronde 1 — position d'attaque juste avant 25.Cf1 : dame en h3, tour en f4,
    // toute l'armée noire braquée sur le roque blanc. Trait aux Blancs.
    file: 'ronde1-diagramme3.svg',
    fen: '5r1k/1p4b1/p1n3p1/2pQp3/P4r2/2PPN2q/1P2RP2/R1B3K1 w - - 0 25',
    lastMove: 'h7g6',
  },
  {
    // Ronde 1 — le mat final : 27...Tf1#. La Ta1 ne peut pas reprendre (Fc1 bloque),
    // la Dh3 verrouille g2/h2.
    file: 'ronde1-diagramme4.svg',
    fen: '7k/1p4b1/p1n3p1/2pQp3/P7/2PP2Pq/1P2R3/R1B2rK1 w - - 0 28',
    lastMove: 'f8f1',
  },

  // ─── Ronde 2 (Blancs vs Tripulskiy) ──────────────────────────────
  {
    // R2 — structure type Benoni/Est-Indienne. Blancs très confiants,
    // sur le point de lancer le plan ambitieux g4!?.
    file: 'ronde2-diagramme1.svg',
    fen: 'r1bqk1nr/pp1nppbp/3p2p1/2pP4/4P3/2N5/PPP1BPPP/R1BQK1NR w KQkq - 0 8',
  },
  {
    // R2 — après l'échange du fou noir contre le Cc3 : le précis Db1,
    // jouant sur le manque de cases de la dame noire.
    file: 'ronde2-diagramme2.svg',
    fen: 'r1b1k2r/pp1npp1p/1n1p2p1/q1pP2P1/4P3/2P2P2/P1PBB2P/R2QK1NR w KQkq - 0 15',
  },
  {
    // R2 — le trop ambitieux a4 : très chaud pour tout le monde, calcul précis.
    file: 'ronde2-diagramme3.svg',
    fen: 'r1b1k2r/pp1npp1p/1n1p2p1/2qP2P1/PQp1P3/2P2P2/2PBB2P/R3K1NR b KQkq - 0 18',
    lastMove: 'a2a4',
  },
  {
    // R2 — pièces replacées après avoir rendu le pion a4 (Cxa4) : les Blancs
    // vont finir d'améliorer avec Ff3!.
    file: 'ronde2-diagramme4.svg',
    fen: '4kr2/1pqbpp1p/rn1p2p1/p2P2PP/n1pBPP2/2P5/2PNB3/1RQ1K2R w - - 0 24',
  },
  {
    // R2 — après Ff3! et de petits échanges : la rupture précise e5!!.
    file: 'ronde2-diagramme5.svg',
    fen: '4kr2/1pqb4/rn1pp1p1/p5P1/n1pBPP2/2P2B2/2PN4/1RQ1K2R w - - 0 27',
  },
  {
    // R2 — pièces noires horribles : Fg4, Cf3-h4, roi centralisé. Les pions tombent.
    file: 'ronde2-diagramme6.svg',
    fen: '4b1r1/1p5R/rnk1p1p1/p2pP1P1/n1pB1PBN/2P1K3/2P5/1R6 w - - 0 35',
  },

  // ─── Ronde 3 (Noirs vs MI Utsab Chatterjee, 2353) — vues côté Noirs ─────
  {
    // R3 — la position rare : les Blancs vont gagner un pion (Fxb7) et gardent
    // la paire de fous en position ouverte, et pourtant c'est égal (activité noire).
    // Les Noirs trouvent Cf6! Fxb7 Ta7 Fe4 Td7!.
    file: 'ronde3-diagramme1.svg',
    fen: 'r2q1rk1/1p2bppp/p7/3npP2/3nB3/2N5/PP3PPP/R1BQR1K1 b - - 0 18',
    flip: true,
  },
  {
    // R3 — ici il fallait stabiliser avec Dc8 ; les Noirs partent pour Fc5.
    file: 'ronde3-diagramme2.svg',
    fen: '3q1rk1/3rbppp/p4n2/4pP2/Q2nB3/2N5/PP3PPP/R1B1R1K1 b - - 0 21',
    flip: true,
  },
  {
    // R3 — le choix : Cb5 (position inférieure mais tenable) ou l'all-in Cf3.
    file: 'ronde3-diagramme3.svg',
    fen: '5rk1/5p1p/pq1r1p2/2b1pP2/2QnB3/2N5/PP3PPP/3RR1K1 b - - 0 25',
    flip: true,
  },
  {
    // R3 — après l'all-in : les Noirs espèrent du jeu, les Blancs jouent très pratique.
    file: 'ronde3-diagramme4.svg',
    fen: '6k1/5p1p/pq3p2/4pP2/3r4/5B2/PPQ1K1PP/3N4 b - - 0 28',
    flip: true,
  },
  {
    // R3 — la fin : les Noirs tentent Dg5, les Blancs répondent Df4, abandon.
    file: 'ronde3-diagramme5.svg',
    fen: '3r2k1/5p2/5p2/5P2/P3Q2p/1P1B1K2/3q1NP1/8 b - - 0 35',
    flip: true,
  },

  // ─── Ronde 4 (Blancs vs Tiwari, 2039) — vues côté Blancs ─────────────
  {
    // R4 — une sorte de française d'avance. Les Blancs ne sont pas sereins sur
    // le très tendu ...f6!, mais les Noirs préfèrent ...Fd7.
    file: 'ronde4-diagramme1.svg',
    fen: 'r1bq1rk1/pp3ppp/1bn1p1n1/3pP3/8/1N1B1N2/PPP2PPP/R1BQR1K1 b - - 0 14',
    lastMove: 'f1e1',
  },
  {
    // R4 — les Blancs « abusent » avec Fc5, sûrs que les Noirs ne sacrifieront
    // jamais la qualité. Suivra ...Te8 puis Fd4.
    file: 'ronde4-diagramme2.svg',
    fen: 'r2q1rk1/ppbb1ppp/2n1p1n1/2BpP3/8/1N1B1N2/PPP2PPP/R2QR1K1 b - - 0 16',
  },
  {
    // R4 — les Blancs songent à Cg5 (rien de clair) et jouent De2. ...Dd8 Tad1 Cb4?.
    file: 'ronde4-diagramme3.svg',
    fen: 'rq2r1k1/ppbb1ppp/2n1p1n1/3pP3/3B4/1N1B1N2/PPP2PPP/R2QR1K1 w - - 0 19',
    lastMove: 'd8b8',
  },
  {
    // R4 — après ...Cb4? : Fxg6 fxg6?! c3 Ca6 Cc5! et les Blancs ont toutes les bonnes pièces.
    file: 'ronde4-diagramme4.svg',
    fen: 'r2qr1k1/ppbb1ppp/4p1n1/3pP3/1n1B4/1N1B1N2/PPP1QPPP/3RR1K1 w - - 0 20',
    lastMove: 'c6b4',
  },
  {
    // R4 — le cavalier s'est échangé en c5 : il reste cavalier (d4) contre le très
    // mauvais fou français (d7), pendant que les Noirs massent sur la colonne c.
    file: 'ronde4-diagramme4b.svg',
    fen: '2r3k1/ppqb2pp/4p1p1/3pP3/2rN4/2P5/PP1Q1PPP/2R1R1K1 w - - 0 23',
  },
  {
    // R4 — trop ambitieux : au lieu de Td3 qui tient tout, les Blancs partent sur Cb3.
    file: 'ronde4-diagramme5.svg',
    fen: '2r3k1/3b2pp/4p1p1/ppqpP3/2rN4/P1P1R3/1P1Q1PPP/4R1K1 w - - 0 27',
    lastMove: 'a6a5',
  },
  {
    // R4 — « il en veut trop » avec ...Da2 ; les Blancs prennent le pion, ...Da7 vise f2.
    file: 'ronde4-diagramme6.svg',
    fen: '2r3k1/3b2p1/4p1pp/3pP3/1pr5/2P2R1P/qP1QNPP1/4R1K1 w - - 0 31',
  },
  {
    // R4 — l'effondrement : De3? (Te1! gagnait pour les Noirs), ...Tf1!, puis b4 et les
    // Blancs ratent que ...Da1 Cg1 d4! gagne. Il fallait tenir avec Tc1.
    file: 'ronde4-diagramme7.svg',
    fen: '5rk1/1RR2bp1/4pqp1/3p4/5P1p/1P5P/3QN1PK/r7 w - - 0 38',
  },

  // ─── Ronde 5 (Noirs vs CM Fernandez-Diaz, 1983) — vues côté Noirs ─────
  {
    // R5 — anglaise à double fianchetto. Trait aux Blancs ; les Noirs estiment
    // avoir déjà réglé tous leurs problèmes (un poil mieux, facile à jouer).
    file: 'ronde5-diagramme1.svg',
    fen: 'rn3rk1/pp2ppbp/5np1/7q/2N3b1/1P2P1P1/PB3PBP/RN1Q1RK1 w - - 0 15',
    flip: true,
  },
  {
    // R5 — les Blancs jouent Cba3 : le cavalier va sur la mauvaise case et ne sert à rien.
    file: 'ronde5-diagramme2.svg',
    fen: 'rn3rk1/pp2ppbp/4bnp1/7q/2N5/NP2PPP1/PB4BP/R2Q1RK1 b - - 0 17',
    lastMove: 'b1a3',
    flip: true,
  },
  {
    // R5 — les Noirs trouvent ...b5!, calculé jusqu'au bout (20 min de vérification).
    file: 'ronde5-diagramme3.svg',
    fen: 'r2r2k1/pp2ppbp/2n1bnp1/7q/2N5/NP2PPP1/PB1R2BP/R2Q2K1 b - - 0 20',
    flip: true,
  },
  {
    // R5 — la pointe finale : après Dxb5, les Noirs jouent ...Ce4! et la partie est finie.
    file: 'ronde5-diagramme4.svg',
    fen: 'r5k1/p3ppbp/4bnp1/1Q6/8/NP2qPP1/PB4BP/R6K b - - 0 25',
    flip: true,
  },

  // ─── Ronde 6 (Blancs vs Lunin Bohomolets, 1872) — vues côté Blancs ────
  {
    // R6 — Caro-Kann. Les Noirs reroutent le cavalier par a6 vers c7 (...Ca6!?),
    // les Blancs comptent échanger leur mauvaise pièce par Ff4 mais les Noirs n'ont plus de pbs.
    file: 'ronde6-diagramme1.svg',
    fen: 'r1bqk2r/ppn2pp1/2pb1p2/7p/3P4/2PB4/PPQ1PPPP/R1B1K2R w KQkq - 0 12',
  },
  {
    // R6 — ici il fallait h3 pour garder un micro-plus ; les Blancs jouent Te1 et
    // encaissent ...Fg4 Fxe2 (les Noirs sont très bien).
    file: 'ronde6-diagramme2.svg',
    fen: 'r1b4r/pp2npk1/2pq1pp1/7p/2PP4/P2B4/1PQ1NPPP/3R1RK1 w - - 0 18',
    lastMove: 'g8g7',
  },
  {
    // R6 — la gaffe : b4 De6 Ff3?? Txd4, les Blancs perdent un pion (paire d'échecs e1/e5).
    file: 'ronde6-diagramme3.svg',
    fen: '3r4/pp1rnpk1/2pq1pp1/7p/2PPB3/P6P/1PQR1PP1/3R2K1 w - - 0 22',
  },
  {
    // R6 — les Noirs (un pion en plus, mais pions f doublés) forcent l'échange des dames.
    file: 'ronde6-diagramme4.svg',
    fen: '8/pp3pk1/2p2pp1/5n1p/1PPq4/P2Q1PP1/5P2/3B2K1 w - - 0 30',
  },
  {
    // R6 — finale F vs C : les Noirs ratent ...c5! (chances de gain) ; une finale de pions
    // gagnerait pour les Blancs (pions f doublés). Il reste solide, nulle.
    file: 'ronde6-diagramme5.svg',
    fen: '8/pp3p2/2p2kp1/5p1p/1PPn1P2/P3K1PP/8/3B4 b - - 0 40',
  },

  // ─── Ronde 7 (Noirs vs Carrasco Holgado, 2040) — vues côté Noirs ─────
  {
    // R7 — le français à la ...a6 (ligne du GM Bauer) : une belle française d'avance.
    file: 'ronde7-diagramme1.svg',
    fen: 'rn1qkb1r/1pp1nppp/p3p3/3pP3/P2P4/2P2N2/1P3PPP/R1BQKN1R b KQkq - 0 8',
    lastMove: 'g1f1',
    flip: true,
  },
  {
    // R7 — le très intéressant ...c4 : figer la structure et viser les cases faibles b3/a4/a5.
    file: 'ronde7-diagramme2.svg',
    fen: 'r3kb1r/1p1qnppp/p1n1p3/2ppP3/P2P4/2P2NN1/1P3PPP/R1BQR1K1 b - - 0 13',
    lastMove: 'f1e1',
    flip: true,
  },
  {
    // R7 — le fort ...0-0-0 : Cxg7 est contré par ...h5 et ...Tg8! (le cavalier est piégé).
    file: 'ronde7-diagramme3.svg',
    fen: 'r3k2r/1p1qbppp/p1n1p3/Pn1pP2N/3P3P/1QP2N2/3B1PP1/R3R1K1 b kq - 0 17',
    lastMove: 'g3h5',
    flip: true,
  },
  {
    // R7 — il tente de s'activer en donnant a5 ; ...Cxa5 gagne le pion, puis on stabilise.
    file: 'ronde7-diagramme4.svg',
    fen: '1kr4r/1pq1bppp/p3p3/nn1pP3/3P3P/1QPNBN2/5PP1/R1R3K1 w - - 0 22',
    lastMove: 'c6a5',
    flip: true,
  },
  {
    // R7 — la conclusion : ...Cc2 vient croquer d4, la partie est pliée.
    file: 'ronde7-diagramme5.svg',
    fen: '1k6/1p3pp1/p3p2p/1n1pP3/1P1P1BPP/n1r2N2/R4PK1/8 b - - 0 30',
    lastMove: 'g1g2',
    flip: true,
  },

  // ─── Ronde 8 (Blancs vs MI Perpinya, 2149) — vues côté Blancs ────────
  {
    // R8 — sicilienne, structure Boleslavsky (Fe2) ; les Noirs choisissent ...Cxd4 et ...g6.
    // Les Blancs tankent 40 min sur e5 (SUPER intéressant... ce qui n'est pas bon).
    file: 'ronde8-diagramme1.svg',
    fen: 'r1bqkb1r/pp2pppp/3p1n2/8/3QP3/2N5/PPP1BPPP/R1B1K2R b KQkq - 0 8',
    lastMove: 'd1d4',
  },
  {
    // R8 — un pion de moins, mais de bonnes compensations.
    file: 'ronde8-diagramme2.svg',
    fen: 'r1k2b1r/pp2pp1p/4bnp1/4p1B1/8/2N5/PPP1BPPP/2KR3R w - - 0 13',
  },
  {
    // R8 — les Blancs construisent tranquillement (Ca4 en route vers c5).
    file: 'ronde8-diagramme3.svg',
    fen: 'r6r/1pk1ppbp/p3bnp1/4p3/N7/4B3/PPPRBPPP/2KR4 b - - 0 18',
    lastMove: 'c3a4',
  },
  {
    // R8 — position dont Alex est très satisfait ; il pèse a4/b4/g4 et joue le brutal g4.
    file: 'ronde8-diagramme4.svg',
    fen: 'r6r/2k1ppbp/p4np1/1pN1pb2/8/4B3/PPPRBPPP/2KR4 w - - 0 22',
    lastMove: 'e6f5',
  },
  {
    // R8 — la gaffe : Cd7?? au lieu du parfait Cb7 (b6 était couvert). Simplification,
    // finale de tours perdante.
    file: 'ronde8-diagramme5.svg',
    fen: 'r1k4r/4Rp1p/p5pb/1pN1pb2/8/4P3/PPP1B2P/2KR4 w - - 0 30',
  },
];

console.log('Generating Badalona diagrams with official Lichess cburnett pieces…\n');
for (const d of diagrams) {
  const svg = generateBoardSVG(d.fen, { lastMove: d.lastMove, flip: d.flip });
  writeFileSync(join(outDir, d.file), svg, 'utf-8');
  console.log(`  ✓ ${d.file}`);
}
console.log(`\nDone! public/images/blog/badalona-2026/`);
