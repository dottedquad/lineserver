package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type IndexedLineWriter struct {
	reader          io.ReadSeeker
	indexReadWriter io.ReadWriteSeeker
	indexStride     int
}

func NewIndexedLineWriter(reader io.ReadSeeker, indexReadWriter io.ReadWriteSeeker, indexStride int) *IndexedLineWriter {
	ilg := &IndexedLineWriter{reader, indexReadWriter, indexStride}
	ilg.createIndex()
	return ilg
}

func (ilg *IndexedLineWriter) createIndex() {
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
					fmt.Printf("Wrote %v to index %v for line %v\n", curPos, posbinary, curLine)
				} else {
					fmt.Printf("Skipping writing %v due to stride\n", curLine)
				}
			}
			curPos++
		}
	}
}

func (ilg *IndexedLineWriter) WriteLine(lineNum int64, writer io.Writer) error {

	if lineNum <= 0 {
		writer.Write([]byte("ERR\r\n"))
		return errors.New("Invalid Line Number")
	}

	indexPos := 8 * ((lineNum - 1) / int64(ilg.indexStride))
	fmt.Printf("indexPos %v\n", indexPos)
	ilg.indexReadWriter.Seek(indexPos, io.SeekStart)
	posbinary := make([]byte, 8)
	bytesread, err := ilg.indexReadWriter.Read(posbinary)
	fmt.Printf("bytesread %v\n", bytesread)
	if err != nil || bytesread != 8 {
		writer.Write([]byte("ERR\r\n"))
		return errors.New("Invalid Line Number")
	}
	filepos := int64(binary.LittleEndian.Uint64(posbinary))
	fmt.Printf("filepos %v\n", filepos)
	ilg.reader.Seek(filepos, io.SeekStart)
	//TODO share code with DirectFileReader?

	curLineNum := int64(0)
	writer.Write([]byte("OK\r\n"))
	buf := make([]byte, 1024)
	for {
		n, err := ilg.reader.Read(buf)
		bidx := 0
		eidx := 0
		done := false
		for i := 0; !done && i < n; i++ {
			if (lineNum-1)%int64(ilg.indexStride) == curLineNum {
				eidx = i + 1
			}
			if buf[i] == '\n' {
				if (lineNum-1)%int64(ilg.indexStride) == curLineNum {
					done = true
				}
				curLineNum++
				if (lineNum-1)%int64(ilg.indexStride) == curLineNum {
					bidx = i + 1
					eidx = i + 1
				}
			}
		}
		writer.Write(buf[bidx:eidx])
		if done {
			return nil
		}
		if err == io.EOF {
			break
		}
	}
	writer.Write([]byte("ERR\r\n"))
	return errors.New("Past end of file")
}
