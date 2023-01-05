import { useCallback, useState } from "react";
import { Position } from "./types";

export type CurrentPositionHook = {
  data: {
    position: Position | null;
    error: string | null;
  };
  request: () => void;
  reset: () => void;
};

export const useCurrentPosition = (): CurrentPositionHook => {
  const [location, setPosition] = useState<Position | null>(null);
  const [error, setError] = useState<string | null>(null);

  const request = useCallback(() => {
    if (!navigator.geolocation) {
      setPosition(null);
      setError("failed to retrieve your location: browser error");
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        setPosition({
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
        });
        setError(null);
      },
      () => {
        setError("failed to retrieve your location: permission denied");
      }
    );
  }, [setPosition, setError]);

  const reset = useCallback(() => {
    setPosition(null);
    setError(null);
  }, [setPosition, setError]);

  return {
    data: {
      position: location,
      error: error,
    },
    request,
    reset,
  };
};
