package client

import (
	"bytes"

	"github.com/Tnze/go-mc/nbt"
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/data"
	"github.com/golangmc/minecraft-server/apis/data/msgs"
	"github.com/golangmc/minecraft-server/apis/ents"
	"github.com/golangmc/minecraft-server/apis/game"
	"github.com/golangmc/minecraft-server/apis/uuid"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/data/client"
	"github.com/golangmc/minecraft-server/impl/data/plugin"
)

type PacketOChatMessage struct {
	Message         msgs.Message
	MessagePosition msgs.MessagePosition
}

func (p *PacketOChatMessage) UUID() int32 {
	return 0x0F
}

func (p *PacketOChatMessage) Push(writer buff.Buffer, conn base.Connection) {
	message := p.Message

	if p.MessagePosition == msgs.HotBarText {
		message = *msgs.New(message.AsText())
	}

	writer.PushTxt(message.AsJson())
	writer.PushByt(byte(p.MessagePosition))
}

type PacketOJoinGame struct {
	EntityID           int32
	Hardcore           bool
	DimensionNames     []string
	MaxPlayers         int
	ViewDistance       int32
	SimulationDistance int32
	ReduceDebug        bool
	RespawnScreen      bool
	DoLimitedCrafting  bool
	DimensionType      int32
	DimensionName      string
	HashedSeed         int64
	GameMode           game.GameMode
	PreviousGameMode   game.GameMode
	IsDebug            bool
	IsFlat             bool
	HasDeathLocation   bool
	DeathDimensionName string
	DeathLocation      data.Location
	PortalCooldown     int32
	SeaLevel           int32
	EnforceSecureChat  bool
}

func (p *PacketOJoinGame) UUID() int32 {
	return 0x2c
}

func (p *PacketOJoinGame) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushI32(p.EntityID)
	writer.PushBit(p.Hardcore)
	writer.PushVrI(int32(len(p.DimensionNames)))
	for _, dimensionName := range p.DimensionNames {
		writer.PushTxt(dimensionName)
	}
	writer.PushVrI(int32(p.MaxPlayers))
	writer.PushVrI(p.ViewDistance)
	writer.PushVrI(p.SimulationDistance)
	writer.PushBit(p.ReduceDebug)
	writer.PushBit(p.RespawnScreen)
	writer.PushBit(p.DoLimitedCrafting)
	writer.PushVrI(p.DimensionType)
	writer.PushTxt(p.DimensionName)
	writer.PushI64(p.HashedSeed)
	writer.PushByt(p.GameMode.Encoded(p.Hardcore /* pull this value from somewhere */))
	writer.PushByt(p.PreviousGameMode.Encoded(p.Hardcore /* pull this value from somewhere */))
	writer.PushBit(p.IsDebug)
	writer.PushBit(p.IsFlat)
	writer.PushBit(false)
	// writer.PushBit(p.HasDeathLocation)
	// writer.PushTxt(p.DeathDimensionName)
	// writer.PushF64(p.DeathLocation.X)
	// writer.PushF64(p.DeathLocation.Y)
	// writer.PushF64(p.DeathLocation.Z)
	// writer.PushF32(p.DeathLocation.AxisX)
	// writer.PushF32(p.DeathLocation.AxisY)
	writer.PushVrI(p.PortalCooldown)
	writer.PushVrI(p.SeaLevel)
	writer.PushBit(p.EnforceSecureChat)
}

type PacketOPluginMessage struct {
	Message plugin.Message
}

func (p *PacketOPluginMessage) UUID() int32 {
	return 0x19
}

func (p *PacketOPluginMessage) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushTxt(p.Message.Chan())
	p.Message.Push(writer)
}

type PacketOPlayerLocation struct {
	Location data.Location
	Relative client.Relativity

	ID int32
}

func (p *PacketOPlayerLocation) UUID() int32 {
	return 0x36
}

func (p *PacketOPlayerLocation) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushF64(p.Location.X)
	writer.PushF64(p.Location.Y)
	writer.PushF64(p.Location.Z)

	writer.PushF32(p.Location.AxisX)
	writer.PushF32(p.Location.AxisY)

	p.Relative.Push(writer)

	writer.PushVrI(p.ID)
}

type PacketOKeepAlive struct {
	KeepAliveID int64
}

func (p *PacketOKeepAlive) UUID() int32 {
	return 0x27
}

