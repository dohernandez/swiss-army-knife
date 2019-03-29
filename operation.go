package technical_test

import (
	"context"
	"fmt"
)

// Operation apply logic (decorate/filter/modify) to the input data.
// Returns any error that occurred, otherwise value processed.
type Operation func(ctx context.Context, value interface{}) (interface{}, error)

type (
	// Key represent a key name
	Key string
	// Value type is string intentionally. Anyway it has to be string during comparison.
	Value string
)

func (k Key) String() string {
	return string(k)
}

func (v Value) String() string {
	return string(v)
}

// PairKeyValue represents key/value pair.
type PairKeyValue struct {
	Key   Key
	Value Value
}

// NewFilteringOperation creates a filtering Operation based on pairs.
// The PairKeyValue is used as a criteria to filter out the value. In case of using multiple PairKeyValue, it behave as an AND.
// All the PairKeyValue must be in the value otherwise does not match.
//
// Accepts only value as a map[string]interface{} type.
//
// ErrTypeMismatch is returned if casting value interface{} to a map[string]interface{} fails.
// ErrDoNotEmit is returned when all PairKeyValue matched, allowing the value to be skipped.
// value is returned when one PairKeyValue do not match, value will be emit.
//
// Common initialization example:
//
//      operation := NewFilteringOperation(
// 			context.TODO(),
// 			[]PairKeyValue{
//				{
//					Key:   "id",
//					Value: "1629",
//				},
//			},
// 		)
//
func NewFilteringOperation(_ context.Context, pairs []PairKeyValue) Operation {
	return func(ctx context.Context, value interface{}) (interface{}, error) {
		m, ok := value.(map[string]interface{})
		if !ok {
			return nil, ErrTypeMismatch
		}

		match := true
		for _, pair := range pairs {
			// comparison is done using string (used fmt.Sprint) to avoid untyped constant
			// which is the type to which the value is implicitly converted
			if fmt.Sprint(m[pair.Key.String()]) != pair.Value.String() {
				match = false

				break
			}
		}

		if match {
			// PairKeyValue matched, value must be skipped
			return nil, ErrDoNotEmit
		}

		return value, nil
	}
}

// NewAppendInformationOperation creates an append information Operation based on pairs.
// The PairKeyValue is used to add an extra information to the value or replacing information, depending
// if the key exists or not.
//
// Accepts only value as a map[string]interface{} type.
//
// ErrTypeMismatch is returned if casting value interface{} to a map[string]interface{} fails.
// value is returned with all PairKeyValue appended.
//
// Common initialization example:
//
//      operation := NewAppendInformationOperation(
// 			context.TODO(),
// 			[]PairKeyValue{
//				{
//					Key:   "id",
//					Value: "1629",
//				},
//			},
// 		)
//
func NewAppendInformationOperation(_ context.Context, pairs []PairKeyValue) Operation {
	return func(ctx context.Context, value interface{}) (interface{}, error) {
		m, ok := value.(map[string]interface{})
		if !ok {
			return nil, ErrTypeMismatch
		}

		for _, pair := range pairs {
			m[pair.Key.String()] = pair.Value.String()
		}

		return m, nil
	}
}

// NewRemoveInformationOperation creates a remove information Operation based on key.
// The Key is used to remove information from the value.
//
// Accepts only value as a map[string]interface{} type.
//
// ErrTypeMismatch is returned if casting value interface{} to a map[string]interface{} fails.
// value is returned with all Key removed.
//
// Common initialization example:
//
//      operation := NewRemoveInformationOperation(
// 			context.TODO(),
// 			[]Key{"id"},
// 		)
//
func NewRemoveInformationOperation(_ context.Context, keys []Key) Operation {
	return func(ctx context.Context, value interface{}) (interface{}, error) {
		m, ok := value.(map[string]interface{})
		if !ok {
			return nil, ErrTypeMismatch
		}

		for _, key := range keys {
			delete(m, key.String())
		}

		return m, nil
	}
}

// PairKeyPrefix represents key/prefix pair.
type PairKeyPrefix struct {
	Key    Key
	Prefix string
}

// NewPrefixKeyOperation creates a prefix key Operation based on key/prefix pair.
// The PairKeyPrefix is used to prefix an extra information from the value.
//
// Accepts only value as a map[string]interface{} type.
//
// ErrTypeMismatch is returned if casting value interface{} to a map[string]interface{} fails.
// value is returned with all Key prefixed.
//
// Common initialization example:
//
//      operation := NewPrefixKeyOperation(
// 			context.TODO(),
// 			[]PairKeyPrefix{
//				{
//					Key:    "id",
//					Prefix: "_",
//				},
//			},
// 		)
//
func NewPrefixKeyOperation(_ context.Context, pairs []PairKeyPrefix) Operation {
	return func(ctx context.Context, value interface{}) (interface{}, error) {
		m, ok := value.(map[string]interface{})
		if !ok {
			return nil, ErrTypeMismatch
		}

		for _, pair := range pairs {
			if v, ok := m[pair.Key.String()]; ok {
				delete(m, pair.Key.String())

				k := pair.Prefix + pair.Key.String()
				m[k] = v
			}
		}

		return m, nil
	}
}
