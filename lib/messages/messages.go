package messages

import (
	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

type MessageType int

const (
	ClMsgHello MessageType = iota
	ClMsgJoinGame
	ClMsgMove
	ClMsgBrake
	ClMsgTeleport
)

const (
	SrvMsgYourID MessageType = iota
	SrvMsgYouJoined
	SrvMsgPlayerJoined
	SrvMsgGameState
	SrvMsgPlayerMoved
	SrvMsgPlayerBraked
	SrvMsgPlayerTeleported
)

type Message struct {
	T MessageType `msg:"type"`
	B MessageBody `msg:"body"`
}

type ClientUDPMessage struct {
	Message
	ID string `msg:"id"`
}

func New(t MessageType, data msgp.Marshaler) *Message {
	body, err := (data).MarshalMsg(nil)
	if err != nil {
		return nil
	}
	return &Message{T: t, B: body}
}

func UDP(t MessageType, id string, data msgp.Marshaler) *ClientUDPMessage {
	m := New(t, data)
	return &ClientUDPMessage{Message: *m, ID: id}
}

type Empty struct {
}

type YourIDMsg struct {
	ID string `msg:"id"`
}

type JoinGameMsg struct {
	Name string `msg:"name"`
}

type MoveMsg struct {
	Dir string `msg:"dir"`
}

type UdpMoveMsg struct {
	ClientUDPMessage
	MoveMsg
}

type PlayerMovedMsg struct {
	ID  string `msg:"id"`
	Dir string `msg:"dir"`
}

type PlayerBrakedMsg struct {
	ID string `msg:"id"`
}

type PlayerTeleportedMsg struct {
	ID string `msg:"id"`
}

func Unmarshal[T msgp.Unmarshaler](msg T, b []byte) (T, error) {
	_, err := msg.UnmarshalMsg(b)
	if err != nil {
		return msg, err
	}
	return msg, nil
}
