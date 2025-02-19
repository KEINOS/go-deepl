package main

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  Type: Translator
// ----------------------------------------------------------------------------

// Translator is the main type of this package. It holds the options and the
// authentication key needed to request the DeepL translation API.
//
// All the fields are private and can be accessed via the Property method. For
// the available properties, see the example of Property method.
// To instantiate a Translator object, use NewTranslator (constructor) function.
type Translator struct {
	authKey            string
	formality          string
	glossaryID         string
	sourceLang         string
	splitSentences     string
	tagHandling        string
	targetLang         string
	ignoreTags         []string
	nonSplittingTags   []string
	splittingTags      []string
	outlineDetection   bool
	preserveFormatting bool
}

// NewTranslator is a constructor of Translator. Required options are WithAuthKey and
// WithtargetLang. Other options are optional.
func NewTranslator(options ...Option) *Translator {
	translator := &Translator{}

	// Set default
	WithDefault()(translator)

	// Set user options (skip if empty)
	for _, option := range options {
		option(translator)
	}

	return translator
}

// Property returns the value of the field specified by the nameKey.
func (trans Translator) Property(nameKey string) any {
	switch nameKey {
	case "authKey":
		return trans.authKey
	case "formality":
		return trans.formality
	case "glossaryID":
		return trans.glossaryID
	case "ignoreTags":
		return trans.ignoreTags
	case "nonSplittingTags":
		return trans.nonSplittingTags
	case "outlineDetection":
		return trans.outlineDetection
	case "preserveFormatting":
		return trans.preserveFormatting
	case "sourceLang":
		return trans.sourceLang
	case "splitSentences":
		return trans.splitSentences
	case "splittingTags":
		return trans.splittingTags
	case "tagHandling":
		return trans.tagHandling
	case "targetLang":
		return trans.targetLang
	default:
		return nil
	}
}

// IsKeyFree returns true if the authentication key is for free account.
func (trans Translator) IsKeyFree() bool {
	// Prefix of free account.
	// - Ref: https://www.deepl.com/docs-api/api-access/api-versions/#authentication
	const prefix = ":fx"

	return trans.authKey[len(trans.authKey)-len(prefix):] == prefix
}

// TranslateText is similar to TranslateString but it takes a slice of strings
// and returns a slice of Translations objects or an error if any.
func (trans *Translator) TranslateText(input []string, options ...Option) ([]Translations, error) {
	Translations := make([]Translations, len(input))

	for index, text := range input {
		result, err := trans.TranslateString(text, options...)
		if err != nil {
			return nil, err
		}

		Translations[index] = result
	}

	return Translations, nil
}

// TranslateString translates the given text and returns the Translations text
// (Translations object) or an error if any.
func (trans *Translator) TranslateString(input string, options ...Option) (Translations, error) {
	if strings.TrimSpace(input) == "" {
		return Translations{}, nil
	}

	// Set user options
	if len(options) != 0 {
		for _, option := range options {
			option(trans)
		}
	}

	return DoTranslate(input)
}

func DoTranslate(text string) (Translations, error) {
	return Translations{
		[]Translation{
			{
				Text:            "Translations text: " + text,
				DetectedSrcLang: "JA",
			},
		},
	}, nil
}

// ----------------------------------------------------------------------------
//  Type: Translations
// ----------------------------------------------------------------------------

// Translation is the type of the response from the DeepL translation API.
type Translations struct {
	Translations []Translation `json:"translations"`
}

// Translation
type Translation struct {
	Text            string `json:"text"`
	DetectedSrcLang string `json:"detected_source_language"`
}

func NewTranslations(responseJSON []byte) (*Translations, error) {
	translations := new(Translations)

	if err := json.Unmarshal(responseJSON, translations); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal JSON")
	}

	if translations.Translations == nil {
		return nil, errors.Errorf(
			"given response is not in the expected JSON format: %v", string(responseJSON))
	}

	return translations, nil
}
