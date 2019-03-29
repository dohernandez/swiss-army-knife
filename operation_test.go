package technical_test_test

import (
	"context"
	"encoding/json"
	"testing"

	technical_test "github.com/heetch/Darien-technical-test"
	"github.com/stretchr/testify/assert"
)

func TestFilteringOperation(t *testing.T) {
	// An artificial value source.
	var value interface{}

	// The value is already Unmarshal to make easy the test
	err := json.Unmarshal([]byte(`{"id":1629,"lat":48.83168740132889,"lng":2.2485795413465577,"created_at":"2016-12-14 18:48:11"}`), &value)
	assert.NoError(t, err)

	testCases := []struct {
		scenario string
		pairs    []technical_test.PairKeyValue
		result   map[string]interface{}
	}{
		{
			scenario: "Operation filter out by single key/value pair",
			pairs: []technical_test.PairKeyValue{
				{
					Key:   "id",
					Value: "1629",
				},
			},
		},
		{
			scenario: "Operation do not filter out by single key/value pair",
			pairs: []technical_test.PairKeyValue{
				{
					Key:   "id",
					Value: "1874",
				},
			},
			result: value.(map[string]interface{}),
		},
		{
			scenario: "Operation filter out by multiple key/value pair",
			pairs: []technical_test.PairKeyValue{
				{
					Key:   "lat",
					Value: "48.83168740132889",
				},
				{
					Key:   "lng",
					Value: "2.2485795413465577",
				},
			},
		},
		{
			scenario: "Operation do not filter out by multiple key/value pair",
			pairs: []technical_test.PairKeyValue{
				{
					Key:   "lat",
					Value: "50.62508560995884",
				},
				{
					Key:   "lng",
					Value: "2.2485795413465577",
				},
			},
			result: value.(map[string]interface{}),
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint.
		t.Run(tc.scenario, func(t *testing.T) {
			ctx := context.TODO()
			operation := technical_test.NewFilteringOperation(ctx, tc.pairs)

			r, err := operation(ctx, value)
			if tc.result == nil {
				assert.EqualError(t, err, technical_test.ErrDoNotEmit.Error())
				assert.Empty(t, r)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.result, r)
			}
		})
	}
}

func TestAppendInformationOperation(t *testing.T) {
	// An artificial value source.
	var value interface{}

	// The value is already Unmarshal to make easy the test
	err := json.Unmarshal([]byte(`{"id":1629,"lat":48.83168740132889,"lng":2.2485795413465577,"created_at":"2016-12-14 18:48:11"}`), &value)
	assert.NoError(t, err)

	testCases := []struct {
		scenario string
		pairs    []technical_test.PairKeyValue
		result   map[string]interface{}
	}{
		{
			scenario: "Operation append country information by single key/value pair",
			pairs: []technical_test.PairKeyValue{
				{
					Key:   "country",
					Value: "fr",
				},
			},
			result: func() map[string]interface{} {
				// value type does not vary, ok to rely on panic
				// nolint:errcheck
				v := value.(map[string]interface{})
				r := copyMap(v)

				r["country"] = "fr"

				return r
			}(),
		},
		{
			scenario: "Operation append/modify country, city and lat by multiple key/value pair",
			pairs: []technical_test.PairKeyValue{
				{
					Key:   "lat",
					Value: "50.62508560995884",
				},
				{
					Key:   "country",
					Value: "fr",
				},
				{
					Key:   "name",
					Value: "Darien",
				},
			},
			result: func() map[string]interface{} {
				// value type does not vary, ok to rely on panic
				// nolint:errcheck
				v := value.(map[string]interface{})
				r := copyMap(v)

				r["lat"] = "50.62508560995884"
				r["country"] = "fr"
				r["name"] = "Darien"

				return r
			}(),
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint.
		value := value
		t.Run(tc.scenario, func(t *testing.T) {
			ctx := context.TODO()
			operation := technical_test.NewAppendInformationOperation(ctx, tc.pairs)

			// copyMap the value to keep the original.
			r, err := operation(ctx, copyMap(value.(map[string]interface{})))
			if tc.result == nil {
				assert.Error(t, err)
				assert.Empty(t, r)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.result, r)
			}
		})
	}
}

