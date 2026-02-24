package lemin

import (
	"fmt"
	"lem-inV2/internal/types"
	"sort"
	"strconv"
	"strings"
)

func DistributeAnts(ants int, paths []types.Path) [][]int {
	dist := make([][]int, len(paths))
	for i := 1; i <= ants; i++ {
		bestPath := 0
		minVal := len(paths[0]) + len(dist[0])
		for j := 1; j < len(paths); j++ {
			val := len(paths[j]) + len(dist[j])
			if val < minVal {
				minVal = val
				bestPath = j
			}
		}
		dist[bestPath] = append(dist[bestPath], i)
	}
	return dist
}

func MoveAnts(totalAnts int, paths []types.Path, distribution [][]int) {
	type ActiveAnt struct {
		id    int
		path  types.Path
		index int
	}
	var activeAnts []*ActiveAnt
	finished := 0
	for finished < totalAnts {
		var moves []string
		var stillMoving []*ActiveAnt
		for _, aa := range activeAnts {
			aa.index++
			moves = append(moves, fmt.Sprintf("L%d-%s", aa.id, aa.path[aa.index]))
			if aa.index == len(aa.path)-1 {
				finished++
			} else {
				stillMoving = append(stillMoving, aa)
			}
		}
		activeAnts = stillMoving
		for pIdx := range paths {
			if len(distribution[pIdx]) > 0 {
				antID := distribution[pIdx][0]
				distribution[pIdx] = distribution[pIdx][1:]
				newAA := &ActiveAnt{id: antID, path: paths[pIdx], index: 1}
				moves = append(moves, fmt.Sprintf("L%d-%s", newAA.id, newAA.path[1]))
				if newAA.index == len(newAA.path)-1 {
					finished++
				} else {
					activeAnts = append(activeAnts, newAA)
				}
			}
		}
		if len(moves) > 0 {
			sort.Slice(moves, func(i, j int) bool {
				idI, _ := strconv.Atoi(strings.Split(strings.Split(moves[i], "-")[0], "L")[1])
				idJ, _ := strconv.Atoi(strings.Split(strings.Split(moves[j], "-")[0], "L")[1])
				return idI < idJ
			})
			fmt.Println(strings.Join(moves, " "))
		}
	}
}