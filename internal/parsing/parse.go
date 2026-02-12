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
	fmt.Println("file name is", filename)
	if err != nil {
		fmt.Println("can't open file")
		return nil, err
	}

	graph := &types.GraphType{
		Rooms: make(map[string]*types.Room),
	}

	fileData := string(data)
	lines := strings.Split(string(fileData), "\n")
	fmt.Println("data", lines)
	var startNext, endNext bool

	// Read number of ants
	if len(lines) == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	firstLine := strings.TrimSpace(lines[0])
	if firstLine == "" || strings.HasPrefix(firstLine, "#") {
		return nil, fmt.Errorf("first line must be number of ants")
	}

	fmt.Println("firstline is", lines)
	types.Ant_number, err = strconv.Atoi(firstLine)
	if err != nil {
		return nil, fmt.Errorf("invalid number of ants: %v", err)
	}

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
				return nil, fmt.Errorf("invalid coordinates for room %s", line)
			}
			room := &types.Room{
				Name:         parts[0],
				X:            x,
				Y:            y,
				Neighborhood: []*types.Room{},
			}
			graph.Rooms[room.Name] = room

			if startNext {
				graph.Start = room
				startNext = false
			} else if endNext {
				graph.End = room
				endNext = false
			}
			continue
		}

		// Link: room1-room2
		if parts := strings.Split(line, "-"); len(parts) == 2 {
			room1Name := parts[0]
			room2Name := parts[1]

			room1, ok1 := graph.Rooms[room1Name]
			room2, ok2 := graph.Rooms[room2Name]

			if !ok1 || !ok2 {
				return nil, fmt.Errorf("link references unknown room: %s-%s", room1Name, room2Name)
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
			continue
		}

		return nil, fmt.Errorf("invalid line: %s", line)
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
