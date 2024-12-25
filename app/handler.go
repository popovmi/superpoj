package main

import (
	"log/slog"

	"wars/lib/game"
	"wars/lib/messages"
)

func (c *gameClient) handleMessage(msg messages.Message) error {
	switch msg.T {
	case messages.SrvMsgYourID:
		userIdMsg, err := messages.Unmarshal(&messages.YourIDMsg{}, msg.B)
		if err != nil {
			return err
		}

		slog.Info("received user id", "ID", userIdMsg.ID)

		c.id = userIdMsg.ID
		c.ui.screen = screenMain

	case messages.SrvMsgYouJoined:
		state, err := messages.Unmarshal(&warsgame.Game{}, msg.B)
		if err != nil {
			return err
		}
		c.game = state
		for _, player := range c.game.Players {
			c.createPlayerImgs(player.ID, player.Color.ToColorRGBA())
		}
		c.openUDPConnection()
		c.drawPortals()
		c.drawBricks()
		c.ui.screen = screenGame

	case messages.SrvMsgPlayerJoined:
		player, err := messages.Unmarshal(&warsgame.Player{}, msg.B)
		if err != nil {
			return err
		}
		c.game.Players[player.ID] = player
		c.createPlayerImgs(player.ID, player.Color.ToColorRGBA())

	case messages.SrvMsgGameState:
		state, err := messages.Unmarshal(&warsgame.Game{}, msg.B)
		if err != nil {
			return err
		}
		c.game.CId = state.CId
		for k, player := range c.game.Players {
			if updatedPlayer, ok := state.Players[k]; ok {
				player.X = updatedPlayer.X
				player.Y = updatedPlayer.Y
				player.Vx = updatedPlayer.Vx
				player.Vy = updatedPlayer.Vy
				player.Direction = updatedPlayer.Direction
				player.ChaseCount = updatedPlayer.ChaseCount
				delete(state.Players, k)
			} else {
				delete(c.game.Players, k)
				delete(c.ui.playerImgs, k)
			}
		}
		for k, player := range state.Players {
			c.game.Players[k] = player
			c.createPlayerImgs(player.ID, player.Color.ToColorRGBA())
		}

	case messages.SrvMsgPlayerMoved:
		movedMsg, err := messages.Unmarshal(&messages.PlayerMovedMsg{}, msg.B)
		if err != nil {
			return err
		}
		if movedMsg.ID != c.id {
			c.game.Players[movedMsg.ID].Move(movedMsg.Dir)
		}

	case messages.SrvMsgPlayerBraked:
		brakedMsg, err := messages.Unmarshal(&messages.PlayerBrakedMsg{}, msg.B)
		if err != nil {
			return err
		}

		if brakedMsg.ID != c.id {
			c.game.Players[brakedMsg.ID].Brake()
		}

	case messages.SrvMsgPlayerTeleported:
		portedMsg, err := messages.Unmarshal(&messages.PlayerTeleportedMsg{}, msg.B)
		if err != nil {
			return err
		}

		if portedMsg.ID != c.id {
			c.game.Teleport(portedMsg.ID)
		}

	default:
		return nil
	}

	return nil
}
