package subtypes

import (
	"bytes"

	"github.com/Tnze/go-mc/nbt"
	"github.com/golangmc/minecraft-server/apis/buff"
)

type NbtTextMessage struct {
	Type  string `nbt:"type,omitempty"`
	Text  string `nbt:"text,omitempty"`
	Color string `nbt:"color,omitempty"`

	Font          string `nbt:"font,omitempty"`
	Bold          bool   `nbt:"bold,omitempty"`
	Italic        bool   `nbt:"italic,omitempty"`
	Underlined    bool   `nbt:"underlined,omitempty"`
	Strikethrough bool   `nbt:"strikethrough,omitempty"`
	Obfuscated    bool   `nbt:"obfuscated,omitempty"`

	Extra []NbtTextMessage `nbt:"extra,omitempty"`
}

func (p *NbtTextMessage) Push(writer buff.Buffer) {
	buf := bytes.NewBuffer(nil)
	enc := nbt.NewEncoder(buf)
	enc.NetworkFormat(true)
	err := enc.Encode(p, "")
	if err != nil {
		panic(err)
	}
	writer.PushUAS(buf.Bytes(), false)
}
