package stateplay

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/data"
	"github.com/golangmc/minecraft-server/impl/base"
)

type PacketOBlockUpdate struct {
	Location data.PositionI
	BlockId  int32
}

func (p *PacketOBlockUpdate) UUID() int32 {
	return 0x09
}

func (p *PacketOBlockUpdate) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushPos(p.Location)
	writer.PushVrI(p.BlockId)
}
