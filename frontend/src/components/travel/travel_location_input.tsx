import {
  TravelLocation,
  usePlace,
  usePlaceAutoComplete,
  useStorage,
} from "core";
import { useEffect, useState } from "react";
import { TravelLocationResult, TravelLocationResults } from "components/travel";

const requestCurrentLocation = (setPlaceId: React.Dispatch<string>) => {
  navigator.geolocation.getCurrentPosition((position) => {
    let geocoder = new google.maps.Geocoder();
    geocoder.geocode(
      {
        location: new google.maps.LatLng({
          lat: position.coords.latitude,
          lng: position.coords.longitude,
        }),
      },
      (res) => {
        if (res === null || res.length === 0) {
          return;
        }
        if (res.length > 1) {
          setPlaceId(res[1].place_id);
        } else {
          setPlaceId(res[0].place_id);
        }
      }
    );
  });
};

type TravelLocationInputProps = {
  setTravelLocation: React.Dispatch<TravelLocation>;
};

export const TravelLocationInput: React.FC<TravelLocationInputProps> = ({
  setTravelLocation,
}) => {
  const { addHistory, favourites, history } = useStorage();
  const { predictions, search: predictionsRequest } = usePlaceAutoComplete({
    position: null,
  });

  // selected place
  const [prevPlaceId, setPrevPlaceId] = useState<string | null>(null);
  const [placeId, setPlaceId] = useState<string | null>(null);
  const { place } = usePlace(placeId);

  useEffect(() => {
    if (place === null) {
      return;
    }
    if (place.id === prevPlaceId) {
      return;
    }
    setPrevPlaceId(place.id);
    setTravelLocation(place);
    addHistory(place);
  }, [place, addHistory, setTravelLocation, prevPlaceId, setPrevPlaceId]);

  return (
    <>
      <div className="mb-1">
        <button
          className="text-primary-500 mr-5 text-sm"
          onClick={() => requestCurrentLocation(setPlaceId)}
        >
          Current Location
        </button>
      </div>
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
        favourites.length + history.length > 0 && (
          <div>
            <h2 className="mt-2 text-sm text-gray-700 font-semibold">
              Suggested
            </h2>
            {favourites.map((favourite) => (
              <TravelLocationResult
                pred={{ ...favourite, distance: undefined }}
                setPlaceId={setPlaceId}
              />
            ))}

            {history.map((recent) => (
              <TravelLocationResult
                pred={{ ...recent, distance: undefined }}
                setPlaceId={setPlaceId}
              />
            ))}
          </div>
        )}
    </>
  );
};
