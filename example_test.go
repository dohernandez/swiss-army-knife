package technical_test_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	technical_test "github.com/heetch/Darien-technical-test"
	ttio "github.com/heetch/Darien-technical-test/io"
)

func Example() {
	// An artificial input source.
	const stdin = `{"id":4649,"lat":49.01249051526539,"lng":2.0403327446430257,"created_at":"2016-12-14 07:00:00"}
{"id":10086,"lat":48.907344373066344,"lng":2.3638633128958166,"created_at":"2016-12-14 07:00:00"}
{"id":1874,"lat":48.95913471644928,"lng":2.240928289825033,"created_at":"2016-12-14 07:00:01"}`

	// create input Stdin
	scanner := bufio.NewScanner(strings.NewReader(stdin))
	input := ttio.NewStdinInput(scanner)

	// create output Stdout
	output := ttio.StdoutOutput{}

	p := technical_test.ChannelConveyorProcessor{}

	if err := p.Process(context.TODO(), input, &output); err != nil {
		panic(err)
	}

	// Output:
	// {"id":4649,"lat":49.01249051526539,"lng":2.0403327446430257,"created_at":"2016-12-14 07:00:00"}
	// {"id":10086,"lat":48.907344373066344,"lng":2.3638633128958166,"created_at":"2016-12-14 07:00:00"}
	// {"id":1874,"lat":48.95913471644928,"lng":2.240928289825033,"created_at":"2016-12-14 07:00:01"}
}

func Example_with_filter_operation() {
	// An artificial input source.
	const stdin = `{"id":4649,"lat":49.01249051526539,"lng":2.0403327446430257,"created_at":"2016-12-14 07:00:00"}
{"id":10086,"lat":48.907344373066344,"lng":2.3638633128958166,"created_at":"2016-12-14 07:00:00"}
{"id":1874,"lat":48.95913471644928,"lng":2.240928289825033,"created_at":"2016-12-14 07:00:01"}`

	filterOut := 10086

	// create input from Stdin
	scanner := bufio.NewScanner(strings.NewReader(stdin))
	input := ttio.NewStdinInput(scanner)
	// add unmarshal to decode input value
	input.WithUnmarshaling(func(_ context.Context, i string) (interface{}, error) {
		var a interface{}

		if err := json.Unmarshal([]byte(i), &a); err != nil {
			return nil, err
		}

		return a, nil
	})

	// create output Stdout
	output := ttio.StdoutOutput{}
	// add marshal to encode output value
	output.WithMarshaling(func(_ context.Context, i interface{}) (string, error) {
		r, err := json.Marshal(i)
		if err != nil {
			return "", err
		}

		return string(r), nil
	})

	p := technical_test.ChannelConveyorProcessor{}

	if err := p.Process(context.TODO(), input, &output, func(ctx context.Context, value interface{}) (interface{}, error) {
		m, ok := value.(map[string]interface{})
		if !ok {
			return nil, technical_test.ErrTypeMismatch
		}

		// comparison is done converting values to string using fmt.Sprint to avoid untyped constant
		// which is the type to which the constant is implicitly converted.
		if fmt.Sprint(m["id"]) == fmt.Sprint(filterOut) {
			return nil, technical_test.ErrDoNotEmit
		}

		return value, nil

	}); err != nil {
		panic(err)
	}

	// Output:
	// {"created_at":"2016-12-14 07:00:00","id":4649,"lat":49.01249051526539,"lng":2.0403327446430257}
	// {"created_at":"2016-12-14 07:00:01","id":1874,"lat":48.95913471644928,"lng":2.240928289825033}
}
