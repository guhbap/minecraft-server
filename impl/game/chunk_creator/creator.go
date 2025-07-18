package chunkcreator

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/Tnze/go-mc/nbt"
	"github.com/aquilax/go-perlin"
	"github.com/golangmc/minecraft-server/impl/conn"
	"github.com/golangmc/minecraft-server/impl/game/chunk_utils"
)

var (
	EmptySection *SendingSection
)

type PerlinSetting struct {
	Alpha          float64
	Beta           float64
	N              int32
	Scale          float64
	YMin           int
	YMax           int
	Seed           int64
	PerlinInstance *perlin.Perlin
}

func NewPerlinSetting(seed int64) *PerlinSetting {
	return &PerlinSetting{
		Seed:  seed,
		Alpha: 1.5,
		Beta:  3.0,
		N:     3,
		Scale: 0.001,
		YMin:  50,
		YMax:  200,
	}
}

func (s *PerlinSetting) getPerlinY(x, z int) int {
	if s.PerlinInstance == nil {
		s.PerlinInstance = perlin.NewPerlin(s.Alpha, s.Beta, s.N, s.Seed)
	}
	noise := s.PerlinInstance.Noise2D(float64(x)*s.Scale, float64(z)*s.Scale)
	return normalize(noise, float64(s.YMin), float64(s.YMax))
}
func normalize(value, min, max float64) int {
	// Приводим значение к диапазону [0, 1]
	norm := (value + 1) / 2
	// Масштабируем в нужный диапазон
	return int(norm*float64(max-min)) + int(min)
}

func init() {
	return
	// var alpha float64 = 1.5
	// var beta float64 = 2.0
	// GlobalPerlin = perlin.NewPerlin(alpha, beta, 3, 123)

	longCount := len(EmptySection.Blocks)
	fmt.Println("longCount", longCount)
	// for _, long := range EmptySection.Blocks {
	// 	fmt.Println("long", long)
	// }
	fmt.Println("longCount", longCount)
	fmt.Println("BitsPerBlock", EmptySection.Pallete.BitsPerBlock)
	fmt.Println("Blocks", len(EmptySection.Pallete.Blocks))
}

func CreatePallete() *chunk_utils.Pallete {
	pallete := chunk_utils.NewPallete(
		[]chunk_utils.ChunkSectionBlocksPaletteNbt{
			{
				Name: "minecraft:air",
			},
			{
				Name: "minecraft:stone",
			},
			{
				Name: "minecraft:dirt",
			},
			{
				Name: "minecraft:red_sand",
			},
			{
				Name: "minecraft:gravel",
			},
			{
				Name: "minecraft:sand",
			},
			{
				Name: "minecraft:red_sandstone",
			},
			{
				Name: "minecraft:redstone_block",
			},
			{
				Name: "minecraft:coal_ore",
			},
			{
				Name: "minecraft:iron_ore",
			},
			{
				Name: "minecraft:gold_ore",
			},
			{
				Name: "minecraft:diamond_ore",
			},
			{
				Name: "minecraft:lapis_ore",
			},
			{
				Name: "minecraft:emerald_ore",
			},
			{
				Name: "minecraft:nether_quartz_ore",
			},
			// {
			// 	Name: "minecraft:grass_block",
			// 	Properties: map[string]string{
			// 		"snowy": "true",
			// 	},
			// },
			// {
			// 	Name: "minecraft:grass_block",
			// 	Properties: map[string]string{
			// 		"snowy": "false",
			// 	},
			// },
		},
	)
	return pallete
}

type SendingChunk struct {
	Sections []*SendingSection
	X        int
	Z        int
}

func (c *SendingChunk) GeneratePerlin(setting *PerlinSetting) {
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			height := setting.getPerlinY(c.X*16+x, c.Z*16+z)
			for y := height; y >= -60; y-- {
				c.SetBlock(x, y, z, 1)
			}
		}
	}
}

func (c *SendingChunk) SetBlock(x, y, z int, blockIndex byte) {
	if y < -64 {
		panic("y < -64")
	}
	if y > 319 {
		panic("y > 319")
	}

	sectionIndex := (y + 64) / 16 // от 0 до 23
	localY := y % 16
	if localY < 0 {
		localY += 16
	}

	if sectionIndex < 0 || sectionIndex >= len(c.Sections) {
		panic(fmt.Sprintf("Invalid section index: %d", sectionIndex))
	}

	section := c.Sections[sectionIndex]
	section.SetBlocks(x, localY, z, blockIndex)
	gb := section.GetBlock(x, localY, z)
	if gb != blockIndex {
		panic("block not set at " + strconv.Itoa(x) + " " + strconv.Itoa(y) + " " + strconv.Itoa(z))
	}
}

