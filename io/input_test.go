package io_test

import (
	"bufio"
	"context"
	"io"
	"strings"
	"testing"

	ttio "github.com/heetch/Darien-technical-test/io"
	"github.com/stretchr/testify/assert"
)

const stdinInput = `{"id":4649,"lat":49.01249051526539,"lng":2.0403327446430257,"created_at":"2016-12-14 07:00:00"}
{"id":10086,"lat":48.907344373066344,"lng":2.3638633128958166,"created_at":"2016-12-14 07:00:00"}
{"id":1874,"lat":48.95913471644928,"lng":2.240928289825033,"created_at":"2016-12-14 07:00:01"}`

func TestStdinInputNext(t *testing.T) {
	testCases := []struct {
		scenario string
		assert   func(t *testing.T, ctx context.Context, input ttio.Input)
	}{
		{
			scenario: "Next first time from stdin successful",
			assert: func(t *testing.T, ctx context.Context, input ttio.Input) {
				r, err := input.Next(ctx)
				assert.NoError(t, err)

				assert.Equal(
					t,
					"{\"id\":4649,\"lat\":49.01249051526539,\"lng\":2.0403327446430257,\"created_at\":\"2016-12-14 07:00:00\"}",
					r,
				)
			},
		},
		{
			scenario: "Next two time from stdin successful",
			assert: func(t *testing.T, ctx context.Context, input ttio.Input) {
				_, err := input.Next(ctx)
				assert.NoError(t, err)

				r, err := input.Next(ctx)
				assert.NoError(t, err)

				assert.Equal(
					t,
					"{\"id\":10086,\"lat\":48.907344373066344,\"lng\":2.3638633128958166,\"created_at\":\"2016-12-14 07:00:00\"}",
					r,
				)
			},
		},
		{
			scenario: "Next 4 time from stdin failed, io.Err",
			assert: func(t *testing.T, ctx context.Context, input ttio.Input) {
				_, err := input.Next(ctx)
				assert.NoError(t, err)

				_, err = input.Next(ctx)
				assert.NoError(t, err)

				_, err = input.Next(ctx)
				assert.NoError(t, err)

				_, err = input.Next(ctx)
				assert.EqualError(t, err, io.EOF.Error())
			},
		},
		{
			scenario: "Next from stdin failed, scanner error",
			assert: func(t *testing.T, ctx context.Context, input ttio.Input) {
				_, err := input.Next(ctx)
				assert.NoError(t, err)

				_, err = input.Next(ctx)
				assert.NoError(t, err)

				_, err = input.Next(ctx)
				assert.NoError(t, err)

				_, err = input.Next(ctx)
				assert.EqualError(t, err, io.EOF.Error())
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint.
		t.Run(tc.scenario, func(t *testing.T) {
			ctx := context.TODO()

			scanner := bufio.NewScanner(strings.NewReader(stdinInput))
			input := ttio.NewStdinInput(scanner)

			tc.assert(t, ctx, input)
		})
	}
}
