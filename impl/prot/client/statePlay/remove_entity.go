package stateplay

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/impl/base"
)

// Remove Entities

// Sent by the server when an entity is to be destroyed on the client.
// Packet ID 	State 	Bound To 	Field Name 	Field Type 	Notes
// protocol:
// 0x47

// resource:
// remove_entities 	Play 	Client 	Entity IDs 	Prefixed Array of VarInt 	The list of entities to destroy.

type PacketORemoveEntity struct {
	EntityIDs []int32
}

func (p *PacketORemoveEntity) UUID() int32 {
	return 0x47
}

func (p *PacketORemoveEntity) Push(writer buff.Buffer, conn base.Connection) {

	writer.PushVrI(int32(len(p.EntityIDs)))
	for _, entityID := range p.EntityIDs {
		writer.PushVrI(entityID)
	}
}
