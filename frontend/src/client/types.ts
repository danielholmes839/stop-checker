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

export type PageInfo = {
  __typename?: 'PageInfo';
  cursor: Scalars['Int'];
  remaining: Scalars['Int'];
};

export type PageInput = {
  limit: Scalars['Int'];
  skip: Scalars['Int'];
};

export type Query = {
  __typename?: 'Query';
  searchStopLocation: StopSearchPayload;
  searchStopText: StopSearchPayload;
  stop?: Maybe<Stop>;
  stopRoute?: Maybe<StopRoute>;
  travelPlanner: TravelPayload;
  travelPlannerFixedRoute: TravelPayload;
};


export type QuerySearchStopLocationArgs = {
  location: LocationInput;
  page: PageInput;
  radius: Scalars['Float'];
};


export type QuerySearchStopTextArgs = {
  page: PageInput;
  text: Scalars['String'];
};


export type QueryStopArgs = {
  id: Scalars['ID'];
};


export type QueryStopRouteArgs = {
  route: Scalars['ID'];
  stop: Scalars['ID'];
};


export type QueryTravelPlannerArgs = {
  destination: Scalars['ID'];
  options: TravelScheduleOptions;
  origin: Scalars['ID'];
};


export type QueryTravelPlannerFixedRouteArgs = {
  options: TravelScheduleOptions;
  route: Array<TravelLegInput>;
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

export enum ScheduleMode {
  ArriveBy = 'ARRIVE_BY',
  DepartAt = 'DEPART_AT'
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
  after?: InputMaybe<Scalars['DateTime']>;
  limit: Scalars['Int'];
};


export type StopRouteScheduleOnArgs = {
  date: Scalars['Date'];
};

export type StopSearchPayload = {
  __typename?: 'StopSearchPayload';
  page: PageInfo;
  results: Array<Stop>;
};

export type StopTime = {
  __typename?: 'StopTime';
  id: Scalars['ID'];
  sequence: Scalars['Int'];
  stop: Stop;
  time: Scalars['Time'];
  trip: Trip;
};

export type Transit = {
  __typename?: 'Transit';
  arrival: StopTime;
  departure: StopTime;
  route: Route;
  trip: Trip;
};

export type TravelLegInput = {
  destination: Scalars['ID'];
  origin: Scalars['ID'];
  route?: InputMaybe<Scalars['ID']>;
};

export type TravelPayload = {
  __typename?: 'TravelPayload';
  errors: Array<UserError>;
  schedule?: Maybe<TravelSchedule>;
};

export type TravelSchedule = {
  __typename?: 'TravelSchedule';
  arrival: Scalars['DateTime'];
  departure: Scalars['DateTime'];
  destination: Stop;
  duration: Scalars['Int'];
  legs: Array<TravelScheduleLeg>;
  origin: Stop;
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

export type TravelScheduleOptions = {
  datetime?: InputMaybe<Scalars['DateTime']>;
  mode: ScheduleMode;
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

export type UserError = {
  __typename?: 'UserError';
  field: Scalars['String'];
  message: Scalars['String'];
};

export type LocationSearchQueryVariables = Exact<{
  location: LocationInput;
  page: PageInput;
}>;


export type LocationSearchQuery = { __typename?: 'Query', searchStopLocation: { __typename?: 'StopSearchPayload', results: Array<{ __typename?: 'Stop', id: string, name: string, code: string, location: { __typename?: 'Location', latitude: number, longitude: number } }> } };

export type StopPageQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type StopPageQuery = { __typename?: 'Query', stop?: { __typename?: 'Stop', id: string, name: string, code: string, routes: Array<{ __typename?: 'StopRoute', headsign: string, route: { __typename?: 'Route', name: string, text: any, background: any }, schedule: { __typename?: 'StopRouteSchedule', next: Array<{ __typename?: 'StopTime', id: string, time: any }> } }> } | null };

export type StopPreviewQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type StopPreviewQuery = { __typename?: 'Query', stop?: { __typename?: 'Stop', id: string, name: string, code: string, routes: Array<{ __typename?: 'StopRoute', headsign: string, route: { __typename?: 'Route', name: string, text: any, background: any } }> } | null };

export type TextSearchQueryVariables = Exact<{
  text: Scalars['String'];
  page: PageInput;
}>;


export type TextSearchQuery = { __typename?: 'Query', searchStopText: { __typename?: 'StopSearchPayload', results: Array<{ __typename?: 'Stop', id: string, name: string, code: string, routes: Array<{ __typename?: 'StopRoute', direction: string, headsign: string, route: { __typename?: 'Route', name: string, background: any, text: any } }> }> } };

export type TravelPlannerQueryVariables = Exact<{
  origin: Scalars['ID'];
  destination: Scalars['ID'];
  options: TravelScheduleOptions;
}>;


export type TravelPlannerQuery = { __typename?: 'Query', travelPlanner: { __typename?: 'TravelPayload', errors: Array<{ __typename?: 'UserError', field: string, message: string }>, schedule?: { __typename?: 'TravelSchedule', departure: any, arrival: any, duration: number, legs: Array<{ __typename?: 'TravelScheduleLeg', departure: any, arrival: any, duration: number, walk: boolean, distance: number, origin: { __typename?: 'Stop', id: string, name: string, code: string }, destination: { __typename?: 'Stop', id: string, name: string, code: string }, transit?: { __typename?: 'Transit', route: { __typename?: 'Route', id: string, name: string, type: RouteType, text: any, background: any }, trip: { __typename?: 'Trip', headsign: string, stopTimes: Array<{ __typename?: 'StopTime', id: string, sequence: number, time: any, stop: { __typename?: 'Stop', name: string } }> }, arrival: { __typename?: 'StopTime', sequence: number }, departure: { __typename?: 'StopTime', sequence: number } } | null }> } | null } };

export type TravelScheduleLegDefaultFragment = { __typename?: 'TravelScheduleLeg', departure: any, arrival: any, duration: number, walk: boolean, distance: number, origin: { __typename?: 'Stop', id: string, name: string, code: string }, destination: { __typename?: 'Stop', id: string, name: string, code: string }, transit?: { __typename?: 'Transit', route: { __typename?: 'Route', id: string, name: string, type: RouteType, text: any, background: any }, trip: { __typename?: 'Trip', headsign: string, stopTimes: Array<{ __typename?: 'StopTime', id: string, sequence: number, time: any, stop: { __typename?: 'Stop', name: string } }> }, arrival: { __typename?: 'StopTime', sequence: number }, departure: { __typename?: 'StopTime', sequence: number } } | null };

export type TravelPlannerDeparturesQueryVariables = Exact<{
  stop: Scalars['ID'];
  route: Scalars['ID'];
  after: Scalars['DateTime'];
  limit: Scalars['Int'];
}>;


export type TravelPlannerDeparturesQuery = { __typename?: 'Query', stopRoute?: { __typename?: 'StopRoute', schedule: { __typename?: 'StopRouteSchedule', next: Array<{ __typename?: 'StopTime', id: string, time: any }> } } | null };

export const TravelScheduleLegDefaultFragmentDoc = gql`
    fragment TravelScheduleLegDefault on TravelScheduleLeg {
  departure
  arrival
  duration
  origin {
    id
    name
    code
  }
  destination {
    id
    name
    code
  }
  walk
  distance
  transit {
    route {
      id
      name
      type
      text
      background
    }
    trip {
      headsign
      stopTimes {
        id
        sequence
        time
        stop {
          name
        }
      }
    }
    arrival {
      sequence
    }
    departure {
      sequence
    }
  }
}
    `;
export const LocationSearchDocument = gql`
    query LocationSearch($location: LocationInput!, $page: PageInput!) {
  searchStopLocation(location: $location, radius: 1000, page: $page) {
    results {
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
export const StopPreviewDocument = gql`
    query StopPreview($id: ID!) {
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
    }
  }
}
    `;

export function useStopPreviewQuery(options: Omit<Urql.UseQueryArgs<StopPreviewQueryVariables>, 'query'>) {
  return Urql.useQuery<StopPreviewQuery, StopPreviewQueryVariables>({ query: StopPreviewDocument, ...options });
};
export const TextSearchDocument = gql`
    query TextSearch($text: String!, $page: PageInput!) {
  searchStopText(text: $text, page: $page) {
    results {
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
}
    `;

export function useTextSearchQuery(options: Omit<Urql.UseQueryArgs<TextSearchQueryVariables>, 'query'>) {
  return Urql.useQuery<TextSearchQuery, TextSearchQueryVariables>({ query: TextSearchDocument, ...options });
};
export const TravelPlannerDocument = gql`
    query TravelPlanner($origin: ID!, $destination: ID!, $options: TravelScheduleOptions!) {
  travelPlanner(origin: $origin, destination: $destination, options: $options) {
    errors {
      field
      message
    }
    schedule {
      legs {
        ...TravelScheduleLegDefault
      }
      departure
      arrival
      duration
    }
  }
}
    ${TravelScheduleLegDefaultFragmentDoc}`;

export function useTravelPlannerQuery(options: Omit<Urql.UseQueryArgs<TravelPlannerQueryVariables>, 'query'>) {
  return Urql.useQuery<TravelPlannerQuery, TravelPlannerQueryVariables>({ query: TravelPlannerDocument, ...options });
};
export const TravelPlannerDeparturesDocument = gql`
    query TravelPlannerDepartures($stop: ID!, $route: ID!, $after: DateTime!, $limit: Int!) {
  stopRoute(stop: $stop, route: $route) {
    schedule {
      next(limit: $limit, after: $after) {
        id
        time
      }
    }
  }
}
    `;

export function useTravelPlannerDeparturesQuery(options: Omit<Urql.UseQueryArgs<TravelPlannerDeparturesQueryVariables>, 'query'>) {
  return Urql.useQuery<TravelPlannerDeparturesQuery, TravelPlannerDeparturesQueryVariables>({ query: TravelPlannerDeparturesDocument, ...options });
};