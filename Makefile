include $(GOROOT)/src/Make.inc

TARG=redis
GOFILES=\
	reflect.go\
	future.go\
	specification.go\
	error.go\
	redis.go\
	redisfutures.go\
	protocol.go\
	connection.go\
	synchclient.go\
    asynchclient.go\
        
include $(GOROOT)/src/Make.pkg

