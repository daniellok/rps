package server

import (
	"net"
	"bufio"
	"encoding/gob"
	"encoding/binary"
	"github.com/daniellok/rps/types"
)

var lobbyList types.SafeLobbyList

func handleInitialConnect(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		b, err := reader.ReadByte()
		if err != nil {
			logger.Println(conn.RemoteAddr(), "disconnected")
			return
		}

		if b == types.CREATE_LOBBY {
			logger.Println("Lobby created:", conn.RemoteAddr())
			go createLobby(conn)
			return
		} else if b == types.JOIN_LOBBY {
			go joinLobby(conn)
			return
		} else {
			writer.WriteByte(types.INVALID_CHOICE)
		}
	}
}

func createLobby(conn net.Conn) {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	
	err := writer.WriteByte(types.LOBBY_CREATED)
	if err != nil {
		logger.Println(conn.RemoteAddr(), "disconnected")
		return
	}
	err = writer.Flush()
	if err != nil {
		logger.Println(conn.RemoteAddr(), "disconnected")
		return
	}

	matchIdMutex.Lock()
	id := currentMatchId
	currentMatchId += 1
	matchIdMutex.Unlock()
	
	lobby := types.Lobby{
		Player : conn,
		Name : "Harold Lobby",
		Id : id,
	}

	lobbyList.AddLobby(lobby)

	b, err := reader.ReadByte()
	if err != nil || b != types.RECEIVED_MATCH {
		logger.Println(conn.RemoteAddr(), "disconnected, cleaning up lobby...")
		lobbyList.RemoveLobbyWithId(id)
	}	
}

func joinLobby(conn net.Conn) {
	sendLobbyList(conn)
	waitForLobbyChoice(conn)
}

func sendLobbyList(conn net.Conn) {
	encoder := gob.NewEncoder(conn)
	marshallable := lobbyList.ToMarshallable()
	encoder.Encode(marshallable)
}

func waitForLobbyChoice(conn net.Conn) {
	reader  := bufio.NewReader(conn)
	
	lobbyIdBytes := make([]byte, 8)
	_, err       := reader.Read(lobbyIdBytes)
	if err != nil {
		panic(err)
	}
	lobbyId, _   := binary.Uvarint(lobbyIdBytes)

	logger.Println("Lobby", lobbyId, "chosen")
}

func joinFirstLobby(conn net.Conn) {
	lobbyList.Mutex.Lock()
	if len(lobbyList.Lobbies) > 0 {
		lobby := lobbyList.Lobbies[0]
		lobbyList.Lobbies = lobbyList.Lobbies[1:]
		logger.Println("Lobby joined:", lobby.Player.RemoteAddr(),
			"vs", conn.RemoteAddr())
		go executeGame(lobby.Player, conn)
	} else {
		writer := bufio.NewWriter(conn)
		err := writer.WriteByte(types.NO_LOBBIES)
		if err != nil {
			logger.Println(conn.RemoteAddr(), "disconnected")
			return
		}
		err = writer.Flush()
		if err != nil {
			logger.Println(conn.RemoteAddr(), "disconnected")
			return
		}
		go handleInitialConnect(conn)
	}
	lobbyList.Mutex.Unlock()
}
