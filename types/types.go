package types

import "net"

// moves
const (
	ROCK = iota
	PAPER
	SCISSORS
	DISCONNECT
)

type Move struct {
	Move   byte
	Player net.Conn
}

type Result struct {
	Result byte
	Player net.Conn
}
	
// client internal
const (
	SCAN_ERROR = iota
)

// client -> server choices
const (
	CREATE_LOBBY = iota
	JOIN_LOBBY
	RECEIVED_MATCH
)

// server -> client messages
const (
	DRAW = iota
	WIN
	LOSE
	LOBBY_CREATED
	LOBBY_JOINED
	NO_LOBBIES
	GAME_START
	INVALID_CHOICE
	OPPONENT_DC
	SERVER_DC
)

