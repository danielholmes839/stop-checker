import React from "react";
import { Container } from "components/util";
import { TravelLocationInput } from "components/travel";
import { useNavigate } from "react-router-dom";

export const Nearby: React.FC = () => {
  const nav = useNavigate();

  return (
    <Container>
      <h1 className="text-3xl font-bold mt-3">Enter Location</h1>
      <div className="mt-2">
        <TravelLocationInput
          setTravelLocation={(location) => {
            nav(`/search/${location.id}`);
          }}
        />
      </div>
    </Container>
  );
};
