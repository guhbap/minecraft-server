package game

import (
	"sync"

	"github.com/golangmc/minecraft-server/apis/uuid"
)

type ChunkPos struct {
	X int32
	Z int32
}

type Profile struct {
	UUID     uuid.UUID
	EntityID int32
	Name     string

	MaxChunksCount     int
	SendedChunks       map[ChunkPos]bool
	Spawned            bool
	SendedChunksL      sync.RWMutex
	ChunksCacheCenterX int32
	ChunksCacheCenterZ int32
	Properties         []*ProfileProperty
	OtherData          map[string]any

	PosInfo PosInfo
}

func (p *Profile) SetPosInfo(x, y, z float64, yaw, pitch float32) {
	p.PosInfo.X = x
	p.PosInfo.Y = y
	p.PosInfo.Z = z
	p.PosInfo.Yaw = yaw
	p.PosInfo.Pitch = pitch
}
func (p *Profile) UpdateYawPitch(yaw, pitch float32) {
	p.PosInfo.Yaw = yaw
	p.PosInfo.Pitch = pitch
}
func (p *Profile) UpdatePos(x, y, z float64) {
	p.PosInfo.X = x
	p.PosInfo.Y = y
	p.PosInfo.Z = z
}
func (p *Profile) GetPosInfo() *PosInfo {
	return &p.PosInfo
}

type PosInfo struct {
	X     float64
	Y     float64
	Z     float64
	Yaw   float32
	Pitch float32
}

type ProfileProperty struct {
	Name      string
	Value     string
	Signature *string
}
