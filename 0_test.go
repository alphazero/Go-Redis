package redis

import (
	"testing"
	"testing/quick"
)

// map of various named datasets
var testdata = make(map[string]interface{})

const (
	_testdata_keys     string = "keys"     // []string of keys used in k/v
	_testdata_kv              = "k/v"      // a map of strings -> []byte
	_testdata_counters        = "counters" // a map of strings -> int64
	_testdata_nums            = "counters" // []int64

)

// Its a hack to use the testing/quick functions to load up some
// test data.  And that is it.
//
// NOTE:
// TestFu sets up the test data and it is expected to be in
// 0<whatever>_test.go.  So don't move it from 0_test.go!
func TestFu(t *testing.T) {

	config := &quick.Config{}
	config.MaxCount = 100

	// collect random k/vs
	testdata[_testdata_keys] = make([]string, 0)
	testdata[_testdata_kv] = make(map[string][]byte)
	setup_kv := func(k string, v []byte) bool {
		dataset := testdata[_testdata_kv].(map[string][]byte)
		dataset[k] = v
		keyset := testdata[_testdata_keys].([]string)
		testdata[_testdata_keys] = append(keyset, k)
		return true // always
	}
	quick.Check(setup_kv, config)

	// collect random counter
	testdata[_testdata_counters] = make(map[string]int64)
	setup_counters := func(k string, v int64) bool {
		dataset := testdata[_testdata_counters].(map[string]int64)
		dataset[k] = v
		return true // always
	}
	quick.Check(setup_counters, config)

	// collect random numbers
	testdata[_testdata_nums] = make([]int64, 0)
	setup_nums := func(v int64) bool {
		dataset := testdata[_testdata_nums].([]int64)
		dataset = append(dataset, v)
		return true // always
	}
	quick.Check(setup_nums, config)
}

func _test_getDefaultClient() (Client, error) {
	spec := DefaultSpec()
	spec.Db(13).Password("go-redis")
	return NewSynchClientWithSpec(spec)
}
