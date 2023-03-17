package deepl_test

import (
	"fmt"

	"github.com/KEINOS/go-deepl/deepl"
)

// ============================================================================
//  Error Handling Functions
// ============================================================================
//  Note: Be careful that the error output contains the file name and line number
//  of the caller. It requires the numbers to be static. If you change the code,
//  you may need to update the "Output:" of each Example test accordingly.

func ExampleNewErr() {
	if err := deepl.NewErr(); err == nil {
		fmt.Println("empty args returns nil")
	}

	// Note the output contains the file name and line number of the caller.
	if err := deepl.NewErr("simple error message at line 22"); err != nil {
		fmt.Println(err)
	}

	if err := deepl.NewErr("%v error message at line 26", "formatted"); err != nil {
		fmt.Println(err)
	}

	if err := deepl.NewErr("%v error message(s) at line 30", 3); err != nil {
		fmt.Println(err)
	}

	if err := deepl.NewErr(1, 2, 3); err != nil {
		fmt.Println(err)
	}
	// Output:
	// empty args returns nil
	// simple error message at line 22 (file: examples_test.go, line: 22)
	// formatted error message at line 26 (file: examples_test.go, line: 26)
	// 3 error message(s) at line 30 (file: examples_test.go, line: 30)
	// 1 2 3 (file: examples_test.go, line: 34)
}

func ExampleWrapIfErr() {
	var err error

	// deepl.WrapIfErr returns nil if err is nil
	fmt.Println("err is nil:", deepl.WrapIfErr(err, "error at line 49"))

	// Cause err to be non-nil
	err = deepl.NewErr("error occurred at line 52")
	// Wrap with no additional message
	fmt.Println("err is non-nil:\n", deepl.WrapIfErr(err))
	// Wrap with additional message
	fmt.Println("err is non-nil:\n", deepl.WrapIfErr(err, "wrapped at line 56"))
	// Output:
	// err is nil: <nil>
	// err is non-nil:
	//  error occurred at line 52 (file: examples_test.go, line: 52)
	// err is non-nil:
	//  wrapped at line 56 (file: examples_test.go, line: 56): error occurred at line 52 (file: examples_test.go, line: 52)
}

func ExampleWrapIfErr_disable_error_position() {
	// Backup and defer restore the original value of deepl.AppendErrPos
	oldAppendErrPos := deepl.AppendErrPos
	defer func() {
		deepl.AppendErrPos = oldAppendErrPos
	}()

	{
		deepl.AppendErrPos = false // Disable appending the error position

		err := deepl.NewErr("error occurred at line 75")
		fmt.Println(deepl.WrapIfErr(err, "wrapped at line 76"))
	}
	{
		deepl.AppendErrPos = true // Enable appending the error position (default)

		err := deepl.NewErr("error occurred at line 81")
		fmt.Println(deepl.WrapIfErr(err, "wrapped at line 82"))
	}
	// Output:
	// wrapped at line 76: error occurred at line 75
	// wrapped at line 82 (file: examples_test.go, line: 82): error occurred at line 81 (file: examples_test.go, line: 81)
}

// ============================================================================
//  The below examples do not require the line be static. So feel free to change
//  or refactor the code.
// ============================================================================

// ----------------------------------------------------------------------------
//  Type: APIType
// ----------------------------------------------------------------------------

//nolint:dupword // duplication in Output is intentional
func ExampleAPIType() {
	fmt.Println(deepl.APIFree)
	fmt.Println(deepl.APIFree.BaseURL())

	fmt.Println(deepl.APIPro)
	fmt.Println(deepl.APIPro.BaseURL())
	// Output:
	// https://api-free.deepl.com
	// https://api-free.deepl.com
	// https://api.deepl.com
	// https://api.deepl.com
}

//nolint:dupword // duplication in Output is intentional
func ExampleAPIType_custom_url() {
	// Force to use the custom URL
	deepl.SetCustomURL("http://localhost:8080")

	defer func() {
		deepl.SetCustomURL("")
	}()

	fmt.Println(deepl.APIFree)
	fmt.Println(deepl.APIPro)
	fmt.Println(deepl.APICustom)
	// Output:
	// http://localhost:8080
	// http://localhost:8080
	// http://localhost:8080
}

// ----------------------------------------------------------------------------
//  New
// ----------------------------------------------------------------------------

func ExampleNew() {
	// Create a new client for the free account of DeepL API
	// and use the default logger (stderr) by passing nil.
	cli, err := deepl.New(deepl.APIFree, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(cli.BaseURL)
	// Output: https://api-free.deepl.com
}
