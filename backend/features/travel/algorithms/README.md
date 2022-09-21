# Travel Planner Algorithms

The stop-checker.com travel planner works using graph algorithms. 
Currently using Dijkstra's shortest path with a small heuristic to account for transfering buses.

## Expanding Nodes

- Each bus stop is a node
- Expand a node first with transit options then by walking
- Expand by transit should include the fastest route to each stop that's reachable
- Expand by walking only to stops that have a new route to use