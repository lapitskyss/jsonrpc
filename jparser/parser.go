package jparser

import (
	"bytes"
	"errors"
)

var (
	ErrParseJSON               = errors.New("parse error")
	ErrIncorrectFieldType      = errors.New("incorrect field type")
	UnknownValueTypeError      = errors.New("unknown value type")
	MalformedStringError       = errors.New("value is string, but can't find closing '\"' symbol")
	MalformedArrayError        = errors.New("value is array, but can't find closing ']' symbol")
	MalformedObjectError       = errors.New("value looks like object, but can't find closing '}' symbol")
	MalformedValueError        = errors.New("value looks like Number/Boolean/None, but can't find its end: ',' or '}' symbol")
	MalformedStringEscapeError = errors.New("encountered an invalid escape sequence in a string")
)

// How much stack space to allocate for unescaping JSON strings; if a string longer
// than this needs to be escaped, it will result in a heap allocation
const unescapeStackBufSize = 64

// ValueType Data types available in valid JSON data.
type ValueType int

const (
	NotExist = ValueType(iota)
	String
	Number
	Object
	Array
	Boolean
	Null
	Unknown
)

func (vt ValueType) String() string {
	switch vt {
	case NotExist:
		return "non-existent"
	case String:
		return "string"
	case Number:
		return "number"
	case Object:
		return "object"
	case Array:
		return "array"
	case Boolean:
		return "boolean"
	case Null:
		return "null"
	default:
		return "unknown"
	}
}

var (
	trueLiteral  = []byte("true")
	falseLiteral = []byte("false")
	nullLiteral  = []byte("null")
)

func tokenEnd(data []byte) int {
	for i, c := range data {
		switch c {
		case ' ', '\n', '\r', '\t', ',', '}', ']':
			return i
		}
	}

	return len(data)
}

// Find position of next character which is not whitespace
func nextToken(data []byte) int {
	for i, c := range data {
		switch c {
		case ' ', '\n', '\r', '\t':
			continue
		default:
			return i
		}
	}

	return -1
}

// Tries to find the end of string
// Support if string contains escaped quote symbols.
func stringEnd(data []byte) (int, bool) {
	escaped := false
	for i, c := range data {
		if c == '"' {
			if !escaped {
				return i + 1, false
			} else {
				j := i - 1
				for {
					if j < 0 || data[j] != '\\' {
						return i + 1, true // even number of backslashes
					}
					j--
					if j < 0 || data[j] != '\\' {
						break // odd number of backslashes
					}
					j--

				}
			}
		} else if c == '\\' {
			escaped = true
		}
	}

	return -1, escaped
}

// Find end of the data structure, array or object.
// For array openSym and closeSym will be '[' and ']', for object '{' and '}'
func blockEnd(data []byte, openSym byte, closeSym byte) int {
	level := 0
	i := 0
	ln := len(data)

	for i < ln {
		switch data[i] {
		case '"': // If inside string, skip it
			se, _ := stringEnd(data[i+1:])
			if se == -1 {
				return -1
			}
			i += se
		case openSym: // If open symbol, increase level
			level++
		case closeSym: // If close symbol, increase level
			level--

			// If we have returned to the original level, we're done
			if level == 0 {
				return i + 1
			}
		}
		i++
	}

	return -1
}

func getType(data []byte) ([]byte, ValueType, int, error) {
	var dataType ValueType
	endOffset := 0

	// if string value
	if data[0] == '"' {
		dataType = String
		if idx, _ := stringEnd(data[1:]); idx != -1 {
			endOffset += idx + 1
		} else {
			return nil, dataType, 0, MalformedStringError
		}
	} else if data[0] == '[' { // if array value
		dataType = Array
		// break label, for stopping nested loops
		endOffset = blockEnd(data[0:], '[', ']')

		if endOffset == -1 {
			return nil, dataType, 0, MalformedArrayError
		}
	} else if data[0] == '{' { // if object value
		dataType = Object
		// break label, for stopping nested loops
		endOffset = blockEnd(data[0:], '{', '}')

		if endOffset == -1 {
			return nil, dataType, 0, MalformedObjectError
		}
	} else {
		// Number, Boolean or None
		end := tokenEnd(data[endOffset:])

		if end == -1 {
			return nil, dataType, 0, MalformedValueError
		}

		value := data[:endOffset+end]

		switch data[0] {
		case 't', 'f': // true or false
			if bytes.Equal(value, trueLiteral) || bytes.Equal(value, falseLiteral) {
				dataType = Boolean
			} else {
				return nil, Unknown, 0, UnknownValueTypeError
			}
		case 'u', 'n': // undefined or null
			if bytes.Equal(value, nullLiteral) {
				dataType = Null
			} else {
				return nil, Unknown, 0, UnknownValueTypeError
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			dataType = Number
		default:
			return nil, Unknown, 0, UnknownValueTypeError
		}

		endOffset += end
	}
	return data[:endOffset], dataType, endOffset, nil
}

func ArrayLength(data []byte) int {
	i := 0
	ln := len(data)
	arrLen := 0

	for i < ln {
		switch data[i] {
		case ',', '[', ']', ' ', '\n', '\r', '\t':
			i++
			continue
		case '{':
			endOffset := blockEnd(data[i:], '{', '}')
			if endOffset == -1 {
				return 0
			}

			arrLen++
			i += endOffset

		default:
			return 0
		}
	}

	return arrLen
}

func ArrayElement(data []byte, index int) []byte {
	i := 0
	ln := len(data)
	arrLen := 0

	for i < ln {
		switch data[i] {
		case ',', '[', ']', ' ', '\n', '\r', '\t':
			i++
			continue
		case '{':
			endOffset := blockEnd(data[i:], '{', '}')
			if endOffset == -1 {
				return nil
			}

			if arrLen == index {
				return data[i : i+endOffset]
			}

			arrLen++
			i += endOffset

		default:
			return nil
		}
	}

	return nil
}

// IsArray check is provided json is array
func IsArray(data []byte) bool {
	i := 0
	ln := len(data)

	for i < ln {
		switch data[i] {
		case ' ', '\n', '\r', '\t':
			i++
			continue
		case '[':
			return true
		default:
			return false
		}
	}

	return false
}
