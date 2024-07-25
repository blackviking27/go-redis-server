package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Defining const identifiers
const (
	STRING  = "+"
	ERROR   = "-"
	INTEGER = ":"
	BULK    = "$"
	ARRAY   = "*"
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

// Method to read a line from input
func (r *Resp) readLine() (line []byte, n int, err error) {
	// reading individual bytes until we reach the CRLF (\r\n) which is the end of line
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1                 // increasing the count of bytes read
		line = append(line, b) // appending the read bytes to slice

		// when we reach the \n byte we check weather the prev byte is \r or not.
		// if \r is present then we know that we have reached the end of line
		// and can break the loop
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

// Method to read integer from input
func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return int(i64), n, nil
}

// Method to read the input array as whole
func (r *Resp) Read() (Value, error) {
	// reading the first byte of the string
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch string(_type) {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Println("Unknown Type %v", string(_type))
		return Value{}, nil
	}
}

// Method to parse array
// func (r *Resp) readArray() (Value, error) {
// 	v := Value{}
// 	v.typ = "array"

// }

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}
