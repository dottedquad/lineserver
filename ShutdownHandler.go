package main

type ShutdownHandler struct {
}

func (ec *ShutdownHandler) Handle(args []string) (string, Disposition) {
	return "", Exit
}
