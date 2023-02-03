// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: MIT

// Package verify provides utility functions for tests.
package verify

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Assert fails the test if the condition is false.
func Assert(t *testing.T, condition bool, msg string, v ...interface{}) {
	t.Helper()
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: "+msg+"\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		t.FailNow()
	}
}

// Ok fails the test if an err is not nil.
func Ok(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: unexpected error: %s\n\n", filepath.Base(file), line, err.Error())
		t.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(t *testing.T, exp, act interface{}) {
	t.Helper()
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}

// NotNil fails if act is nil
func NotNil(t *testing.T, act interface{}, msg string) {
	t.Helper()
	if act == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: "+msg+"\n\n", append([]interface{}{filepath.Base(file), line}))
		t.FailNow()
	}
}