func (c *SendingChunk) Push(buf *conn.ConnBuffer) {
	buf.PushI32(int32(c.X))
	buf.PushI32(int32(c.Z))

	{ // Heightmaps
		nbtBuf := bytes.NewBuffer(make([]byte, 0))
		enc := nbt.NewEncoder(nbtBuf)
		enc.NetworkFormat(true)
		err := enc.Encode(&chunk_utils.Heightmap{}, "")
		if err != nil {
			panic(err)
		}
		buf.PushUAS(nbtBuf.Bytes(), false)
	}

	tmpbuf := conn.ConnBuffer{}
	for _, section := range c.Sections {
		section.Push(&tmpbuf)
	}
	buf.PushUAS(tmpbuf.UAS(), true)
	buf.PushByt(0)

	buf.PushByt(0) // Empty Sky Light Mask

	buf.PushByt(0) // Empty Block Light Mask

	buf.PushByt(0) // Empty Sky Light Mask

	buf.PushByt(0) // Empty Block Light Mask

	buf.PushVrI(0) // Sky Light arrays length

	buf.PushVrI(0) // Block Light arrays length

}

type SendingSection struct {
	Y       int
	Blocks  []int64
	Pallete *chunk_utils.Pallete
}

func (section *SendingSection) SetBlocks(x, y, z int, blockIndex byte) {
	if y < 0 {
		panic("y < 0")
	}
	if y > 15 {
		panic("y > 15")
	}
	if z < 0 {
		panic("z < 0")
	}
	if z > 15 {
		panic("z > 15")
	}
	if x < 0 {
		panic("x < 0")
	}
	if x > 15 {
		panic("x > 15")
	}

	x = 15 - x
	// Проверяем валидность координат
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return // или можно выбросить ошибку, если требуется
	}

	// Рассчитываем индекс блока в одномерном массиве
	blockIdx := x + z*16 + y*16*16

	// Параметры палитры
	bitsPerBlock := int(section.Pallete.BitsPerBlock)
	blocksPerLong := 64 / bitsPerBlock

	// Индекс long в массиве Blocks
	longIdx := blockIdx / blocksPerLong
	// Позиция блока внутри long (в блоках)
	blockOffset := blockIdx % blocksPerLong
	// Начальная позиция битов для блока
	bitPos := blockOffset * bitsPerBlock

	// Маска для очистки старых бит
	mask := int64((1<<bitsPerBlock - 1) << (64 - bitsPerBlock - bitPos))
	// Очищаем биты для текущего блока
	section.Blocks[longIdx] &^= mask
	// Устанавливаем новые биты
	section.Blocks[longIdx] |= int64(blockIndex) << (64 - bitsPerBlock - bitPos)
}
func (section *SendingSection) GetBlock(x, y, z int) byte {
	x = 15 - x

	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return 0
	}

	blockIdx := x + z*16 + y*16*16

	bitsPerBlock := int(section.Pallete.BitsPerBlock)
	if bitsPerBlock <= 0 || bitsPerBlock > 8 {
		return 0
	}
	blocksPerLong := 64 / bitsPerBlock

	longIdx := blockIdx / blocksPerLong
	blockOffset := blockIdx % blocksPerLong
	bitPos := blockOffset * bitsPerBlock
	shift := 64 - bitsPerBlock - bitPos

	mask := int64((1 << bitsPerBlock) - 1)
	value := (section.Blocks[longIdx] >> shift) & mask

	return byte(value)
}
func (section *SendingSection) Optimize() {
	uniqueIndexes := extractUniqueBlocks(section.Blocks, int(section.Pallete.BitsPerBlock))
	if len(uniqueIndexes) != 1 {
		return
	}

	var onlyIndex uint8
	for idx := range uniqueIndexes {
		onlyIndex = idx
		break
	}

	block := section.Pallete.Blocks[onlyIndex]

	section.Pallete = &chunk_utils.Pallete{
		Blocks:       []int{block},
		BitsPerBlock: 0,
		AirIndex: func() int {
			if block == 0 {
				return 0
			}
			return -1
		}(),
	}
	section.Blocks = nil
}

