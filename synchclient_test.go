package redis

import (
	"log"
	"testing"

//	"testing/quick"
)

func _test_getDefConnSpec() *ConnectionSpec {

	host := "localhost"
	port := 6379
	db := 13
	password := "go-redis"

	connspec := DefaultSpec().Host(host).Port(port).Db(db).Password(password)
	return connspec
}

func flushAndQuitOnCompletion(t *testing.T, client Client) {
	// flush it
	e := client.Flushdb()
	if e != nil {
		t.Errorf("on Flushdb", e)
	}
	e = client.Quit()
	if e != nil {
		t.Errorf("on Quit", e)
	}
}

// Check that connection is actually passing passwords from spec
// and catching AUTH ERRs.
func TestClientConnectWithBadSpec(t *testing.T) {
	spec := _test_getDefConnSpec()
	spec.Password("bad-password")
	c, expected := NewSynchClientWithSpec(spec)
	if expected == nil {
		t.Error("BUG: Expected a RedisError")
	}
	if c != nil {
		t.Error("BUG: client reference on error MUST be nil")
	}
}

// Check that connection is actually passing passwords from spec
func TestClientConnectWithSpec(t *testing.T) {
	spec := _test_getDefConnSpec()

	c, err := NewSynchClientWithSpec(spec)
	if err != nil {
		t.Error("failed to create client with spec. Error: ", err.Message())
	} else if c == nil {
		t.Error("BUG: client is nil")
	}
}

func TestPing(t *testing.T) {
	client, e := _test_getDefaultClient()
	if e != nil {
		t.Fatalf("on getDefaultClient", e)
	}

	if e := client.Ping(); e != nil {
		t.Errorf("on Ping() - %s", e)
	}

	flushAndQuitOnCompletion(t, client)
}

func TestSetGet(t *testing.T) {
	client, e := _test_getDefaultClient()
	if e != nil {
		t.Fatalf("on getDefaultClient", e)
	}

	for k, v := range testdata[_testdata_kv].(map[string][]byte) {
		if e := client.Set(k, v); e != nil {
			t.Errorf("on Set(%s, %s) - %s", k, v, e)
		}
		got, e := client.Get(k)
		if e != nil {
			t.Errorf("on Get(%s) - %s", k, e)
		}
		if got == nil {
			t.Errorf("on Get(%s) - got nil", k)
		}
		if !compareByteArrays(got, v) {
			t.Errorf("on Get(%s) - got:%s expected:%s", k, got, v)
		}
	}

	flushAndQuitOnCompletion(t, client)
}

func TestGetType(t *testing.T) {
	client, e := _test_getDefaultClient()
	if e != nil {
		t.Fatalf("on getDefaultClient", e)
	}

	var expected, got KeyType
	var key string

	key = "string-key"
	if e := client.Set(key, []byte("woof")); e != nil {
		t.Errorf("on Set - %s", e)
	}
	expected = RT_STRING
	got, e = client.Type(key)
	if e != nil {
		t.Errorf("on Type(%s) - %s", key, e)
	}
	if got != expected {
		t.Errorf("on Type() expected:%d got:%d - %s", expected, got)
	}

	key = "set-key"
	if _, e := client.Sadd(key, []byte("woof")); e != nil {
		t.Errorf("on Sadd - %s", e)
	}
	expected = RT_SET
	got, e = client.Type(key)
	if e != nil {
		t.Errorf("on Type(%s) - %s", key, e)
	}
	if got != expected {
		t.Errorf("on Type() expected:%d got:%d - %s", expected, got)
	}

	key = "list-key"
	if e = client.Lpush(key, []byte("woof")); e != nil {
		t.Errorf("on Lpush - %s", e)
	}
	expected = RT_LIST
	got, e = client.Type(key)
	if e != nil {
		t.Errorf("on Type(%s) - %s", key, e)
	}
	if got != expected {
		t.Errorf("on Type() expected:%d got:%d - %s", expected, got)
	}

	key = "zset-key"
	if _, e = client.Zadd(key, float64(0), []byte("woof")); e != nil {
		t.Errorf("on Zadd - %s", e)
	}
	expected = RT_ZSET
	got, e = client.Type(key)
	if e != nil {
		t.Errorf("on Type(%s) - %s", key, e)
	}
	if got != expected {
		t.Errorf("on Type() expected:%d got:%d - %s", expected, got)
	}

	flushAndQuitOnCompletion(t, client)
}

func TestSave(t *testing.T) {
	client, e := _test_getDefaultClient()
	if e != nil {
		t.Fatalf("on getDefaultClient", e)
	}

	if e := client.Save(); e != nil {
		t.Errorf("on Save() - %s", e)
	}

	flushAndQuitOnCompletion(t, client)
}

func TestAllKeys(t *testing.T) {
	client, e := _test_getDefaultClient()
	if e != nil {
		t.Fatalf("on getDefaultClient", e)
	}

	kvmap := testdata[_testdata_kv].(map[string][]byte)
	for k, v := range kvmap {
		if e := client.Set(k, v); e != nil {
			t.Errorf("on Set(%s, %s) - %s", k, v, e)
		}
	}
	got, e := client.AllKeys()
	if e != nil {
		t.Errorf("on AllKeys() - %s", e)
	}
	if got == nil {
		t.Errorf("on AllKeys() - got nil")
	}

	// if same length and all elements in kvmap, its ok
	if len(got) != len(kvmap) {
		t.Errorf("on AllKeys() - Len mismatch - got:%d expected:%d", len(got), len(kvmap))
	}
	for _, k := range got {
		if kvmap[k] == nil {
			t.Errorf("on AllKeys() - key %s is not in original kvmap", k)
		}
	}

	flushAndQuitOnCompletion(t, client)
}

func TestExists(t *testing.T) {
	client, e := _test_getDefaultClient()
	if e != nil {
		t.Fatalf("on getDefaultClient", e)
	}

	kvmap := testdata[_testdata_kv].(map[string][]byte)
	for k, v := range kvmap {
		res, e := client.Exists(k)
		if e != nil {
			t.Errorf("on Exists(%s) - %s", k, e)
		}
		if res {
			t.Errorf("on Exists(%s) - unexpected res True", k)
		}
		if e = client.Set(k, v); e != nil {
			t.Errorf("on Set(%s, %s) - %s", k, v, e)
		}
		res, e = client.Exists(k)
		if e != nil {
			t.Errorf("on Exists(%s) - %s", k, e)
		}
		if !res {
			t.Errorf("on Exists(%s)", k)
		}
	}

	flushAndQuitOnCompletion(t, client)
}

func compareStringArrays(got, expected []string) bool {
	if len(got) != len(expected) {
		return false
	}
	for i, b := range expected {
		if got[i] != b {
			return false
		}
	}
	return true
}
func compareByteArrays(got, expected []byte) bool {
	if len(got) != len(expected) {
		return false
	}
	for i, b := range expected {
		if got[i] != b {
			return false
		}
	}
	return true
}

/* --------------- KEEP THIS AS LAST FUNCTION -------------- */
func TestEnd_sct(t *testing.T) {
	log.Println("synchclient test completed")
}
