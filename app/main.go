package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinylib/msgp/msgp"

	warscolor "wars/lib/color"
	"wars/lib/game"
	"wars/lib/messages"
	"wars/lib/ui"
)

const (
	defaultWindowWidth  = 800
	defaultWindowHeight = 600
)

var tcpAddr, udpAddr string

func main() {
	msgp.RegisterExtension(98, func() msgp.Extension { return new(messages.MessageBody) })
	msgp.RegisterExtension(99, func() msgp.Extension { return new(warscolor.RGBA) })

	ui.InitFont()

	c := newGameClient()

	go c.openTCPConnection()
	defer func() {
		c.TCPConn.Close()
		if c.UDPConn != nil {
			c.UDPConn.Close()
		}
	}()

	ebiten.SetWindowTitle("WARS")
	ebiten.SetWindowSize(defaultWindowWidth, defaultWindowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(warsgame.TPS)

	if err := ebiten.RunGame(c); err != nil {
		log.Fatal(err)
	}
}
