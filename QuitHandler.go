package main

import (
	"io"
)

// QuitHandler Handles the QUIT command
type QuitHandler struct {
}

// Handle the QUIT Command
func (qc *QuitHandler) Handle(args []string, writer io.Writer) Disposition {

	return Return
}
