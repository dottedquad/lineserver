package main

import (
	"encoding/binary"
	"errors"
	"io"
)

// IndexedLineWriter Implementation of LineWriter that uses an index file
type IndexedLineWriter struct {
	reader          io.ReadSeeker
	indexReadWriter io.ReadWriteSeeker
	indexStride     int64
}

// NewIndexedLineWriter Factory function
func NewIndexedLineWriter(reader io.ReadSeeker, indexReadWriter io.ReadWriteSeeker, indexStride int) *IndexedLineWriter {
	ilg := &IndexedLineWriter{reader, indexReadWriter, int64(indexStride)}
	ilg.createIndex()
	return ilg
}

// Create an index file that holds the byte position of every ilg.indexStride'th line
func (ilg *IndexedLineWriter) createIndex() {
	ilg.reader.Seek(0, io.SeekStart)
	buf := make([]byte, 1024)
	var err error
	curPos := int64(0)
	curLine := int64(0)
	nextIsBeginning := true
	for err == nil {
		var bytesread int
		bytesread, err = ilg.reader.Read(buf)
		for i := 0; i < bytesread; i++ {
			if buf[i] == '\n' {
				nextIsBeginning = true
				curLine++
			} else if nextIsBeginning {
				nextIsBeginning = false
				if curLine%ilg.indexStride == 0 {
					posbinary := make([]byte, 8)
					binary.LittleEndian.PutUint64(posbinary, uint64(curPos))
					ilg.indexReadWriter.Write(posbinary)

				}
			}
			curPos++
		}
	}
}

// WriteLine write a single line at lineNum from the reader into the writer
func (ilg *IndexedLineWriter) WriteLine(lineNum int64, writer io.Writer) error {

	if lineNum <= 0 {
		writer.Write([]byte("ERR\r\n"))
		return errors.New("Invalid Line Number")
	}

	indexPos := 8 * ((lineNum - 1) / ilg.indexStride)

	ilg.indexReadWriter.Seek(indexPos, io.SeekStart)
	posbinary := make([]byte, 8)
	bytesread, err := ilg.indexReadWriter.Read(posbinary)

	if err != nil || bytesread != 8 {
		writer.Write([]byte("ERR\r\n"))
		return errors.New("Invalid Line Number")
	}
	filepos := int64(binary.LittleEndian.Uint64(posbinary))

	ilg.reader.Seek(filepos, io.SeekStart)
	//TODO share code with DirectFileReader?

	// How many lines from the lookup position to find the desired line
	lineOffsetNum := (lineNum - 1) % ilg.indexStride

	curLineNum := int64(0)
	writer.Write([]byte("OK\r\n"))
	buf := make([]byte, 1024)
	for {
		n, err := ilg.reader.Read(buf)
		bidx := 0
		eidx := 0
		done := false
		for i := 0; !done && i < n; i++ {
			if lineOffsetNum == curLineNum {
				eidx = i + 1
			}
			if buf[i] == '\n' {
				if lineOffsetNum == curLineNum {
					done = true
				}
				curLineNum++
				if lineOffsetNum == curLineNum {
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
