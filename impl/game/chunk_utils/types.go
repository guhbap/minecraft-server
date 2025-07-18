package chunk_utils

import (
	"bytes"
	"strconv"

	"github.com/Tnze/go-mc/nbt"
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/impl/conn"
	biomeFinder "github.com/golangmc/minecraft-server/impl/game/registry/biome"
	blockFinder "github.com/golangmc/minecraft-server/impl/game/registry/block"
)

type ChunkNbt struct {
	DataVersion   int32             `nbt:"DataVersion"`
	XPos          int32             `nbt:"xPos"`
	ZPos          int32             `nbt:"zPos"`
	YPos          int32             `nbt:"yPos"`
	Status        string            `nbt:"Status"`
	LastUpdate    int64             `nbt:"LastUpdate"`
	Sections      []ChunkSectionNbt `nbt:"sections"`
	BlockEntities []any             `nbt:"block_entities"`
	Heightmaps    Heightmap         `nbt:"Heightmaps"`
}
type ChunkSectionNbt struct {
	Y           int8                       `nbt:"Y"`
	BlockStates ChunkSectionBlockStatesNbt `nbt:"block_states"`
	Biomes      ChunkSectionBiomesNbt      `nbt:"biomes"`
	BlockLight  []byte                     `nbt:"BlockLight"`
	SkyLight    []byte                     `nbt:"SkyLight"`
}
type ChunkSectionBlockStatesNbt struct {
	Palette []ChunkSectionBlocksPaletteNbt `nbt:"palette"`
	Data    []int64                        `nbt:"data"`
}

type ChunkSectionBlocksPaletteNbt struct {
	Name       string            `nbt:"Name"`
	Properties map[string]string `nbt:"Properties"`
}
type ChunkSectionBiomesPaletteNbt struct {
	Name string `nbt:"Name"`
}

type ChunkSectionBiomesNbt struct {
	Palette []string `nbt:"palette"`
	Data    []int64  `nbt:"data"`
}
type Heightmap struct {
	MotionBlocking []int64 `nbt:"MOTION_BLOCKING"`
	WorldSurface   []int64 `nbt:"WORLD_SURFACE"`
}

func (h *Heightmap) DecodeFrom(buf *conn.ConnBuffer) {
	copyBuf := buf.CopyI()
	dec := nbt.NewDecoder(bytes.NewReader(copyBuf.UAS()))
	dec.NetworkFormat(true)
	_, err := dec.Decode(&h)
	if err != nil {
		panic(err)
	}
}
func (h *Heightmap) EncodeTo(buf *conn.ConnBuffer) {
	enc := nbt.NewEncoder(buf)
	enc.NetworkFormat(true)
	err := enc.Encode(h, "Heightmaps")
	if err != nil {
		panic(err)
	}
}

type Pallete struct {
	Blocks    []int
	strBlocks []string
	AirIndex  int

	BitsPerBlock byte
}

func NewPallete(blocks []ChunkSectionBlocksPaletteNbt) *Pallete {
	palette := &Pallete{}
	palette.AirIndex = -1
	for _, block := range blocks {
		bl, err := blockFinder.GetBlockID(block.Name, block.Properties)
		if err != nil {
			panic(err)
		}
		if block.Name == "minecraft:air" {
			palette.AirIndex = len(palette.Blocks)
		}
		palette.Blocks = append(palette.Blocks, bl)
		palette.strBlocks = append(palette.strBlocks, block.Name)
	}
	if len(palette.Blocks) == 1 {
		palette.BitsPerBlock = 0
	} else if len(palette.Blocks) < 17 {
		palette.BitsPerBlock = 4
	} else if len(palette.Blocks) < 33 {
		palette.BitsPerBlock = 5
	} else if len(palette.Blocks) < 65 {
		palette.BitsPerBlock = 6
	} else if len(palette.Blocks) < 129 {
		palette.BitsPerBlock = 7
	} else {
		panic("not implemented for len(palette.Blocks) = " + strconv.Itoa(len(palette.Blocks)))
	}
	return palette
}

func NewBiomesPallete(biomes []string) *Pallete {
	palette := &Pallete{}
	for _, biome := range biomes {
		palette.Blocks = append(palette.Blocks, biomeFinder.GetBiomeId(biome))
	}
	if len(palette.Blocks) == 1 {
		palette.BitsPerBlock = 0
	} else if len(palette.Blocks) < 17 {
		palette.BitsPerBlock = 4
	} else if len(palette.Blocks) < 33 {
		palette.BitsPerBlock = 5
	} else if len(palette.Blocks) < 65 {
		palette.BitsPerBlock = 6
	} else if len(palette.Blocks) < 129 {
		palette.BitsPerBlock = 7
	} else {
		panic("not implemented for len(palette.Blocks) = " + strconv.Itoa(len(palette.Blocks)))
	}
	return palette
}

