package deepl

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
//  New
// ----------------------------------------------------------------------------

//nolint:paralleltest // do not parallelize due to global variable change
func TestNew_bad_custom_url(t *testing.T) {
	defer func() {
		SetCustomURL("")
	}()

	SetCustomURL("http://badurl.with.control.char/\n")

	cli, err := New(APICustom, nil)

	require.Error(t, err,
		"malformed URL should return an error")
	require.Nil(t, cli,
		"returned client should be nil on error")
	assert.Contains(t, err.Error(), "failed to parse URL",
		"it should contain the error reason")
	assert.Contains(t, err.Error(), "invalid control character in URL",
		"it should contain the underlying error reason")
}

// ============================================================================
//  Methods
// ============================================================================

// ----------------------------------------------------------------------------
//  Client.GetAccountStatus
// ----------------------------------------------------------------------------

//nolint:paralleltest // do not parallelize due to temporary env var change
func TestClient_GetAccountStatus(t *testing.T) {
	t.Setenv(NameEnvKeyAPI, dummyAuthKey) // Set dummy DeepL API key

	//nolint:varnamelen // tt is a test case by habit
	for index, tt := range dataGetAccountStatus {
		nameTest := fmt.Sprintf("test #%d: %s", index+1, tt.name)

		t.Run(nameTest, func(t *testing.T) {
			cli, teardown := spawnTestServer(
				t,
				tt.mockResponseHeaderFile,
				tt.mockResponseBodyFile,
				tt.expectMethod,
				tt.expectRequestPath,
				tt.expectRawQuery,
			)
			defer teardown()

			actualResponse, err := cli.GetAccountStatus(context.Background())
			if tt.expectErrMessage == "" {
				require.NoError(t, err, "response error should be nil on success")
				require.Equal(t, tt.expectResponse.CharacterCount, actualResponse.CharacterCount, "response items wrong")
				require.Equal(t, tt.expectResponse.CharacterLimit, actualResponse.CharacterLimit, "response items wrong")
			} else {
				require.Error(t, err, "response error should not be nil on error case")
			}
		})
	}
}

//nolint:paralleltest // do not parallelize due to global variable change
func TestClient_GetAccountStatus_no_env_set(t *testing.T) {
	oldNameEnvKeyAPI := NameEnvKeyAPI
	defer func() {
		NameEnvKeyAPI = oldNameEnvKeyAPI
	}()

	NameEnvKeyAPI = "ENV_VAL_WITH_NO_NAME" + t.Name()

	cli, err := New(APICustom, nil)
	require.NoError(t, err, "failed to create client")

	act, err := cli.GetAccountStatus(context.TODO())

	require.Error(t, err,
		"should return an error if env variable is not set")
	require.Nil(t, act,
		"should return nil on error")
	assert.Contains(t, err.Error(), "failed to get API key",
		"it should contain the error reason")
	assert.Contains(t, err.Error(), "env variable for API key not set",
		"it should contain the underlying error reason")
}

//nolint:paralleltest // do not parallelize due to temporary env var change
func TestClient_GetAccountStatus_bad_scheme(t *testing.T) {
	SetCustomURL("http://0.0.0.0:0") // Set dummy base URL

	defer func() {
		SetCustomURL("")
	}()

	t.Setenv(NameEnvKeyAPI, dummyAuthKey) // Set dummy DeepL API key

	badURL := new(url.URL)
	badURL.Host = "\n" // Set invalid host

	cli := &Client{
		BaseURL:    badURL,
		HTTPClient: http.DefaultClient,
		Logger:     &log.Logger{},
	}
	act, err := cli.GetAccountStatus(context.TODO())

	require.Error(t, err,
		"it should return an error if creating the request fails")
	require.Nil(t, act,
		"returned value should be nil on error")
	assert.Contains(t, err.Error(), "failed to create request",
		"it should contain the error reason")
}

//nolint:paralleltest // do not parallelize due to temporary env var change
func TestClient_GetAccountStatus_fail_request(t *testing.T) {
	SetCustomURL("http://0.0.0.0:0") // Set dummy base URL

	defer func() {
		SetCustomURL("")
	}()

	t.Setenv(NameEnvKeyAPI, dummyAuthKey) // Set dummy DeepL API key

	cli := &Client{
		BaseURL:    new(url.URL),
		HTTPClient: http.DefaultClient,
		Logger:     &log.Logger{},
	}

	act, err := cli.GetAccountStatus(context.TODO())

	require.Error(t, err,
		"it should return an error if sending request fails")
	require.Nil(t, act,
		"returned value should be nil on error")
	assert.Contains(t, err.Error(), "failed to send http request",
		"it should contain the error reason")
}

