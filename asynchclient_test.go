package redis

import (
	//	"fmt"
	"log"
	"testing"
)

// Check that connection is actually passing passwords from spec
// and catching AUTH ERRs.
func TestAsyncClientConnectWithBadSpec(t *testing.T) {
	spec := _test_getDefConnSpec()
	spec.Password("bad-password")
	client, expected := NewAsynchClientWithSpec(spec)
	if expected == nil {
		t.Error("BUG: Expected a RedisError")
	}
	if client != nil {
		t.Error("BUG: async client reference on error MUST be nil")
	}
}

// Check that connection is actually passing passwords from spec
func TestAsyncClientConnectWithSpec(t *testing.T) {
	spec := _test_getDefConnSpec()

	client, err := NewAsynchClientWithSpec(spec)
	if err != nil {
		t.Fatalf("failed to create client with spec. Error: %s ", err)
	} else if client == nil {
		t.Fatal("BUG: client is nil")
	}

	// quit once -- OK
	futureBool, err := client.Quit()
	if err != nil {
		t.Errorf("BUG - initial Quit on asyncClient should not return error - %s ", err)
	}
	if futureBool == nil {
		t.Errorf("BUG - non-error asyncClient response should not return nil future")
	}
	// block until we get results
	ok, fe := futureBool.Get()
	if fe != nil {
		t.Errorf("BUG - non-Error Quit future result get must never return error - got: %s", fe)
	}
	if !ok {
		t.Errorf("BUG - non-Error Quit future result must always be true ")
	}

	// subsequent quit should raise error
	for i := 0; i < 10; i++ {
		futureBool, err = client.Quit()
		if err == nil {
			t.Errorf("BUG - Quit on shutdown asyncClient should return error")
		}
		if futureBool != nil {
			t.Errorf("BUG - Quit on shutdown asyncClient should not return future. got: %s", futureBool)
		}
	}
}

func TestAsyncMget(t *testing.T) {
	client, e := _test_getDefaultAsyncClient()
	if e != nil {
		t.Fatalf("on getDefaultClient - %s", e)
	}

	client.Quit()

}

/* --------------- KEEP THIS AS LAST FUNCTION -------------- */
func TestEnd_asct(t *testing.T) {
	log.Println("-- asynchclient test completed")
}
