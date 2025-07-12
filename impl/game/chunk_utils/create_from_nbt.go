package chunk_utils

import (
	"bytes"
	"fmt"

	"github.com/Tnze/go-mc/nbt"
	"github.com/golangmc/minecraft-server/impl/conn"
	"github.com/golangmc/minecraft-server/impl/game/registry/biome"
	"github.com/golangmc/minecraft-server/impl/game/utils"
)

func init() {

	needData := utils.ReadHexFile("testData/levelFuullHeight.hex")
	// os.WriteFile("testData/levelFuullHeight.bin", needData, 0644)

	needDataBuff := conn.ConnBuffer{}
	needDataBuff.PushUAS(needData, false)

	needDataNbt := LoadChunk(0, 0)
	fmt.Println("needDataNbt", needDataNbt)
	fmt.Println("--------------------------------")
	fmt.Println("")
	// ParseChunk(needDataBuff)
	var hexString string
	hexString = utils.HexToString(CreateFromNbt(*needDataNbt))
	if hexString != "" {
		fmt.Println("hexString", hexString)
	}

	bool, int, string := deepCompareByteArrays(needData[1:], CreateFromNbt(*needDataNbt))
	fmt.Println("bool", bool)
	fmt.Println("int", int)
	fmt.Println("string", string)

}

func TestHeightmapSerialize() {
	needDataNbt := LoadChunk(0, 0)
	buf := conn.ConnBuffer{}
	needDataNbt.Heightmaps.EncodeTo(&buf)
	fmt.Println("buf", buf.UAS())

	copyBuf := buf.Copy()
	dec := nbt.NewDecoder(bytes.NewReader(copyBuf.UAS()))
	dec.NetworkFormat(true)
	var h2 Heightmap
	_, err := dec.Decode(&h2)
	if err != nil {
		panic(err)
	}

	buf2 := conn.ConnBuffer{}
	h2.EncodeTo(&buf2)

	fmt.Println(bytes.Equal(buf.UAS(), buf2.UAS()))
}

func ParseChunk(buf conn.ConnBuffer) ChunkNbt {
	chunk := ChunkNbt{}
	chunk.XPos = buf.PullI32()
	chunk.ZPos = buf.PullI32()
	chunk.Heightmaps = Heightmap{}

	copyBuf := buf.CopyI()
	dec := nbt.NewDecoder(bytes.NewReader(copyBuf.UAS()))
	dec.NetworkFormat(true)
	_, err := dec.Decode(&chunk.Heightmaps)
	if err != nil {
		panic(err)
	}
	fmt.Println("chunk.Heightmaps", chunk.Heightmaps)

	return chunk
}

func CreateFromNbt(nbtData ChunkNbt) []byte {
	buf := conn.ConnBuffer{}
	buf.PushI32(int32(nbtData.XPos))
	buf.PushI32(int32(nbtData.ZPos))
	{ // data
		{ // Heightmaps
			nbtBuf := bytes.NewBuffer(make([]byte, 0))
			enc := nbt.NewEncoder(nbtBuf)
			enc.NetworkFormat(true)
			fmt.Println("nbtData.Heightmaps", nbtData.Heightmaps)
			err := enc.Encode(&nbtData.Heightmaps, "")
			if err != nil {
				panic(err)
			}
			buf.PushUAS(nbtBuf.Bytes(), false)
		}

		tmpBuf := conn.ConnBuffer{}

		skys := [][]byte{}
		blys := [][]byte{}

		skyBitMask := int64(0)
		emptySkyBitMask := int64(0)
		blockBitMask := int64(0)
		emptyBlockBitMask := int64(0)

		for _, section := range nbtData.Sections[:] {
			if section.Y != -5 {
				pushSection(&tmpBuf, section)
			}
		}
		for i, section := range nbtData.Sections[:] {
			if section.SkyLight != nil {
				skys = append(skys, section.SkyLight)
				skyBitMask |= 1 << int64(i)
				// fmt.Println("section.SkyLight", section.SkyLight)
			} else {
				emptySkyBitMask |= 1 << int64(i)
			}
			if section.BlockLight != nil {
				blys = append(blys, section.BlockLight)
				blockBitMask |= 1 << int64(i)
			} else {
				emptyBlockBitMask |= 1 << int64(i)
			}
		}

		fmt.Println("skysLength", len(skys))

		buf.PushVrI(tmpBuf.Len())
		buf.PushUAS(tmpBuf.UAS(), false)
		buf.PushByt(0)

		if skyBitMask != 0 {
			buf.PushByt(1) // Sky Light Mask
			// buf.PushI64(1048575)
			buf.PushI64(skyBitMask)

		} else {
			buf.PushByt(0) // Empty Sky Light Mask
		}

		if blockBitMask != 0 {
			buf.PushByt(1) // Block Light Mask
			// buf.PushI64(6)
			buf.PushI64(blockBitMask)
		} else {
			buf.PushByt(0) // Empty Block Light Mask
		}

		// buf.PushByt(1) // Block Light Mask
		// buf.PushI64(6)
		// buf.PushI64(blockBitMask)

		buf.PushByt(0) // Empty Sky Light Mask

		if emptyBlockBitMask != 0 {
			buf.PushByt(1) // Empty Block Light Mask
			buf.PushI64(emptyBlockBitMask)
		} else {
			buf.PushByt(0) // Empty Block Light Mask
		}

		buf.PushVrI(int32(len(skys))) // Sky Light arrays length
		for _, sky := range skys {
			buf.PushUAS(sky, true)
		}
		buf.PushVrI(int32(len(blys))) // Block Light arrays length
		for _, bly := range blys {
			buf.PushUAS(bly, true)
		}

	}

	return buf.UAS()
}

