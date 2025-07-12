package mode

import (
	"github.com/golangmc/minecraft-server/impl/conn"
)

func init() {
}

// func ParseChunk(data []byte) {
// 	buf := conn.ConnBuffer{}
// 	buf.PushUAS(data, false)
// 	x := buf.PullI32()
// 	z := buf.PullI32()
// 	fmt.Println("x", x, "z", z)

// 	bytesData := data[buf.InI():buf.InO()]
// 	dec := nbt.NewDecoder(bytes.NewReader(bytesData))
// 	dec.NetworkFormat(true)
// 	root := Root{
// 		Heightmap: Heightmap{},
// 	}
// 	// // считывает 1 nbt запись
// 	strerr, err := dec.Decode(&root.Heightmap)
// 	if err != nil {
// 		fmt.Println("err heightmap", err, strerr)
// 		os.Exit(1)
// 	}
// 	fmt.Println("root.Heightmap", root.Heightmap)

// 	tmpBuf := bytes.NewBuffer(make([]byte, 0))
// 	enc := nbt.NewEncoder(tmpBuf)
// 	enc.NetworkFormat(true)
// 	err = enc.Encode(root.Heightmap, "")
// 	if err != nil {
// 		fmt.Println("err ", err)
// 		os.Exit(1)
// 	}
// 	buf.SkpLen(int32(tmpBuf.Len()))
// 	chdatalen := buf.PullVrI()
// 	fmt.Println("chunk data len", chdatalen)
// 	startValue := buf.InI()
// 	counter := 0
// 	for {
// 		counter++
// 		if counter >= 25 {
// 			break
// 		}
// 		fmt.Println("counter", counter)
// 		temp3 := buf.PullU16()
// 		fmt.Println("Block count", temp3)
// 		{ // pallete

// 			temp4 := buf.PullByt()
// 			fmt.Println("Bits Per Entry", temp4)
// 			pal := Palette{}
// 			if temp4 == 0 {
// 				palVal := buf.PullVrI()
// 				fmt.Println("palVal", palVal)
// 				// pal.Blocks = []byte{buf.PullByt()}
// 				// fmt.Println("pal", pal, pal.Print())
// 			} else if temp4 >= 4 && temp4 <= 8 {
// 				pal.Pull(&buf)
// 				fmt.Println("pal", pal, pal.Print())
// 			} else {
// 				fmt.Println("temp4", temp4, buf.PullByt(), buf.PullByt(), buf.PullByt(), buf.PullByt())
// 				panic("sec temp4 not 0 or 4-8: " + strconv.Itoa(int(temp4)))
// 			}
// 			longsCount := buf.PullVrI()
// 			fmt.Println("лонгов данных", longsCount)
// 			for i := 0; i < int(longsCount); i++ {
// 				blocks := buf.PullI64()
// 				// firts 4 bits
// 				binaryStr := fmt.Sprintf("%064b", blocks)

// 				for i := 0; i < 64; i += 4 {
// 					group := binaryStr[i : i+4]
// 					// Преобразуем группу в число
// 					val, _ := strconv.ParseInt(group, 2, 64)
// 					val += 0
// 					// Выводим группу и её значение
// 					fmt.Printf("%s(%s) \t", group, registry.Blocks[int(pal.Blocks[int(val)])])
// 				}
// 				fmt.Println()
// 			}

// 		}
// 		{ // pallete

// 			temp4 := buf.PullByt()
// 			fmt.Println("Bits Per Entry", temp4)
// 			pal := Palette{}
// 			if temp4 == 0 {
// 				palVal := buf.PullVrI()
// 				fmt.Println("palVal", palVal)
// 				// pal.Blocks = []byte{buf.PullByt()}
// 				// fmt.Println("pal", pal, pal.Print())
// 			} else if temp4 >= 4 && temp4 <= 8 {
// 				pal.Pull(&buf)
// 				fmt.Println("pal", pal, pal.Print())
// 			} else {
// 				fmt.Println("temp4", temp4, buf.PullByt(), buf.PullByt(), buf.PullByt(), buf.PullByt())
// 				panic("temp4 not 0 or 4-8: " + strconv.Itoa(int(temp4)))
// 			}
// 			longsCount := buf.PullVrI()
// 			fmt.Println("лонгов данных", longsCount)
// 			for i := 0; i < int(longsCount); i++ {
// 				blocks := buf.PullI64()
// 				// firts 4 bits
// 				binaryStr := fmt.Sprintf("%064b", blocks)

// 				for i := 0; i < 64; i += 4 {
// 					group := binaryStr[i : i+4]
// 					// Преобразуем группу в число
// 					val, _ := strconv.ParseInt(group, 2, 64)
// 					val += 0
// 					// Выводим группу и её значение
// 					fmt.Printf("%s(%s) \t", group, registry.Blocks[int(pal.Blocks[int(val)])])
// 				}
// 				fmt.Println()
// 			}

// 		}
// 		// os.Exit(0)
// 	}
// 	buf.PullByt()
// 	endValue := buf.InI()
// 	fmt.Println("endValue", endValue)
// 	fmt.Println("endValue - startValue", endValue-startValue)

