import { PlacePrediction, useStorage } from "core";
import { SearchLocationIcon } from "components/travel";
import { useState } from "react";

const requestCurrentLocation = (
  setPlaceId: React.Dispatch<string>,
  setCurrentLocationError: React.Dispatch<string | null>
) => {
  const getCurrentLocation = () => {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        let geocoder = new google.maps.Geocoder();
        geocoder.geocode(
          {
            location: new google.maps.LatLng({
              lat: position.coords.latitude,
              lng: position.coords.longitude,
            }),
          },
          (res) => {
            if (res === null || res.length === 0) {
              setCurrentLocationError(
                "Sorry, we couldn't find the address of your current location. Please enter your location manually."
              );
              return;
            }
            if (res.length > 1) {
              setPlaceId(res[1].place_id);
            } else {
              setPlaceId(res[0].place_id);
            }
          }
        );
      },
      (error) => {
        setCurrentLocationError(`${error.message} (${error.code})`);
      }
    );
  };

  navigator.permissions.query({ name: "geolocation" }).then((res) => {
    if (res.state === "denied") {
      setCurrentLocationError(
        "Sorry, we could not access your current location. Please enter your location manually."
      );
      return;
    }
    getCurrentLocation();
  });
};

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
  setPlaceId: React.Dispatch<React.SetStateAction<string | null>>;
}> = ({ setPlaceId }) => {
  const [currentLocationError, setCurrentLocationError] = useState<
    string | null
  >(null);

  return (
    <div
      className="px-3 py-2 mt-2 bg-gray-50 rounded border-b hover:bg-primary-100 cursor-pointer"
      onClick={() => {
        requestCurrentLocation(setPlaceId, setCurrentLocationError);
      }}
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
        <div>
          <span>Current Location</span>
        </div>
        {currentLocationError && (
          <div>
            <span className="text-sm text-red-600">{currentLocationError}</span>
          </div>
        )}
      </div>
    </div>
  );
};
