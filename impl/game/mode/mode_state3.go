package mode

import (
	"fmt"
	"time"

	"github.com/golangmc/minecraft-server/apis"
	"github.com/golangmc/minecraft-server/apis/game"
	"github.com/golangmc/minecraft-server/apis/logs"
	"github.com/golangmc/minecraft-server/apis/task"
	"github.com/golangmc/minecraft-server/apis/util"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/conf"
	"github.com/golangmc/minecraft-server/impl/data/values"

	"github.com/golangmc/minecraft-server/impl/game/commands"
	impl_event "github.com/golangmc/minecraft-server/impl/game/event"

	client_packet "github.com/golangmc/minecraft-server/impl/prot/client"
	server_packet "github.com/golangmc/minecraft-server/impl/prot/server"
)

func HandleState3(watcher util.Watcher, logger *logs.Logging, tasking *task.Tasking, join chan base.PlayerAndConnection, quit chan base.PlayerAndConnection, serverInfo *conf.ServerInfo) {
	fmt.Println("state 3")
	tasking.EveryTime(10, time.Second, func(task *task.Task) {

		api := apis.MinecraftServer()

		// I hate this, add a functional method for player iterating
		for _, player := range api.Players() {

			// also probably add one that returns both the player and their connection
			conn := api.ConnByUUID(player.UUID())

			// keep player connection alive via keep alive
			conn.SendPacket(&client_packet.PacketOKeepAlive{KeepAliveID: time.Now().UnixNano() / 1e6})
		}
	})

	watcher.SubAs(func(packet *server_packet.PacketIKeepAlive, conn base.Connection) {
		logger.DataF("player %s is being kept alive", conn.Address())
	})

	watcher.SubAs(func(packet *server_packet.PacketIPluginMessage, conn base.Connection) {
		api := apis.MinecraftServer()

		player := api.PlayerByConn(conn)
		if player == nil {
			return // log no player found?
		}

		api.Watcher().PubAs(impl_event.PlayerPluginMessagePullEvent{
			Conn: base.PlayerAndConnection{
				Connection: conn,
				Player:     player,
			},
			Channel: packet.Message.Chan(),
			Message: packet.Message,
		})
	})

	watcher.SubAs(func(packet *server_packet.PacketIChatSessionUpdate, conn base.Connection) {
		logger.DataF("player %s is updating their chat session", conn.Address())
		logger.DataF("session id: %s", packet.SessionID)
		logger.DataF("public key: %v", packet.PublicKey)
	})

	watcher.SubAs(func(packet *server_packet.PacketIChatCommand, conn base.Connection) {
		logger.DataF("player %s is sending a chat command: %s", conn.Address(), packet.Command)
		conn.SendPacket(&client_packet.PacketOSystemChat{
			Message: packet.Command,
			Overlay: false,
		})
		commands.CommandHandler(packet, conn)
	})

	watcher.SubAs(func(packet *server_packet.PacketIMovePlayerPos, conn base.Connection) {
		res := sendChunk(conn, int(packet.Position.X/16), int(packet.Position.Z/16), conn.Profile().MaxChunksCount)
		if res > 0 {
			fmt.Println("sent chunks", res)
		}
	})
	watcher.SubAs(func(packet *server_packet.PacketIMovePlayerPosRot, conn base.Connection) {
		res := sendChunk(conn, int(packet.Position.X/16), int(packet.Position.Z/16), conn.Profile().MaxChunksCount)
		if res > 0 {
			fmt.Println("sent chunks", res)
		}
	})

	watcher.SubAs(func(packet *server_packet.PacketIPlayerLoaded, conn base.Connection) {
		fmt.Println("player is loaded")

	})

	watcher.SubAs(func(packet *server_packet.PacketIChatMessage, conn base.Connection) {
	})

	go func() {
		for conn := range join {
			apis.MinecraftServer().Watcher().PubAs(impl_event.PlayerConnJoinEvent{Conn: conn})

			conn.SendPacket(&client_packet.PacketOJoinGame{
				EntityID:           int32(conn.EntityUUID()),
				Hardcore:           false,
				DoLimitedCrafting:  true,
				DimensionType:      0,
				DimensionName:      "minecraft:overworld",
				DeathDimensionName: "minecraft:overworld",
				PortalCooldown:     0,
				SeaLevel:           64,
				EnforceSecureChat:  true,
				PreviousGameMode:   game.SPECTATOR,
				IsDebug:            false,
				IsFlat:             false,
				HasDeathLocation:   false,
				GameMode:           game.SPECTATOR,
				DimensionNames:     []string{"minecraft:overworld"},
				HashedSeed:         values.DefaultWorldHashedSeed,
				MaxPlayers:         10,
				SimulationDistance: 12,
				ViewDistance:       12,
				ReduceDebug:        false,
				RespawnScreen:      false,
			})

			serverInfo.DynamicServerInfo.Online[conn.Profile().UUID.String()] = conf.SamplePlayer{
				Name: conn.Profile().Name,
				ID:   conn.Profile().UUID.String(),
			}

			prof := conn.Profile()
			prof.EntityID = int32(conn.EntityUUID())
			conn.SetProfile(prof)
			conn.SendPacket(&client_packet.PacketOEntityEvent{
				EntityID: int32(conn.EntityUUID()),
				EventID:  byte(28),
			})
			conn.SendPacket(&client_packet.PacketOInitializeBorder{
				X:                      0,
				Z:                      0,
				OldDiameter:            59999968,
				NewDiameter:            59999968,
				Speed:                  0.0,
				PortalTeleportBoundary: 29999984,
				WarningBlocks:          5,
				WarningTime:            15,
			})
			conn.SendPacket(&client_packet.PacketOGameEvent{EventID: 13, Data: 0.0})

			conn.SendPacket(&client_packet.PacketOSetChunkCacheCenter{
				X: 0,
				Z: 0,
			})
			conn.SendPacket(
				&client_packet.PacketOSetEntityMetadata{
					EntityID: int32(conn.Profile().EntityID),
					Metadata: []byte{9, 3, 0x41, 0xa0, 0, 0, 0x11, 0, 0x7f, 0xff},
				},
			)
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
					}, {
						ID:        21,
						Value:     0.1,
						Modifiers: []client_packet.ModifierData{},
					},
				},
			})

			// conn.SendPacket(&client_packet.PacketOPlayerInfoUpdate{
			// 	Actions: 0xff,
			// 	Values: []client.PlayerInfo{},
			// })

			// conn.SendPacket(&client_packet.PacketOPlayerInfoUpdate{
			// 	Action: 0xff,
			// 	Values: []client.PlayerInfo{},
			// })

			sendChunk(conn, 0, 0, conn.Profile().MaxChunksCount)

		}
	}()

	go func() {
		for conn := range quit {
			apis.MinecraftServer().Watcher().PubAs(impl_event.PlayerConnQuitEvent{Conn: conn})
			if conn.Profile() != nil {
				fmt.Println("player quit", conn.Profile().UUID.String())
				delete(serverInfo.DynamicServerInfo.Online, conn.Profile().UUID.String())
			}
		}
	}()
}

