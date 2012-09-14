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
		t.Error("failed to create client with spec. Error: ", err.Message())
	} else if client == nil {
		t.Error("BUG: client is nil")
	}
	client.Quit()
}

/* --------------- KEEP THIS AS LAST FUNCTION -------------- */
func TestEnd_asct(t *testing.T) {
	log.Println("-- asynchclient test completed")
}
