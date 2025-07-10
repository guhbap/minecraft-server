package server

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/data"
	"github.com/golangmc/minecraft-server/apis/game"
	"github.com/golangmc/minecraft-server/apis/uuid"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/data/client"
	"github.com/golangmc/minecraft-server/impl/data/plugin"
)

type PacketIKeepAlive struct {
	KeepAliveID int64
}

func (p *PacketIKeepAlive) UUID() int32 {
	return 0x1A
}

func (p *PacketIKeepAlive) Pull(reader buff.Buffer, conn base.Connection) {
	p.KeepAliveID = reader.PullI64()
}

type PacketIChatMessage struct {
	Message string
}

func (p *PacketIChatMessage) UUID() int32 {
	return 0x03
}

func (p *PacketIChatMessage) Pull(reader buff.Buffer, conn base.Connection) {
	p.Message = reader.PullTxt()
}

type PacketITeleportConfirm struct {
	TeleportID int32
}

func (p *PacketITeleportConfirm) UUID() int32 {
	return 0x00
}

func (p *PacketITeleportConfirm) Pull(reader buff.Buffer, conn base.Connection) {
	p.TeleportID = reader.PullVrI()
}

type PacketIQueryBlockNBT struct {
	TransactionID int32
	Position      data.PositionI
}

func (p *PacketIQueryBlockNBT) UUID() int32 {
	return 0x01
}

func (p *PacketIQueryBlockNBT) Pull(reader buff.Buffer, conn base.Connection) {
	p.TransactionID = reader.PullVrI()
	p.Position = reader.PullPos()
}

type PacketISetDifficulty struct {
	Difficult game.Difficulty
}

func (p *PacketISetDifficulty) UUID() int32 {
	return 0x02
}

func (p *PacketISetDifficulty) Pull(reader buff.Buffer, conn base.Connection) {
	p.Difficult = game.DifficultyValueOf(reader.PullByt())
}

type PacketIPluginMessage struct {
	Message plugin.Message
}

func (p *PacketIPluginMessage) UUID() int32 {
	return 0x0B
}

func (p *PacketIPluginMessage) Pull(reader buff.Buffer, conn base.Connection) {
	channel := reader.PullTxt()
	message := plugin.GetMessageForChannel(channel)

	if message == nil {
		return // log unregistered channel?
	}

	message.Pull(reader)

	p.Message = message
}

type PacketIClientStatus struct {
	Action client.StatusAction
}

func (p *PacketIClientStatus) UUID() int32 {
	return 0x04
}

func (p *PacketIClientStatus) Pull(reader buff.Buffer, conn base.Connection) {
	p.Action = client.StatusAction(reader.PullVrI())
}

type PacketIPlayerAbilities struct {
	Flags byte
}

func (p *PacketIPlayerAbilities) UUID() int32 {
	return 0x26
}

func (p *PacketIPlayerAbilities) Pull(reader buff.Buffer, conn base.Connection) {
	p.Flags = reader.PullByt()
}

type PacketIPlayerPosition struct {
	Position data.PositionF
	OnGround bool
}

func (p *PacketIPlayerPosition) UUID() int32 {
	return 0x11
}

func (p *PacketIPlayerPosition) Pull(reader buff.Buffer, conn base.Connection) {
	p.Position = data.PositionF{
		X: reader.PullF64(),
		Y: reader.PullF64(),
		Z: reader.PullF64(),
	}

	p.OnGround = reader.PullBit()
}

type PacketIPlayerLocation struct {
	Location data.Location
	OnGround bool
}

func (p *PacketIPlayerLocation) UUID() int32 {
	return 0x12
}

func (p *PacketIPlayerLocation) Pull(reader buff.Buffer, conn base.Connection) {
	p.Location = data.Location{
		PositionF: data.PositionF{
			X: reader.PullF64(),
			Y: reader.PullF64(),
			Z: reader.PullF64(),
		},
		RotationF: data.RotationF{
			AxisX: reader.PullF32(),
			AxisY: reader.PullF32(),
		},
	}

	p.OnGround = reader.PullBit()
}

type PacketIPlayerRotation struct {
	Rotation data.RotationF
	OnGround bool
}

func (p *PacketIPlayerRotation) UUID() int32 {
	return 0x13
}

func (p *PacketIPlayerRotation) Pull(reader buff.Buffer, conn base.Connection) {
	p.Rotation = data.RotationF{
		AxisX: reader.PullF32(),
		AxisY: reader.PullF32(),
	}

	p.OnGround = reader.PullBit()
}

type PacketIChatSessionUpdate struct {
	SessionID uuid.UUID
	PublicKey struct {
		ExpyredAt    int64
		Key          []byte
		KeySignature []byte
	}
}

func (p *PacketIChatSessionUpdate) UUID() int32 {
	return 0x08
}

