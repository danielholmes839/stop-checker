import React, { useEffect } from "react";
import { Container } from "components/util";
import { useNavigate, useParams } from "react-router-dom";

import { TravelLocation, usePlace } from "core";
import { PlaceIcon, TravelLocationInput } from "components/travel";
import { ScheduleMode, useTravelPlannerQuery } from "client/types";
import { Instructions } from "./instructions";

export const TravelLocationDisplay: React.FC<{
  travelLocation: TravelLocation | null;
  symbol: string;
}> = ({ travelLocation, symbol }) => {
  return (
    <div className="px-3 py-2 bg-gray-50 rounded border-b">
      <div className="inline-block align-middle text-4xl font-bold mr-2">
        <PlaceIcon placeId={travelLocation ? travelLocation.id : null} />
      </div>
      <div
        className="pl-2 border-l border-gray-300 inline-block align-middle"
        style={{ maxWidth: "90%" }}
      >
        {travelLocation ? (
          <div>
            <h2>{travelLocation.title}</h2>
            <span className="text-xs">{travelLocation.description}</span>
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

export const TravelDestinationInput: React.FC = () => {
  const nav = useNavigate();
  const onTravelLocationChange = (location: TravelLocation) => {
    nav(`/travel/p/${location.id}`);
  };

  return (
    <Container>
      <h1 className="text-3xl font-bold font mt-3">Travel Planner</h1>
      {/* <div className="mt-2">
        <TravelLocationDisplay travelLocation={travelLocation} symbol={"A"} />
      </div> */}
      <div className="mt-1">
        <h2 className="text-xl mt-2">Where do you want to go?</h2>
        <div className="mt-1">
          <TravelLocationInput
            setTravelLocation={onTravelLocationChange}
            suggestCurrentLocation={false}
          />
        </div>
      </div>
    </Container>
  );
};

export const TravelOriginInput: React.FC = () => {
  const nav = useNavigate();
  const { destinationId } = useParams();
  const destination = usePlace(destinationId ? destinationId : null);

  const onTravelLocationChange = (location: TravelLocation) => {
    nav(`/travel/p/${destinationId}/${location.id}`);
  };

  if (destination === null) {
    return <></>;
  }

  return (
    <Container>
      <h1 className="text-3xl font-bold font mt-3">Travel Planner</h1>
      <div className="mt-2">
        <TravelLocationDisplay travelLocation={destination} symbol={"D"} />
      </div>
      <div className="mt-1">
        <h2 className="text-xl mt-2">Where are you starting from?</h2>
        <div className="mt-1">
          <TravelLocationInput
            setTravelLocation={onTravelLocationChange}
            suggestCurrentLocation={true}
            suggestionFilter={(suggestion) => suggestion.id !== destination.id}
          />
        </div>
      </div>
    </Container>
  );
};

export const TravelSchedule: React.FC = () => {
  const { originId, destinationId } = useParams();
  const origin = usePlace(originId ? originId : null);
  const destination = usePlace(destinationId ? destinationId : null);

  console.log(originId, origin, destination, destinationId);
  if (origin === null || destination === null) {
    return <></>;
  }

  return <TravelScheduleQuery origin={origin} destination={destination} />;
};

export const TravelScheduleQuery: React.FC<{
  origin: TravelLocation;
  destination: TravelLocation;
}> = ({ origin, destination }) => {
  const {
    data,
    fetching,
    error: err,
  } = useTravelPlannerQuery({
    variables: {
      origin: origin.position,
      destination: destination.position,
      options: {
        mode: ScheduleMode.DepartAt,
        datetime: null,
      },
    },
  })[0];

  if (fetching) {
    return <></>;
  }

  if (err) {
    return (
      <>
        {err.name}: {err.message}
      </>
    );
  }

  if (data) {
    return (
      <Instructions
        origin={origin}
        destination={destination}
        payload={data.travelPlanner}
      />
    );
  }

  return <></>;
};
