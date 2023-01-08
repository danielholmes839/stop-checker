import React, { useState } from "react";
import { Polyline } from "@react-google-maps/api";
import {
  ScheduleNodeFragment,
  SchedulePayloadFragment,
  ScheduleTransitFragment,
  ScheduleWalkFragment,
  useStopRouteDetailsQuery,
} from "client/types";
import { SimpleMap } from "components/search/map";
import { Sign } from "components/util";
import { Position, TravelLocation } from "core";
import { formatDistance, formatDistanceShort, formatTime } from "helper";

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

const TransitExtraInstruction: React.FC<{
  originId: string;
  destinationId: string;
  routeId: string;
  after: any;
}> = ({ originId, destinationId, routeId, after }) => {
  const [{ data, fetching }, refetch] = useStopRouteDetailsQuery({
    variables: {
      routeId: routeId,
      originId: originId,
      destinationId: destinationId,
      after: after,
    },
    requestPolicy: "network-only",
  });

  if (fetching) {
    return <>Loading</>;
  }

  if (!data) {
    return <></>;
  }

  if (!data.stopRoute) {
    return <></>;
  }

  let schedule = data.stopRoute.scheduleReaches;
  let buses = data.stopRoute.liveBuses.filter((bus) => {
    let busArrival = new Date(bus.arrival);
    let busArrivalLowerBound = new Date(new Date(after).valueOf() - 600_000); // 10 minutes
    let busArrivalUpperBound = new Date(new Date(after).valueOf() + 14_400_000); // 4 hours
    return (
      busArrival >= busArrivalLowerBound && busArrival <= busArrivalUpperBound
    );
  });

  return (
    <div>
      {schedule && (
        <>
          <h1 className="mt-1 font-semibold text-sm">Schedule</h1>
          <p className="text-sm">
            {schedule.next.map((res, i) => (
              <span key={i} className="mr-2">
                {formatTime(res.datetime)}
              </span>
            ))}
          </p>
        </>
      )}

      <h1 className="mt-1 font-semibold text-sm">Live Data</h1>
      <p className="text-sm">
        {buses.length > 0 ? (
          buses.map((bus, i) => {
            return (
              <span key={i} className="mr-2">
                {formatTime(bus.arrival)}{" "}
                {bus.distance && (
                  <span>({formatDistanceShort(bus.distance)})</span>
                )}
              </span>
            );
          })
        ) : (
          <span>Not Available</span>
        )}
      </p>
      <button
        className="text-primary-700 bg-primary-100 hover:bg-primary-200 px-2 mt-2 rounded text-sm"
        onClick={() =>
          refetch({
            routeId: routeId,
            originId: originId,
            destinationId: destinationId,
            after: after,
          })
        }
        disabled={fetching}
      >
        Refresh
      </button>
    </div>
  );
};

const TransitInstructions: React.FC<{
  origin: ScheduleNodeFragment;
  destination: ScheduleNodeFragment;
  transit: ScheduleTransitFragment;
}> = ({ origin, destination, transit }) => {
  const [showStopsBetween, setShowStopsBetween] = useState(false);
  const [showStopRouteDetails, setShowStopRouteDetails] = useState(false);

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
        <button
          className="text-primary-500 text-xs mt-1"
          onClick={() => setShowStopRouteDetails(!showStopRouteDetails)}
        >
          {showStopRouteDetails ? "Less" : "More"}
        </button>
        {showStopRouteDetails && (
          <div className="mt-2 pt-1 border-t">
            <TransitExtraInstruction
              originId={origin.stop.id}
              destinationId={destination.stop.id}
              routeId={transit.route.id}
              after={origin.arrival}
            />
          </div>
        )}
      </InstructionContainer>
      <InstructionContainer color={transit.route.background} dotted={false}>
        <InstructionTitle>
          Exit the <Sign props={transit.route} /> at {destination.stop.name}
        </InstructionTitle>
        <InstructionSubtitle>
          After {stopsBetween[j - i - 1].stop.name}
        </InstructionSubtitle>
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
            {stopsBetween.map((stoptime, i) => {
              return (
                <p className="text-sm mt-1" key={stoptime.id}>
                  {stoptime.time} - {stoptime.stop.name}{" "}
                  {i === 0 && <span>(Board)</span>}
                  {i === stopsBetween.length - 1 && <span>(Exit)</span>}
                </p>
              );
            })}
          </div>
        )}
      </InstructionContainer>
    </>
  );
};

