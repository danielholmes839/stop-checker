import {
  TravelPlannerDeparturesQueryVariables,
  TravelScheduleFragment,
  TravelScheduleLegDefaultFragment,
  useTravelPlannerDeparturesQuery,
} from "client/types";
import { Sign } from "components";
import { formatDistance, formatTime } from "helper";
import { SimpleMap } from "pages/search/map";
import { useMemo, useState } from "react";
import { Polyline } from "@react-google-maps/api";
import { useStorage } from "providers/storage";

type InstructionProps = { leg: TravelScheduleLegDefaultFragment };

export const Instruction: React.FC = ({ children }) => {
  return <div className="border-b pb-3 mt-3">{children}</div>;
};

export const InstructionTitle: React.FC = ({ children }) => {
  return <h1 className="font-semibold">{children}</h1>;
};

export const InstructionSubtitle: React.FC = ({ children }) => {
  return <h2 className="text-xs text-gray-800 font-semibold">{children}</h2>;
};

const MoreDepartures: React.FC<{
  input: TravelPlannerDeparturesQueryVariables;
}> = ({ input }) => {
  const { data } = useTravelPlannerDeparturesQuery({
    variables: input,
  })[0];

  if (!data || !data.stopRoute) {
    return <>{JSON.stringify(data, undefined, 4)}</>;
  }

  return (
    <div className="border-t mt-1 pt-1">
      <p className="text-xs">
        {data.stopRoute.schedule &&
          data.stopRoute.schedule.next
            .filter((_, i) => i > 0)
            .map(({ stoptime }) => (
              <span key={stoptime.id} className="mr-2">
                {stoptime.time}
              </span>
            ))}
      </p>
    </div>
  );
};
export const BoardInstructions: React.FC<InstructionProps> = ({ leg }) => {
  const { origin, destination, departure, transit } = leg;
  const [showMoreDepartures, setShowMoreDepartures] = useState(false);

  if (!transit) {
    return <></>;
  }

  const { route, trip } = transit;

  return (
    <Instruction>
      <div
        className="border-l-4 pl-3"
        style={{ borderColor: route.background }}
      >
        <InstructionTitle>
          Board the{" "}
          <span className="text-sm">
            <Sign props={route} />
          </span>{" "}
          at {origin.name}
        </InstructionTitle>
        <InstructionSubtitle>Towards {trip.headsign}</InstructionSubtitle>

        <p className="text-sm mt-2">
          Scheduled to depart at {formatTime(departure)}
        </p>
        <button
          onClick={() => setShowMoreDepartures(!showMoreDepartures)}
          className="text-primary-500 text-xs"
        >
          More
        </button>
        {showMoreDepartures && (
          <MoreDepartures
            input={{
              after: departure,
              limit: 4,
              route: transit.route.id,
              origin: origin.id,
              destination: destination.id,
            }}
          />
        )}
      </div>
    </Instruction>
  );
};

export const RideInstructions: React.FC<InstructionProps> = ({ leg }) => {
  const { destination, duration, arrival, transit } = leg;
  const [showStopTimes, setShowStopTimes] = useState(false);

  let prevStopTime = useMemo(() => {
    if (!transit) {
      return;
    }

    // all previous stop times
    let prev = transit.trip.stoptimes.filter((st) => {
      return st.sequence < transit.arrival.sequence;
    });

    // the previous stop times
    return prev[prev.length - 1];
  }, [transit]);

  if (!transit || !prevStopTime) {
    return <></>;
  }

  const { route } = transit;

  return (
    <Instruction>
      <div
        className="border-l-4 pl-3"
        style={{ borderColor: route.background }}
      >
        <InstructionTitle>
          Exit the{" "}
          <span className="text-sm">
            <Sign props={route} />
          </span>{" "}
          at {destination.name}
        </InstructionTitle>
        <InstructionSubtitle>
          After {prevStopTime.stop.name}
        </InstructionSubtitle>
        <p className="text-sm mt-2">
          Scheduled to arrive at {formatTime(arrival)}
        </p>
        <button
          className="text-primary-500 text-xs"
          onClick={() => {
            setShowStopTimes(!showStopTimes);
          }}
        >
          Ride {transit.arrival.sequence - transit.departure.sequence} stops (
          {duration} min)
        </button>
        {showStopTimes && (
          <div className="border-t mt-1 pt-1">
            {transit.trip.stoptimes
              .filter((stoptime) => {
                return (
                  stoptime.sequence >= transit.departure.sequence &&
                  stoptime.sequence <= transit.arrival.sequence
                );
              })
              .map((stoptime) => {
                return (
                  <p key={stoptime.id} className="text-sm mt-1">
                    {stoptime.time} - {stoptime.stop.name}
                  </p>
                );
              })}
          </div>
        )}
      </div>
    </Instruction>
  );
};

