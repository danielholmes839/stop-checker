import { StopPageQuery, useStopPageQuery } from "client/types";
import { Container, QueryResponseWrapper, Sign, Card } from "components";
import { useParams } from "react-router-dom";

const StopPageResponse: React.FC<{ data: StopPageQuery }> = ({ data }) => {
  const { stop } = data;
  if (!stop) {
    return <div>stop not found</div>;
  }

  return (
    <div className="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
      {stop.routes
        .sort((a, b) => a.route.name.localeCompare(b.route.name))
        .map((stopRoute) => {
          const { route, headsign, schedule } = stopRoute;
          return (
            <Card>
              <Sign
                props={{
                  background: route.background,
                  text: route.text,
                  name: route.name,
                }}
              />
              <div>
                {schedule.next.length > 0 ? (
                  schedule.next.map((stopTime) => (
                    <span className="mr-2 text-sm">{stopTime.time}</span>
                  ))
                ) : (
                  <span className="text-sm">
                    No more stops today or tomorrow
                  </span>
                )}
              </div>
            </Card>
          );
        })}
    </div>
  );
};

export const StopPage: React.FC = () => {
  const { id } = useParams();
  const [response, _] = useStopPageQuery({
    variables: { id: id === undefined ? "" : id },
  });

  return (
    <Container>
      <QueryResponseWrapper response={response}>
        {response.data && <StopPageResponse data={response.data} />}
      </QueryResponseWrapper>
    </Container>
  );
};
