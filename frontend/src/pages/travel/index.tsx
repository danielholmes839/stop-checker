import { Container } from "components";
import { Search } from "../search";
import { Wizard, useWizard } from "react-use-wizard";
import { useEffect, useState } from "react";

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
  const { nextStep, activeStep } = useWizard();

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
    ? "border-b-2 bg-gray-50 px-3 rounded-t-sm border-indigo-500 py-1 cursor-pointer"
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

export const TravelPage: React.FC = () => {
  const [origin, setOrigin] = useState<SelectedStop>(null);
  const [destination, setDestination] = useState<SelectedStop>(null);

  return (
    <Container>
      <h1 className="text-4xl mt-3">Travel Planner</h1>
      <Wizard
        header={<WizardHeader destination={destination} origin={origin} />}
      >
        <SelectStop name="Origin" setter={setOrigin} />
        <SelectStop name="Destination" setter={setDestination} />
        {origin && destination && (
          <div>
            You selected {origin.name} and {destination.name}
          </div>
        )}
      </Wizard>
    </Container>
  );
};
