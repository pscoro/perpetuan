package main

type setting struct {
	key          string
	portals      []portal
	sizeX, sizeY int
	imagePath    string
	grid         [][]gridLocation
}

type portal struct {
	in, out position
}
