package main

// Room keeps information about particular room
type Room struct {
	ID   int
	Name string
}

var rooms = []Room{
	Room{1, "room1"},
	Room{2, "room2"},
	Room{3, "room3"},
	Room{4, "room4"},
}
