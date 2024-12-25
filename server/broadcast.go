package main

import (
	"log/slog"

	"wars/lib/messages"
)

func (srv *server) broadcastState() {
	b, err := messages.New(messages.SrvMsgGameState, srv.game).MarshalMsg(nil)
	if err != nil {
		slog.Error("could not marshal state", err.Error())
		return
	}
	srv.broadcastUDP(b)
}

func (srv *server) broadcastUDP(b []byte) {
	for _, player := range srv.game.Players {
		_ = srv.clients[player.ID].sendUDPBytes(b)
	}
}
