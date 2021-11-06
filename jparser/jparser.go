package jparser

type JParser struct {
	data []byte

	ID     []byte
	IDType ValueType

	Version     []byte
	VersionType ValueType

	Method     []byte
	MethodType ValueType

	Params     []byte
	ParamsType ValueType

	err error
}

// Parse jsonrpc request data.
func Parse(data []byte) *JParser {
	jParser := &JParser{
		data: data,
	}

	jParser.parse()

	return jParser
}

// GetId get jsonrpc id as string.
func (j *JParser) GetId() string {
	if j.IDType == String {
		id := j.ID[1 : len(j.ID)-1]
		return string(id)
	}
	return string(j.ID)
}

// GetVersion get jsonrpc version as string.
func (j *JParser) GetVersion() string {
	return string(j.Version)
}

// GetMethod get jsonrpc method as string.
func (j *JParser) GetMethod() string {
	return string(j.Method)
}

// GetMethod return error from parsing.
func (j *JParser) Error() error {
	return j.err
}

// parse jsonrpc.
func (j *JParser) parse() {
	j.IDType = NotExist
	j.VersionType = NotExist
	j.MethodType = NotExist
	j.ParamsType = NotExist

	level := 0
	i := 0
	ln := len(j.data)
	var stackbuf [unescapeStackBufSize]byte

	for i < ln {
		switch j.data[i] {
		case '"':
			i++
			keyBegin := i

			strEnd, keyEscaped := stringEnd(j.data[i:])
			if strEnd == -1 {
				j.err = ErrParseJSON
				return
			}
			i += strEnd
			keyEnd := i - 1

			valueOffset := nextToken(j.data[i:])
			if valueOffset == -1 {
				j.err = ErrParseJSON
				return
			}

			i += valueOffset

			// if string is a key
			if j.data[i] == ':' {
				if level < 1 {
					j.err = ErrParseJSON
					return
				}

				key := j.data[keyBegin:keyEnd]

				// for unescape: if there are no escape sequences, this is cheap; if there are, it is a
				// bit more expensive, but causes no allocations unless len(key) > unescapeStackBufSize
				var keyUnesc []byte
				if !keyEscaped {
					keyUnesc = key
				} else if ku, err := Unescape(key, stackbuf[:]); err != nil {
					j.err = ErrParseJSON
					return
				} else {
					keyUnesc = ku
				}

				if level == 1 {
					i++

					nt := nextToken(j.data[i:])
					if nt == -1 {
						j.err = ErrParseJSON
						return
					}

					i += nt
					value, dataType, endOffset, err := getType(j.data[i:])
					if err != nil {
						j.err = ErrParseJSON
						return
					}

					if string(keyUnesc) == "id" {
						if dataType == String || dataType == Number || dataType == Null {
							j.ID = value
							j.IDType = dataType
						} else {
							j.err = ErrIncorrectFieldType
							return
						}
					} else if string(keyUnesc) == "jsonrpc" {
						if dataType != String {
							j.err = ErrIncorrectFieldType
							return
						}
						j.Version = value[1 : len(value)-1]
						j.VersionType = dataType
					} else if string(keyUnesc) == "method" {
						if dataType != String {
							j.err = ErrIncorrectFieldType
							return
						}
						j.Method = value[1 : len(value)-1]
						j.MethodType = dataType
					} else if string(keyUnesc) == "params" {
						j.Params = value
						j.ParamsType = dataType
					}

					i += endOffset
				} else {
					j.err = ErrParseJSON
					return
				}
			} else {
				i--
			}
		case '{':
			level++
			// skip all, except first level
			if level > 1 {
				j.err = ErrParseJSON
				return
			}
		case '}':
			level--
		case '[':
			j.err = ErrParseJSON
			return
		case ':':
			j.err = ErrParseJSON
			return
		}

		i++
	}
}
