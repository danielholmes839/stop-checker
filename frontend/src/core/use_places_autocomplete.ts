import { useMemo } from "react";
import usePlacesAutocompleteService from "react-google-autocomplete/lib/usePlacesAutocompleteService";
import { Position } from "./types";

export type PlacePrediction = {
  placeId: string;
  title: string;
  description: string;
  distance: number | undefined;
};

export type PlaceAutoCompleteParams = {
  position: Position | null;
};

export type PlaceAutoCompleteHook = {
  predictions: {
    loading: boolean;
    data: PlacePrediction[];
  };
  search: (text: string) => void;
};

export const usePlaceAutoComplete = (
  params: PlaceAutoCompleteParams
): PlaceAutoCompleteHook => {
  const { placePredictions, getPlacePredictions, isPlacePredictionsLoading } =
    usePlacesAutocompleteService({
      debounce: 200,
    });

  const data: PlacePrediction[] = useMemo(() => {
    return placePredictions.map((item) => {
      return {
        placeId: item.place_id,
        title: item.structured_formatting.main_text,
        description: item.structured_formatting.secondary_text,
        distance: item.distance_meters,
      };
    });
  }, [placePredictions]);

  return {
    predictions: {
      data: data,
      loading: isPlacePredictionsLoading,
    },
    search: (text) => {
      getPlacePredictions({
        input: text,
        location: new google.maps.LatLng({ lat: 45.419003, lng: -75.698142 }),
        radius: 75000,
        origin:
          params.position !== null
            ? new google.maps.LatLng({
                lat: params.position.latitude,
                lng: params.position.longitude,
              })
            : undefined,
      });
    },
  };
};
