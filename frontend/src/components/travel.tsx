import React, { useEffect, useState } from "react";
import { Container } from "./util";

import { formatDistance } from "helper";
import {
  TravelLocation,
  useCurrentPosition,
  usePlace,
  usePlaceAutoComplete,
} from "core";

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

export const Travel: React.FC = () => {
  const [travelLocation, setTravelLocation] = useState<TravelLocation | null>(
    null
  );

  const { data, request: positionRequest } = useCurrentPosition();
  const position = data.position;

  const { predictions, search: predictionsRequest } = usePlaceAutoComplete({
    position: position,
  });

  const [placeId, setPlaceId] = useState<string | null>(null);
  const place = usePlace(placeId);

  useEffect(() => {
    if (position === null) {
      return;
    }

    setTravelLocation({
      id: undefined,
      title: "Current Location",
      description: "Your location from TO:DOpm",
      position: position,
    });
  }, [position]);

  useEffect(() => {
    if (place === null) {
      return;
    }
    setTravelLocation(place);
  }, [place]);

  return (
    <Container>
      {travelLocation && (
        <div className="py-3 px-3 mt-5 bg-gray-50 rounded border-b">
          <div className="inline-block align-middle text-4xl font-bold mr-2">
            A
          </div>
          <div
            className="pl-2 border-l border-gray-300 inline-block align-middle"
            style={{ maxWidth: "85%" }}
          >
            <h2>{travelLocation.title}</h2>
            <span className="text-xs">{travelLocation.description}</span>
          </div>
        </div>
      )}

      <h2 className="text-2xl font-bold mt-5">Search</h2>
      <input
        className="my-1 bg-gray-50 border-b rounded w-full p-3 px-5 focus:outline-none focus:border-b focus:border-gray-200 focus:border-0 focus:shadow"
        placeholder="Search"
        onChange={(e) => {
          predictionsRequest(e.target.value);
        }}
      />

      <button
        className="text-primary-500 mr-5 text-sm"
        onClick={positionRequest}
      >
        Current Location
      </button>

      <button className="text-primary-500 text-sm">Saved Locations</button>

      {!predictions.loading && (
        <div>
          {predictions.data.map((pred, i) => (
            <div
              key={i}
              className="py-3 px-3 mt-2 bg-gray-50 rounded border-b hover:bg-primary-100 cursor-pointer"
              onClick={() => setPlaceId(pred.placeId)}
            >
              <h1>
                <Marker />{" "}
                <span className="inline-block align-middle">{pred.title}</span>
              </h1>
              <p className="text-xs mt-1">
                {pred.distance && (
                  <span className="font-semibold">
                    {formatDistance(pred.distance)}.
                  </span>
                )}{" "}
                {pred.description}
              </p>
            </div>
          ))}
        </div>
      )}
    </Container>
  );
};