func (p *Pallete) Push(buf buff.Buffer, isBiome bool) {
	if p.BitsPerBlock == 0 {
		buf.PushByt(0)
		buf.PushVrI(int32(p.Blocks[0]))
	} else {
		var typeByte byte
		blockCount := len(p.Blocks)
		if isBiome {
			switch {
			case blockCount < 3:
				typeByte = 1
			case blockCount < 5:
				typeByte = 2
			case blockCount < 9:
				typeByte = 3
			}
		}
		if typeByte == 0 {
			switch {
			case blockCount < 17:
				typeByte = 4
			case blockCount < 33:
				typeByte = 5
			case blockCount < 65:
				typeByte = 6
			case blockCount < 129:
				typeByte = 7
			}
		}
		if typeByte == 0 {
			panic("not implemented for len(p.Blocks) = " + strconv.Itoa(len(p.Blocks)))
		}

		buf.PushByt(typeByte)
		buf.PushVrI(int32(len(p.Blocks)))
		for _, block := range p.Blocks {
			buf.PushVrI(int32(block))
		}
	}

	// if isBiome {
	// 	if len(p.Blocks) > 1 && len(p.Blocks) < 3 {
	// 		buf.PushByt(1)
	// 		buf.PushVrI(int32(len(p.Blocks)))
	// 		for _, block := range p.Blocks {
	// 			buf.PushVrI(int32(block))
	// 		}
	// 	} else if len(p.Blocks) > 3 && len(p.Blocks) < 5 {
	// 		buf.PushByt(2)
	// 		buf.PushVrI(int32(len(p.Blocks)))
	// 		for _, block := range p.Blocks {
	// 			buf.PushVrI(int32(block))
	// 		}
	// 	} else if len(p.Blocks) > 5 && len(p.Blocks) < 17 {
	// 		buf.PushByt(3)
	// 		buf.PushVrI(int32(len(p.Blocks)))
	// 		for _, block := range p.Blocks {
	// 			buf.PushVrI(int32(block))
	// 		}
	// 	}
	// }
	// if len(p.Blocks) > 1 && len(p.Blocks) < 17 {
	// 	buf.PushByt(4)
	// 	buf.PushVrI(int32(len(p.Blocks)))
	// 	for _, block := range p.Blocks {
	// 		buf.PushVrI(int32(block))
	// 	}
	// } else if len(p.Blocks) > 16 && len(p.Blocks) < 33 {
	// 	buf.PushByt(5)
	// 	buf.PushVrI(int32(len(p.Blocks)))
	// 	for _, block := range p.Blocks {
	// 		buf.PushVrI(int32(block))
	// 	}
	// } else if len(p.Blocks) > 32 && len(p.Blocks) < 65 {
	// 	buf.PushByt(6)
	// 	buf.PushVrI(int32(len(p.Blocks)))
	// 	for _, block := range p.Blocks {
	// 		buf.PushVrI(int32(block))
	// 	}
	// } else if len(p.Blocks) == 1 {
	// 	buf.PushByt(0)
	// 	buf.PushVrI(int32(p.Blocks[0]))
	// } else {
	// 	panic("not implemented for len(p.Blocks) = " + strconv.Itoa(len(p.Blocks)))
	// }
}

type Pallete0 struct {
	Block int
}

func (p *Pallete0) Push(buf *conn.ConnBuffer) {
	buf.PushByt(0)
	buf.PushVrI(int32(p.Block))
}
func (p *Pallete0) Pull(buf *conn.ConnBuffer) {
	p.Block = int(buf.PullVrI())
}

type Palette4 struct {
	Blocks []int
}

func (p *Palette4) Push(buf *conn.ConnBuffer) {
	buf.PushByt(4)
	buf.PushVrI(int32(len(p.Blocks)))
	for _, block := range p.Blocks {
		buf.PushVrI(int32(block))
	}
}
func (p *Palette4) Pull(buf *conn.ConnBuffer) {
	len := int(buf.PullVrI())
	for i := 0; i < len; i++ {
		p.Blocks = append(p.Blocks, int(buf.PullVrI()))
	}
}

// func (p *Palette4) Print() string {
// 	blocks := []string{}
// 	for _, block := range p.Blocks {
// 		blocks = append(blocks, registry.Blocks[int(block)])
// 	}
// 	return fmt.Sprintf("Palette: %v", blocks)
// }
