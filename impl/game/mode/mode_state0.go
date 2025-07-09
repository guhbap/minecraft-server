package mode

import (
	"fmt"

	"github.com/golangmc/minecraft-server/apis/util"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/conf"
	"github.com/golangmc/minecraft-server/impl/prot/server"
)

/**
 * handshake
 */

func HandleState0(watcher util.Watcher, serverInfo *conf.ServerInfo) {

	watcher.SubAs(func(packet *server.PacketIHandshake, conn base.Connection) {
		fmt.Println("handshake")
		fmt.Println("state", packet.State)
		conn.SetState(packet.State)
	})

}
