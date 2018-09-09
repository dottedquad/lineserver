package main

import "io"

type ShutdownHandler struct {
}

func (ec *ShutdownHandler) Handle(args []string, writer io.Writer) Disposition {
	return Exit
}
