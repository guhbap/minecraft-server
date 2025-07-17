package stateplay

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/impl/base"
)

type Animation byte

// 0 	Swing main arm
// 2 	Leave bed
// 3 	Swing offhand
// 4 	Critical effect
// 5 	Magic critical effect
const (
	AnimationSwing Animation = iota
	AnimationLeaveBed
	AnimationSwingOffhand
	AnimationCriticalEffect
	AnimationMagicCriticalEffect
)

type PacketOAnimate struct {
	EntityID  int32
	Animation Animation
}

func (p *PacketOAnimate) UUID() int32 {
	return 0x03
}

func (p *PacketOAnimate) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(p.EntityID)
	writer.PushByt(byte(p.Animation))
}
