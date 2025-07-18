package prot

import (
	"github.com/golangmc/minecraft-server/apis/logs"
	"github.com/golangmc/minecraft-server/apis/task"
	"github.com/golangmc/minecraft-server/apis/util"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/conf"
	"github.com/golangmc/minecraft-server/impl/game/mode"
	"github.com/golangmc/minecraft-server/impl/prot/server"
	stateplay "github.com/golangmc/minecraft-server/impl/prot/server/statePlay"
)

type packets struct {
	util.Watcher

	logger  *logs.Logging
	packetI map[base.PacketState]map[int32]func() base.PacketI // UUID to I server_data

}

func NewPackets(serverInfo *conf.ServerInfo, tasking *task.Tasking, join chan base.PlayerAndConnection, quit chan base.PlayerAndConnection) base.Packets {
	packets := &packets{
		Watcher: util.NewWatcher(),

		logger:  logs.NewLogging("protocol", logs.EveryLevel...),
		packetI: createPacketI(),
	}

	mode.HandleState0(packets, serverInfo)
	mode.HandleState1(packets, serverInfo)
	mode.HandleState2(packets, serverInfo)
	mode.HandleStateConfiguration(packets, join)
	mode.HandleState3(packets, packets.logger, tasking, join, quit, serverInfo)

	return packets
}

func (p *packets) GetPacketI(uuid int32, state base.PacketState) base.PacketI {
	creator := p.packetI[state][uuid]
	if creator == nil {
		return nil
	}

	return creator()
}

func createPacketI() map[base.PacketState]map[int32]func() base.PacketI {
	return map[base.PacketState]map[int32]func() base.PacketI{
		base.SHAKE: {
			0x00: func() base.PacketI {
				return &server.PacketIHandshake{}
			},
		},
		base.STATUS: {
			0x00: func() base.PacketI {
				return &server.PacketIRequest{}
			},
			0x01: func() base.PacketI {
				return &server.PacketIPing{}
			},
		},
		base.LOGIN: {
			0x00: func() base.PacketI {
				return &server.PacketILoginStart{}
			},
			0x01: func() base.PacketI {
				return &server.PacketIEncryptionResponse{}
			},
			0x02: func() base.PacketI {
				return &server.PacketILoginPluginResponse{}
			},
			0x03: func() base.PacketI {
				return &server.PacketILoginAcknowledged{}
			},
		},
		base.CONFIGURATION: {
			0x00: func() base.PacketI {
				return &server.PacketIClientInformation{}
			},
			0x02: func() base.PacketI {
				return &server.PacketICustomPayload{}
			},
			0x03: func() base.PacketI {
				return &server.PacketIFinishConfiguration{}
			},
			0x07: func() base.PacketI {
				return &server.PacketISelectKnownPacks{}
			},
		},
		base.PLAY: {
			0x00: func() base.PacketI {
				return &server.PacketIAcceptTeleportation{}
			},
			0x05: func() base.PacketI {
				return &server.PacketIChatCommand{}
			},
			0x08: func() base.PacketI {
				return &server.PacketIChatSessionUpdate{}
			},
			0x09: func() base.PacketI {
				return &server.PacketIChunkBatchReceived{}
			},
			0x0B: func() base.PacketI {
				return &server.PacketIClientTickEnd{}
			},
			0x11: func() base.PacketI {
				return &server.PacketIContainerClose{}
			},
			0x1A: func() base.PacketI {
				return &server.PacketIKeepAlive{}
			},
			0x1C: func() base.PacketI {
				return &server.PacketIMovePlayerPos{}
			},
			0x1E: func() base.PacketI {
				return &server.PacketIMovePlayerRot{}
			},
			0x1D: func() base.PacketI {
				return &server.PacketIMovePlayerPosRot{}
			},
			0x1F: func() base.PacketI {
				return &stateplay.PacketIMovePlayerStatusOnly{}
			},
			0x26: func() base.PacketI {
				return &server.PacketIPlayerAbilities{}
			},
			0x27: func() base.PacketI {
				return &stateplay.PacketIPlayerAction{}
			},
			0x28: func() base.PacketI {
				return &server.PacketIPlayerCommand{}
			},
			0x29: func() base.PacketI {
				return &server.PacketIPlayerInput{}
			},
			0x2a: func() base.PacketI {
				return &server.PacketIPlayerLoaded{}
			},
			0x3A: func() base.PacketI {
				return &stateplay.PacketISwing{}
			},
		},
	}
}
