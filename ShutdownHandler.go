package main

import "io"

// ShutdownHandler Handle the SHUTDOWN command
type ShutdownHandler struct {
}

// Handle the SHUTDOWN command
func (ec *ShutdownHandler) Handle(args []string, writer io.Writer) Disposition {
	return Exit
}
