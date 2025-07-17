package stateplay

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/impl/base"
)

type PacketOBlockChangedAck struct {
	Sequence int32
}

func (p *PacketOBlockChangedAck) UUID() int32 {
	return 0x05
}

func (p *PacketOBlockChangedAck) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(p.Sequence)
}
