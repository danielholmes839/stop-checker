import { LiveDataFragment, useStopRouteQuery } from "client/types";
import { Container, Sign } from "components";
import { currentDate, formatDistance } from "helper";
import { format } from "path";
import React from "react";
import { useParams } from "react-router-dom";

export const StopRoutePage: React.FC = () => {
  const { stop: stopId, route: routeId } = useParams();
  const [{ data }, _] = useStopRouteQuery({
    variables: {
      stop: stopId as string,
      route: routeId as string,
      today: currentDate(),
      tomorrow: currentDate(1)
    },
  });

  if (!data || !data.stopRoute) {
    return <></>;
  }

  const { stop, route, headsign, liveMap, liveBuses, schedule } =
    data.stopRoute;

  return (
    <Container>
      <h1 className="mb-1 mt-3 text-3xl">
        <span className="text-xl">
          <Sign props={route} />{" "}
        </span>{" "}
        {headsign}
      </h1>
      <h2 className="text-sm font-semibold text-gray-700">
        {stop.name} #{stop.code}
      </h2>

      <div>
        {liveBuses.length > 0 && (
          <div className="mt-3">
            <p className="mb-3 text-sm border-l-2 pl-1 border-primary-500">
              Estimated arrival times and GPS data from OC Transpo updated every
              30 seconds
            </p>
            {liveBuses.map((bus, i) => {
              return (
                <div key={i} className="mb-3 pb-3 border-b">
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
            {liveMap && <img className="mt-3 shadow" src={liveMap} />}
          </div>
        )}
      </div>

      <div>
        <h1 className="mb-1 text-2xl">Schedule</h1>
        <div className="pb-10">
          <h2>Today</h2>
          <p className="text-xs">
            {schedule.today.map(({ stoptime }) => (
              <span key={stoptime.id} className="mr-3 inline-block">
                {stoptime.time}
              </span>
            ))}
          </p>
        </div>
      </div>
    </Container>
  );
};
