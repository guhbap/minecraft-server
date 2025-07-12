package chunk_utils

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"

	"github.com/Tnze/go-mc/nbt"
	"github.com/golangmc/minecraft-server/impl/conn"
)

const B = 1
const KB = 1024 * B

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

func LoadChunk(x, z int) *ChunkNbt {
	offset := CalculateOffset(x, z)
	fmt.Println("offset", offset)
	region := fmt.Sprintf("world/region/r.%d.%d.mca", x>>5, z>>5)
	fmt.Println("region", region)
	datas, err := os.ReadFile(region)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("region not found")
			return nil
		}
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
	if pos == 0 {
		return nil
	}
	fmt.Println("size", headReader.PullByt())
	si := int32(pos*4*KB - 4*KB)
	fmt.Println("si", si)
	chunkReader.SetIndex(si)

	// chunkReader.SkpLen(int32(pos * 4 * KB))
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
		return &root
	}
	fmt.Println("realLen", realLen)
	fmt.Println("compressScheme", compressScheme)
	return nil

	// bytesData := chunkReader.UAS()
	// fmt.Println("bytesData", bytesData)
}

func deepCompareByteArrays(test, real []byte) (bool, int, string) {
	for i := 0; i < len(real); i++ {
		if test[i] != real[i] {
			return false, i, fmt.Sprintf("test[%d] = 0x%x, real[%d] = 0x%x (0x%x, 0x%x, 0x%x)", i, test[i], i, real[i], test[i], test[i+1], test[i+2])
		}
	}
	return true, 0, ""
}
