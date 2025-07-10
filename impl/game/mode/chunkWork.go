package mode

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/Tnze/go-mc/nbt"
	"github.com/golangmc/minecraft-server/impl/conn"
	"github.com/golangmc/minecraft-server/impl/game/registry"
)

type Heightmap struct {
	MotionBlocking []int64 `nbt:"MOTION_BLOCKING"`
	WorldSurface   []int64 `nbt:"WORLD_SURFACE"`
}
type Root struct {
	Heightmap Heightmap `nbt:""`
}

func init() {
	// datas := LoadChunk(0, 0)
	// ParseChunk(datas)
	return

	// filename := "testData/chunk-stone-and-cobble-0-0.hex"
	filename := "testData/chunk-stone-cobble-grass-0-0.hex"
	// filename := "testData/chunk-air-only--1-0.hex"

	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	strData := string(data)

	bytesData := []byte{}

	hexes := strings.Split(strData, " ") // one hex is string like 0xa1
	for _, hexStr := range hexes {
		cleanHex := strings.TrimPrefix(hexStr, "0x")
		// Разбираем строку в число
		value, err := strconv.ParseUint(cleanHex, 16, 8)
		if err != nil {
			fmt.Println("Ошибка при парсинге:", err)
			continue
		}
		bytesData = append(bytesData, byte(value))
	}
	// os.WriteFile("testData/chunk-stone-and-cobble-0-0.bin", bytesData, 0644)
	ParseChunk(bytesData)
	// ParseChunk(CreateChunk(0, 0))
}
func CreateChunk(nbtData ChunkNbt) []byte {
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
		// data len будет после цикла
		{ // Data
			// Chunk Section
			// Block count
			// Block states
			// Biomes

			chSectionBuf := conn.ConnBuffer{}
			fmt.Println("nbtData.Sections", len(nbtData.Sections[:1]))
			for _, section := range nbtData.Sections[:1] {
				break
				chSectionBuf.PushI16(1000) // todo реальное количество невоздушных блоков

				{ // Block states

					pal := section.BlockStates.Palette
					// fmt.Println("section", section.Y)
					// fmt.Println("pal", pal)
					pallete := Palette{
						Blocks: []int{},
					}
					for _, palBlock := range pal {
						pallete.Blocks = append(pallete.Blocks, int(registry.InvertedBlocks[palBlock.Name]))
					}
					fmt.Println("pallete", pallete)
					pallete.Push(&chSectionBuf)
					longsCount := len(section.BlockStates.Data)
					chSectionBuf.PushVrI(int32(longsCount))
					for _, block := range section.BlockStates.Data {
						chSectionBuf.PushI64(block)
					}
				}
				{ // Biomes
					bpal := section.Biomes.Palette
					fmt.Println("bpal", bpal)
					bpallete := Palette{
						Blocks: []int{},
					}

					for idx, _ := range bpal {
						bpallete.Blocks = append(bpallete.Blocks, int(idx)) // todo тут временный костыль из-за того что нет четкого регистра биомов
					}
					bpallete.Push(&chSectionBuf)
					biomesCount := len(section.Biomes.Data)
					chSectionBuf.PushVrI(int32(biomesCount))
					for _, biome := range section.Biomes.Data {
						chSectionBuf.PushI64(biome)
					}
				}

				// break // todo пока что для ттеста отправляю только одну секцию
				// biomesBuf := conn.ConnBuffer{}
				// biomesBuf.PushVrI(int32(len(bpal)))
				// for _, biome := range bpal {
				// 	biomesBuf.PushString(biome)
				// }
			}

			buf.PushUAS(chSectionBuf.UAS(), true)

		}
		{ // BlockEntities
			buf.PushByt(0)

		}
	}
	{ // light
		// chSectionBuf.PushUAS(section.BlockLight, false)
		// chSectionBuf.PushUAS(section.SkyLight, false)
		skyLight := NewSkyLight()
		buf.PushByt(3)
		skyLight.Push(&buf)
		skyLight.Push(&buf)
		skyLight.Push(&buf)
		buf.PushByt(0)

	}

	// return buf.UAS()

	return buf.UAS()
	// buf.PushString(nbtData.Status)
	// buf.PushI64(nbtData.LastUpdate)
	// buf.PushUAS(nbtData.Sections, false)
}

