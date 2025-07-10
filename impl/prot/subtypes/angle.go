package subtypes

import "github.com/golangmc/minecraft-server/apis/buff"

type Angle byte

func (a Angle) Push(writer buff.Buffer) {
	writer.PushByt(byte(a))
}
