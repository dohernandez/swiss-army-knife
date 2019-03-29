package io_test

import (
	"context"
	"strings"
	"testing"

	ttio "github.com/heetch/Darien-technical-test/io"
	"github.com/heetch/Darien-technical-test/test"
	"github.com/stretchr/testify/assert"
)

const stdoutOutput = `{"id":7064,"lat":48.88340457471041,"lng":2.3952910238105294,"created_at":"2016-12-14 18:48:10"}
{"id":11426,"lat":48.927968740518686,"lng":2.2497446977911437,"created_at":"2016-12-14 18:48:10"}
{"id":1629,"lat":48.83168740132889,"lng":2.2485795413465577,"created_at":"2016-12-14 18:48:11"}`

func TestStdoutOutputWrite(t *testing.T) {
	assert := func(t *testing.T, ctx context.Context, output ttio.Output) {
		for _, v := range strings.Split(stdoutOutput, "\n") {
			output.Append(ctx, v)
		}

		stdout := test.CaptureStdOut(func() {
			err := output.Write(ctx)
			assert.NoError(t, err)
		})

		assert.Equal(
			t,
			stdoutOutput,
			stdout,
		)
	}

	testCases := []struct {
		scenario string
		marshall ttio.MarshalOutput
		assert   func(t *testing.T, ctx context.Context, output ttio.Output)
	}{
		{
			scenario: "Write to Stdout successful",
			assert:   assert,
		},
		{
			scenario: "Write to Stdout successful with marshal func",
			marshall: func(_ context.Context, i interface{}) (string, error) {
				return i.(string), nil
			},
			assert: assert,
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint.
		t.Run(tc.scenario, func(t *testing.T) {
			ctx := context.TODO()

			output := ttio.StdoutOutput{}

			if tc.marshall != nil {
				output.WithMarshaling(tc.marshall)
			}

			tc.assert(t, ctx, &output)
		})
	}
}
