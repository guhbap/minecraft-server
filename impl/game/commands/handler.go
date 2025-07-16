package commands

import (
	"errors"
	"math/rand"
	"strconv"

	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/data"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/data/client"
	client_packet "github.com/golangmc/minecraft-server/impl/prot/client"
	server_packet "github.com/golangmc/minecraft-server/impl/prot/server"
)

func CommandHandler(packet *server_packet.PacketIChatCommand, conn base.Connection) {

	command, args := ParseCommand(packet.Command)

	switch command {
	case "tp":
		TpCommand(args, conn)
	case "gamemode":
		err := GameModeCommand(args, conn)
		if err != nil {
			conn.SendPacket(&client_packet.PacketOSystemChat{
				Message: client_packet.NbtTextMessage{
					Type:  "text",
					Text:  err.Error(),
					Color: "red",
				},
				Overlay: false,
			})
		}
	case "spawn":
		err := SpawnEntity(conn, args)
		if err != nil {
			conn.SendPacket(&client_packet.PacketOSystemChat{
				Message: client_packet.NbtTextMessage{
					Type:  "text",
					Text:  err.Error(),
					Color: "red",
				},
				Overlay: false,
			})
		}
	}
}

var gamemodes = map[string]int{
	"survival":  0,
	"creative":  1,
	"adventure": 2,
	"spectator": 3,
}

func SpawnEntity(conn base.Connection, args []string) error {
	if len(args) < 3 {
		return errors.New("not enough arguments")
	}

	return nil
}

func GameModeCommand(args []string, conn base.Connection) error {
	if len(args) < 1 {
		return errors.New("not enough arguments")
	}
	mode_str := args[0]
	mode, ok := gamemodes[mode_str]
	if !ok {
		return errors.New("invalid gamemode")
	}

	if mode == 3 {
		conn.SendPacket(
			&client_packet.PacketOPlayerAbilities{
				Abilities: client.PlayerAbilities{
					Invulnerable: true,
					Flying:       true,
					AllowFlight:  true,
					InstantBuild: false,
				},
				FlyingSpeed: 0.05,
				FieldOfView: 0.1,
			})
		properties := make([]client_packet.Property, len(conn.Profile().Properties))
		for i, prop := range conn.Profile().Properties {
			properties[i] = client_packet.Property{
				Name:      prop.Name,
				Value:     prop.Value,
				Signature: prop.Signature,
			}
		}
		conn.SendPacket(&client_packet.PacketOPlayerInfoUpdate{
			Actions: 0x01,
			Players: []client_packet.PlayerInfoUpdatePlayer{
				{
					UUID: conn.Profile().UUID,
					Actions: []func(buff.Buffer){
						client_packet.ADD_PLAYER_ACTION(conn.Profile().Name, properties),
					},
				},
			},
		})

		conn.SendPacket(&client_packet.PacketOGameEvent{EventID: 3, Data: 3})
		conn.SendPacket(
			&client_packet.PacketOPlayerAbilities{
				Abilities: client.PlayerAbilities{
					Invulnerable: true,
					Flying:       true,
					AllowFlight:  true,
					InstantBuild: false,
				},
				FlyingSpeed: 0.05,
				FieldOfView: 0.1,
			})
		conn.SendPacket(&client_packet.PacketOSetEntityMetadata{
			EntityID: int32(conn.Profile().EntityID), // todo
			// Metadata: []byte{0, 0, 32, 255},
		})
		conn.SendPacket(&client_packet.PacketOUpdateAttributes{
			EntityID: int32(conn.Profile().EntityID),
			Attributes: []client_packet.AttrProperty{
				{
					ID:        6,
					Value:     4.5,
					Modifiers: []client_packet.ModifierData{},
				},
				{
					ID:        9,
					Value:     3,
					Modifiers: []client_packet.ModifierData{},
				},
			},
		})
	}

	// conn.SendPacket(&client_packet.PacketOGameEvent{
	// 	EventID: 3,
	// 	Data:    float32(mode),
	// })
	return nil
}

func TpCommand(args []string, conn base.Connection) {
	if len(args) < 3 {
		return
	}
	x_str := args[0]
	y_str := args[1]
	z_str := args[2]
	x, _ := strconv.ParseFloat(x_str, 64)
	y, _ := strconv.ParseFloat(y_str, 64)
	z, _ := strconv.ParseFloat(z_str, 64)

	conn.SendPacket(&client_packet.PacketOPlayerPosition{
		TpId:     int32(rand.Intn(1000000)),
		Position: data.PositionF{X: float64(x), Y: float64(y), Z: float64(z)},
		Speed:    data.PositionF{X: 0, Y: 0, Z: 0},
		Yaw:      0,
		Pitch:    0,
		Flags:    0,
	})
}
