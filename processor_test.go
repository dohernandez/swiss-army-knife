package technical_test_test

import (
	"bufio"
	"context"
	"errors"
	"strings"
	"testing"

	technical_test "github.com/heetch/Darien-technical-test"
	ttio "github.com/heetch/Darien-technical-test/io"
	"github.com/heetch/Darien-technical-test/test"
	"github.com/stretchr/testify/assert"
)

const stdinInput = `{"id":7064,"lat":48.88340457471041,"lng":2.3952910238105294,"created_at":"2016-12-14 18:48:10"}
{"id":11426,"lat":48.927968740518686,"lng":2.2497446977911437,"created_at":"2016-12-14 18:48:10"}
{"id":1629,"lat":48.83168740132889,"lng":2.2485795413465577,"created_at":"2016-12-14 18:48:11"}`

func TestChannelConveyorProcessor(t *testing.T) {
	testCases := []struct {
		scenario   string
		operations []technical_test.Operation
		output     string
		errors     []error
	}{
		{
			scenario: "Process data successful, without operation",
			output:   stdinInput,
		},
		{
			scenario: "Process data successful, with operation",
			operations: []technical_test.Operation{
				func(_ context.Context, value interface{}) (interface{}, error) {
					return value, nil
				},
			},
			output: stdinInput,
		},
		{
			scenario: "Process data unsuccessful, with operation fails",
			operations: []technical_test.Operation{
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
			input := ttio.NewStdinInput(scanner)

			output := ttio.StdoutOutput{}

			p := technical_test.ChannelConveyorProcessor{}

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
