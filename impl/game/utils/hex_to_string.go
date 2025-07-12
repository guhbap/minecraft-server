package utils

import "fmt"

func HexToString(data []byte) string {
	result := "\n"
	counter := 0
	for i := 0; i < len(data); i++ {
		if counter%16 == 0 {
			result += fmt.Sprintf("\n%08X ", counter)
		}
		result += fmt.Sprintf("%02X ", data[i])
		counter++
	}
	return result
}