func (p *PacketOKeepAlive) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushI64(p.KeepAliveID)
}

type PacketOServerDifficulty struct {
	Difficulty game.Difficulty
	Locked     bool // should probably always be true
}

func (p *PacketOServerDifficulty) UUID() int32 {
	return 0x0E
}

func (p *PacketOServerDifficulty) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushByt(byte(p.Difficulty))
	writer.PushBit(p.Locked)
}

type PacketOPlayerAbilities struct {
	Abilities   client.PlayerAbilities
	FlyingSpeed float32
	FieldOfView float32
}

func (p *PacketOPlayerAbilities) UUID() int32 {
	return 0x3a
}

func (p *PacketOPlayerAbilities) Push(writer buff.Buffer, conn base.Connection) {
	p.Abilities.Push(writer)

	writer.PushF32(p.FlyingSpeed)
	writer.PushF32(p.FieldOfView)
}

type PacketOHeldItemChange struct {
	Slot client.HotBarSlot
}

func (p *PacketOHeldItemChange) UUID() int32 {
	return 0x40
}

func (p *PacketOHeldItemChange) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushByt(byte(p.Slot))
}

type PacketODeclareRecipes struct {
	// Recipes []*Recipe // this doesn't exist yet ;(
	RecipeCount int32
}

func (p *PacketODeclareRecipes) UUID() int32 {
	return 0x5B
}

func (p *PacketODeclareRecipes) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(p.RecipeCount)
	// when recipes are implemented, instead of holding a recipe count, simply write the size of the slice, Recipe will implement BufferPush
}

type PalettedContainer struct {
	BitsPerEntry int8
	Palette      []int32
	DataArray    []int64
}

func (p *PalettedContainer) Push(writer buff.Buffer) {
	writer.PushByt(byte(p.BitsPerEntry))
	writer.PushVrI(int32(len(p.Palette)))
	for _, palette := range p.Palette {
		writer.PushVrI(palette)
	}
	writer.PushVrL(int64(len(p.DataArray)))
	for _, data := range p.DataArray {
		writer.PushI64(data)
	}
}

type ChunkSection struct {
	BlockCount  int16
	BlockStates PalettedContainer
	Biomes      PalettedContainer
}

func (p *ChunkSection) Push(writer buff.Buffer) {
	writer.PushI16(p.BlockCount)
	p.BlockStates.Push(writer)
	p.Biomes.Push(writer)
}

type Heightmap struct {
	Type int32
	Data []int64
}

func (p *Heightmap) Push(writer buff.Buffer) {
	writer.PushVrI(p.Type)
	writer.PushVrL(int64(len(p.Data)))
	for _, data := range p.Data {
		writer.PushI64(data)
	}
}

type ChunkData struct {
	Heightmaps    []Heightmap
	Data          [][]ChunkSection
	BlockEntities []BlockEntity
}

func (p *ChunkData) Push(writer buff.Buffer) {
	writer.PushVrI(int32(len(p.Heightmaps)))
	for _, heightmap := range p.Heightmaps {
		heightmap.Push(writer)
	}
	writer.PushVrI(int32(len(p.Data)))
	for _, section1 := range p.Data {
		for _, section2 := range section1 {
			section2.Push(writer)
		}
	}
	writer.PushVrI(int32(len(p.BlockEntities)))
	for _, blockEntity := range p.BlockEntities {
		blockEntity.Push(writer)
	}
}

type BitSet struct {
	Bits []int64
}

func (p *BitSet) Push(writer buff.Buffer) {
	writer.PushVrL(int64(len(p.Bits)))
	for _, bit := range p.Bits {
		writer.PushI64(bit)
	}
}

type BlockEntity struct {
	PackedXZ int8
	Y        int16
	Type     int32
	Data     []byte // nbt
}

func (p *BlockEntity) Push(writer buff.Buffer) {
	writer.PushByt(byte(p.PackedXZ))
	writer.PushI16(p.Y)
	writer.PushVrI(p.Type)
	writer.PushUAS(p.Data, true)
}

type ModifierData struct {
	ID        int32
	Amount    float64
	Operation byte
}

func (p *ModifierData) Push(writer buff.Buffer) {
	writer.PushVrI(p.ID)
	writer.PushF64(p.Amount)
	writer.PushByt(p.Operation)
}

type AttrProperty struct {
	ID        int32
	Value     float64
	Modifiers []ModifierData
}

