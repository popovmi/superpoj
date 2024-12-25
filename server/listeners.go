package main

import (
	"log/slog"
	"net"
)

func (srv *server) listen(tcpAddr, udpAddr string) error {
	tcpStarted := make(chan error)
	udpStarted := make(chan error)

	go srv.startTCPServer(tcpAddr, tcpStarted)
	err := <-tcpStarted
	if err != nil {
		return err
	}

	go srv.startUDPServer(udpAddr, udpStarted)
	return <-udpStarted
}

func (srv *server) startTCPServer(addr string, tcpStarted chan<- error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		slog.Error("could not resolve TCP address", err.Error())
		tcpStarted <- err
		return
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		slog.Error("could not start TCP listener", err.Error())
		tcpStarted <- err
		return
	}

	slog.Info("TCP listener started", "addr", tcpListener.Addr())

	srv.tcp = tcpListener
	tcpStarted <- nil

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			slog.Error("could not create TCP connection", err.Error())
			continue
		}
		go srv.handleTCPConnection(conn)
	}
}

func (srv *server) startUDPServer(addr string, udpStarted chan<- error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		slog.Error("could not resolve UDP address", err.Error())
		udpStarted <- err
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		slog.Error("could not start UDP listener", err.Error())
		udpStarted <- err
		return
	}

	slog.Info("UDP listener started", "addr", udpAddr)

	srv.udp = udpConn
	udpStarted <- nil

	buf := make([]byte, 1024)
	for {
		n, addr, err := srv.udp.ReadFromUDP(buf)
		if err != nil {
			slog.Error("could not read UDP data", err.Error())
			continue
		}

		err = srv.handleUDPData(addr, buf, n)
		if err != nil {
			continue
		}
	}
}

func (srv *server) close() {
	if srv.udp != nil {
		err := srv.udp.Close()
		if err != nil {
			slog.Error("could not close UDP", err.Error())
		}
	}

	if srv.tcp != nil {
		err := srv.tcp.Close()
		if err != nil {
			slog.Error("could not close UDP", err.Error())
		}
	}
}
