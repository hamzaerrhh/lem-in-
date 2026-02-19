package types

var (
	Ant_number int
	Graph      *GraphType
	Tunnels    []Tunnel
)

type Tunnel struct {
	Weight  int
	Name    string
	Roadmap []*Room
}

type GraphType struct {
	Rooms map[string]*Room
	Start *Room
	End   *Room
}

type Room struct {
	Name         string
	X, Y         int
	Neighborhood []*Room
}
