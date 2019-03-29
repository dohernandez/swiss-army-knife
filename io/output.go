package io

import (
	"context"
	"fmt"
	"os"
)

// Output defines a contract for output data target.
type Output interface {
	// Append adds output data.
	Append(ctx context.Context, output interface{})

	// Write writes the output data into the target.
	// Returns any error that occurred.
	Write(ctx context.Context) error
}

// MarshalOutput function to marshal an object.
//
// Returns error if marshal fails.
type MarshalOutput func(ctx context.Context, i interface{}) (string, error)

// StdoutOutput write the output data to the os.Stdout.
type StdoutOutput struct {
	marshalOutput MarshalOutput

	output []interface{}
}

var _ Output = new(StdoutOutput)

// Append adds output data to be printed to io.Stdout.
func (o *StdoutOutput) Append(_ context.Context, output interface{}) {
	o.output = append(o.output, output)
}

// Write writes the output into the io.Stdout.
//
// Returns any error that occurred.
func (o *StdoutOutput) Write(ctx context.Context) error {
	var newLine bool
	for _, out := range o.output {
		if newLine {
			fmt.Fprintf(os.Stdout, "\n")
		}

		if o.marshalOutput != nil {
			r, err := o.marshalOutput(ctx, out)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "%s", r)

			newLine = true

			continue
		}

		fmt.Fprintf(os.Stdout, "%v", out)

		newLine = true
	}

	return nil
}

// WithMarshaling set MarshalOutput func into StdoutOutput.
func (o *StdoutOutput) WithMarshaling(marshalOutput MarshalOutput) *StdoutOutput {
	o.marshalOutput = marshalOutput

	return o
}
