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

export type Bus = {
  __typename?: 'Bus';
  arrival: Scalars['Time'];
  distance?: Maybe<Scalars['Float']>;
  headsign: Scalars['String'];
  lastUpdated: Scalars['Time'];
  lastUpdatedMessage: Scalars['String'];
  lastUpdatedMinutes: Scalars['Int'];
  location?: Maybe<Location>;
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
  travelPlanner: TravelSchedulePayload;
  travelPlannerFixedRoute: TravelSchedulePayload;
  travelRoute: TravelRoutePayload;
};


export type QuerySearchStopLocationArgs = {
  location: LocationInput;
  page: PageInput;
  radius: Scalars['Float'];
  sorted: Scalars['Boolean'];
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
  options: TravelOptions;
  origin: Scalars['ID'];
};


export type QueryTravelPlannerFixedRouteArgs = {
  input: Array<TravelLegInput>;
  options: TravelOptions;
};


export type QueryTravelRouteArgs = {
  input?: InputMaybe<Array<TravelLegInput>>;
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

export type ScheduleResult = {
  __typename?: 'ScheduleResult';
  datetime: Scalars['DateTime'];
  stoptime: StopTime;
};

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
  liveBuses: Array<Bus>;
  liveMap?: Maybe<Scalars['String']>;
  reaches: Array<Stop>;
  route: Route;
  schedule: StopRouteSchedule;
  scheduleReaches?: Maybe<StopRouteSchedule>;
  stop: Stop;
};


export type StopRouteReachesArgs = {
  forward: Scalars['Boolean'];
};


export type StopRouteScheduleReachesArgs = {
  destination: Scalars['ID'];
};

