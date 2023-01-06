import { useEffect, useState } from "react";
import usePlacesAutocompleteService from "react-google-autocomplete/lib/usePlacesAutocompleteService";
import { useStorage } from "./storage";
import { TravelLocation } from "./types";

export const usePlace = (placeId: string | null) => {
  const { getFavourite } = useStorage();
  const [place, setPlace] = useState<TravelLocation | null>(null);
  const { placesService: service } = usePlacesAutocompleteService({
    debounce: 200,
  });
  useEffect(() => {
    if (placeId === null) {
      return;
    }

    let fav = getFavourite(placeId);
    if (fav) {
      setPlace(fav);
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
  }, [placeId, service]);

  return place;
};
