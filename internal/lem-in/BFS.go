package lemin

import (
	"fmt"

	"lemin/internal/types"
)

type Edge struct {
	to       int
	rev      int
	capacity int
}

type FlowGraph struct {
	adj [][]*Edge
}


func FindAllPaths(graph *types.GraphType) error {
	if graph.Start == nil || graph.End == nil {
		return fmt.Errorf("graph must have Start and End rooms")
	}

	// Create consistent room indexing
	roomIndex := make(map[*types.Room]int)
	indexRoom := make(map[int]*types.Room)
	idx := 0
	for _, room := range graph.Rooms {
		roomIndex[room] = idx
		indexRoom[idx] = room
		idx++
	}

	flowGraph := buildFlowGraph(graph, roomIndex)
	
	s := roomIndex[graph.Start]
	t := roomIndex[graph.End]
	source := s*2 + 1  // output node of start
	sink := t * 2      // input node of end
	
	maxFlow(flowGraph, source, sink)

	paths := extractPaths(flowGraph, graph, graph.Start, graph.End, roomIndex, indexRoom)

	if len(paths) == 0 {
		return fmt.Errorf("no paths found")
	}

	var tunnels []types.Tunnel
	for i, p := range paths {
		tunnels = append(tunnels, types.Tunnel{
			Name:    fmt.Sprintf("Path %d", i+1),
			Weight:  len(p),
			Roadmap: p,
		})
	}

	// Sort paths by length (weight)
	for i := 0; i < len(tunnels)-1; i++ {
		for j := i + 1; j < len(tunnels); j++ {
			if tunnels[i].Weight > tunnels[j].Weight {
				tunnels[i], tunnels[j] = tunnels[j], tunnels[i]
			}
		}
	}

	// Select optimal paths based on number of ants
	optimalPaths := SelectOptimalPaths(tunnels, types.Ant_number)
	if len(optimalPaths) == 0 {
		return fmt.Errorf("no valid paths found for %d ants", types.Ant_number)
	}

	types.Tunnels = optimalPaths
	return nil
}



// bfs finds a single shortest path avoiding used rooms
func bfs(graph *types.GraphType, usedRooms map[*types.Room]bool) ([]*types.Room, bool) {
	start := graph.Start
	end := graph.End

	queue := []*types.Room{start}
	visited := make(map[*types.Room]bool)
	parent := make(map[*types.Room]*types.Room)

	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == end {
			// reconstruct path
			path := []*types.Room{}
			for r := end; r != nil; r = parent[r] {
				path = append([]*types.Room{r}, path...)
			}
			return path, true
		}

		for _, neighbor := range current.Neighborhood {
			if visited[neighbor] {
				continue
			}
			if usedRooms[neighbor] && neighbor != end && neighbor != start {
				continue // skip rooms already used in another path
			}
			visited[neighbor] = true
			parent[neighbor] = current
			queue = append(queue, neighbor)
		}
	}

	return nil, false
}

func SelectOptimalPaths(paths []types.Tunnel, ants int) []types.Tunnel {
	n := len(paths)

	bestTurns := -1
	bestK := 0

	for k := 1; k <= n; k++ {
		totalLen := 0
		for i := 0; i < k; i++ {
			totalLen += paths[i].Weight
		}

		// Compute minimal turns for k paths
		turns := (ants + totalLen - k) / k

		valid := true
		for i := 0; i < k; i++ {
			if turns < paths[i].Weight-1 {
				valid = false
				break
			}
		}

		if valid {
			if bestTurns == -1 || turns < bestTurns {
				bestTurns = turns
				bestK = k
			}
		}
	}

	if bestK == 0 {
		return nil
	}

	return paths[:bestK]
}

func buildFlowGraph(graph *types.GraphType, roomIndex map[*types.Room]int) *FlowGraph {

    n := len(graph.Rooms) * 2
    fg := &FlowGraph{adj: make([][]*Edge, n)}

    addEdge := func(u, v, cap int) {
        fg.adj[u] = append(fg.adj[u], &Edge{v, len(fg.adj[v]), cap})
        fg.adj[v] = append(fg.adj[v], &Edge{u, len(fg.adj[u]) - 1, 0})
    }

    // Create in/out nodes for each room
    for room, i := range roomIndex {
        in := i * 2
        out := i*2 + 1

        // Set capacity based on room type
        if room == graph.Start || room == graph.End {
            addEdge(in, out, types.Ant_number) // infinite capacity for start/end
        } else {
            addEdge(in, out, 1) // capacity 1 for intermediate rooms
        }
    }

    // Connect rooms based on neighborhood
    for room, i := range roomIndex {
        out := i*2 + 1  // output node of current room

        for _, neigh := range room.Neighborhood {
            j := roomIndex[neigh]
            inNeigh := j * 2  // input node of neighbor
            addEdge(out, inNeigh, 1) // capacity 1 for each edge
        }
    }

    return fg
}