// var sendedChunks = make(map[string]bool)

func chunkKey(x, z int) string {
	return fmt.Sprintf("%d:%d", x, z)
}
func sendChunk(conn base.Connection, x, z int, maxChunksCount int) int {
	sendedChunks := conn.Profile().SendedChunks
	chunksToSend := []client_packet.PacketOLevelChunkWithLightFake{}

	radius := 10

	for i := -radius; i <= radius; i++ {
		for j := -radius; j <= radius; j++ {
			// Проверка: находится ли точка внутри круга
			if i*i+j*j > radius*radius {
				continue
			}

			key := chunkKey(x+i, z+j)
			if sendedChunks[key] {
				continue
			}
			chunksToSend = append(chunksToSend, client_packet.PacketOLevelChunkWithLightFake{
				Data: CreateChunk(x+i, z+j),
			})
			fmt.Println("sending chunk", key)
			sendedChunks[key] = true
			conn.Profile().SendedChunks = sendedChunks
			if len(chunksToSend) >= maxChunksCount {
				break
			}
		}
	}
	if len(chunksToSend) > 0 {
		conn.SendPacket(&client_packet.PacketOChunkBatchStart{})
		for _, chunk := range chunksToSend {
			conn.SendPacket(&chunk)
		}
		conn.SendPacket(&client_packet.PacketOChunkBatchFinished{
			BatchSize: int32(len(chunksToSend)),
		})
	}
	return len(chunksToSend)
}
