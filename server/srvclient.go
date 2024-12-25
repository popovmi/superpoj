package main

import (
	"log"
	"log/slog"
	"net"

	"github.com/tinylib/msgp/msgp"

	"wars/lib/game"
	"wars/lib/messages"
)

type srvClient struct {
	*warsgame.Player

	ip      string
	tcp     net.Conn
	udpAddr *net.UDPAddr
	udp     *net.UDPConn
}

func (c *srvClient) sendTCP(t messages.MessageType) error {
	msg := messages.New(t, &messages.Empty{})
	if err := msgp.Encode(c.tcp, msg); err != nil {
		log.Println("could not encode message", err)
		return err
	}
	return nil
}

func (c *srvClient) sendTCPWithBody(t messages.MessageType, data msgp.Marshaler) error {
	msg := messages.New(t, data)
	if err := msgp.Encode(c.tcp, msg); err != nil {
		log.Println("could not encode message", err)
		return err
	}
	return nil
}

func (c *srvClient) sendTCPBytes(b []byte) error {
	if _, err := c.tcp.Write(b); err != nil {
		slog.Error("could not send TCP message", err)
		return err
	}
	return nil
}

func (c *srvClient) sendUDPBytes(b []byte) error {
	if c.udp != nil {
		if _, err := c.udp.WriteToUDP(b, c.udpAddr); err != nil {
			slog.Error("could not send UDP message", err)
			return err
		}
	}
	return nil
}
