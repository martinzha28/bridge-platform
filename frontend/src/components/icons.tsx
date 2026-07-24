import type { ReactNode, SVGProps } from "react";

type IconProps = SVGProps<SVGSVGElement>;

/** Minimal hairline glyphs — 24-box, 1.6 stroke, currentColor. */
function Svg({ children, ...rest }: IconProps) {
  return (
    <svg
      viewBox="0 0 24 24"
      width={20}
      height={20}
      fill="none"
      stroke="currentColor"
      strokeWidth={1.6}
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
      {...rest}
    >
      {children}
    </svg>
  );
}

export function SpadeIcon(p: IconProps) {
  return (
    <Svg {...p}>
      <path d="M12 3c3 4 7 6 7 10a4 4 0 0 1-6.5 3.1c.3 1.6.8 2.6 1.5 3.4h-4c.7-.8 1.2-1.8 1.5-3.4A4 4 0 0 1 5 13c0-4 4-6 7-10Z" />
    </Svg>
  );
}

const PATHS: Record<string, ReactNode> = {
  home: (
    <>
      <path d="M4 11 12 4l8 7" />
      <path d="M6 10v9h12v-9" />
    </>
  ),
  play: <path d="M8 5v14l11-7-11-7Z" />,
  history: (
    <>
      <path d="M4 12a8 8 0 1 0 3-6.2" />
      <path d="M4 4v4h4" />
      <path d="M12 8v4l3 2" />
    </>
  ),
  puzzles: (
    <path d="M9 4h6v3a2 2 0 0 0 4 0v0h1v6h-3a2 2 0 0 0 0 4h3v3H9v-3a2 2 0 0 0-4 0v0H4V8h3a2 2 0 0 0 0-4v0Z" />
  ),
  learn: (
    <>
      <path d="M12 4 2 9l10 5 10-5-10-5Z" />
      <path d="M6 11v5c0 1.5 2.7 3 6 3s6-1.5 6-3v-5" />
    </>
  ),
  watch: (
    <>
      <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7S2 12 2 12Z" />
      <circle cx="12" cy="12" r="3" />
    </>
  ),
  news: (
    <>
      <path d="M4 5h13v14H5a2 2 0 0 1-2-2V7" />
      <path d="M17 9h3v8a2 2 0 0 1-4 0" />
      <path d="M7 9h6M7 13h6M7 17h4" />
    </>
  ),
  friends: (
    <>
      <circle cx="9" cy="8" r="3" />
      <path d="M3 20c0-3.3 2.7-6 6-6s6 2.7 6 6" />
      <path d="M16 6a3 3 0 0 1 0 6M21 20c0-2.5-1.5-4.7-3.6-5.6" />
    </>
  ),
  profile: (
    <>
      <circle cx="12" cy="8" r="4" />
      <path d="M4 21c0-4.4 3.6-8 8-8s8 3.6 8 8" />
    </>
  ),
  settings: (
    <>
      <circle cx="12" cy="12" r="3" />
      <path d="M12 2v3M12 19v3M4.2 4.2l2.1 2.1M17.7 17.7l2.1 2.1M2 12h3M19 12h3M4.2 19.8l2.1-2.1M17.7 6.3l2.1-2.1" />
    </>
  ),
};

export function RailIcon({ name, ...rest }: IconProps & { name: string }) {
  return <Svg {...rest}>{PATHS[name] ?? null}</Svg>;
}
