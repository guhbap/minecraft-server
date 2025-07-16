package chunk_utils

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/Tnze/go-mc/nbt"
	"github.com/golangmc/minecraft-server/impl/conn"
	"github.com/golangmc/minecraft-server/impl/game/utils"
)

func init() {

	needData := utils.ReadHexFile("testData/levelFuullHeight.hex")
	// os.WriteFile("testData/levelFuullHeight.bin", needData, 0644)

	needDataBuff := conn.ConnBuffer{}
	needDataBuff.PushUAS(needData, false)

	needDataNbt := LoadChunk(0, 0)
	// fmt.Println("needDataNbt", needDataNbt)
	fmt.Println("--------------------------------")
	notAirBlocks := CalculateNotAirBlocksNew(4, 9, needDataNbt.Sections[0].BlockStates.Data)
	fmt.Println("notAirBlocks", notAirBlocks)
	tBuf := conn.ConnBuffer{}
	tBuf.PushI16(int16(notAirBlocks))
	fmt.Println("tBuf", hex.EncodeToString(tBuf.UAS()))
	// return
	// ParseChunk(needDataBuff)
	var hexString string
	hexString = utils.HexToString(CreateFromNbt(*needDataNbt))
	if hexString != "" {
		fmt.Println("hexString", hexString)
	}
	res, intt, str := DeepCompareByteArrays(needData, CreateFromNbt(*needDataNbt), 3)
	fmt.Println("res", res, intt, str)

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

		// buf.PushVrI(200)
		buf.PushVrI(tmpBuf.Len())
		buf.PushUAS(tmpBuf.UAS(), false)
		buf.PushByt(0)
		// return buf.UAS()

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
			// buf.PushVrI(2048)
			// for i := 0; i < 2048; i++ {
			// 	fmt.Println("sky", sky[i])
			// 	buf.PushByt(0xff)
			// }
		}
		buf.PushVrI(int32(len(blys))) // Block Light arrays length
		for _, bly := range blys {
			buf.PushUAS(bly, true)
			// buf.PushVrI(2048)
			// for i := 0; i < 2048; i++ {
			// 	fmt.Println("bly", bly[i])
			// 	buf.PushByt(0x00)
			// }
		}

	}

	return buf.UAS()
}

func pushSection(buf *conn.ConnBuffer, section ChunkSectionNbt) {

	tempBuf := conn.ConnBuffer{}
	pallete := *NewPallete(section.BlockStates.Palette)
	// fmt.Println("section.BlockStates.Palette", section.BlockStates.Palette)
	// fmt.Println("Y", section.Y)
	pallete.Push(&tempBuf, false)

	// if len(pallete.Blocks) > 16 {
	// 	fmt.Println("pallete", pallete)
	// 	t2buf := conn.ConnBuffer{}
	// 	pallete.Push(&t2buf)
	// 	fmt.Println("t2buf", hex.EncodeToString(t2buf.UAS()))
	// }

	// return

	bitsPerBlock := uint(pallete.BitsPerBlock)
	if bitsPerBlock == 0 {
		bitsPerBlock = 1
	}
	// Подсчет непустых блоков
	nonAirCount := 0
	blocksPerLong := uint(64 / bitsPerBlock)
	mask := (1 << bitsPerBlock) - 1 // Маска для извлечения битов
	airIndex := pallete.AirIndex

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

	pallete = *NewBiomesPallete(section.Biomes.Palette)
	// pallete = Pallete{
	// 	Blocks: []int{6},
	// }
	pallete.Push(&tempBuf, true)
	// if section.Y == -2 {
	// 	tempBuf.PushByt(0xff)
	// }

	// if len(section.Biomes.Palette) > 1 {
	tempBuf.PushVrI(int32(len(section.Biomes.Data)))
	for _, biome := range section.Biomes.Data {
		tempBuf.PushI64(biome)
	}
	// if section.Y == -1 {
	// 	buf.PushByt(0xAA)
	// }
	buf.PushI16(int16(nonAirCount)) // непонятно нихуя. Сколько блоков надо обновить?
	buf.PushUAS(tempBuf.UAS(), false)

}

type SkyLight struct {
	Light []byte
}

func (p *SkyLight) Push(buf *conn.ConnBuffer) {
	buf.PushUAS(p.Light, true)
}

func CalculateNotAirBlocks(bitsPerBlock uint, airIndex int, data []int64) int {
	nonAirCount := 0
	blocksPerLong := uint(64 / bitsPerBlock)
	mask := (1 << bitsPerBlock) - 1 // Маска для извлечения битов

	for _, block := range data {
		for i := uint(0); i < blocksPerLong; i++ {
			blockIndex := int((block >> (i * bitsPerBlock)) & int64(mask))
			if blockIndex != airIndex {
				nonAirCount++
			}
		}
	}
	return nonAirCount
}
func CalculateNotAirBlocksNew(bitsPerBlock uint, airIndex int, data []int64) int {
	nonAirCount := 0
	blocksPerLong := uint(64 / bitsPerBlock)
	mask := (1 << bitsPerBlock) - 1          // Маска для извлечения битов
	remainderBits := uint(64 % bitsPerBlock) // Остаточные биты

	for _, block := range data {
		// Обработка полных блоков
		for i := uint(0); i < blocksPerLong; i++ {
			blockIndex := int((block >> (i * bitsPerBlock)) & int64(mask))
			if blockIndex != airIndex {
				nonAirCount++
			}
		}

		// Обработка остаточных битов, если они есть
		if remainderBits > 0 {
			// Сдвиг для получения остаточных битов
			lastBlockIndex := int((block >> (blocksPerLong * bitsPerBlock)) & int64((1<<remainderBits)-1))
			// Проверяем, не является ли последний неполный блок airIndex
			// Учитываем, что неполный блок может быть меньше bitsPerBlock
			if lastBlockIndex != airIndex {
				nonAirCount++
			}
		}
	}
	return nonAirCount
}
