package technical_test

import (
	"bytes"
	"encoding/json"
	"io"
)

// Conveyor interface defines contract for conveying data between operations.
type Conveyor interface {
	Accept(v interface{}) error
	Emit(v interface{}) error
}

// ChannelConveyor interface defines contract for conveyor data between operations using channels.
type ChannelConveyor interface {
	Conveyor

	ChainNext() ChannelConveyor

	Close()
}

type channelConveyor struct {
	inputCh  chan interface{}
	outputCh chan interface{}
}

// NewChannelConveyor creates new conveyor that conveys values with channels.
//
// Common initialization example:
//
//      inputs := make(chan interface{})
//
//		// create a ChannelConveyor.
//		cc := NewChannelConveyor(inputs)
//
func NewChannelConveyor(input chan interface{}) ChannelConveyor {
	return channelConveyor{
		inputCh: input,
		// to limit the amount of work that is queued up.
		outputCh: make(chan interface{}, 1024),
	}
}

// Close closes the channel sending an io.EOF signal.
func (c channelConveyor) Close() {
	close(c.outputCh)
}

// Accept accepts a pointer to an object you want the receive the data into.
// Currently you *have* to pass a pointer to a object. Accept will marshal the received
// object into json and than unmarshal back into the object you have provided. Be wary
// of your json tags and private fields.
//
// NOTE: marshal/unmarshal mechanism will change in future version.
func (c channelConveyor) Accept(v interface{}) error {
	item, ok := <-c.inputCh
	if !ok {
		return io.EOF
	}

	// TODO: Make this pluggable.
	// JSON encode -> decode is just a temporary implementation.
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(item); err != nil {
		return err
	}

	if err := json.NewDecoder(&buf).Decode(v); err != nil {
		return err
	}

	return nil
}

// Emit emits the data on the channel to be accepted by next operation.
func (c channelConveyor) Emit(v interface{}) error {
	c.outputCh <- v

	return nil
}

// ChainNext initiates a new conveyor (B) with the output of this consumer (A) being
// the input of the new consumer. In effect chaining them A->B with the arrow showing direction
// of the items data passed.
func (c channelConveyor) ChainNext() ChannelConveyor {
	return NewChannelConveyor(c.outputCh)
}
