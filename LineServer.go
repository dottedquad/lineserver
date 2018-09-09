package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

// Disposition Once a command completes, this is what we should do about it
type Disposition int

const (
	// Continue receiving new commands
	Continue Disposition = 0
	// Return and end current connection
	Return Disposition = 1
	// Exit the server completely
	Exit Disposition = 2
)

// Command Is an interface to provide implementation for different line-based TCP commands.
type Command interface {
	Handle(args []string) Disposition
}

type GetCommand struct {
}

func (gc *GetCommand) Handle(args []string) Disposition {
	return Continue
}

type QuitCommand struct {
}

func (qc *QuitCommand) Handle(args []string) Disposition {

	return Return
}

type ShutdownCommand struct {
}

func (ec *ShutdownCommand) Handle(args []string) Disposition {
	return Exit
}

func main() {

	args := os.Args
	fmt.Println(args)

	// Listen for incoming connections.
	listener, err := net.Listen("tcp", CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer listener.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	commandDispatch := make(map[string]*Command)

	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, commandDispatch)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, commandDispatch map[string]*Command) Disposition {

	reader := bufio.NewReader(conn)

	defer conn.Close()

	for {
		// read one line (ended with \n or \r\n)
		line, err := reader.ReadString('\n')
		fmt.Printf("Line: %v\n", line)
		// do something with data here, concat, handle and etc...
		commandargs := strings.Fields(line)

		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		if len(commandargs) > 0 {
			disposition := Continue
			if val, ok := commandDispatch[commandargs[0]]; ok {
				disposition = val.Handle(commandargs)
			}

			switch disposition {
			case Continue:
				// Do nothing
				break
			case Return:
				// Return
				return
				break
			case Exit:
				// TODO Do something special
				return disposition
				break
			}
		}
	}

	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
}