func (p *AttrProperty) Push(writer buff.Buffer) {
	writer.PushVrI(p.ID)
	writer.PushF64(p.Value)
	writer.PushVrI(int32(len(p.Modifiers)))
	for _, modifier := range p.Modifiers {
		modifier.Push(writer)
	}
}

type PacketOUpdateAttributes struct {
	EntityID   int32
	Attributes []AttrProperty
}

func (p *PacketOUpdateAttributes) UUID() int32 {
	return 0x7c
}

func (p *PacketOUpdateAttributes) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(p.EntityID)
	writer.PushVrI(int32(len(p.Attributes)))
	for _, attribute := range p.Attributes {
		attribute.Push(writer)
	}
}

type PacketOEntityEvent struct {
	EntityID int32
	EventID  byte
}

func (p *PacketOEntityEvent) UUID() int32 {
	return 0x1f
}

func (p *PacketOEntityEvent) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushI32(p.EntityID)
	writer.PushByt(p.EventID)
}

type PacketOSetEntityMetadata struct {
	EntityID int32
	Metadata []byte
}

func (p *PacketOSetEntityMetadata) UUID() int32 {
	return 0x5d
}

func (p *PacketOSetEntityMetadata) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(p.EntityID)
	writer.PushUAS(p.Metadata, false)
}

type LightData struct {
	SkyLightMask        BitSet
	BlockLightMask      BitSet
	EmptySkyLightMask   BitSet
	EmptyBlockLightMask BitSet
	SkyLightArray       [][]byte
	BlockLightArray     [][]byte
}

func (p *LightData) Push(writer buff.Buffer) {
	p.SkyLightMask.Push(writer)
	p.BlockLightMask.Push(writer)
	p.EmptySkyLightMask.Push(writer)
	p.EmptyBlockLightMask.Push(writer)
	writer.PushVrI(int32(len(p.SkyLightArray)))
	for _, skyLight := range p.SkyLightArray {
		writer.PushUAS(skyLight, true)
	}
	writer.PushVrI(int32(len(p.BlockLightArray)))
	for _, blockLight := range p.BlockLightArray {
		writer.PushUAS(blockLight, true)
	}
}

type PacketOLevelChunkWithLight struct {
	ChunkX int
	ChunkZ int
	Data   ChunkData
	Light  LightData
}

func (p *PacketOLevelChunkWithLight) UUID() int32 {
	return 0x28
}

func (p *PacketOLevelChunkWithLight) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushI32(int32(p.ChunkX))
	writer.PushI32(int32(p.ChunkZ))

	p.Data.Push(writer)

	p.Light.Push(writer)
}

type PacketOPlayerInfo struct {
	Action client.PlayerInfoAction
	Values []client.PlayerInfo
}

func (p *PacketOPlayerInfo) UUID() int32 {
	return 0x34
}

func (p *PacketOPlayerInfo) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(int32(p.Action))
	writer.PushVrI(int32(len(p.Values)))

	for _, value := range p.Values {
		value.Push(writer)
	}
}

type PacketOEntityMetadata struct {
	Entity ents.Entity
}

func (p *PacketOEntityMetadata) UUID() int32 {
	return 0x44
}

func (p *PacketOEntityMetadata) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(int32(p.Entity.EntityUUID())) // questionable...

	// only supporting player metadata for now
	_, ok := p.Entity.(ents.Player)
	if ok {

		writer.PushByt(16) // index | displayed skin parts
		writer.PushVrI(0)  // type | byte

		skin := client.SkinParts{
			Cape: true,
			Head: true,
			Body: true,
			ArmL: true,
			ArmR: true,
			LegL: true,
			LegR: true,
		}

		skin.Push(writer)
	}

	writer.PushByt(0xFF)
}

type PacketOChunkBatchStart struct {
}

func (p *PacketOChunkBatchStart) UUID() int32 {
	return 0x0d
}

func (p *PacketOChunkBatchStart) Push(writer buff.Buffer, conn base.Connection) {

}

type PacketOChunkBatchFinished struct {
	BatchSize int32
}

func (p *PacketOChunkBatchFinished) UUID() int32 {
	return 0x0c
}

func (p *PacketOChunkBatchFinished) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(p.BatchSize)
}

type PacketOGameEvent struct {
	EventID byte
	Data    float32
}

func (p *PacketOGameEvent) UUID() int32 {
	return 0x23
}

