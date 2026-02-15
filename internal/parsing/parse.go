package parsing

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	types "lemin/internal/types"
)

func ParseFile(filename string) (*types.GraphType, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format")
	}

	graph := &types.GraphType{
		Rooms: make(map[string]*types.Room),
	}

	fileData := string(data)
	lines := strings.Split(string(fileData), "\n")
	var startNext, endNext bool

	// Read number of ants
	if len(lines) == 0 {
		return nil, fmt.Errorf("ERROR: invalid data format")
	}

	firstLine := strings.TrimSpace(lines[0])
	if firstLine == "" || strings.HasPrefix(firstLine, "#") {
		return nil, fmt.Errorf("ERROR: invalid data format")
	}

	types.Ant_number, err = strconv.Atoi(firstLine)
	if err != nil || types.Ant_number <= 0 {
		return nil, fmt.Errorf("ERROR: invalid data format")
	}
	
	// Print the input file contents
	fmt.Println(types.Ant_number)

	// Parse the rest of the file
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##") {
			continue // skip comments
		}

		if line == "##start" {
			startNext = true
			continue
		} else if line == "##end" {
			endNext = true
			continue
		}

		// Room: name x y
		if parts := strings.Fields(line); len(parts) == 3 {
			x, errX := strconv.Atoi(parts[1])
			y, errY := strconv.Atoi(parts[2])
			if errX != nil || errY != nil {
				return nil, fmt.Errorf("ERROR: invalid data format")
			}
			room := &types.Room{
				Name:         parts[0],
				X:            x,
				Y:            y,
				Neighborhood: []*types.Room{},
			}
			graph.Rooms[room.Name] = room

			if startNext {
				fmt.Println("##start")
				graph.Start = room
				startNext = false
			} else if endNext {
				fmt.Println("##end")
				graph.End = room
				endNext = false
			}
			fmt.Println(line)
			continue
		}

		// Link: room1-room2
		if parts := strings.Split(line, "-"); len(parts) == 2 {
			room1Name := parts[0]
			room2Name := parts[1]

			room1, ok1 := graph.Rooms[room1Name]
			room2, ok2 := graph.Rooms[room2Name]

			if !ok1 || !ok2 {
				return nil, fmt.Errorf("ERROR: invalid data format")
			}

			// Ignore self-links
			if room1Name == room2Name {
				continue
			}

			// Add neighbors (no duplicates)
			if !contains(room1.Neighborhood, room2) {
				room1.Neighborhood = append(room1.Neighborhood, room2)
			}
			if !contains(room2.Neighborhood, room1) {
				room2.Neighborhood = append(room2.Neighborhood, room1)
			}
			fmt.Println(line)
			continue
		}

		return nil, fmt.Errorf("ERROR: invalid data format")
	}

	return graph, nil
}

// Helper to check if a room slice already contains a room
func contains(slice []*types.Room, r *types.Room) bool {
	for _, room := range slice {
		if room == r {
			return true
		}
	}
	return false
}
