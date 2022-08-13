import { Stop } from "client/types";

export type SearchConfig = {
  actionName: string;
  action: (stop: Stop) => void;
};