func NbtTest() {
	root := ChunkNbt{
		Heightmaps: Heightmap{
			MotionBlocking: []int64{
				2526915989746289285,
				2400603680826791563,
				2508901178114640517,
				2400604093949221004,
				2400603680825215621,
				2508936844464429195,
				2400603680825215621,
				2526951242839427717,
				2400603680825218699,
				2508901590431500933,
				2400603681632097932,
				2400603680825215621,
				2418653676830529163,
				2400603680825215621,
				2418653332560743557,
				2400603680825216134,
				2418653332560743557,
				2400603749679172742,
				2418653332560743045,
				2418653332560743558,
				2418653332426263174,
				2418653332560743558,
				2418618079334960262,
				2418653332560743558,
				2400603749679172742,
				2418653332560743557,
				2400638934051261574,
				2418653332560480901,
				2418653332560743558,
				2418653263706786438,
				2418653332560743558,
				2400603680825478278,
				2418653332560743558,
				2400603749679172742,
				2418653332560743045,
				2418653332560743558,
				17885891205},
			WorldSurface: []int64{2526915989746551941,
				2418653263842580107,
				2508901246834117254,
				2400604093949221004,
				2418618148054174342,
				2508936844464429195,
				2418618148188654214,
				2526951242839428230,
				2418618079334962827,
				2508901590565981317,
				2418618148861056652,
				2418653263841266309,
				2436703259712099979,
				2400603680959433349,
				2436702915576532613,
				2400638865197304966,
				2436702915442576517,
				2400603749679435399,
				2436667731070487174,
				2418688516932832391,
				2436702984161528455,
				2418653332560744071,
				2436632546564181127,
				2436667731070487686,
				2400603818398649479,
				2436667731070226054,
				2400638934185479303,
				2418653401414175365,
				2418688516933095046,
				2418653332560480903,
				2436702984296270983,
				2400603749545217158,
				2418653401280482438,
				2418653332560743559,
				2418653332560743558,
				2418653401280220294,
				17886153350}},
	}
	nbtBuf := bytes.NewBuffer(make([]byte, 0))
	enc := nbt.NewEncoder(nbtBuf)
	enc.NetworkFormat(true)
	err := enc.Encode(&root.Heightmaps, "")
	if err != nil {
		fmt.Println("err", err)
		os.Exit(1)
	}
	bytesData := nbtBuf.Bytes()
	dec := nbt.NewDecoder(bytes.NewReader(bytesData))
	dec.NetworkFormat(true)
	root2 := Root{
		Heightmap: Heightmap{},
	}
	dec.Decode(&root2.Heightmap)
	fmt.Println("root2", root2)

	os.Exit(0)

}
func ParseChunk(data []byte) {
	buf := conn.ConnBuffer{}
	buf.PushUAS(data, false)
	x := buf.PullI32()
	z := buf.PullI32()
	fmt.Println("x", x, "z", z)

	bytesData := data[buf.InI():buf.InO()]
	dec := nbt.NewDecoder(bytes.NewReader(bytesData))
	dec.NetworkFormat(true)
	root := Root{
		Heightmap: Heightmap{},
	}
	// // считывает 1 nbt запись
	strerr, err := dec.Decode(&root.Heightmap)
	if err != nil {
		fmt.Println("err heightmap", err, strerr)
		os.Exit(1)
	}
	fmt.Println("root.Heightmap", root.Heightmap)

	tmpBuf := bytes.NewBuffer(make([]byte, 0))
	enc := nbt.NewEncoder(tmpBuf)
	enc.NetworkFormat(true)
	err = enc.Encode(root.Heightmap, "")
	if err != nil {
		fmt.Println("err ", err)
		os.Exit(1)
	}
	buf.SkpLen(int32(tmpBuf.Len()))
	chdatalen := buf.PullVrI()
	fmt.Println("chunk data len", chdatalen)
	startValue := buf.InI()
	counter := 0
	for {
		counter++
		if counter >= 25 {
			break
		}
		fmt.Println("counter", counter)
		temp3 := buf.PullU16()
		fmt.Println("Block count", temp3)
		{ // pallete

			temp4 := buf.PullByt()
			fmt.Println("Bits Per Entry", temp4)
			pal := Palette{}
			if temp4 == 0 {
				palVal := buf.PullVrI()
				fmt.Println("palVal", palVal)
				// pal.Blocks = []byte{buf.PullByt()}
				// fmt.Println("pal", pal, pal.Print())
			} else if temp4 >= 4 && temp4 <= 8 {
				pal.Pull(&buf)
				fmt.Println("pal", pal, pal.Print())
			} else {
				fmt.Println("temp4", temp4, buf.PullByt(), buf.PullByt(), buf.PullByt(), buf.PullByt())
				panic("sec temp4 not 0 or 4-8: " + strconv.Itoa(int(temp4)))
			}
			longsCount := buf.PullVrI()
			fmt.Println("лонгов данных", longsCount)
			for i := 0; i < int(longsCount); i++ {
				blocks := buf.PullI64()
				// firts 4 bits
				binaryStr := fmt.Sprintf("%064b", blocks)

				for i := 0; i < 64; i += 4 {
					group := binaryStr[i : i+4]
					// Преобразуем группу в число
					val, _ := strconv.ParseInt(group, 2, 64)
					val += 0
					// Выводим группу и её значение
					fmt.Printf("%s(%s) \t", group, registry.Blocks[int(pal.Blocks[int(val)])])
				}
				fmt.Println()
			}

		}
		{ // pallete

			temp4 := buf.PullByt()
			fmt.Println("Bits Per Entry", temp4)
			pal := Palette{}
			if temp4 == 0 {
				palVal := buf.PullVrI()
				fmt.Println("palVal", palVal)
				// pal.Blocks = []byte{buf.PullByt()}
				// fmt.Println("pal", pal, pal.Print())
			} else if temp4 >= 4 && temp4 <= 8 {
				pal.Pull(&buf)
				fmt.Println("pal", pal, pal.Print())
			} else {
				fmt.Println("temp4", temp4, buf.PullByt(), buf.PullByt(), buf.PullByt(), buf.PullByt())
				panic("temp4 not 0 or 4-8: " + strconv.Itoa(int(temp4)))
			}
			longsCount := buf.PullVrI()
			fmt.Println("лонгов данных", longsCount)
			for i := 0; i < int(longsCount); i++ {
				blocks := buf.PullI64()
				// firts 4 bits
				binaryStr := fmt.Sprintf("%064b", blocks)

				for i := 0; i < 64; i += 4 {
					group := binaryStr[i : i+4]
					// Преобразуем группу в число
					val, _ := strconv.ParseInt(group, 2, 64)
					val += 0
					// Выводим группу и её значение
					fmt.Printf("%s(%s) \t", group, registry.Blocks[int(pal.Blocks[int(val)])])
				}
				fmt.Println()
			}

		}
		// os.Exit(0)
	}
	buf.PullByt()
	endValue := buf.InI()
	fmt.Println("endValue", endValue)
	fmt.Println("endValue - startValue", endValue-startValue)

	// light data
	SkyLightLen := buf.PullVrI()
	fmt.Println("SkyLightLen", SkyLightLen)
	for i := 0; i < int(SkyLightLen); i++ {
		skyLight := buf.PullI64()
		fmt.Println("skyLightMask", skyLight)
	}
	BlockLightLen := buf.PullVrI()
	fmt.Println("BlockLightLen", BlockLightLen)
	for i := 0; i < int(BlockLightLen); i++ {
		blockLight := buf.PullI64()
		fmt.Println("blockLightMask", blockLight)
	}
	EmptySkyLightLen := buf.PullVrI()
	for i := 0; i < int(EmptySkyLightLen); i++ {
		emptySkyLight := buf.PullI64()
		fmt.Println("emptySkyLight", emptySkyLight)
	}
	EmptyBlockLightLen := buf.PullVrI()
	for i := 0; i < int(EmptyBlockLightLen); i++ {
		emptyBlockLight := buf.PullI64()
		fmt.Println("emptyBlockLight", emptyBlockLight)
	}
	SkyLightArrays := [][]byte{}
	// BlockLightArrays := [][]byte{}
	SkyLightArraysCount := buf.PullVrI()
	for i := 0; i < int(SkyLightArraysCount); i++ {
		fmt.Println("SkyLightArraysCount", i)
		arr := []byte{}
		len := buf.PullVrI()
		for j := 0; j < int(len); j++ {
			arr = append(arr, buf.PullByt())
		}
		SkyLightArrays = append(SkyLightArrays, arr)
	}
	fmt.Println("SkyLightArraysCount", SkyLightArraysCount)
	// fmt.Println("SkyLightArrays", SkyLightArrays)
	BlockLightArraysCount := buf.PullVrI()
	for i := 0; i < int(BlockLightArraysCount); i++ {
		// BlockLightArrays = append(BlockLightArrays, buf.PullUAS(int(temp)))
	}
	fmt.Println("BlockLightArraysCount", BlockLightArraysCount)

	fmt.Println(buf.InO())

	os.Exit(0)

}
func CalculateOffset(x, z int) int {
	xMod := x % 32
	if xMod < 0 {
		xMod += 32
	}

	zMod := z % 32
	if zMod < 0 {
		zMod += 32
	}

	return 4 * (xMod + zMod*32)
}
func LoadChunk(x, z int) []byte {
	x, z = 0, 0
	offset := CalculateOffset(x, z)
	region := fmt.Sprintf("world/region/r.%d.%d.mca", x>>5, z>>5)
	fmt.Println("region", region)
	datas, err := os.ReadFile(region)
	if err != nil {
		panic(err)
	}

	// headReader := bytes.NewReader(datas)
	// headReader := bufio.NewReader(bytes.NewReader(datas))
	headReader := conn.ConnBuffer{}
	headReader.PushUAS(datas[0:4*KB], false)
	headReader.SkpLen(int32(offset))
	chunkReader := conn.ConnBuffer{}
	chunkReader.PushUAS(datas[4*KB:], false)
	fmt.Println("--------------------------------")

	pos := headReader.PullI24()
	fmt.Println("pos", pos)
	fmt.Println("size", headReader.PullByt())
	si := int32(pos*4*KB - 4*KB)
	fmt.Println("si", si)
	chunkReader.SetIndex(si)

	chunkReader.SkpLen(int32(pos * 4 * KB))
	// rl1 := chunkReader.PullByt()
	// rl2 := chunkReader.PullByt()
	// rl3 := chunkReader.PullByt()
	// rl4 := chunkReader.PullByt()
	// fmt.Println("rl1", rl1, "rl2", rl2, "rl3", rl3, "rl4", rl4)
	realLen := chunkReader.PullI32()
	compressScheme := chunkReader.PullByt()
	compressedData := []byte{}
	for i := 0; i < int(realLen); i++ {
		compressedData = append(compressedData, chunkReader.PullByt())
	}
	if compressScheme == 2 {
		// zlib decompressing
		realData, err := zlib.NewReader(bytes.NewReader(compressedData))
		if err != nil {
			panic(err)
		}

		dec := nbt.NewDecoder(realData)
		root := ChunkNbt{}
		_, err = dec.Decode(&root)
		if err != nil {
			panic(err)
		}
		// fmt.Println("root", root)
		return CreateChunk(root)
	}
	return nil
	fmt.Println("realLen", realLen)
	fmt.Println("compressScheme", compressScheme)

	// bytesData := chunkReader.UAS()
	// fmt.Println("bytesData", bytesData)
	os.Exit(0)

	return datas
}

