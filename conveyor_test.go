package swissarmyknife_test

import (
	"io"
	"testing"

	swiss_army_knife "github.com/dohernandez/swiss-army-knife"
	"github.com/stretchr/testify/assert"
)

type Item struct {
	F1 string
	F2 int
	F3 int
}

var items = []Item{
	{"asd", 1, 3},
	{"asd", 21313, 3},
	{"asd", 1, 3},
	{"asd", 1, 3},
	{"asd", 1, 3},
	{"asd", 1, 3},
	{"asd", 1, 3},
	{"asd", 1, 3},
	{"asd", 1, 3},
	{"asd", 1, 3},
	{"asd", 1, 3},
}

func newChanWithItems(ii []Item) chan interface{} {
	ch := make(chan interface{})

	go func() {
		for _, i := range ii {
			ch <- i
		}
		close(ch)
	}()

	return ch
}

func expectItemsOnAccept(t *testing.T, c swiss_army_knife.Conveyor, ii []Item) {
	n := 0
	for {
		val := Item{}
		if err := c.Accept(&val); err != nil {
			if err != io.EOF {
				t.Fatal("Error on accept", err)
			}

			return
		}

		assert.Equal(t, val.F1, ii[n].F1, "Accept: item field values do not match")
		assert.Equal(t, val.F2, ii[n].F2, "Accept: item field values do not match")
		assert.Equal(t, val.F3, ii[n].F3, "Accept: item field values do not match")

		n++

		if n > len(ii) {
			t.Fatal("To many items on accept")
			return
		}

	}
}

func TestChannelConveyor(t *testing.T) {
	testCases := []struct {
		scenario string
		assert   func(t *testing.T, cc swiss_army_knife.ChannelConveyor, ii []Item)
	}{
		{
			scenario: "Coveyor check input is correct",
			assert: func(t *testing.T, cc swiss_army_knife.ChannelConveyor, ii []Item) {
				expectItemsOnAccept(t, cc, items)

				cc.Close()
			},
		},
		{
			scenario: "Coveyor check input in chain next is correct",
			assert: func(t *testing.T, cc swiss_army_knife.ChannelConveyor, ii []Item) {
				for {
					val := Item{}
					if err := cc.Accept(&val); err != nil {
						if err != io.EOF {
							t.Fatal("Error on accept", err)
						}

						break
					}

					if err := cc.Emit(&val); err != nil {
						t.Fatal("Error on emit", err)
					}
				}

				cc.Close()
				ccNext := cc.ChainNext()

				expectItemsOnAccept(t, ccNext, items)

				ccNext.Close()
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint.
		t.Run(tc.scenario, func(t *testing.T) {
			input := newChanWithItems(items)

			cc := swiss_army_knife.NewChannelConveyor(input)

			tc.assert(t, cc, items)
		})
	}
}
