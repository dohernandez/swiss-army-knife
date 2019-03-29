package io_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ttio "github.com/heetch/Darien-technical-test/io"
)

func ExampleStdinInput_Next() {
	// An artificial input source.
	const stdin = `{"id":4649,"lat":49.01249051526539,"lng":2.0403327446430257,"created_at":"2016-12-14 07:00:00"}
{"id":10086,"lat":48.907344373066344,"lng":2.3638633128958166,"created_at":"2016-12-14 07:00:00"}
{"id":1874,"lat":48.95913471644928,"lng":2.240928289825033,"created_at":"2016-12-14 07:00:01"}`

	scanner := bufio.NewScanner(strings.NewReader(stdin))
	input := ttio.NewStdinInput(scanner)

	if r, err := input.Next(context.TODO()); err == nil {
		fmt.Printf("%v", r)
	}

	// Output:
	// {"id":4649,"lat":49.01249051526539,"lng":2.0403327446430257,"created_at":"2016-12-14 07:00:00"}
}

func ExampleStdinInput_WithUnmarshaling() {
	// An artificial input source.
	const stdin = `{"id":4649,"lat":49.01249051526539,"lng":2.0403327446430257,"created_at":"2016-12-14 07:00:00"}
{"id":10086,"lat":48.907344373066344,"lng":2.3638633128958166,"created_at":"2016-12-14 07:00:00"}
{"id":1874,"lat":48.95913471644928,"lng":2.240928289825033,"created_at":"2016-12-14 07:00:01"}`

	type stdinItem struct {
		ID        int64   `json:"id"`
		Lat       float64 `json:"lat"`
		Lng       float64 `json:"lng"`
		CreatedAt string  `json:"created_at"`
	}

	var stdItem stdinItem

	scanner := bufio.NewScanner(strings.NewReader(stdin))
	input := ttio.NewStdinInput(scanner)
	input.WithUnmarshaling(func(_ context.Context, i string) (interface{}, error) {
		if err := json.Unmarshal([]byte(i), &stdItem); err != nil {
			return nil, err
		}

		return stdItem, nil
	})

	if r, err := input.Next(context.TODO()); err == nil {
		fmt.Printf("%+v", r)
	}

	// Output:
	// {ID:4649 Lat:49.01249051526539 Lng:2.0403327446430257 CreatedAt:2016-12-14 07:00:00}
}

func ExampleStdoutOutput_Write() {
	ctx := context.TODO()
	output := ttio.StdoutOutput{}

	output.Append(ctx, `{"id":4649,"lat":49.01249051526539,"lng":2.0403327446430257,"created_at":"2016-12-14 07:00:00"}`)

	if err := output.Write(ctx); err == nil {
		panic(err)
	}

	// Output:
	// {"id":4649,"lat":49.01249051526539,"lng":2.0403327446430257,"created_at":"2016-12-14 07:00:00"}
}

func ExampleStdoutOutput_WithMarshaling() {
	type stdinItem struct {
		ID        int64   `json:"id"`
		Lat       float64 `json:"lat"`
		Lng       float64 `json:"lng"`
		CreatedAt string  `json:"created_at"`
	}

	stdItem := stdinItem{
		ID:        4649,
		Lat:       49.01249051526539,
		Lng:       2.0403327446430257,
		CreatedAt: "2016-12-14 07:00:00",
	}
	ctx := context.TODO()
	output := ttio.StdoutOutput{}
	output.WithMarshaling(func(_ context.Context, i interface{}) (string, error) {
		r, err := json.Marshal(i)
		if err != nil {
			return "", err
		}

		return string(r), nil
	})

	output.Append(ctx, stdItem)

	if err := output.Write(ctx); err == nil {
		panic(err)
	}

	// Output:
	// {"id":4649,"lat":49.01249051526539,"lng":2.0403327446430257,"created_at":"2016-12-14 07:00:00"}
}
