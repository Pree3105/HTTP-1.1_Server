package response

import (
	"fmt"
	"httpServer_go/internal/headers"
	"io"
)

type StatusCode int

const (
	OK          StatusCode = 200
	BadRequest  StatusCode = 400
	ServerError StatusCode = 500
)

type Response struct {
	RequestLine string
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var reasonPhrase string
	switch statusCode {
	case OK:
		reasonPhrase = "OK"
	case BadRequest:
		reasonPhrase = "Bad Request"
	case ServerError:
		reasonPhrase = "Internal Server Error"
	default:
		reasonPhrase = ""
	}
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s", statusCode, reasonPhrase)
	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return *h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	headers.ForEach(func(n, v string) {
		fmt.Sprintf("%s %s\r\n", n, v)
	})
	return nil
}
