package deepl

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

const (
	// NameEnvKeyAPIDefault is default environment key name of DeepL API key to search.
	NameEnvKeyAPIDefault = "DEEPL_API_KEY"
	// StatusQuotaExceeded is the status code of DeepL API when the character limit
	// for translation has been reached.
	StatusQuotaExceeded = 456
	// UserAgentDefault is the default user agent string.
	UserAgentDefault = "Deepl-Go-Client"
)

var (
	// NameEnvKeyAPI is the environment variable name of DeepL API key to search.
	// By default, it is "DEEPL_API_KEY".
	NameEnvKeyAPI = NameEnvKeyAPIDefault
	// UserAgent is the user agent used in HTTP requests. By default, it is
	// "Deepl-Go-Client".
	UserAgent = UserAgentDefault
)

// ioReadAll is a copy of io.ReadAll for ease testing.
var ioReadAll = io.ReadAll

// ----------------------------------------------------------------------------
//  Types (Structs for JSON Unmarshaling)
// ----------------------------------------------------------------------------

type AccountStatus struct {
	CharacterCount int `json:"character_count"`
	CharacterLimit int `json:"character_limit"`
}

type ErrorResponse struct {
	ErrMessage string `json:"message"`
}

type TranslateResponse struct {
	Translations []translation `json:"translations"`
}

type translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

// ----------------------------------------------------------------------------
//  Private Functions
// ----------------------------------------------------------------------------

// decodeBody decodes JSON bytes to the given struct.
func decodeBody(bodyBytes []byte, outStruct interface{}) error {
	if err := json.Unmarshal(bodyBytes, outStruct); err != nil {
		return WrapIfErr(err, "failed to decode json")
	}

	return nil
}

// getAPIKey returns DeepL API key from environment variable.
// If the environment variable is not set or empty, it returns an error.
func getAPIKey() (string, error) {
	apiKey, ok := os.LookupEnv(NameEnvKeyAPI)
	if !ok {
		return "", NewErr("env variable for API key not set: %s", NameEnvKeyAPI)
	}

	if apiKey == "" {
		return "", NewErr("env var is set but empty: %s", NameEnvKeyAPI)
	}

	return apiKey, nil
}

// responseParse parses the response from DeepL API.
func responseParse(resp *http.Response, outStruct interface{}) error {
	if resp == nil || outStruct == nil {
		return NewErr("the input was nil")
	}

	bodyBytes, err := ioReadAll(resp.Body)
	if err != nil {
		return WrapIfErr(err, "failed to read response")
	}

	return treatBodyAsErr(resp.StatusCode, bodyBytes, outStruct)
}

// treatBodyAsErr treats the response body as an error message if the status code
// is not 200.
//
//nolint:cyclop // due to switch statement allow cyclomatic complexity be 16/10.
func treatBodyAsErr(status int, body []byte, outStruct interface{}) error {
	var (
		errResp    ErrorResponse
		errMessage string
	)

	// Capture the response body as an error message if the status code is not
	// 200.
	if status != http.StatusOK && len(body) != 0 {
		err := decodeBody(body, &errResp)
		if err != nil {
			return WrapIfErr(err, "failed to decode error response")
		}

		errMessage = errResp.ErrMessage
	}

	switch status {
	case http.StatusOK:
		err := decodeBody(body, &outStruct)

		return WrapIfErr(err, "failed to parse JSON response")
	case http.StatusBadRequest:
		return NewErr(
			"Bad request. Please check the error message and your parameters. Returned message: %s",
			errMessage,
		)
	case http.StatusUnauthorized:
		return NewErr(
			"Unauthorized. Please check your API key. Returned message: %s",
			errMessage,
		)
	case http.StatusForbidden:
		return NewErr("Authorization failed. Please supply a valid auth_key parameter. Returned message: %s",
			errMessage)
	case http.StatusNotFound:
		return NewErr("Not found. The requested resource clould not be found. Returned message: %s",
			errMessage)
	case http.StatusRequestEntityTooLarge:
		return NewErr("Request entity too large. The entity size exceeds the limit of each request. Returned message: %s",
			errMessage)
	case http.StatusTooManyRequests:
		return NewErr("Too many requests. Please wait and resend your request. Returned message: %s",
			errMessage)
	case StatusQuotaExceeded:
		return NewErr("Quota exceeded. The character limit has been reached. Returned message: %s",
			errMessage)
	case http.StatusServiceUnavailable:
		return NewErr("Service currently unavailable. Try again later. Returned message: %s",
			errMessage)
	default:
		// Internal error (5**) other than "http.StatusServiceUnavailable"(503)
		if 599 >= status && status >= 500 {
			return NewErr("Internal server error. Please try again later. Status code: %d, Returned message: %s",
				status, errMessage)
		}
	}

	return NewErr("Unexpected error. Status code: %d, Returned message: %s",
		status, errMessage)
}
