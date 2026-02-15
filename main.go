package main

import (
	"fmt"
	"os"

	lemin "lemin/internal/lem-in"
	parsing "lemin/internal/parsing"
)

func main() {
	args := os.Args
	// handle the args

	if len(args) < 2 {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}

	fileName := args[1]
	coulounie, err := parsing.ParseFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	errPath := lemin.FindAllPaths(coulounie)
	if errPath != nil {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}

	fmt.Println()

	lemin.TravelAnt()
}
