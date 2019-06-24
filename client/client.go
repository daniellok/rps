package client

import (
	"os"
	"log"
	"fmt"
	"net"
	"bufio"
)

var conn net.Conn
var reader *bufio.Reader
var writer *bufio.Writer
var logger = log.New(os.Stdin, "RPS client | ", log.Ldate|log.Ltime)

func handleError(err error) {
	if err != nil {
		fmt.Println("Error occured!", err)
		os.Exit(1)
	}
}

func RunClient() {
	var err error
	
	conn, err = net.Dial("tcp", "localhost:8787")
	handleError(err)

	reader = bufio.NewReader(conn)
	writer = bufio.NewWriter(conn)
	
	prompt()
}

func prompt() {
	var choice int

	fmt.Println("What would you like to do?\n  [1] Create a lobby\n  [2] Join a lobby")
	
	for {
		fmt.Print("> ")
		_, err := fmt.Scan(&choice)
		handleError(err)
		
		if choice == 1 {
			createLobby()
			break
		} else if choice == 2 {
			joinLobby()
			break
		} else {
			fmt.Println("Please choose properly")
		}
	}
}
