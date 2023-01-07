import { createClient } from "urql";

export const googleMapsKey = "";

export const client = createClient({
  url:
    process.env.NODE_ENV === "production"
      ? "https://api.stop-checker.com/graphql"
      : "http://192.168.0.214:3001/graphql",
});
