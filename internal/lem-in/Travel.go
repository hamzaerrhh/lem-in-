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

	var ants []*Ant
	nextAntID := 1
	finished := 0

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

		// Spawn new ants intelligently across paths
		// We spawn an ant in a path if it can accommodate one without blocking
		for pathIdx := range types.Tunnels {
			if nextAntID > types.Ant_number {
				break
			}

			tunnel := &types.Tunnels[pathIdx]
			
			// Skip tunnels that are too short
			if len(tunnel.Roadmap) < 2 {
				continue
			}

			// Check if we can spawn an ant in this tunnel
			// We can spawn if the path is not blocked (no ant at position 1)
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
					position: 0, // Start at position 0 (start room)
				}

				ants = append(ants, ant)
				nextAntID++

				// Move the newly spawned ant to first room
				if len(tunnel.Roadmap) > 1 {
					ant.position = 1
					room := tunnel.Roadmap[1]
					moves = append(moves, "L"+strconv.Itoa(ant.id)+"-"+room.Name)

					if ant.position == len(tunnel.Roadmap)-1 {
						finished++
					}
				}
			}
		}

		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		} else if finished < types.Ant_number {
			// If no moves but not all finished, something is wrong
			break
		}
	}
}
