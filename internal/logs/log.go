package logLemin

import (
	"fmt"
	"lemin/internal/types"
)

func LogTunnels() {

	fmt.Println("tunnels numbers",len(types.Tunnels))
	for _, tunnel := range types.Tunnels {
    
		fmt.Printf("Tunnel: %s, Length: %d\n", tunnel.Name, len(tunnel.Roadmap))
		for index, room := range tunnel.Roadmap {

			fmt.Print( room.Name)
			if index < len(tunnel.Roadmap)-1 {
				fmt.Print(" -> ")
			}
		}
		fmt.Println()
	}
}