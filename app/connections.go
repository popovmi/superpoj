package main

import (
	"errors"
	"log"
	"log/slog"
	"net"

	"github.com/tinylib/msgp/msgp"

	"wars/lib/messages"
)

func (c *gameClient) openTCPConnection() {
	conn, err := net.Dial("tcp", c.tcpAddr)
	if err != nil {
		log.Fatal("Dial error:", err)
	}

	c.TCPConn = conn
	go c.handleTCP()
}

func (c *gameClient) handleTCP() {
	for {
		var msg messages.Message
		if err := msgp.Decode(c.TCPConn, &msg); err != nil {
			slog.Error("could not decode TCP message:", err.Error())
			break
		}

		err := c.handleMessage(msg)
		if err != nil {
			slog.Error("could not handle TCP message", err.Error())
			break
		}
	}
}

func (c *gameClient) openUDPConnection() {
	udpAddr, err := net.ResolveUDPAddr("udp", c.udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	c.UDPConn = conn
	if err := c.sendUDP(messages.ClMsgHello); err != nil {
		log.Fatal(err)
	}

	go c.handleUDP()
}

func (c *gameClient) handleUDP() {
	var msg messages.Message
	for {
		if err := msgp.Decode(c.UDPConn, &msg); err != nil {
			slog.Error("could not decode UDP message", err.Error())
			break
		}
		if err := c.handleMessage(msg); err != nil {
			slog.Error("could not handle UDP message", err.Error())
			break
		}
	}
}

func (c *gameClient) sendTCP(t messages.MessageType) error {
	return c.sendTCPWithBody(t, &messages.Empty{})
}

func (c *gameClient) sendUDP(t messages.MessageType) error {
	return c.sendUDPWithBody(t, &messages.Empty{})
}

func (c *gameClient) sendTCPWithBody(t messages.MessageType, data msgp.Marshaler) error {
	if err := msgp.Encode(c.TCPConn, messages.New(t, data)); err != nil {
		slog.Error("could not send TCP message", err.Error())
		return err
	}
	return nil
}

func (c *gameClient) sendUDPWithBody(t messages.MessageType, data msgp.Marshaler) error {
	if err := msgp.Encode(c.UDPConn, messages.UDP(t, c.id, data)); err != nil {
		slog.Error("could not send UDP message", err.Error())
		return err
	}
	return nil
}

func (c *gameClient) sendMsg(conType string, t messages.MessageType) error {
	switch conType {
	case "tcp":
		return c.sendTCP(t)
	case "udp":
		return c.sendUDP(t)
	default:
		return errors.New("unknown con type")
	}
}

func (c *gameClient) sendMsgWithBody(conType string, t messages.MessageType, data msgp.Marshaler) error {
	switch conType {
	case "tcp":
		return c.sendTCPWithBody(t, data)
	case "udp":
		return c.sendUDPWithBody(t, data)
	default:
		return errors.New("unknown con type")
	}
}
