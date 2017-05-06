package setting

import (
	"errors"
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

func TestInitConfigError(t *testing.T) {
	origTestingMode := testingMode
	testingMode = ""
	origLogFatalf := logFatalf
	origGodotenvLoad := godotenvLoad

	defer func() { logFatalf = origLogFatalf }()
	defer func() { godotenvLoad = origGodotenvLoad }()
	defer func() { testingMode = origTestingMode }()

	var errorsInfo []string
	logFatalf = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errorsInfo = append(errorsInfo, fmt.Sprintf(format, args))
		} else {
			errorsInfo = append(errorsInfo, format)
		}
	}

	godotenvLoad = func(...string) error {
		return errors.New("err")
	}

	InitConfig()

	if len(errorsInfo) == 0 {
		t.Fail()
	}
}
