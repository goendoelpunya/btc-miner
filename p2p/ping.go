package p2p

import (
	"encoding/binary"
	"errors"
	"math/rand"
)

const (
	PingPongPayloadLen = 8
)

type PingPayload struct {
	Nonce uint64
}

func NewPingMsg() *Msg {
	nonce := uint64(rand.Int63())
	var buf [PingPongPayloadLen]byte
	binary.LittleEndian.PutUint64(buf[:], nonce)
	msg,_ :=NewMsg("ping", buf[:])//因为是自己组装的消息，所以一定不会出错
	return msg
}

func NewPongMsg(nonce uint64) *Msg {
	var buf [PingPongPayloadLen]byte
	binary.LittleEndian.PutUint64(buf[:], nonce)
	msg,_ :=NewMsg("pong", buf[:])//因为是自己组装的消息，所以一定不会出错
	return msg
}

func (p *PingPayload) Serialize() []byte {
	var data [8]byte
	binary.LittleEndian.PutUint64(data[:], p.Nonce)
	return data[:]
}

func (p *PingPayload) Parse(data []byte) error {
	if len(data) != 8 {
		return errors.New("length of data is wrong(not 8)")
	}
	p.Nonce = binary.LittleEndian.Uint64(data)
	return nil
}
func (p *PingPayload) Len() int {
	return 8
}

func (node *Node) HandlePing(peer *Peer, payload []byte) error {
	var err error
	var msgPong *Msg = nil
	if msgPong, err = NewMsg("pong", payload); err != nil {
		return err
	}
	if err = MustWrite(peer.Conn, msgPong.Serialize()); err != nil {
		return err
	}
	return nil
}

func (node *Node) HandlePong(peer *Peer, payload []byte) error {
	//实现心跳机制
	peer.Alive <- true

	return nil
}
