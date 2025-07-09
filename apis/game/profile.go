package game

import "github.com/golangmc/minecraft-server/apis/uuid"

type Profile struct {
	UUID     uuid.UUID
	EntityID int32
	Name     string

	MaxChunksCount int
	SendedChunks   map[string]bool
	Properties     []*ProfileProperty
}

type ProfileProperty struct {
	Name      string
	Value     string
	Signature *string
}
