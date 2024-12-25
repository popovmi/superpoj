package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	warsgame "wars/lib/game"
	"wars/lib/ui"
)

func (c *gameClient) Update() error {
	switch c.ui.screen {

	case screenMain:
		if c.ui.nameInput == nil {
			c.ui.nameInput = ui.NewTextField(
				image.Rect(
					c.ui.windowW/2,
					c.ui.windowH/2-warsgame.TextFieldHeight,
					c.ui.windowW/2+warsgame.TextFieldWidth,
					c.ui.windowH/2,
				),
				false,
				ui.FontFace22,
			)
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if c.ui.nameInput.Contains(x, y) {
				c.ui.nameInput.Focus()
				c.ui.nameInput.SetSelectionStartByCursorPosition(x, y)
			} else {
				c.ui.nameInput.Blur()
			}
		}

		if err := c.ui.nameInput.Update(); err != nil {
			return err
		}

		x, y := ebiten.CursorPosition()
		if c.ui.nameInput.Contains(x, y) {
			ebiten.SetCursorShape(ebiten.CursorShapeText)
		} else {
			ebiten.SetCursorShape(ebiten.CursorShapeDefault)
		}

		if c.ui.nameInput.IsFocused() && ebiten.IsKeyPressed(ebiten.KeyEnter) {
			if err := c.joinGame(c.ui.nameInput.Value()); err != nil {
				return err
			}
		}

	case screenGame:
		c.game.Tick()
		err := c.handleMovement()
		if err != nil {
			return err
		}
		if c.game.Players[c.id] == nil {
			break
		}

		if c.ui.windowW < warsgame.FieldWidth {
			c.ui.cameraX = c.game.Players[c.id].X - float64(c.ui.windowW)/2
			c.ui.cameraX = clamp(c.ui.cameraX, -50, float64(warsgame.FieldWidth-c.ui.windowW)+50)
		}

		if c.ui.windowH < warsgame.FieldHeight {
			c.ui.cameraY = c.game.Players[c.id].Y - float64(c.ui.windowH)/2
			c.ui.cameraY = clamp(c.ui.cameraY, -50, float64(warsgame.FieldHeight-c.ui.windowH)+50)
		}

	}

	return nil
}

func (c *gameClient) handleMovement() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return c.teleport()
	}

	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		return c.brake()
	}

	h := chooseDir(isRight(), isLeft(), "r", "l")
	v := chooseDir(isUp(), isDown(), "u", "d")
	return c.move(h + v)
}

func chooseDir(m, om bool, d, od string) string {
	if m != om {
		if m {
			return d
		}
		if om {
			return od
		}
	}
	return ""
}

func isUp() bool {
	return ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp)
}

func isDown() bool {
	return ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown)
}

func isLeft() bool {
	return ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
}

func isRight() bool {
	return ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
