package lemin

import (
	"lem-inV2/internal/types"
)

func FindOptimalPaths(rooms map[string]*types.Room, start, end string, totalAnts int) []types.Path {
	capacity := make(map[types.Edge]int)

	for name, r := range rooms {
		for _, link := range r.Links {
			capacity[types.Edge{From:name,To: link}] = 1
		}
	}

	var bestPathSet []types.Path

	for {
		parent := make(map[string]string)
		queue := []string{start}
		visited := map[string]bool{start: true}

		found := false
		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]
			if curr == end {
				found = true
				break
			}
			for _, next := range rooms[curr].Links {
				// Residual graph logic: check if capacity is available
				if !visited[next] && capacity[types.Edge{From:curr,To: next}] > 0 {
					visited[next] = true
					parent[next] = curr
					queue = append(queue, next)
				}
			}
		}

		if !found {
			break
		}

		// Update capacities (Augmenting Path)
		curr := end
		for curr != start {
			p := parent[curr]
			capacity[types.Edge{From:p,To: curr}] -= 1
			capacity[types.Edge{From:curr,To: p}] += 1
			curr = p
		}

		// Calculate efficiency of the CURRENT flow state
		currentPaths := ExtractPaths(rooms, start, end, capacity)

	
			bestPathSet = currentPaths
		
	}
	return bestPathSet
}

func ExtractPaths(rooms map[string]*types.Room, start, end string, cap map[types.Edge]int) []types.Path {
	var paths []types.Path
	// Create a working copy of capacities to mark types.Edges as used during extraction
	tempCap := make(map[types.Edge]int)
	for k, v := range cap {
		tempCap[k] = v
	}

	for {
		path := []string{start}
		curr := start
		found := false

		// Follow types.Edges where capacity was reduced (flow was pushed)
		for {
			nextStep := ""
			for _, next := range rooms[curr].Links {
				// In Edmonds-Karp, a used types.Edge has capacity 0 (from original 1)
				// and its reverse types.Edge in the residual graph has capacity 2 (from original 1)
				if tempCap[types.Edge{From:curr,To: next}] == 0 && tempCap[types.Edge{From:next,To: curr}] == 2 {
					nextStep = next
					break
				}
			}

			if nextStep == "" {
				break
			}

			path = append(path, nextStep)
			// "Clean" the flow so we don't re-extract the same path
			tempCap[types.Edge{From:curr,To: nextStep}] = 1
			tempCap[types.Edge{From:nextStep,To: curr}] = 1
			curr = nextStep

			if curr == end {
				found = true
				break
			}
		}

		if !found {
			break
		}
		paths = append(paths, path)
	}
	return paths
}

