import type { SVGProps } from 'react'

/*
  Inline SVG icon set — no icon-font, no network, no CLS (intrinsic 24×24
  viewBox). Stroke icons use currentColor so they inherit text color. Purely
  decorative by default (aria-hidden); pass a `title`/aria-label at the call
  site if an icon ever carries meaning.
*/

type IconProps = SVGProps<SVGSVGElement> & { size?: number | string }

function Base({ size = 24, children, ...rest }: IconProps & { children: React.ReactNode }) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={1.6}
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
      focusable="false"
      {...rest}
    >
      {children}
    </svg>
  )
}

/** Chess knight — the brand mark, drawn as a solid silhouette. */
export function IconKnight({ size = 24, className, ...rest }: IconProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="currentColor"
      className={className}
      aria-hidden="true"
      focusable="false"
      {...rest}
    >
      <path d="M8.2 3.3c.5 1 .4 1.8-.1 2.6C6.7 7.7 5.4 9.4 5 11.6c-.2 1.1.5 1.9 1.5 1.6.9-.3 1.6-1 2.3-1.7.4.9-.2 1.6-1 2.3-1.4 1.3-2.3 2.7-2.5 4.6-.1.7.3 1.2 1 1.2h8.3c.6 0 1-.4 1-1 .1-2.7-.2-5-1.1-7.3-1-2.6-2.7-4.6-5.3-5.5.3-.8.4-1.6.1-2.5-.1-.4-.4-.6-.8-.6-.4 0-.7.3-.6.7Z" />
      <rect x="5.4" y="19.4" width="12.2" height="1.9" rx="0.9" />
    </svg>
  )
}

export function IconTarget(p: IconProps) {
  return (
    <Base {...p}>
      <circle cx="12" cy="12" r="8.2" />
      <circle cx="12" cy="12" r="4.2" />
      <circle cx="12" cy="12" r="0.4" fill="currentColor" />
    </Base>
  )
}

export function IconBoard(p: IconProps) {
  return (
    <Base {...p}>
      <rect x="3.5" y="3.5" width="17" height="17" rx="2" />
      <path d="M3.5 9h17M3.5 14.5h17M9 3.5v17M14.5 3.5v17" strokeWidth={1.2} />
    </Base>
  )
}

export function IconTrophy(p: IconProps) {
  return (
    <Base {...p}>
      <path d="M7 4.5h10v3a5 5 0 0 1-10 0v-3Z" />
      <path d="M7 5.5H4.5a2.5 2.5 0 0 0 3 2.4M17 5.5h2.5a2.5 2.5 0 0 1-3 2.4" />
      <path d="M12 12.5v3M9 19.5h6M9.5 19.5c0-1.6.9-2.4 2.5-2.4s2.5.8 2.5 2.4" />
    </Base>
  )
}

export function IconClock(p: IconProps) {
  return (
    <Base {...p}>
      <circle cx="12" cy="12" r="8.2" />
      <path d="M12 7.6V12l3 1.8" />
    </Base>
  )
}

export function IconRoute(p: IconProps) {
  return (
    <Base {...p}>
      <circle cx="6.5" cy="18" r="2.2" />
      <circle cx="17.5" cy="6" r="2.2" />
      <path d="M8.7 18h5.3a3 3 0 0 0 0-6H10a3 3 0 0 1 0-6h5.3" />
    </Base>
  )
}

export function IconUsers(p: IconProps) {
  return (
    <Base {...p}>
      <circle cx="9" cy="8" r="3" />
      <path d="M3.5 19.5a5.5 5.5 0 0 1 11 0" />
      <path d="M16 5.5a3 3 0 0 1 0 5.9M17 14.2a5.5 5.5 0 0 1 3.5 5.3" />
    </Base>
  )
}

export function IconMonitor(p: IconProps) {
  return (
    <Base {...p}>
      <rect x="3" y="4.5" width="18" height="12" rx="2" />
      <path d="M9 20h6M12 16.5V20" />
    </Base>
  )
}

export function IconGraduation(p: IconProps) {
  return (
    <Base {...p}>
      <path d="M2.5 8.5 12 4.5l9.5 4-9.5 4-9.5-4Z" />
      <path d="M6.5 10.5v4.2c0 1.3 2.5 2.3 5.5 2.3s5.5-1 5.5-2.3v-4.2M21.5 8.5v4.5" />
    </Base>
  )
}

export function IconSpark(p: IconProps) {
  return (
    <Base {...p}>
      <path d="M12 3.5c.4 3.4 1.6 4.6 5 5-3.4.4-4.6 1.6-5 5-.4-3.4-1.6-4.6-5-5 3.4-.4 4.6-1.6 5-5Z" />
      <path d="M18.5 14c.2 1.6.8 2.2 2.4 2.4-1.6.2-2.2.8-2.4 2.4-.2-1.6-.8-2.2-2.4-2.4 1.6-.2 2.2-.8 2.4-2.4Z" />
    </Base>
  )
}

export function IconCheck(p: IconProps) {
  return (
    <Base strokeWidth={2} {...p}>
      <path d="M4.5 12.5 9.5 17.5 20 6.5" />
    </Base>
  )
}

export function IconArrowRight(p: IconProps) {
  return (
    <Base strokeWidth={1.8} {...p}>
      <path d="M4.5 12h15M13 5.5 19.5 12 13 18.5" />
    </Base>
  )
}

export function IconMail(p: IconProps) {
  return (
    <Base {...p}>
      <rect x="3" y="5" width="18" height="14" rx="2.5" />
      <path d="m4 7.5 8 5.2 8-5.2" />
    </Base>
  )
}

export function IconPhone(p: IconProps) {
  return (
    <Base {...p}>
      <path d="M6.5 4h3l1.5 4-2 1.4a11 11 0 0 0 5.1 5.1L15.5 12l4 1.5v3a2 2 0 0 1-2.2 2A15.5 15.5 0 0 1 4.5 6.2 2 2 0 0 1 6.5 4Z" />
    </Base>
  )
}

export function IconPin(p: IconProps) {
  return (
    <Base {...p}>
      <path d="M12 21c4-4.2 6-7.4 6-10.2A6 6 0 0 0 6 10.8C6 13.6 8 16.8 12 21Z" />
      <circle cx="12" cy="10.6" r="2.3" />
    </Base>
  )
}

export function IconQuote({ size = 40, className, ...rest }: IconProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 40 40"
      fill="currentColor"
      className={className}
      aria-hidden="true"
      focusable="false"
      {...rest}
    >
      <path d="M6 22c0-6 3.4-10.4 9-12l1 2.6c-3.2 1.2-4.9 3-5.1 5.4H16v8H6v-4Zm16 0c0-6 3.4-10.4 9-12l1 2.6c-3.2 1.2-4.9 3-5.1 5.4H32v8H22v-4Z" />
    </svg>
  )
}
