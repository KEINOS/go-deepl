# go-deepl

[go-deepl](https://github.com/KEINOS/go-deepl) is a simple Go library for DeepL API client.

> __Note__: It is a fork from [deepl-go](https://github.com/shopper29/deepl-go) by [shopper29](https://github.com/shopper29/) with some modifications. Such as security updates, replacing deprecated modules, code coverage, etc.

## Usage

```go
go get github.com/KEINOS/go-deepl
```

```go
import "github.com/KEINOS/go-deepl/deepl"
```

### Requirements

- You need an account of [DeepL API Free or Pro](https://www.deepl.com/pro#developer).
- The environment variable `DEEPL_API_KEY` and a valid API key ("Authentication Key for DeepL API" from [your account settings](https://www.deepl.com/account/summary)) must be set.

## Examples

```go
package main

import (
    "context"
    "fmt"

    "github.com/KEINOS/go-deepl/deepl"
)

func main() {
    // Create a client for free account of DeepL API (choices: deepl.APIFree,
    // deepl.APIPro, deepl.APICustom). The second arg is the logger. If nil,
    // the default logger is used. Which logs to stderr.
    cli, err := deepl.New(deepl.APIFree, nil)
    if err != nil {
        log.Fatal(err)
    }

    translateResponse, err := cli.TranslateSentence(
        context.Background(),
        "Hello", // Phrase to translate
        "EN",    // from English
        "JA",    // to Japanese
    )

    if err != nil {
        log.Fatal(err)
    } else {
        fmt.Printf("%+v\n", translateResponse)
    }
    // Output:
    // &{Translations:[{DetectedSourceLanguage:EN Text:こんにちは}]}
}
```

```console
```

## License and Authors

- [MIT License](https://github.com/KEINOS/go-deepl/blob/main/LICENSE.md). Copyright (c) 2023 [shopper29](https://github.com/shopper29/), [KEINOS and the go-deepl contributors](https://github.com/KEINOS/go-deepl/graphs/contributors).
