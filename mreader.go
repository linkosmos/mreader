package mreader

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html/charset"
)

// -
var (
	ErrResponseIsNil              = errors.New("Supplied response is nil")
	ErrResponseBodyIsNil          = errors.New("Response body is nil")
	ErrResponseContentTypeNotHTML = errors.New("Response header 'Content-Type' is not 'text/html'")
)

// FromHTMLResponse - takes response Body, converts to UTF-8, reads into buffer
// and wraps into bytes Reader
func FromHTMLResponse(response *http.Response) ([]byte, *bytes.Reader, error) {
	if response == nil {
		return nil, nil, ErrResponseIsNil
	}

	if response.Body == nil {
		return nil, nil, ErrResponseBodyIsNil
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return nil, nil, ErrResponseContentTypeNotHTML
	}

	// Converts HTML body to UTF-8
	utf8, err := charset.NewReader(
		response.Body,
		contentType, // Takes response content type charset to determine encoding
	)
	if err != nil {
		return nil, nil, err
	}

	return FromReader(utf8)
}

// FromReader - buffers io.Reader and wraps with bytes Reader
func FromReader(input io.Reader) ([]byte, *bytes.Reader, error) {

	// Replace with more efficient
	// growable buffer with inverse (diminishing) backoff
	// https://golang.org/src/io/io.go?s=10995:11049#L304
	// io.ReadAtLeast
	// http://openmymind.net/Go-Slices-And-The-Case-Of-The-Missing-Memory/
	buffer, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, nil, err
	}

	return buffer, bytes.NewReader(buffer), nil
}