export const WalkInstructions: React.FC<InstructionProps> = ({ leg }) => {
  const { destination, duration, distance } = leg;

  return (
    <Instruction>
      <div className="border-l-4 pl-3 border-gray-300 border-dashed">
        <InstructionTitle>
          <span className="align-text-bottom">Walk to {destination.name}</span>
        </InstructionTitle>
        <InstructionSubtitle>
          {formatDistance(distance)} ({duration} min)
        </InstructionSubtitle>
      </div>
    </Instruction>
  );
};

export const Instructions: React.FC<{
  data: TravelScheduleFragment;
}> = ({ data }) => {
  const { schedule, error } = data;
  const { has, add, remove } = useStorage();

  if (!schedule) {
    return (
      <div>
        <p>
          Failed to create a travel plan. This can occur when there's no way to
          create a travel plan within 3 days of the departure/arrival time.
          {error && <span> Error: {error}</span>}
        </p>
      </div>
    );
  }

  const { arrival, departure, duration, legs } = schedule;
  const route = legs.map((leg) => {
    return {
      origin: leg.origin.id,
      destination: leg.destination.id,
      route: leg.transit ? leg.transit.route.id : null,
    };
  });

  return (
    <div>
      <div className="mb-1 mt-3">
        <SimpleMap origin={legs[0].origin.location}>
          {legs.map((leg) => {
            let path = leg.shape.map(({ latitude, longitude }) => {
              return { lat: latitude, lng: longitude };
            });

            let options = leg.transit
              ? {
                  strokeColor: leg.transit.route.background,
                }
              : {
                  strokeColor: "#d1d5db",
                };

            console.log(path);

            return (
              <>
                <Polyline
                  path={path}
                  options={{
                    strokeWeight: 5,
                    strokeOpacity: 1,
                    geodesic: true,
                    clickable: true,
                    strokeColor: "#000000",
                  }}
                />
                <Polyline
                  path={path}
                  options={{
                    strokeWeight: 4,
                    strokeOpacity: 1,
                    geodesic: true,
                    clickable: true,
                    ...options,
                  }}
                />
              </>
            );
          })}
        </SimpleMap>
      </div>
      {legs.map((leg, i) => {
        return leg.walk ? (
          <WalkInstructions key={i} leg={leg} />
        ) : (
          <div key={i}>
            <BoardInstructions leg={leg} />
            <RideInstructions leg={leg} />
          </div>
        );
      })}
      <div className="mb-10">
        <h1 className="font-semibold mt-3">You've reached your destination</h1>
        <h2 className="text-xs text-gray-700 font-semibold">
          Departure {formatTime(departure)} - Arrival {formatTime(arrival)} (
          {duration} min)
        </h2>
        <div className="mt-2">
          {has(route) ? (
            <button
              className="text-red-600 text-sm"
              onClick={() => remove(route)}
            >
              Remove from Dashboard
            </button>
          ) : (
            <button
              className="text-primary-500 text-sm"
              onClick={() => add(route)}
            >
              Add to Dashboard
            </button>
          )}
        </div>
      </div>
    </div>
  );
};
