package deepl

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dummyAuthKey is a dummy DeepL API key for testing.
const dummyAuthKey = "12345678-1234-abcd-abcd-123456789abc"

// ----------------------------------------------------------------------------
//  Helpers
// ----------------------------------------------------------------------------

// DummyReadCloser is a dummy implementation of io.ReadCloser for testing.
type DummyReadCloser struct {
	DummyBody   []byte
	ForcedError error
}

// Read writes the DummyBody to the given buffer if ForcedError is nil. If not,
// it returns ForcedError.
func (d *DummyReadCloser) Read(buf []byte) (int, error) {
	if d.ForcedError != nil {
		return 0, WrapIfErr(d.ForcedError, "forced error")
	}

	return copy(buf, d.DummyBody), io.EOF
}

// Close is the dummy implementation of io.ReadCloser.
func (d *DummyReadCloser) Close() error {
	return nil
}

// ============================================================================
//  Test Cases (test data follows at the bottom of this file)
// ============================================================================

// ----------------------------------------------------------------------------
//  Private Functions
// ----------------------------------------------------------------------------

func Test_decodeBody_nil_body(t *testing.T) {
	t.Parallel()

	err := decodeBody(nil, nil)
	require.Error(t, err,
		"it should return error on nil body")
	assert.Contains(t, err.Error(), "failed to decode json",
		"returned error should contain the reason")
}

//nolint:paralleltest // do not parallelize due to global variable change
func Test_getAPIKey(t *testing.T) {
	oldNameEnvKeyAPI := NameEnvKeyAPI
	defer func() {
		NameEnvKeyAPI = oldNameEnvKeyAPI
	}()

	NameEnvKeyAPI = "UNDEFINED_ENV_KEY" + t.Name()

	t.Run("env key is not set", func(t *testing.T) {
		val, err := getAPIKey()
		require.Error(t, err,
			"undefined env key should return error")
		require.Empty(t, val,
			"returned value should be empty on error")
		assert.Contains(t, err.Error(), "env variable for API key not set",
			"returned error should contain the reason")
	})

	t.Run("env key is set but empty", func(t *testing.T) {
		t.Setenv(NameEnvKeyAPI, "")

		val, err := getAPIKey()
		require.Error(t, err,
			"empty env value should return error")
		require.Empty(t, val,
			"returned value should be empty on error")
		assert.Contains(t, err.Error(), "env var is set but empty",
			"returned error should contain the reason")
	})
}

func Test_responseParse_nil_input(t *testing.T) {
	t.Parallel()

	err := responseParse(nil, nil)

	require.Error(t, err,
		"it should return error on nil body")
	require.Contains(t, err.Error(), "the input was nil",
		"returned error should contain the reason")
}

//nolint:paralleltest // do not parallelize due to global variable change
func Test_responseParse_read_body_fail(t *testing.T) {
	// Backup and defer restore of ioReadAll
	oldIOReadAll := ioReadAll
	defer func() {
		ioReadAll = oldIOReadAll
	}()

	// Mock ioReadAll to force return error
	ioReadAll = func(r io.Reader) ([]byte, error) {
		return nil, NewErr("forced error")
	}

	resp := new(http.Response)
	accountStatusResp := new(AccountStatus)
	err := responseParse(resp, accountStatusResp)

	require.Error(t, err,
		"it should return error on invalid input")
	require.Contains(t, err.Error(), "failed to read response",
		"returned error should contain the reason")
	require.Contains(t, err.Error(), "forced error",
		"returned error should contain the underlying error")
}

func Test_responseParse_fail_decode_response(t *testing.T) {
	t.Parallel()

	content := "this is: invalid json format"

	resp := new(http.Response)

	resp.Status = "400 Bad Request"
	resp.StatusCode = http.StatusBadRequest
	resp.Body = &DummyReadCloser{
		DummyBody:   []byte(content),
		ForcedError: nil,
	}
	resp.ContentLength = int64(len(content))

	accountStatusResp := new(AccountStatus)

	err := responseParse(resp, accountStatusResp)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decode error response",
		"returned error should contain the reason")
	require.Contains(t, err.Error(), "failed to decode json",
		"returned error should contain the underlying error")
}

