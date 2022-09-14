import { ScheduleMode, TravelOptions } from "client/types";
import React, { useState } from "react";
import DateTimePicker from "react-datetime-picker";

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

export const useOptions = () => React.useContext(OptionContext);

export const OptionProvider: React.FC<{
  initial?: TravelOptions;
}> = ({ children, initial }) => {
  const [mode, setMode] = useState<ScheduleMode>(
    initial ? initial.mode : ScheduleMode.DepartAt
  );
  const [date, setDate] = useState<Date>(
    initial ? initial.datetime : new Date()
  );

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

export const OptionInput: React.FC = () => {
  const { options, setMode, setDate } = useOptions();
  return (
    <div>
      <button
        className={`bg-${
          options.mode === ScheduleMode.DepartAt ? "primary" : "gray"
        }-200 px-3 py-1 mr-1 rounded-full text-xs font-semibold mb-1`}
        onClick={() => setMode(ScheduleMode.DepartAt)}
      >
        Depart At
      </button>
      <button
        className={`bg-${
          options.mode === ScheduleMode.ArriveBy ? "primary" : "gray"
        }-200 px-3 py-1 mr-1 rounded-full text-xs font-semibold mb-1`}
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
