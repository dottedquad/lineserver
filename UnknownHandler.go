package main

type UnknownHandler struct {
}

func (ec *UnknownHandler) Handle(args []string) (string, Disposition) {
	return "Unknown command" + args[0], Exit
}
