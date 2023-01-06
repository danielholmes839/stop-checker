import { useEffect, useState } from "react";
import usePlacesAutocompleteService from "react-google-autocomplete/lib/usePlacesAutocompleteService";
import { useStorage } from "./storage_provider";
import { TravelLocation } from "./types";

type UsePlaceHook = {
  place: TravelLocation | null;
  loading: boolean;
};

export const usePlace = (placeId: string | null): UsePlaceHook => {
  const { getFavourite, getRecent } = useStorage();
  const [place, setPlace] = useState<TravelLocation | null>(null);
  const [loading, setLoading] = useState(true);

  const { placesService: service } = usePlacesAutocompleteService({
    debounce: 200,
  });
  useEffect(() => {
    if (placeId === null) {
      setLoading(false);
      return;
    }

    let fav = getFavourite(placeId);
    if (fav) {
      setPlace(fav);
      setLoading(false);
      return;
    }

    let recent = getRecent(placeId);
    if (recent) {
      setPlace(recent);
      setLoading(false);
      return;
    }

    service?.getDetails(
      {
        placeId: placeId,
        fields: ["name", "formatted_address", "geometry"],
      },
      (placeDetails) => {
        setLoading(false);
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
  }, [placeId, service, getFavourite]);

  return { place, loading };
};
