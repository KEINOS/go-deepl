package deepl

// APIType is an enum type for the API version.
type APIType int

const (
	// APICustom represents the test API in the local server.
	// To use this, you need to set the base URL of the local server
	// via SetCustomURL() method.
	APICustom APIType = iota
	// APIFree represents the DeepL API for free account. It has a limitation
	// of 500,000 characters per month.
	APIFree
	// APIPro represents the DeepL API for pro/paid account. It has no limitation
	// of characters per month.
	APIPro
)

// SetCustomURL sets the base URL to be forced to use.
func SetCustomURL(url string) {
	baseURLCustom = url
}

// BaseURLCustom is the base URL of the test API in the local server.
var baseURLCustom string

// String is a Stringer implementation for APIType. Which returns the base URL
// of the API as it's string representation.
func (a APIType) String() string {
	return a.BaseURL()
}

// BaseURL returns the base URL of the API.
func (a APIType) BaseURL() string {
	baseURL := baseURLCustom

	switch a {
	case APIPro:
		baseURL = "https://api.deepl.com"
	case APIFree:
		baseURL = "https://api-free.deepl.com"
	case APICustom:
	default:
	}

	// Force to use the custom URL if set
	if baseURLCustom != "" {
		baseURL = baseURLCustom
	}

	return baseURL
}
