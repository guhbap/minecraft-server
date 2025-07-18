package blockid

import (
	"encoding/json"
	"fmt"
	"os"
)

// BlockState описывает состояние блока в JSON
type BlockState struct {
	ID         int               `json:"id"`
	Properties map[string]string `json:"properties"`
	Default    bool              `json:"default,omitempty"`
}

// BlockData описывает блок в JSON
type BlockData struct {
	Definition struct {
		Type string `json:"type"`
	} `json:"definition"`
	Properties map[string][]string `json:"properties"`
	States     []BlockState        `json:"states"`
}

// BlocksMap хранит данные блоков из JSON
type BlocksMap map[string]BlockData

// BlockIDFinder управляет поиском ID блоков
type BlockIDFinder struct {
	blocks BlocksMap
}

var defaultBlockIDFinder *BlockIDFinder

func init() {
	data, err := os.ReadFile("registry/blocks.json")
	if err != nil {
		panic(err)
	}

	var blocks BlocksMap
	if err := json.Unmarshal(data, &blocks); err != nil {
		panic(err)
	}

	defaultBlockIDFinder = &BlockIDFinder{blocks: blocks}
}

func GetBlockID(blockName string, properties map[string]string) (int, error) {

	return defaultBlockIDFinder.GetBlockID(blockName, properties)
}

// GetBlockID возвращает ID блока по его имени и свойствам
func (f *BlockIDFinder) GetBlockID(blockName string, properties map[string]string) (int, error) {
	block, exists := f.blocks[blockName]
	if !exists {
		return 0, fmt.Errorf("block %s not found", blockName)
	}

	// Проверяем, что все переданные свойства существуют в блоке
	for prop := range properties {
		if _, ok := block.Properties[prop]; !ok {
			return 0, fmt.Errorf("property %s not found for block %s", prop, blockName)
		}
	}

	// Ищем состояние, соответствующее переданным свойствам
	for _, state := range block.States {
		matches := true
		// Если свойства nil или пустые, они должны совпадать
		if (properties == nil && state.Properties != nil) ||
			(properties != nil && state.Properties == nil) {
			matches = false
			continue
		}

		// Проверяем совпадение количества свойств
		if len(properties) != len(state.Properties) {
			matches = false
			continue
		}

		// Проверяем совпадение значений свойств
		for key, value := range properties {
			if stateValue, exists := state.Properties[key]; !exists || stateValue != value {
				matches = false
				break
			}
		}

		if matches {
			return state.ID, nil
		}
	}

	// Если точное совпадение не найдено, возвращаем ошибку
	return 0, fmt.Errorf("no state found for block %s with properties %v", blockName, properties)
}