func copyMap(originalMap map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{}, len(originalMap))
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}

func TestRemoveInformationOperation(t *testing.T) {
	// An artificial value source.
	var value interface{}

	// The value is already Unmarshal to make easy the test
	err := json.Unmarshal([]byte(`{"id":1629,"lat":48.83168740132889,"lng":2.2485795413465577,"created_at":"2016-12-14 18:48:11"}`), &value)
	assert.NoError(t, err)

	testCases := []struct {
		scenario string
		keys     []technical_test.Key
		result   map[string]interface{}
	}{
		{
			scenario: "Operation remove created_at information by single key",
			keys:     []technical_test.Key{"created_at"},
			result: func() map[string]interface{} {
				// value type does not vary, ok to rely on panic
				// nolint:errcheck
				v := value.(map[string]interface{})
				r := copyMap(v)

				delete(r, "created_at")

				return r
			}(),
		},
		{
			scenario: "Operation remove lat and lng information by multiple key",
			keys:     []technical_test.Key{"lat", "lng"},
			result: func() map[string]interface{} {
				// value type does not vary, ok to rely on panic
				// nolint:errcheck
				v := value.(map[string]interface{})
				r := copyMap(v)

				delete(r, "lat")
				delete(r, "lng")

				return r
			}(),
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint.
		value := value
		t.Run(tc.scenario, func(t *testing.T) {
			ctx := context.TODO()
			operation := technical_test.NewRemoveInformationOperation(ctx, tc.keys)

			// copyMap the value to keep the original.
			r, err := operation(ctx, copyMap(value.(map[string]interface{})))
			if tc.result == nil {
				assert.Error(t, err)
				assert.Empty(t, r)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.result, r)
			}
		})
	}
}

func TestPrefixKeyOperation(t *testing.T) {
	// An artificial value source.
	var value interface{}

	// The value is already Unmarshal to make easy the test
	err := json.Unmarshal([]byte(`{"id":1629,"lat":48.83168740132889,"lng":2.2485795413465577,"created_at":"2016-12-14 18:48:11"}`), &value)
	assert.NoError(t, err)

	testCases := []struct {
		scenario string
		pairs    []technical_test.PairKeyPrefix
		result   map[string]interface{}
	}{
		{
			scenario: "Operation prefix id (_id) by single key/prefix pair",
			pairs: []technical_test.PairKeyPrefix{
				{
					Key:    "id",
					Prefix: "_",
				},
			},
			result: func() map[string]interface{} {
				// value type does not vary, ok to rely on panic
				// nolint:errcheck
				v := value.(map[string]interface{})
				r := copyMap(v)

				r["_id"] = v["id"]

				delete(r, "id")

				return r
			}(),
		},
		{
			scenario: "Operation prefix id (_id), lat (c_lat) and lng (c_lng) by multiple key/prefix pair",
			pairs: []technical_test.PairKeyPrefix{
				{
					Key:    "lat",
					Prefix: "c_",
				},
				{
					Key:    "lng",
					Prefix: "c_",
				},
			},
			result: func() map[string]interface{} {
				// value type does not vary, ok to rely on panic
				// nolint:errcheck
				v := value.(map[string]interface{})
				r := copyMap(v)

				r["c_lat"] = v["lat"]
				r["c_lng"] = v["lng"]

				delete(r, "lat")
				delete(r, "lng")

				return r
			}(),
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint.
		value := value
		t.Run(tc.scenario, func(t *testing.T) {
			ctx := context.TODO()
			operation := technical_test.NewPrefixKeyOperation(ctx, tc.pairs)

			// copyMap the value to keep the original.
			r, err := operation(ctx, copyMap(value.(map[string]interface{})))
			if tc.result == nil {
				assert.Error(t, err)
				assert.Empty(t, r)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.result, r)
			}
		})
	}
}
