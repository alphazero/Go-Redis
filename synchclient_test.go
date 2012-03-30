package redis

import (
	"log"
	"testing"
)

type tspec_sct struct {
	host     string
	port     int
	password string
	db       int
	connspec *ConnectionSpec
}

func testspec_sct() tspec_sct {
	var spec tspec_sct

	spec.host = "localhost"
	spec.port = 6379
	spec.db = 13
	spec.password = "go-redis"

	spec.connspec = DefaultSpec().Host(spec.host).Port(spec.port).Db(spec.db).Password(spec.password)
	return spec
}

func TestClientConnectWithSpec(t *testing.T) {
	spec := testspec_sct()

	c, err := NewSynchClientWithSpec(spec.connspec)
	if err != nil {
		t.Error("failed to create client with spec. Error: %s", err.Message())
	} else if c == nil {
		t.Error("failed to create client with spec. Error: %s")
	}
}
func TestEnd_sct(t *testing.T) {
	log.Println("synchclient test")
}
