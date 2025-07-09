package conn

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"

	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/data"
	"github.com/golangmc/minecraft-server/apis/uuid"
)

/*
 Language used:
	- Len = Length
	- Arr = Array

	- Bit = Boolean
	- Byt = Byte
	- Int = int
	- VrI = VarInt
	- Srt = Short
	- Txt = String
*/

type ConnBuffer struct {
	iIndex int32
	oIndex int32

	bArray []byte
}

func (b *ConnBuffer) String() string {
	return fmt.Sprintf("Buffer[%d](i: %d, o: %d)%v", b.Len(), b.iIndex, b.oIndex, b.bArray)
}

func (b *ConnBuffer) HexString() string {
	return hex.EncodeToString(b.bArray)
}

// new
func NewBuffer() buff.Buffer {
	return NewBufferWith(make([]byte, 0))
}

func NewBufferWith(bArray []byte) buff.Buffer {
	return &ConnBuffer{bArray: bArray}
}

// server_data
func (b *ConnBuffer) Len() int32 {
	return int32(len(b.bArray))
}

func (b *ConnBuffer) SAS() []int8 {
	return asSArray(b.bArray)
}

func (b *ConnBuffer) UAS() []byte {
	return b.bArray
}

func (b *ConnBuffer) InI() int32 {
	return b.iIndex
}

func (b *ConnBuffer) InO() int32 {
	return b.oIndex
}

func (b *ConnBuffer) SkpAll() {
	b.SkpLen(b.Len() - 1)
}

func (b *ConnBuffer) SkpLen(delta int32) {
	b.iIndex += delta
}

// pull
func (b *ConnBuffer) PullBit() bool {
	return b.pullNext() != 0
}

func (b *ConnBuffer) PullByt() byte {
	return b.pullNext()
}

func (b *ConnBuffer) PullI16() int16 {
	return int16(binary.BigEndian.Uint16(b.pullSize(4)))
}

func (b *ConnBuffer) PullU16() uint16 {
	return uint16(b.pullNext())<<8 | uint16(b.pullNext())
}

func (b *ConnBuffer) PullI32() int32 {
	return int32(binary.BigEndian.Uint32(b.pullSize(4)))
}

func (b *ConnBuffer) PullI64() int64 {
	return int64(b.PullU64())
}

func (b *ConnBuffer) PullU64() uint64 {
	return binary.BigEndian.Uint64(b.pullSize(8))
}

func (b *ConnBuffer) PullF32() float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(b.pullSize(4)))
}

func (b *ConnBuffer) PullF64() float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(b.pullSize(8)))
}

func (b *ConnBuffer) PullVrI() int32 {
	return int32(b.pullVariable(5))
}

func (b *ConnBuffer) PullVrL() int64 {
	return b.pullVariable(10)
}

func (b *ConnBuffer) PullTxt() string {
	return string(b.PullUAS())
}

func (b *ConnBuffer) PullUAS() []byte {
	sze := b.PullVrI()
	arr := b.bArray[b.iIndex : b.iIndex+sze]

	b.iIndex += sze

	return arr
}

func (b *ConnBuffer) PullSAS() []int8 {
	return asSArray(b.PullUAS())
}

func (b *ConnBuffer) PullUID() uuid.UUID {
	data, _ := uuid.BitsToUUID(b.PullI64(), b.PullI64())

	return data
}

func (b *ConnBuffer) PullPos() data.PositionI {
	val := b.PullU64()

	x := int64(val) >> 38
	y := int64(val) & 0xFFF
	z := int64(val) << 26 >> 38

	return data.PositionI{
		X: x,
		Y: y,
		Z: z,
	}
}

// push
func (b *ConnBuffer) PushBit(data bool) {
	if data {
		b.pushNext(byte(0x01))
	} else {
		b.pushNext(byte(0x00))
	}
}

func (b *ConnBuffer) PushByt(data byte) {
	b.pushNext(data)
}

func (b *ConnBuffer) PushI16(data int16) {
	b.pushNext(
		byte(data>>8),
		byte(data))
}

func (b *ConnBuffer) PushI32(data int32) {
	b.pushNext(
		byte(data>>24),
		byte(data>>16),
		byte(data>>8),
		byte(data))
}

func (b *ConnBuffer) PushI64(data int64) {
	b.pushNext(
		byte(data>>56),
		byte(data>>48),
		byte(data>>40),
		byte(data>>32),
		byte(data>>24),
		byte(data>>16),
		byte(data>>8),
		byte(data))
}

func (b *ConnBuffer) PushF32(data float32) {
	b.PushI32(int32(math.Float32bits(data)))
}

func (b *ConnBuffer) PushF64(data float64) {
	b.PushI64(int64(math.Float64bits(data)))
}

func (b *ConnBuffer) PushVrI(data int32) {
	for {
		temp := data & 0x7F
		data >>= 7

		if data != 0 {
			temp |= 0x80
		}

		b.pushNext(byte(temp))

		if data == 0 {
			break
		}
	}
}

func (b *ConnBuffer) PushVrL(data int64) {
	for {
		temp := data & 0x7F
		data >>= 7

		if data != 0 {
			temp |= 0x80
		}

		b.pushNext(byte(temp))

		if data == 0 {
			break
		}
	}
}

func (b *ConnBuffer) PushTxt(data string) {
	b.PushUAS([]byte(data), true)
}

func (b *ConnBuffer) PushUAS(data []byte, prefixWithLen bool) {
	if prefixWithLen {
		b.PushVrI(int32(len(data)))
	}

	b.pushNext(data...)
}

func (b *ConnBuffer) PushSAS(data []int8, prefixWithLen bool) {
	b.PushUAS(asUArray(data), prefixWithLen)
}

func (b *ConnBuffer) PushUID(data uuid.UUID) {
	msb, lsb := uuid.SigBits(data)

	b.PushI64(msb)
	b.PushI64(lsb)
}

func (b *ConnBuffer) PushPos(data data.PositionI) {
	b.PushI64(((data.X & 0x3FFFFFF) << 38) | ((data.Z & 0x3FFFFFF) << 12) | (data.Y & 0xFFF))
}

// internal
func (b *ConnBuffer) pullNext() byte {

	if b.iIndex >= b.Len() {
		return 0
		// panic("reached end of buffer")
	}

	next := b.bArray[b.iIndex]
	b.iIndex++

	if b.oIndex > 0 {
		b.oIndex--
	}

	return next
}

func (b *ConnBuffer) pullSize(next int) []byte {
	bytes := make([]byte, next)

	for i := 0; i < next; i++ {
		bytes[i] = b.pullNext()
	}

	return bytes
}

func (b *ConnBuffer) pushNext(bArray ...byte) {
	b.oIndex += int32(len(bArray))
	b.bArray = append(b.bArray, bArray...)
}

func (b *ConnBuffer) pullVariable(max int) int64 {
	var num int
	var res int64

	for {
		tmp := int64(b.pullNext())
		res |= (tmp & 0x7F) << uint(num*7)

		if num++; num > max {
			panic("VarInt > " + strconv.Itoa(max))
		}

		if tmp&0x80 != 0x80 {
			break
		}
	}

	return res
}

func asSArray(bytes []byte) []int8 {
	array := make([]int8, 0)

	for _, b := range bytes {
		array = append(array, int8(b))
	}

	return array
}

func asUArray(bytes []int8) []byte {
	array := make([]byte, 0)

	for _, b := range bytes {
		array = append(array, byte(b))
	}

	return array
}
