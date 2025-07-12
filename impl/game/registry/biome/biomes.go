package biome

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var Biomes = map[int]string{}

var InvertedBiomes = map[string]int{}

func init() {
	// todo это работает неправильно, нужно исправить. Ломается цвет травы

	// read all files in dir /registry/biome
	files, err := os.ReadDir("registry/biome")
	if err != nil {
		panic(err)
	}
	biomes := []string{}
	for _, file := range files {
		// add filename without extension
		biomes = append(biomes, "minecraft:"+strings.TrimSuffix(file.Name(), ".json"))
	}
	// sort by alphabet
	sort.Strings(biomes)
	for i, biome := range biomes {
		Biomes[i] = biome
		InvertedBiomes[biome] = i
	}
	fmt.Println(Biomes)
}

func GetBiomeId(biomeName string) int {
	if id, ok := InvertedBiomes[biomeName]; ok {
		return id
	}
	panic(fmt.Sprintf("biome %s not found", biomeName))
}
