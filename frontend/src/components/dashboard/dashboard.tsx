import { useStorage } from "core";
import React from "react";

export const Dashboard: React.FC = () => {
  const { clear, clearHistory } = useStorage();
  return (
    <>
      <button onClick={clear}>Clear All</button>
      <button onClick={clearHistory}>Clear History</button>
    </>
  );
};
