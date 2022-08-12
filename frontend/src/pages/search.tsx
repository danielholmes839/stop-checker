import {
  LocationInput,
  StopPageQuery,
  TextSearchQuery,
  useLocationSearchQuery,
  useStopPageQuery,
  useTextSearchQuery,
} from "client/types";
import { Sign, Container, Card, QueryResponseWrapper } from "components";
import React, { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useDebounce } from "use-debounce";

import {
  Circle,
  withGoogleMap,
  withScriptjs,
  GoogleMap,
  Marker,
  InfoWindow,
} from "react-google-maps";

const Map: React.FC = () => {
  const [location, setLocation] = useState<LocationInput>({
    latitude: 45.4211,
    longitude: -75.6903,
  });

  const [locationDebounced, __] = useDebounce(location, 10);

  const [{ data }, _] = useLocationSearchQuery({
    variables: {
      location: locationDebounced,
    },
  });

  return (
    <>
      <GoogleMap
        defaultZoom={10}
        defaultCenter={{ lat: 45.4211, lng: -75.6903 }}
        options={{
          maxZoom: 18,
          minZoom: 9,
          streetViewControl: false,
          mapTypeControl: false,
          mapTypeId: google.maps.MapTypeId.ROADMAP,
          scrollwheel: true,
        }}
        onClick={(e) => {
          setLocation({
            latitude: e.latLng.lat(),
            longitude: e.latLng.lng(),
          });
        }}
      >
        <Circle
          key={1}
          options={{
            center: {
              lat: location.latitude,
              lng: location.longitude,
            },
            radius: 1500,
            fillColor: "#c7d2fe",
            fillOpacity: 0.2,
            strokeWeight: 1,
            strokeColor: "#6366f1",
          }}
        />
        {data &&
          data.searchStopLocation.map(({ stop }) => (
            <Marker
              position={{
                lat: stop.location.latitude,
                lng: stop.location.longitude,
              }}
              key={stop.id}
            ></Marker>
          ))}
      </GoogleMap>
    </>
  );
};

const MapWrapped = withScriptjs(withGoogleMap(Map));

const SearchMap: React.FC = () => {
  return (
    <div style={{ width: "100%", height: "400px" }}>
      <MapWrapped
        googleMapURL={`https://maps.googleapis.com/maps/api/js?v=3.exp&libraries=geometry,drawing,places&key=AIzaSyAvCHchRFUDqVPHSs5jpR74ehIY7A5WBIY`}
        loadingElement={<div style={{ height: `100%` }} />}
        containerElement={<div style={{ height: `100%` }} />}
        mapElement={<div style={{ height: `100%` }} />}
      />
    </div>
  );
};

const TextSearchQueryResponse: React.FC<{ data: TextSearchQuery }> = ({
  data,
}) => {
  const navigate = useNavigate();
  return (
    <div>
      {data.searchStopText.map(({ id, name, code, routes }) => (
        <div className="mb-3 cursor-pointer hover:shadow-sm">
          <Card key={id} onClick={() => navigate(`/stop/${id}`)}>
            <h1 className="font-semibold mb-1">
              {name} <span className="float-right text-gray-700">#{code}</span>
            </h1>

            <div>
              {routes.map(({ direction, headsign, route }) => (
                <Sign
                  key={route.name + direction}
                  props={{
                    background: route.background,
                    headsign: headsign,
                    name: route.name,
                    text: route.text,
                  }}
                />
              ))}
            </div>
          </Card>
        </div>
      ))}
    </div>
  );
};

export const SearchPage: React.FC = () => {
  const [search, setTextSearch] = useState("");
  const [searchDebounced] = useDebounce(search, 200);

  const [response, _fetch] = useTextSearchQuery({
    variables: {
      text: searchDebounced,
    },
  });

  return (
    <Container>
      <input
        className="bg-gray-50 border-gray-100 w-full mb-3 mt-1 p-3 rounded-sm focus:outline-none border-b-2 focus:border-indigo-500"
        value={search}
        type="text"
        placeholder="Search by stop name or code Ex. Rideau A, O-Train, 3000"
        onChange={(event) => setTextSearch(event.target.value)}
      />
      <SearchMap />
      <QueryResponseWrapper response={response}>
        {response.data && <TextSearchQueryResponse data={response.data} />}
      </QueryResponseWrapper>
    </Container>
  );
};

const StopPageResponse: React.FC<{ data: StopPageQuery }> = ({ data }) => {
  const { stop } = data;
  if (stop == undefined) {
    return <div>stop not found</div>;
  }

  return (
    <div className="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
      {stop.routes
        .sort((a, b) => a.route.name.localeCompare(b.route.name))
        .map((stopRoute) => {
          const { route, headsign, schedule } = stopRoute;
          return (
            <Card>
              <Sign
                props={{
                  background: route.background,
                  text: route.text,
                  headsign: headsign,
                  name: route.name,
                }}
              />
              <div>
                {schedule.next.length > 0 ? (
                  schedule.next.map((stopTime) => (
                    <span className="mr-2 text-sm">{stopTime.time}</span>
                  ))
                ) : (
                  <span className="text-sm">
                    No more stops today or tomorrow
                  </span>
                )}
              </div>
            </Card>
          );
        })}
    </div>
  );
};

export const StopPage: React.FC = () => {
  const { id } = useParams();
  const [response, _] = useStopPageQuery({
    variables: { id: id === undefined ? "" : id },
  });

  return (
    <Container>
      <QueryResponseWrapper response={response}>
        {response.data && <StopPageResponse data={response.data} />}
      </QueryResponseWrapper>
    </Container>
  );
};