const B = 1
const KB = 1024 * B

func _CreateChunk(x, z int) []byte {
	buf := conn.ConnBuffer{}
	buf.PushI32(int32(x))
	buf.PushI32(int32(z))

	baseHM := []int64{}
	for range 37 {
		baseHM = append(baseHM, 0)
	}

	root := Root{
		Heightmap: Heightmap{
			MotionBlocking: baseHM,
			WorldSurface:   baseHM,
		},
	}
	var nbtBuf bytes.Buffer
	enc := nbt.NewEncoder(&nbtBuf)
	enc.NetworkFormat(true)
	if err := enc.Encode(root.Heightmap, ""); err != nil {
		panic(err)
	}
	buf.PushUAS(nbtBuf.Bytes(), false)
	buf.PushVrI(2245)
	// buf.PushByt(1)
	// buf.PushByt(0)
	buf.PushI16(256)
	buf.PushByt(4)
	pal := Palette{
		Blocks: []int{0, 9, 14, 1},
	}

	pal.Push(&buf)
	chunk := NewChunkSection()
	// chunk.Blocks[0][0][0] = 0x01
	var air = byte(0b0000)
	var grass = byte(0b0001)
	var cobble = byte(0b0010)
	var stone = byte(0b0011)
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			chunk.SetBlock(i, 2, j, grass)
			if i == 0 || i == 15 || j == 0 || j == 15 {
				chunk.SetBlock(i, 2, j, cobble)
			}
		}
	}
	chunk.SetBlock(2, 2, 5, air)
	chunk.SetBlock(2, 2, 4, stone)

	chunk.Push(&buf)
	buf.PushByt(0)
	// buf.PushByt(57) // биом the void
	buf.PushByt(byte(rand.Intn(10) + 50)) // биом the void
	pal = Palette{
		Blocks: []int{57},
	}
	pal.Push(&buf)
	buf.PushUAS(randomBiomeData, false)
	skyLight := NewSkyLight()
	buf.PushByt(3)
	skyLight.Push(&buf)
	skyLight.Push(&buf)
	skyLight.Push(&buf)
	buf.PushByt(0)
	return buf.UAS()
}