func (p *PacketOGameEvent) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushByt(p.EventID)
	writer.PushF32(p.Data)
}

type PacketOPlayerPosition struct {
	TpId     int32
	Position data.PositionF
	Speed    data.PositionF
	Yaw      float32
	Pitch    float32
	Flags    int32
}

func (p *PacketOPlayerPosition) UUID() int32 {
	return 0x42
}

func (p *PacketOPlayerPosition) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(p.TpId)
	writer.PushF64(p.Position.X)
	writer.PushF64(p.Position.Y)
	writer.PushF64(p.Position.Z)
	writer.PushF64(p.Speed.X)
	writer.PushF64(p.Speed.Y)
	writer.PushF64(p.Speed.Z)
	writer.PushF32(p.Yaw)
	writer.PushF32(p.Pitch)
	writer.PushI32(p.Flags)
}

type ChunkBiomeData struct {
	Z    int
	X    int
	Data []byte
}

func (p *ChunkBiomeData) Push(writer buff.Buffer) {
	writer.PushI32(int32(p.Z))
	writer.PushI32(int32(p.X))
	writer.PushUAS(p.Data, true)
}

type PacketOChunkBiomes struct {
	Biomes []ChunkBiomeData
}

func (p *PacketOChunkBiomes) UUID() int32 {
	return 0x0E
}

func (p *PacketOChunkBiomes) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(int32(len(p.Biomes)))
	for _, biome := range p.Biomes {
		biome.Push(writer)
	}
}

type PacketOLevelChunkWithLightFake struct {
	Data []byte
}

func (p *PacketOLevelChunkWithLightFake) UUID() int32 {
	return 0x28
}

func (p *PacketOLevelChunkWithLightFake) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushUAS(p.Data, false)
}

type PacketOInitializeBorder struct {
	X                      float64
	Z                      float64
	OldDiameter            float64
	NewDiameter            float64
	Speed                  int64
	PortalTeleportBoundary int32
	WarningBlocks          int32
	WarningTime            int32
}

func (p *PacketOInitializeBorder) UUID() int32 {
	return 0x26
}

func (p *PacketOInitializeBorder) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushF64(p.X)
	writer.PushF64(p.Z)
	writer.PushF64(p.OldDiameter)
	writer.PushF64(p.NewDiameter)
	writer.PushVrL(p.Speed)
	writer.PushVrI(p.PortalTeleportBoundary)
	writer.PushVrI(p.WarningBlocks)
	writer.PushVrI(p.WarningTime)
}

type PacketOSetChunkCacheCenter struct {
	X int32
	Z int32
}

func (p *PacketOSetChunkCacheCenter) UUID() int32 {
	return 0x58
}

func (p *PacketOSetChunkCacheCenter) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(p.X)
	writer.PushVrI(p.Z)
}

type PlayerInfoUpdatePlayers struct {
	UUID uuid.UUID
}

func (p *PlayerInfoUpdatePlayers) Push(writer buff.Buffer) {
	writer.PushUID(p.UUID)
}

type PacketOPlayerInfoUpdate struct {
	Actions byte
	Players []PlayerInfoUpdatePlayers
}

func (p *PacketOPlayerInfoUpdate) UUID() int32 {
	return 0x40
}

func (p *PacketOPlayerInfoUpdate) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushByt(p.Actions)
	writer.PushVrI(int32(len(p.Players)))
	// writer.PushUID(p.PlayerUUID)

	for _, player := range p.Players {
		player.Push(writer)
	}
}

type PacketOSystemChat struct {
	Message string
	Overlay bool
}

func (p *PacketOSystemChat) UUID() int32 {
	return 0x73
}

func (p *PacketOSystemChat) Push(writer buff.Buffer, conn base.Connection) {

	message := NbtTextMessage{
		Text: p.Message,
	}
	message.Push(writer)
	// writer.PushTxt(`[{"text": "A", "color": "red"}, "B", "C"]`)
	writer.PushBit(p.Overlay)
}

type NbtTextMessage struct {
	Text string `nbt:"text"`
}

func (p *NbtTextMessage) Push(writer buff.Buffer) {
	buf := bytes.NewBuffer(nil)
	enc := nbt.NewEncoder(buf)
	enc.NetworkFormat(true)
	err := enc.Encode(p, "")
	if err != nil {
		panic(err)
	}
	writer.PushUAS(buf.Bytes(), false)
}
