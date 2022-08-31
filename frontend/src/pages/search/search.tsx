import React, { useState } from "react";
import {
  Stop,
  StopPreviewFragment,
  TextSearchQuery,
  useStopPreviewQuery,
  useTextSearchQuery,
} from "client/types";
import { Sign, Card } from "components";
import { useDebounce } from "use-debounce";

import { SearchMap } from "./map";
import { Link } from "react-router-dom";

type Config = {
  enableMap: boolean;
  enableStopRouteLinks: boolean;
  Actions: StopPreviewActions;
};

type StopPreviewActions = React.FC<{ stop: StopPreviewFragment }>;

const FlagLink: React.FC<{ to: string; enabled: boolean }> = ({
  to,
  enabled,
  children,
}) => {
  if (enabled) {
    return (
      <Link to={to} className="text-primary-700">
        {children}
      </Link>
    );
  }
  return <>{children}</>;
};

const StopPreview: React.FC<{
  stop: StopPreviewFragment;
  config: Config;
}> = ({ stop, config }) => {
  const { id, name, code, routes } = stop;
  const { Actions, enableStopRouteLinks } = config;
  return (
    <Card key={id}>
      <h1 className="font-semibold mb-1">
        {name} <span className="float-right text-gray-700">#{code}</span>
      </h1>

      <div>
        {routes.map(({ headsign, route }, i) => (
          <span key={i} className="text-sm mr-3 inline-block">
            <FlagLink
              enabled={enableStopRouteLinks}
              to={`/stop/${stop.id}/route/${route.id}`}
            >
              <span className="text-xs">
                <Sign
                  key={route.name}
                  props={{
                    background: route.background,
                    name: route.name,
                    text: route.text,
                  }}
                />
              </span>{" "}
              {headsign}
            </FlagLink>
          </span>
        ))}
      </div>
      <div>
        <Actions stop={stop} />
      </div>
    </Card>
  );
};

export const StopPreviewDefaultActions: StopPreviewActions = ({ stop }) => {
  const to =
    stop.routes.length > 1
      ? `/stop/${stop.id}`
      : `/stop/${stop.id}/route/${stop.routes[0].route.id}`;

  return (
    <Link to={to}>
      <button className="border border-primary-500 px-5 py-1 mt-2 hover:bg-primary-500 hover:text-white text-primary-500 text-sm rounded-sm">
        View
      </button>
    </Link>
  );
};

const SelectedStopPreview: React.FC<{
  id: string;
  config: Config;
}> = ({ id, config }) => {
  const [{ data }, _] = useStopPreviewQuery({
    variables: {
      id: id,
    },
  });

  const [debounced, __] = useDebounce(data, 100);
  if (!debounced) {
    return <></>;
  }

  if (!debounced.stop) {
    return <></>;
  }

  return (
    <div className="mb-3">
      <StopPreview stop={debounced.stop} config={config} />
      <hr className="bg-primary-500 rounded-full mt-3" style={{ height: 2 }} />
    </div>
  );
};

const SearchResults: React.FC<{
  data: TextSearchQuery | undefined;
  config: Config;
}> = ({ data, config }) => {
  if (!data) {
    return <></>;
  }

  return (
    <div>
      {data &&
        data.searchStopText.results.map((stop) => {
          return (
            <div key={stop.id} className="mb-3">
              <StopPreview stop={stop as Stop} config={config} />
            </div>
          );
        })}
    </div>
  );
};

export const Search: React.FC<{ config: Config }> = ({ config }) => {
  const { enableMap } = config;
  const [selected, setSelected] = useState<Stop | null>(null);
  const [searchText, setSearchText] = useState("");
  const [searchTextDebounced] = useDebounce(searchText, 200);

  const [{ data }, _] = useTextSearchQuery({
    variables: {
      text: searchTextDebounced,
      page: {
        limit: 15,
        skip: 0,
      },
    },
  });

  return (
    <>
      <input
        className="bg-gray-50 border-b-2 border-gray-200 w-full mb-1 p-3 rounded-t-sm focus:outline-none focus:border-b-1 focus:border-primary-500 sm:text-xs md:text-sm"
        value={searchText}
        type="text"
        placeholder="Search by stop name or code Ex. Rideau A, O-Train, 3000"
        onChange={(event) => setSearchText(event.target.value)}
      />
      {enableMap && (
        <>
          <SearchMap selected={selected} setSelected={setSelected} />
          {selected && <SelectedStopPreview config={config} id={selected.id} />}
        </>
      )}

      <SearchResults data={data} config={config} />
    </>
  );
};