//nolint:paralleltest // do not parallelize due to temporary env var change
func TestClient_GetAccountStatus_bad_response_format(t *testing.T) {
	t.Setenv(NameEnvKeyAPI, dummyAuthKey) // Set dummy DeepL API key

	cli, teardown := spawnTestServer(
		t,
		"testdata/GetAccountStatus/success-header",
		"testdata/GetAccountStatus/malformed-body",
		http.MethodPost,
		"/v2/usage",
		fmt.Sprintf("auth_key=%s", dummyAuthKey),
	)
	defer teardown()

	actualResponse, err := cli.GetAccountStatus(context.Background())

	require.Error(t, err, "malformed response should return an error")
	require.Nil(t, actualResponse, "returned response should be nil on error")
	assert.Contains(t, err.Error(), "failed to parse response to AccountStatus",
		"it should contain the error reason")
	assert.Contains(t, err.Error(), "failed to parse JSON response",
		"it should contain the underlying error reason")
}

// ----------------------------------------------------------------------------
//  Client.TranslateSentence
// ----------------------------------------------------------------------------

//nolint:paralleltest // do not parallelize due to temporary env var change
func TestClient_TranslateSentence(t *testing.T) {
	t.Setenv(NameEnvKeyAPI, dummyAuthKey) // Set dummy DeepL API key

	//nolint:varnamelen // tt is a test case by habit
	for index, tt := range dataTranslateSentence {
		nameTest := fmt.Sprintf("test #%d: %s", index+1, tt.name)

		t.Run(nameTest, func(t *testing.T) {
			cli, teardown := spawnTestServer(
				t,
				tt.mockResponseHeaderFile,
				tt.mockResponseBodyFile,
				tt.expectMethod,
				tt.expectRequestPath,
				tt.expectRawQuery,
			)
			defer teardown()

			actualResponse, err := cli.TranslateSentence(
				context.TODO(),
				tt.inputText,
				tt.inputSourceLang,
				tt.inputTargetLang,
			)

			if tt.expectErrMessage == "" {
				require.NoError(t, err, "response error should be nil on success")

				for index, transVal := range actualResponse.Translations {
					require.Equal(
						t,
						tt.expectResponse.Translations[index].DetectedSourceLanguage,
						transVal.DetectedSourceLanguage,
						"response items wrong",
					)
					require.Equal(
						t,
						tt.expectResponse.Translations[index].Text,
						transVal.Text,
						"response items wrong",
					)
				}
			} else {
				require.Error(t, err,
					"response error should not be nil on error case")
				require.Contains(t, err.Error(), tt.expectErrMessage,
					"response error message does not contain expect message")
			}
		})
	}
}

//nolint:paralleltest // do not parallelize due to temporary env var change
func TestClient_TranslateSentence_no_env_set(t *testing.T) {
	oldNameEnvKeyAPI := NameEnvKeyAPI
	defer func() {
		NameEnvKeyAPI = oldNameEnvKeyAPI
	}()

	NameEnvKeyAPI = "ENV_VAL_WITH_NO_NAME" + t.Name()

	cli, err := New(APICustom, nil)
	require.NoError(t, err, "failed to create client")

	act, err := cli.TranslateSentence(
		context.TODO(),
		"hello",
		"EN",
		"JA",
	)

	require.Error(t, err,
		"should return an error if env variable is not set")
	require.Nil(t, act,
		"should return nil on error")
	assert.Contains(t, err.Error(), "failed to get API key",
		"it should contain the error reason")
	assert.Contains(t, err.Error(), "env variable for API key not set",
		"it should contain the underlying error reason")
}

//nolint:paralleltest // do not parallelize due to temporary env var change
func TestClient_TranslateSentence_bad_scheme(t *testing.T) {
	SetCustomURL("http://0.0.0.0:0") // Set dummy base URL

	defer func() {
		SetCustomURL("")
	}()

	t.Setenv(NameEnvKeyAPI, dummyAuthKey) // Set dummy DeepL API key

	badURL := new(url.URL)
	badURL.Host = "\n" // Set invalid host

	cli := &Client{
		BaseURL:    badURL,
		HTTPClient: http.DefaultClient,
		Logger:     &log.Logger{},
	}
	act, err := cli.TranslateSentence(
		context.TODO(),
		"hello",
		"EN",
		"JA",
	)

	require.Error(t, err,
		"it should return an error if creating the request fails")
	require.Nil(t, act,
		"returned value should be nil on error")
	assert.Contains(t, err.Error(), "failed to create request",
		"it should contain the error reason")
}

