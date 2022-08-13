import { Container } from "components";
import { useNavigate } from "react-router-dom";
import { Search } from "./search";

export { Search } from "./search";

export const SearchPage = () => {
  const navigate = useNavigate();
  return (
    <Container>
      <Search
        config={{
          action: ({ id }) => navigate(`/stop/${id}`),
          actionName: "View",
        }}
      />
    </Container>
  );
};
