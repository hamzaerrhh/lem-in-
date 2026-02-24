package main

import (
	"fmt"
	lemin "lem-inV2/internal/lem-in"
	"lem-inV2/internal/parsing"
	"os"
)



func main() {

	if len(os.Args) != 2 {
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("ERROR: invalid data format")
		return
	}
	defer file.Close()

	ants, rooms, startNode, endNode, err := parsing.ParseInput(file)
	if err != nil {
		fmt.Println("ERROR: invalid data format")
		return
	}

	paths := lemin.FindOptimalPaths(rooms, startNode, endNode, ants)
	if len(paths) == 0 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// 2. Distribute and move
	distribution := lemin.DistributeAnts(ants, paths)
	lemin.MoveAnts(ants, paths, distribution)
}



