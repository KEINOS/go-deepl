package deepl

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// createTranslateResponse is a short hand function to create a TranslateResponse
// as a dummy response during testing.
func createTranslateResponse(detectLang string, text string) *TranslateResponse {
	resp := &TranslateResponse{
		[]translation{
			{
				DetectedSourceLanguage: detectLang,
				Text:                   text,
			},
		},
	}

	return resp
}

// spawnTestServer returns a test server and a teardown function.
func spawnTestServer(
	t *testing.T,
	mockResponseHeaderFile,
	mockResponseBodyFile,
	expectedMethod,
	expectedRequestPath,
	expectedRawQuery string,
) (*Client, func()) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(respWriter http.ResponseWriter, req *http.Request) {
		require.Equal(t, expectedMethod, req.Method, "request method is wrong")
		require.Equal(t, expectedRequestPath, req.URL.Path, "request path is wrong")
		require.Equal(t, expectedRawQuery, req.URL.RawQuery, "request query is wrong")

		headerBytes, err := os.ReadFile(mockResponseHeaderFile)
		require.NoError(t, err, "failed to read header '%s'", mockResponseHeaderFile)

		firstLine := strings.Split(string(headerBytes), "\n")[0]

		statusCode, err := strconv.Atoi(strings.Fields(firstLine)[1])
		require.NoError(t, err, "failed to extract status code from header")

		respWriter.WriteHeader(statusCode)

		bodyBytes, err := os.ReadFile(mockResponseBodyFile)
		require.NoError(t, err, "failed to read body '%s'", mockResponseBodyFile)

		_, err = respWriter.Write(bodyBytes)
		require.NoError(t, err, "failed to write response body")
	}))

	serverURL, err := url.Parse(server.URL)
	require.NoError(t, err, "failed to get mock server URL")

	cli := &Client{
		BaseURL:    serverURL,
		HTTPClient: server.Client(),
		Logger:     nil,
	}
	teardown := func() {
		server.Close()
	}

	return cli, teardown
}