export type StopRouteSchedule = {
  __typename?: 'StopRouteSchedule';
  next: Array<ScheduleResult>;
  on: Array<ScheduleResult>;
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
  overflow: Scalars['Boolean'];
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

export type TravelOptions = {
  datetime?: InputMaybe<Scalars['DateTime']>;
  mode: ScheduleMode;
};

export type TravelRouteLeg = {
  __typename?: 'TravelRouteLeg';
  destination: Stop;
  distance: Scalars['Float'];
  origin: Stop;
  stopRoute?: Maybe<StopRoute>;
  walk: Scalars['Boolean'];
};

export type TravelRoutePayload = {
  __typename?: 'TravelRoutePayload';
  error?: Maybe<Scalars['String']>;
  route?: Maybe<Array<TravelRouteLeg>>;
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
  shape: Array<Location>;
  transit?: Maybe<Transit>;
  walk: Scalars['Boolean'];
};

export type TravelSchedulePayload = {
  __typename?: 'TravelSchedulePayload';
  error?: Maybe<Scalars['String']>;
  schedule?: Maybe<TravelSchedule>;
};

export type Trip = {
  __typename?: 'Trip';
  direction: Scalars['ID'];
  headsign: Scalars['String'];
  id: Scalars['ID'];
  route: Route;
  service: Service;
  shape: Array<Location>;
  stoptimes: Array<StopTime>;
};

export type LocationSearchQueryVariables = Exact<{
  location: LocationInput;
  page: PageInput;
}>;


export type LocationSearchQuery = { __typename?: 'Query', searchStopLocation: { __typename?: 'StopSearchPayload', results: Array<{ __typename?: 'Stop', id: string, name: string, code: string, location: { __typename?: 'Location', latitude: number, longitude: number } }> } };

export type StopPageQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type StopPageQuery = { __typename?: 'Query', stop?: { __typename?: 'Stop', id: string, name: string, code: string, routes: Array<{ __typename?: 'StopRoute', headsign: string, route: { __typename?: 'Route', id: string, name: string, text: any, background: any }, schedule: { __typename?: 'StopRouteSchedule', next: Array<{ __typename?: 'ScheduleResult', stoptime: { __typename?: 'StopTime', id: string, time: any } }> }, liveBuses: Array<{ __typename?: 'Bus', headsign: string }> }> } | null };

export type StopExploreQueryVariables = Exact<{
  origin: Scalars['ID'];
}>;


export type StopExploreQuery = { __typename?: 'Query', stop?: { __typename?: 'Stop', id: string, name: string, code: string, routes: Array<{ __typename?: 'StopRoute', headsign: string, route: { __typename?: 'Route', id: string, name: string, text: any, background: any }, destinations: Array<{ __typename?: 'Stop', id: string, name: string, code: string }> }>, location: { __typename?: 'Location', latitude: number, longitude: number } } | null };

export type StopExploreWalkQueryVariables = Exact<{
  location: LocationInput;
}>;


export type StopExploreWalkQuery = { __typename?: 'Query', searchStopLocation: { __typename?: 'StopSearchPayload', page: { __typename?: 'PageInfo', cursor: number, remaining: number }, results: Array<{ __typename?: 'Stop', id: string, name: string, code: string, location: { __typename?: 'Location', latitude: number, longitude: number, distance: number }, routes: Array<{ __typename?: 'StopRoute', headsign: string, route: { __typename?: 'Route', id: string, name: string, text: any, background: any } }> }> } };

export type StopExploreFragment = { __typename?: 'Stop', id: string, name: string, code: string, location: { __typename?: 'Location', latitude: number, longitude: number } };

export type StopRouteExploreFragment = { __typename?: 'StopRoute', headsign: string, route: { __typename?: 'Route', id: string, name: string, text: any, background: any }, destinations: Array<{ __typename?: 'Stop', id: string, name: string, code: string }> };

export type StopPreviewQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type StopPreviewQuery = { __typename?: 'Query', stop?: { __typename?: 'Stop', id: string, name: string, code: string, routes: Array<{ __typename?: 'StopRoute', headsign: string, route: { __typename?: 'Route', id: string, name: string, text: any, background: any } }> } | null };

export type StopPreviewFragment = { __typename?: 'Stop', id: string, name: string, code: string, routes: Array<{ __typename?: 'StopRoute', headsign: string, route: { __typename?: 'Route', id: string, name: string, text: any, background: any } }> };

export type StopRouteQueryVariables = Exact<{
  stop: Scalars['ID'];
  route: Scalars['ID'];
  today: Scalars['Date'];
  tomorrow: Scalars['Date'];
}>;


export type StopRouteQuery = { __typename?: 'Query', stopRoute?: { __typename?: 'StopRoute', headsign: string, liveMap?: string | null, stop: { __typename?: 'Stop', id: string, name: string, code: string }, route: { __typename?: 'Route', id: string, name: string, text: any, background: any }, schedule: { __typename?: 'StopRouteSchedule', today: Array<{ __typename?: 'ScheduleResult', stoptime: { __typename?: 'StopTime', id: string, time: any } }>, tomorrow: Array<{ __typename?: 'ScheduleResult', stoptime: { __typename?: 'StopTime', id: string, time: any } }>, next: Array<{ __typename?: 'ScheduleResult', stoptime: { __typename?: 'StopTime', id: string, time: any, overflow: boolean, trip: { __typename?: 'Trip', service: { __typename?: 'Service', saturday: boolean, sunday: boolean, monday: boolean, start: any, end: any } } } }> }, liveBuses: Array<{ __typename?: 'Bus', headsign: string, arrival: any, lastUpdated: any, lastUpdatedMessage: string, lastUpdatedMinutes: number, distance?: number | null, location?: { __typename?: 'Location', latitude: number, longitude: number } | null }> } | null };

export type LiveDataFragment = { __typename?: 'StopRoute', liveMap?: string | null, liveBuses: Array<{ __typename?: 'Bus', headsign: string, arrival: any, lastUpdated: any, lastUpdatedMessage: string, lastUpdatedMinutes: number, distance?: number | null, location?: { __typename?: 'Location', latitude: number, longitude: number } | null }> };

export type TextSearchQueryVariables = Exact<{
  text: Scalars['String'];
  page: PageInput;
}>;


export type TextSearchQuery = { __typename?: 'Query', searchStopText: { __typename?: 'StopSearchPayload', results: Array<{ __typename?: 'Stop', id: string, name: string, code: string, routes: Array<{ __typename?: 'StopRoute', headsign: string, route: { __typename?: 'Route', id: string, name: string, text: any, background: any } }> }> } };

export type TravelPlannerQueryVariables = Exact<{
  origin: Scalars['ID'];
  destination: Scalars['ID'];
  options: TravelOptions;
}>;


export type TravelPlannerQuery = { __typename?: 'Query', travelPlanner: { __typename?: 'TravelSchedulePayload', error?: string | null, schedule?: { __typename?: 'TravelSchedule', departure: any, arrival: any, duration: number, legs: Array<{ __typename?: 'TravelScheduleLeg', departure: any, arrival: any, duration: number, walk: boolean, distance: number, origin: { __typename?: 'Stop', id: string, name: string, code: string, location: { __typename?: 'Location', latitude: number, longitude: number } }, destination: { __typename?: 'Stop', id: string, name: string, code: string }, shape: Array<{ __typename?: 'Location', latitude: number, longitude: number }>, transit?: { __typename?: 'Transit', route: { __typename?: 'Route', id: string, name: string, type: RouteType, text: any, background: any }, trip: { __typename?: 'Trip', headsign: string, stoptimes: Array<{ __typename?: 'StopTime', id: string, sequence: number, time: any, stop: { __typename?: 'Stop', name: string } }> }, arrival: { __typename?: 'StopTime', sequence: number }, departure: { __typename?: 'StopTime', sequence: number } } | null }> } | null } };

export type TravelPlannerFixedRouteQueryVariables = Exact<{
  input: Array<TravelLegInput> | TravelLegInput;
  options: TravelOptions;
}>;


export type TravelPlannerFixedRouteQuery = { __typename?: 'Query', travelPlannerFixedRoute: { __typename?: 'TravelSchedulePayload', error?: string | null, schedule?: { __typename?: 'TravelSchedule', departure: any, arrival: any, duration: number, legs: Array<{ __typename?: 'TravelScheduleLeg', departure: any, arrival: any, duration: number, walk: boolean, distance: number, origin: { __typename?: 'Stop', id: string, name: string, code: string, location: { __typename?: 'Location', latitude: number, longitude: number } }, destination: { __typename?: 'Stop', id: string, name: string, code: string }, shape: Array<{ __typename?: 'Location', latitude: number, longitude: number }>, transit?: { __typename?: 'Transit', route: { __typename?: 'Route', id: string, name: string, type: RouteType, text: any, background: any }, trip: { __typename?: 'Trip', headsign: string, stoptimes: Array<{ __typename?: 'StopTime', id: string, sequence: number, time: any, stop: { __typename?: 'Stop', name: string } }> }, arrival: { __typename?: 'StopTime', sequence: number }, departure: { __typename?: 'StopTime', sequence: number } } | null }> } | null } };

export type TravelScheduleFragment = { __typename?: 'TravelSchedulePayload', error?: string | null, schedule?: { __typename?: 'TravelSchedule', departure: any, arrival: any, duration: number, legs: Array<{ __typename?: 'TravelScheduleLeg', departure: any, arrival: any, duration: number, walk: boolean, distance: number, origin: { __typename?: 'Stop', id: string, name: string, code: string, location: { __typename?: 'Location', latitude: number, longitude: number } }, destination: { __typename?: 'Stop', id: string, name: string, code: string }, shape: Array<{ __typename?: 'Location', latitude: number, longitude: number }>, transit?: { __typename?: 'Transit', route: { __typename?: 'Route', id: string, name: string, type: RouteType, text: any, background: any }, trip: { __typename?: 'Trip', headsign: string, stoptimes: Array<{ __typename?: 'StopTime', id: string, sequence: number, time: any, stop: { __typename?: 'Stop', name: string } }> }, arrival: { __typename?: 'StopTime', sequence: number }, departure: { __typename?: 'StopTime', sequence: number } } | null }> } | null };

export type TravelScheduleLegDefaultFragment = { __typename?: 'TravelScheduleLeg', departure: any, arrival: any, duration: number, walk: boolean, distance: number, origin: { __typename?: 'Stop', id: string, name: string, code: string, location: { __typename?: 'Location', latitude: number, longitude: number } }, destination: { __typename?: 'Stop', id: string, name: string, code: string }, shape: Array<{ __typename?: 'Location', latitude: number, longitude: number }>, transit?: { __typename?: 'Transit', route: { __typename?: 'Route', id: string, name: string, type: RouteType, text: any, background: any }, trip: { __typename?: 'Trip', headsign: string, stoptimes: Array<{ __typename?: 'StopTime', id: string, sequence: number, time: any, stop: { __typename?: 'Stop', name: string } }> }, arrival: { __typename?: 'StopTime', sequence: number }, departure: { __typename?: 'StopTime', sequence: number } } | null };

export type TravelPlannerDeparturesQueryVariables = Exact<{
  origin: Scalars['ID'];
  destination: Scalars['ID'];
  route: Scalars['ID'];
  after: Scalars['DateTime'];
  limit: Scalars['Int'];
}>;


export type TravelPlannerDeparturesQuery = { __typename?: 'Query', stopRoute?: { __typename?: 'StopRoute', schedule?: { __typename?: 'StopRouteSchedule', next: Array<{ __typename?: 'ScheduleResult', stoptime: { __typename?: 'StopTime', id: string, time: any } }> } | null } | null };

export type TravelRouteQueryVariables = Exact<{
  input: Array<TravelLegInput> | TravelLegInput;
}>;


export type TravelRouteQuery = { __typename?: 'Query', travelRoute: { __typename?: 'TravelRoutePayload', route?: Array<{ __typename?: 'TravelRouteLeg', walk: boolean, distance: number, origin: { __typename?: 'Stop', id: string, name: string, code: string }, destination: { __typename?: 'Stop', id: string, name: string, code: string }, stopRoute?: { __typename?: 'StopRoute', headsign: string, route: { __typename?: 'Route', id: string, name: string, text: any, background: any } } | null }> | null } };

export const StopExploreFragmentDoc = gql`
    fragment StopExplore on Stop {
  id
  name
  code
  location {
    latitude
    longitude
  }
}
    `;
export const StopRouteExploreFragmentDoc = gql`
    fragment StopRouteExplore on StopRoute {
  headsign
  route {
    id
    name
    text
    background
  }
  destinations: reaches(forward: true) {
    id
    name
    code
  }
}
    `;
export const StopPreviewFragmentDoc = gql`
    fragment StopPreview on Stop {
  id
  name
  code
  routes {
    headsign
    route {
      id
      name
      text
      background
    }
  }
}
    `;
export const LiveDataFragmentDoc = gql`
    fragment LiveData on StopRoute {
  liveMap
  liveBuses {
    headsign
    arrival
    lastUpdated
    lastUpdatedMessage
    lastUpdatedMinutes
    distance
    location {
      latitude
      longitude
    }
  }
}
    `;
export const TravelScheduleLegDefaultFragmentDoc = gql`
    fragment TravelScheduleLegDefault on TravelScheduleLeg {
  departure
  arrival
  duration
  origin {
    id
    name
    code
    location {
      latitude
      longitude
    }
  }
  destination {
    id
    name
    code
  }
  walk
  distance
  shape {
    latitude
    longitude
  }
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
      stoptimes {
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
export const TravelScheduleFragmentDoc = gql`
    fragment TravelSchedule on TravelSchedulePayload {
  error
  schedule {
    legs {
      ...TravelScheduleLegDefault
    }
    departure
    arrival
    duration
  }
}
    ${TravelScheduleLegDefaultFragmentDoc}`;
export const LocationSearchDocument = gql`
    query LocationSearch($location: LocationInput!, $page: PageInput!) {
  searchStopLocation(
    location: $location
    radius: 1000
    page: $page
    sorted: false
  ) {
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
        id
        name
        text
        background
      }
      schedule {
        next(limit: 3) {
          stoptime {
            id
            time
          }
        }
      }
      liveBuses {
        headsign
      }
    }
  }
}
    `;

export function useStopPageQuery(options: Omit<Urql.UseQueryArgs<StopPageQueryVariables>, 'query'>) {
  return Urql.useQuery<StopPageQuery, StopPageQueryVariables>({ query: StopPageDocument, ...options });
};
export const StopExploreDocument = gql`
    query StopExplore($origin: ID!) {
  stop(id: $origin) {
    ...StopExplore
    routes {
      ...StopRouteExplore
    }
  }
}
    ${StopExploreFragmentDoc}
