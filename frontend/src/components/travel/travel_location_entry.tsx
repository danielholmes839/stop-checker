import React, { useEffect, useState } from "react";
import { Container } from "components/util";

import {
  PlacePrediction,
  TravelLocation,
  usePlace,
  usePlaceAutoComplete,
  useStorage,
} from "core";

const MarkerIcon: React.FC = () => {
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

const HistoryIcon: React.FC = () => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 20 20"
      fill="currentColor"
      className="w-6 h-6 inline-block align-middle"
    >
      <path
        fillRule="evenodd"
        d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z"
        clipRule="evenodd"
      />
    </svg>
  );
};

const FavouriteIcon: React.FC = () => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 20 20"
      fill="currentColor"
      className="w-6 h-6 inline-block align-middle"
    >
      <path
        fillRule="evenodd"
        d="M10 2c-1.716 0-3.408.106-5.07.31C3.806 2.45 3 3.414 3 4.517V17.25a.75.75 0 001.075.676L10 15.082l5.925 2.844A.75.75 0 0017 17.25V4.517c0-1.103-.806-2.068-1.93-2.207A41.403 41.403 0 0010 2z"
        clipRule="evenodd"
      />
    </svg>
  );
};

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

export const TravelLocationDisplay: React.FC<{
  travelLocation: TravelLocation | null;
  symbol: string;
}> = ({ travelLocation, symbol }) => {
  const { isFavourite, deleteFavourite, addFavourite } = useStorage();
  return (
    <div className="px-3 py-2 bg-gray-50 rounded border-b">
      <div className="inline-block align-middle text-4xl font-bold mr-2">
        {symbol}
      </div>
      <div
        className="pl-2 border-l border-gray-300 inline-block align-middle"
        style={{ maxWidth: "90%" }}
      >
        {travelLocation ? (
          <div>
            {" "}
            <h2>{travelLocation.title}</h2>
            <span className="text-xs">{travelLocation.description}</span>
            {/* <div>
              {isFavourite(travelLocation.id) ? (
                <button
                  onClick={() => deleteFavourite(travelLocation.id)}
                  className="text-red-500 text-sm mt-1"
                >
                  Delete Favourite
                </button>
              ) : (
                <button
                  onClick={() => addFavourite(travelLocation)}
                  className="text-primary-500 text-sm mt-1"
                >
                  Add Favourite
                </button>
              )}
            </div> */}
          </div>
        ) : (
          <div>
            <h2>Not Selected</h2>
          </div>
        )}
      </div>
    </div>
  );
};
export const Travel: React.FC = () => {
  const [travelLocation, setTravelLocation] = useState<TravelLocation | null>(
    null
  );

  const { clear, clearHistory } = useStorage();

  return (
    <Container>
      <h1 className="text-3xl font-bold font mt-3">Travel Planner</h1>
      <div className="mt-2">
        <TravelLocationDisplay travelLocation={travelLocation} symbol={"A"} />
      </div>
      {/* <div className="mt-2">
        <TravelLocationDisplay travelLocation={null} symbol={"B"} />
      </div> */}
      <div className="mt-3">
        <h2 className="text-xl font-bold">Where do you want to go?</h2>
        <div className="mt-1">
          <TravelLocationInput setTravelLocation={setTravelLocation} />
        </div>
      </div>
    </Container>
  );
};

type TravelLocationInputProps = {
  setTravelLocation: React.Dispatch<TravelLocation>;
};

const TravelLocationIcon: React.FC<{ placeId: string }> = ({ placeId }) => {
  const { isFavourite, isHistory } = useStorage();
  if (isFavourite(placeId)) {
    return <FavouriteIcon />;
  } else if (isHistory(placeId)) {
    return <HistoryIcon />;
  }
  return <MarkerIcon />;
};

const TravelLocationResult: React.FC<{
  setPlaceId: React.Dispatch<string>;
  pred: PlacePrediction;
}> = ({ pred, setPlaceId }) => {
  return (
    <div
      key={pred.id}
      className="px-3 py-2 mt-2 bg-gray-50 rounded border-b hover:bg-primary-100 cursor-pointer"
      onClick={() => setPlaceId(pred.id)}
    >
      <div className="inline-block align-middle">
        <TravelLocationIcon placeId={pred.id} />
      </div>
      <div
        className="pl-2 inline-block align-middle"
        style={{ maxWidth: "90%" }}
      >
        <span>
          {pred.title} {pred.id}
        </span>
        <p className="text-xs mt-1">{pred.description}</p>
      </div>
    </div>
  );
};

const TravelLocationResults: React.FC<{
  setPlaceId: React.Dispatch<string>;
  predictions: PlacePrediction[];
}> = ({ setPlaceId, predictions }) => {
  const { getFavourite } = useStorage();
  return (
    <div>
      {predictions.map((p) => {
        let pred = p;
        let fav = getFavourite(pred.id);
        if (fav) {
          pred = {
            ...p,
            title: fav.title,
          };
        }
        return <TravelLocationResult pred={pred} setPlaceId={setPlaceId} />;
      })}
    </div>
  );
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