func deepCompareByteArrays(test, real []byte) (bool, int, string) {
	for i := 0; i < len(real); i++ {
		if test[i] != real[i] {
			return false, i, fmt.Sprintf("test[%d] = 0x%x, real[%d] = 0x%x (0x%x, 0x%x, 0x%x)", i, test[i], i, real[i], test[i], test[i+1], test[i+2])
		}
	}
	return true, 0, ""
}

type Palette struct {
	Blocks []int
}

func (p *Palette) Push(buf *conn.ConnBuffer) {

	// оставить только уникальные значения

	uniqueBlocks := []int{}
	for _, block := range p.Blocks {
		if !slices.Contains(uniqueBlocks, block) {
			uniqueBlocks = append(uniqueBlocks, block)
		}
	}

	if len(uniqueBlocks) > 1 {
		if len(uniqueBlocks) > 16 {
			panic("palette blocks count > 16")
		}
		buf.PushByt(4)
		for _, block := range p.Blocks {
			buf.PushVrI(int32(block))
		}
	} else {
		buf.PushByt(0)
		buf.PushVrI(int32(p.Blocks[0]))
	}
	// buf.PushUAS(p.Blocks, false)
}
func (p *Palette) Pull(buf *conn.ConnBuffer) {
	len := int(buf.PullVrI())
	for i := 0; i < len; i++ {
		p.Blocks = append(p.Blocks, int(buf.PullVrI()))
	}
}
func (p *Palette) Print() string {
	blocks := []string{}
	for _, block := range p.Blocks {
		blocks = append(blocks, registry.Blocks[int(block)])
	}
	return fmt.Sprintf("Palette: %v", blocks)
}

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

var randomBiomeData = []byte{
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x39, 0x00,
	// 0x00, 0x01, 0x00, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x07, 0x00, 0x00,
	// 0x01, 0x00, 0x00, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	0x07,
}

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
	BlockLight  []byte                     `nbt:"block_light"` // todo эти поля не используются, надо будет как то решать
	SkyLight    []byte                     `nbt:"sky_light"`   // todo эти поля не используются, надо будет как то решать
}

type ChunkSectionBlockStatesNbt struct {
	Palette []ChunkSectionBlocksPaletteNbt `nbt:"palette"`
	Data    []int64                        `nbt:"data"`
}

type ChunkSectionBlocksPaletteNbt struct {
	Name       string `nbt:"Name"`
	Properties struct {
		Name string `nbt:"Name"`
	} `nbt:"Properties"`
}
type ChunkSectionBiomesPaletteNbt struct {
	Name string `nbt:"Name"`
}

type ChunkSectionBiomesNbt struct {
	Palette []string `nbt:"palette"`
	Data    []int64  `nbt:"data"`
}