func (p *PacketIChatSessionUpdate) Pull(reader buff.Buffer, conn base.Connection) {
	p.SessionID = reader.PullUID()
	p.PublicKey.ExpyredAt = reader.PullI64()
	p.PublicKey.Key = reader.PullUAS()
	p.PublicKey.KeySignature = reader.PullUAS()
}

type PacketIClientTickEnd struct {
	TickDelta int32
}

func (p *PacketIClientTickEnd) UUID() int32 {
	return 0x0b
}

func (p *PacketIClientTickEnd) Pull(reader buff.Buffer, conn base.Connection) {
	p.TickDelta = reader.PullVrI()
}

type PacketIMovePlayerPosRot struct {
	Position data.PositionF
	Rotation data.RotationF
	Flags    byte
}

func (p *PacketIMovePlayerPosRot) UUID() int32 {
	return 0x1d
}

func (p *PacketIMovePlayerPosRot) Pull(reader buff.Buffer, conn base.Connection) {
	p.Position = data.PositionF{
		X: reader.PullF64(),
		Y: reader.PullF64(),
		Z: reader.PullF64(),
	}

	p.Rotation = data.RotationF{
		AxisX: reader.PullF32(),
		AxisY: reader.PullF32(),
	}

	p.Flags = reader.PullByt()
}

// type PacketIMovePlayerRot struct {
// 	Yaw   float32
// 	Pitch float32
// 	Flags byte
// }

// func (p *PacketIMovePlayerRot) UUID() int32 {
// 	return 0x1e
// }
// func (p *PacketIMovePlayerRot) Pull(reader buff.Buffer, conn base.Connection) {
// 	p.Yaw = reader.PullF32()
// 	p.Pitch = reader.PullF32()
// 	p.Flags = reader.PullByt()
// }

type PacketIMovePlayerPos struct {
	Position data.PositionF
	Flags    byte
}

func (p *PacketIMovePlayerPos) UUID() int32 {
	return 0x1c
}

func (p *PacketIMovePlayerPos) Pull(reader buff.Buffer, conn base.Connection) {
	p.Position = data.PositionF{
		X: reader.PullF64(),
		Y: reader.PullF64(),
		Z: reader.PullF64(),
	}

	p.Flags = reader.PullByt()
}

type PacketIPlayerInput struct {
	Flags byte
}

func (p *PacketIPlayerInput) UUID() int32 {
	return 0x29
}

func (p *PacketIPlayerInput) Pull(reader buff.Buffer, conn base.Connection) {
	p.Flags = reader.PullByt()
}

type PacketIChatCommand struct {
	Command string
}

func (p *PacketIChatCommand) UUID() int32 {
	return 0x05
}

func (p *PacketIChatCommand) Pull(reader buff.Buffer, conn base.Connection) {
	p.Command = reader.PullTxt()
}

type PacketIPlayerLoaded struct {
}

func (p *PacketIPlayerLoaded) UUID() int32 {
	return 0x2a
}

func (p *PacketIPlayerLoaded) Pull(reader buff.Buffer, conn base.Connection) {
}

type PacketIPlayerCommand struct {
	EntityID  int32
	ActionId  int32
	JumpBoost int32
}

func (p *PacketIPlayerCommand) UUID() int32 {
	return 0x28
}

func (p *PacketIPlayerCommand) Pull(reader buff.Buffer, conn base.Connection) {
	p.EntityID = reader.PullVrI()
	p.ActionId = reader.PullVrI()
	p.JumpBoost = reader.PullVrI()
}

type PacketIMovePlayerRot struct {
	Yaw   float32
	Pitch float32
	Flags byte
}

func (p *PacketIMovePlayerRot) UUID() int32 {
	return 0x1e
}

func (p *PacketIMovePlayerRot) Pull(reader buff.Buffer, conn base.Connection) {
	p.Yaw = reader.PullF32()
	p.Pitch = reader.PullF32()
	p.Flags = reader.PullByt()
}

type PacketIContainerClose struct {
	ContainerID int32
}

func (p *PacketIContainerClose) UUID() int32 {
	return 0x11
}

func (p *PacketIContainerClose) Pull(reader buff.Buffer, conn base.Connection) {
	p.ContainerID = reader.PullVrI()
}

type PacketIBundleDelimiter struct {
}

func (p *PacketIBundleDelimiter) UUID() int32 {
	return 0x00
}

func (p *PacketIBundleDelimiter) Pull(reader buff.Buffer, conn base.Connection) {
}

type PacketIChunkBatchReceived struct {
	ChunkOnTick float32
}

func (p *PacketIChunkBatchReceived) UUID() int32 {
	return 0x09
}
func (p *PacketIChunkBatchReceived) Pull(reader buff.Buffer, conn base.Connection) {
	p.ChunkOnTick = reader.PullF32()
}
