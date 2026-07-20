/**
 * Decorative chessboard/knight artwork, inline SVG (no network, no CLS —
 * dimensions are intrinsic via viewBox + aspect ratio). Purely visual, so it's
 * hidden from assistive tech and excluded from the accessibility tree.
 */
export function ChessMotif({ className = '' }: { className?: string }) {
  const squares = []
  for (let r = 0; r < 8; r++) {
    for (let c = 0; c < 8; c++) {
      if ((r + c) % 2 === 1) {
        squares.push(<rect key={`${r}-${c}`} x={c * 40} y={r * 40} width={40} height={40} fill="currentColor" />)
      }
    }
  }
  return (
    <svg
      viewBox="0 0 320 320"
      className={className}
      role="presentation"
      aria-hidden="true"
      focusable="false"
    >
      <defs>
        <clipPath id="board-round">
          <rect x="0" y="0" width="320" height="320" rx="24" />
        </clipPath>
        <linearGradient id="board-sheen" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0" stopColor="#ffffff" stopOpacity="0.5" />
          <stop offset="0.5" stopColor="#ffffff" stopOpacity="0" />
        </linearGradient>
        <radialGradient id="knight-glow" cx="0.5" cy="0.4" r="0.6">
          <stop offset="0" stopColor="#d4af37" stopOpacity="0.28" />
          <stop offset="1" stopColor="#d4af37" stopOpacity="0" />
        </radialGradient>
      </defs>
      <g clipPath="url(#board-round)">
        <rect x="0" y="0" width="320" height="320" fill="#faf7f0" />
        <g className="text-ink-900/85" opacity="0.88">{squares}</g>
        <rect x="0" y="0" width="320" height="320" fill="url(#board-sheen)" />
        <circle cx="160" cy="150" r="110" fill="url(#knight-glow)" />
        {/* Gold knight, echoing the favicon. */}
        <g transform="translate(96 78) scale(2.4)" fill="#d4af37">
          <path d="M24 15c6-1 15 3 17 14 1.6 5.7 2 11 2 17H21c0-3 1-5 3-7 2-1.8 3-3 2.5-5-1.6 1.8-4 3.4-6.6 3.4-3 0-4.4-2-3.4-4.6 1-2.7 4-5.4 4-8.2 0-1.6-.8-2.8-2-3.6 2-2.4 5-4.4 9.5-6.4Z" />
          <circle cx="27" cy="24" r="1.7" fill="#0f172a" />
          <rect x="19" y="47" width="26" height="4" rx="1.5" />
        </g>
      </g>
      <rect x="1.5" y="1.5" width="317" height="317" rx="23" fill="none" stroke="#d4af37" strokeWidth="2" strokeOpacity="0.55" />
    </svg>
  )
}
