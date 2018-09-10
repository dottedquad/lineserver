package main

import "io"

//LineWriter Provides a method that will write the lineNum'th line (1 based) into the writer
type LineWriter interface {
	WriteLine(lineNum int64, writer io.Writer) error
}