//nolint:paralleltest,varnamelen // do not parallelize due to range looping
func Test_responseParse_status_msg(t *testing.T) {
	for index, tt := range dataResponseParse {
		nameTest := fmt.Sprintf("test #%d: %s", index+1, tt.name)

		t.Run(nameTest, func(t *testing.T) {
			// Server response
			resp := new(http.Response)

			resp.Status = tt.respStatus
			resp.StatusCode = tt.statusCode
			resp.Body = &DummyReadCloser{
				DummyBody:   []byte(tt.respBody),
				ForcedError: nil,
			}

			// Response struct to be parsed
			accountStatusResp := new(AccountStatus)

			err := responseParse(resp, accountStatusResp)

			if tt.expectMsgCore != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectMsgCore,
					"returned error should contain the reason")
				assert.Contains(t, err.Error(), tt.expectMsgSub,
					"returned error should contain the response body")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ============================================================================
//  Data Providers
// ============================================================================

var dataResponseParse = []struct {
	name          string
	respStatus    string
	respBody      string
	expectMsgCore string // expected msg to contain in error
	expectMsgSub  string // expected msg to contain in response
	statusCode    int
}{
	{
		name:          "success",
		respStatus:    "200 OK",
		respBody:      `{"character_count": 10, "character_limit": 20}`,
		expectMsgCore: "",
		expectMsgSub:  "",
		statusCode:    http.StatusOK,
	},
	{
		name:          "bad request",
		respStatus:    "400 Bad Request",
		respBody:      `{"message": "BAD REQUEST"}`,
		expectMsgCore: "Bad request. Please check the error message and your parameters.",
		expectMsgSub:  "BAD REQUEST",
		statusCode:    http.StatusBadRequest,
	},
	{
		name:          "unauthorized",
		respStatus:    "401 Unauthorized",
		respBody:      `{"message": "UNAUTHORIZED"}`,
		expectMsgCore: "Unauthorized. Please check your API key.",
		expectMsgSub:  "UNAUTHORIZED",
		statusCode:    http.StatusUnauthorized,
	},
	{
		name:          "forbidden",
		respStatus:    "403 Forbidden",
		respBody:      `{"message": "FORBIDDEN"}`,
		expectMsgCore: "Authorization failed. Please supply a valid auth_key parameter.",
		expectMsgSub:  "FORBIDDEN",
		statusCode:    http.StatusForbidden,
	},
	{
		name:          "not found",
		respStatus:    "404 Not Found",
		respBody:      `{"message": "NOT FOUND"}`,
		expectMsgCore: "Not found. The requested resource clould not be found.",
		expectMsgSub:  "NOT FOUND",
		statusCode:    http.StatusNotFound,
	},
	{
		name:          "too large request",
		respStatus:    "413 Request Entity Too Large",
		respBody:      `{"message": "TOO LARGE REQUEST"}`,
		expectMsgCore: "Request entity too large. The entity size exceeds the limit of each request.",
		expectMsgSub:  "TOO LARGE REQUEST",
		statusCode:    http.StatusRequestEntityTooLarge,
	},
	{
		name:          "too many requests",
		respStatus:    "429 Too Many Requests",
		respBody:      `{"message": "TOO MANY REQUESTS"}`,
		expectMsgCore: "Too many requests. Please wait and resend your request.",
		expectMsgSub:  "TOO MANY REQUESTS",
		statusCode:    http.StatusTooManyRequests,
	},
	{
		name:          "quota exceeded",
		respStatus:    "456 Quota Exceeded",
		respBody:      `{"message": "QUOTA EXCEEDED"}`,
		expectMsgCore: "Quota exceeded. The character limit has been reached.",
		expectMsgSub:  "QUOTA EXCEEDED",
		statusCode:    StatusQuotaExceeded,
	},
	{
		name:          "service unavailable",
		respStatus:    "503 Service Unavailable",
		respBody:      `{"message": "SERVICE UNAVAILABLE"}`,
		expectMsgCore: "Service currently unavailable. Try again later.",
		expectMsgSub:  "SERVICE UNAVAILABLE",
		statusCode:    http.StatusServiceUnavailable,
	},
	{
		name:          "internal server error",
		respStatus:    "500 Internal Server Error",
		respBody:      `{"message": "INTERNAL SERVER ERROR"}`,
		expectMsgCore: "Internal server error. Please try again later.",
		expectMsgSub:  "INTERNAL SERVER ERROR",
		statusCode:    http.StatusInternalServerError,
	},
	{
		name:          "unknown expected error",
		respStatus:    "999 Unknown Error",
		respBody:      `{"message": "UNKNOWN ERROR"}`,
		expectMsgCore: "Unexpected error.",
		expectMsgSub:  "UNKNOWN ERROR",
		statusCode:    999,
	},
}
