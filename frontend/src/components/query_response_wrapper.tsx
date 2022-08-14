import { UseQueryState } from "urql";

export const QueryResponseWrapper: React.FC<{ response: UseQueryState }> = ({
  response,
  children,
}) => {
  const { data, fetching, error } = response;
  if (fetching || error || !data) {
    return <></>;
  }

  return <>{children}</>;
};
