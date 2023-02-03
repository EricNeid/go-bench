// Package verify provides utility functions for tests.
package verify

import (
	"encoding/json"
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
func Equals(t *testing.T, exp, act any) {
	t.Helper()

	if !deepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}

func deepEqual(v1, v2 interface{}) bool {
	if reflect.DeepEqual(v1, v2) {
		return true
	}
	var x1 interface{}
	bytesA, _ := json.Marshal(v1)
	_ = json.Unmarshal(bytesA, &x1)
	var x2 interface{}
	bytesB, _ := json.Marshal(v2)
	_ = json.Unmarshal(bytesB, &x2)
	return reflect.DeepEqual(x1, x2)
}

// NotNil fails if act is nil
func NotNil(t *testing.T, act interface{}, msg string) {
	t.Helper()

	if act == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: "+msg+"\n\n", filepath.Base(file), line)
		t.FailNow()
	}
}
