import React from "react";
import { useEffect, useState } from "react";
import usePlacesAutocompleteService from "react-google-autocomplete/lib/usePlacesAutocompleteService";
import { useStorage } from "./storage_provider";
import { TravelLocation } from "./types";

export const usePlace = (placeId: string | null): TravelLocation | null => {
  const { getFavourite, getRecent } = useStorage();
  const [place, setPlace] = useState<TravelLocation | null>(() => {
    if (placeId === null) {
      return null;
    }

    let fav = getFavourite(placeId);
    if (fav) {
      return fav;
    }

    let recent = getRecent(placeId);
    if (recent) {
      return recent;
    }

    return null;
  });

  const { placesService: service } = usePlacesAutocompleteService({
    debounce: 200,
  });

  useEffect(() => {
    if (placeId === null) {
      setPlace(null);
      return;
    }

    let fav = getFavourite(placeId);
    if (fav) {
      setPlace(fav);
      return;
    }

    let recent = getRecent(placeId);
    if (recent) {
      setPlace(recent);
      return;
    }

    service?.getDetails(
      {
        placeId: placeId,
        fields: ["name", "formatted_address", "geometry"],
      },
      (placeDetails) => {
        if (!placeDetails) {
          return;
        }
        setPlace({
          id: placeId,
          title: placeDetails.name as string,
          description: placeDetails.formatted_address as string,
          position: {
            latitude: placeDetails.geometry?.location?.lat() as number,
            longitude: placeDetails.geometry?.location?.lng() as number,
          },
        });
      }
    );
  }, [placeId, service, getFavourite, getRecent]);

  return place;
};

// https://www.30secondsofcode.org/react/s/use-timeout
const useTimeout = (callback: any, delay: number) => {
  const savedCallback: any = React.useRef();

  React.useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  React.useEffect(() => {
    const tick = () => {
      savedCallback.current();
    };
    if (delay !== null) {
      let id = setTimeout(tick, delay);
      return () => clearTimeout(id);
    }
  }, [delay]);
};
