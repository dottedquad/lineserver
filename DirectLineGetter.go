package main

import (
	"bufio"
	"errors"
	"io"
)

type DirectLineGetter struct {
	reader io.ReadSeeker
}

func (dlg *DirectLineGetter) GetLine(lineNum int64, writer io.Writer) error {

	if lineNum <= 0 {
		return errors.New("Invalid Line Number")
	}
	dlg.reader.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(dlg.reader)
	curLineNum := int64(0)
	for scanner.Scan() {
		curLineNum++
		if lineNum == curLineNum {
			writer.Write([]byte(scanner.Text()))
			return scanner.Err()
		}
	}
	return errors.New("Past end of file")
}
