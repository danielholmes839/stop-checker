import {
  LocationInput,
  StopRouteExploreFragment,
  TravelLegInput,
  useStopExploreQuery,
  useStopExploreWalkQuery,
  useTravelRouteQuery,
} from "client/types";
import { Container, Sign } from "components";
import { formatDistance, formatDistanceShort } from "helper";
import { Search, StopPreviewActions } from "pages/search";
import { encodeRoute } from "providers";
import React, { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import {
  Instruction,
  InstructionSubtitle,
  InstructionTitle,
} from "./instructions";

const activeCss =
  "py-1 px-2 bg-primary-100 rounded-sm mr-2 mb-1 border border-primary-200";
const normalCss =
  "py-1 px-2 bg-gray-50 rounded-sm mr-2 mb-1 border border-gray-200 hover:bg-gray-100";

type LegContextValue = {
  current: string; // stop id
  legs: TravelLegInput[];
  add: (leg: TravelLegInput) => void;
  back: () => void;
  taken: (route: string) => boolean;
};

const LegContext = React.createContext<LegContextValue>({
  current: "",
  legs: [],
  taken: (route) => false,
  add: (leg) => {},
  back: () => {},
});

const useLegContext = () => React.useContext(LegContext);

const LegContextProvider: React.FC<{ origin: string }> = ({
  origin,
  children,
}) => {
  const [current, setCurrent] = useState(origin);
  const [legs, setLegs] = useState<TravelLegInput[]>([]);

  const back = () => {
    if (legs.length === 0) {
      return;
    }

    let legsCopy = legs.map((leg) => leg);
    let leg = legsCopy.pop();

    setCurrent((leg as TravelLegInput).origin);
    setLegs(legsCopy);
  };

  const taken = (route: string): boolean => {
    for (let leg of legs) {
      if (route === leg.route) {
        return true;
      }
    }
    return false;
  };

  const add = (leg: TravelLegInput) => {
    if (legs.length > 0) {
      let lastLeg = legs[legs.length - 1];

      if (leg.route === null && lastLeg.route === null) {
        let legsCopy = legs.map((leg) => leg);
        legsCopy[legs.length - 1].destination = leg.destination;
        setCurrent(leg.destination);
        setLegs(legsCopy);
        return;
      }
    }

    setCurrent(leg.destination);
    setLegs([...legs, leg]);
  };

  return (
    <LegContext.Provider
      value={{
        add: add,
        back: back,
        taken: taken,
        legs: legs,
        current: current,
      }}
    >
      {children}
    </LegContext.Provider>
  );
};

const Remove: React.FC = () => {
  const { back } = useLegContext();
  return (
    <button className="text-red-600 text-sm mt-2" onClick={back}>
      Remove from Route
    </button>
  );
};
const Current: React.FC = () => {
  const { legs } = useLegContext();

  const [{ data }, fetch] = useTravelRouteQuery({
    variables: {
      input: legs,
    },
    pause: true,
  });

  useEffect(() => {
    if (legs.length > 0) {
      fetch();
    }
  }, [legs, fetch]);

  if (!data || !data.travelRoute.route || legs.length === 0) {
    return <></>;
  }

  return (
    <div className="mb-3">
      {data.travelRoute.route.map(
        ({ origin, destination, stopRoute, distance }, i) => {
          let isLast = i === legs.length - 1;

          if (stopRoute) {
            let route = stopRoute.route;
            return (
              <Instruction key={origin.id + destination.id}>
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
                  <InstructionSubtitle>
                    Towards {stopRoute.headsign}
                  </InstructionSubtitle>
                  <div className="mt-2">
                    <p className="text-sm">
                      Exit at{" "}
                      <span className="font-semibold text-gray-700">
                        {destination.name}
                      </span>
                    </p>
                  </div>
                  {isLast && <Remove />}
                </div>
              </Instruction>
            );
          }

          return (
            <Instruction key={origin.id + destination.id}>
              <div className="border-l-4 pl-3 border-gray-300 border-dashed">
                <InstructionTitle>
                  <span className="align-text-bottom">
                    Walk to {destination.name}
                  </span>
                </InstructionTitle>
                <InstructionSubtitle>
                  {formatDistance(distance)}
                </InstructionSubtitle>
                {isLast && <Remove />}
              </div>
            </Instruction>
          );
        }
      )}
      <div className="flex">
        <Link
          className="border border-primary-500 py-1 text-center py-0 mt-2 hover:bg-primary-500 hover:text-white text-primary-500 text-sm rounded-sm w-full"
          to={`/travel/r/${encodeRoute(legs)}`}
        >
          Done
        </Link>
      </div>
    </div>
  );
};

const WalkSelect: React.FC<{ location: LocationInput; origin: string }> = ({
  location,
  origin,
}) => {
  const { add, current } = useLegContext();
  const { data } = useStopExploreWalkQuery({
    variables: {
      location: location,
    },
  })[0];

  if (!data) {
    return <></>;
  }
  return (
    <div className="mt-3">
      {data.searchStopLocation.results.map(({ id, name, code, location }) => {
        if (id === current) {
          return <></>;
        }
        return (
          <button
            key={id}
            onClick={() =>
              add({
                origin: origin,
                destination: id,
                route: null,
              })
            }
            className={`text-xs ${normalCss}`}
          >
            {name} #{code} ({formatDistanceShort(location.distance)})
          </button>
        );
      })}
    </div>
  );
};

const Select: React.FC = () => {
  const { current, add, taken } = useLegContext();

  const { data, error } = useStopExploreQuery({
    variables: {
      origin: current,
    },
  })[0];

  const [transit, setTransit] = useState(true);
  const [activeStopRoute, setActiveStopRoute] =
    useState<StopRouteExploreFragment | null>(null);

  const update = (stopRoute: StopRouteExploreFragment | null) => {
    setTransit(stopRoute !== null);
    setActiveStopRoute(stopRoute);
  };

  useEffect(() => {
    if (!data || !data.stop) {
      return;
    }
    for (let stopRoute of data.stop.routes) {
      if (stopRoute.destinations.length > 0 && !taken(stopRoute.route.id)) {
        update(stopRoute);
        break;
      }
      update(null);
    }
  }, [data, taken]);

  if (!data || !data.stop || error) {
    return <></>;
  }

  const { stop } = data;

  return (
    <div>
      <h1 className="font-semibold">
        {stop.name} #{stop.code}
      </h1>
      <p className="text-sm">
        Choose your destination from{" "}
        <span className="underline">{stop.name}</span>
      </p>
      <div className="mt-2">
        <button
          className={transit ? normalCss : activeCss}
          onClick={() => update(null)}
        >
          <span className="text-sm">Walk</span>
        </button>
        {stop.routes.map((stopRoute, i) => {
          let { route, headsign, destinations } = stopRoute;
          if (destinations.length === 0 || taken(route.id)) {
            return <></>;
          }

          return (
            <button
              key={i}
              className={
                activeStopRoute && route.id === activeStopRoute.route.id
                  ? activeCss
                  : normalCss
              }
              onClick={() => update(stopRoute)}
            >
              <span className="text-xs">
                <Sign props={route} />
              </span>{" "}
              <span className="text-sm">{headsign}</span>
            </button>
          );
        })}
      </div>
      <div>
        {transit && activeStopRoute && (
          <div className="mt-3">
            {activeStopRoute.destinations.map((destination, i) => {
              return (
                <button
                  key={destination.id}
                  onClick={() => {
                    add({
                      origin: current,
                      destination: destination.id,
                      route: activeStopRoute.route.id,
                    });
                    update(null);
                  }}
                  className={`text-xs ${normalCss}`}
                >
                  {i + 1}. {destination.name} #{destination.code}
                </button>
              );
            })}
          </div>
        )}
        {!transit && !activeStopRoute && (
          <WalkSelect location={data.stop.location} origin={current} />
        )}
      </div>
    </div>
  );
};

const Actions: StopPreviewActions = ({ stop }) => {
  return (
    <div className="mt-3">
      <Link to={`/travel/m/${stop.id}`}>
        <button
          className="mr-3 text-primary-500 underline text-sm"
          onClick={() => {}}
        >
          Set as Origin
        </button>
      </Link>
    </div>
  );
};

export const ManualOriginInput: React.FC = () => {
  return (
    <Container>
      <div className="my-3">
        <h1 className="text-3xl font-semibold">Select an Origin</h1>
      </div>
      <Search
        config={{
          Actions: Actions,
          enableMap: true,
          enableStopRouteLinks: false,
        }}
      />
    </Container>
  );
};

export const ManualLegInput: React.FC = () => {
  const { origin } = useParams();
  return (
    <Container>
      <div className="my-3">
        <h1 className="text-3xl font-semibold">Route Planner</h1>
      </div>
      <LegContextProvider origin={origin ? origin : ""}>
        <Current />
        <Select />
      </LegContextProvider>
    </Container>
  );
};
