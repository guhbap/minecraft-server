package conn

import (
	"fmt"
	"net"
	"reflect"
	"slices"
	"strconv"

	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/logs"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/data/system"
)

type network struct {
	host string
	port int

	logger  *logs.Logging
	packets base.Packets

	join chan base.PlayerAndConnection
	quit chan base.PlayerAndConnection

	report chan system.Message
}

func NewNetwork(host string, port int, packet base.Packets, report chan system.Message, join chan base.PlayerAndConnection, quit chan base.PlayerAndConnection) base.Network {
	return &network{
		host: host,
		port: port,

		join: join,
		quit: quit,

		report: report,

		logger:  logs.NewLogging("network", logs.EveryLevel...),
		packets: packet,
	}
}

func (n *network) Load() {
	if err := n.startListening(); err != nil {
		n.report <- system.Make(system.FAIL, err)
		return
	}
}

func (n *network) Kill() {

}

func (n *network) startListening() error {
	ser, err := net.ResolveTCPAddr("tcp", n.host+":"+strconv.Itoa(n.port))
	if err != nil {
		return fmt.Errorf("address resolution failed [%v]", err)
	}

	tcp, err := net.ListenTCP("tcp", ser)
	if err != nil {
		return fmt.Errorf("failed to bind [%v]", err)
	}

	n.logger.InfoF("listening on %s:%d", n.host, n.port)

	go func() {
		for {
			con, err := tcp.AcceptTCP()

			if err != nil {
				n.report <- system.Make(system.FAIL, err)
				break
			}

			_ = con.SetNoDelay(true)
			_ = con.SetKeepAlive(true)

			go handleConnect(n, NewConnection(con))
		}
	}()

	return nil
}

func handleConnect(network *network, conn base.Connection) {
	network.logger.DataF("New Connection from &6%v", conn.Address())

	var inf = make([]byte, 1024)
	var leftover []byte // 🔸 буфер для хранения "обрезков"

	for {
		sze, err := conn.Pull(inf)

		if err != nil && err.Error() == "EOF" {
			network.quit <- base.PlayerAndConnection{Player: nil, Connection: conn}
			break
		}

		if err != nil || sze == 0 {
			_ = conn.Stop()
			network.quit <- base.PlayerAndConnection{Player: nil, Connection: conn}
			break
		}

		// 🔸 объединяем прошлые остатки и новый приходящий пакет
		data := append(leftover, conn.Decrypt(inf[:sze])...)
		buf := NewBufferWith(data)

		for buf.InI() < buf.Len() {
			if buf.UAS()[buf.InI()] == 0xFE {
				fmt.Println("LEGACY PING")
				buf.SkpLen(1) // Пропускаем байт LEGACY PING
				continue
			}

			startIndex := buf.InI()

			packetLen := buf.PullVrI()
			if packetLen == 0 {
				// 🔸 Не смогли прочитать длину — ждём следующую порцию
				break
			}

			if buf.Len()-buf.InI() < packetLen {
				// 🔸 Неполный пакет — сохраняем остаток и ждём продолжения
				leftover = buf.UAS()[startIndex:]
				break
			}

			packetData := buf.UAS()[buf.InI() : buf.InI()+packetLen]
			buf.SkpLen(packetLen)

			bufI := NewBufferWith(packetData)
			bufO := NewBuffer()

			handleReceive(network, conn, bufI)

			if bufO.Len() > 1 {
				temp := NewBuffer()
				temp.PushVrI(bufO.Len())

				comp := NewBuffer()
				comp.PushUAS(conn.Deflate(bufO.UAS()), false)

				temp.PushUAS(comp.UAS(), false)

				_, err := conn.Push(conn.Encrypt(temp.UAS()))
				if err != nil {
					network.logger.Fail("Failed to push client bound packet: %v", err)
				}
			}

			// После успешной обработки очищаем leftover
			leftover = nil
		}
	}
}

func handleReceive(network *network, conn base.Connection, bufI buff.Buffer) {
	uuid := bufI.PullVrI()

	packetI := network.packets.GetPacketI(uuid, conn.GetState())
	if packetI == nil {
		network.logger.DataF("unable to decode %v packet with uuid: 0x%02x", conn.GetState(), uuid)
		return
	}
	silentkList := []int32{

		// } //
		0x0b, 0x1c, 0x1d, 0x1e, 0x09, 0x29}

	if !slices.Contains(silentkList, uuid) {
		network.logger.DataF("GET packet: 0x%02x %d | %v | %v", packetI.UUID(), uuid, reflect.TypeOf(packetI), conn.GetState())
	}

	// populate incoming packet
	packetI.Pull(bufI, conn)

	network.packets.PubAs(packetI)
	network.packets.PubAs(packetI, conn)
}
