import React, { useCallback, useEffect, useMemo, useState } from "react";
import { Container } from "./util";

import usePlacesService from "react-google-autocomplete/lib/usePlacesAutocompleteService";
import { formatDistance } from "helper";

const storageContext = React.createContext({});

const Marker: React.FC = () => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 20 20"
      fill="currentColor"
      className="w-6 h-6 inline-block align-middle"
    >
      <path
        fillRule="evenodd"
        d="M9.69 18.933l.003.001C9.89 19.02 10 19 10 19s.11.02.308-.066l.002-.001.006-.003.018-.008a5.741 5.741 0 00.281-.14c.186-.096.446-.24.757-.433.62-.384 1.445-.966 2.274-1.765C15.302 14.988 17 12.493 17 9A7 7 0 103 9c0 3.492 1.698 5.988 3.355 7.584a13.731 13.731 0 002.273 1.765 11.842 11.842 0 00.976.544l.062.029.018.008.006.003zM10 11.25a2.25 2.25 0 100-4.5 2.25 2.25 0 000 4.5z"
        clipRule="evenodd"
      />
    </svg>
  );
};

type Location = {
  title: string;
  description: string;
  position: Position;
};

type Position = {
  latitude: number;
  longitude: number;
};

type CurrentPositionHook = {
  data: {
    position: Position | null;
    error: string | null;
  };
  request: () => void;
  reset: () => void;
};

type PlacePrediction = {
  title: string;
  description: string;
  distance: number | undefined;
};

type PlaceAutoCompleteParams = {
  position: Position | undefined;
};

type PlaceAutoCompleteHook = {
  predictions: {
    loading: boolean;
    data: PlacePrediction[];
  };
  search: (text: string) => void;
};

const usePlaceAutoComplete = ({
  position,
}: PlaceAutoCompleteParams): PlaceAutoCompleteHook => {
  const { placePredictions, getPlacePredictions, isPlacePredictionsLoading } =
    usePlacesService({
      debounce: 200,
    });

  const data: PlacePrediction[] = useMemo(() => {
    return placePredictions.map((item) => {
      return {
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
        location: new google.maps.LatLng(45.419003, -75.698142),
        radius: 75000,
        origin: position
          ? new google.maps.LatLng(position.latitude, position.longitude)
          : undefined,
      });
    },
  };
};

const useCurrentPosition = (): CurrentPositionHook => {
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

export const Travel: React.FC = () => {
  const {
    placesService,
    placePredictions,
    getPlacePredictions,
    isPlacePredictionsLoading,
  } = usePlacesService({
    debounce: 200,
  });

  const {
    data: { position, error: positionError },
    request,
  } = useCurrentPosition();

  const [location, setLocation] = useState<Location | null>(null);
  const [placeId, setPlaceId] = useState<string | null>(null);
  // const [place, setPlace] = useState<any>(null);

  useEffect(() => {
    if (position === null) {
      return;
    }

    setLocation({
      title: "Current Location",
      description: "Your location from 6:45pm",
      position: position,
    });
  }, [position]);

  useEffect(() => {
    // fetch place details for the first element in placePredictions array
    if (placeId !== null) {
      placesService?.getDetails(
        {
          placeId: placeId,
          fields: ["name", "formatted_address", "geometry"],
        },
        (placeDetails) => {
          if (!placeDetails) {
            return;
          }

          setLocation({
            title: placeDetails.name as string,
            description: placeDetails.formatted_address as string,
            position: {
              latitude: placeDetails.geometry?.location?.lat() as number,
              longitude: placeDetails.geometry?.location?.lng() as number,
            },
          });
        }
      );
    }
  }, [placeId]);

  return (
    <Container>
      {location && (
        <div className="py-3 px-3 mt-5 bg-gray-50 rounded border-b">
          <div className="inline-block align-middle text-4xl font-bold mr-2">
            A
          </div>
          <div
            className="pl-2 border-l border-gray-300 inline-block align-middle"
            style={{ maxWidth: "85%" }}
          >
            <h2>{location.title}</h2>
            <span className="text-xs">{location.description}</span>
          </div>
        </div>
      )}

      <h2 className="text-2xl font-bold mt-5">Search</h2>
      <input
        className="my-1 bg-gray-50 border-b rounded w-full p-3 px-5 focus:outline-none focus:border-b focus:border-gray-200 focus:border-0 focus:shadow"
        placeholder="Search"
        onChange={(e) => {
          getPlacePredictions({
            input: e.target.value,
            location: new google.maps.LatLng(45.419003, -75.698142),
            radius: 75000,
            origin:
              position !== null
                ? new google.maps.LatLng(position.latitude, position.longitude)
                : undefined,
          });
        }}
      />

      <button className="text-primary-500 mr-5 text-sm" onClick={request}>
        Current Location
      </button>

      <button className="text-primary-500 text-sm">Saved Locations</button>

      {!isPlacePredictionsLoading && (
        <div>
          {placePredictions.map((item, i) => (
            <div
              key={i}
              className="py-3 px-3 mt-2 bg-gray-50 rounded border-b hover:bg-primary-100 cursor-pointer"
              onClick={() => {
                setPlaceId(item.place_id);
              }}
            >
              <h1>
                <Marker />{" "}
                <span className="inline-block align-middle">
                  {item.structured_formatting.main_text}
                </span>
              </h1>
              <p className="text-xs mt-1">
                {item.distance_meters && (
                  <span className="font-semibold">
                    {formatDistance(item.distance_meters)}.
                  </span>
                )}{" "}
                {item.structured_formatting.secondary_text}
              </p>
            </div>
          ))}
        </div>
      )}
    </Container>
  );
};
