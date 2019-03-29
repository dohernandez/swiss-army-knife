package swissarmyknife_test

import (
	"bufio"
	"context"
	"errors"
	"strings"
	"testing"

	swiss_army_knife "github.com/dohernandez/swiss-army-knife"
	sakio "github.com/dohernandez/swiss-army-knife/io"
	"github.com/dohernandez/swiss-army-knife/test"
	"github.com/stretchr/testify/assert"
)

const stdinInput = `{"id":7064,"lat":48.88340457471041,"lng":2.3952910238105294,"created_at":"2016-12-14 18:48:10"}
{"id":11426,"lat":48.927968740518686,"lng":2.2497446977911437,"created_at":"2016-12-14 18:48:10"}
{"id":1629,"lat":48.83168740132889,"lng":2.2485795413465577,"created_at":"2016-12-14 18:48:11"}`

func TestChannelConveyorProcessor(t *testing.T) {
	testCases := []struct {
		scenario   string
		operations []swiss_army_knife.Operation
		output     string
		errors     []error
	}{
		{
			scenario: "Process data successful, without operation",
			output:   stdinInput,
		},
		{
			scenario: "Process data successful, with operation",
			operations: []swiss_army_knife.Operation{
				func(_ context.Context, value interface{}) (interface{}, error) {
					return value, nil
				},
			},
			output: stdinInput,
		},
		{
			scenario: "Process data unsuccessful, with operation fails",
			operations: []swiss_army_knife.Operation{
				func(_ context.Context, _ interface{}) (interface{}, error) {
					return nil, errors.New("operation fails")
				},
			},
			output: "",
			errors: []error{
				errors.New("operation fails"),
				errors.New("operation fails"),
				errors.New("operation fails"),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint
		t.Run(tc.scenario, func(t *testing.T) {
			ctx := context.TODO()

			scanner := bufio.NewScanner(strings.NewReader(stdinInput))
			input := sakio.NewStdinInput(scanner)

			output := sakio.StdoutOutput{}

			p := swiss_army_knife.ChannelConveyorProcessor{}

			stdout := test.CaptureStdOut(func() {
				err := p.Process(ctx, input, &output, tc.operations...)
				assert.NoError(t, err)
			})

			assert.Equal(
				t,
				tc.output,
				stdout,
			)
			assert.EqualValues(t, tc.errors, p.Errors())
		})
	}
}
