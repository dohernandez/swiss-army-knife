package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	swiss_army_knife "github.com/dohernandez/swiss-army-knife"
	sakio "github.com/dohernandez/swiss-army-knife/io"
	"github.com/dohernandez/swiss-army-knife/version"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	filterKey    = "filter"
	appendKey    = "append"
	removeKey    = "remove"
	prefixingKey = "prefix"
)

var binaryName = "swiss-army-knife"

var errInvalidPairKeyValue = errors.New("invalid pair key/value. Valid format key:value")

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version.Info().String())
		return
	}

	ctx, cancelCtx := context.WithCancel(context.TODO())
	defer cancelCtx()

	app := cli.NewApp()
	app.Version = version.Info().Version
	app.Name = binaryName

	app.Usage = "To give some background, the stream of JSON objects can be locations updates from drivers, comments about rides etc. "
	app.UsageText = fmt.Sprintf("%s [arguments]", binaryName)
	app.HideVersion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  filterKey + ", f",
			Usage: "Filter out base on key/value pair. Valid format key:value;keyn:valuen. Example id:347.",
		},
		cli.StringFlag{
			Name:  appendKey + ", a",
			Usage: "Append key/value pair. Valid format key:value;keyn:valuen. Example id:347.",
		},
		cli.StringFlag{
			Name:  removeKey + ", r",
			Usage: "Remove a key. Valid format key:value;keyn:valuen. Example id:347.",
		},
		cli.StringFlag{
			Name:  prefixingKey + ", p",
			Usage: "Prefixing a key. Valid format key:value;keyn:valuen. Example id:347.",
		},
	}

	app.Action = func(cliCtx *cli.Context) error {
		input := initInput()

		output := initOutput()

		p := swiss_army_knife.ChannelConveyorProcessor{}

		// init operations
		var operations []swiss_army_knife.Operation

		// Filter out base on key/value pair.
		if cliCtx.String(filterKey) != "" {
			value := cliCtx.String(filterKey)

			kvs, err := splitPairs(value)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("%s (%s)", filterKey, value))
			}

			var pairs []swiss_army_knife.PairKeyValue
			for _, pair := range kvs {
				pairs = append(pairs, swiss_army_knife.PairKeyValue{
					Key:   swiss_army_knife.Key(pair[0]),
					Value: swiss_army_knife.Value(pair[1]),
				})
			}
			operations = append(operations, swiss_army_knife.NewFilteringOperation(ctx, pairs))
		}

		if cliCtx.String(appendKey) != "" {
			value := cliCtx.String(appendKey)

			kvs, err := splitPairs(value)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("%s (%s)", appendKey, value))
			}

			var pairs []swiss_army_knife.PairKeyValue
			for _, pair := range kvs {
				pairs = append(pairs, swiss_army_knife.PairKeyValue{
					Key:   swiss_army_knife.Key(pair[0]),
					Value: swiss_army_knife.Value(pair[1]),
				})
			}
			operations = append(operations, swiss_army_knife.NewAppendInformationOperation(ctx, pairs))
		}

		if cliCtx.String(removeKey) != "" {
			var keys []swiss_army_knife.Key
			for _, key := range strings.Split(cliCtx.String(removeKey), ":") {
				keys = append(keys, swiss_army_knife.Key(key))
			}
			operations = append(operations, swiss_army_knife.NewRemoveInformationOperation(ctx, keys))
		}

		// Prefixing a key.
		if cliCtx.String(prefixingKey) != "" {
			value := cliCtx.String(prefixingKey)

			kvs, err := splitPairs(value)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("%s (%s)", prefixingKey, value))
			}

			var pairs []swiss_army_knife.PairKeyPrefix
			for _, pair := range kvs {
				pairs = append(pairs, swiss_army_knife.PairKeyPrefix{
					Key:    swiss_army_knife.Key(pair[0]),
					Prefix: pair[1],
				})
			}
			operations = append(operations, swiss_army_knife.NewPrefixKeyOperation(ctx, pairs))
		}

		// Process data
		if err := p.Process(ctx, input, output, operations...); err != nil {
			return err
		}

		// Checking if there were any error while processing data
		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				// TODO For the purpose of this iteration is enough printing the errors.
				fmt.Println(err)
			}
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func initInput() *sakio.StdinInput {
	// create input Stdin
	scanner := bufio.NewScanner(os.Stdin)
	input := sakio.NewStdinInput(scanner)
	// add unmarshal to decode input value
	input.WithUnmarshaling(func(_ context.Context, i string) (interface{}, error) {
		var a interface{}

		if err := json.Unmarshal([]byte(i), &a); err != nil {
			return nil, err
		}

		return a, nil
	})

	return input
}

func initOutput() *sakio.StdoutOutput {
	// create output Stdout
	output := sakio.StdoutOutput{}
	// add marshal to encode output value
	output.WithMarshaling(func(_ context.Context, i interface{}) (string, error) {
		r, err := json.Marshal(i)
		if err != nil {
			return "", err
		}

		return string(r), nil
	})

	return &output
}

func splitPairs(value string) (pairs [][]string, err error) {
	kvs := strings.Split(value, ";")
	for _, kv := range kvs {
		pair := strings.Split(kv, ":")

		if len(pair) != 2 {
			return nil, errInvalidPairKeyValue
		}

		pairs = append(pairs, []string{pair[0], pair[1]})
	}

	return pairs, nil
}
