package utils

import (
	"net"
	"sync"
)

type Lobby struct {
	Player net.Conn
	Id     uint64
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
