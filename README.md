# Jettison
[![Go Report Card](https://goreportcard.com/badge/github.com/golang-standards/project-layout?style=flat-square)](https://goreportcard.com/report/github.com/luno/jettison)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/luno/jettison)

## What is it?
Jettison is a library providing structured logs and errors in a way that's
interoperable with gRPC. It does this under the hood by serialising message
details to protobuf messages and attaching them to gRPC `Status` objects - 
see [the gRPC docs](https://godoc.org/google.golang.org/grpc/status) for 
details of how this works. Jettison is also compatible with the Go 2 error
spec, which can be found [in the draft design page](https://go.googlesource.com/proposal/+/master/design/go2draft.md).

## Features
Jettison is in alpha, but the following features are planned:

* [✓] Simple, gRPC-compatible utilities for building up structured errors/logs.
* [✓] Support for error identification and unwrapping as per the Go 2 spec.
* [✓] Structured error/log formatting utilities for both machines (JSON) and 
      humans.
* [✕] Key/value pair type support for compatibility with Elasticsearch.

## API
### Errors
The `jettison/errors` package provides functions for creating and working with
gRPC-compatible error values. You can attach arbitrary metadata to jettison
errors, decorate them with source/stacktrace information and wrap errors to
form a chain as the stack unwinds. Passing jettison errors over gRPC preserves 
all metadata structure, in contrast to other error types (which get marshalled 
to a string by default).

Jettison also provides gRPC middleware that automatically groups the errors 
in a chain by the gRPC server (or "hop") they originated from.

See the `jettison/example` package for a more complete usage example, including
a gRPC server/client passing jettison errors over the wire.

```GO
import (
    "github.com/luno/jettison"
    "github.com/luno/jettison/errors"
)

func Example() error {
	// Construct your error as usual, with additional metadata.
	err := errors.New("something went wrong",
		jettison.WithKeyValueString("key", "value"),
		jettison.WithSource("Example()"))

	// Wrap errors with additional metadata as they get passed down the stack.
	err = errors.Wrap(err, "something else went wrong",
		jettison.WithKeyValueString("another_key", "another_value"))

	// Pass it around - including over gRPC - like you would any other error.
	return err
}
```

### Logs
The `jettison/log` package provides structured JSON logging, with additional
utilities for logging jettison errors. You can attach metadata to logs in the
same manner as you attach metadata to errors.

```GO
import (
    "context"

    "github.com/luno/jettison"
    "github.com/luno/jettison/errors"
    "github.com/luno/jettison/log"
)

func Example(ctx context.Context) {
	// You can log general info as you normally would.
	log.Info(ctx, "entering the example function",
		jettison.WithKeyValueString("key", "value"))

	err := errors.New("a jettison error",
		jettison.WithKeyValueString("key", "value"),
		jettison.WithSource("Example()"))

	// Errors can be logged separately, with metadata marshalled to JSON in
	// a machine-friendly manner.
	log.Error(ctx, err,
		jettison.WithKeyValueString("another_key", "another_value"))
}
```

An example log written via `log.Info`:
```JSON
{
  "message": "entering the example function",
  "source": "jettison/example/example.go:9",
  "level": "info",
  "parameters": [
    {
      "key": "key",
      "value": "value"
    }
  ]
}
```

An example log written via `log.Error`:
```JSON
{
  "message": "a jettison error",
  "source": "jettison/example/example.go:18",
  "level": "error",
  "hops": [
    {
      "binary": "example",
      "errors": [
        {
          "message": "a jettison error",
          "source": "jettison/example/example.go:14",
          "parameters": [
            {
              "key": "key",
              "value": "value"
            }
          ]
        }
      ]
    }
  ],
  "parameters": [
    {
      "key": "key",
      "value": "value"
    },
    {
      "key": "another_key",
      "value": "another_value"
    }
  ]
}
```

### Utilities
The `jettison/j` package contains aliases for common jettison options,
saving you a couple of keystrokes:

```GO
import (
    "github.com/luno/jettison/errors"
    "github.com/luno/jettison/j"
)

func Example() error {
	return errors.New("using the aliases",
		j.KV("int_key", 1),
		j.KS("string_key", "value"))
}
```
