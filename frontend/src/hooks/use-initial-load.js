import { useEffect, useEffectEvent } from "react";

// Runs the latest loader once on mount without forcing every screen to duplicate
// effect wiring or accidentally refetch because an inline callback changed identity.
export function useInitialLoad(loader) {
  const runLoader = useEffectEvent(loader);

  useEffect(() => {
    void runLoader();
  }, []);
}
