package main

import (
	"bufio"
	"errors"
	"io"
)

type DirectLineGetter struct {
	reader io.ReadSeeker
}

func (dlg *DirectLineGetter) GetLine(lineNum int64) (string, error) {

	if lineNum <= 0 {
		return "", errors.New("Invalid Line Number")
	}
	dlg.reader.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(dlg.reader)
	curLineNum := int64(0)
	for scanner.Scan() {
		curLineNum++
		if lineNum == curLineNum {
			return scanner.Text(), scanner.Err()
		}
	}
	return "", errors.New("Past end of file")
}
