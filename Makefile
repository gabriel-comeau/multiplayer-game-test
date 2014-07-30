# Makefile for Multiplayer Game Test

main:
	go install github.com/gabriel-comeau/multiplayer-game-test/mpgtclient

server:
	go install github.com/gabriel-comeau/multiplayer-game-test/mpgtserver

tests:
	go install github.com/gabriel-comeau/multiplayer-game-test/loadtester

libs:
	go install github.com/gabriel-comeau/multiplayer-game-test/shared
	go install github.com/gabriel-comeau/multiplayer-game-test/protocol
	go install github.com/gabriel-comeau/multiplayer-game-test/texturemanager

clean:
	rm -f "$(GOPATH)/bin/mpgtserver"
	rm -f "$(GOPATH)/bin/mpgtclient"
	rm -f "$(GOPATH)/bin/loadtester"
	rm -rf "$(GOPATH)/pkg/linux_amd64/github.com/gabriel-comeau/multiplayer-game-test"

dep-clean:
	rm -rf "$(GOPATH)/pkg/linux_amd64/bitbucket.org/krepa098"

all:
	make clean
	make libs
	make
	make serv
