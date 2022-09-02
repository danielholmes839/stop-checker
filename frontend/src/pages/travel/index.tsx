import { Container } from "components";
import React from "react";
import {
  ScheduleMode,
  TravelScheduleFragment,
  useTravelPlannerQuery,
} from "client/types";
import { formatTime } from "helper";
import {
  BoardInstructions,
  RideInstructions,
  WalkInstructions,
} from "./instructions";
import { Link } from "react-router-dom";
import { Manual } from "./manual";

export const Travel: React.FC = () => {
  return <></>;
};

const TravelSchedule: React.FC<TravelScheduleFragment> = ({ schedule }) => {
  if (!schedule) {
    return (
      <div>
        Failed to create a travel schedule. This can occur when there's no way
        to create a travel schedule within 3 days of the departure/arrival time
      </div>
    );
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

const TravelPlanner: React.FC = () => {
  const [{ data }, _] = useTravelPlannerQuery({
    variables: {
      origin: "AK151",
      destination: "CD998",
      options: {
        datetime: "2022-08-31T11:03:00Z",
        mode: ScheduleMode.ArriveBy,
      },
    },
  });

  if (!data) {
    return <>Loading...</>;
  }

  return <TravelSchedule schedule={data.travelPlanner.schedule} />;
};

export const TravelPage: React.FC = () => {
  return (
    <Container>
      <h1 className="text-4xl mt-3">Travel Planner</h1>
      {/* <div className="grid sm:grid-cols-1 md:grid-cols-2 gap-3">
        <div className="p-3 rounded-sm border">
          <h2 className="text-lg font-semibold">Automatic</h2>
          <p className="text-sm text-gray-600 mb-3">
            Let stop-checker plan your route automatically. Just enter the
            origin and destination stops.
          </p>
          <Link
            to="/travel/manual"
            className="text-white bg-primary-500 px-5 py-1 rounded-sm text-sm font-semibold tracking-wide"
          >
            Start!
          </Link>
        </div>
        <div className="p-3 rounded-sm border">
          <h2 className="text-lg font-semibold">Manual</h2>
          <p className="text-sm text-gray-600 mb-3">
            Plan your route yourself! Just enter the origin the stop.
          </p>
          <Link
            to="/travel/manual"
            className="text-gray-800 bg-gray-200 px-5 py-1 rounded-sm text-sm font-semibold tracking-wide"
          >
            Select
          </Link>
        </div>
      </div> */}
      <Manual origin="AK151" />
      <TravelPlanner />
    </Container>
  );
};
