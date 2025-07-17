package stateplay

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/impl/base"
)

type SwingHand byte

const (
	SwingHandMainHand SwingHand = iota
	SwingHandOffhand
)

type PacketISwing struct {
	Hand SwingHand
}

func (p *PacketISwing) UUID() int32 {
	return 0x3a
}

func (p *PacketISwing) Pull(reader buff.Buffer, conn base.Connection) {
	p.Hand = SwingHand(reader.PullByt())
}
