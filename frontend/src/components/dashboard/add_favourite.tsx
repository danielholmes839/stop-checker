import { TravelLocationInput } from "components/travel";
import { Container } from "components/util";
import { useStorage } from "core";
import { useNavigate } from "react-router-dom";

export const DashboardAddFavourite: React.FC = () => {
  const nav = useNavigate();
  const { addFavourite } = useStorage();
  return (
    <Container>
      <h1 className="text-3xl font-bold font mt-3 mb-1">New Favourite</h1>
      <TravelLocationInput
        suggestFavourites={false}
        setTravelLocation={(location) => {
          addFavourite(location);
          nav(`/dashboard/favourite/${location.id}`);
        }}
      />
    </Container>
  );
};
