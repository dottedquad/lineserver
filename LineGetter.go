package main

import "io"

type LineWriter interface {
	WriteLine(lineNum int64, writer io.Writer) error
}
