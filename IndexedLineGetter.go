package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type IndexedLineGetter struct {
	reader          io.ReadSeeker
	indexReadWriter io.ReadWriteSeeker
	indexStride     int
}

func NewIndexedLineGetter(reader io.ReadSeeker, indexReadWriter io.ReadWriteSeeker, indexStride int) *IndexedLineGetter {
	ilg := &IndexedLineGetter{reader, indexReadWriter, indexStride}
	ilg.createIndex()
	return ilg
}

func (ilg *IndexedLineGetter) createIndex() {
	ilg.reader.Seek(0, io.SeekStart)
	//scanner := bufio.NewScanner(ilg.reader)
	buf := make([]byte, 1024)
	var err error = nil
	curPos := int64(0)
	curLine := int64(0)
	nextIsBeginning := true
	for err == nil {
		var bytesread int = 0
		bytesread, err = ilg.reader.Read(buf)
		for i := 0; i < bytesread; i++ {
			if buf[i] == '\n' {
				nextIsBeginning = true
				curLine++
			} else if nextIsBeginning {
				nextIsBeginning = false
				if curLine%int64(ilg.indexStride) == 0 {
					posbinary := make([]byte, 8)
					binary.LittleEndian.PutUint64(posbinary, uint64(curPos))
					ilg.indexReadWriter.Write(posbinary)
					fmt.Printf("Wrote %v to index %v for line %v", curPos, posbinary, curLine)
				} else {
					fmt.Printf("Skipping writing %v due to stride", curLine)
				}
			}
			curPos++
		}
	}
}

func (ilg *IndexedLineGetter) GetLine(lineNum int64) (string, error) {

	if lineNum <= 0 {
		return "", errors.New("Invalid Line Number")
	}

	indexPos := 8 * ((lineNum - 1) / int64(ilg.indexStride))
	fmt.Printf("indexPos %v\n", indexPos)
	ilg.indexReadWriter.Seek(indexPos, io.SeekStart)
	posbinary := make([]byte, 8)
	bytesread, err := ilg.indexReadWriter.Read(posbinary)
	fmt.Printf("bytesread %v\n", bytesread)
	if err != nil || bytesread != 8 {
		return "", errors.New("Invalid Line Number")
	}
	filepos := int64(binary.LittleEndian.Uint64(posbinary))
	fmt.Printf("filepos %v\n", filepos)
	ilg.reader.Seek(filepos, io.SeekStart)
	//TODO share code with DirectFileReader?
	scanner := bufio.NewScanner(ilg.reader)
	curLineNum := int64(0)
	for scanner.Scan() {
		if (lineNum-1)%int64(ilg.indexStride) == curLineNum {
			fmt.Printf("(lineNum%v-1)%%int64(ilg.indexStride%v) == curLineNum%v TRUE\n", lineNum, ilg.indexStride, curLineNum)
			return scanner.Text(), scanner.Err()
		} else {
			fmt.Printf("(lineNum%v-1)%%int64(ilg.indexStride%v) == curLineNum%v FALSE\n", lineNum, ilg.indexStride, curLineNum)

		}
		curLineNum++
	}
	return "", errors.New("Past end of file")
}
