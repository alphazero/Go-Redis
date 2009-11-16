include $(GOROOT)/src/Make.$(GOARCH)

TARG=redis
REDIS_FILES=\
	reflect_utils.go\
	specification.go\
	error.go\
	redis.go\
	protocol.go\
	connection.go\
	synchclient.go\

redis.$(O): $(REDIS_FILES)
	$(GC) -o redis.$(O) $(REDIS_FILES)
        
include $(GOROOT)/src/Make.pkg

