package mode

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"time"

	"github.com/golangmc/minecraft-server/apis"
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/data"
	"github.com/golangmc/minecraft-server/apis/game"
	"github.com/golangmc/minecraft-server/apis/logs"
	"github.com/golangmc/minecraft-server/apis/task"
	"github.com/golangmc/minecraft-server/apis/util"
	"github.com/golangmc/minecraft-server/apis/uuid"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/conf"
	t_conn "github.com/golangmc/minecraft-server/impl/conn"
	"github.com/golangmc/minecraft-server/impl/data/values"

	chunkcreator "github.com/golangmc/minecraft-server/impl/game/chunk_creator"
	"github.com/golangmc/minecraft-server/impl/game/chunk_utils"
	"github.com/golangmc/minecraft-server/impl/game/commands"
	impl_event "github.com/golangmc/minecraft-server/impl/game/event"
	block_registry "github.com/golangmc/minecraft-server/impl/game/registry/block"

	client_packet "github.com/golangmc/minecraft-server/impl/prot/client"
	c_stateplay "github.com/golangmc/minecraft-server/impl/prot/client/statePlay"
	server_packet "github.com/golangmc/minecraft-server/impl/prot/server"
	s_stateplay "github.com/golangmc/minecraft-server/impl/prot/server/statePlay"
	"github.com/golangmc/minecraft-server/impl/prot/subtypes"
	"github.com/golangmc/minecraft-server/impl/prot/subtypes/entityMetadata"
)

