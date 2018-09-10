package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {

	args := os.Args
	fmt.Println(args)

	// If I wanted to spend more time on this, I'd handle the args better:
	if len(args) < 3 {
		fmt.Printf("Not enough arguments")
	}

	if args[1] != "-p" {
		fmt.Printf("first arument must be -p")
	}

	port := args[2]

	// filename and stride arguments are optional and must be the 3rd and 4th arguments if present
	stride := 8
	filename := "john.txt"
	if len(args) >= 4 {
		filename = args[3]
	}
	if len(args) >= 5 {
		stridestr := args[4]
		var err error
		stride, err = strconv.Atoi(stridestr)
		if err != nil {
			fmt.Printf("Error while parsing stride\n")
			return
		}
	}

	// Listen for incoming connections.
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer listener.Close()
	fmt.Printf("Listening on 0.0.0.0:%v\n", port)

	commandDispatch := make(map[string]Handler)

	// Set up the dependencies
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open file: %v\n", filename)
		return
	}

	tmpfile, err := ioutil.TempFile("", "lineserver")
	if err != nil {
		fmt.Printf("Could not open temp file\n")
		return
	}

	defer os.Remove(tmpfile.Name()) // clean up the temp file

	commandDispatch["GET"] = &GetHandler{NewIndexedLineWriter(file, tmpfile, stride)}
	commandDispatch["QUIT"] = &QuitHandler{}
	commandDispatch["SHUTDOWN"] = &ShutdownHandler{}
	unknownHandler := &UnknownHandler{}

	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			// This will happen when the lister is closed by the Shutdown routine.
			break
		}

		// Handle connections in a new goroutine.
		go handleRequest(conn, commandDispatch, unknownHandler, listener)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, commandDispatch map[string]Handler, unknownHandler Handler, closer io.Closer) Disposition {

	reader := bufio.NewReader(conn)

	defer conn.Close()

	for {
		// read one line (ended with \n or \r\n)
		line, err := reader.ReadString('\n')
		// do something with data here, concat, handle and etc...
		commandargs := strings.Fields(line)

		if err != nil {
			fmt.Printf("Error reading: %v\n", err.Error())
			return Return
		}
		if len(commandargs) > 0 {
			disposition := Continue
			if val, ok := commandDispatch[commandargs[0]]; ok {
				disposition = val.Handle(commandargs, conn)
			} else {
				disposition = unknownHandler.Handle(commandargs, conn)
			}

			switch disposition {
			case Continue:
				// Do nothing
				break
			case Return:
				return disposition
			case Exit:
				closer.Close()
				return disposition
			}
		}
	}
}
