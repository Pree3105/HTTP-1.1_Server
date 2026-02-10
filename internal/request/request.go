package request

import (
	"bytes"
	"fmt"
	"httpServer_go/internal/headers"
	"io"
	"strconv"
	"strings"
)

type parserState string

const (
	StateInit                  parserState = "init"
	StateDone                  parserState = "done"
	StateError                 parserState = "error"
	requestStateParsingHeaders parserState = "headers"
	StateBody                  parserState = "body"
)

type Request struct {
	RequestLine  RequestLine
	Headers      *headers.Headers
	Body         string
	RequestState parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func newRequest() *Request {
	return &Request{
		RequestState: StateInit,
		Headers:      headers.NewHeaders(),
	}
}
func getInt(headers *headers.Headers, name string, defaultval int) int {
	valueStr, exists := headers.Get(strings.ToLower(name))
	if !exists {
		return defaultval
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultval
	}
	return value
}

func (r *Request) hasBody() bool {
	//todo:refactor during chunked encoding
	lengthStr := getInt(r.Headers, "content-length", 0)
	return lengthStr > 0
}

var separator = []byte("\r\n")
var incorrect_request = fmt.Errorf("INCORRECT REQUEST STRUCTURE")
var ErrorRequestInErrorState = fmt.Errorf("REQUEST IN ERROR STATE")
var wrong_version = fmt.Errorf("INCORRECT HTTP VERSION")

// takes r /request(or start)line , get reqline(GET /home HTTP/1.1\r\n) till \r\n, split into 3 parts based on " "
func parseRequestLine(r []byte) (*RequestLine, int, error) {
	i := bytes.Index(r, separator)
	if i == -1 {
		return nil, 0, nil
	}
	reqline := r[:i]
	read := i + len(separator) //u have read till \r\n

	parts := bytes.Split(reqline, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, incorrect_request
	}

	meth := string(parts[0])
	reqTarget := string(parts[1])
	httpversion := string(parts[2])

	httpParts := bytes.Split([]byte(httpversion), []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, wrong_version
	}
	rl := &RequestLine{ //returning a pointer to RequestLine struct
		HttpVersion:   string(httpParts[1]),
		RequestTarget: reqTarget,
		Method:        meth,
	}

	return rl, read, nil
}

func (r *Request) parse(data []byte) (int, error) {
	//if u get a very large piece of data that has many state transitions, so u can keep looping and parse the start line, header, body,blah
	read := 0
	for {
		fmt.Printf("State: %s, Read so far: %d, Data length: %d\n", r.RequestState, read, len(data[read:]))
		if r.RequestState == StateInit {
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				return read, nil
			}
			r.RequestLine = *rl
			read += n
			r.RequestState = requestStateParsingHeaders
		} else if r.RequestState == requestStateParsingHeaders {
			n, done, err := r.Headers.Parse(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				return read, nil
			}
			read += n
			if done {
				if r.hasBody() {
					r.RequestState = StateBody
				} else {
					r.RequestState = StateDone
				}
				continue
			}
			return read, nil
		} else if r.RequestState == StateBody {
			lengthStr := getInt(r.Headers, "content-length", 0)
			if lengthStr == 0 {
				panic("chunked encoding not implemented")
			}
			currentData := data[read:]
			if len(currentData) == 0 {
				return read, nil
			}
			remaining := lengthStr - len(r.Body) //so we dont read more than whats left in buffer
			if remaining > len(currentData) {    //if there is more data in the body then the remaining bytes the buffer has read(currentData)
				r.Body += string(currentData) //then append only till only those bytes that buffer has read
				read += len(currentData)
				return read, nil //keep reading
			} else {
				r.Body += string(currentData[:remaining]) //otherwise append all the remaining data of the body
				read += remaining
				r.RequestState = StateDone //body parsing done
				return read, nil
			}
			//r.Body += string(currentData[:remaining])
		} else if r.RequestState == StateDone {
			return read, nil //returning how much we have read
		} else if r.RequestState == StateError {
			return 0, ErrorRequestInErrorState
		}
	}
}

func (r *Request) done() bool {
	return r.RequestState == StateDone || r.RequestState == StateError
}

// func (r *Request) error() bool{
// 	return r.RequestState==StateError
// }

func RequestFromReader(reader io.Reader) (*Request, error) {
	//data, err := io.ReadAll(reader) //read all is not a good practice as it can hang, we only need to read the headers, body parsing later

	request := newRequest()
	//here buffer could overrun if header is bigger
	buf := make([]byte, 1024) //we will be reading the req in chunks of buf
	bufIdx := 0
	for !request.done() {
		n, err := reader.Read(buf[bufIdx:]) //filling buf from bufIdx to end,....it returns n/the last index num it read
		if n > 0 {
			readN, err := request.parse(buf[:bufIdx+n]) //returns num of successfully parsed bytes from leftover unparsed bytes from previous reads(buf[:bufIdx]) & newly read bytes(buf[bufIdx:bufIdx+n])
			if err != nil {
				return nil, err
			}
			copy(buf, buf[readN:bufIdx+n]) //copying leftovers that didnt get parsed into buf/ shortening buf
			bufIdx = bufIdx + n - readN
		}
		if err == io.EOF {
			if !request.done() {
				return nil, fmt.Errorf("unexpected EOF") // only an error if request is NOT finished
			}
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return request, nil
}
