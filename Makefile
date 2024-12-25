# Makefile

SERVER_TCP_ADDRESS ?= localhost:4200
SERVER_UDP_ADDRESS ?= localhost:4201
PUBLIC_TCP_ADDRESS ?= localhost:4200
PUBLIC_UDP_ADDRESS ?= localhost:4201

.PHONY: dev_server
dev_server:
	go build -o build/dev_server ./server
	build/dev_server -tcpAddr $(SERVER_TCP_ADDRESS) -udpAddr $(SERVER_UDP_ADDRESS)

.PHONY: dev_app
dev_app:
	go build -ldflags="-X 'main.tcpAddr=$(PUBLIC_TCP_ADDRESS)' -X 'main.udpAddr=$(PUBLIC_UDP_ADDRESS)'" -o build/dev_app ./app
	build/dev_app

.PHONY: clean
clean:
	-rm -f ./build