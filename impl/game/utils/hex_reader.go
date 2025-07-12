package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadHexFile(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	strData := strings.TrimSpace(string(data))

	bytesData := []byte{}

	hexes := strings.Split(strData, " ") // one hex is string like 0xa1
	for _, hexStr := range hexes {
		cleanHex := strings.TrimPrefix(hexStr, "0x")
		// Разбираем строку в число
		value, err := strconv.ParseUint(cleanHex, 16, 8)
		if err != nil {
			fmt.Println("Ошибка при парсинге:", err)
			continue
		}
		bytesData = append(bytesData, byte(value))
	}
	return bytesData
}
