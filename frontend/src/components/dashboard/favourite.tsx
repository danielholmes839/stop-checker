import { FavouriteIcon } from "components/travel";
import { FavouriteTravelLocation } from "core";
import { useNavigate } from "react-router-dom";

export const Favourite: React.FC<{
  favourite: FavouriteTravelLocation;
}> = ({ favourite }) => {
  const nav = useNavigate();
  return (
    <div className="px-3 py-2 bg-gray-50 rounded border-b">
      <div>
        <div className="inline-block align-middle">
          <FavouriteIcon placeId={favourite.id} />
        </div>
        <div
          className="pl-2 inline-block align-middle"
          style={{ maxWidth: "85%" }}
        >
          <span>{favourite.title}</span>
          <p className="text-xs">{favourite.description}</p>
        </div>
      </div>
      <div className="mt-2">
        <button
          className="text-primary-700 bg-primary-100 hover:bg-primary-200 px-2 rounded text-sm mr-2"
          onClick={() => nav(`/p/${favourite.id}`)}
        >
          Directions
        </button>
        <button
          className="text-primary-700 bg-primary-100 hover:bg-primary-200 px-2 rounded text-sm mr-2"
          onClick={() => nav(`/dashboard/favourite/${favourite.id}`)}
        >
          Edit
        </button>
      </div>
    </div>
  );
};
