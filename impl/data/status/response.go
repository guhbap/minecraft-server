package status

import (
	"github.com/golangmc/minecraft-server/impl/conf"
)

type Response struct {
	Version     Version `json:"version"`
	Players     Players `json:"players"`
	Description Message `json:"description"`
	Favicon     string  `json:"favicon"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Players struct {
	Max    int                 `json:"max"`
	Online int                 `json:"online"`
	Sample []conf.SamplePlayer `json:"sample"`
}

type Message struct {
	Text string `json:"text"`
}
