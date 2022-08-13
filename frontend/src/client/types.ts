import gql from 'graphql-tag';
import * as Urql from 'urql';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  Color: any;
  Date: any;
  DateTime: any;
  Time: any;
};

export type Error = {
  __typename?: 'Error';
  field: Scalars['String'];
  message: Scalars['String'];
};

export type Location = {
  __typename?: 'Location';
  distance: Scalars['Float'];
  latitude: Scalars['Float'];
  longitude: Scalars['Float'];
};


export type LocationDistanceArgs = {
  location: LocationInput;
};

export type LocationInput = {
  latitude: Scalars['Float'];
  longitude: Scalars['Float'];
};

export type PageInput = {
  limit: Scalars['Int'];
  skip: Scalars['Int'];
};

export type Query = {
  __typename?: 'Query';
  searchStopLocation: Array<StopLocationResult>;
  searchStopText: Array<Stop>;
  stop?: Maybe<Stop>;
  travelRoutePlanner: TravelRoutePayload;
  travelSchedulePlanner: TravelSchedulePayload;
};


export type QuerySearchStopLocationArgs = {
  location: LocationInput;
  radius: Scalars['Float'];
};


export type QuerySearchStopTextArgs = {
  text: Scalars['String'];
};


export type QueryStopArgs = {
  id: Scalars['ID'];
};


export type QueryTravelRoutePlannerArgs = {
  input: TravelRoutePlannerInput;
};


export type QueryTravelSchedulePlannerArgs = {
  input: TravelSchedulePlannerInput;
};

export type Route = {
  __typename?: 'Route';
  background: Scalars['Color'];
  id: Scalars['ID'];
  name: Scalars['String'];
  text: Scalars['Color'];
  type: RouteType;
};

export enum RouteType {
  Bus = 'BUS',
  Train = 'TRAIN'
}

export type Service = {
  __typename?: 'Service';
  end: Scalars['Date'];
  exceptions: Array<ServiceException>;
  friday: Scalars['Boolean'];
  monday: Scalars['Boolean'];
  saturday: Scalars['Boolean'];
  start: Scalars['Date'];
  sunday: Scalars['Boolean'];
  thursday: Scalars['Boolean'];
  tuesday: Scalars['Boolean'];
  wednesday: Scalars['Boolean'];
};

export type ServiceException = {
  __typename?: 'ServiceException';
  added: Scalars['Boolean'];
  date: Scalars['Date'];
};

export type Stop = {
  __typename?: 'Stop';
  code: Scalars['String'];
  id: Scalars['ID'];
  location: Location;
  name: Scalars['String'];
  routes: Array<StopRoute>;
};

export type StopLocationResult = {
  __typename?: 'StopLocationResult';
  distance: Scalars['Float'];
  stop: Stop;
};

export type StopRoute = {
  __typename?: 'StopRoute';
  direction: Scalars['ID'];
  headsign: Scalars['String'];
  route: Route;
  schedule: StopRouteSchedule;
  stop: Stop;
};

export type StopRouteSchedule = {
  __typename?: 'StopRouteSchedule';
  next: Array<StopTime>;
  on: Array<StopTime>;
};


export type StopRouteScheduleNextArgs = {
  limit: Scalars['Int'];
};


export type StopRouteScheduleOnArgs = {
  date: Scalars['Date'];
};

export type StopTime = {
  __typename?: 'StopTime';
  id: Scalars['ID'];
  sequence: Scalars['Int'];
  time: Scalars['Time'];
  trip: Trip;
};

export type Transit = {
  __typename?: 'Transit';
  route: Route;
  trip: Trip;
};

export type TravelLegInput = {
  destination: Scalars['ID'];
  origin: Scalars['ID'];
  route?: InputMaybe<Scalars['ID']>;
};

export type TravelRoute = {
  __typename?: 'TravelRoute';
  legs: Array<TravelRouteLeg>;
};

export type TravelRouteLeg = {
  __typename?: 'TravelRouteLeg';
  destination: Stop;
  distance: Scalars['Float'];
  origin: Stop;
  route?: Maybe<Route>;
  walk: Scalars['Boolean'];
};

export type TravelRoutePayload = {
  __typename?: 'TravelRoutePayload';
  errors: Array<Error>;
  route?: Maybe<TravelRoute>;
};

