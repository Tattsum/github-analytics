import { useEffect, useState } from "react";
import { breakpoints } from "../styles/breakpoints";

const QUERY = `(max-width: ${breakpoints.mobile}px)`;

// Tracks whether the viewport is at or below the mobile breakpoint. Charts use
// this to thin/angle dense X-axis labels only on narrow screens; on desktop the
// Recharts defaults are kept untouched (which is what the visual baseline pins).
export function useIsNarrow(): boolean {
  const [isNarrow, setIsNarrow] = useState(() =>
    typeof window === "undefined" ? false : window.matchMedia(QUERY).matches,
  );

  useEffect(() => {
    const mql = window.matchMedia(QUERY);
    const onChange = (event: MediaQueryListEvent) => setIsNarrow(event.matches);
    mql.addEventListener("change", onChange);
    setIsNarrow(mql.matches);
    return () => mql.removeEventListener("change", onChange);
  }, []);

  return isNarrow;
}