const (
	InGameChunkRadius = 10
	SpawnChunkRadius  = 2
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
	watcher.SubAs(func(packet *s_stateplay.PacketISwing, conn base.Connection) {
		fmt.Println("player is swinging: ", packet.Hand)
		needHand := c_stateplay.AnimationSwing
		switch packet.Hand {
		case s_stateplay.SwingHandOffhand:
			needHand = c_stateplay.AnimationSwingOffhand
		case s_stateplay.SwingHandMainHand:
			needHand = c_stateplay.AnimationSwing
		}

		BroadcastPacket(serverInfo, &c_stateplay.PacketOAnimate{
			EntityID:  int32(conn.Profile().EntityID),
			Animation: needHand,
		}, conn.Profile().UUID)
	})

	watcher.SubAs(func(packet *s_stateplay.PacketIPlayerAction, conn base.Connection) {
		conn.SendPacket(&c_stateplay.PacketOBlockChangedAck{
			Sequence: packet.Sequence,
		})
		if packet.Status == s_stateplay.PlayerActionStatusStartedDigging {
			newBlockId, _ := block_registry.GetBlockID("minecraft:air", nil)
			BroadcastPacket(serverInfo, &c_stateplay.PacketOBlockUpdate{
				Location: packet.Location,
				BlockId:  int32(newBlockId),
			})
		}
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
			Message: subtypes.NbtTextMessage{
				Type:  "text",
				Text:  packet.Command,
				Color: "white",
			},
			Overlay: false,
		})
		commands.CommandHandler(packet, conn)
	})

	watcher.SubAs(func(packet *server_packet.PacketIMovePlayerPos, conn base.Connection) {
		x := int(packet.Position.X / 16)
		z := int(packet.Position.Z / 16)
		if calculateDistance(x, z, int(conn.Profile().ChunksCacheCenterX), int(conn.Profile().ChunksCacheCenterZ)) > 8 {
			conn.SendPacket(&client_packet.PacketOSetChunkCacheCenter{
				X: int32(x),
				Z: int32(z),
			})
			conn.Profile().ChunksCacheCenterX = int32(x)
			conn.Profile().ChunksCacheCenterZ = int32(z)
		}
		sendChunksInRadius(conn, int(packet.Position.X/16), int(packet.Position.Z/16), conn.Profile().MaxChunksCount, InGameChunkRadius)
		if packet.Position.Y < -100 {
			conn.SendPacket(&client_packet.PacketOPlayerPosition{
				TpId:     int32(rand.Intn(1000000)),
				Position: data.PositionF{X: packet.Position.X, Y: 10, Z: packet.Position.Z},
				Speed:    data.PositionF{X: 0, Y: 0, Z: 0},
				Yaw:      0,
				Pitch:    0,
				Flags:    0,
			})
		}
		conn.Profile().UpdatePos(packet.Position.X, packet.Position.Y, packet.Position.Z)
		BroadcastPacket(serverInfo, &client_packet.PacketOEntityPositionSync{
			EntityID: int32(conn.Profile().EntityID),
			X:        packet.Position.X,
			Y:        packet.Position.Y,
			Z:        packet.Position.Z,
			VelX:     0,
			VelY:     0,
			VelZ:     0,
			Yaw:      conn.Profile().GetPosInfo().Yaw,
			Pitch:    conn.Profile().GetPosInfo().Pitch,
			OnGround: true,
		}, conn.Profile().UUID)
	})
	watcher.SubAs(func(packet *server_packet.PacketIMovePlayerRot, conn base.Connection) {
		// conn.Profile().UpdateYawPitch(packet.Yaw, packet.Pitch)

		// BroadcastPacket(serverInfo, &client_packet.PacketOMoveEntityRot{
		// 	EntityID: int32(conn.Profile().EntityID),
		// 	Yaw:      subtypes.Angle(packet.Yaw),
		// 	Pitch:    subtypes.Angle(packet.Pitch),
		// 	OnGround: packet.Flags&0x01 != 0,
		// }, conn.Profile().UUID)
		// из значений от 0 до 360 в значение от 0 до 255
		yaw := float32(packet.Yaw) / 360 * 255
		BroadcastPacket(serverInfo, &client_packet.PacketORotateHead{
			EntityID: int32(conn.Profile().EntityID),
			Yaw:      subtypes.Angle(yaw),
		}, conn.Profile().UUID)
		conn.Profile().UpdateYawPitch(float32(packet.Yaw), float32(packet.Pitch))
		BroadcastPacket(serverInfo, &client_packet.PacketOEntityPositionSync{
			EntityID: int32(conn.Profile().EntityID),
			X:        conn.Profile().GetPosInfo().X,
			Y:        conn.Profile().GetPosInfo().Y,
			Z:        conn.Profile().GetPosInfo().Z,
			VelX:     0,
			VelY:     0,
			VelZ:     0,
			Yaw:      conn.Profile().GetPosInfo().Yaw,
			Pitch:    conn.Profile().GetPosInfo().Pitch,
			OnGround: true,
		}, conn.Profile().UUID)
	})
	watcher.SubAs(func(packet *server_packet.PacketIMovePlayerPosRot, conn base.Connection) {
		x := int(packet.Position.X / 16)
		z := int(packet.Position.Z / 16)
		if calculateDistance(x, z, int(conn.Profile().ChunksCacheCenterX), int(conn.Profile().ChunksCacheCenterZ)) > 8 {
			conn.SendPacket(&client_packet.PacketOSetChunkCacheCenter{
				X: int32(x),
				Z: int32(z),
			})
			conn.Profile().ChunksCacheCenterX = int32(x)
			conn.Profile().ChunksCacheCenterZ = int32(z)
		}
		sendChunksInRadius(conn, int(packet.Position.X/16), int(packet.Position.Z/16), conn.Profile().MaxChunksCount, InGameChunkRadius)
		conn.Profile().UpdatePos(packet.Position.X, packet.Position.Y, packet.Position.Z)
		conn.Profile().UpdateYawPitch(float32(packet.Rotation.AxisX), float32(packet.Rotation.AxisY))
		BroadcastPacket(serverInfo, &client_packet.PacketOEntityPositionSync{
			EntityID: int32(conn.Profile().EntityID),
			X:        conn.Profile().GetPosInfo().X,
			Y:        conn.Profile().GetPosInfo().Y,
			Z:        conn.Profile().GetPosInfo().Z,
			VelX:     0,
			VelY:     0,
			VelZ:     0,
			Yaw:      conn.Profile().GetPosInfo().Yaw,
			Pitch:    conn.Profile().GetPosInfo().Pitch,
			OnGround: true,
		}, conn.Profile().UUID)
	})

	watcher.SubAs(func(packet *server_packet.PacketIPlayerLoaded, conn base.Connection) {
		fmt.Println("player is loaded")

	})

	watcher.SubAs(func(packet *server_packet.PacketIChatMessage, conn base.Connection) {
	})
	watcher.SubAs(func(packet *server_packet.PacketIChunkBatchReceived, conn base.Connection) {
		fmt.Println("player is receiving chunk batch: ", packet.ChunkOnTick)
		prof := conn.Profile()
		if !prof.Spawned {
			prof.Spawned = true
			conn.SendPacket(&client_packet.PacketOPlayerPosition{
				TpId:     int32(rand.Intn(1000000)),
				Position: data.PositionF{X: 0, Y: 62, Z: 0},
				Speed:    data.PositionF{X: 0, Y: 0, Z: 0},
				Yaw:      0,
				Pitch:    0,
				Flags:    0,
			})
		}
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
				PreviousGameMode:   0xff,
				IsDebug:            false,
				IsFlat:             false,
				HasDeathLocation:   false,
				GameMode:           game.CREATIVE,
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
			player := conf.PlayerData{
				UUID:      conn.Profile().UUID,
				Name:      conn.Profile().Name,
				EntityID:  int32(conn.Profile().EntityID),
				OtherData: make(map[string]any),
			}

			properties := make([]client_packet.Property, len(prof.Properties))

			for i, prop := range prof.Properties {
				properties[i] = client_packet.Property{
					Name:      prop.Name,
					Value:     prop.Value,
					Signature: prop.Signature,
				}
			}
			playerInfoPacket := &client_packet.PacketOPlayerInfoUpdate{
				Actions: 0x01 | 0x02 | 0x04 | 0x08 | 0x10 | 0x20 | 0x40 | 0x80,
				Players: []client_packet.PlayerInfoUpdatePlayer{
					{
						UUID: player.UUID,
						Actions: []func(buff.Buffer){
							client_packet.ADD_PLAYER_ACTION(player.Name, properties),
							client_packet.INITIALIZE_CHAT(player.UUID),
							client_packet.UPDATE_GAME_MODE(1),
							client_packet.UPDATE_LISTED(true),
							client_packet.UPDATE_LATENCY(0),
							client_packet.UPDATE_DISPLAY_NAME(player.Name),
							client_packet.UPDATE_LIST_PRIORITY(0),
							client_packet.UPDATE_HAT(false),
						},
					},
				},
			}

			addEntityPacket := &c_stateplay.PacketOAddEntity{
				EntityID:   player.EntityID,
				EntityUUID: player.UUID,
				Type:       subtypes.EntityTypesRegistry["minecraft:player"].Index,
				X:          7,
				Y:          64,
				Z:          24,
				Pitch:      0,
				Yaw:        0,
				HeadYaw:    0,
				Data:       0,
				VelocityX:  0,
				VelocityY:  0,
				VelocityZ:  0,
			}
			conn.SendPacket(playerInfoPacket)
			serverInfo.DynamicServerInfo.Players[conn.Profile().UUID.String()] = &player
			player.OtherData["playerInfoPacket"] = playerInfoPacket
			player.OtherData["addEntityPacket"] = addEntityPacket
			for _, player := range serverInfo.DynamicServerInfo.Players {
				if player.UUID == conn.Profile().UUID {
					continue
				}
				packet, ok := player.OtherData["playerInfoPacket"].(*client_packet.PacketOPlayerInfoUpdate)
				if !ok {
					continue
				}
				fmt.Println("sending packet about", player.Name)
				conn.SendPacket(packet)
				conn.SendPacket(&client_packet.PacketOBundle{})

				packet2, ok := player.OtherData["addEntityPacket"].(*c_stateplay.PacketOAddEntity)
				if !ok {
					continue
				}
				conn.SendPacket(packet2)
				conn.SendPacket(&client_packet.PacketOBundle{})
			}

			BroadcastAddPlayer(serverInfo, &player)
			sendChunksInRadius(conn, 0, 0, conn.Profile().MaxChunksCount, SpawnChunkRadius)

			ef := entityMetadata.GetLivingEntityFields()
			ef.Health.Value = float32(20.0)
			pf := entityMetadata.GetPlayerFields()
			pf.TheDisplayedSkinPartsBitMaskThat.Value = byte(0)
			conn.SendPacket(
				&client_packet.PacketOSetEntityMetadata{
					EntityID: int32(conn.Profile().EntityID),
					Metadata: []entityMetadata.EntityField{
						ef.Health,
						pf.TheDisplayedSkinPartsBitMaskThat,
					},
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
			BroadcastPacket(serverInfo, &client_packet.PacketOSystemChat{
				Message: subtypes.NbtTextMessage{
					Type:  "text",
					Text:  "player ",
					Color: "white",
					Extra: []subtypes.NbtTextMessage{
						{
							Type:   "text",
							Text:   conn.Profile().Name,
							Color:  "green",
							Italic: true,
						},
						{
							Type:  "text",
							Text:  " joined the game",
							Color: "white",
						},
					},
				},
				Overlay: false,
			})

			// conn.SendPacket(&client_packet.PacketOPlayerInfoUpdate{
			// 	Actions: 0xff,
			// 	Values: []client.PlayerInfo{},
			// })

			// conn.SendPacket(&client_packet.PacketOPlayerInfoUpdate{
			// 	Action: 0xff,
			// 	Values: []client.PlayerInfo{},
			// })

		}
	}()

	go func() {
		for conn := range quit {
			apis.MinecraftServer().Watcher().PubAs(impl_event.PlayerConnQuitEvent{Conn: conn})
			if conn.Profile() != nil {
				fmt.Println("player quit", conn.Profile().UUID.String())
				delete(serverInfo.DynamicServerInfo.Online, conn.Profile().UUID.String())
				delete(serverInfo.DynamicServerInfo.Players, conn.Profile().UUID.String())
				BroadcastPacket(serverInfo, &c_stateplay.PacketORemoveEntity{
					EntityIDs: []int32{conn.Profile().EntityID},
				})
				BroadcastPacket(serverInfo, &c_stateplay.PacketOPlayerInfoRemove{
					UUIDs: []uuid.UUID{conn.Profile().UUID},
				})
			}
		}
	}()
}

// var sendedChunks = make(map[string]bool)

var perlinSetting = chunkcreator.NewPerlinSetting(123)

func chunkKey(x, z int) game.ChunkPos {
	return game.ChunkPos{X: int32(x), Z: int32(z)}
}

var (
	defaultChunkPallete = chunk_utils.NewPallete(
		[]chunk_utils.ChunkSectionBlocksPaletteNbt{
			{
				Name: "minecraft:air",
			},
			{
				Name: "minecraft:grass_block",
				Properties: map[string]string{
					"snowy": "false",
				},
			},
		},
	)
)

func sendChunksInRadius(conn base.Connection, x, z int, maxChunksCount int, radius int) int {
	conn.Profile().SendedChunksL.Lock()
	sendedChunks := conn.Profile().SendedChunks
	chunksToSend := []client_packet.PacketOLevelChunkWithLightFake{}
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

			if true {
				// fmt.Println("sending chunk", x+i, z+j)
				sch := chunkcreator.SendingChunk{
					X: x + i,
					Z: z + j,
				}
				for secIndex := 0; secIndex < 24; secIndex++ {
					section := chunkcreator.CreateEmptySection(defaultChunkPallete)
					sch.Sections = append(sch.Sections, section)
					section.Y = secIndex - 4
				}
				tmpBuf := t_conn.ConnBuffer{}
				sch.GeneratePerlin(perlinSetting)
				// for y := 0; y < 20; y++ {
				// 	sch.SetBlock(5, y, 7, 1)
				// }
				sch.Push(&tmpBuf)
				chunksToSend = append(chunksToSend, client_packet.PacketOLevelChunkWithLightFake{
					Data: tmpBuf.UAS(),
				})
				tmpBuf.Reset()

			} else {

				ch := chunk_utils.LoadChunk(x+i, z+j)
				// ch := chunk_utils.LoadChunk(0, 0)
				if ch != nil {
					chunksToSend = append(chunksToSend, client_packet.PacketOLevelChunkWithLightFake{
						Data: chunk_utils.CreateFromNbt(*ch),
					})
				}
			}
			sendedChunks[key] = true
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

	conn.Profile().SendedChunksL.Unlock()

	return len(chunksToSend)
}

func BroadcastPacket(serverInfo *conf.ServerInfo, packet base.PacketO, exclude ...uuid.UUID) {
	for _, player := range serverInfo.DynamicServerInfo.Online {
		uuid, err := uuid.TextToUUID(player.ID)
		if err != nil {
			continue
		}
		if slices.Contains(exclude, uuid) {
			continue
		}
		conn := apis.MinecraftServer().ConnByUUID(uuid)
		if conn != nil {
			conn.SendPacket(packet)
		}
	}
}

func BroadcastAddPlayer(serverInfo *conf.ServerInfo, player *conf.PlayerData) {
	playerInfoPacket, ok := player.OtherData["playerInfoPacket"].(*client_packet.PacketOPlayerInfoUpdate)
	if !ok {
		return
	}
	addEntityPacket, ok := player.OtherData["addEntityPacket"].(*c_stateplay.PacketOAddEntity)
	if !ok {
		return
	}
	BroadcastPacket(serverInfo, playerInfoPacket, player.UUID)

	BroadcastPacket(serverInfo, &client_packet.PacketOBundle{}, player.UUID)
	BroadcastPacket(serverInfo, addEntityPacket, player.UUID)
	lf := entityMetadata.GetLivingEntityFields()
	lf.Health.Value = float32(20)
	pf := entityMetadata.GetPlayerFields()
	pf.TheDisplayedSkinPartsBitMaskThat.Value = byte(127)
	BroadcastPacket(serverInfo, &client_packet.PacketOSetEntityMetadata{
		EntityID: player.EntityID,
		Metadata: []entityMetadata.EntityField{
			lf.Health,
			pf.TheDisplayedSkinPartsBitMaskThat,
		},
	}, player.UUID)
	BroadcastPacket(serverInfo, &client_packet.PacketOBundle{}, player.UUID)
	BroadcastPacket(serverInfo, &client_packet.PacketOMoveEntityPos{
		EntityID: player.EntityID,
		DeltaX:   0,
		DeltaY:   0,
		DeltaZ:   0,
		OnGround: true,
	}, player.UUID)
}

func calculateDistance(x1, z1, x2, z2 int) float64 {
	return math.Sqrt(float64((x1-x2)*(x1-x2) + (z1-z2)*(z1-z2)))
}
