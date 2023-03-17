package deepl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ============================================================================
//  Tests
// ============================================================================

func TestNewErr_no_args(t *testing.T) {
	t.Parallel()

	require.Nil(t, NewErr(), "empty args should return nil")
}

func Test_fmtArgs(t *testing.T) {
	t.Parallel()

	require.Empty(t, fmtArgs(),
		"no args should return empty string")

	require.Equal(t, "10 lines found", fmtArgs("%v lines found", 10),
		"if the first arg is a string, it should be formatted with the rest of the args")

	require.Equal(t, "1 2 3", fmtArgs(1, 2, 3),
		"if the first arg is not a string, it should return the concatenation of all args")
}

//nolint:paralleltest // do not parallelize due to global variable change
func Test_getErrorPos(t *testing.T) {
	// Backup and defer restore the original value of AppendErrPos
	oldAppendErrPos := AppendErrPos
	defer func() {
		AppendErrPos = oldAppendErrPos
	}()

	//nolint:gocritic
	fn1 := func() string {
		return getErrorPos() // capture line number of the caller
	}

	fn2 := func() string {
		return fn1() // should return line 22
	}

	t.Run("disable append", func(t *testing.T) {
		// Disable append the error position info
		AppendErrPos = false

		result := fn2()

		require.Empty(t, result, "setting AppendErrPos to false should return empty string")
	})

	t.Run("enable append (default)", func(t *testing.T) {
		// Enable append the error position info
		AppendErrPos = true

		result := fn2()

		require.Contains(t, result, "file: errors_test.go",
			"on AppendErrPos=true, the caller's file name should be included")
		require.Contains(t, result, "line:",
			"on AppendErrPos=true, the caller's line number should be included")
	})
}
