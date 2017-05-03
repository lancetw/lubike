package setting

import (
	"fmt"
	"testing"
)

func TestInitConfigWithoutGopath(t *testing.T) {
	origLogFatalf := logFatalf
	origGopath := gopath

	defer func() { logFatalf = origLogFatalf }()
	defer func() { gopath = origGopath }()

	var errors []string
	logFatalf = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			errors = append(errors, format)
		}
	}

	gopath = ""

	InitConfig()

	if len(errors) != 0 {
		t.Fail()
	}
}

func TestInitConfig(t *testing.T) {
	origLogFatalf := logFatalf
	origGopath := gopath

	defer func() { logFatalf = origLogFatalf }()
	defer func() { gopath = origGopath }()

	var errors []string
	logFatalf = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			errors = append(errors, format)
		}
	}

	InitConfig()

	if len(errors) != 0 {
		t.Fail()
	}
}
