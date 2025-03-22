package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRequest(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}

	r, err := RequestFromReader(reader)

	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
}

func TestGetRequestWithPath(t *testing.T) {
	reader := &chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}

	r, err := RequestFromReader(reader)

	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
}

func TestPostRequest(t *testing.T) {
	reader := &chunkReader{
		data:            "POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}

	r, err := RequestFromReader(reader)

	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
}

func TestInvalidMethod(t *testing.T) {
	reader := &chunkReader{
		data:            "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}

	_, err := RequestFromReader(reader)

	require.Error(t, err)
}

func TestInvalidNumberOfParts(t *testing.T) {
	reader := &chunkReader{
		data:            "GET /coffee\r\n",
		numBytesPerRead: 8,
	}

	_, err := RequestFromReader(reader)

	require.Error(t, err)
}

func TestInvalidMethodOutOfOrder(t *testing.T) {
	reader := &chunkReader{
		data:            "/coffee GET HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 4,
	}

	_, err := RequestFromReader(reader)

	require.Error(t, err)
}

func TestInvalidVersion(t *testing.T) {
	reader := &chunkReader{
		data:            "GET /coffee HTTP/420.69\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}

	_, err := RequestFromReader(reader)

	require.Error(t, err)
}

func TestRequestWithHeaders(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}

	r, err := RequestFromReader(reader)

	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
	assert.Equal(t, "localhost:42069", r.Headers.Get("Host"))
	assert.Equal(t, "curl/7.81.0", r.Headers.Get("User-Agent"))
	assert.Equal(t, "*/*", r.Headers.Get("Accept"))
	assert.Equal(t, r.Status, RequestStateParsingBody)
}

func TestRequestWithHeadersMultipleValuesForOneKey(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*; application/json;\r\n\r\n",
		numBytesPerRead: 3,
	}

	r, err := RequestFromReader(reader)

	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
	assert.Equal(t, "localhost:42069", r.Headers.Get("Host"))
	assert.Equal(t, "curl/7.81.0", r.Headers.Get("User-Agent"))
	assert.Equal(t, "*/*; application/json;", r.Headers.Get("Accept"))
	assert.Equal(t, r.Status, RequestStateParsingBody)
}

func TestInvalidRequestWithHeaders(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
		numBytesPerRead: 3,
	}

	_, err := RequestFromReader(reader)

	require.Error(t, err)
}
