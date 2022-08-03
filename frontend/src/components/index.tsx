import React from "react";

type SignProps = {
  text: string;
  background: string;
  name: string; // route name
  headsign: string; // trip headsign
};

export const Sign: React.FC<SignProps> = ({
  text,
  background,
  name,
  headsign,
}) => {
  return (
    <span className="text-sm mr-3 inline-block">
      <span
        style={{
          color: text,
          background: background,
        }}
        className="px-1 rounded font-semibold text-xs"
      >
        {name}
      </span>{" "}
      {headsign}
    </span>
  );
};