// 	// light data
// 	SkyLightLen := buf.PullVrI()
// 	fmt.Println("SkyLightLen", SkyLightLen)
// 	for i := 0; i < int(SkyLightLen); i++ {
// 		skyLight := buf.PullI64()
// 		fmt.Println("skyLightMask", skyLight)
// 	}
// 	BlockLightLen := buf.PullVrI()
// 	fmt.Println("BlockLightLen", BlockLightLen)
// 	for i := 0; i < int(BlockLightLen); i++ {
// 		blockLight := buf.PullI64()
// 		fmt.Println("blockLightMask", blockLight)
// 	}
// 	EmptySkyLightLen := buf.PullVrI()
// 	for i := 0; i < int(EmptySkyLightLen); i++ {
// 		emptySkyLight := buf.PullI64()
// 		fmt.Println("emptySkyLight", emptySkyLight)
// 	}
// 	EmptyBlockLightLen := buf.PullVrI()
// 	for i := 0; i < int(EmptyBlockLightLen); i++ {
// 		emptyBlockLight := buf.PullI64()
// 		fmt.Println("emptyBlockLight", emptyBlockLight)
// 	}
// 	SkyLightArrays := [][]byte{}
// 	// BlockLightArrays := [][]byte{}
// 	SkyLightArraysCount := buf.PullVrI()
// 	for i := 0; i < int(SkyLightArraysCount); i++ {
// 		fmt.Println("SkyLightArraysCount", i)
// 		arr := []byte{}
// 		len := buf.PullVrI()
// 		for j := 0; j < int(len); j++ {
// 			arr = append(arr, buf.PullByt())
// 		}
// 		SkyLightArrays = append(SkyLightArrays, arr)
// 	}
// 	fmt.Println("SkyLightArraysCount", SkyLightArraysCount)
// 	// fmt.Println("SkyLightArrays", SkyLightArrays)
// 	BlockLightArraysCount := buf.PullVrI()
// 	for i := 0; i < int(BlockLightArraysCount); i++ {
// 		// BlockLightArrays = append(BlockLightArrays, buf.PullUAS(int(temp)))
// 	}
// 	fmt.Println("BlockLightArraysCount", BlockLightArraysCount)

// 	fmt.Println(buf.InO())

// 	os.Exit(0)

// }

// func _CreateChunk(x, z int) []byte {
// 	buf := conn.ConnBuffer{}
// 	buf.PushI32(int32(x))
// 	buf.PushI32(int32(z))

// 	baseHM := []int64{}
// 	for range 37 {
// 		baseHM = append(baseHM, 0)
// 	}

// 	root := Root{
// 		Heightmap: Heightmap{
// 			MotionBlocking: baseHM,
// 			WorldSurface:   baseHM,
// 		},
// 	}
// 	var nbtBuf bytes.Buffer
// 	enc := nbt.NewEncoder(&nbtBuf)
// 	enc.NetworkFormat(true)
// 	if err := enc.Encode(root.Heightmap, ""); err != nil {
// 		panic(err)
// 	}
// 	buf.PushUAS(nbtBuf.Bytes(), false)
// 	// buf.PushVrI(2245)

// 	chunkTempBuf := conn.ConnBuffer{}
// 	// buf.PushByt(1)
// 	// buf.PushByt(0)
// 	chunkTempBuf.PushI16(256)
// 	chunkTempBuf.PushByt(4)
// 	pal := Palette{
// 		Blocks: []int{0, 9, 14, 1},
// 	}

// 	pal.Push(&chunkTempBuf)
// 	chunk := NewChunkSection()
// 	// chunk.Blocks[0][0][0] = 0x01
// 	var air = byte(0b0000)
// 	var grass = byte(0b0001)
// 	var cobble = byte(0b0010)
// 	var stone = byte(0b0011)
// 	for i := 0; i < 16; i++ {
// 		for j := 0; j < 16; j++ {
// 			chunk.SetBlock(i, 2, j, grass)
// 			if i == 0 || i == 15 || j == 0 || j == 15 {
// 				chunk.SetBlock(i, 2, j, cobble)
// 			}
// 		}
// 	}
// 	chunk.SetBlock(2, 2, 5, air)
// 	chunk.SetBlock(2, 2, 4, stone)

// 	chunk.Push(&chunkTempBuf)
// 	chunkTempBuf.PushByt(0)
// 	// buf.PushByt(57) // биом the void
// 	chunkTempBuf.PushByt(byte(rand.Intn(10) + 50)) // биом the void
// 	pal = Palette{
// 		Blocks: []int{57},
// 	}
// 	pal.Push(&chunkTempBuf)
// 	chunkTempBuf.PushUAS(randomBiomeData, false)

// 	buf.PushVrI(int32(len(chunkTempBuf.UAS())))
// 	buf.PushUAS(chunkTempBuf.UAS(), false)

// 	buf.PushByt(0) // block entities Data

// 	// buf.PushUAS(skyLightData, false)

// 	buf.PushByt(1) // sky light Mask
// 	buf.PushI64(7)

// 	buf.PushByt(0) // block light Mask
// 	buf.PushByt(0) // empty sky light
// 	buf.PushByt(1) // empty block light
// 	buf.PushI64(7)

// 	buf.PushByt(3) // Sky Light arrays length

// 	skyLight := NewSkyLight()
// 	skyLight.Push(&buf)
// 	skyLight.Push(&buf)
// 	skyLight.Push(&buf)
// 	buf.PushByt(0)
// 	return buf.UAS()
// }

var skyLightData = []byte{

	// 0x01, 0x00, 0x00,
	// 0x00, 0x00, 0x00, 0x00,
	// 0x00, 0x07, 0x00, 0x00,
	0x01, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x07, 0x03,
}
var randomBiomeData = []byte{
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x39, 0x00,
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
