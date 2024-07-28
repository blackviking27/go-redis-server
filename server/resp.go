package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Defining const identifiers
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// Value struct - convert resp to struct and vice versa
type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

// Marhsal the data to send to the user
// As part of the value struct
func (v Value) Marshall() []byte {
	switch v.typ {
	case "array":
		return v.marshallArray()
	case "bulk":
		return v.marshallBulk()
	case "string":
		return v.marshallString()
	case "null":
		return v.marshallNull()
	case "error":
		return v.marshallError()
	default:
		return []byte{}
	}
}

// Marshalling simple string
func (v Value) marshallString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// Marshalling Bulk string
func (v Value) marshallBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// Marshalling array
func (v Value) marshallArray() []byte {
	len := len(v.array)
	var bytes []byte

	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Marshall()...)
	}

	return bytes
}

// Marhsalling error
func (v Value) marshallError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// Marshalling Null
func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}

// RESP strucutre - contrcut RESP objects
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

	switch _type {
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
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// read the length of the array
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// parsing each line and reading the value
	v.array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		// append the value to array
		v.array = append(v.array, val)
	}
	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, nil
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	// read the trailing crlf
	r.readLine()

	return v, nil

}

// Returns the RESP reader object - used to read RESP commands
func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// Writer struct - use to send RESP response
type Writer struct {
	writer io.Writer
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshall()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// Return a RESP response object - use to send RESP response
func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}
