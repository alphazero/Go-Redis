include $(GOROOT)/src/Make.$(GOARCH)

TARG=redis
GOFILES=\
	reflect.go\
	future.go\
	specification.go\
	error.go\
	redis.go\
	protocol.go\
	connection.go\
	synchclient.go\
    asynchclient.go\

$(TARG).$(O): $(GOFILES)
	$(GC) -o $(TARG).$(O) $(GOFILES)
        
include $(GOROOT)/src/Make.pkg

