import React, { useState } from "react";
import { Container } from "components/util";
import { useNavigate, useParams } from "react-router-dom";

import { TravelLocation, usePlace } from "core";
import { PlaceIcon, TravelLocationInput } from "components/travel";
import { ScheduleMode, useTravelPlannerQuery } from "client/types";
import { Instructions } from "./instructions";
import { formatDateTime } from "helper";
import DateTimePicker from "react-datetime-picker";

export const TravelLocationDisplay: React.FC<{
  travelLocation: TravelLocation | null;
  prefix: string;
}> = ({ travelLocation, prefix }) => {
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
            <h2>
              {prefix} - {travelLocation.title}
            </h2>
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
        <TravelLocationDisplay
          prefix={"Destination"}
          travelLocation={destination}
        />
      </div>
      <div className="mt-1">
        <h2 className="text-2xl mt-5 font-semibold text-gray-800">
          Where are you starting from?
        </h2>
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

  if (origin === null || destination === null) {
    return <></>;
  }

  return (
    <div>
      <TravelScheduleQuery origin={origin} destination={destination} />
    </div>
  );
};

export const TravelScheduleQuery: React.FC<{
  origin: TravelLocation;
  destination: TravelLocation;
}> = ({ origin, destination }) => {
  const [mode, setMode] = useState(ScheduleMode.DepartAt);
  const [date, setDate] = useState<Date>(new Date());
  const {
    data,
    fetching,
    error: err,
  } = useTravelPlannerQuery({
    variables: {
      origin: origin.position,
      destination: destination.position,
      options: {
        mode: mode,
        datetime: date ? formatDateTime(date) : null,
      },
    },
  })[0];

  if (err) {
    return (
      <>
        {err.name}: {err.message}
      </>
    );
  }

  return (
    <Container>
      <h1 className="text-3xl font-bold mt-3">Travel Planner</h1>
      <div className="mt-2">
        <TravelLocationDisplay prefix={"Origin"} travelLocation={origin} />
      </div>
      <div className="mt-2">
        <TravelLocationDisplay
          prefix={"Destination"}
          travelLocation={destination}
        />
      </div>
      <div className="mt-2">
        <div className="inline-block">
          <button
            className={
              mode === ScheduleMode.DepartAt
                ? "text-primary-700 bg-primary-100 hover:bg-primary-200 px-2 py-1 rounded text-sm mr-2"
                : "text-gray-800 bg-gray-100 hover:bg-primary-100 hover:text-primary-700 px-2 py-1 rounded text-sm mr-2"
            }
            onClick={() => setMode(ScheduleMode.DepartAt)}
          >
            Depart At
          </button>
          <button
            className={
              mode === ScheduleMode.ArriveBy
                ? "text-primary-700 bg-primary-100 px-2 py-1 rounded text-sm mr-2"
                : "text-gray-800 bg-gray-100 hover:bg-primary-100 hover:text-primary-700 px-2 py-1 rounded text-sm mr-2"
            }
            onClick={() => setMode(ScheduleMode.ArriveBy)}
          >
            Arrive By
          </button>
        </div>
        <div className="inline-block">
          <DateTimePicker
            onChange={(date: Date | null) => {
              if (date) {
                setDate(date);
              }
            }}
            disableClock={true}
            required={true}
            value={date}
            clearIcon={null}
            calendarIcon={
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="currentColor"
                className="w-6 h-6"
              >
                <path d="M12.75 12.75a.75.75 0 11-1.5 0 .75.75 0 011.5 0zM7.5 15.75a.75.75 0 100-1.5.75.75 0 000 1.5zM8.25 17.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0zM9.75 15.75a.75.75 0 100-1.5.75.75 0 000 1.5zM10.5 17.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0zM12 15.75a.75.75 0 100-1.5.75.75 0 000 1.5zM12.75 17.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0zM14.25 15.75a.75.75 0 100-1.5.75.75 0 000 1.5zM15 17.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0zM16.5 15.75a.75.75 0 100-1.5.75.75 0 000 1.5zM15 12.75a.75.75 0 11-1.5 0 .75.75 0 011.5 0zM16.5 13.5a.75.75 0 100-1.5.75.75 0 000 1.5z" />
                <path
                  fillRule="evenodd"
                  d="M6.75 2.25A.75.75 0 017.5 3v1.5h9V3A.75.75 0 0118 3v1.5h.75a3 3 0 013 3v11.25a3 3 0 01-3 3H5.25a3 3 0 01-3-3V7.5a3 3 0 013-3H6V3a.75.75 0 01.75-.75zm13.5 9a1.5 1.5 0 00-1.5-1.5H5.25a1.5 1.5 0 00-1.5 1.5v7.5a1.5 1.5 0 001.5 1.5h13.5a1.5 1.5 0 001.5-1.5v-7.5z"
                  clipRule="evenodd"
                />
              </svg>
            }
            className="text-sm"
            autoFocus={false}
          />
        </div>
      </div>
      <div className="mt-2 pt-2 border-t">
        {data && (
          <>
            <h1 className="text-2xl mt-2 mb-2">Instructions</h1>
            <Instructions
              origin={origin}
              destination={destination}
              payload={data.travelPlanner}
            />
          </>
        )}
      </div>
    </Container>
  );
};
