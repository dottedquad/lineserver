package main

import (
	"fmt"
	"io"
)

//UnknownHandler handle unknown commands
type UnknownHandler struct {
}

// Handle Unknown commands
func (ec *UnknownHandler) Handle(args []string, writer io.Writer) Disposition {
	writer.Write([]byte(fmt.Sprintf("ERR Unknown command - %v\n", args[0])))
	return Return
}
