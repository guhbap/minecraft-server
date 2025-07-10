package stateplay

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/uuid"
	"github.com/golangmc/minecraft-server/impl/base"
)

// Player Info Remove

// Used by the server to remove players from the player list.
// Packet ID 	State 	Bound To 	Field Name 	Field Type 	Notes
// protocol:
// 0x3F

// resource:
// player_info_remove 	Play 	Client 	UUIDs 	Prefixed Array of UUID 	UUIDs of players to remove.

type PacketOPlayerInfoRemove struct {
	UUIDs []uuid.UUID
}

func (p *PacketOPlayerInfoRemove) UUID() int32 {
	return 0x3F
}

func (p *PacketOPlayerInfoRemove) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(int32(len(p.UUIDs)))
	for _, uuid := range p.UUIDs {
		writer.PushUID(uuid)
	}
}
