package headers

import (
	"bytes"
	"fmt"

	// "log/slog"
	"strings"
)

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func isValidToken(name []byte) bool {
	for _, ch := range name {
		if ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' {
			continue
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			continue
		}
		return false
	}
	return true
}

func (h *Headers) Get(name string) (string, bool) {
	str, ok := h.headers[strings.ToLower(name)] //even if uppercase got stored we still do lookup only after converting to lowercase
	return str, ok
}
func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)
	if v, ok := h.headers[name]; ok {
		h.headers[name] = fmt.Sprintf("%s,%s", v, value)
	} else {
		h.headers[strings.ToLower(name)] = value
	}
}

func (h *Headers) ForEach(cb func(n, v string)) {
	for n, v := range h.headers {
		cb(n, v)
	}
}

// field-line   = field-name ":" OWS field-value OWS
var headerSeparator = []byte("\r\n")

func parseHeaders(fieldLine []byte) (string, string, error) { //key value error
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed header")
	}
	value := string(bytes.TrimSpace(parts[1]))
	name := strings.ToLower(string(parts[0]))
	if bytes.HasSuffix([]byte(name), []byte(" ")) { //if there is space after the name then error
		return "", "", fmt.Errorf("malformed header")
	}
	return string(name), value, nil
}

func (h *Headers) Parse(data []byte) (n int, done bool, err error) {
	// Look for the end of the header section
	end := bytes.Index(data, []byte("\r\n\r\n"))
	if end == -1 {
		// Header not fully received yet
		return 0, false, nil
	}

	headerBlock := data[:end]
	lines := bytes.Split(headerBlock, []byte("\r\n"))

	for _, line := range lines {
		name, value, err := parseHeaders(line)
		if err != nil {
			return 0, false, err
		}
		h.Set(name, value)
	}

	// Return number of bytes consumed (headers + \r\n\r\n)
	return end + 4, true, nil
}
