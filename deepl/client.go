package deepl

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

// ----------------------------------------------------------------------------
//  Type: Client
// ----------------------------------------------------------------------------

// Client is a struct that holds the basic client information.
type Client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	Logger     *log.Logger
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New returns a new Client instance.
// It will request to the given rawBaseURL and use the given logger. If the logger
// is nil, it will use the default logger which simply logs to stderr.
func New(apiType APIType, logger *log.Logger) (*Client, error) {
	rawBaseURL := apiType.BaseURL()

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, WrapIfErr(err, "failed to parse URL")
	}

	if logger == nil {
		logger = log.New(os.Stderr, "[Log]", log.LstdFlags)
	}

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
		Logger:     logger,
	}, nil
}

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

// GetAccountStatus returns the account status.
// The API key is retrieved via getAPIKey() function which tries to get the API
// key from the environment variable "DEEPL_API_KEY".
func (c *Client) GetAccountStatus(ctx context.Context) (*AccountStatus, error) {
	apiKey, err := getAPIKey()
	if err != nil {
		return nil, WrapIfErr(err, "failed to get API key")
	}

	// Set endpoint path of the API
	reqURL := *c.BaseURL
	reqURL.Path = path.Join(reqURL.Path, "v2", "usage")

	// Set query parameters
	urlVal := reqURL.Query()

	urlVal.Add("auth_key", apiKey)

	reqURL.RawQuery = urlVal.Encode()

	// Make new request
	req, err := http.NewRequest(http.MethodPost, reqURL.String(), nil)
	if err != nil {
		return nil, WrapIfErr(err, "failed to create request")
	}

	// Set header
	req.Header.Set("User-Agent", UserAgent)

	// Set context
	req = req.WithContext(ctx)

	// Request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, WrapIfErr(err, "failed to send http request")
	}

	defer resp.Body.Close()

	var accountStatusResp AccountStatus
	if err := responseParse(resp, &accountStatusResp); err != nil {
		return nil, WrapIfErr(err, "failed to parse response to AccountStatus")
	}

	return &accountStatusResp, nil
}

// TranslateSentence translates the given text from the sourceLang to the targetLang.
func (c *Client) TranslateSentence(
	ctx context.Context,
	text string,
	sourceLang string,
	targetLang string,
) (*TranslateResponse, error) {
	apiKey, err := getAPIKey()
	if err != nil {
		return nil, WrapIfErr(err, "failed to get API key")
	}

	// Set endpoint path of the API
	reqURL := *c.BaseURL
	reqURL.Path = path.Join(reqURL.Path, "v2", "translate")

	// Set query parameters
	urlVal := reqURL.Query()

	urlVal.Add("auth_key", apiKey)
	urlVal.Add("text", text)
	urlVal.Add("target_lang", targetLang)
	urlVal.Add("source_lang", sourceLang)

	reqURL.RawQuery = urlVal.Encode()

	// Make new request
	req, err := http.NewRequest(http.MethodPost, reqURL.String(), nil)
	if err != nil {
		return nil, WrapIfErr(err, "failed to create request")
	}

	// Set header
	req.Header.Set("User-Agent", UserAgent)

	// Set context
	req = req.WithContext(ctx)

	// Request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, WrapIfErr(err, "failed to send http request")
	}

	defer resp.Body.Close()

	var transResp TranslateResponse
	if err := responseParse(resp, &transResp); err != nil {
		return nil, WrapIfErr(err, "failed to parse response to TranslateResponse")
	}

	return &transResp, nil
}
