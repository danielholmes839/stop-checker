import { TravelScheduleLegDefaultFragment } from "client/types";
import { Sign } from "components";
import { formatDistance, formatTime } from "format";
import { useMemo, useState } from "react";

type InstructionProps = { leg: TravelScheduleLegDefaultFragment };

export const Instruction: React.FC = ({ children }) => {
  return <div className="border-b pb-3 mt-3">{children}</div>;
};

export const InstructionTitle: React.FC = ({ children }) => {
  return <h1 className="font-semibold">{children}</h1>;
};

export const InstructionSubtitle: React.FC = ({ children }) => {
  return <h2 className="text-xs text-gray-700 font-semibold">{children}</h2>;
};

export const BoardInstructions: React.FC<InstructionProps> = ({ leg }) => {
  const { origin, departure, transit } = leg;

  if (!transit) {
    return <></>;
  }

  const { route, trip } = transit;

  return (
    <Instruction>
      <div
        className="border-l-2 pl-3"
        style={{ borderColor: route.background }}
      >
        <InstructionTitle>
          Board the <Sign props={route} /> at {origin.name}
        </InstructionTitle>
        <InstructionSubtitle>Towards {trip.headsign}</InstructionSubtitle>

        <p className="text-sm mt-2">
          Scheduled to depart at {formatTime(departure)}
        </p>
        <button className="text-primary-500 text-xs">More Departures</button>
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
    let prev = transit.trip.stopTimes.filter((st) => {
      return st.sequence < transit.arrival.sequence;
    });

    // the previous stop times
    return prev[prev.length - 1];
  }, []);

  if (!transit || !prevStopTime) {
    return <></>;
  }

  const { route } = transit;

  return (
    <Instruction>
      <div
        className="border-l-2 pl-3"
        style={{ borderColor: route.background }}
      >
        <InstructionTitle>
          Exit the <Sign props={route} /> at {destination.name}
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
            {transit.trip.stopTimes
              .filter((st) => {
                let originSeq = transit.departure.sequence;
                let destinationSeq = transit.arrival.sequence;
                console.log(originSeq, destinationSeq, st.sequence);
                return (
                  st.sequence >= originSeq && st.sequence <= destinationSeq
                );
              })
              .map((st) => {
                return (
                  <p className="text-sm mt-1">
                    {st.time} - {st.stop.name}
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
          Distance of {formatDistance(distance)} ({duration} min)
        </InstructionSubtitle>
      </div>
    </Instruction>
  );
};
