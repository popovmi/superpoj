package main

func (srv *server) initTickers() {
	go func() {
		for {
			select {
			case <-srv.rateTicker.C:
				srv.game.Tick()
			case <-srv.fpsTicker.C:
				srv.broadcastState()
			case <-srv.quit:
				srv.rateTicker.Stop()
				srv.fpsTicker.Stop()
				return
			}
		}
	}()
}
