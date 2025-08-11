package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		// the empty line
		// headers are done, consume the CRLF
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := string(parts[0])

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)

	notAllowed := `(),/:;<=>?@[]{}\`

	for _, c := range notAllowed {
		if strings.ContainsRune(key, c) {
			return 0, false, fmt.Errorf("invalid character '%c' in header name: %s", c, key)
		}
	}

	if _, exists := h[strings.ToLower(key)]; exists {
		h.Set(strings.ToLower(key), h[strings.ToLower(key)]+", "+string(value))
	} else {
		h.Set(strings.ToLower(key), string(value))
	}
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	h[key] = value
}
