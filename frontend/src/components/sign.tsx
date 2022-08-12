import React from "react";

type SignProps = {
  text: string;
  background: string;
  name: string; // route name
  headsign: string; // trip headsign
  css?: string;
};

export const Sign: React.FC<{ props: SignProps }> = ({ props }) => {
  const { text, background, name, headsign, css = "text-sm" } = props;
  return (
    <span className={`${css} mr-3 inline-block`}>
      <span
        style={{
          color: text,
          background: background,
        }}
        className="px-1 rounded font-semibold text-xs tracking-wider"
      >
        {name}
      </span>{" "}
      {headsign}
    </span>
  );
};