//nolint:paralleltest // do not parallelize due to temporary env var change
func TestClient_TranslateSentence_fail_request(t *testing.T) {
	SetCustomURL("http://0.0.0.0:0") // Set dummy base URL

	defer func() {
		SetCustomURL("")
	}()

	t.Setenv(NameEnvKeyAPI, dummyAuthKey) // Set dummy DeepL API key

	cli := &Client{
		BaseURL:    new(url.URL),
		HTTPClient: http.DefaultClient,
		Logger:     &log.Logger{},
	}

	act, err := cli.TranslateSentence(
		context.TODO(),
		"hello",
		"EN",
		"JA",
	)

	require.Error(t, err,
		"it should return an error if sending request fails")
	require.Nil(t, act,
		"returned value should be nil on error")
	assert.Contains(t, err.Error(), "failed to send http request",
		"it should contain the error reason")
}

// ============================================================================
//  Data Providers
// ============================================================================

// dataGetAccountStatus is a data provider for TestClient_GetAccountStatus.
var dataGetAccountStatus = []struct {
	name string

	mockResponseHeaderFile string
	mockResponseBodyFile   string

	expectMethod      string
	expectRequestPath string
	expectRawQuery    string
	expectResponse    *AccountStatus
	expectErrMessage  string
}{
	{
		name: "success",

		mockResponseHeaderFile: "testdata/GetAccountStatus/success-header",
		mockResponseBodyFile:   "testdata/GetAccountStatus/success-body",

		expectMethod:      http.MethodPost,
		expectRequestPath: "/v2/usage",
		expectRawQuery:    fmt.Sprintf("auth_key=%s", dummyAuthKey),
		expectResponse:    &AccountStatus{CharacterCount: 30315, CharacterLimit: 1000000},
	},
}

// dataTranslateSentence is a data provider for TestClient_TranslateSentence.
var dataTranslateSentence = []struct {
	name string

	inputText       string
	inputSourceLang string
	inputTargetLang string

	mockResponseHeaderFile string
	mockResponseBodyFile   string

	expectMethod      string
	expectRequestPath string
	expectRawQuery    string
	expectResponse    *TranslateResponse
	expectErrMessage  string
}{
	{
		name: "success",

		inputText:       "hello",
		inputSourceLang: "EN",
		inputTargetLang: "JA",

		mockResponseHeaderFile: "testdata/TranslateText/success-header",
		mockResponseBodyFile:   "testdata/TranslateText/success-body",

		expectMethod:      http.MethodPost,
		expectRequestPath: "/v2/translate",
		expectRawQuery:    fmt.Sprintf("auth_key=%s&source_lang=EN&target_lang=JA&text=hello", dummyAuthKey),
		expectResponse:    createTranslateResponse("EN", "こんにちわ"),
	},
	{
		name: "misssing target_lang",

		inputText:       "hello",
		inputSourceLang: "EN",
		inputTargetLang: "",

		mockResponseHeaderFile: "testdata/TranslateText/missing-target_lang-header",
		mockResponseBodyFile:   "testdata/TranslateText/missing-target_lang-body",

		expectMethod:      http.MethodPost,
		expectRequestPath: "/v2/translate",
		expectRawQuery:    fmt.Sprintf("auth_key=%s&source_lang=EN&target_lang=&text=hello", dummyAuthKey),
		expectErrMessage:  "Bad request.",
	},
	{
		name: "unsuport target_lang",

		inputText:       "hello",
		inputSourceLang: "EN",
		inputTargetLang: "AA",

		mockResponseHeaderFile: "testdata/TranslateText/unsuport-target_lang-header",
		mockResponseBodyFile:   "testdata/TranslateText/unsuport-target_lang-body",

		expectMethod:      http.MethodPost,
		expectRequestPath: "/v2/translate",
		expectRawQuery:    fmt.Sprintf("auth_key=%s&source_lang=EN&target_lang=AA&text=hello", dummyAuthKey),
		expectErrMessage:  "Bad request.",
	},
	{
		name: "wrong api key",

		inputText:       "hello",
		inputSourceLang: "EN",
		inputTargetLang: "JA",

		mockResponseHeaderFile: "testdata/TranslateText/wrong-apikey-header",
		mockResponseBodyFile:   "testdata/TranslateText/wrong-apikey-body",

		expectMethod:      http.MethodPost,
		expectRequestPath: "/v2/translate",
		expectRawQuery:    fmt.Sprintf("auth_key=%s&source_lang=EN&target_lang=JA&text=hello", dummyAuthKey),
		expectErrMessage:  "Authorization failed.",
	},
}
