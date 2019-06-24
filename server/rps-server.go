package server

import (
	"os"
	"log"
	"net"
	"sync"
	"github.com/daniellok/rps/utils"
)

var currentMatchId uint64
var matchIdMutex sync.Mutex
var logger = log.New(os.Stdin, "RPS server | ", log.Ldate|log.Ltime)

func RunServer() {
	lobbyList = utils.SafeLobbyList{
		Lobbies : []utils.Lobby{},
	}

	currentMatchId = 0
	
	ln, err := net.Listen("tcp", ":8787")
	if err != nil {
		logger.Println("Error occured! ", err)
	}
	logger.Println("Listening for connections on port 8787...")
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Println("Error occured! ", err)
		}
		logger.Println("Accepted connection from ", conn.RemoteAddr())
		go handleInitialConnect(conn)
	}
}
