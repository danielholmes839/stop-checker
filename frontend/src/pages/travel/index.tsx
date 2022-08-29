import { Container } from "components";
import React from "react";
import { ScheduleMode, useTravelPlannerQuery } from "client/types";
import { formatTime } from "helper";
import {
  BoardInstructions,
  RideInstructions,
  WalkInstructions,
} from "./instructions";

const TravelPlanner: React.FC = () => {
  const [{ data }, _] = useTravelPlannerQuery({
    variables: {
      origin: "AK151",
      destination: "CD998",
      options: {
        datetime: "2022-08-25T11:57:00Z",
        mode: ScheduleMode.DepartAt,
      },
    },
  });

  if (!data) {
    return <></>;
  }

  const { schedule, errors } = data.travelPlanner;

  if (errors.length > 0) {
    return <pre>{JSON.stringify(errors, undefined, 4)}</pre>;
  }

  if (!schedule) {
    return <p>Failed to create a schedule...</p>;
  }

  const { arrival, departure, duration } = schedule;

  return (
    <div>
      {schedule.legs.map((leg, i) => {
        return leg.walk ? (
          <WalkInstructions key={i} leg={leg} />
        ) : (
          <div key={i}>
            <BoardInstructions leg={leg} />
            <RideInstructions leg={leg} />
          </div>
        );
      })}
      <div>
        <h1 className="font-semibold mt-3">You've reached your destination</h1>
        <h2 className="text-xs text-gray-700 font-semibold">
          Departure {formatTime(departure)} - Arrival {formatTime(arrival)} (
          {duration} min)
        </h2>
      </div>
    </div>
  );
};

export const TravelPage: React.FC = () => {
  return (
    <Container>
      <h1 className="text-4xl mt-3">Travel Planner</h1>
      {/* <Wizard
        header={<WizardHeader destination={destination} origin={origin} />}
      >
        <SelectStop name="Origin" setter={setOrigin} />
        <SelectStop name="Destination" setter={setDestination} />
        {origin && destination && (
          <div>
            You selected {origin.name} and {destination.name}
          </div>
        )}
      </Wizard> */}
      <TravelPlanner />
    </Container>
  );
};
