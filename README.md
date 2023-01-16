# stop-checker.com

stop-checker.com provides a travel planner, bus schedules, and real-time bus GPS positions for Ottawa's public transit system.
The backend is a GraphQL API written in Go. The frontend was created using React and Typescript.

## Travel Planner

The travel planner is available here: https://www.stop-checker.com/travel

The travel planner is split into two parts. The first part is the "planner" which is responsible for determining the optimal travel plan between the origin and destination locations using either an "arrive by" or "depart at" mode. The output is a travel plan ([`model.TravelPlan`](https://github.com/danielholmes839/stop-checker.com-2/blob/18503348549fbd9791376ca73fa5b786b2e91d25/backend/db/model/travel.go#L8-L18)) which contains a list of route IDs and origin/destination IDs. The second part is the "scheduler" which is responsible for creating a schedule from a travel plan. 
This is done separately for a few reasons. The first reason was to create a feature to save and reuse specific travel plans. The second reason is just separation of concerns and allowing for flexibility later on. For example the current travel planner attempts to find the optimal travel plan however, in the future the travel planner may generate alternative travel plans. 
The interfaces for the travel planner and scheduler can be found [here](https://github.com/danielholmes839/stop-checker.com-2/blob/18503348549fbd9791376ca73fa5b786b2e91d25/backend/application/services/services.go#L10-L18).

### Planner Algorithm

The travel planner algorithm uses A* search with the following heuristics:

1. Walking distance. The algorithm tracks the cumulative walking distance each path takes. This heuristic is meant to increase the quality of the solutions. Before adding this heuristic the algorithm would often find paths that exit buses too soon.

2. Bus transfers. The algorithm tracks the cumulative number of buses each path takes. The first bus does not count as a transfer. This heuristic is meant to increase the quality of the solutions. Before adding this heuristic the algorithm would suggest plans with too many or unecessary transfers.

3. The distance of each node to the destination. Harversine distance is used to calculate the distance between two coordinates. This heuristic is meant to improve the performance of the algorithm by reducing the number of nodes the travel planner needs to explore before reaching a solution.

The algorithm explores nodes by transit and walking. Nodes can be explored by walking as long as the previous node was reached by transit. 
Each node also contains a set of route IDs used as "blockers". This set ensures the algorithm does take a route "twice in a row".

The [Open Source Routing Machine](https://project-osrm.org/) project was used to get accurate walking directions between bus stops and origin/destination locations. 
Walking directions between neighboring bus stops are cached to increase the performance of the algorithm. Walking directions are requested for all stops within 1km of the origin.
If a node being explored is within 1km of the destination, then walking directions are requested and a destination node is added to the priority queue

The algorithm stops when a destination node is removed from the priority queue.
