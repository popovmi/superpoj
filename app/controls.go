package main

import (
	"wars/lib/messages"
)

func (c *gameClient) joinGame(name string) error {
	c.ui.screen = screenWait
	err := c.sendTCPWithBody(messages.ClMsgJoinGame, &messages.JoinGameMsg{Name: name})
	if err != nil {
		return err
	}
	return nil
}

func (c *gameClient) move(dir string) error {
	p, ok := c.game.Players[c.id]
	if ok && dir != p.Direction {
		err := c.sendUDPWithBody(messages.ClMsgMove, &messages.MoveMsg{Dir: dir})
		if err != nil {
			return err
		}
		p.Move(dir)
	}
	return nil
}

func (c *gameClient) brake() error {
	p, ok := c.game.Players[c.id]
	if ok {
		err := c.sendUDP(messages.ClMsgBrake)
		if err != nil {
			return err
		}
		p.Brake()
	}
	return nil
}

func (c *gameClient) teleport() error {
	_, ok := c.game.Players[c.id]
	if ok {
		err := c.sendTCP(messages.ClMsgTeleport)
		if err != nil {
			return err
		}
		c.game.Teleport(c.id)
	}
	return nil
}
