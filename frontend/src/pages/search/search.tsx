import React, { useState } from "react";
import { Stop, TextSearchQuery, useTextSearchQuery } from "client/types";
import { Sign, Container, Card, QueryResponseWrapper } from "components";
import { useNavigate } from "react-router-dom";
import { useDebounce } from "use-debounce";

import { SearchConfig } from "./config";
import { SearchMap } from "./map";

const TextSearchQueryResponse: React.FC<{
  data: TextSearchQuery;
  config: SearchConfig;
}> = ({ data, config }) => {
  return (
    <div>
      {data.searchStopText.map((stop) => {
        const { id, name, code, routes } = stop;
        return (
          <div className="mb-3 cursor-pointer">
            <Card key={id}>
              <h1 className="font-semibold mb-1">
                {name}{" "}
                <span className="float-right text-gray-700">#{code}</span>
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
          </div>
        );
      })}
    </div>
  );
};

export const Search: React.FC<{ config: SearchConfig }> = ({ config }) => {
  const [searchText, setSearchText] = useState("");
  const [searchTextDebounced] = useDebounce(searchText, 200);

  const [response, _] = useTextSearchQuery({
    variables: {
      text: searchTextDebounced,
    },
  });

  return (
    <>
      <input
        className="bg-gray-50 border-gray-100 w-full mb-3 mt-1 p-3 rounded-sm focus:outline-none border-b-2 focus:border-indigo-500 text-sm"
        value={searchText}
        type="text"
        placeholder="Search by stop name or code Ex. Rideau A, O-Train, 3000"
        onChange={(event) => setSearchText(event.target.value)}
      />
      <SearchMap config={config} />
      <QueryResponseWrapper response={response}>
        <TextSearchQueryResponse
          data={response.data as TextSearchQuery}
          config={config}
        />
      </QueryResponseWrapper>
    </>
  );
};
