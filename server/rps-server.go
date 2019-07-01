package server

import (
	"os"
	"log"
	"net"
	"net/rpc"
	"sync"
	"github.com/daniellok/rps/types"
)

var currentMatchId uint64
var matchIdMutex sync.Mutex
var logger = log.New(os.Stdin, "RPS server | ", log.Ldate|log.Ltime)

func RunServer() {
	lobbyList = types.SafeLobbyList{
		Lobbies : []types.Lobby{},
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
		go rpc.ServeConn(conn)
		go handleInitialConnect(conn)
	}
}
