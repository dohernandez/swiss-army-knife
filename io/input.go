package io

import (
	"bufio"
	"context"
	"io"
)

// Input defines a contract for input data source.
type Input interface {
	// Next returns the next record of the input source.
	// Starting from the first record when it is call the first time.
	// Returns any error that occurred, including io.EOF when no more record is available.
	Next(ctx context.Context) (interface{}, error)
}

// UnmarshalInput function to unmarshal a stream object.
//
// Returns error if unmarshal fails
type UnmarshalInput func(ctx context.Context, i string) (interface{}, error)

// StdinInput reads the input data coming from os.Stdin.
type StdinInput struct {
	scanner        *bufio.Scanner
	unmarshalInput UnmarshalInput
}

var _ Input = new(StdinInput)

// NewStdinInput create an instance of StdinInput.
func NewStdinInput(scanner *bufio.Scanner) *StdinInput {
	return &StdinInput{
		scanner: scanner,
	}
}

// Next returns the next record of the io.Stdin. If unmarshalInput is set, the record will be unmarshaled.
// Starting from the first record when it is call the first time.
//
// Returns any error that occurred, including io.EOF when no more record is available.
func (i *StdinInput) Next(ctx context.Context) (interface{}, error) {
	if !i.scanner.Scan() {
		return nil, io.EOF
	}

	if err := i.scanner.Err(); err != nil {
		return nil, err
	}

	if i.unmarshalInput != nil {
		return i.unmarshalInput(ctx, i.scanner.Text())
	}

	return i.scanner.Text(), nil
}

// WithUnmarshaling set UnmarshalInput func into StdinInput.
func (i *StdinInput) WithUnmarshaling(unmarshalInput UnmarshalInput) *StdinInput {
	i.unmarshalInput = unmarshalInput

	return i
}
