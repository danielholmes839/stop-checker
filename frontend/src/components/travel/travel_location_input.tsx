import {
  TravelLocation,
  usePlace,
  usePlaceAutoComplete,
  useStorage,
} from "core";
import React, { useEffect, useState } from "react";
import {
  TravelCurrentLocationOption,
  TravelLocationResult,
  TravelLocationResults,
} from "components/travel";

type TravelLocationInputProps = {
  setTravelLocation: React.Dispatch<TravelLocation>;
  suggestCurrentLocation?: boolean;
  suggestFavourites?: boolean;
  suggestHistory?: boolean;
  suggestionFilter?: (travelLocation: TravelLocation) => boolean;
};

export const TravelLocationInput: React.FC<TravelLocationInputProps> = ({
  setTravelLocation,
  suggestCurrentLocation = true,
  suggestFavourites = true,
  suggestHistory = true,
  suggestionFilter = (location) => true,
}) => {
  const { addHistory, favourites, history } = useStorage();
  const { predictions, search: predictionsRequest } = usePlaceAutoComplete({
    position: null,
  });

  // selected place
  const [prevPlaceId, setPrevPlaceId] = useState<string | null>(null);
  const [placeId, setPlaceId] = useState<string | null>(null);
  const place = usePlace(placeId);

  useEffect(() => {
    if (place === null) {
      return;
    }
    if (place.id === prevPlaceId) {
      return;
    }
    addHistory(place);
    setPrevPlaceId(place.id);
    setTravelLocation(place);
  }, [place, addHistory, setTravelLocation, prevPlaceId, setPrevPlaceId]);

  return (
    <>
      <div>
        <input
          className="bg-gray-50 border-b rounded w-full p-3 focus:outline-none focus:border-b focus:border-gray-200 focus:border-0 focus:shadow text-lg"
          placeholder="Search"
          onChange={(e) => {
            predictionsRequest(e.target.value);
          }}
        />
      </div>
      {/* action buttons */}
      {!predictions.loading && predictions.data.length > 0 && (
        <div>
          <TravelLocationResults
            predictions={predictions.data}
            setPlaceId={setPlaceId}
          />
        </div>
      )}
      {!predictions.loading &&
        predictions.data.length === 0 &&
        ((favourites.length > 0 && suggestFavourites) ||
          (history.length > 0 && suggestHistory) ||
          suggestCurrentLocation) && (
          <div>
            <h2 className="mt-3 text-sm text-gray-700 font-semibold">
              Suggested
            </h2>
            {suggestCurrentLocation && (
              <TravelCurrentLocationOption setPlaceId={setPlaceId} />
            )}
            {suggestFavourites &&
              favourites
                .filter(suggestionFilter)
                .map((favourite) => (
                  <TravelLocationResult
                    key={favourite.id}
                    pred={{ ...favourite, distance: undefined }}
                    setPlaceId={setPlaceId}
                  />
                ))}

            {history.filter(suggestionFilter).map((recent) => (
              <TravelLocationResult
                key={recent.id}
                pred={{ ...recent, distance: undefined }}
                setPlaceId={setPlaceId}
              />
            ))}
          </div>
        )}
    </>
  );
};
