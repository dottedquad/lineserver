package main

import "strconv"

type GetHandler struct {
	lineGetter LineGetter
}

func (gc *GetHandler) Handle(args []string) (string, Disposition) {

	if len(args) != 2 {
		return "ERR", Continue
	}
	lineNum, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return "ERR", Continue
	}
	line, err := gc.lineGetter.GetLine(lineNum)
	if err == nil {
		return line, Continue
	} else {
		return "ERR", Continue
	}
}
