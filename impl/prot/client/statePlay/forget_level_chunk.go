package stateplay

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/impl/base"
)

type PacketOForgetLevelChunk struct {
	X int32
	Z int32
}

func (p *PacketOForgetLevelChunk) UUID() int32 {
	return 0x22
}

func (p *PacketOForgetLevelChunk) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushI32(p.Z)
	writer.PushI32(p.X)
}
