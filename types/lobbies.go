package types

import (
	"net"
	"sync"
)

type Lobbies int

type Lobby struct {
	Player net.Conn
	Name   string
	Id     uint64
}

type MarshallableLobby struct {
	Id   uint64
	Name string
}

type SafeLobbyList struct {
	Lobbies []Lobby
	Mutex   sync.Mutex
}

func (sll *SafeLobbyList) RemoveLobbyWithId(id uint64) {
	sll.Mutex.Lock()
	for i, lobby := range sll.Lobbies {
		if lobby.Id == id {
			lobbies := append(sll.Lobbies[:i], sll.Lobbies[i+1:]...)
			sll.Lobbies = lobbies
			sll.Mutex.Unlock()
			return
		}
	}
	sll.Mutex.Unlock()
}

func (sll *SafeLobbyList) AddLobby(lobby Lobby) {
	sll.Mutex.Lock()
	sll.Lobbies = append(sll.Lobbies, lobby)
	sll.Mutex.Unlock()
}

func (sll *SafeLobbyList) ToMarshallable() []MarshallableLobby {
	result := []MarshallableLobby{}
	for _, lobby := range sll.Lobbies {
		marshallable := MarshallableLobby{
			Id   : lobby.Id,
			Name : lobby.Name,
		}
		result = append(result, marshallable)
	}
	return result
}
