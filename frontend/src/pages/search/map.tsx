import React, { useCallback, useRef, useState } from "react";
import {
  Circle,
  GoogleMap,
  Marker,
  InfoWindow,
  useLoadScript,
} from "@react-google-maps/api";
import { LocationInput, Stop, useLocationSearchQuery } from "client/types";
import { Link } from "react-router-dom";
import { SearchConfig } from "./config";

const libraries = ["geometry", "drawing"];
const mapOptions = {
  mapTypeControl: false,
  streetViewControl: false,
  styles: [
    {
      featureType: "poi.business",
      stylers: [
        {
          visibility: "off",
        },
      ],
    },
    {
      featureType: "poi.park",
      elementType: "labels.text",
      stylers: [
        {
          visibility: "off",
        },
      ],
    },
  ],
};

const initialLocation = {
  latitude: 45.4211,
  longitude: -75.6903,
};

export const SearchMap: React.FC<{ config: SearchConfig }> = ({ config }) => {
  // load google maps
  const { isLoaded, loadError } = useLoadScript({
    googleMapsApiKey: "AIzaSyAvCHchRFUDqVPHSs5jpR74ehIY7A5WBIY",
    libraries: libraries as any,
  });

  // with the initial location / center  as downtown ottawa
  const [location, setLocation] = useState<LocationInput>(initialLocation);

  const [selected, setSelected] = useState<Stop | null>(null);

  // map element ref
  const mapRef = useRef<google.maps.Map>();

  // set map element ref
  const onMapLoad = useCallback((map: google.maps.Map) => {
    mapRef.current = map;
    mapRef.current.setCenter({
      lat: 45.4211,
      lng: -75.6903,
    });
  }, []);

  const onMapDragEnd = useCallback(() => {
    if (!mapRef.current) {
      return;
    }

    let center = mapRef.current.getCenter();
    if (!center) {
      return;
    }

    setLocation({
      latitude: center.lat(),
      longitude: center.lng(),
    });
  }, []);

  const [{ data }, _] = useLocationSearchQuery({
    variables: {
      location: location,
    },
  });

  if (!isLoaded) {
    return <></>;
  }

  if (loadError) {
    return <div>{loadError.message}</div>;
  }

  return (
    <div className="mb-3">
      <div style={{ height: 450 }}>
        <button
          className="absolute text-white bg-indigo-500 hover:bg-indigo-600 z-10 mt-3 ml-3 px-3 py-1 shadow-lg rounded-sm text-sm font-semibold"
          onClick={() => {
            navigator.geolocation.getCurrentPosition(
              (e) => {
                if (!mapRef.current) {
                  return;
                }
                console.log(e.coords);
                setLocation({
                  latitude: e.coords.latitude,
                  longitude: e.coords.longitude,
                });
                mapRef.current.panTo({
                  lat: e.coords.latitude,
                  lng: e.coords.longitude,
                });

                mapRef.current.setZoom(14);
              },
              (e) => {
                console.log(e);
              },
              {
                enableHighAccuracy: true,
                maximumAge: 0,
              }
            );
          }}
        >
          Where's My Location?
        </button>
        <GoogleMap
          onDragEnd={onMapDragEnd}
          onLoad={onMapLoad}
          mapContainerStyle={{
            width: "100%",
            height: "450px",
          }}
          zoom={14}
          center={{
            lat: mapRef.current?.getCenter()?.lat() as any,
            lng: mapRef.current?.getCenter()?.lng() as any,
          }}
          options={{ ...mapOptions, gestureHandling: "greedy" }}
        >
          <Circle
            center={{
              lat: location.latitude,
              lng: location.longitude,
            }}
            options={{
              fillOpacity: 0.2,
              fillColor: "#c7d2fe",
              strokeColor: "#6366f1",
              strokeWeight: 0.75,
              radius: 1000,
            }}
          />
          {data &&
            data.searchStopLocation.map(({ stop }) => (
              <Marker
                key={stop.id}
                position={{
                  lat: stop.location.latitude,
                  lng: stop.location.longitude,
                }}
                icon={{
                  url: "stop.svg",
                  scaledSize: new google.maps.Size(20, 20),
                  origin: new google.maps.Point(0, 0),
                  anchor: new google.maps.Point(10, 10),
                }}
                onClick={(e) => {
                  setSelected(stop as Stop);
                }}
              />
            ))}
          {selected && (
            <InfoWindow
              position={{
                lat: selected.location.latitude,
                lng: selected.location.longitude,
              }}
              onCloseClick={() => setSelected(null)}
            >
              <div>
                <h1 className="text-gray-700 font-semibold tracking-wide">
                  {selected.name} #{selected.code}
                </h1>
                <button
                  onClick={() => {
                    config.action(selected as Stop);
                  }}
                  className="border border-indigo-500 px-5 py-1 mt-2 hover:bg-indigo-500 hover:text-white text-indigo-500 text-sm rounded-sm"
                >
                  {config.actionName}
                </button>
              </div>
            </InfoWindow>
          )}
        </GoogleMap>
      </div>
    </div>
  );
};
