package mode

import (
	"github.com/golangmc/minecraft-server/impl/conn"
)

type ChunkSection struct {
	Blocks [][][]byte
}

func (p *ChunkSection) GetBlock(x, y, z int) byte {
	return p.Blocks[y][x][7-z]
}

func (p *ChunkSection) SetBlock(x, y, z int, block byte) {
	x = 15 - x
	realX := x / 2
	isHigh := x%2 == 1 // старшая или младшая половина байта

	if isHigh {
		p.Blocks[y][z][realX] &= 0xF0 // 11110000
		p.Blocks[y][z][realX] |= block & 0x0F
	} else {
		p.Blocks[y][z][realX] &= 0x0F // 00001111
		p.Blocks[y][z][realX] |= (block & 0x0F) << 4
	}

}

func NewChunkSection() *ChunkSection {
	chs := &ChunkSection{}
	chs.Blocks = make([][][]byte, 16)
	for i := 0; i < 16; i++ {
		chs.Blocks[i] = make([][]byte, 16)
		for j := 0; j < 16; j++ {
			chs.Blocks[i][j] = make([]byte, 8)
		}
	}
	return chs
}
func (chs *ChunkSection) Push(buf *conn.ConnBuffer) {
	buf.PushVrI(256)
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			for k := 0; k < 8; k++ {
				buf.PushByt(chs.Blocks[i][j][k])
			}
		}
	}
}
