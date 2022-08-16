import { Container, Sign } from "components";
import { Search } from "../search";
import { Wizard, useWizard } from "react-use-wizard";
import React, { useState } from "react";
import {
  ScheduleMode,
  TravelScheduleLegDefaultFragment,
  useTravelPlannerQuery,
} from "client/types";
import { formatDistance, formatTime } from "format";
import {
  BoardInstructions,
  RideInstructions,
  WalkInstructions,
} from "./instructions";

type SelectedStop = {
  id: string;
  name: string;
  code: string;
} | null;

type SelectStopProps = {
  name: string;
  setter: (stop: SelectedStop) => void;
};

const SelectStop: React.FC<SelectStopProps> = ({ name, setter }) => {
  const { nextStep } = useWizard();

  return (
    <>
      <Search
        config={{
          action: (stop) => {
            setter({ ...stop });
            nextStep();
          },
          actionName: `Select ${name}`,
        }}
      />
    </>
  );
};

type WizardHeaderProps = {
  origin: SelectedStop;
  destination: SelectedStop;
};

const WizardStopHeader: React.FC<{
  stop: SelectedStop;
  name: string;
  wizardPos: number;
}> = ({ stop, name, wizardPos }) => {
  const { goToStep, activeStep } = useWizard();
  const active = wizardPos === activeStep;

  const divCss = active
    ? "border-b-2 bg-gray-50 px-3 rounded-t-sm border-primary-500 py-1 cursor-pointer"
    : "border-b-2 bg-gray-50 px-3 rounded-t-sm border-gray-200 py-1 cursor-pointer";

  return (
    <div className={divCss} onClick={() => goToStep(wizardPos)}>
      <span className="font-semibold text-gray-700">
        {wizardPos + 1}. {name}
      </span>
      {stop && (
        <span className="text-sm ml-2">
          {stop.name} #{stop.code}
        </span>
      )}
      {activeStep === 2 && (
        <button
          className="float-right"
          onClick={() => {
            goToStep(wizardPos);
          }}
        >
          Edit
        </button>
      )}
    </div>
  );
};

const WizardHeader: React.FC<WizardHeaderProps> = ({ origin, destination }) => {
  return (
    <div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3 my-5">
        <WizardStopHeader name="Origin" stop={origin} wizardPos={0} />
        <WizardStopHeader name="Destination" stop={destination} wizardPos={1} />
      </div>
    </div>
  );
};

const TravelPlanner: React.FC = () => {
  const [{ data }, _] = useTravelPlannerQuery({
    variables: {
      origin: "AK151",
      destination: "CD998",
      datetime: "2022-08-15T10:45:00Z",
      mode: ScheduleMode.DepartAt,
    },
  });

  if (!data) {
    return <></>;
  }

  const { travelRoute, errors } = data.travelPlanner;

  if (errors.length > 0) {
    return <pre>{JSON.stringify(errors, undefined, 4)}</pre>;
  }

  if (!travelRoute || !travelRoute.travelSchedule) {
    return <pre>Unknown server error. failed to create route.</pre>;
  }

  const schedule = travelRoute.travelSchedule;
  const { arrival, departure, duration } = travelRoute.travelSchedule;

  return (
    <div>
      {schedule.legs.map((leg) => {
        return leg.walk ? (
          <WalkInstructions leg={leg} />
        ) : (
          <>
            <BoardInstructions leg={leg} />
            <RideInstructions leg={leg} />
          </>
        );
      })}
      <div>
        <h1 className="font-semibold mt-3">You've reached your destination</h1>
        <h2 className="text-xs text-gray-700 font-semibold">
          Departure {formatTime(departure)} - Arrival {formatTime(arrival)} (
          {duration} min)
        </h2>
      </div>
    </div>
  );
};

export const TravelPage: React.FC = () => {
  const [origin, setOrigin] = useState<SelectedStop>(null);
  const [destination, setDestination] = useState<SelectedStop>(null);

  return (
    <Container>
      <h1 className="text-4xl mt-3">Travel Planner</h1>
      {/* <Wizard
        header={<WizardHeader destination={destination} origin={origin} />}
      >
        <SelectStop name="Origin" setter={setOrigin} />
        <SelectStop name="Destination" setter={setDestination} />
        {origin && destination && (
          <div>
            You selected {origin.name} and {destination.name}
          </div>
        )}
      </Wizard> */}
      <TravelPlanner />
    </Container>
  );
};
