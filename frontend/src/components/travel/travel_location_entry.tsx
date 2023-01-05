import React, { useEffect, useMemo, useState } from "react";
import { Container } from "../util";

import { formatDistanceShort } from "helper";
import {
  Position,
  TravelLocation,
  useCurrentPosition,
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
      <path d="M9.653 16.915l-.005-.003-.019-.01a20.759 20.759 0 01-1.162-.682 22.045 22.045 0 01-2.582-1.9C4.045 12.733 2 10.352 2 7.5a4.5 4.5 0 018-2.828A4.5 4.5 0 0118 7.5c0 2.852-2.044 5.233-3.885 6.82a22.049 22.049 0 01-3.744 2.582l-.019.01-.005.003h-.002a.739.739 0 01-.69.001l-.002-.001z" />
    </svg>
  );
};

const useReverseGeocode = (position: Position | null): string | null => {
  const [placeId, setPlaceId] = useState<string | null>(null);
  const geocoder = useMemo(() => new google.maps.Geocoder(), []);
  useEffect(() => {
    if (position === null) {
      return;
    }
    setPlaceId(null);
    geocoder.geocode(
      {
        location: new google.maps.LatLng({
          lat: position.latitude,
          lng: position.longitude,
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
  }, [position, geocoder]);

  return placeId;
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
        style={{ maxWidth: "85%" }}
      >
        {travelLocation ? (
          <div>
            {" "}
            <h2>{travelLocation.title}</h2>
            <span className="text-xs">{travelLocation.description}</span>
            <div>
              {isFavourite(travelLocation.id) ? (
                <button
                  onClick={() => deleteFavourite(travelLocation.id)}
                  className="text-primary-500 text-sm mt-1"
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
            </div>
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
      <h1 className="text-3xl font-bold tracking-wide font mt-4">
        Travel Planner
      </h1>
      <div className="mt-2">
        <TravelLocationDisplay travelLocation={travelLocation} symbol={"A"} />
      </div>
      <div className="mt-2">
        <TravelLocationDisplay travelLocation={null} symbol={"B"} />
      </div>
      <div>
        <button onClick={clearHistory}>Clear History</button>{" "}
        <button onClick={clear}>Clear All Data</button>
      </div>
      <div className="mt-2">
        <h2 className="text-xl font">Where do you want to go?</h2>
        <div className="mt-3">
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

export const TravelLocationInput: React.FC<TravelLocationInputProps> = ({
  setTravelLocation,
}) => {
  const { addHistory } = useStorage();
  const { data, request: positionRequest } = useCurrentPosition();
  const position = data.position;

  const { predictions, search: predictionsRequest } = usePlaceAutoComplete({
    position: position,
  });

  // selected place
  const [placeId, setPlaceId] = useState<string | null>(null);
  const selectedPlace = usePlace(placeId);

  // current gps place
  const currentPlace = usePlace(useReverseGeocode(position));

  // place given to the parent
  const [inputPlace, setInputPlace] = useState<TravelLocation | null>(null);

  useEffect(() => {
    if (inputPlace === null) {
      return;
    }
    setTravelLocation(inputPlace);
    addHistory(inputPlace);
  }, [inputPlace, addHistory, setTravelLocation]);

  useEffect(() => {
    setInputPlace(selectedPlace);
  }, [selectedPlace]);

  useEffect(() => {
    if (currentPlace === null) {
      return;
    }
    setInputPlace(currentPlace);
  }, [currentPlace]);

  return (
    <>
      {/* input */}
      <div>
        <input
          className="bg-gray-50 border-b rounded w-full p-3 focus:outline-none focus:border-b focus:border-gray-200 focus:border-0 focus:shadow"
          placeholder="Search"
          onChange={(e) => {
            predictionsRequest(e.target.value);
          }}
        />
      </div>
      {/* action buttons */}
      <div className="mt-1">
        <button
          className="text-primary-500 mr-5 text-sm"
          onClick={positionRequest}
        >
          Current Location
        </button>

        <button className="text-primary-500 text-sm">Saved Locations</button>
      </div>
      {/* prediction results */}
      {!predictions.loading && (
        <div>
          {predictions.data.map((pred, i) => (
            <div
              key={i}
              className="px-3 py-2 mt-2 bg-gray-50 rounded border-b hover:bg-primary-100 cursor-pointer"
              onClick={() => {
                if (
                  selectedPlace === null ||
                  selectedPlace.id !== pred.placeId
                ) {
                  // we load the place
                  setPlaceId(pred.placeId);
                } else {
                  // we just set the place
                  setInputPlace(selectedPlace);
                }
              }}
            >
              <h1>
                <TravelLocationIcon placeId={pred.placeId} />{" "}
                <span className="inline-block align-middle">{pred.title}</span>
              </h1>
              <p className="text-xs mt-1">
                {pred.distance && (
                  <span className="font-semibold">
                    {formatDistanceShort(pred.distance)}.
                  </span>
                )}{" "}
                {pred.description}
              </p>
            </div>
          ))}
        </div>
      )}
    </>
  );
};
