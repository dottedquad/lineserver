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

	// This is a little confusing, please read the comments
	for err == nil {

		// Read the intput file 1024 bytes at a time, and find all the newline character
		var bytesread int
		bytesread, err = ilg.reader.Read(buf)
		for i := 0; i < bytesread; i++ {
			if buf[i] == '\n' {
				// When we find a newline, that means that the next character is the beginning of a line, and that the current line increases
				nextIsBeginning = true
				curLine++
			} else if nextIsBeginning {
				// If nextIsBeginning is set, it means that this character is the start of a new line. Reset the flag.
				nextIsBeginning = false
				if curLine%ilg.indexStride == 0 {
					// See if we need a write an entry to the index based on the stride and the line number.
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

	// Find the closest line before in the index
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

	// How many lines from the lookup position to find the desired line
	lineOffsetNum := (lineNum - 1) % ilg.indexStride

	curLineNum := int64(0)

	// I don't like this OK being here, but I couldn't figure out how to move it out without requiring
	// copying the entire line output buffer rather than the more efficient writing directly to the listener.
	// Since only this layer knows if we should write OK or ERR
	writer.Write([]byte("OK\r\n"))
	buf := make([]byte, 1024)

	// Read in the file 1024 bytes at a time from the closest indexed point.
	for {
		n, err := ilg.reader.Read(buf)
		bidx := 0
		eidx := 0
		done := false
		// Look through the buffer and try to find the line we're looking for (LineOffsetNum from the indexed line
		// Set the bidx and eidx indices to mark which part of this buffer needs to be copied
		for i := 0; !done && i < n; i++ {
			if lineOffsetNum == curLineNum {
				// Include one more character from the buffer
				eidx = i + 1
			}
			if buf[i] == '\n' {
				if lineOffsetNum == curLineNum {
					// This is the newline for the line we're writing. We can be done now.
					done = true
				}
				curLineNum++
				if lineOffsetNum == curLineNum {
					// Begin the buffer here if the conditions are right
					bidx = i + 1
					eidx = i + 1
				}
			}
		}
		// Write the appropriate slice of the buffer to the listener
		writer.Write(buf[bidx:eidx])
		if done {
			return nil
		}
		if err == io.EOF {
			// We hit an unexpected EOF.
			break
		}
	}
	writer.Write([]byte("ERR\r\n"))
	return errors.New("Past end of file")
}
