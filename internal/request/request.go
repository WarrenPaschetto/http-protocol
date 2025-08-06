package request

import (
	"errors"
	"io"
	"log"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	reqBytes, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("error reading request line: %s", err)
		return nil, err
	}
	reqString := string(reqBytes)
	r, err := parseRequestLine(reqString)
	if err != nil {
		log.Printf("error parsing request line: %s", err)
		return nil, err
	}
	return r, nil
}

func parseRequestLine(request string) (*Request, error) {
	lines := strings.Split(request, "\r\n")
	if len(lines) < 1 || len(lines[0]) == 0 {
		return nil, errors.New("empty request line")
	}

	parts := strings.Split(lines[0], " ")
	if len(parts) != 3 {
		return nil, errors.New("malformed request line: must contain method, version, and target")
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	// Validate order of request line
	if strings.Contains(method, "HTTP") || strings.Contains(target, "GET") {
		return nil, errors.New("invalid method. Request line must be in order: method, version, target")
	}

	// Validate method: must be uppercase alphabetic only
	for _, ch := range method {
		if !unicode.IsUpper(ch) {
			return nil, errors.New("invalid method: must be uppercase letters only")
		}
	}

	// Validate version
	if version != "HTTP/1.1" {
		return nil, errors.New("unsupported HTTP version: only HTTP/1.1 supported")
	}

	return &Request{
		RequestLine: RequestLine{
			HttpVersion:   version,
			RequestTarget: target,
			Method:        method,
		},
	}, nil
}
