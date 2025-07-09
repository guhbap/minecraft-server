package mode

import (
	"github.com/golangmc/minecraft-server/apis/data"
	"github.com/golangmc/minecraft-server/apis/util"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/conf"
	"github.com/golangmc/minecraft-server/impl/data/status"
	"github.com/golangmc/minecraft-server/impl/prot/client"
	"github.com/golangmc/minecraft-server/impl/prot/server"
)

/**
 * status
 */

func HandleState1(watcher util.Watcher, serverInfo *conf.ServerInfo) {

	watcher.SubAs(func(packet *server.PacketIRequest, conn base.Connection) {
		response := client.PacketOResponse{Status: status.Response{
			Version: status.Version{
				Name:     "GoLang Server",
				Protocol: data.CurrentProtocol.Protocol(),
			},
			Players: status.Players{
				Max:    serverInfo.MaxPlayers,
				Online: len(serverInfo.DynamicServerInfo.Online),
				Sample: make([]conf.SamplePlayer, 0),
			},
			Favicon: serverInfo.ServerIcon,
			Description: status.Message{
				Text: serverInfo.ServerMotd,
			},
		}}
		for _, player := range serverInfo.DynamicServerInfo.Online {
			response.Status.Players.Sample = append(response.Status.Players.Sample, player)
		}
		conn.SendPacket(&response)
	})

	watcher.SubAs(func(packet *server.PacketIPing, conn base.Connection) {
		response := client.PacketOPong{Ping: packet.Ping}
		conn.SendPacket(&response)
	})

}
