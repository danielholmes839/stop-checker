import { StopPreviewFragment } from "client/types";
import { Search } from "pages/search";
import { StopPreviewActions } from "pages/search/search";
import React, { useState } from "react";
import { Link } from "react-router-dom";
import { Container } from "components";

const Actions: StopPreviewActions = ({ stop }) => {
  const { origin, setOrigin, destination, setDestination } = useAutomatic();

  const selectedAsOrigin = origin !== null && stop.id === origin.id;
  const selectedAsDestination =
    destination !== null && stop.id === destination.id;

  return (
    <div className="mt-3">
      <>
        <button
          disabled={selectedAsOrigin}
          className={
            selectedAsOrigin
              ? "mr-3 text-gray-500 text-sm"
              : "mr-3 text-primary-500 underline text-sm"
          }
          onClick={() => setOrigin(stop)}
        >
          {selectedAsOrigin ? "Selected as Origin" : "Set as Origin"}
        </button>
        <button
          className={
            selectedAsDestination
              ? "mr-3 text-gray-500 text-sm"
              : "mr-3 text-primary-500 underline text-sm"
          }
          onClick={() => setDestination(stop)}
        >
          {selectedAsDestination
            ? "Selected as Destination"
            : "Set as Destination"}
        </button>
      </>
    </div>
  );
};

type AutomaticContextValue = {
  origin: StopPreviewFragment | null;
  setOrigin: (stop: StopPreviewFragment) => void;
  destination: StopPreviewFragment | null;
  setDestination: (stop: StopPreviewFragment) => void;
  swap: () => void;
};

const AutomaticContext = React.createContext<AutomaticContextValue>({
  origin: null,
  setOrigin: (stop) => {},
  destination: null,
  setDestination: (stop) => {},
  swap: () => {},
});

const useAutomatic = () => React.useContext(AutomaticContext);

export const AutomaticProvider: React.FC = ({ children }) => {
  const [origin, setOrigin] = useState<StopPreviewFragment | null>(null);
  const [destination, setDestination] = useState<StopPreviewFragment | null>(
    null
  );
  const swap = () => {
    let originCopy = origin ? { ...origin } : null;
    setOrigin(destination);
    setDestination(originCopy);
  };

  return (
    <AutomaticContext.Provider
      value={{
        origin: origin,
        setOrigin: (stop) => {
          if (destination && destination.id === stop.id) {
            swap();
            return;
          }
          setOrigin(stop);
        },
        destination: destination,
        setDestination: (stop) => {
          if (origin && origin.id === stop.id) {
            swap();
            return;
          }
          setDestination(stop);
        },
        swap: swap,
      }}
    >
      {children}
    </AutomaticContext.Provider>
  );
};

const Current: React.FC = () => {
  const { origin, destination } = useAutomatic();

  return (
    <div className="mt-3">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-x-5 gap-y-3">
        <div className="border-b pb-1 mb-2">
          <h1>
            <div className="bg-red-600 rounded-full px-2 text-white text-xs inline mr-2 font-bold">
              A
            </div>
            {origin ? (
              <span className="text-sm font-semibold">
                {origin.name} (Origin)
              </span>
            ) : (
              <span className="text-sm text-gray-600 font-semibold">
                Origin
              </span>
            )}
          </h1>
        </div>
        <div className="border-b pb-2 mb-2">
          <h1>
            <div className="bg-red-600 rounded-full px-2 text-white text-xs inline mr-2 font-bold">
              B
            </div>{" "}
            {destination ? (
              <span className="text-sm font-semibold">
                {destination.name} (Destination)
              </span>
            ) : (
              <span className="text-sm text-gray-600 font-semibold">
                Destination
              </span>
            )}
          </h1>
        </div>
      </div>
      <div className="mb-3">
        {origin && destination ? (
          <Link
            className="text-primary-500 text-sm"
            to={`/travel/${origin.id}/${destination.id}`}
          >
            Next
          </Link>
        ) : (
          <p className="text-xs text-red-600">
            Please select an origin and destination using the search below
          </p>
        )}
      </div>
    </div>
  );
};

export const AutomaticInput: React.FC = () => {
  return (
    <Container>
      <AutomaticProvider>
        <div className="mt-3">
          <h1 className="text-3xl font-semibold">Travel Planner</h1>
        </div>

        <Current />

        <div className="mt-3">
          <Search
            config={{
              Actions: Actions,
              enableMap: true,
              enableStopRouteLinks: false,
            }}
          />
        </div>
      </AutomaticProvider>
    </Container>
  );
};
