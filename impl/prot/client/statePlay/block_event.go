package stateplay

// todo
import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/data"
	"github.com/golangmc/minecraft-server/impl/base"
)

type BlockEventAction byte

const (
	BlockEventActionBreak BlockEventAction = iota
	BlockEventActionPlace
)

type BlockType byte

const (
	BlockTypeAir BlockType = iota
	BlockTypeStone
	BlockTypeGrass
)

type PacketOBlockEvent struct {
	Location        data.PositionI
	ActionId        BlockEventAction
	ActionParameter BlockEventAction
	BlockType       BlockType
}

func (p *PacketOBlockEvent) UUID() int32 {
	return 0x08
}

func (p *PacketOBlockEvent) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushPos(p.Location)
	writer.PushI32(int32(p.ActionId))
	writer.PushI32(int32(p.ActionParameter))
	writer.PushI32(int32(p.BlockType))
}
