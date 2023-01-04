import React from "react";
import { Container } from "./util";

import usePlacesService from "react-google-autocomplete/lib/usePlacesAutocompleteService";

export const Travel: React.FC = () => {
  return (
    <Container>
      <h1 className="text-3xl mt-5">Travel</h1>

      <input className="bg-gray-100" />
    </Container>
  );
};
