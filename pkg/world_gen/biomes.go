package worldgen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
)

type ClimateData struct {
	Temperature float64 `json:"temperature"`
	Downfall    float64 `json:"downfall"`
	// Добавьте другие поля, если они есть в JSON
}

func loadClosestFile(dir string, targetTemp, targetDownfall float64) (*ClimateData, string, error) {
	var closestData *ClimateData
	var closestFile string
	minDistance := math.MaxFloat64

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var climate ClimateData
		if err := json.Unmarshal(data, &climate); err != nil {
			return err
		}

		distance := math.Pow(climate.Temperature-targetTemp, 2) + math.Pow(climate.Downfall-targetDownfall, 2)
		if distance < minDistance {
			minDistance = distance
			closestData = &climate
			closestFile = path
		}
		return nil
	})

	if err != nil {
		return nil, "", err
	}

	if closestData == nil {
		return nil, "", fmt.Errorf("не найдено подходящих файлов")
	}

	return closestData, closestFile, nil
}
