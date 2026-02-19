package lemin

import (
	"fmt"
	"strconv"
	"strings"

	types "lemin/internal/types"
)

func TravelAnt() {
	type Ant struct {
		id       int
		tunnel   *types.Tunnel
		position int
	}

	type PathInfo struct {
		tunnel   *types.Tunnel
		length   int
		assigned int // total ants assigned to this path
		sent     int // ants already spawned
	}

	var ants []*Ant
	nextAntID := 1
	finished := 0
	var paths []PathInfo

	// ----------------------------
	// 1️⃣ Build paths
	// ----------------------------
	for i := range types.Tunnels {
		t := &types.Tunnels[i]
		if len(t.Roadmap) >= 2 {
			paths = append(paths, PathInfo{
				tunnel: t,
				length: len(t.Roadmap),
			})
		}
	}

	// ----------------------------
	// 2️⃣ Distribute ants optimally
	// ----------------------------
	remaining := types.Ant_number
	for remaining > 0 {
		best := 0
		for i := 1; i < len(paths); i++ {
			if paths[i].length+paths[i].assigned <
				paths[best].length+paths[best].assigned {
				best = i
			}
		}
		paths[best].assigned++
		remaining--
	}

	// ----------------------------
	// 3️⃣ Simulation
	// ----------------------------
	for finished < types.Ant_number {
		moves := []string{}

		// Move existing ants
		for _, ant := range ants {
			if ant.position >= len(ant.tunnel.Roadmap)-1 {
				continue
			}

			ant.position++
			room := ant.tunnel.Roadmap[ant.position]
			moves = append(moves, "L"+strconv.Itoa(ant.id)+"-"+room.Name)

			if ant.position == len(ant.tunnel.Roadmap)-1 {
				finished++
			}
		}

		// Spawn ants according to balanced distribution
		for i := range paths {
			if nextAntID > types.Ant_number {
				break
			}

			// Skip if this path already sent all its assigned ants
			if paths[i].sent >= paths[i].assigned {
				continue
			}

			tunnel := paths[i].tunnel

			// Check if first room is free
			canSpawn := true
			for _, ant := range ants {
				if ant.tunnel == tunnel && ant.position == 1 {
					canSpawn = false
					break
				}
			}

			if canSpawn {
				ant := &Ant{
					id:       nextAntID,
					tunnel:   tunnel,
					position: 1,
				}

				ants = append(ants, ant)
				paths[i].sent++
				nextAntID++

				room := tunnel.Roadmap[1]
				moves = append(moves, "L"+strconv.Itoa(ant.id)+"-"+room.Name)

				if ant.position == len(tunnel.Roadmap)-1 {
					finished++
				}
			}
		}

		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		} else if finished < types.Ant_number {
			break
		}
	}
}
