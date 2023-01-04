import { Container } from "components/util";
import { Search, StopPreviewDefaultActions } from "./search";

export * from "./search";

export const SearchPage = () => {
  return (
    <Container>
      <div className="mt-3">
        <h1 className="text-3xl font-semibold mb-3">Search</h1>
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
