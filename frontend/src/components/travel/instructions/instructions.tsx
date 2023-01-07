import React, { useState } from "react";
import { Polyline } from "@react-google-maps/api";
import {
  ScheduleNodeFragment,
  SchedulePayloadFragment,
  ScheduleTransitFragment,
  ScheduleWalkFragment,
} from "client/types";
import { SimpleMap } from "components/search/map";
import { Container, Sign } from "components/util";
import { Position, TravelLocation } from "core";
import { formatDistance, formatTime } from "helper";
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
                  {i === 0 && <span className="font-semibold">(Board)</span>}
                  {i === stopsBetween.length - 1 && (
                    <span className="font-semibold">(Exit)</span>
                  )}
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

export const EndpointLocationDisplay: React.FC<{
  travelLocation: TravelLocation | null;
  endpoint: string;
}> = ({ travelLocation, endpoint }) => {
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
            <h1 className="text-sm font-bold text-gray-800">{endpoint}</h1>
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

export const Instructions: React.FC<{
  payload: SchedulePayloadFragment;
  origin: TravelLocation;
  destination: TravelLocation;
}> = ({ payload, origin, destination }) => {
  let [openMap, setOpenMap] = useState(false);

  if (payload.error || !payload.schedule) {
    return <>{payload.error}</>;
  }

  let { schedule } = payload;

  return (
    <div>
      <h3 className="text-sm font-semibold text-gray-800">
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