export type TravelRoutePlannerInput = {
  departure?: InputMaybe<Scalars['DateTime']>;
  destination: Scalars['ID'];
  origin: Scalars['ID'];
};

export type TravelSchedule = {
  __typename?: 'TravelSchedule';
  arrival: Scalars['DateTime'];
  departure: Scalars['DateTime'];
  duration: Scalars['Int'];
  legs: Array<TravelScheduleLeg>;
};

export type TravelScheduleLeg = {
  __typename?: 'TravelScheduleLeg';
  arrival: Scalars['DateTime'];
  departure: Scalars['DateTime'];
  destination: Stop;
  distance: Scalars['Float'];
  duration: Scalars['Int'];
  origin: Stop;
  transit?: Maybe<Transit>;
  walk: Scalars['Boolean'];
};

export type TravelSchedulePayload = {
  __typename?: 'TravelSchedulePayload';
  errors: Array<Error>;
  schedule?: Maybe<TravelSchedule>;
};

export type TravelSchedulePlannerInput = {
  arrival?: InputMaybe<Scalars['DateTime']>;
  departure?: InputMaybe<Scalars['DateTime']>;
  legs: Array<TravelLegInput>;
};

export type Trip = {
  __typename?: 'Trip';
  direction: Scalars['ID'];
  headsign: Scalars['String'];
  id: Scalars['ID'];
  route: Route;
  service: Service;
  stopTimes: Array<StopTime>;
};

export type LocationSearchQueryVariables = Exact<{
  location: LocationInput;
}>;


export type LocationSearchQuery = { __typename?: 'Query', searchStopLocation: Array<{ __typename?: 'StopLocationResult', distance: number, stop: { __typename?: 'Stop', id: string, name: string, code: string, location: { __typename?: 'Location', latitude: number, longitude: number } } }> };

export type StopPageQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type StopPageQuery = { __typename?: 'Query', stop?: { __typename?: 'Stop', id: string, name: string, code: string, routes: Array<{ __typename?: 'StopRoute', headsign: string, route: { __typename?: 'Route', name: string, text: any, background: any }, schedule: { __typename?: 'StopRouteSchedule', next: Array<{ __typename?: 'StopTime', id: string, time: any }> } }> } | null };

export type TextSearchQueryVariables = Exact<{
  text: Scalars['String'];
}>;


export type TextSearchQuery = { __typename?: 'Query', searchStopText: Array<{ __typename?: 'Stop', id: string, name: string, code: string, routes: Array<{ __typename?: 'StopRoute', direction: string, headsign: string, route: { __typename?: 'Route', name: string, background: any, text: any } }> }> };


export const LocationSearchDocument = gql`
    query LocationSearch($location: LocationInput!) {
  searchStopLocation(location: $location, radius: 1000) {
    distance
    stop {
      id
      name
      code
      location {
        latitude
        longitude
      }
    }
  }
}
    `;

export function useLocationSearchQuery(options: Omit<Urql.UseQueryArgs<LocationSearchQueryVariables>, 'query'>) {
  return Urql.useQuery<LocationSearchQuery, LocationSearchQueryVariables>({ query: LocationSearchDocument, ...options });
};
export const StopPageDocument = gql`
    query StopPage($id: ID!) {
  stop(id: $id) {
    id
    name
    code
    routes {
      headsign
      route {
        name
        text
        background
      }
      schedule {
        next(limit: 3) {
          id
          time
        }
      }
    }
  }
}
    `;

export function useStopPageQuery(options: Omit<Urql.UseQueryArgs<StopPageQueryVariables>, 'query'>) {
  return Urql.useQuery<StopPageQuery, StopPageQueryVariables>({ query: StopPageDocument, ...options });
};
export const TextSearchDocument = gql`
    query TextSearch($text: String!) {
  searchStopText(text: $text) {
    id
    name
    code
    routes {
      direction
      headsign
      route {
        name
        background
        text
      }
    }
  }
}
    `;

export function useTextSearchQuery(options: Omit<Urql.UseQueryArgs<TextSearchQueryVariables>, 'query'>) {
  return Urql.useQuery<TextSearchQuery, TextSearchQueryVariables>({ query: TextSearchDocument, ...options });
};