const InstructionsMap: React.FC<{
  schedule: SchedulePayloadFragment;
  origin: TravelLocation;
}> = ({ schedule, origin }) => {
  return (
    <SimpleMap origin={origin.position}>
      {schedule.schedule?.legs.map((leg, i) => {
        let path: any = [];

        let options = leg.transit
          ? {
              strokeColor: leg.transit.route.background,
            }
          : {
              strokeColor: "#d1d5db",
            };

        if (leg.walk) {
          path = leg.walk.path.map(({ latitude, longitude }) => {
            return { lat: latitude, lng: longitude };
          });
        }
        if (leg.transit) {
          path = transitPath(
            leg.origin.location,
            leg.destination.location,
            leg.transit.trip.shape
          );
        }

        return (
          <div key={i}>
            <Polyline
              path={path}
              options={{
                strokeWeight: 5,
                strokeOpacity: 1,
                geodesic: true,
                clickable: false,
                strokeColor: "#000000",
              }}
            />
            <Polyline
              path={path}
              options={{
                strokeWeight: 4,
                strokeOpacity: 1,
                geodesic: true,
                clickable: false,
                ...options,
              }}
            />
          </div>
        );
      })}
    </SimpleMap>
  );
};

export const Instructions: React.FC<{
  payload: SchedulePayloadFragment;
  origin: TravelLocation;
  destination: TravelLocation;
}> = ({ payload, origin, destination }) => {
  let [openMap, setOpenMap] = useState(false);

  if (payload.error || !payload.schedule) {
    return (
      <div className="bg-red-200 text-red-800 p-3 rounded font-medium">
        Sorry, we couldn't create a travel plan for this request. Please make
        sure there's an OC Transpo bus stop within at least 1km of the origin
        and destination.
      </div>
    );
  }

  let { schedule } = payload;

  return (
    <div>
      <h2 className="text-2xl">Instructions</h2>
      <h3 className="mt-1 text-sm font-semibold text-gray-800">
        Leave by {formatTime(schedule.origin.arrival)}. Arrive at{" "}
        {formatTime(schedule.destination.arrival)} ({schedule.duration} min)
      </h3>
      <div className="mt-2">
        {openMap ? (
          <>
            <button
              className="text-primary-700 bg-primary-100 hover:bg-primary-200 px-2 mb-2 rounded text-sm"
              onClick={() => setOpenMap(!openMap)}
            >
              Hide Map
            </button>
            <InstructionsMap origin={origin} schedule={payload} />
          </>
        ) : (
          <button
            className="text-primary-700 bg-primary-100 hover:bg-primary-200 px-2 rounded text-sm"
            onClick={() => setOpenMap(!openMap)}
          >
            Open Map
          </button>
        )}
      </div>
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
      <div className="py-5"></div>
    </div>
  );
};

const transitPath = (
  origin: Position,
  destination: Position,
  shape: Position[]
): any[] => {
  let i = 0;
  let minOriginDistance = getDistanceFromLatLonInKm(
    shape[0].latitude,
    shape[0].longitude,
    origin.latitude,
    origin.longitude
  );

  let j = 0;
  let minDestinationDistance = getDistanceFromLatLonInKm(
    shape[0].latitude,
    shape[0].longitude,
    destination.latitude,
    destination.longitude
  );

  for (let k = 0; k < shape.length; k += 1) {
    let location = shape[k];
    let originDistance = getDistanceFromLatLonInKm(
      location.latitude,
      location.longitude,
      origin.latitude,
      origin.longitude
    );

    let destinationDistance = getDistanceFromLatLonInKm(
      location.latitude,
      location.longitude,
      destination.latitude,
      destination.longitude
    );

    if (originDistance < minOriginDistance) {
      minOriginDistance = originDistance;
      i = k;
    }

    if (destinationDistance < minDestinationDistance) {
      minDestinationDistance = destinationDistance;
      j = k;
    }
  }

  let path = [origin, ...shape.slice(i, j + 1), destination];

  return path.map((position) => {
    return {
      lat: position.latitude,
      lng: position.longitude,
    };
  });
};

function getDistanceFromLatLonInKm(
  lat1: number,
  lon1: number,
  lat2: number,
  lon2: number
) {
  var R = 6371; // Radius of the earth in km
  var dLat = deg2rad(lat2 - lat1); // deg2rad below
  var dLon = deg2rad(lon2 - lon1);
  var a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(deg2rad(lat1)) *
      Math.cos(deg2rad(lat2)) *
      Math.sin(dLon / 2) *
      Math.sin(dLon / 2);
  var c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  var d = R * c; // Distance in km
  return d;
}

function deg2rad(deg: number) {
  return deg * (Math.PI / 180);
}
