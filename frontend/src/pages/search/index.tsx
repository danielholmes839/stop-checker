import { Container } from "components";
import { Search, StopPreviewDefaultActions } from "./search";

export { Search } from "./search";

export const SearchPage = () => {
  return (
    <Container>
      <div className="mt-3">
        <Search
          config={{
            Actions: StopPreviewDefaultActions,
            enableMap: true,
            enableStopRouteLinks: true,
          }}
        />
      </div>
    </Container>
  );
};
