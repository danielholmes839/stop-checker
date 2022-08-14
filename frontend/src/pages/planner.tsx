import { Container } from "components";
import { Search } from "./search";

export const PlannerPage: React.FC = () => {
  return (
    <Container>
      <h1 className="text-4xl">Travel Planner</h1>
      <Search
        config={{
          action: (stop) => {},
          actionName: "Set Origin",
        }}
      />
    </Container>
  );
};
