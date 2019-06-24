package client

import (
	"os"
	"fmt"
	"encoding/gob"
	"github.com/daniellok/rps/types"
)

func createLobby() {
	writer.WriteByte(types.CREATE_LOBBY)
	writer.Flush()
	fmt.Println("Waiting for response from server...")
	
	b, err := reader.ReadByte()
	handleError(err)
	
	if b == types.LOBBY_CREATED {
		fmt.Println("Lobby created! Waiting for someone to join...")
	} else {
		fmt.Println("Something went wrong!")
		os.Exit(1)
	}

	waitForGameStart()
}

func joinLobby() {
	writer.WriteByte(types.JOIN_LOBBY)
	writer.Flush()

	fmt.Println("Waiting to join a lobby...")
	waitForGameStart()
}

func retrieveLobbies() {
	var lobbies []types.Lobby
	dec := gob.NewDecoder(conn)
	
	
}

func waitForGameStart() {
	b, err := reader.ReadByte()
	handleError(err)
	
	if b == types.LOBBY_JOINED {
		writer.WriteByte(types.RECEIVED_MATCH)
		writer.Flush()

		b, err = reader.ReadByte()
		handleError(err)
		
		fmt.Println("Game start!")
		playGame()
	} else if b == types.NO_LOBBIES {
		fmt.Println("No lobbies :(.")
		prompt()
	}
}
