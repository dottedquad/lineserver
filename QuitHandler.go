package main

import (
	"io"
)

type QuitHandler struct {
}

func (qc *QuitHandler) Handle(args []string, writer io.Writer) Disposition {

	return Return
}
