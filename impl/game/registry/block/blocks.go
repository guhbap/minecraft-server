package blockid

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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

func GetBlockID(blockName string, properties map[string]string) (int, error) {
	var err error
	if defaultBlockIDFinder == nil {
		defaultBlockIDFinder, err = NewBlockIDFinder("registry/blocks.json")
		if err != nil {
			return 0, err
		}
	}
	return defaultBlockIDFinder.GetBlockID(blockName, properties)
}

// NewBlockIDFinder создает новый экземпляр BlockIDFinder, загружая данные из JSON-файла
func NewBlockIDFinder(filePath string) (*BlockIDFinder, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var blocks BlocksMap
	if err := json.Unmarshal(data, &blocks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &BlockIDFinder{blocks: blocks}, nil
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
		if reflect.DeepEqual(state.Properties, properties) {
			return state.ID, nil
		}
	}

	// Если точное совпадение не найдено, возвращаем ошибку
	return 0, fmt.Errorf("no state found for block %s with properties %v", blockName, properties)
}
