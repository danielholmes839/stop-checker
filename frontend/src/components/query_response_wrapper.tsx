import { UseQueryState } from "urql";

export const QueryResponseWrapper: React.FC<{ response: UseQueryState }> = ({
  response,
  children,
}) => {
  const { data, fetching, error } = response;
  if (fetching) {
    return <></>;
  }

  if (error) {
    return <div>{error.toString()}</div>;
  }

  if (data === null || data === undefined) {
    return <div>no data</div>;
  }

  return <>{children}</>;
};
