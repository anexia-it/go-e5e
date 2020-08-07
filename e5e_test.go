package e5e

import (
	"fmt"
	"io"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type SumData struct {
	A int `json:"a"`
	B int `json:"b"`
}

type SumEvent struct {
	Event
	Data SumData `json:"data"`
}

type testfunction func()
type entrypoints struct{}

func (f *entrypoints) SimpleEntrypoint(event Event, context Context) (*Return, error) {
	return nil, nil
}

func (f *entrypoints) SumEntrypoint(event SumEvent, context Context) (*Return, error) {
	return &Return{
		Data: event.Data.A + event.Data.B,
	}, nil
}

func (f *entrypoints) PrintStdOutEntrypoint(event Event, context Context) (*Return, error) {
	fmt.Print("print")
	return nil, nil
}

func (f *entrypoints) ErrorEntrypoint(event Event, context Context) (*Return, error) {
	return nil, fmt.Errorf("error")
}

func (f *entrypoints) InvalidParametersEntrypoint() (*Return, error) {
	return nil, nil
}

func (f *entrypoints) InvalidReturnEntrypoint(event Event, context Context) {
	return
}

func (f *entrypoints) InvalidReturnValueEntrypoint(event Event, context Context) (*Return, error) {
	return &Return{
		Data: math.Inf(1),
	}, nil
}

func (f *entrypoints) InvalidErrorReturnValueEntrypoint(event Event, context Context) (*Return, int) {
	return nil, 1
}

func TestStartSimpleEntrypoint(t *testing.T) {
	stdOut, stdErr, exitCode := runMocked(func() {
		if err := Start(&entrypoints{}); err != nil {
			fmtFprint(os.Stderr, fmt.Sprintf("%s", err))
			osExit(-255)
		}
	}, "SimpleEntrypoint", "{}", "{}")

	require.EqualValues(t, "{\"output\":\"\",\"result\":null}", stdOut)
	require.EqualValues(t, "", stdErr)
	require.EqualValues(t, 0, exitCode)
}

func TestStartInvalidSimpleEntrypoint(t *testing.T) {
	stdOut, stdErr, exitCode := runMocked(func() {
		if err := Start(&entrypoints{}); err != nil {
			fmtFprint(os.Stderr, fmt.Sprintf("%s", err))
			osExit(-255)
		}
	}, "InvalidSimpleEntrypointt", "{}", "{}")

	require.EqualValues(t, "", stdOut)
	require.EqualValues(t, "invalid entrypoint name", stdErr)
	require.EqualValues(t, -255, exitCode)
}

func TestStartSumEntrypoint(t *testing.T) {
	stdOut, stdErr, exitCode := runMocked(func() {
		if err := Start(&entrypoints{}); err != nil {
			fmtFprint(os.Stderr, fmt.Sprintf("%s", err))
			osExit(-255)
		}
	}, "SumEntrypoint", "{\"data\":{\"a\": 2, \"b\": 3}}", "{}")

	require.EqualValues(t, "{\"output\":\"\",\"result\":{\"data\":5}}", stdOut)
	require.EqualValues(t, "", stdErr)
	require.EqualValues(t, 0, exitCode)
}

func TestStartPrintStdOutEntrypoint(t *testing.T) {
	stdOut, stdErr, exitCode := runMocked(func() {
		if err := Start(&entrypoints{}); err != nil {
			fmtFprint(os.Stderr, fmt.Sprintf("%s", err))
			osExit(-255)
		}
	}, "PrintStdOutEntrypoint", "{}", "{}")

	require.EqualValues(t, "{\"output\":\"print\",\"result\":null}", stdOut)
	require.EqualValues(t, "", stdErr)
	require.EqualValues(t, 0, exitCode)
}

func TestStartErrorEntrypoint(t *testing.T) {
	stdOut, stdErr, exitCode := runMocked(func() {
		if err := Start(&entrypoints{}); err != nil {
			fmtFprint(os.Stderr, fmt.Sprintf("%s", err))
			osExit(-255)
		}
	}, "ErrorEntrypoint", "{}", "{}")

	require.EqualValues(t, "{\"output\":\"\",\"result\":null}", stdOut)
	require.EqualValues(t, "error", stdErr)
	require.EqualValues(t, -1, exitCode)
}

func TestStartInvalidArgumentsSimpleEntrypoint(t *testing.T) {
	stdOut, stdErr, exitCode := runMocked(func() {
		if err := Start(&entrypoints{}); err != nil {
			fmtFprint(os.Stderr, fmt.Sprintf("%s", err))
			osExit(-255)
		}
	}, "SimpleEntrypoint")

	require.EqualValues(t, "", stdOut)
	require.EqualValues(t, "invalid number of process arguments", stdErr)
	require.EqualValues(t, -255, exitCode)
}

func TestStartInvalidEventSimpleEntrypoint(t *testing.T) {
	stdOut, stdErr, exitCode := runMocked(func() {
		if err := Start(&entrypoints{}); err != nil {
			fmtFprint(os.Stderr, fmt.Sprintf("%s", err))
			osExit(-255)
		}
	}, "SimpleEntrypoint", "{...}", "{}")

	require.EqualValues(t, "", stdOut)
	require.EqualValues(t, "cannot apply event object to 'Event' type", stdErr)
	require.EqualValues(t, -255, exitCode)
}

func TestStartInvalidContextSimpleEntrypoint(t *testing.T) {
	stdOut, stdErr, exitCode := runMocked(func() {
		if err := Start(&entrypoints{}); err != nil {
			fmtFprint(os.Stderr, fmt.Sprintf("%s", err))
			osExit(-255)
		}
	}, "SimpleEntrypoint", "{}", "{...}")

	require.EqualValues(t, "", stdOut)
	require.EqualValues(t, "cannot apply context object to 'Context' type", stdErr)
	require.EqualValues(t, -255, exitCode)
}

func TestStartInvalidParametersEntrypoint(t *testing.T) {
	stdOut, stdErr, exitCode := runMocked(func() {
		if err := Start(&entrypoints{}); err != nil {
			fmtFprint(os.Stderr, fmt.Sprintf("%s", err))
			osExit(-255)
		}
	}, "InvalidParametersEntrypoint", "{}", "{}")

	require.EqualValues(t, "", stdOut)
	require.EqualValues(t, "invalid number of entrypoint parameters", stdErr)
	require.EqualValues(t, -255, exitCode)
}

func TestStartInvalidReturnEntrypoint(t *testing.T) {
	stdOut, stdErr, exitCode := runMocked(func() {
		if err := Start(&entrypoints{}); err != nil {
			fmtFprint(os.Stderr, fmt.Sprintf("%s", err))
			osExit(-255)
		}
	}, "InvalidReturnEntrypoint", "{}", "{}")

	require.EqualValues(t, "", stdOut)
	require.EqualValues(t, "invalid number of entrypoint return values", stdErr)
	require.EqualValues(t, -255, exitCode)
}

func TestStartInvalidReturnValueEntrypoint(t *testing.T) {
	stdOut, stdErr, exitCode := runMocked(func() {
		if err := Start(&entrypoints{}); err != nil {
			fmtFprint(os.Stderr, fmt.Sprintf("%s", err))
			osExit(-255)
		}
	}, "InvalidReturnValueEntrypoint", "{}", "{}")

	require.EqualValues(t, "", stdOut)
	require.EqualValues(t, "cannot marshal return value", stdErr)
	require.EqualValues(t, -255, exitCode)
}

func TestStartInvalidErrorReturnValueEntrypoint(t *testing.T) {
	stdOut, stdErr, exitCode := runMocked(func() {
		if err := Start(&entrypoints{}); err != nil {
			fmtFprint(os.Stderr, fmt.Sprintf("%s", err))
			osExit(-255)
		}
	}, "InvalidErrorReturnValueEntrypoint", "{}", "{}")

	require.EqualValues(t, "", stdOut)
	require.EqualValues(t, "invalid error return value", stdErr)
	require.EqualValues(t, -255, exitCode)
}

func runMocked(f testfunction, args ...string) (stdOut string, stdErr string, exitCode int) {
	// Mock `os.Exit`
	exitCode = 1024
	origOsExit := osExit
	defer func() {
		osExit = origOsExit
	}()
	osExit = func(code int) {
		exitCode = code
	}

	// Mock `fmt.Fprint`
	stdOut = ""
	stdErr = ""
	origFmtFprint := fmtFprint
	defer func() {
		fmtFprint = origFmtFprint
	}()
	fmtFprint = func(w io.Writer, a ...interface{}) (int, error) {
		switch w {
		case os.Stdout:
			stdOut += fmt.Sprint(a...)
		case os.Stderr:
			stdErr += fmt.Sprint(a...)
		}
		return 0, nil
	}

	// Mock `os.Args`
	origOsArgs := os.Args
	defer func() {
		os.Args = origOsArgs
	}()
	os.Args = append([]string{"cmd"}, args...)

	f()

	return stdOut, stdErr, exitCode
}
