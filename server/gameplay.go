package server

import (
	"net"
	"bufio"
	"github.com/daniellok/rps/types"
)

func executeGame(player1 net.Conn, player2 net.Conn) {
	setupGame(player1, player2)
	move1, move2 := collectMoves(player1, player2)
	if move1.Move == types.DISCONNECT {
		disconnector := move1.Player
		otherPlayer  := player1
		if otherPlayer == disconnector {
			otherPlayer = player2
		}
		handleDisconnects(disconnector, otherPlayer)
		return
	}
	
	res1, res2 := calcResult(move1, move2)
	logger.Println("Results calculated")
	sendResults(res1, res2)
}

func setupGame(player1 net.Conn, player2 net.Conn) {
	var err error
	
	reader2 := bufio.NewReader(player2)
	writer1 := bufio.NewWriter(player1)	
	writer2 := bufio.NewWriter(player2)

	err = writer1.WriteByte(types.LOBBY_JOINED)
	writer1.Flush()
	if err != nil {
		logger.Println(player1.RemoteAddr(), "disconnected")
	}
	err = writer2.WriteByte(types.LOBBY_JOINED)
	if err != nil {
		logger.Println(player2.RemoteAddr(), "disconnected")
	}
	writer2.Flush()
	b, err := reader2.ReadByte()
	if err != nil || b != types.RECEIVED_MATCH {
		logger.Println(player2.RemoteAddr(), "disconnected")
	}
	logger.Println("acknowledgement received")
	
	err = writer1.WriteByte(types.GAME_START)
	writer1.Flush()
	err = writer2.WriteByte(types.GAME_START)
	writer2.Flush()
}

// TODO the problem is in here
func collectMoves(player1 net.Conn, player2 net.Conn) (types.Move, types.Move) {
	ch := make(chan types.Move)
	go waitForMove(player1, ch)
	go waitForMove(player2, ch)

	move1 := <-ch
	if move1.Move == types.DISCONNECT {
		move2 := types.Move{
			Move   : types.ROCK,
			Player : player1,
		}
		if move1.Player == player1 {
			move2.Player = player2
		}
		return move1, move2
	}
	move2 := <-ch
	if move2.Move == types.DISCONNECT {
		return move2, move1
	}
	return move1, move2
}

func handleDisconnects(disconnector net.Conn, otherPlayer net.Conn) {
	logger.Println(disconnector.RemoteAddr(), "has disconnected")
	writer := bufio.NewWriter(otherPlayer)
	
	writer.WriteByte(types.OPPONENT_DC)
	writer.Flush()

	disconnector.Close()
	otherPlayer.Close()
}

func waitForMove(conn net.Conn, out chan types.Move) {
	reader := bufio.NewReader(conn)
	b, err := reader.ReadByte()
	if err != nil {
		logger.Println(conn.RemoteAddr(), "disconnected")
		b = types.DISCONNECT
	}
	move := types.Move{
		Move   : b,
		Player : conn,
	}
	out <- move

	b, err = reader.ReadByte()
	if err != nil {
		move.Move = types.DISCONNECT
	}
	out <- move
}

func calcResult(move1 types.Move, move2 types.Move) (types.Result, types.Result) {
	result  := int(move1.Move) - int(move2.Move)
	result1 := types.Result{
		Result : types.DRAW,
		Player : move1.Player,
	}
	result2 := types.Result{
		Result : types.DRAW,
		Player : move2.Player,
	}
	if result == -1 || result == 2 {
		result1.Result = types.LOSE
		result2.Result = types.WIN
	} else if result == 1 || result == -2 {
		result1.Result = types.WIN
		result2.Result = types.LOSE
	}
	return result1, result2
}

func sendResults(result1 types.Result, result2 types.Result) {
	writer1 := bufio.NewWriter(result1.Player)
	writer2 := bufio.NewWriter(result2.Player)
	
	writer1.WriteByte(result1.Result)
	writer1.Flush()
	writer2.WriteByte(result2.Result)
	writer2.Flush()

	logger.Println("Results sent")
	result1.Player.Close()
	result2.Player.Close()
}
