import { Container, Sign } from "components";
import { Search } from "../search";
import { Wizard, useWizard } from "react-use-wizard";
import React, { useState } from "react";
import {
  RouteType,
  ScheduleMode,
  TravelScheduleLegDefaultFragment,
  useTravelPlannerQuery,
} from "client/types";
import { formatDistance, formatRouteType, formatTime } from "format";

import { ReactComponent as Bus } from "./icons/bus.svg";
import { ReactComponent as Train } from "./icons/train.svg";
import Walk from "./icons/walk.png";

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

type InstructionProps = { leg: TravelScheduleLegDefaultFragment };

const BoardInstructions: React.FC<InstructionProps> = ({ leg }) => {
  const { origin, destination, duration, departure, transit } = leg;

  if (!transit) {
    return <></>;
  }

  const { route, trip } = transit;

  return (
    <div className="border-l-4 pl-3" style={{ borderColor: route.background }}>
      <h1 className="font-semibold">
        Board the{" "}
        <Sign
          props={{
            background: route.background,
            name: route.name,
            text: route.text,
          }}
        />{" "}
        at {origin.name} #{origin.code}
        {route.type === RouteType.Train ? (
          <Train
            className="float-right w-5 h-5 mt-1 mr-1 inline"
            stroke="#3730a3"
            fill="#3730a3"
          />
        ) : (
          <Bus
            className="float-right w-5 h-5 mt-1 mr-1 inline"
            stroke="#3730a3"
            fill="#3730a3"
          />
        )}
      </h1>
      <h2 className="text-xs text-gray-700 font-semibold">
        Towards {trip.headsign}
      </h2>

      <p className="text-sm mt-2">Scheduled at {formatTime(departure)}</p>
      <button className="text-primary-500 text-xs">More</button>
    </div>
  );
};

const RideInstructions: React.FC<InstructionProps> = ({ leg }) => {
  const { destination, duration, arrival, transit } = leg;
  if (!transit) {
    return <></>;
  }

  const { route } = transit;

  return (
    <div className="border-l-4 pl-3" style={{ borderColor: route.background }}>
      <h1 className="font-semibold">
        Exit the{" "}
        <Sign
          props={{
            background: route.background,
            name: route.name,
            text: route.text,
          }}
        />{" "}
        at {destination.name} #{destination.code}
      </h1>
      <h2 className="text-xs text-gray-700 font-semibold">
        Arrival {formatTime(arrival)}
      </h2>
      <button className="text-primary-500 text-xs">
        {transit.arrival.sequence - transit.departure.sequence} stops (
        {duration} min)
      </button>
    </div>
  );
};

const WalkInstructions: React.FC<InstructionProps> = ({ leg }) => {
  const { origin, destination, duration, departure, arrival, distance } = leg;

  return (
    <div className="border-l-4 pl-3 border-gray-300 border-dashed">
      <h1 className="font-semibold">
        <span className="align-text-bottom">
          Walk to {destination.name} #{destination.code}
        </span>
      </h1>
      <h2 className="text-xs text-gray-700 font-semibold">
        {formatDistance(distance)} ({duration} min)
      </h2>
    </div>
  );
};

const Instruction: React.FC = ({ children }) => {
  return <div className="border-b pb-3 mt-3">{children}</div>;
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
  const { arrival, duration } = travelRoute.travelSchedule;

  return (
    <div>
      {schedule.legs.map((leg) => {
        return leg.walk ? (
          <Instruction>
            <WalkInstructions leg={leg} />
          </Instruction>
        ) : (
          <>
            <Instruction>
              <BoardInstructions leg={leg} />
            </Instruction>
            <Instruction>
              <RideInstructions leg={leg} />
            </Instruction>
          </>
        );
      })}
      <div>
        <h1 className="font-semibold mt-3">You've reached your destination</h1>
        <h2 className="text-xs text-gray-700 font-semibold">
          Arrival {formatTime(arrival)} ({duration} min)
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
