import { Container } from "components";
import { OptionInput, OptionProvider } from "providers";
import { useStorage } from "providers/storage";
import React from "react";

export const Dashboard: React.FC = () => {
  const { routes } = useStorage();
  return (
    <Container>
      <OptionProvider>
        <div className="my-3">
          <h1 className="text-3xl font-semibold">Dashboard</h1>
        </div>
        <OptionInput />
        {routes.map((route) => (
          <pre>{JSON.stringify(route, undefined, 4)}</pre>
        ))}
      </OptionProvider>
    </Container>
  );
};
