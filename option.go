package main

// Option is a functional option pattern. It holds the function that sets the
// options to the given Translator.
type Option func(*Translator)

// ----------------------------------------------------------------------------
//  Functional Option Pattern
// ----------------------------------------------------------------------------

// WithAuthKey is a required option to specify the authentication key.
func WithAuthKey(key string) func(*Translator) {
	return func(translator *Translator) {
		translator.authKey = key
	}
}

// WithAutomaticFormatting is an option to specify whether the translation engine
// corrects text formatting automatically. Default is false ("off").
func WithAutomaticFormatting(automaticFormat bool) func(*Translator) {
	return func(translator *Translator) {
		translator.preserveFormatting = automaticFormat
	}
}

// WithDefault is an option that sets the default values.
func WithDefault() func(*Translator) {
	return func(translator *Translator) {
		translator.targetLang = ""
		translator.sourceLang = ""
		translator.splitSentences = "on"
		translator.formality = "default"
		translator.glossaryID = ""
		translator.tagHandling = ""
		translator.nonSplittingTags = []string{}
		translator.outlineDetection = true
		translator.preserveFormatting = false
	}
}

// WithFormality is an option to specify the formality of the translation if available,
// otherwise fallback to default formality. Default is "default". Possible values are:
//   - "default":
//     Default behavior of the translation engine.
//   - "informal":
//     For a less formal language.
//   - "formal":
//     For a more formal language.
func WithFormality(level string) func(*Translator) {
	switch level {
	case "formal":
		return func(translator *Translator) {
			translator.formality = "prefer_more"
		}
	case "informal":
		return func(translator *Translator) {
			translator.formality = "prefer_less"
		}
	default:
		return func(translator *Translator) {
			translator.formality = "default"
		}
	}
}

// WithGlossaryID is an option to specify the glossary to use for the translation.
//
// Note that this option requires that the WithSourceLang and WithTargetLang
// options are explicitly defined and match the language pair in the glossary.
func WithGlossaryID(id string) func(*Translator) {
	return func(translator *Translator) {
		translator.glossaryID = id
	}
}

// WithIgnoreTags is an option that specifies a list of XML tags which never be
// translated.
func WithIgnoreTags(ignoreTags []string) func(*Translator) {
	return func(translator *Translator) {
		translator.ignoreTags = ignoreTags
	}
}

// WithNonSplittingTags is an option that specifies a list of XML tags which never
// be used to split sentences.
//
// Some XML files do not provide the best results when searching for tags that
// contain text content or when splitting text based on those tags. In that case
// the indicated tags here are not used in the splitting statement.
func WithNonSplittingTags(tags []string) func(*Translator) {
	return func(translator *Translator) {
		translator.nonSplittingTags = tags
	}
}

// WithOutlineDetection is an option to specify automatic detection of the XML
// structure. Default is true ("on", "1").
//
// Disabling this option and manually setting parameters like WithTagHandling,
// WithSentenceSplitting, and WithSplittingTags options provides greater control
// over translation output.
func WithOutlineDetection(autoDetect bool) func(*Translator) {
	return func(translator *Translator) {
		translator.outlineDetection = autoDetect
	}
}

// WithSourceLang is an option to specify the source language code.
// By default is "auto-detect".
func WithSourceLang(lang string) func(*Translator) {
	return func(translator *Translator) {
		translator.sourceLang = lang
	}
}

// WithSplitSentences is an option to specify how input text should be split into
// sentences. Default is "on". Possible values are:
//   - "on":
//     Input text will be split into sentences using both newlines and punctuation.
//   - "off":
//     Input text will not be split into sentences. Use this for applications
//     where each input text contains only one sentence.
//   - "nonewlines":
//     Input text will be split into sentences using punctuation but not newlines.
func WithSplitSentences(splitType string) func(*Translator) {
	return func(translator *Translator) {
		translator.splitSentences = splitType
	}
}

// WithSplittingTags is an option that specifies a list of XML tags which always
// cause splits.
func WithSplittingTags(tags []string) func(*Translator) {
	return func(translator *Translator) {
		translator.splittingTags = tags
	}
}

// WithTagHandling is an option to set which kind of tags should be handled.
// Options currently available:
//   - "xml": Enable XML tag handling.
//   - "html": Enable HTML tag handling.
func WithTagHandling(tagType string) func(*Translator) {
	return func(translator *Translator) {
		translator.tagHandling = tagType
	}
}

// WithTargetLang is a required option to pecify the target language to translate
// the text into.
func WithTargetLang(lang string) func(*Translator) {
	return func(translator *Translator) {
		translator.targetLang = lang
	}
}
