import { useTextSearchQuery } from "client/types";
import { Sign } from "components";
import React from "react";

const App: React.FC = () => {
  const [{ data, fetching, error }, _] = useTextSearchQuery({
    variables: {
      text: "pleasant park arch",
    },
  });

  if (fetching) {
    return <>loading</>;
  }

  if (error) {
    return <>{error.toString()}</>;
  }

  if (data === undefined) {
    return <>no data</>;
  }

  return (
    <div className="container mx-auto px-10 xl:w-1/3">
      <h1 className="text-4xl">stop-checker.com</h1>
      <div>
        {data.searchStopText.map(({ id, name, code, routes }) => (
          <div
            key={id}
            className="border border-gray-100 rounded px-5 py-3 mb-3 bg-white hover:bg-gray-50 hover:shadow-sm cursor-pointer"
          >
            <h1 className="">
              {name} #{code}
            </h1>

            <div>
              {routes.map(({ direction, headsign, route }) => (
                <Sign
                  key={route.name + direction}
                  background={route.background}
                  headsign={headsign}
                  name={route.name}
                  text={route.text}
                />
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default App;
