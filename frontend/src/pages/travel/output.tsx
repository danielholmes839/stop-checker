import {
  ScheduleMode,
  TravelOptions,
  useTravelPlannerQuery,
} from "client/types";
import { Container } from "components";
import { formatDateTime } from "helper";
import React, { useState } from "react";
import DateTimePicker from "react-datetime-picker";
import { useParams } from "react-router-dom";
import { Instructions } from "./instructions";

type OptionContextValue = {
  options: TravelOptions;
  setDate: React.Dispatch<React.SetStateAction<Date>>;
  setMode: React.Dispatch<React.SetStateAction<ScheduleMode>>;
};
const OptionContext = React.createContext<OptionContextValue>({
  options: {
    mode: ScheduleMode.DepartAt,
    datetime: new Date(),
  },
  setDate: (date) => {},
  setMode: (mode) => {},
});

const useOptions = () => React.useContext(OptionContext);

const OptionProvider: React.FC = ({ children }) => {
  const [mode, setMode] = useState(ScheduleMode.DepartAt);
  const [date, setDate] = useState(new Date());
  return (
    <OptionContext.Provider
      value={{
        options: {
          mode: mode,
          datetime: date,
        },
        setDate: setDate,
        setMode: setMode,
      }}
    >
      {children}
    </OptionContext.Provider>
  );
};

const OptionSelect: React.FC = () => {
  const { options, setMode, setDate } = useOptions();
  return (
    <div>
      <button
        className={`bg-${
          options.mode === ScheduleMode.DepartAt ? "primary" : "gray"
        }-200 px-3 py-1 mr-1 rounded-full text-xs font-semibold`}
        onClick={() => setMode(ScheduleMode.DepartAt)}
      >
        Depart At
      </button>
      <button
        className={`bg-${
          options.mode === ScheduleMode.ArriveBy ? "primary" : "gray"
        }-200 px-3 py-1 mr-1 rounded-full text-xs font-semibold`}
        onClick={() => setMode(ScheduleMode.ArriveBy)}
      >
        Arrive By
      </button>
      <DateTimePicker
        onChange={(date: Date | null) => {
          if (date) {
            setDate(date);
          }
        }}
        disableClock={true}
        required={true}
        value={options.datetime}
        className="text-sm"
      />
    </div>
  );
};

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

  return <Instructions data={data.travelPlanner} />;
};

export const AutomaticOutput: React.FC = () => {
  return (
    <Container>
      <div className="my-3">
        <h1 className="text-3xl font-semibold">Travel Schedule</h1>
      </div>
      <OptionProvider>
        <OptionSelect />
        <Display />
      </OptionProvider>
    </Container>
  );
};
