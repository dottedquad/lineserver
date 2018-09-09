package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

const (
	ConnHost = "localhost"
	ConnPort = "3333"
	ConnType = "tcp"
)

func main() {

	args := os.Args
	fmt.Println(args)

	// Listen for incoming connections.
	listener, err := net.Listen("tcp", ConnHost+":"+ConnPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer listener.Close()
	fmt.Println("Listening on " + ConnHost + ":" + ConnPort)

	commandDispatch := make(map[string]Handler)

	// Set up the dependencies
	file, err := os.Open("john.txt")
	if err != nil {
		fmt.Errorf("Could not open file")
		return
	}

	tmpfile, err := ioutil.TempFile("", "lineserver")
	if err != nil {
		fmt.Errorf("Could not open temp file")
		return
	}

	defer os.Remove(tmpfile.Name()) // clean up the temp file

	commandDispatch["GET"] = &GetHandler{NewIndexedLineWriter(file, tmpfile, 4)}
	commandDispatch["QUIT"] = &QuitHandler{}
	commandDispatch["SHUTDOWN"] = &ShutdownHandler{}

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
func handleRequest(conn net.Conn, commandDispatch map[string]Handler) Disposition {

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
			return Return
		}
		if len(commandargs) > 0 {
			disposition := Continue
			if val, ok := commandDispatch[commandargs[0]]; ok {
				disposition = val.Handle(commandargs, conn)
			} else {
				disposition = (&UnknownHandler{}).Handle(commandargs, conn)
			}

			switch disposition {
			case Continue:
				// Do nothing
				break
			case Return:
				// Return
				return disposition
			case Exit:
				// TODO Do something special
				return disposition
			}
		}
	}
}
