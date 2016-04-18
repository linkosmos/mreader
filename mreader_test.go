package mreader

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromReader(t *testing.T) {
	expected := []byte("Some text")
	inputAsIoReader := bytes.NewReader(expected)

	raw, reader, err := FromReader(inputAsIoReader)
	assert.Nil(t, err)
	assert.Equal(t, expected, raw)
	assert.IsType(t, &bytes.Reader{}, reader)
}

func TestFromHTMLResponse(t *testing.T) {
	raw, reader, err := FromHTMLResponse(nil)
	assert.Nil(t, raw)
	assert.Nil(t, reader)
	assert.Equal(t, ErrResponseIsNil, err)

	raw, reader, err = FromHTMLResponse(&http.Response{})
	assert.Nil(t, raw)
	assert.Nil(t, reader)
	assert.Equal(t, ErrResponseBodyIsNil, err)

	response := &http.Response{}
	response.Header = make(http.Header)
	response.Header.Set("Content-Type", "not-HTML")
	response.Body = ioutil.NopCloser(&bytes.Reader{})
	raw, reader, err = FromHTMLResponse(response)
	assert.Nil(t, raw)
	assert.Nil(t, reader)
	assert.Equal(t, ErrResponseContentTypeNotHTML, err)

	response = &http.Response{}
	response.Header = make(http.Header)
	response.Header.Set("Content-Type", "text/html")
	expected := []byte("Some bytes as mocked HTML document")
	response.Body = ioutil.NopCloser(bytes.NewReader(expected))
	raw, reader, err = FromHTMLResponse(response)
	assert.Equal(t, expected, raw)
	assert.IsType(t, &bytes.Reader{}, reader)
	assert.Nil(t, err)
}
