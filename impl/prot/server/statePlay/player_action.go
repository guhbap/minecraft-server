package stateplay

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/data"
	"github.com/golangmc/minecraft-server/impl/base"
)

type PlayerActionStatus int

const (
	PlayerActionStatusStartedDigging = iota
	PlayerActionStatusCancelledDigging
	PlayerActionStatusFinishedDigging
	PlayerActionStatusDropItemStack
	PlayerActionStatusDropItem
	PlayerActionStatusShootArrow
	PlayerActionStatusSwapItemInHand
)

type PlayerActionFace byte

const (
	PlayerActionFaceBottom = iota
	PlayerActionFaceTop    // 1
	PlayerActionFaceNorth  // 2
	PlayerActionFaceSouth  // 3
	PlayerActionFaceWest   // 4
	PlayerActionFaceEast   // 5
)

type PacketIPlayerAction struct {
	Status   PlayerActionStatus
	Location data.PositionI
	Face     PlayerActionFace
	Sequence int32
}

func (p *PacketIPlayerAction) Pull(reader buff.Buffer, conn base.Connection) {
	p.Status = PlayerActionStatus(reader.PullVrI())
	p.Location = reader.PullPos()
	p.Face = PlayerActionFace(reader.PullByt())
	p.Sequence = reader.PullVrI()
}

func (p *PacketIPlayerAction) UUID() int32 {
	return 0x27
}
