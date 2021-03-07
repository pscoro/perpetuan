package main

import "github.com/google/uuid"

type player struct {
	ID       uuid.UUID
	Username string
	Pos      position
	Room     string
}

