import { Container } from "components";

export const PlannerPage: React.FC = () => {
  return (
    <Container>
      <h1 className="text-4xl">Travel Planner</h1>
      <ul>
        <li>Origin</li>
        <li>Destination</li>
        <li>View Plan</li>
      </ul>
    </Container>
  );
};
