package types

type Room struct {
	Name  string
	Links []string
}
type Edge struct {
	From, To string
}

type Path []string