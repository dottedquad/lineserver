package main

import (
	"io"
	"strconv"
)

// GetHandler Handle the GET command
type GetHandler struct {
	lineWriter LineWriter
}

// Handle the GET command
func (gc *GetHandler) Handle(args []string, writer io.Writer) Disposition {

	if len(args) != 2 {
		writer.Write([]byte("ERR\r\n"))
		return Continue
	}
	lineNum, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		writer.Write([]byte("ERR\r\n"))
		return Continue
	}
	gc.lineWriter.WriteLine(lineNum, writer)
	return Continue
}
