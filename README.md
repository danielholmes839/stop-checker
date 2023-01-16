# stop-checker.com

stop-checker.com provides a travel planner, bus schedules, and real-time bus GPS positions for Ottawa's public transit system.
The backend is a GraphQL API written in Go. The frontend was created using React and Typescript.

## Travel Planner

**https://www.stop-checker.com/travel**

The travel planner is split into two parts. The first part is the "planner" which is responsible for determining the optimal travel plan between the origin and destination locations. The travel planner has an "arrive by" and "depart at" mode. The output is a travel plan which contains a list of route IDs and origin/destination IDs.

The second part is the "scheduler" which is responsible for creating a schedule from a travel plan.
This is done separately for a few reasons. The first reason was to create a feature to save and reuse specific travel plans. The second reason is just separation of concerns and allowing for flexibility later on. For example the current travel planner attempts to find one optimal travel plan however, in the future the travel planner may generate many alternative travel plans.
The interfaces for the travel planner and scheduler can be found [here](https://github.com/danielholmes839/stop-checker.com-2/blob/18503348549fbd9791376ca73fa5b786b2e91d25/backend/application/services/services.go#L10-L18). The travel package can be found [here](https://github.com/danielholmes839/stop-checker.com-2/tree/main/backend/features/travel).

### Planner Algorithm

The travel planner algorithm uses A\* search with the following heuristics:

1. Walking distance. The algorithm tracks the cumulative walking distance each path takes. This heuristic is meant to increase the quality of the solutions. Before adding this heuristic the algorithm would often find paths that exit buses too soon.

2. Bus transfers. The algorithm tracks the cumulative number of buses each path takes. The first bus does not count as a transfer. This heuristic is meant to increase the quality of the solutions. Before adding this heuristic the algorithm would suggest plans with too many or unecessary transfers.

3. The distance of each node to the destination. Harversine distance is used to calculate the distance between two coordinates. This heuristic is meant to improve the performance of the algorithm by reducing the number of nodes the travel planner needs to explore before reaching a solution.

The algorithm explores nodes by transit and walking. A node can be explored by walking as long as the node was reached by transit.
Each node also contains a set of route IDs used as "blockers". This set ensures the algorithm does take a route "twice in a row".

The [Open Source Routing Machine](https://project-osrm.org/) project was used to get accurate walking directions between bus stops and origin/destination locations.
Walking directions between neighboring bus stops are cached to increase the performance of the algorithm. Walking directions are requested for all stops within 1km of the origin.
If a node being explored is within 1km of the destination, then walking directions are requested and a destination node is added to the priority queue

The algorithm ends when a destination node is removed from the priority queue.

### Planner GraphQL API

The response from the GraphQL API might look like this:

```json
{
  "data": {
    "travelPlanner": {
      "schedule": {
        "duration": 28,
        "origin": {
          "arrival": "2023-01-16T23:26:00Z",
          "stop": null
        },
        "destination": {
          "arrival": "2023-01-16T23:54:00Z",
          "stop": null
        },
        "legs": [
          {
            "origin": {
              "arrival": "2023-01-16T23:26:00Z",
              "stop": null
            },
            "destination": {
              "arrival": "2023-01-16T23:31:00Z",
              "stop": {
                "id": "CG600",
                "name": "Campus / Commons",
                "code": "6612"
              }
            },
            "transit": null,
            "duration": 5,
            "walk": {
              "distance": 395.4215701870218
            }
          },
          {
            "origin": {
              "arrival": "2023-01-16T23:31:00Z",
              "stop": {
                "id": "CG600",
                "name": "Campus / Commons",
                "code": "6612"
              }
            },
            "destination": {
              "arrival": "2023-01-16T23:43:00Z",
              "stop": {
                "id": "AF920",
                "name": "Hurdman B",
                "code": "3023"
              }
            },
            "transit": {
              "route": {
                "id": "10-343",
                "name": "10",
                "text": "#FFFFFF",
                "background": "#D74100"
              },
              "trip": {
                "id": "88428008-JAN23-JANDA23-Weekday-31",
                "headsign": "Hurdman"
              },
              "departure": "2023-01-16T23:31:00Z",
              "duration": 12,
              "wait": 0
            },
            "duration": 12,
            "walk": null
          },
          {
            "origin": {
              "arrival": "2023-01-16T23:43:00Z",
              "stop": {
                "id": "AF920",
                "name": "Hurdman B",
                "code": "3023"
              }
            },
            "destination": {
              "arrival": "2023-01-16T23:46:00Z",
              "stop": {
                "id": "AF990",
                "name": "Hurdman O-Train West / Ouest",
                "code": "3023"
              }
            },
            "transit": null,
            "duration": 3,
            "walk": {
              "distance": 271.22632171994127
            }
          },
          {
            "origin": {
              "arrival": "2023-01-16T23:46:00Z",
              "stop": {
                "id": "AF990",
                "name": "Hurdman O-Train West / Ouest",
                "code": "3023"
              }
            },
            "destination": {
              "arrival": "2023-01-16T23:53:00Z",
              "stop": {
                "id": "CD998",
                "name": "Uottawa O-Train West / Ouest",
                "code": "3021"
              }
            },
            "transit": {
              "route": {
                "id": "1-343",
                "name": "1",
                "text": "#FFFFFF",
                "background": "#DA291C"
              },
              "trip": {
                "id": "87476259-JAN23-MThCFD-Weekday-01",
                "headsign": "Tunney's Pasture"
              },
              "departure": "2023-01-16T23:49:00Z",
              "duration": 4,
              "wait": 3
            },
            "duration": 7,
            "walk": null
          },
          {
            "origin": {
              "arrival": "2023-01-16T23:53:00Z",
              "stop": {
                "id": "CD998",
                "name": "Uottawa O-Train West / Ouest",
                "code": "3021"
              }
            },
            "destination": {
              "arrival": "2023-01-16T23:54:00Z",
              "stop": null
            },
            "transit": null,
            "duration": 1,
            "walk": {
              "distance": 51.519678130388236
            }
          }
        ]
      },
      "error": null
    }
  }
}
```
