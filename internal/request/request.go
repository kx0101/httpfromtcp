package request

import "io"

const (
	BufferSize = 8
)

type Request struct {
	RequestLine RequestLine
	Status      Status
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Status int

const (
	Initialized Status = iota
	Done
)

var statusName = map[Status]string{
	Initialized: "initialized",
	Done:        "done",
}

func (s Status) String() string {
	return statusName[s]
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, BufferSize, BufferSize)
	readToIndex := 0

	r := Request{
		RequestLine: RequestLine{},
		Status:      Initialized,
	}

	for {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)

			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		readToIndex += n

		bytesParsed, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		if bytesParsed > 0 {
			copy(buf, buf[bytesParsed:])
			readToIndex -= bytesParsed
		}

		if r.Status == Done {
			break
		}
	}

	return &r, nil
}
