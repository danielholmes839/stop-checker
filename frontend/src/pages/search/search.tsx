import React, { useState } from "react";
import {
  Stop,
  TextSearchQuery,
  useStopPreviewQuery,
  useTextSearchQuery,
} from "client/types";
import { Sign, Card, QueryResponseWrapper } from "components";
import { useDebounce } from "use-debounce";

import { SearchConfig } from "./config";
import { SearchMap } from "./map";

const StopPreview: React.FC<{ stop: Stop; config: SearchConfig }> = ({
  stop,
  config,
}) => {
  const { id, name, code, routes } = stop;
  return (
    <Card key={id}>
      <h1 className="font-semibold mb-1">
        {name} <span className="float-right text-gray-700">#{code}</span>
      </h1>

      <div>
        {routes.map(({ direction, headsign, route }) => (
          <Sign
            key={route.name + direction}
            props={{
              background: route.background,
              headsign: headsign,
              name: route.name,
              text: route.text,
            }}
          />
        ))}
      </div>
      <div>
        <button
          className="border border-indigo-500 px-5 py-1 mt-2 hover:bg-indigo-500 hover:text-white text-indigo-500 text-sm rounded-sm"
          onClick={() => {
            config.action(stop as Stop);
          }}
        >
          {config.actionName}
        </button>
      </div>
    </Card>
  );
};
const TextSearchQueryResponse: React.FC<{
  data: TextSearchQuery;
  config: SearchConfig;
}> = ({ data, config }) => {
  return (
    <div>
      {data.searchStopText.results.map((stop) => {
        return (
          <div className="mb-3">
            <StopPreview stop={stop as Stop} config={config} />
          </div>
        );
      })}
    </div>
  );
};

const SelectedStopPreview: React.FC<{
  id: string;
  config: SearchConfig;
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

  return (
    <div className="mb-3">
      <StopPreview stop={debounced.stop as Stop} config={config} />
      <hr className="bg-indigo-500 rounded-full mt-3" style={{ height: 2 }} />
    </div>
  );
};

export const Search: React.FC<{ config: SearchConfig }> = ({ config }) => {
  const [selected, setSelected] = useState<Stop | null>(null);
  const [searchText, setSearchText] = useState("");
  const [searchTextDebounced] = useDebounce(searchText, 200);

  const [response, _] = useTextSearchQuery({
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
        className="bg-gray-50 border-b-2 border-gray-200 w-full mb-1 p-3 rounded-t-sm focus:outline-none focus:border-b-1 focus:border-indigo-500 sm:text-xs md:text-sm"
        value={searchText}
        type="text"
        placeholder="Search by stop name or code Ex. Rideau A, O-Train, 3000"
        onChange={(event) => setSearchText(event.target.value)}
      />
      <SearchMap
        config={config}
        selected={selected}
        setSelected={setSelected}
      />
      {selected && <SelectedStopPreview config={config} id={selected.id} />}
      <QueryResponseWrapper response={response}>
        <TextSearchQueryResponse
          data={response.data as TextSearchQuery}
          config={config}
        />
      </QueryResponseWrapper>
    </>
  );
};
