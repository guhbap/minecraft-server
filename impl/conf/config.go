package conf

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/disintegration/imaging"
	"github.com/golangmc/minecraft-server/apis/uuid"
	"gopkg.in/yaml.v2"
)

type ServerInfo struct {
	Network Network

	// read yaml config
	// allow_unlicensed: true
	// port: 25565
	// server_motd: "Hello, World!"
	// max_players: 20
	// server_icon: "image.png"

	AllowUnlicensed bool   `yaml:"allow_unlicensed"`
	Port            int    `yaml:"port"`
	ServerMotd      string `yaml:"server_motd"`
	MaxPlayers      int    `yaml:"max_players"`
	ServerIconPath  string `yaml:"server_icon_path"`
	ServerIcon      string

	DynamicServerInfo *DynamicServerInfo
}

func NewServerInfo(configFileName string) *ServerInfo {
	var res ServerInfo
	fmt.Println("reading config file", configFileName)
	data, err := os.ReadFile(configFileName)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &res)
	if err != nil {
		panic(err)
	}
	fmt.Println("config file readed", res)

	res.Network = Network{
		Host: "0.0.0.0",
		Port: res.Port,
	}

	res.DynamicServerInfo = NewDynamicServerInfo(res.ServerMotd, []SamplePlayer{})

	fmt.Println("reading icon file", res.ServerIconPath)

	// convert any type of file to png 64x64
	res.ServerIcon = convertToBase64Icon(res.ServerIconPath)

	return &res
}

type Network struct {
	Host string
	Port int
}

type SamplePlayer struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type PlayerData struct {
	UUID      uuid.UUID
	Name      string
	EntityID  int32
	OtherData map[string]any
}

type DynamicServerInfo struct {
	ServerMotd string
	Online     map[string]SamplePlayer
	Players    map[string]*PlayerData
}

func NewDynamicServerInfo(serverMotd string, online []SamplePlayer) *DynamicServerInfo {
	return &DynamicServerInfo{
		ServerMotd: serverMotd,
		Online:     make(map[string]SamplePlayer),
		Players:    make(map[string]*PlayerData),
	}
}

func convertToBase64Icon(path string) string {
	// Открытие файла вручную (чтобы можно было ловить более точные ошибки)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to open server icon:", err)
		return ""
	}
	defer file.Close()

	// Декодирование изображения (любой поддерживаемый тип)
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Failed to decode server icon:", err)
		return ""
	}

	// Изменение размера
	resized := imaging.Resize(img, 64, 64, imaging.Lanczos)

	// Кодирование в PNG
	var buf bytes.Buffer
	err = imaging.Encode(&buf, resized, imaging.PNG)
	if err != nil {
		fmt.Println("Failed to encode server icon:", err)
		return ""
	}

	// Возврат строки с base64 Data URI
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}
