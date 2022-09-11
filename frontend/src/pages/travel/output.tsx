import { useTravelPlannerQuery } from "client/types";
import { Container } from "components";
import { formatDateTime } from "helper";
import { OptionInput, OptionProvider, useOptions } from "providers";
import React from "react";
import { useParams } from "react-router-dom";
import { Instructions } from "./instructions";

const Display: React.FC = () => {
  const { origin, destination } = useParams();
  const { options } = useOptions();

  const [{ data, fetching }] = useTravelPlannerQuery({
    variables: {
      origin: origin ? origin : "",
      destination: destination ? destination : "",
      options: {
        mode: options.mode,
        datetime: formatDateTime(options.datetime),
      },
    },
  });

  if (!origin || !destination) {
    return <>Invalid origin or destination</>;
  }

  if (fetching) {
    return <>Loading...</>;
  }

  if (!data) {
    return <>server error</>;
  }

  return (
    <>
      <Instructions data={data.travelPlanner} />
    </>
  );
};

export const AutomaticOutput: React.FC = () => {
  return (
    <Container>
      <div className="my-3">
        <h1 className="text-3xl font-semibold">Travel Schedule</h1>
      </div>
      <OptionProvider>
        <OptionInput />
        <Display />
      </OptionProvider>
    </Container>
  );
};
