import {
  ScheduleMode,
  TravelLegInput,
  TravelOptions,
  useTravelPlannerFixedRouteQuery,
} from "client/types";
import { Container } from "components";
import { formatDateTime } from "helper";
import {
  decodeRoute,
  OptionProvider,
  OptionInput,
  useOptions,
} from "providers";
import React from "react";
import { useLocation, useParams } from "react-router-dom";
import { Instructions } from "./instructions";

export const FixedRoute: React.FC = () => {
  const { state } = useLocation();
  const { encoded } = useParams();

  const options: TravelOptions =
    state === undefined || state === null
      ? {
          mode: ScheduleMode.ArriveBy,
          datetime: new Date(),
        }
      : (state as { options: TravelOptions }).options;

  if (encoded === undefined) {
    return <>Invalid route</>;
  }

  const route = decodeRoute(encoded);

  return (
    <Container>
      <div className="my-3">
        <h1 className="text-3xl font-semibold">Travel Plan</h1>
      </div>
      <OptionProvider initial={options}>
        <OptionInput />
        <FixedRouteResults route={route} />
      </OptionProvider>
    </Container>
  );
};

const FixedRouteResults: React.FC<{ route: TravelLegInput[] }> = ({
  route,
}) => {
  const { options } = useOptions();
  const { data, fetching, error } = useTravelPlannerFixedRouteQuery({
    variables: {
      input: route,
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

  return <Instructions data={data.travelPlannerFixedRoute} />;
};
