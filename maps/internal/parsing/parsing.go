package parsing

import (
	"bufio"
	"fmt"
	"lem-inV2/internal/types"
	"os"
	"strconv"

	"strings"
)

func ParseInput(file *os.File) (int, map[string]*types.Room, string, string, error) {
	scanner := bufio.NewScanner(file)
	rooms := make(map[string]*types.Room)
	var ants int
	var startNode, endNode string
	var mode string

	lineCount := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || (strings.HasPrefix(line, "#") && line != "##start" && line != "##end") {
			continue
		}
		if lineCount == 0 {
			var err error
			ants, err = strconv.Atoi(line)
			if err != nil || ants <= 0 {
				return 0, nil, "", "", fmt.Errorf("invalid ants")
			}
			lineCount++
			continue
		}
		if line == "##start" {
			mode = "start"
			continue
		} else if line == "##end" {
			mode = "end"
			continue
		}
		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				continue
			}
			if r1, ok := rooms[parts[0]]; ok {
				r1.Links = append(r1.Links, parts[1])
			}
			if r2, ok := rooms[parts[1]]; ok {
				r2.Links = append(r2.Links, parts[0])
			}
		} else {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				name := parts[0]
				rooms[name] = &types.Room{Name: name}
				if mode == "start" {
					startNode = name
					mode = ""
				}
				if mode == "end" {
					endNode = name
					mode = ""
				}
			}
		}
	}
	return ants, rooms, startNode, endNode, nil
}