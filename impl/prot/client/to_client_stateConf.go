package client

import (
	"bytes"

	"github.com/Tnze/go-mc/nbt"
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/impl/base"
)

type PacketOFinishConfiguration struct {
}

func (p *PacketOFinishConfiguration) UUID() int32 {
	return 0x03
}

func (p *PacketOFinishConfiguration) Push(writer buff.Buffer, conn base.Connection) {
}

type PacketOCustomPayload struct {
	Channel string
	Data    []byte
}

func (p *PacketOCustomPayload) UUID() int32 {
	return 0x01
}

func (p *PacketOCustomPayload) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushTxt(p.Channel)
	writer.PushUAS(p.Data, true)
}

type PacketOUpdateEnabledFeatures struct {
	Features []string
}

func (p *PacketOUpdateEnabledFeatures) UUID() int32 {
	return 0x0C
}

func (p *PacketOUpdateEnabledFeatures) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(int32(len(p.Features)))
	for _, feature := range p.Features {
		writer.PushTxt(feature)
	}
}

type KnownPack struct {
	Namespace string
	Id        string
	Version   string
}

func (p *KnownPack) Push(writer buff.Buffer) {
	writer.PushTxt(p.Namespace)
	writer.PushTxt(p.Id)
	writer.PushTxt(p.Version)
}

type PacketOSelectKnownPacks struct {
	KnownPacks []KnownPack
}

func (p *PacketOSelectKnownPacks) UUID() int32 {
	return 0x0E
}

func (p *PacketOSelectKnownPacks) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(int32(len(p.KnownPacks)))
	for _, pack := range p.KnownPacks {
		pack.Push(writer)
	}
}

type RegistryEntry struct {
	Id    string
	Value interface{}
}

func (e *RegistryEntry) Push(writer buff.Buffer) {
	// writer.PushBit(true)
	writer.PushTxt(e.Id)
	writer.PushBit(true)
	buf := bytes.Buffer{}
	enc := nbt.NewEncoder(&buf)
	enc.NetworkFormat(true)
	err := enc.Encode(e.Value, e.Id)
	if err != nil {
		panic(err)
	}
	writer.PushUAS(buf.Bytes(), false)
}

type PacketORegistryData struct {
	Id      string
	Entries []RegistryEntry
}

func (p *PacketORegistryData) UUID() int32 {
	return 0x07
}

func (p *PacketORegistryData) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushTxt(p.Id)
	writer.PushVrI(int32(len(p.Entries)))
	for _, entry := range p.Entries {
		entry.Push(writer)
	}
}

type PacketOLightUpdate struct {
	ChunkX              int32
	ChunkZ              int32
	SkyLightMask        int32
	BlockLightMask      int32
	EmptySkyLightMask   int32
	EmptyBlockLightMask int32
	SkyLight            int32
	BlockLight          int32
}

func (p *PacketOLightUpdate) UUID() int32 {
	return 0x2A
}

func (p *PacketOLightUpdate) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushI32(p.ChunkX)
	writer.PushI32(p.ChunkZ)
	writer.PushI32(p.SkyLightMask)
	writer.PushI32(p.BlockLightMask)
	writer.PushI32(p.EmptySkyLightMask)
	writer.PushI32(p.EmptyBlockLightMask)
	writer.PushI32(p.SkyLight)
	writer.PushI32(p.BlockLight)
}

type PacketOUpdateTags struct {
	RawData []byte
}

func (p *PacketOUpdateTags) UUID() int32 {
	return 0x0d
}

func (p *PacketOUpdateTags) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushUAS(p.RawData, false)
}
