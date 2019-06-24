package client

import (
	"os"
	"fmt"
	"github.com/daniellok/rps/types"
)

func playGame() {
	fmt.Print("What move would you like to make?\n  [1] Rock 🤜\n  [2] Paper ✋\n  [3] Scissors ✌️\n> ")

	choiceChannel := make(chan int)
	serverChannel := make(chan byte)

	go listenForInput(choiceChannel)
	go listenFromServer(serverChannel)

	select {
	case choice := <-choiceChannel:
		switch choice {
		case 1:
			writer.WriteByte(types.ROCK)
			writer.Flush()
		case 2:
			writer.WriteByte(types.PAPER)
			writer.Flush()
		case 3:
			writer.WriteByte(types.SCISSORS)
			writer.Flush()
		}
	case b := <-serverChannel:
		if b == types.OPPONENT_DC {
			fmt.Println("Your opponent has disconnected")
			os.Exit(1)
		} else if b == types.SERVER_DC {
			fmt.Println("Lost connection to the server.")
			os.Exit(1)
		}
	}

	fmt.Println("Move sent! Waiting for your opponent...")
	
	// there's no way to listen from server and listen for input
	// at the same time in one process, so all listening is
	// transferred to the `listenFromServer` goroutine.
	waitForResult(serverChannel)
}

func listenFromServer(ch chan byte) {
	b, err := reader.ReadByte()
	if err != nil {
		ch <- types.SERVER_DC
	}
	ch <- b
}

func listenForInput(ch chan int) {
	var choice int
	for {
		_, err := fmt.Scan(&choice)
		if err != nil {
			ch <- types.SCAN_ERROR
			return
		}
		if choice > 0 && choice <= 3 {
			break
		} else {
			fmt.Print("Please choose properly\n> ")
		}
	}
	ch <- choice
}

func waitForResult(serverChannel chan byte) {
	result := <-serverChannel

	switch result {
	case types.DRAW:
		fmt.Println("It was a draw!")
	case types.WIN:
		fmt.Println("You won! Congratulations.")
	case types.LOSE:
		fmt.Println("You lost :(. Better luck next time.")
	case types.OPPONENT_DC:
		fmt.Println("Your opponent disconnected.")
	case types.SERVER_DC:
		fmt.Println("Lost connection to sever.")
	default:
		fmt.Println("Whoops! Something went wrong.")
	}
}
