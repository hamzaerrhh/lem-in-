package lemin

import (
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

		// 2. Spawn new ants (one per tunnel per turn)
		for _, tunnel := range types.Tunnels {

	if nextAntID > types.Ant_number {
		break
	}

	// Skip tunnels that are too short
	if len(tunnel.Roadmap) < 2 {
		continue
	}

	ant := &Ant{
		id:       nextAntID,
		tunnel:   &tunnel,
		position: 1, // move directly to first real room
	}

	ants = append(ants, ant)

	room := tunnel.Roadmap[1]
	moves = append(moves, "L"+strconv.Itoa(ant.id)+"-"+room.Name)

	if ant.position == len(tunnel.Roadmap)-1 {
		finished++
	}

	nextAntID++
}


		if len(moves) > 0 {
			println(strings.Join(moves, " "))
		}
	}
}
