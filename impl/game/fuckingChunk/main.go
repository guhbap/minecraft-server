package fuckingChunk

type Chunk struct {
	X, Z        int32
	DataVersion int32
	Status      string
	LastUpdate  int64

	Sections []Section

	Heightmaps     map[string][]int32
	Entities       []Entity // можно оставлять raw NBT
	BlockEntities  []BlockEntity
	ScheduledTicks []TickEntry
}

type Section struct {
	Y            int8
	BlockPalette []PaletteEntry
	BlockData    []uint64 // packed indices
	BiomePalette []string
	BiomeData    []uint64
	LightBlock   []byte
	LightSky     []byte
}

type PaletteEntry struct {
	Name       string
	Properties map[string]string
}
type Entity struct {
	ID       string
	Pos      [3]float64
	Motion   [3]float64
	Rotation [2]float32
	NBT      map[string]interface{} // остальные поля как raw NBT
}
type BlockEntity struct {
	ID      string
	X, Y, Z int32
	NBT     map[string]interface{}
}

type TickEntry struct {
	SectionY  int8
	PosPacked uint16
	Type      string
	Priority  int32
	Delay     int32 // задержка до обновления (tick)
}
