package commands

import "strings"

func ParseCommand(command string) (string, []string) {
	parts := strings.Split(command, " ")
	return parts[0], parts[1:]
}
