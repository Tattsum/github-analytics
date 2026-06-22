// Extra Recharts <XAxis> props applied only on narrow screens so dense date /
// category labels angle and auto-thin instead of overlapping. Spread these
// behind useIsNarrow(); on desktop the axis keeps the Recharts defaults so the
// non-regression baseline stays byte-identical.
export const narrowXAxisProps = {
  angle: -35,
  textAnchor: "end",
  height: 56,
  minTickGap: 8,
} as const;
