import { createClient } from "urql";

export const client = createClient({
  url: process.env.NODE_ENV === "production" ? "https://api.stop-checker.com/graphql" : "http://localhost:3001/graphql",
});
