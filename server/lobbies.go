package server

import (
	"net"
	"bufio"
	"strings"
	"encoding/gob"
	"encoding/binary"
	"github.com/daniellok/rps/types"
)

var lobbyList types.SafeLobbyList

func handleInitialConnect(conn net.Conn) {
	reader := bufio.NewReader(conn)

	b, err := reader.ReadByte()
	if err != nil {
		logger.Println(conn.RemoteAddr(), "disconnected")
		return
	}

	if b == types.CREATE_LOBBY {
		createLobby(conn)
	} else if b == types.JOIN_LOBBY {
		joinLobby(conn)
	} else {
		logger.Println("Someone is trolling")
		conn.Close()
	}
}

func createLobby(conn net.Conn) {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	
	err := writer.WriteByte(types.LOBBY_NAME)
	writer.Flush()
	if err != nil {
		logger.Println(conn.RemoteAddr(), "disconnected")
		return
	}
	
	name, err := reader.ReadString('\n')
	if err != nil {
		logger.Println(conn.RemoteAddr(), "disconnected")
		return
	}
	name = strings.TrimSuffix(name, "\n")
	logger.Println("Creating lobby {", currentMatchId, name, "}")

	matchIdMutex.Lock()
	id := currentMatchId
	currentMatchId += 1
	matchIdMutex.Unlock()
	
	lobby := types.Lobby{
		Player : conn,
		Name : name,
		Id : id,
	}

	lobbyList.AddLobby(lobby)

	err = writer.WriteByte(types.LOBBY_CREATED)
	writer.Flush()
	if err != nil {
		logger.Println(conn.RemoteAddr(), "disconnected")
		lobbyList.RemoveLobbyWithId(id)
		return
	}

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
	reader := bufio.NewReader(conn)
	
	lobbyIdBytes := make([]byte, 8)
	_, err       := reader.Read(lobbyIdBytes)
	if err != nil {
		logger.Println(err)
		return
	}
	lobbyId, _ := binary.Uvarint(lobbyIdBytes)
	lobby, err := lobbyList.RemoveLobbyWithId(lobbyId)
	if err != nil {
		logger.Println("Lobby", lobbyId, "is no longer available")
	}

	executeGame(lobby.Player, conn)
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
