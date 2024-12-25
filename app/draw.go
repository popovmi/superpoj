package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	warsgame "wars/lib/game"
	"wars/lib/ui"
)

func (c *gameClient) Draw(screen *ebiten.Image) {
	switch c.ui.screen {
	case screenMain:
		label := "Input name:"
		textW, textH := text.Measure(label, ui.FontFace22, warsgame.LineSpacing)
		middleW, middleH := float64(c.ui.windowW/2), float64(c.ui.windowH/2)

		op := &text.DrawOptions{}
		op.LineSpacing = warsgame.LineSpacing
		op.GeoM.Translate(middleW-textW, middleH-textH)
		op.ColorScale.ScaleWithColor(ui.White.ToColorRGBA())
		text.Draw(screen, label, ui.FontFace22, op)

		c.ui.nameInput.Draw(screen)

	case screenGame:
		img := ebiten.NewImage(warsgame.FieldWidth, warsgame.FieldHeight)
		imgOp := &ebiten.DrawImageOptions{}
		imgOp.GeoM.Translate(-c.ui.cameraX, -c.ui.cameraY)
		img.DrawImage(c.ui.worldImg, imgOp)

		for _, p := range c.game.Players {
			c.drawPlayer(img, p)
		}
		c.drawPlayerList(screen)
		c.drawChaser(screen)
		screen.DrawImage(img, &ebiten.DrawImageOptions{})
	default:

	}
}

func (c *gameClient) drawPlayer(screen *ebiten.Image, p *warsgame.Player) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-c.ui.cameraX, -c.ui.cameraY)
	op.GeoM.Translate(p.X-warsgame.Radius, p.Y-warsgame.Radius)

	textOp := &text.DrawOptions{}

	nameStr := p.Name
	if p.ID == c.id {
		nameStr += " (you)"
	}
	textOp.ColorScale.ScaleWithColor(p.Color.ToColorRGBA())
	textW, textH := text.Measure(nameStr, ui.FontFace16, warsgame.LineSpacing)
	textOp.LineSpacing = warsgame.LineSpacing
	textOp.GeoM.Translate(-c.ui.cameraX, -c.ui.cameraY)
	textOp.GeoM.Translate(p.X-textW/2, p.Y-warsgame.Radius-textH)
	text.Draw(screen, nameStr, ui.FontFace16, textOp)

	if p.ID == c.game.CId {
		screen.DrawImage(c.ui.playerImgs[p.ID].chaseImg, op)
	} else {
		screen.DrawImage(c.ui.playerImgs[p.ID].baseImg, op)
	}
}

func (c *gameClient) drawPortals() {
	for _, link := range c.game.PortalLinks {
		for _, p := range []*warsgame.Portal{link.P1, link.P2} {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(p.X-warsgame.PortalRadius, p.Y-warsgame.PortalRadius)
			c.ui.worldImg.DrawImage(c.ui.portalImg, op)
		}
	}
}

func (c *gameClient) drawBricks() {
	for _, brick := range c.game.Bricks {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(brick.X, brick.Y)
		img := c.ui.horBrickImg
		if brick.D == "V" {
			img = c.ui.verBrickImg
		}
		c.ui.worldImg.DrawImage(img, op)
	}
}

func (c *gameClient) drawChaser(screen *ebiten.Image) {
	chaser := c.game.Players[c.game.CId]
	label := fmt.Sprintf("Chaser: %s", chaser.Name)
	textW, textH := text.Measure(label, ui.FontFace18, warsgame.LineSpacing)
	middleW, middleH := float64(c.ui.windowW/2), 25.0
	textOp := &text.DrawOptions{}
	textOp.LineSpacing = warsgame.LineSpacing
	textOp.GeoM.Translate(middleW-textW/2, middleH-textH)
	textOp.ColorScale.ScaleWithColor(chaser.Color.ToColorRGBA())
	text.Draw(screen, label, ui.FontFace18, textOp)
}

func (c *gameClient) drawPlayerList(screen *ebiten.Image) {
	t := time.Now().Unix()
	i := 1

	sorted := make([]*warsgame.Player, 0)
	for _, v := range c.game.Players {
		sorted = append(sorted, v)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].JoinedAt > sorted[j].JoinedAt
	})

	for _, player := range sorted {
		playerStr := fmt.Sprintf("%s | %dm | %d", player.Name, (t-player.JoinedAt)/60, player.ChaseCount)
		textW, textH := text.Measure(playerStr, ui.FontFace18, warsgame.LineSpacing)
		op := &text.DrawOptions{}
		op.LineSpacing = warsgame.LineSpacing
		op.GeoM.Translate(float64(c.ui.windowW)-textW, float64(c.ui.windowH)-textH*float64(i))
		op.ColorScale.ScaleWithColor(player.Color.ToColorRGBA())
		text.Draw(screen, playerStr, ui.FontFace18, op)
		i++
	}
}
