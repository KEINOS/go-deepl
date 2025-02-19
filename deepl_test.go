package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTranslator_no_option(t *testing.T) {
	require.NotPanics(t, func() {
		translator := NewTranslator()

		require.NotNil(t, translator,
			"NewTranslator() should return a Translator object")
	}, "it should not panic when no option is given")
}

func ExampleNewTranslator_golden() {
	translator := NewTranslator(
		WithAuthKey("myAuthKey"),
		WithAutomaticFormatting(true),
		WithFormality("formal"),
		WithGlossaryID("myGlossaryID"),
		WithIgnoreTags([]string{"tag1", "tag2"}),
		WithNonSplittingTags([]string{"tag3", "tag4"}),
		WithOutlineDetection(false),
		WithSourceLang("EN"),
		WithSplitSentences("nonewlines"),
		WithSplittingTags([]string{"tag5", "tag6"}),
		WithTagHandling("xml"),
		WithTargetLang("DE"),
	)

	for _, key := range []string{
		// Available property names
		"authKey",
		"formality",
		"glossaryID",
		"ignoreTags",
		"nonSplittingTags",
		"outlineDetection",
		"preserveFormatting",
		"sourceLang",
		"splitSentences",
		"splittingTags",
		"tagHandling",
		"targetLang",
	} {
		value := translator.Property(key)

		fmt.Printf("NewTranslator.Property(\"%s\") = %#v\n", key, value)
	}
	// Output:
	// NewTranslator.Property("authKey") = "myAuthKey"
	// NewTranslator.Property("formality") = "prefer_more"
	// NewTranslator.Property("glossaryID") = "myGlossaryID"
	// NewTranslator.Property("ignoreTags") = []string{"tag1", "tag2"}
	// NewTranslator.Property("nonSplittingTags") = []string{"tag3", "tag4"}
	// NewTranslator.Property("outlineDetection") = false
	// NewTranslator.Property("preserveFormatting") = true
	// NewTranslator.Property("sourceLang") = "EN"
	// NewTranslator.Property("splitSentences") = "nonewlines"
	// NewTranslator.Property("splittingTags") = []string{"tag5", "tag6"}
	// NewTranslator.Property("tagHandling") = "xml"
	// NewTranslator.Property("targetLang") = "DE"
}

func TestWithFormality(t *testing.T) {
	tmpTranslator := NewTranslator()

	for _, test := range []struct {
		level  string
		expect string
	}{
		{"formal", "prefer_more"},
		{"informal", "prefer_less"},
		{"default", "default"},
	} {
		fnOpt := WithFormality(test.level)
		fnOpt(tmpTranslator)

		expect := test.expect
		actual := tmpTranslator.Property("formality")

		assert.Equal(t, expect, actual,
			"WithFormality(\"%s\") did not set the expected value", test.level)
	}
}

func TestTranslator_Property_unknown(t *testing.T) {
	tmpTranslator := NewTranslator()

	require.Nil(t, tmpTranslator.Property("unknown_property"),
		"Undefined property should return nil")
}

func TestTranslator_IsKeyFree(t *testing.T) {
	for _, test := range []struct {
		authKey string
		expect  bool
	}{
		// Test value taken from:
		//   https://www.deepl.com/docs-api/api-access/api-versions/#authentication
		{"279a2e9d-83b3-c416-7e2d-f721593e42a0:fx", true},
		{"279a2e9d-83b3-c416-7e2d-f721593e42a0", false},
	} {
		tmpTranslator := NewTranslator(WithAuthKey(test.authKey))

		expect := test.expect
		actual := tmpTranslator.IsKeyFree()

		require.Equal(t, expect, actual,
			"auth key \"%s\" should be %v", test.authKey, expect)
	}
}

func ExampleNewTranslations() {
	// Example response taken from: https://www.deepl.com/docs-api
	responseJSON := []byte(`{
		"translations": [
			{
				"detected_source_language": "EN",
				"text": "Hallo, Welt!"
			}
	]}`)

	// Parse the JSON response.
	response, err := NewTranslations(responseJSON)
	if err != nil {
		log.Fatal(err)
	}

	// Print the first translation's text and detected source language.
	fmt.Println("Translation:", response.Translations[0].Text)
	fmt.Println("Detected source language:", response.Translations[0].DetectedSrcLang)
	// Output:
	// Translation: Hallo, Welt!
	// Detected source language: EN
}

func TestNewTranslations_malformed_json(t *testing.T) {
	responseJSON := []byte(`helllo:	world`)

	response, err := NewTranslations(responseJSON)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to unmarshal JSON",
		"error message should contain the error reason")
	require.Contains(t, err.Error(), "invalid character",
		"error message should contain the underlying error message")
	require.Nil(t, response,
		"returned response should be nil on error")
}

func TestNewTranslations_other_format_json(t *testing.T) {
	// JSON but not the expected format.
	responseJSON := []byte(`{
		"foo": [
			{
				"bar": "hoge",
				"buz": 1
			}
	]}`)

	response, err := NewTranslations(responseJSON)

	require.Error(t, err)
	require.Contains(t, err.Error(), "given response is not in the expected JSON format",
		"error message should contain the error reason")
	require.Contains(t, err.Error(), "foo",
		"error message should contain the original JSON")
	require.Nil(t, response,
		"returned response should be nil on error")
}
