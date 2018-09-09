package main

import (
	"fmt"
	"io"
)

type UnknownHandler struct {
}

func (ec *UnknownHandler) Handle(args []string, writer io.Writer) Disposition {
	writer.Write([]byte(fmt.Sprintf("ERR Unknown command - %v\n", args[0])))
	return Return
}
