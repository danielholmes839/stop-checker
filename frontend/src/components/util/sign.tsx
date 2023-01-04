import React from "react";

type SignProps = {
  text: string; // color
  background: string; // color
  name: string; // route name
};

export const Sign: React.FC<{ props: SignProps }> = ({ props }) => {
  const { text, background, name } = props;
  const px = name.length === 1 ? "px-2" : "px-1";

  return (
    <span
      style={{
        color: text,
        background: background,
      }}
      className={`${px} rounded font-semibold tracking-wider`}
    >
      {name}
    </span>
  );
};
