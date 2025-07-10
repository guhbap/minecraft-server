package stateplay

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/uuid"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/prot/subtypes"
)

// protocol:
// 0x01

// resource:
// add_entity 	Play 	Client 	Entity ID 	VarInt 	A unique integer ID mostly used in the protocol to identify the entity.
// Entity UUID 	UUID 	A unique identifier that is mostly used in persistence and places where the uniqueness matters more.
// Type 	VarInt 	ID in the minecraft:entity_type registry (see "type" field in Entity metadata#Entities).
// X 	Double
// Y 	Double
// Z 	Double
// Pitch 	Angle 	To get the real pitch, you must divide this by (256.0F / 360.0F)
// Yaw 	Angle 	To get the real yaw, you must divide this by (256.0F / 360.0F)
// Head Yaw 	Angle 	Only used by living entities, where the head of the entity may differ from the general body rotation.
// Data 	VarInt 	Meaning dependent on the value of the Type field, see Object Data for details.
// Velocity X 	Short 	Same un
type PacketOAddEntity struct {
	EntityID   int32
	EntityUUID uuid.UUID
	Type       int32
	X          float64
	Y          float64
	Z          float64
	Pitch      subtypes.Angle
	Yaw        subtypes.Angle
	HeadYaw    subtypes.Angle
	Data       int32
	VelocityX  int16
	VelocityY  int16
	VelocityZ  int16
}

func (p *PacketOAddEntity) UUID() int32 {
	return 0x01
}

func (p *PacketOAddEntity) Push(writer buff.Buffer, conn base.Connection) {

	writer.PushVrI(p.EntityID)
	writer.PushUID(p.EntityUUID)
	writer.PushVrI(p.Type)
	writer.PushF64(p.X)
	writer.PushF64(p.Y)
	writer.PushF64(p.Z)
	p.Pitch.Push(writer)
	p.Yaw.Push(writer)
	p.HeadYaw.Push(writer)
	writer.PushVrI(p.Data)
	writer.PushI16(p.VelocityX)
	writer.PushI16(p.VelocityY)
	writer.PushI16(p.VelocityZ)
}
