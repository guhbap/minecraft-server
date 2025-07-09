package server

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/data/client"
)

type PacketICustomPayload struct {
	Channel string
	Data    []byte
}

func (p *PacketICustomPayload) UUID() int32 {
	return 0x02
}

func (p *PacketICustomPayload) Pull(reader buff.Buffer, conn base.Connection) {
	p.Channel = reader.PullTxt()
	p.Data = reader.PullUAS()
}

type KnownPack struct {
	Namespace string
	Id        string
	Version   string
}

type PacketISelectKnownPacks struct {
	KnownPacks []KnownPack
}

func (p *PacketISelectKnownPacks) UUID() int32 {
	return 0x07
}

func (p *PacketISelectKnownPacks) Pull(reader buff.Buffer, conn base.Connection) {
	count := reader.PullVrI()
	for i := 0; i < int(count); i++ {
		pack := KnownPack{}
		pack.Namespace = reader.PullTxt()
		pack.Id = reader.PullTxt()
		pack.Version = reader.PullTxt()
		p.KnownPacks = append(p.KnownPacks, pack)
	}
}

type PacketIFinishConfiguration struct {
}

func (p *PacketIFinishConfiguration) UUID() int32 {
	return 0x03
}

func (p *PacketIFinishConfiguration) Push(writer buff.Buffer, conn base.Connection) {
}

func (p *PacketIFinishConfiguration) Pull(reader buff.Buffer, conn base.Connection) {
}

type PacketIClientInformation struct {
	Locale              string
	ViewDistance        byte
	ChatMode            client.ChatMode
	ChatColors          bool // if false, strip messages of colors before sending
	SkinParts           client.SkinParts
	MainHand            client.MainHand
	EnableTextFiltering bool
	AllowServerListings bool
	ParticleStatus      int32
}

func (p *PacketIClientInformation) UUID() int32 {
	return 0x05
}

func (p *PacketIClientInformation) Pull(reader buff.Buffer, conn base.Connection) {
	p.Locale = reader.PullTxt()
	p.ViewDistance = reader.PullByt()
	p.ChatMode = client.ChatMode(reader.PullVrI())
	p.ChatColors = reader.PullBit()

	parts := client.SkinParts{}
	parts.Pull(reader)

	p.SkinParts = parts
	p.MainHand = client.MainHand(reader.PullVrI())
	p.EnableTextFiltering = reader.PullBit()
	p.AllowServerListings = reader.PullBit()
	p.ParticleStatus = reader.PullVrI()
}
