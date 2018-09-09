package main

import "io"

// Disposition Once a command completes, this is what we should do about it
type Disposition int

const (
	// Continue receiving new commands
	Continue Disposition = iota + 1
	// Return and end current connection
	Return
	// Exit the server completely
	Exit
)

// Handler Is an interface to provide implementation for different line-based TCP commands.
type Handler interface {
	Handle(args []string, writer io.Writer) Disposition
}
