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
  Datetime: any;
  Time: any;
};

export type Bus = {
  __typename?: 'Bus';
  arrival: Scalars['Datetime'];
  distance?: Maybe<Scalars['Float']>;
  headsign: Scalars['String'];
  lastUpdated: Scalars['Datetime'];
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

export type Path = {
  __typename?: 'Path';
  distance: Scalars['Float'];
  path: Array<Location>;
};

export type Query = {
  __typename?: 'Query';
  searchStopLocation: StopSearchPayload;
  searchStopText: StopSearchPayload;
  stop?: Maybe<Stop>;
  stopRoute?: Maybe<StopRoute>;
  travelPlanner: TravelSchedulePayload;
  travelPlannerFixedRoute: TravelSchedulePayload;
  travelPlannerFixedRoutes: Array<TravelSchedulePayload>;
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
  routeId: Scalars['ID'];
  stopId: Scalars['ID'];
};


export type QueryTravelPlannerArgs = {
  destination: LocationInput;
  options: TravelPlannerOptions;
  origin: LocationInput;
};


export type QueryTravelPlannerFixedRouteArgs = {
  input: TravelPlanInput;
  options: TravelPlannerOptions;
};


export type QueryTravelPlannerFixedRoutesArgs = {
  input: Array<TravelPlanInput>;
  options: TravelPlannerOptions;
};

export type Route = {
  __typename?: 'Route';
  background: Scalars['Color'];
  id: Scalars['ID'];
  name: Scalars['String'];
  text: Scalars['Color'];
};

export type Schedule = {
  __typename?: 'Schedule';
  next: Array<ScheduleResult>;
  on: Array<ScheduleResult>;
};


export type ScheduleNextArgs = {
  after?: InputMaybe<Scalars['Datetime']>;
  limit: Scalars['Int'];
};


export type ScheduleOnArgs = {
  date: Scalars['Date'];
};

export enum ScheduleMode {
  ArriveBy = 'ARRIVE_BY',
  DepartAt = 'DEPART_AT'
}

export type ScheduleResult = {
  __typename?: 'ScheduleResult';
  datetime: Scalars['Datetime'];
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
  schedule: Schedule;
  scheduleReaches?: Maybe<Schedule>;
  stop: Stop;
};


export type StopRouteReachesArgs = {
  forward: Scalars['Boolean'];
};


export type StopRouteScheduleReachesArgs = {
  destination: Scalars['ID'];
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
  departure: Scalars['Datetime'];
  duration: Scalars['Int'];
  route: Route;
  trip: Trip;
  wait: Scalars['Int'];
};

export type TravelPlanInput = {
  destination: LocationInput;
  legs: Array<InputMaybe<TravelPlanLegInput>>;
  origin: LocationInput;
};

export type TravelPlanLegInput = {
  destinationId: Scalars['ID'];
  originId: Scalars['ID'];
  routeId: Scalars['ID'];
};

export type TravelPlannerOptions = {
  datetime?: InputMaybe<Scalars['Datetime']>;
  mode: ScheduleMode;
};

export type TravelSchedule = {
  __typename?: 'TravelSchedule';
  destination: TravelScheduleNode;
  duration: Scalars['Int'];
  legs: Array<TravelScheduleLeg>;
  origin: TravelScheduleNode;
};

export type TravelScheduleLeg = {
  __typename?: 'TravelScheduleLeg';
  destination: TravelScheduleNode;
  duration: Scalars['Int'];
  origin: TravelScheduleNode;
  transit?: Maybe<Transit>;
  walk?: Maybe<Path>;
};

export type TravelScheduleNode = {
  __typename?: 'TravelScheduleNode';
  arrival: Scalars['Datetime'];
  location: Location;
  stop?: Maybe<Stop>;
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


export type StopPageQuery = { __typename?: 'Query', stop?: { __typename?: 'Stop', id: string, name: string, code: string, routes: Array<{ __typename?: 'StopRoute', headsign: string, route: { __typename?: 'Route', id: string, name: string, text: any, background: any }, schedule: { __typename?: 'Schedule', next: Array<{ __typename?: 'ScheduleResult', stoptime: { __typename?: 'StopTime', id: string, time: any } }> }, liveBuses: Array<{ __typename?: 'Bus', headsign: string }> }> } | null };

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


export type StopRouteQuery = { __typename?: 'Query', stopRoute?: { __typename?: 'StopRoute', headsign: string, liveMap?: string | null, stop: { __typename?: 'Stop', id: string, name: string, code: string }, route: { __typename?: 'Route', id: string, name: string, text: any, background: any }, schedule: { __typename?: 'Schedule', today: Array<{ __typename?: 'ScheduleResult', stoptime: { __typename?: 'StopTime', id: string, time: any } }>, tomorrow: Array<{ __typename?: 'ScheduleResult', stoptime: { __typename?: 'StopTime', id: string, time: any } }>, next: Array<{ __typename?: 'ScheduleResult', stoptime: { __typename?: 'StopTime', id: string, time: any, overflow: boolean, trip: { __typename?: 'Trip', service: { __typename?: 'Service', saturday: boolean, sunday: boolean, monday: boolean, start: any, end: any } } } }> }, liveBuses: Array<{ __typename?: 'Bus', headsign: string, arrival: any, lastUpdated: any, lastUpdatedMessage: string, lastUpdatedMinutes: number, distance?: number | null, location?: { __typename?: 'Location', latitude: number, longitude: number } | null }> } | null };

export type LiveDataFragment = { __typename?: 'StopRoute', liveMap?: string | null, liveBuses: Array<{ __typename?: 'Bus', headsign: string, arrival: any, lastUpdated: any, lastUpdatedMessage: string, lastUpdatedMinutes: number, distance?: number | null, location?: { __typename?: 'Location', latitude: number, longitude: number } | null }> };

export type StopRouteDetailsQueryVariables = Exact<{
  originId: Scalars['ID'];
  destinationId: Scalars['ID'];
  routeId: Scalars['ID'];
  after?: InputMaybe<Scalars['Datetime']>;
}>;


export type StopRouteDetailsQuery = { __typename?: 'Query', stopRoute?: { __typename?: 'StopRoute', scheduleReaches?: { __typename?: 'Schedule', next: Array<{ __typename?: 'ScheduleResult', datetime: any }> } | null, liveBuses: Array<{ __typename?: 'Bus', arrival: any, distance?: number | null, lastUpdatedMinutes: number, lastUpdated: any }> } | null };

export type TextSearchQueryVariables = Exact<{
  text: Scalars['String'];
  page: PageInput;
}>;


export type TextSearchQuery = { __typename?: 'Query', searchStopText: { __typename?: 'StopSearchPayload', results: Array<{ __typename?: 'Stop', id: string, name: string, code: string, routes: Array<{ __typename?: 'StopRoute', headsign: string, route: { __typename?: 'Route', id: string, name: string, text: any, background: any } }> }> } };

export type ScheduleNodeFragment = { __typename?: 'TravelScheduleNode', arrival: any, location: { __typename?: 'Location', latitude: number, longitude: number }, stop?: { __typename?: 'Stop', id: string, name: string, code: string } | null };

export type ScheduleWalkFragment = { __typename?: 'TravelScheduleLeg', duration: number, walk?: { __typename?: 'Path', distance: number, path: Array<{ __typename?: 'Location', latitude: number, longitude: number }> } | null };

export type ScheduleTransitFragment = { __typename?: 'Transit', departure: any, duration: number, wait: number, route: { __typename?: 'Route', id: string, name: string, text: any, background: any }, trip: { __typename?: 'Trip', id: string, headsign: string, shape: Array<{ __typename?: 'Location', latitude: number, longitude: number }>, stoptimes: Array<{ __typename?: 'StopTime', id: string, time: any, stop: { __typename?: 'Stop', id: string, name: string } }> } };

export type ScheduleFragment = { __typename?: 'TravelSchedule', duration: number, origin: { __typename?: 'TravelScheduleNode', arrival: any, location: { __typename?: 'Location', latitude: number, longitude: number }, stop?: { __typename?: 'Stop', id: string, name: string, code: string } | null }, destination: { __typename?: 'TravelScheduleNode', arrival: any, location: { __typename?: 'Location', latitude: number, longitude: number }, stop?: { __typename?: 'Stop', id: string, name: string, code: string } | null }, legs: Array<{ __typename?: 'TravelScheduleLeg', duration: number, origin: { __typename?: 'TravelScheduleNode', arrival: any, location: { __typename?: 'Location', latitude: number, longitude: number }, stop?: { __typename?: 'Stop', id: string, name: string, code: string } | null }, destination: { __typename?: 'TravelScheduleNode', arrival: any, location: { __typename?: 'Location', latitude: number, longitude: number }, stop?: { __typename?: 'Stop', id: string, name: string, code: string } | null }, transit?: { __typename?: 'Transit', departure: any, duration: number, wait: number, route: { __typename?: 'Route', id: string, name: string, text: any, background: any }, trip: { __typename?: 'Trip', id: string, headsign: string, shape: Array<{ __typename?: 'Location', latitude: number, longitude: number }>, stoptimes: Array<{ __typename?: 'StopTime', id: string, time: any, stop: { __typename?: 'Stop', id: string, name: string } }> } } | null, walk?: { __typename?: 'Path', distance: number, path: Array<{ __typename?: 'Location', latitude: number, longitude: number }> } | null }> };

export type SchedulePayloadFragment = { __typename?: 'TravelSchedulePayload', error?: string | null, schedule?: { __typename?: 'TravelSchedule', duration: number, origin: { __typename?: 'TravelScheduleNode', arrival: any, location: { __typename?: 'Location', latitude: number, longitude: number }, stop?: { __typename?: 'Stop', id: string, name: string, code: string } | null }, destination: { __typename?: 'TravelScheduleNode', arrival: any, location: { __typename?: 'Location', latitude: number, longitude: number }, stop?: { __typename?: 'Stop', id: string, name: string, code: string } | null }, legs: Array<{ __typename?: 'TravelScheduleLeg', duration: number, origin: { __typename?: 'TravelScheduleNode', arrival: any, location: { __typename?: 'Location', latitude: number, longitude: number }, stop?: { __typename?: 'Stop', id: string, name: string, code: string } | null }, destination: { __typename?: 'TravelScheduleNode', arrival: any, location: { __typename?: 'Location', latitude: number, longitude: number }, stop?: { __typename?: 'Stop', id: string, name: string, code: string } | null }, transit?: { __typename?: 'Transit', departure: any, duration: number, wait: number, route: { __typename?: 'Route', id: string, name: string, text: any, background: any }, trip: { __typename?: 'Trip', id: string, headsign: string, shape: Array<{ __typename?: 'Location', latitude: number, longitude: number }>, stoptimes: Array<{ __typename?: 'StopTime', id: string, time: any, stop: { __typename?: 'Stop', id: string, name: string } }> } } | null, walk?: { __typename?: 'Path', distance: number, path: Array<{ __typename?: 'Location', latitude: number, longitude: number }> } | null }> } | null };

export type TravelPlannerQueryVariables = Exact<{
  origin: LocationInput;
  destination: LocationInput;
  options: TravelPlannerOptions;
}>;


export type TravelPlannerQuery = { __typename?: 'Query', travelPlanner: { __typename?: 'TravelSchedulePayload', error?: string | null, schedule?: { __typename?: 'TravelSchedule', duration: number, origin: { __typename?: 'TravelScheduleNode', arrival: any, location: { __typename?: 'Location', latitude: number, longitude: number }, stop?: { __typename?: 'Stop', id: string, name: string, code: string } | null }, destination: { __typename?: 'TravelScheduleNode', arrival: any, location: { __typename?: 'Location', latitude: number, longitude: number }, stop?: { __typename?: 'Stop', id: string, name: string, code: string } | null }, legs: Array<{ __typename?: 'TravelScheduleLeg', duration: number, origin: { __typename?: 'TravelScheduleNode', arrival: any, location: { __typename?: 'Location', latitude: number, longitude: number }, stop?: { __typename?: 'Stop', id: string, name: string, code: string } | null }, destination: { __typename?: 'TravelScheduleNode', arrival: any, location: { __typename?: 'Location', latitude: number, longitude: number }, stop?: { __typename?: 'Stop', id: string, name: string, code: string } | null }, transit?: { __typename?: 'Transit', departure: any, duration: number, wait: number, route: { __typename?: 'Route', id: string, name: string, text: any, background: any }, trip: { __typename?: 'Trip', id: string, headsign: string, shape: Array<{ __typename?: 'Location', latitude: number, longitude: number }>, stoptimes: Array<{ __typename?: 'StopTime', id: string, time: any, stop: { __typename?: 'Stop', id: string, name: string } }> } } | null, walk?: { __typename?: 'Path', distance: number, path: Array<{ __typename?: 'Location', latitude: number, longitude: number }> } | null }> } | null } };

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
export const ScheduleNodeFragmentDoc = gql`
    fragment ScheduleNode on TravelScheduleNode {
  arrival
  location {
    latitude
    longitude
  }
  stop {
    id
    name
    code
  }
}
    `;
export const ScheduleTransitFragmentDoc = gql`
    fragment ScheduleTransit on Transit {
  route {
    id
    name
    text
    background
  }
  trip {
    id
    headsign
    shape {
      latitude
      longitude
    }
    stoptimes {
      id
      stop {
        id
        name
      }
      time
    }
  }
  departure
  duration
  wait
}
    `;
export const ScheduleWalkFragmentDoc = gql`
    fragment ScheduleWalk on TravelScheduleLeg {
  duration
  walk {
    distance
    path {
      latitude
      longitude
    }
  }
}
    `;
export const ScheduleFragmentDoc = gql`
    fragment Schedule on TravelSchedule {
  duration
  origin {
    ...ScheduleNode
  }
  destination {
    ...ScheduleNode
  }
  legs {
    origin {
      ...ScheduleNode
    }
    destination {
      ...ScheduleNode
    }
    transit {
      ...ScheduleTransit
    }
    ...ScheduleWalk
  }
}
    ${ScheduleNodeFragmentDoc}
${ScheduleTransitFragmentDoc}
${ScheduleWalkFragmentDoc}`;
export const SchedulePayloadFragmentDoc = gql`
    fragment SchedulePayload on TravelSchedulePayload {
  schedule {
    ...Schedule
  }
  error
}
    ${ScheduleFragmentDoc}`;
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
  stopRoute(stopId: $stop, routeId: $route) {
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
export const StopRouteDetailsDocument = gql`
    query StopRouteDetails($originId: ID!, $destinationId: ID!, $routeId: ID!, $after: Datetime) {
  stopRoute(stopId: $originId, routeId: $routeId) {
    scheduleReaches(destination: $destinationId) {
      next(after: $after, limit: 3) {
        datetime
      }
    }
    liveBuses {
      arrival
      distance
      lastUpdatedMinutes
      lastUpdated
    }
  }
}
    `;

export function useStopRouteDetailsQuery(options: Omit<Urql.UseQueryArgs<StopRouteDetailsQueryVariables>, 'query'>) {
  return Urql.useQuery<StopRouteDetailsQuery, StopRouteDetailsQueryVariables>({ query: StopRouteDetailsDocument, ...options });
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
    query TravelPlanner($origin: LocationInput!, $destination: LocationInput!, $options: TravelPlannerOptions!) {
  travelPlanner(origin: $origin, destination: $destination, options: $options) {
    ...SchedulePayload
  }
}
    ${SchedulePayloadFragmentDoc}`;

export function useTravelPlannerQuery(options: Omit<Urql.UseQueryArgs<TravelPlannerQueryVariables>, 'query'>) {
  return Urql.useQuery<TravelPlannerQuery, TravelPlannerQueryVariables>({ query: TravelPlannerDocument, ...options });
};