package stateplay

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/impl/base"
)

type PacketIMovePlayerStatusOnly struct {
	Flags byte
}

func (p *PacketIMovePlayerStatusOnly) UUID() int32 {
	return 0x1F
}

func (p *PacketIMovePlayerStatusOnly) Pull(reader buff.Buffer, conn base.Connection) {
	p.Flags = reader.PullByt()
}
