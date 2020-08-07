package e5e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
)

// The struct for e5e event instances. Contains all fields but `data`, as the user code is expected to encapsulate this
// struct within it's own struct containing the `data` definition.
type Event struct {
	Params         map[string][]string `json:"params,omitempty"`
	RequestHeaders map[string]string   `json:"request_headers,omitempty"`
	Type           string              `json:"type,omitempty"`
}

// The struct for e5e context instances.
type Context struct {
	Async bool   `json:"async,omitempty"`
	Date  string `json:"date,omitempty"`
	Type  string `json:"type,omitempty"`
}

// The struct for the result value of an entrypoint function.
type Return struct {
	Status          int               `json:"status,omitempty"`
	ResponseHeaders map[string]string `json:"response_headers,omitempty"`
	Data            interface{}       `json:"data,omitempty"`
	Type            string            `json:"type,omitempty"`
}

// The struct for the internal e5e output representation.
type output struct {
	Output string      `json:"output"`
	Result interface{} `json:"result"`
}

// We neet to mock `os.Exit` and `fmt.Fprint` for the tests, so we save the function references here to be able
// to replace the function references with mock versions on the tests.
var (
	osExit    = os.Exit
	fmtFprint = fmt.Fprint
)

// This function takes the struct containing the available entrypoint methods and handles the invocation of the
// entrypoint as well as the communication with the e5e platform itself. Does not return if everything went well but
// does return an error if there are errors in the invocation or the entrypoint signature.
//
// Rules:
//	* Entrypoint functions must take 2 input parameters (Event and Context). Both types may be encapsulated within
//	  an user defined struct type.
//	* Entrypoint functions must return 2 values (Result and error). Type encapsulation is also allowed here.
//	* The input parameters as well as the return values must be compatible with "encoding/json" standard library.
func Start(entrypoints interface{}) error {
	// An e5e custom runtime binary must get 4 arguments when called by the platform. That is the binary name itself,
	// the name of the entrypoint function, the event object as JSON and the context object as JSON. We check for
	// the expected number of arguments here.
	if len(os.Args) != 4 {
		return fmt.Errorf("invalid number of process arguments")
	}

	// From the given struct that contains the possible entrypoint functions, we get the one with the name given
	// by the argument.
	entrypoint := reflect.ValueOf(entrypoints).MethodByName(os.Args[1])

	// Then we check if the entrypoint method we got represents a valid method and has the expected signature.
	if !entrypoint.IsValid() {
		return fmt.Errorf("invalid entrypoint name")
	}
	if entrypoint.Type().NumIn() != 2 {
		return fmt.Errorf("invalid number of entrypoint parameters")
	}
	if entrypoint.Type().NumOut() != 2 {
		return fmt.Errorf("invalid number of entrypoint return values")
	}

	// As we are now as sure as we can get that the method signature is the expected one, we receive the type
	// information of the first and the second parameter, and create a new references to instances of those types.
	eventType := entrypoint.Type().In(0)
	contextType := entrypoint.Type().In(1)
	event := reflect.New(eventType).Interface()
	context := reflect.New(contextType).Interface()

	// Next we try to load the JSON data of the event and context arguments into the previousely created instances.
	if err := json.Unmarshal([]byte(os.Args[2]), event); err != nil {
		return fmt.Errorf("cannot apply event object to '%s' type", eventType.Name())
	}
	if err := json.Unmarshal([]byte(os.Args[3]), context); err != nil {
		return fmt.Errorf("cannot apply context object to '%s' type", contextType.Name())
	}

	// As we want to capture all prints on stdout while the entrypoint function is running, we set this capturing
	// up now. To do so, we save the original stdout in a temporary variable and set the os.Stdout file to a
	// new pipe object. Writes to this file will get collected in a temporary channel within a go routine.
	originalStdout := os.Stdout
	captureStdoutChannel := make(chan string)
	captureStdoutRead, captureStdoutWrite, _ := os.Pipe()
	os.Stdout = captureStdoutWrite

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, captureStdoutRead)
		captureStdoutChannel <- buf.String()
	}()

	// Everythin is set up now, so we can call the entrypoint function now.
	results := entrypoint.Call([]reflect.Value{
		reflect.ValueOf(event).Elem(),
		reflect.ValueOf(context).Elem(),
	})

	// The user code now run, so we can revert the stdout capturing.
	captureStdoutWrite.Close()
	os.Stdout = originalStdout

	// The second return value of the entrypoint function is the error value. If it is not nil, we ignore the
	// result, print the error to stderr and print the output structure to stdout.
	if results[1].Kind() != reflect.Interface || !results[1].IsNil() {
		if err, ok := results[1].Interface().(error); ok {
			fmtFprint(os.Stderr, fmt.Sprintf("%s", err))

			out := output{
				Output: <-captureStdoutChannel,
				Result: nil,
			}

			if marshaled, err := json.Marshal(out); err == nil {
				fmtFprint(os.Stdout, string(marshaled))
				osExit(-1)
				return nil
			} else {
				return fmt.Errorf("cannot marshal return value")
			}
		} else {
			return fmt.Errorf("invalid error return value")
		}
	}

	// No error was returned by the user code, so we go ahead and print the output structure to stdout.
	out := output{
		Output: <-captureStdoutChannel,
		Result: results[0].Interface(),
	}

	if marshaled, err := json.Marshal(out); err == nil {
		fmtFprint(os.Stdout, string(marshaled))
		osExit(0)
		return nil
	} else {
		return fmt.Errorf("cannot marshal return value")
	}
}
