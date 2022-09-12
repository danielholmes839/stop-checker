import { useDashboardQuery } from "client/types";
import { Container, Sign } from "components";
import { formatDateTime, formatTime } from "helper";
import { OptionInput, OptionProvider, useOptions } from "providers";
import { encodeRoute, useStorage } from "providers/storage";
import React from "react";
import { Link } from "react-router-dom";

const DashboardResults: React.FC = () => {
  const { options } = useOptions();
  const { routes } = useStorage();

  const { data, fetching, error } = useDashboardQuery({
    variables: {
      input: routes,
      options: {
        mode: options.mode,
        datetime: formatDateTime(options.datetime),
      },
    },
  })[0];

  if (fetching) {
    return <>Loading...</>;
  }

  if (error) {
    return <>{error.message}</>;
  }

  if (!data) {
    return <>server error</>;
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-3 mt-3">
      {data.travelPlannerFixedRoutes.map(({ schedule, error }) => {
        if (!schedule) {
          return <div>{error}</div>;
        }

        let { legs, arrival, departure, duration, origin, destination } =
          schedule;

        let route = legs.map((leg) => {
          return {
            origin: leg.origin.id,
            destination: leg.destination.id,
            route: leg.transit ? leg.transit.route.id : null,
          };
        });

        return (
          <div className="py-2 px-3 border border-gray-200">
            <h1 className="font-semibold">
              {origin.name} - {destination.name}
            </h1>
            <div>
              {legs.map((leg) => {
                if (!leg.transit) {
                  return <></>;
                }
                return (
                  <span className="mr-1 text-xs">
                    <Sign props={leg.transit.route} />
                  </span>
                );
              })}
            </div>
            <p className="text-sm mt-1">
              {formatTime(departure)} - {formatTime(arrival)} ({duration} min)
            </p>
            <Link
              to={`/travel/r/${encodeRoute(route)}`}
              state={{ options: options }}
            >
              <button className="border border-primary-500 px-5 py-1 mt-2 hover:bg-primary-500 hover:text-white text-primary-500 text-sm rounded-sm">
                View
              </button>
            </Link>
          </div>
        );
      })}
    </div>
  );
};
export const Dashboard: React.FC = () => {
  const { routes } = useStorage();
  return (
    <Container>
      <div className="mt-3 mb-5">
        <h1 className="text-3xl font-semibold">Dashboard</h1>
        <p className="mt-2">
          Add routes to your dashboard using the{" "}
          <Link to="/travel">
            <span className="text-primary-500 hover:text-primary-600 hover:underline">
              Travel Planner
            </span>
          </Link>{" "}
          or enter your route{" "}
          <Link to="/travel/create">
            <span className="text-primary-500 hover:text-primary-600 hover:underline">
              Manually
            </span>
          </Link>
          .
        </p>
      </div>
      {routes.length > 0 && (
        <OptionProvider>
          <OptionInput />
          <DashboardResults />
        </OptionProvider>
      )}
    </Container>
  );
};
