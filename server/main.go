package main

import (
	"flag"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/tinylib/msgp/msgp"

	"wars/lib/color"
	"wars/lib/game"
	"wars/lib/messages"
	"wars/lib/ui"
)

func getColors() map[color.RGBA]bool {
	return map[color.RGBA]bool{
		ui.Green:      false,
		ui.Blue:       false,
		ui.Yellow:     false,
		ui.Purple:     false,
		ui.LightBlue:  false,
		ui.Sky:        false,
		ui.Lime:       false,
		ui.Orange:     false,
		ui.LightGreen: false,
		ui.Brown:      false,
	}
}

type server struct {
	tcp        *net.TCPListener
	udp        *net.UDPConn
	clients    map[string]*srvClient
	game       *warsgame.Game
	rateTicker *time.Ticker
	fpsTicker  *time.Ticker
	quit       chan struct{}
	colors     map[color.RGBA]bool

	mu sync.Mutex
}

func main() {
	msgp.RegisterExtension(98, func() msgp.Extension { return new(messages.MessageBody) })
	msgp.RegisterExtension(99, func() msgp.Extension { return new(color.RGBA) })

	var tcpAddr, udpAddr string
	flag.StringVar(&tcpAddr, "tcpAddr", ":4200", "Server tcp address")
	flag.StringVar(&udpAddr, "udpAddr", ":4201", "Server udp address")
	flag.Parse()

	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	slog.SetDefault(logger)

	srv := &server{
		game:       warsgame.NewGame(),
		clients:    make(map[string]*srvClient),
		colors:     getColors(),
		rateTicker: time.NewTicker(time.Second / warsgame.TPS),
		fpsTicker:  time.NewTicker(time.Second / (warsgame.TPS / 2)),
		quit:       make(chan struct{}),
	}

	err := srv.listen(tcpAddr, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer srv.close()

	srv.initTickers()
	slog.Info("tickers started")
	slog.Info("server started")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	slog.Info("got signal")
	slog.Info("server stopped")
}
