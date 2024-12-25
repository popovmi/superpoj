package main

import (
	"image/color"
	"net"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	warsgame "wars/lib/game"
	"wars/lib/ui"
)

type gameScreen = int

const (
	screenWait gameScreen = iota
	screenMain
	screenGame
)

type playerImg struct {
	baseImg  *ebiten.Image
	chaseImg *ebiten.Image
}

type gameUI struct {
	windowW          int
	windowH          int
	screen           gameScreen
	nameInput        *ui.TextField
	cameraX, cameraY float64
	worldImg         *ebiten.Image
	portalImg        *ebiten.Image
	horBrickImg      *ebiten.Image
	verBrickImg      *ebiten.Image
	playerImgs       map[string]*playerImg
}

type gameClient struct {
	id   string
	game *warsgame.Game
	ui   *gameUI

	tcpAddr string
	udpAddr string
	TCPConn net.Conn
	UDPConn *net.UDPConn

	quit chan struct{}

	mu sync.Mutex
}

func newGameClient() *gameClient {
	ui.InitFont()

	worldBound := ebiten.NewImage(warsgame.FieldWidth, warsgame.FieldHeight)
	vector.StrokeRect(worldBound, 0, 0, warsgame.FieldWidth, warsgame.FieldHeight, 1, color.White, true)

	portalImg := ebiten.NewImage(2*warsgame.PortalRadius, 2*warsgame.PortalRadius)
	vector.DrawFilledCircle(
		portalImg, warsgame.PortalRadius, warsgame.PortalRadius, warsgame.PortalRadius, ui.DarkGray.ToColorRGBA(), true)

	horBrickImg := ebiten.NewImage(200, 40)
	vector.StrokeRect(horBrickImg, 0, 0, 200, 40, 1, ui.Gray.ToColorRGBA(), true)

	verBrickImg := ebiten.NewImage(40, 200)
	vector.StrokeRect(verBrickImg, 0, 0, 40, 200, 1, ui.Gray.ToColorRGBA(), true)

	worldImg := ebiten.NewImage(warsgame.FieldWidth, warsgame.FieldHeight)
	worldImg.DrawImage(worldBound, &ebiten.DrawImageOptions{})

	return &gameClient{
		ui: &gameUI{
			screen:      screenWait,
			worldImg:    worldBound,
			portalImg:   portalImg,
			horBrickImg: horBrickImg,
			verBrickImg: verBrickImg,
			playerImgs:  make(map[string]*playerImg),
		},
		tcpAddr: tcpAddr,
		udpAddr: udpAddr,
		quit:    make(chan struct{}),
	}
}

func (c *gameClient) createPlayerImgs(id string, clr color.RGBA) {
	baseImg := ebiten.NewImage(2*warsgame.Radius, 2*warsgame.Radius)
	vector.DrawFilledCircle(baseImg, warsgame.Radius, warsgame.Radius, warsgame.Radius, clr, true)

	chaseImg := ebiten.NewImage(2*warsgame.Radius, 2*warsgame.Radius)
	vector.StrokeCircle(chaseImg, warsgame.Radius, warsgame.Radius, warsgame.Radius-5, 5, clr, true)

	c.ui.playerImgs[id] = &playerImg{baseImg, chaseImg}
}

func (c *gameClient) Layout(w, h int) (int, int) {
	if c.ui.windowW != w || c.ui.windowH != h {
		c.ui.windowW = w
		c.ui.windowH = h

		if w > warsgame.FieldWidth {
			c.ui.windowW = warsgame.FieldWidth

		}
		if h > warsgame.FieldHeight {
			c.ui.windowH = warsgame.FieldHeight
		}

		c.ui.nameInput = nil
	}
	return c.ui.windowW, c.ui.windowH
}
