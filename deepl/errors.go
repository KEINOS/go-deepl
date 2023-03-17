package deepl

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  This file contains a simple error handling functions.
//
//  Basically, it is an error wrapper but it also appends the file name and line
//  number of the caller to the error message.
// ----------------------------------------------------------------------------

// AppendErrPos is a global flag to disable the file name and line number of the
// caller from the error message. If set to false, it will not be appended.
var AppendErrPos = true

// NewErr returns a new error object with the given message appending the file
// name and line number of the caller.
//
// It is a wrapper of errors.New() and errors.Errorf(). Which is the alternative
// of deprecated github.com/pkg/errors.
func NewErr(msgs ...interface{}) error {
	lenMsgs := len(msgs)
	if lenMsgs == 0 {
		return nil
	}

	errPos := getErrorPos()

	fmtErr, ok := msgs[0].(string)
	if !ok {
		errMsg := fmtArgs(msgs[:]...) //nolint:gocritic // false positive

		return errors.New(errMsg + errPos)
	}

	fmtErr += errPos

	if lenMsgs == 1 {
		return errors.New(fmtErr)
	}

	return errors.Errorf(fmtErr, msgs[1:]...)
}

// WrapIfErr returns nil if err is nil.
//
// Otherwise, it returns an error annotating err with a stack trace at the point
// WrapIfErr is called. The supplied message contains the file name and line
// number of the caller.
//
// Note that if the "msgs" arg is more than one, the first arg is used as a
// format string and the rest are used as arguments.
//
// E.g.
//
//	WrapIfErr(nil, "it wil do nothing")
//	WrapIfErr(err)                                 // returns err as is
//	WrapIfErr(err, "failed to do something")       // eq to errors.Wrap
//	WrapIfErr(err, "failed to do %s", "something") // eq to errors.Wrapf
//
// It is a wrapper of errors.Wrap() and errors.Wrapf(). Which is the alternative
// of deprecated github.com/pkg/errors.
func WrapIfErr(err error, msgs ...interface{}) error {
	if err == nil {
		return nil
	}

	if len(msgs) == 0 {
		return fmt.Errorf("%w", err)
	}

	errMsg := fmtArgs(msgs...)
	errPos := getErrorPos()

	return errors.Wrap(err, errMsg+errPos)
}

// fmtArgs is a shorthand for fmt.Sprintf and fmt.Sprint arguments. It is a helper
// function to format the given arguments.
//
// If the inputs is empty, it returns an empty string.
// If the inputs has only one element, it returns the string representation of
// the element.
// If the inputs has more than one element, the first element is used as a
// format string and the rest are used as arguments.
func fmtArgs(inputs ...interface{}) string {
	lenInput := len(inputs)

	if lenInput == 0 {
		return ""
	}

	if lenInput == 1 {
		return fmt.Sprint(inputs[0])
	}

	format, ok := inputs[0].(string)
	if !ok {
		return fmt.Sprint(inputs[:]...) //nolint:gocritic // false positive
	}

	return fmt.Sprintf(format, inputs[1:]...)
}

// getErrorPos returns a string containing the file name and line number of the
// caller.
func getErrorPos() string {
	grandparent := 2 // 0 = self, 1 = parent, 2 = grandparent

	_, file, line, ok := runtime.Caller(grandparent)
	if !ok || !AppendErrPos {
		return ""
	}

	return fmt.Sprintf(" (file: %s, line: %d)", filepath.Base(file), line)
}
