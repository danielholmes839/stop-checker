import {
  ScheduleNodeFragment,
  SchedulePayloadFragment,
  ScheduleTransitFragment,
  ScheduleWalkFragment,
} from "client/types";
import { Container, Sign } from "components/util";
import { TravelLocation } from "core";
import { formatDistance, formatTime } from "helper";
import React, { useState } from "react";
import { PlaceIcon } from "../place_icon";

const InstructionContainer: React.FC<{ color?: string; dotted?: boolean }> = ({
  children,
  color = "",
  dotted = true,
}) => {
  return (
    <div className="mt-3 pt-3 border-t">
      <div
        className="border-l-4 pl-3"
        style={{ borderColor: color, borderStyle: dotted ? "dashed" : "solid" }}
      >
        {children}
      </div>
    </div>
  );
};
const InstructionTitle: React.FC = ({ children }) => {
  return <h1 className="font-medium">{children}</h1>;
};

const InstructionSubtitle: React.FC = ({ children }) => {
  return <h2 className="text-xs text-gray-800 font-semibold">{children}</h2>;
};

const InstructionText: React.FC = ({ children }) => {
  return <p className="text-sm mt-1">{children}</p>;
};

const WalkInstructions: React.FC<{
  originPlace: TravelLocation;
  destinationPlace: TravelLocation;
  origin: ScheduleNodeFragment;
  destination: ScheduleNodeFragment;
  walk: ScheduleWalkFragment;
}> = ({ originPlace, origin, destinationPlace, destination, walk }) => {
  if (!walk || !walk.walk) {
    return <></>;
  }

  let originName = origin.stop ? origin.stop.name : <u>{originPlace.title}</u>;
  let destinationName = destination.stop ? (
    destination.stop.name
  ) : (
    <u>{destinationPlace.title}</u>
  );

  return (
    <InstructionContainer>
      <InstructionTitle>
        Walk from {originName} to {destinationName}
      </InstructionTitle>
      <InstructionSubtitle>
        {formatDistance(walk.walk.distance)} ({walk.duration} min)
      </InstructionSubtitle>
    </InstructionContainer>
  );
};

const TransitInstructions: React.FC<{
  origin: ScheduleNodeFragment;
  destination: ScheduleNodeFragment;
  transit: ScheduleTransitFragment;
}> = ({ origin, destination, transit }) => {
  const [showStopsBetween, setShowStopsBetween] = useState(false);
  if (!origin.stop || !destination.stop) {
    return <></>;
  }

  let i = transit.trip.stoptimes.findIndex((stoptime) => {
    return origin.stop && stoptime.stop.id === origin.stop.id;
  });

  let j = transit.trip.stoptimes.findIndex((stoptime) => {
    return destination.stop && stoptime.stop.id === destination.stop.id;
  });

  const stopsBetween = transit.trip.stoptimes.slice(i, j + 1);

  return (
    <>
      <InstructionContainer color={transit.route.background} dotted={false}>
        <InstructionTitle>
          Board the <Sign props={transit.route} /> at {origin.stop.name}
        </InstructionTitle>
        <InstructionSubtitle>
          Towards {transit.trip.headsign}
        </InstructionSubtitle>
        <InstructionText>
          Scheduled to depart at {formatTime(transit.departure)}.{" "}
          {transit.wait > 0 && <span>Wait ({transit.wait} min)</span>}
        </InstructionText>
      </InstructionContainer>
      <InstructionContainer color={transit.route.background} dotted={false}>
        <InstructionTitle>
          Exit the <Sign props={transit.route} /> at {destination.stop.name}
        </InstructionTitle>
        <InstructionSubtitle>After Stop</InstructionSubtitle>
        <InstructionText>
          Scheduled to arrive at {formatTime(destination.arrival)}
        </InstructionText>
        <button
          className="text-primary-500 text-xs mt-1"
          onClick={() => setShowStopsBetween(!showStopsBetween)}
        >
          Ride {j - i} Stops ({transit.duration} min)
        </button>
        {showStopsBetween && (
          <div className="mt-2 pt-1 border-t">
            {stopsBetween.map((stoptime) => {
              return (
                <p className="text-sm mt-1">
                  {stoptime.time} - {stoptime.stop.name}
                </p>
              );
            })}
          </div>
        )}
      </InstructionContainer>
    </>
  );
};

export const Instructions: React.FC<{
  payload: SchedulePayloadFragment;
  origin: TravelLocation;
  destination: TravelLocation;
}> = ({ payload, origin, destination }) => {
  if (payload.error || !payload.schedule) {
    return <>{payload.error}</>;
  }

  let { schedule } = payload;

  return (
    <Container>
      <h1 className="text-3xl font-bold font mt-3">Travel Schedule</h1>
      <h2 className="text-xl mt-2">
        <span className="inline-block mt-1">
          <PlaceIcon placeId={origin.id} />
          <span className="align-middle mr-2">{origin.title}</span>
        </span>
        <span className="inline-block mt-1">
          <PlaceIcon placeId={origin.id} />
          <span className="align-middle">{destination.title}</span>
        </span>
      </h2>
      <h3 className="mt-2 text-sm font-semibold text-gray-800">
        Leave by {formatTime(schedule.origin.arrival)} - Arrive at{" "}
        {formatTime(schedule.destination.arrival)} ({schedule.duration} min).
      </h3>

      {schedule.legs.map((leg, i) => {
        if (leg.transit) {
          return (
            <TransitInstructions
              key={i}
              origin={leg.origin}
              destination={leg.destination}
              transit={leg.transit}
            />
          );
        }
        return (
          <WalkInstructions
            key={i}
            destination={leg.destination}
            destinationPlace={destination}
            origin={leg.origin}
            originPlace={origin}
            walk={leg}
          />
        );
      })}
    </Container>
  );
};
