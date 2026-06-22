// Single source of truth for the responsive breakpoints. The layout is
// desktop-first: components keep their existing (wide) values and override them
// inside these max-width media queries, so the two widths never drift across
// files.
export const breakpoints = {
  /** Tablet and below. */
  tablet: 1024,
  /** Phone and below. */
  mobile: 768,
} as const;

export const mq = {
  tablet: `@media (max-width: ${breakpoints.tablet}px)`,
  mobile: `@media (max-width: ${breakpoints.mobile}px)`,
} as const;
