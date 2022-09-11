import { StopPageQuery, useStopPageQuery } from "client/types";
import { Container, Sign, Card } from "components";
import { Link, useParams } from "react-router-dom";

const StopPageResponse: React.FC<{ data: StopPageQuery }> = ({ data }) => {
  const { stop } = data;
  if (!stop) {
    return <div>stop not found</div>;
  }

  return (
    <div>
      <h1 className="mt-3 font-semibold text-lg">
        {stop.name} #{stop.code}
      </h1>
      <div className="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-3 gap-3 mt-3">
        {stop.routes
          .sort((a, b) => a.route.name.localeCompare(b.route.name))
          .map((stopRoute, i) => {
            const { route, headsign, schedule, liveBuses } = stopRoute;
            return (
              <Card key={i}>
                <h1 className="text-xs">
                  <Sign
                    props={{
                      background: route.background,
                      text: route.text,
                      name: route.name,
                    }}
                  />{" "}
                  <span className="text-sm font-semibold">{headsign}</span>
                </h1>
                <div>
                  <div className="mt-1">
                    {schedule.next.length > 0 ? (
                      schedule.next.map(({ stoptime }, i) => (
                        <span key={i} className="mr-2 text-sm inline-block">
                          {stoptime.time}
                        </span>
                      ))
                    ) : (
                      <span className="text-sm">
                        No more stops today or tomorrow
                      </span>
                    )}
                  </div>
                  <div>
                    <span className="text-xs text-gray-500">
                      {liveBuses.length === 1
                        ? "1 Live bus"
                        : `${liveBuses.length} Live buses`}
                    </span>
                  </div>
                  <div>
                    <Link to={`/stop/${stop.id}/route/${route.id}`}>
                      <button
                        className="border border-primary-500 px-5 py-0 mt-2 hover:bg-primary-500 hover:text-white text-primary-500 text-sm rounded-sm w-full"
                        type="button"
                      >
                        View
                      </button>
                    </Link>
                  </div>
                </div>
              </Card>
            );
          })}
      </div>
    </div>
  );
};

export const StopPage: React.FC = () => {
  const { id } = useParams();
  const [{ data, fetching, error }, _ignore] = useStopPageQuery({
    variables: { id: id === undefined ? "" : id },
  });

  if (fetching) {
    return <>Loading...</>;
  }

  if (error) {
    return <>{error.message}</>;
  }

  if (!data) {
    return <></>;
  }

  return (
    <Container>
      <StopPageResponse data={data} />
    </Container>
  );
};
