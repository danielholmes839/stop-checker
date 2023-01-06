import React, { useState } from "react";
import { Container } from "components/util";

import { TravelLocation } from "core";
import { TravelLocationInput } from "components/travel";

export const TravelLocationDisplay: React.FC<{
  travelLocation: TravelLocation | null;
  symbol: string;
}> = ({ travelLocation, symbol }) => {
  return (
    <div className="px-3 py-2 bg-gray-50 rounded border-b">
      <div className="inline-block align-middle text-4xl font-bold mr-2">
        {symbol}
      </div>
      <div
        className="pl-2 border-l border-gray-300 inline-block align-middle"
        style={{ maxWidth: "90%" }}
      >
        {travelLocation ? (
          <div>
            <h2>{travelLocation.title}</h2>
            <span className="text-xs">{travelLocation.description}</span>
          </div>
        ) : (
          <div>
            <h2>Not Selected</h2>
          </div>
        )}
      </div>
    </div>
  );
};
export const Travel: React.FC = () => {
  const [travelLocation, setTravelLocation] = useState<TravelLocation | null>(
    null
  );

  return (
    <Container>
      <h1 className="text-3xl font-bold font mt-3">Travel Planner</h1>
      <div className="mt-2">
        <TravelLocationDisplay travelLocation={travelLocation} symbol={"A"} />
      </div>
      <div className="mt-3">
        <h2 className="text-xl font-bold">Where do you want to go?</h2>
        <div className="mt-1">
          <TravelLocationInput setTravelLocation={setTravelLocation} />
        </div>
      </div>
    </Container>
  );
};
