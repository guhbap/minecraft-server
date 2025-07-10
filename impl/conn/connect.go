package conn

import (
	"bytes"
	"compress/zlib"
	"crypto/cipher"
	"fmt"
	"io"
	"net"
	"slices"

	"github.com/golangmc/minecraft-server/apis/game"
	"github.com/golangmc/minecraft-server/apis/rand"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/conn/crypto"
)

type connection struct {
	new bool
	tcp *net.TCPConn

	state base.PacketState

	certify Certify
	compact Compact
	profile *game.Profile
}

func NewConnection(conn *net.TCPConn) base.Connection {
	return &connection{
		new: true,
		tcp: conn,

		certify: Certify{},
		compact: Compact{},
	}
}
func (c *connection) SetProfile(profile *game.Profile) {
	c.profile = profile
}
func (c *connection) Profile() *game.Profile {
	return c.profile
}

func (c *connection) Address() net.Addr {
	return c.tcp.RemoteAddr()
}

func (c *connection) GetState() base.PacketState {
	return c.state
}

func (c *connection) SetState(state base.PacketState) {
	c.state = state
}

type Certify struct {
	name string

	used bool
	data []byte

	encrypt cipher.Stream
	decrypt cipher.Stream
}

func (c *connection) Encrypt(data []byte) (output []byte) {
	if !c.certify.used {
		return data
	}

	output = make([]byte, len(data))
	c.certify.encrypt.XORKeyStream(output, data)

	return
}

func (c *connection) Decrypt(data []byte) (output []byte) {
	if !c.certify.used {
		return data
	}

	output = make([]byte, len(data))
	c.certify.decrypt.XORKeyStream(output, data)

	return
}

func (c *connection) CertifyName() string {
	return c.certify.name
}

func (c *connection) CertifyData() []byte {
	return c.certify.data
}

func (c *connection) CertifyUpdate(secret []byte) {
	encrypt, decrypt, err := crypto.NewEncryptAndDecrypt(secret)

	c.certify.encrypt = encrypt
	c.certify.decrypt = decrypt

	if err != nil {
		panic(fmt.Errorf("failed to enable encryption for user: %s\n%v", c.CertifyName(), err))
	}

	c.certify.used = true
	c.certify.data = secret
}

func (c *connection) CertifyValues(name string) {
	c.certify.name = name
	c.certify.data = rand.RandomByteArray(4)
}

type Compact struct {
	used bool
	size int32
}

func (c *connection) Deflate(data []byte) (output []byte) {
	if !c.compact.used {
		return data
	}

	var out bytes.Buffer

	writer, _ := zlib.NewWriterLevel(&out, zlib.BestCompression)
	_, _ = writer.Write(data)
	_ = writer.Close()

	output = out.Bytes()

	return
}

func (c *connection) Inflate(data []byte) (output []byte) {
	if !c.compact.used {
		return data
	}

	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	_, _ = io.Copy(&out, reader)

	output = out.Bytes()

	return
}

func (c *connection) Pull(data []byte) (len int, err error) {
	len, err = c.tcp.Read(data)
	return
}

func (c *connection) Push(data []byte) (len int, err error) {
	len, err = c.tcp.Write(data)
	return
}

func (c *connection) Stop() (err error) {
	err = c.tcp.Close()
	return
}

func (c *connection) SendPacket(packet base.PacketO) {
	silentkList := []int32{0x20, 0x4D, 0x28, 0x0d, 0x0c}
	if !slices.Contains(silentkList, packet.UUID()) {
		fmt.Printf("sending packet: 0x%02x\n", packet.UUID())
	}

	bufO := NewBuffer()
	temp := NewBuffer()

	// write buffer
	bufO.PushVrI(packet.UUID())
	packet.Push(bufO, c)
	// fmt.Println("packet: ", hex.EncodeToString(bufO.UAS()))

	temp.PushVrI(bufO.Len())
	temp.PushUAS(bufO.UAS(), false)

	encrypted := c.Encrypt(temp.UAS())
	_, err := c.tcp.Write(encrypted)
	if err != nil {
		fmt.Println("error sending packet: ", err)
	}
}
