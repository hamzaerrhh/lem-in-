package main

import (
	"fmt"
	"os"

	lemin "lemin/internal/lem-in"
	parsing "lemin/internal/parsing"
	types "lemin/internal/types"
)

func main() {
	args := os.Args
	// handle the args

	fileName := args[1]
	coulounie, err := parsing.ParseFile(fileName)
	if err != nil {
		fmt.Println("err in parsing file", err)
		os.Exit(0)
	}

	errPath := lemin.FindAllPaths(coulounie)
	if errPath != nil {
		fmt.Println("no avaliable path")
		return
	}

	fmt.Println("number of ant", types.Ant_number)
	lemin.TravelAnt()
}
