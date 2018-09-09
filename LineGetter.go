package main

type LineGetter interface {
	GetLine(lineNum int64) (string, error)
}
