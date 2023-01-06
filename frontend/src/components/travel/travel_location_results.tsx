import { PlacePrediction, useStorage } from "core";
import { TravelLocationIcon } from "components/travel";

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
        return <TravelLocationResult pred={pred} setPlaceId={setPlaceId} />;
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
        <TravelLocationIcon placeId={pred.id} />
      </div>
      <div
        className="pl-2 inline-block align-middle"
        style={{ maxWidth: "90%" }}
      >
        <span>{pred.title}</span>
        <p className="text-xs mt-1">{pred.description}</p>
      </div>
    </div>
  );
};
