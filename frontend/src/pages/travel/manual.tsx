import { TravelLegInput, useStopExploreQuery } from "client/types";
import { Sign } from "components";
import React, { Children, useState } from "react";

type LegContextValue = {
  current: string; // stop id
  legs: TravelLegInput[];
  add: (leg: TravelLegInput) => void;
};

const LegContext = React.createContext<LegContextValue>({
  current: "",
  legs: [],
  add: (leg) => {},
});

const useLegContext = () => React.useContext(LegContext);

const LegContextProvider: React.FC<{ origin: string }> = ({
  origin,
  children,
}) => {
  const [current, setCurrent] = useState(origin);
  const [legs, setLegs] = useState<TravelLegInput[]>([]);

  const add = (leg: TravelLegInput) => {
    setCurrent(leg.destination);
    setLegs([...legs, leg]);
  };

  return (
    <LegContext.Provider
      value={{
        add: add,
        legs: legs,
        current: current,
      }}
    >
      {children}
    </LegContext.Provider>
  );
};

const Current: React.FC = () => {
  const { legs } = useLegContext();
  if (legs.length === 0) {
    return <></>;
  }
  return (
    <>
      {legs.map(({ origin, destination, route }) => (
        <div>
          {origin} {destination} {route ? route : "walk"}
        </div>
      ))}
    </>
  );
};

const WalkSelect: React.FC = () => {
  return <></>
}

const TransitSelect: React.FC = () => {
  return <></>
}

const Select: React.FC = () => {
  const { current, add } = useLegContext();
  const [{ data, error, fetching }, _] = useStopExploreQuery({
    variables: {
      origin: current,
    },
  });

  const [transit, setTransit] = useState(true);
  const [activeStopRouteIndex, setActiveStopRouteIndex] = useState(0);

  if (fetching) {
    return <>Loading...</>;
  }

  if (!data || !data.stop || error) {
    return <>error {error?.message}</>;
  }

  const { stop } = data;

  let activeCss =
    "py-1 px-2 bg-primary-100 rounded-sm mr-2 mb-1 border border-primary-200";
  let normalCss =
    "py-1 px-2 bg-gray-50 rounded-sm mr-2 mb-1 border border-gray-200 hover:bg-gray-100";

  return (
    <div className="my-5">
      <h1>
        {stop.name} #{stop.code}
      </h1>
      <div className="mt-1">
        <button
          className={transit ? normalCss : activeCss}
          onClick={() => {
            setActiveStopRouteIndex(-1);
            setTransit(false);
          }}
        >
          <span className="text-sm">Walk</span>
        </button>
        {stop.routes.map(({ route, headsign, destinations }, i) => {
          return (
            <button
              className={i === activeStopRouteIndex ? activeCss : normalCss}
              onClick={() => {
                setActiveStopRouteIndex(i);
                setTransit(true);
              }}
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
        {transit && (
          <div className="mt-3">
            {stop.routes[activeStopRouteIndex].destinations.map(
              (destination, i) => {
                return (
                  <button
                    onClick={() =>
                      add({
                        origin: current,
                        destination: destination.id,
                        route: stop.routes[activeStopRouteIndex].route.id,
                      })
                    }
                    className={`text-xs ${normalCss}`}
                  >
                    {i + 1}. {destination.name} #{destination.code}
                  </button>
                );
              }
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export const Manual: React.FC<{ origin: string }> = ({ origin }) => {
  return (
    <LegContextProvider origin={origin}>
      <Current />
      <Select />
    </LegContextProvider>
  );
};
