package query

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gotoxu/assert"
)

type E1 struct {
	F01 int     `url:"f01"`
	F02 int     `url:"-"`
	F03 string  `url:"f03"`
	F04 string  `url:"f04,omitempty"`
	F05 bool    `url:"f05"`
	F06 bool    `url:"f06"`
	F07 *string `url:"f07"`
	F08 *int8   `url:"f08"`
	F09 float64 `url:"f09"`
	F10 func()  `url:"-"`
	F11 inner
}

type inner struct {
	F12 int
}

func TestFilled(t *testing.T) {
	f07 := "seven"
	var f08 int8 = 8
	s := &E1{
		F01: 1,
		F02: 2,
		F03: "three",
		F04: "four",
		F05: true,
		F06: false,
		F07: &f07,
		F08: &f08,
		F09: 1.618,
		F10: func() {},
		F11: inner{12},
	}

	vals, errs := NewEncoder().Encode(s)
	assert.Nil(t, errs)

	valExists(t, "f01", "1", vals)
	valNotExists(t, "f02", vals)
	valExists(t, "f03", "three", vals)
	valExists(t, "f05", "true", vals)
	valExists(t, "f06", "false", vals)
	valExists(t, "f07", "seven", vals)
	valExists(t, "f08", "8", vals)
	valExists(t, "f09", "1.618000", vals)
	valExists(t, "F12", "12", vals)
}

func TestEmpty(t *testing.T) {
	s := &E1{
		F01: 1,
		F02: 2,
		F03: "three",
	}

	vals, err := NewEncoder().Encode(s)
	assert.Nil(t, err)
	valExists(t, "f03", "three", vals)
	valNotExists(t, "f04", vals)
}

func TestSlices(t *testing.T) {
	type oneAsWord int
	ones := []oneAsWord{1, 2}
	s1 := &struct {
		ones  []oneAsWord `url:"ones"`
		ints  []int       `url:"ints"`
		empty []int       `url:"empty,omitempty"`
	}{ones, []int{1, 1}, []int{}}

	encoder := NewEncoder()
	encoder.RegisterEncoder(ones[0], func(v reflect.Value) string { return "one" })

	vals, err := encoder.Encode(s1)
	assert.Nil(t, err)
	valsExists(t, "ones", []string{"one", "one"}, vals)
	valsExists(t, "ints", []string{"1", "1"}, vals)
	valNotExists(t, "empty", vals)
}

func valExists(t *testing.T, key string, expect string, result map[string][]string) {
	valsExists(t, key, []string{expect}, result)
}

func valsExists(t *testing.T, key string, expect []string, result map[string][]string) {
	vals, ok := result[key]
	assert.True(t, ok, fmt.Sprintf("Key not found. Expected: %s", key))
	assert.DeepEqual(t, vals, expect)
}

func valNotExists(t *testing.T, key string, result map[string][]string) {
	vals, ok := result[key]
	assert.False(t, ok, fmt.Sprintf("Key not ommited. Expected: empty; got: %v", vals))
}
