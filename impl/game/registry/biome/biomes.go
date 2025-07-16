package biome

import (
	"encoding/json"
	"fmt"
	"os"
)

type Biome struct {
	Index            int          `nbt:"-"`
	Name             string       `nbt:"-"`
	HasPrecipitation bool         `nbt:"has_precipitation"`
	Temperature      float32      `nbt:"temperature"`
	Downfall         float32      `nbt:"downfall"`
	Effects          BiomeEffects `nbt:"effects" json:"effects"`
}
type BiomeEffects struct {
	SkyColor           int32  `nbt:"sky_color" json:"sky_color"`
	WaterColor         int32  `nbt:"water_color" json:"water_color"`
	WaterFogColor      int32  `nbt:"water_fog_color" json:"water_fog_color"`
	FogColor           int32  `nbt:"fog_color" json:"fog_color"`
	FoliageColor       *int32 `nbt:"foliage_color,omitempty" json:"foliage_color"` // Опционально
	GrassColor         *int32 `nbt:"grass_color,omitempty" json:"grass_color"`     // Опционально
	GrassColorModifier string `nbt:"grass_color_modifier,omitempty" json:"grass_color_modifier"`
}

type PrismarineBiome struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Category         string  `json:"category"`
	Temperature      float64 `json:"temperature"`
	HasPrecipitation bool    `json:"has_precipitation"`
	Dimension        string  `json:"dimension"`
	DisplayName      string  `json:"displayName"`
	Color            int     `json:"color"`
}

var InvertedBiomes = map[string]Biome{}

var IdToTag = map[int]string{}
var TagToId = map[string]int{}

func init() {
	data, err := os.ReadFile("registry/prismarine/biomes.json")
	if err != nil {
		panic(err)
	}
	var biomes []PrismarineBiome
	if err := json.Unmarshal(data, &biomes); err != nil {
		panic(err)
	}

	for _, biome := range biomes {
		IdToTag[biome.ID] = "minecraft:" + biome.Name
		data, err := os.ReadFile("registry/biome/" + biome.Name + ".json")
		if err != nil {
			panic(err)
		}
		var biome2 Biome
		if err := json.Unmarshal(data, &biome2); err != nil {
			panic(err)
		}
		biome2.Name = "minecraft:" + biome.Name
		biome2.Index = biome.ID
		InvertedBiomes[biome2.Name] = biome2
	}
}
func GetBiomeId(biomeName string) int {
	// return 35
	if id, ok := InvertedBiomes[biomeName]; ok {
		return id.Index
	}
	panic(fmt.Sprintf("biome %s not found", biomeName))
}