func extractUniqueBlocks(bitArray []int64, bitsPerBlock int) map[uint8]bool {
	result := make(map[uint8]bool)

	if bitsPerBlock < 0 || bitsPerBlock > 8 {
		panic("bitsPerBlock должен быть от 1 до 8")
	}

	if bitsPerBlock == 0 {
		result[uint8(bitArray[0])] = true
		return result
	}

	mask := int64((1 << bitsPerBlock) - 1)
	totalBits := int64(len(bitArray)) * 64

	for bitIndex := int64(0); bitIndex+int64(bitsPerBlock) <= totalBits; bitIndex += int64(bitsPerBlock) {
		longIndex := bitIndex / 64
		startBit := bitIndex % 64

		var value int64
		if startBit+int64(bitsPerBlock) <= 64 {
			value = (bitArray[longIndex] >> startBit) & mask
		} else {
			low := bitArray[longIndex] >> startBit
			high := bitArray[longIndex+1] << (64 - startBit)
			value = (low | high) & mask
		}

		result[uint8(value)] = true

		if len(result) > 1 {
			break // ранний выход — дальше нет смысла
		}
	}

	return result
}

func (section *SendingSection) Push(buf *conn.ConnBuffer) {
	section.Optimize()
	tempBuf := conn.NewBuffer()
	pallete := section.Pallete
	pallete.Push(tempBuf, false)
	nonAirCount := 0
	bitsPerBlock := uint(pallete.BitsPerBlock)
	if bitsPerBlock == 0 {
		bitsPerBlock = 1
	}
	blocksPerLong := uint(64 / bitsPerBlock)
	mask := (1 << bitsPerBlock) - 1 // Маска для извлечения битов
	airIndex := pallete.AirIndex

	if airIndex == -1 {
		nonAirCount = CHUNK_TOTAL_BLOCKS
	} else {

		tempBuf.PushVrI(int32(len(section.Blocks)))
		for _, block := range section.Blocks {
			tempBuf.PushI64(block)
			for i := uint(0); i < blocksPerLong; i++ {
				blockIndex := int((block >> (i * bitsPerBlock)) & int64(mask))
				if blockIndex != airIndex {
					nonAirCount++
				}
			}
		}
	}
	pallete = chunk_utils.NewBiomesPallete([]string{"minecraft:plains"})
	pallete.Push(tempBuf, true)
	tempBuf.PushVrI(int32(1))
	tempBuf.PushI64(0)
	buf.PushI16(int16(nonAirCount))
	buf.PushUAS(tempBuf.UAS(), false)
}

const (
	CHUNK_TOTAL_BLOCKS = int(4096)
)

func CreateEmptySection(pallete *chunk_utils.Pallete) *SendingSection {
	section := SendingSection{
		Y:       0,
		Pallete: pallete,
	}

	bitsPerBlock := pallete.BitsPerBlock
	totalBits := CHUNK_TOTAL_BLOCKS * int(bitsPerBlock)
	bufferSize := (totalBits + 7) / 8 // округляем вверх, если не кратно 8
	bytesBuffer := make([]byte, bufferSize)

	airIndex := uint64(pallete.AirIndex)

	// позиция в битах
	bitPos := 0
	for i := 0; i < CHUNK_TOTAL_BLOCKS; i++ {
		byteIndex := bitPos / 8
		bitOffset := bitPos % 8

		// Пишем blockBits в нужную позицию буфера
		for b := 0; b < int(bitsPerBlock); b++ {
			// Вычисляем позицию бита в индексе
			bit := (airIndex >> (int(bitsPerBlock) - 1 - int(b))) & 1

			// Позиция в буфере
			bufBitPos := bitOffset + b
			targetByteIndex := byteIndex + bufBitPos/8
			targetBitOffset := 7 - (bufBitPos % 8)

			if bit == 1 {
				bytesBuffer[targetByteIndex] |= 1 << targetBitOffset
			}
		}

		bitPos += int(bitsPerBlock)
	}

	// Рассчитываем количество long (64-битных значений)
	blocksPerLong := 64 / int(bitsPerBlock)
	numLongs := (CHUNK_TOTAL_BLOCKS + blocksPerLong - 1) / blocksPerLong
	longs := make([]int64, numLongs)

	// Заполняем longs из bytesBuffer
	for i := 0; i < len(longs); i++ {
		for b := 0; b < 8 && (i*8+b) < len(bytesBuffer); b++ {
			longs[i] |= int64(bytesBuffer[i*8+b]) << (56 - b*8)
		}
	}

	section.Blocks = longs

	return &section
}
