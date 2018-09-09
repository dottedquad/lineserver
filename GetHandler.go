package main

import (
	"io"
	"strconv"
)

type GetHandler struct {
	lineWriter LineWriter
}

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
