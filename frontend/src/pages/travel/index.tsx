import { Container } from "components";
import React from "react";
import { ScheduleMode, useTravelPlannerQuery } from "client/types";
import { Manual } from "./manual";
import { Instructions } from "./instructions";
import { Automatic } from "./automatic";

const TravelPlanner: React.FC = () => {
  const [{ data }, _] = useTravelPlannerQuery({
    variables: {
      origin: "AK151",
      destination: "CD998",
      options: {
        datetime: "2022-08-31T11:03:00Z",
        mode: ScheduleMode.ArriveBy,
      },
    },
  });

  if (!data) {
    return <></>;
  }

  return <Instructions data={data.travelPlanner} />;
};

export const TravelPage: React.FC = () => {
  return (
    <Container>
      <Automatic />
      {/* <Manual origin="AK151" />
      <TravelPlanner /> */}
    </Container>
  );
};
