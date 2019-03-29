package swissarmyknife

import (
	"context"
	"io"
	"sync"

	sakio "github.com/dohernandez/swiss-army-knife/io"
)

// Processor defines a contract to process data.
type Processor interface {
	Process(ctx context.Context, input sakio.Input, output sakio.Output, operations ...Operation) error
}

// ChannelConveyorProcessor a processor that uses ChannelConveyor to share data between operations.
type ChannelConveyorProcessor struct {
	conveyorErrors []error
}

var _ Processor = new(ChannelConveyorProcessor)

// Process processes the data input thro the operations defines and outputted the result.
// Returns error if outputting the result fails.
func (p *ChannelConveyorProcessor) Process(ctx context.Context, input sakio.Input, output sakio.Output, operations ...Operation) error {
	var wg sync.WaitGroup
	inputs := make(chan interface{})
	operationResults := make(chan error)

	// create a ChannelConveyor.
	cc := NewChannelConveyor(inputs)

	// starts the conveyor.
	p.inputConveyor(ctx, &wg, input, cc, operationResults)
	cc = cc.ChainNext()

	// operate the conveyor.
	for _, op := range operations {
		p.operateConveyor(ctx, &wg, op, cc, operationResults)
		cc = cc.ChainNext()
	}

	// ends the conveyor.
	p.outputConveyor(ctx, &wg, output, cc, operationResults)

	// this along with wg.Wait() are why the error handling works and doesn't deadlock.
	finished := make(chan bool, 1)

	// Wait for all operation to return and then close the result chan.
	go func() {
		wg.Wait()

		close(finished)
	}()

	// handler operation errors.
	for {
		var fin bool

		select {
		case <-finished:
			fin = true
		case err := <-operationResults:
			if err != nil {
				p.conveyorErrors = append(p.conveyorErrors, err)
			}
		}

		if fin {
			break
		}
	}

	return output.Write(ctx)
}

// inputConveyor takes the input one by one and start the conveyor sending the data input to the first
// operation in the list.
// As it is a function that runs in the background - using go routines - error will be sent to the main routine
// thro the channel `operationResults`.
func (p *ChannelConveyorProcessor) inputConveyor(ctx context.Context, wg *sync.WaitGroup, input sakio.Input, cc ChannelConveyor, operationResults chan error) {
	wg.Add(1)

	go func(ctx context.Context, c ChannelConveyor) {
		defer func() {
			c.Close()
			wg.Done()
		}()

		r, err := input.Next(ctx)
		for err == nil {
			err = c.Emit(r)
			if err != nil {
				operationResults <- err
			}

			r, err = input.Next(ctx)
		}
		if err != io.EOF {
			operationResults <- err
		}
	}(ctx, cc)
}

// operateConveyor takes the input an apply the operation. The resulting output is sent either to the next
// operation in the list.
// As it is a function that runs in the background - using go routines - error will be sent to the main routine
// thro the channel `operationResults`.
func (p *ChannelConveyorProcessor) operateConveyor(ctx context.Context, wg *sync.WaitGroup, op Operation, cc ChannelConveyor, operationResults chan error) {
	wg.Add(1)

	go func(ctx context.Context, op Operation, c ChannelConveyor) {
		defer func() {
			c.Close()
			wg.Done()
		}()

		for {
			var input interface{}

			if err := c.Accept(&input); err != nil {
				if err == io.EOF {
					break
				}

				operationResults <- err
				continue
			}

			output, err := op(ctx, input)
			if err != nil {
				if err != ErrDoNotEmit {
					operationResults <- err
				}
				continue
			}

			err = c.Emit(output)
			if err != nil {
				operationResults <- err
			}
		}
	}(ctx, op, cc)
}

// outputConveyor takes the result normally after being processed by the operation (In case there is no operation
// it will take the exact input) and add it to the output.
// As it is a function that runs in the background - using go routines - error will be sent to the main routine
// thro the channel `operationResults`.
func (p *ChannelConveyorProcessor) outputConveyor(ctx context.Context, wg *sync.WaitGroup, output sakio.Output, cc ChannelConveyor, operationResults chan error) {
	wg.Add(1)

	go func(ctx context.Context, c ChannelConveyor) {
		defer func() {
			c.Close()
			wg.Done()
		}()

		for {
			var out interface{}

			if err := c.Accept(&out); err != nil {
				if err == io.EOF {
					break
				}

				operationResults <- err
			}

			output.Append(ctx, out)
		}
	}(ctx, cc)
}

// Errors returns errors that happen during the process in case any error occurred.
func (p *ChannelConveyorProcessor) Errors() []error {
	return p.conveyorErrors
}