func pushSection(buf *conn.ConnBuffer, section ChunkSectionNbt) {

	tempBuf := conn.ConnBuffer{}
	blocks := []ChunkSectionBlocksPaletteNbt{}
	for _, block := range section.BlockStates.Palette {
		blocks = append(blocks, block)
	}
	pallete := *NewPallete(blocks)
	// fmt.Println("section.BlockStates.Palette", section.BlockStates.Palette)
	// fmt.Println("Y", section.Y)
	fmt.Println("pallete", pallete)
	pallete.Push(&tempBuf)

	// if len(pallete.Blocks) > 16 {
	// 	fmt.Println("pallete", pallete)
	// 	t2buf := conn.ConnBuffer{}
	// 	pallete.Push(&t2buf)
	// 	fmt.Println("t2buf", hex.EncodeToString(t2buf.UAS()))
	// }

	// return

	bitsPerBlock := uint(4)
	// Подсчет непустых блоков
	nonAirCount := 0
	blocksPerLong := uint(64 / bitsPerBlock)
	mask := (1 << bitsPerBlock) - 1 // Маска для извлечения битов
	airIndex := 0

	tempBuf.PushVrI(int32(len(section.BlockStates.Data)))
	for _, block := range section.BlockStates.Data {
		tempBuf.PushI64(block)
		for i := uint(0); i < blocksPerLong; i++ {
			blockIndex := int((block >> (i * bitsPerBlock)) & int64(mask))
			if blockIndex != airIndex {
				nonAirCount++
			}
		}
	}
	pallete = Pallete{
		Blocks: []int{},
	}
	for _, block := range section.Biomes.Palette {
		pallete.Blocks = append(pallete.Blocks, int(biome.GetBiomeId(block)))
	}
	pallete.Push(&tempBuf)

	tempBuf.PushByt(0)
	buf.PushI16(int16(nonAirCount)) // непонятно нихуя. Сколько блоков надо обновить?
	buf.PushUAS(tempBuf.UAS(), false)
	return

}

type SkyLight struct {
	Light []byte
}

func (p *SkyLight) Push(buf *conn.ConnBuffer) {
	buf.PushUAS(p.Light, true)
}

func NewSkyLight() *SkyLight {
	skyLight := &SkyLight{}
	skyLight.Light = make([]byte, 2048)
	for i := 0; i < 2048; i++ {
		skyLight.Light[i] = 0xff
	}
	return skyLight
}

func countZeroBitGroups(data []int64, groupSize int) int16 {
	if groupSize <= 0 || groupSize > 64 {
		panic("groupSize должен быть от 1 до 64")
	}

	count := 0
	totalBits := len(data) * 64

	// Проходим по битам с шагом groupSize (например, 2)
	for i := 0; i+groupSize <= totalBits; i += groupSize {
		allZero := true

		for j := 0; j < groupSize; j++ {
			bitIndex := i + j
			word := data[bitIndex/64]
			bit := (word >> (63 - (bitIndex % 64))) & 1

			if bit == 1 {
				allZero = false
				break
			}
		}

		if allZero {
			count++
		}
	}

	return int16(count)
}