${StopRouteExploreFragmentDoc}`;

export function useStopExploreQuery(options: Omit<Urql.UseQueryArgs<StopExploreQueryVariables>, 'query'>) {
  return Urql.useQuery<StopExploreQuery, StopExploreQueryVariables>({ query: StopExploreDocument, ...options });
};
export const StopExploreWalkDocument = gql`
    query StopExploreWalk($location: LocationInput!) {
  searchStopLocation(
    location: $location
    radius: 500
    page: {skip: 0, limit: -1}
    sorted: true
  ) {
    page {
      cursor
      remaining
    }
    results {
      id
      name
      code
      location {
        latitude
        longitude
        distance(location: $location)
      }
      routes {
        headsign
        route {
          id
          name
          text
          background
        }
      }
    }
  }
}
    `;

export function useStopExploreWalkQuery(options: Omit<Urql.UseQueryArgs<StopExploreWalkQueryVariables>, 'query'>) {
  return Urql.useQuery<StopExploreWalkQuery, StopExploreWalkQueryVariables>({ query: StopExploreWalkDocument, ...options });
};
export const StopPreviewDocument = gql`
    query StopPreview($id: ID!) {
  stop(id: $id) {
    ...StopPreview
  }
}
    ${StopPreviewFragmentDoc}`;

export function useStopPreviewQuery(options: Omit<Urql.UseQueryArgs<StopPreviewQueryVariables>, 'query'>) {
  return Urql.useQuery<StopPreviewQuery, StopPreviewQueryVariables>({ query: StopPreviewDocument, ...options });
};
export const StopRouteDocument = gql`
    query StopRoute($stop: ID!, $route: ID!, $today: Date!, $tomorrow: Date!) {
  stopRoute(stop: $stop, route: $route) {
    stop {
      id
      name
      code
    }
    route {
      id
      name
      text
      background
    }
    headsign
    ...LiveData
    schedule {
      today: on(date: $today) {
        stoptime {
          id
          time
        }
      }
      tomorrow: on(date: $tomorrow) {
        stoptime {
          id
          time
        }
      }
      next(limit: 3) {
        stoptime {
          id
          time
          overflow
          trip {
            service {
              saturday
              sunday
              monday
              start
              end
            }
          }
        }
      }
    }
  }
}
    ${LiveDataFragmentDoc}`;

export function useStopRouteQuery(options: Omit<Urql.UseQueryArgs<StopRouteQueryVariables>, 'query'>) {
  return Urql.useQuery<StopRouteQuery, StopRouteQueryVariables>({ query: StopRouteDocument, ...options });
};
export const TextSearchDocument = gql`
    query TextSearch($text: String!, $page: PageInput!) {
  searchStopText(text: $text, page: $page) {
    results {
      ...StopPreview
    }
  }
}
    ${StopPreviewFragmentDoc}`;

export function useTextSearchQuery(options: Omit<Urql.UseQueryArgs<TextSearchQueryVariables>, 'query'>) {
  return Urql.useQuery<TextSearchQuery, TextSearchQueryVariables>({ query: TextSearchDocument, ...options });
};
export const TravelPlannerDocument = gql`
    query TravelPlanner($origin: ID!, $destination: ID!, $options: TravelOptions!) {
  travelPlanner(origin: $origin, destination: $destination, options: $options) {
    ...TravelSchedule
  }
}
    ${TravelScheduleFragmentDoc}`;

export function useTravelPlannerQuery(options: Omit<Urql.UseQueryArgs<TravelPlannerQueryVariables>, 'query'>) {
  return Urql.useQuery<TravelPlannerQuery, TravelPlannerQueryVariables>({ query: TravelPlannerDocument, ...options });
};
export const TravelPlannerFixedRouteDocument = gql`
    query TravelPlannerFixedRoute($input: [TravelLegInput!]!, $options: TravelOptions!) {
  travelPlannerFixedRoute(input: $input, options: $options) {
    ...TravelSchedule
  }
}
    ${TravelScheduleFragmentDoc}`;

export function useTravelPlannerFixedRouteQuery(options: Omit<Urql.UseQueryArgs<TravelPlannerFixedRouteQueryVariables>, 'query'>) {
  return Urql.useQuery<TravelPlannerFixedRouteQuery, TravelPlannerFixedRouteQueryVariables>({ query: TravelPlannerFixedRouteDocument, ...options });
};
export const TravelPlannerDeparturesDocument = gql`
    query TravelPlannerDepartures($origin: ID!, $destination: ID!, $route: ID!, $after: DateTime!, $limit: Int!) {
  stopRoute(stop: $origin, route: $route) {
    schedule: scheduleReaches(destination: $destination) {
      next(limit: $limit, after: $after) {
        stoptime {
          id
          time
        }
      }
    }
  }
}
    `;

export function useTravelPlannerDeparturesQuery(options: Omit<Urql.UseQueryArgs<TravelPlannerDeparturesQueryVariables>, 'query'>) {
  return Urql.useQuery<TravelPlannerDeparturesQuery, TravelPlannerDeparturesQueryVariables>({ query: TravelPlannerDeparturesDocument, ...options });
};
export const TravelRouteDocument = gql`
    query TravelRoute($input: [TravelLegInput!]!) {
  travelRoute(input: $input) {
    route {
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
      stopRoute {
        route {
          id
          name
          text
          background
        }
        headsign
      }
    }
  }
}
    `;

export function useTravelRouteQuery(options: Omit<Urql.UseQueryArgs<TravelRouteQueryVariables>, 'query'>) {
  return Urql.useQuery<TravelRouteQuery, TravelRouteQueryVariables>({ query: TravelRouteDocument, ...options });
};