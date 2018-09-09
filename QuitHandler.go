package main

type QuitHandler struct {
}

func (qc *QuitHandler) Handle(args []string) (string, Disposition) {

	return "", Return
}
