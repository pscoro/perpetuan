package main

import "github.com/google/uuid"

type session struct {
	ID     uuid.UUID
	Player player
	Room string
}