func maxFlow(fg *FlowGraph, source, sink int) int {
	flow := 0

	for {
		parent := make([][2]int, len(fg.adj))
		for i := range parent {
			parent[i] = [2]int{-1, -1}
		}

		queue := []int{source}
		parent[source] = [2]int{source, -1}

		for len(queue) > 0 && parent[sink][0] == -1 {
			u := queue[0]
			queue = queue[1:]

			for i, e := range fg.adj[u] {
				if parent[e.to][0] == -1 && e.capacity > 0 {
					parent[e.to] = [2]int{u, i}
					queue = append(queue, e.to)
				}
			}
		}

		if parent[sink][0] == -1 {
			break
		}

		v := sink
		for v != source {
			u := parent[v][0]
			i := parent[v][1]
			fg.adj[u][i].capacity--
			rev := fg.adj[u][i].rev
			fg.adj[v][rev].capacity++
			v = u
		}

		flow++
	}

	return flow
}

func extractPaths(fg *FlowGraph, graph *types.GraphType, start, end *types.Room, roomIndex map[*types.Room]int, indexRoom map[int]*types.Room) [][]*types.Room {

    paths := [][]*types.Room{}
    
    s := roomIndex[start]
    t := roomIndex[end]
    
    // Start from the output node of start, go to input node of end
    source := s*2 + 1  // output node of start
    sink := t * 2      // input node of end
    
    // Track used edges to find multiple paths
    usedEdges := make(map[[2]int]bool)
    
    for {
        path := []*types.Room{start}
        current := source
        
        visitedNodes := make(map[int]bool)
        visitedNodes[source] = true
        
        // Follow the flow
        for current != sink {
            found := false
            
            // Look for an edge with flow (capacity 0 in residual means flow=1 in original)
            for _, e := range fg.adj[current] {
                edgeKey := [2]int{current, e.to}
                if usedEdges[edgeKey] || visitedNodes[e.to] {
                    continue
                }
                
                // In residual graph, capacity 0 means flow=1 in original
                if e.capacity == 0 {
                    usedEdges[edgeKey] = true
                    visitedNodes[e.to] = true
                    
                    // e.to is an in-node (even index)
                    nextRoomID := e.to / 2
                    if nextRoom, ok := indexRoom[nextRoomID]; ok {
                        // Skip if it's the start room (already in path)
                        if nextRoom == start {
                            continue
                        }
                        
                        // Add the room to path only if not already there
                        alreadyInPath := false
                        for _, r := range path {
                            if r == nextRoom {
                                alreadyInPath = true
                                break
                            }
                        }
                        if alreadyInPath {
                            continue
                        }
                        
                        // Add the room to path
                        path = append(path, nextRoom)
                        
                        // Now we need to go from in-node to out-node
                        // The in->out edge should also have capacity 0 if used
                        inNode := e.to
                        outNode := inNode + 1
                        
                        // Check if in->out edge has flow
                        hasFlow := false
                        for _, e2 := range fg.adj[inNode] {
                            if e2.to == outNode && e2.capacity == 0 {
                                hasFlow = true
                                edgeKey2 := [2]int{inNode, outNode}
                                usedEdges[edgeKey2] = true
                                visitedNodes[outNode] = true
                                break
                            }
                        }
                        
                        if hasFlow || nextRoom == end {
                            // For end room, we don't need to check out-node
                            if nextRoom == end {
                                current = sink
                            } else {
                                current = outNode
                            }
                            found = true
                            break
                        }
                    }
                }
            }
            
            if !found {
                break
            }
        }
        
        // Check if we found a valid path to end
        if len(path) < 2 || path[len(path)-1] != end {
            break
        }
        
        paths = append(paths, path)
    }
    
    return paths
}


