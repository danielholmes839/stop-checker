import { LiveDataFragment, useStopRouteQuery } from "client/types";
import { Container, Sign } from "components";
import { formatDistance } from "helper";
import { format } from "path";
import React from "react";
import { useParams } from "react-router-dom";

export const StopRoutePage: React.FC = () => {
  const { stop: stopId, route: routeId } = useParams();
  const [{ data }, _] = useStopRouteQuery({
    variables: {
      stop: stopId as string,
      route: routeId as string,
    },
  });

  if (!data || !data.stopRoute) {
    return <></>;
  }

  const { stop, route, headsign, liveMap, liveBuses } = data.stopRoute;

  return (
    <Container>
      <h1 className="text-xl font-semibold mt-3">
        {stop.name} <span className="mx-1">/</span> <Sign props={route} />{" "}
        {headsign}
      </h1>

      <div>
        {liveMap && (
          <>
            <img className="mt-3" src={liveMap} />
          </>
        )}

        {liveBuses.length > 0 && (
          <div>
            <h1 className="mt-3 text-2xl">Live Data</h1>
            {liveBuses.map((bus, i) => {
              return (
                <div className="mt-3 pt-3 border-t">
                  <h1 className="font-semibold">
                    <span className="mr-2">Bus #{i + 1} </span>
                    <span className="text-sm">
                      <Sign props={route} />
                    </span>{" "}
                    {bus.headsign}
                  </h1>
                  <h2 className="text-xs font-semibold text-gray-700 mb-1">
                    {bus.lastUpdatedMessage}{" "}
                    {bus.distance ? `(${formatDistance(bus.distance)})` : ""}
                  </h2>
                  <p className="text-sm">Estimated arrival {bus.arrival}</p>
                </div>
              );
            })}
          </div>
        )}
      </div>

      <pre>{JSON.stringify(data, undefined, 4)}</pre>
    </Container>
  );
};
