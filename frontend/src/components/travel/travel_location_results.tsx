import { PlacePrediction, useStorage } from "core";
import { SearchLocationIcon } from "components/travel";

export const TravelLocationResults: React.FC<{
  setPlaceId: React.Dispatch<string>;
  predictions: PlacePrediction[];
}> = ({ setPlaceId, predictions }) => {
  const { getFavourite } = useStorage();
  return (
    <div>
      {predictions.map((p) => {
        let pred = p;
        let fav = getFavourite(pred.id);
        if (fav) {
          pred = {
            ...p,
            title: fav.title,
          };
        }
        return (
          <TravelLocationResult
            key={pred.id}
            pred={pred}
            setPlaceId={setPlaceId}
          />
        );
      })}
    </div>
  );
};

export const TravelLocationResult: React.FC<{
  setPlaceId: React.Dispatch<string>;
  pred: PlacePrediction;
}> = ({ pred, setPlaceId }) => {
  return (
    <div
      key={pred.id}
      className="px-3 py-2 mt-2 bg-gray-50 rounded border-b hover:bg-primary-100 cursor-pointer"
      onClick={() => setPlaceId(pred.id)}
    >
      <div className="inline-block align-middle">
        <SearchLocationIcon placeId={pred.id} />
      </div>
      <div
        className="pl-2 inline-block align-middle"
        style={{ maxWidth: "90%" }}
      >
        <span>{pred.title}</span>
        <p className="text-xs">{pred.description}</p>
      </div>
    </div>
  );
};

export const TravelCurrentLocationOption: React.FC<{
  onClick: () => void;
}> = ({ onClick }) => {
  return (
    <div
      className="px-3 py-2 mt-2 bg-gray-50 rounded border-b hover:bg-primary-100 cursor-pointer"
      onClick={onClick}
    >
      <div className="inline-block align-middle">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 20 20"
          fill="currentColor"
          className="w-6 h-6 inline-block align-middle"
        >
          <path
            fillRule="evenodd"
            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-1.5 0a6.5 6.5 0 11-11-4.69v.447a3.5 3.5 0 001.025 2.475L8.293 10 8 10.293a1 1 0 000 1.414l1.06 1.06a1.5 1.5 0 01.44 1.061v.363a1 1 0 00.553.894l.276.139a1 1 0 001.342-.448l1.454-2.908a1.5 1.5 0 00-.281-1.731l-.772-.772a1 1 0 00-1.023-.242l-.384.128a.5.5 0 01-.606-.25l-.296-.592a.481.481 0 01.646-.646l.262.131a1 1 0 00.447.106h.188a1 1 0 00.949-1.316l-.068-.204a.5.5 0 01.149-.538l1.44-1.234A6.492 6.492 0 0116.5 10z"
            clipRule="evenodd"
          />
        </svg>
      </div>
      <div
        className="pl-2 inline-block align-middle"
        style={{ maxWidth: "90%" }}
      >
        <span>Current Location</span>
      </div>
    </div>
  );
};
