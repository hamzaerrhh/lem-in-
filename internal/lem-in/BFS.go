package lemin

import (
	"fmt"

	"lemin/internal/types"
)

func FindAllPaths(graph *types.GraphType) error {
	if graph.Start == nil || graph.End == nil {
		return fmt.Errorf("graph must have Start and End rooms")
	}

	var tunnels []types.Tunnel
	usedRooms := make(map[*types.Room]bool) // intermediate rooms already in a path

	for {
		path, ok := bfs(graph, usedRooms)
		if !ok {
			break // no more paths
		}

		tunnel := types.Tunnel{
			Weight:  len(path),
			Name:    fmt.Sprintf("Path %d", len(tunnels)+1),
			Roadmap: path,
		}
		tunnels = append(tunnels, tunnel)

		// Mark intermediate rooms as used (excluding start and end)
		for i := 1; i < len(path)-1; i++ {
			usedRooms[path[i]] = true
		}
	}

	if len(tunnels) == 0 {
		return fmt.Errorf("no paths found from Start to End")
	}
	types.Tunnels = tunnels

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
