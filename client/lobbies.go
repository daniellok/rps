package client

import (
	"os"
	"fmt"
	"encoding/gob"
	"encoding/binary"
	"github.com/daniellok/rps/types"
)

func createLobby() {
	var name string
	fmt.Print("What do you want to name your lobby?\n> ")
	_, err := fmt.Scanln(&name)
	handleError(err)
	
	writer.WriteByte(types.CREATE_LOBBY)
	writer.Flush()

	b, err := reader.ReadByte()
	handleError(err)
	if b != types.LOBBY_NAME {
		fmt.Println("Something went wrong!")
		os.Exit(1)
	}

	_, err = writer.WriteString(name + "\n")
	writer.Flush()
	handleError(err)
	
	b, err = reader.ReadByte()
	handleError(err)
	if b != types.LOBBY_CREATED {
		fmt.Println("Something went wrong!")
		os.Exit(1)
	}

	fmt.Println("lobby created!")
	waitForGameStart()
}

func joinLobby() {
	writer.WriteByte(types.JOIN_LOBBY)
	writer.Flush()

	lobbies := retrieveLobbies()
	fmt.Println(lobbies)

	var id uint64
	fmt.Print("Which lobby ID do you want to join?\n> ")
	_, err := fmt.Scan(&id)
	idBytes := make([]byte, 8)
	binary.PutUvarint(idBytes, id)
	_, err = writer.Write(idBytes)
	writer.Flush()
	handleError(err)

	fmt.Println("Waiting to join a lobby...")
	waitForGameStart()
}

func retrieveLobbies() []types.MarshallableLobby {
	lobbies := []types.MarshallableLobby{}

	decoder := gob.NewDecoder(conn)
	err     := decoder.Decode(&lobbies)
	handleError(err)

	return lobbies
}

func waitForGameStart() {
	fmt.Println("Waiting for game start")
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
