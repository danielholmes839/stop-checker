import { SearchLocationIcon } from "components/travel";
import { TravelLocation, useStorage } from "core";
import { useNavigate } from "react-router-dom";

export const History: React.FC<{
  recent: TravelLocation;
}> = ({ recent }) => {
  const nav = useNavigate();
  const { addFavourite } = useStorage();
  return (
    <div className="px-3 py-2 bg-gray-50 rounded border-b">
      <div>
        <div className="inline-block align-middle">
          <SearchLocationIcon placeId={recent.id} />
        </div>
        <div
          className="pl-2 inline-block align-middle"
          style={{ maxWidth: "85%" }}
        >
          <span>{recent.title}</span>
          <p className="text-xs">{recent.description}</p>
        </div>
      </div>
      <div className="mt-2">
        <button
          className="text-primary-700 bg-primary-100 hover:bg-primary-200 px-2 rounded text-sm mr-2"
          onClick={() => nav(`/p/${recent.id}`)}
        >
          Directions
        </button>
        <button
          className="bg-gray-200 border hover:bg-gray-300 px-2 rounded text-sm mr-2"
          onClick={() => {
            addFavourite(recent);
            nav(`/dashboard/favourite/${recent.id}`);
          }}
        >
          Favourite
        </button>
      </div>
    </div>
  );
};
