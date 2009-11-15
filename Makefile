.SUFFIXES: .go .6
 
OBJS=reflect_utils.6 redis.6 protocol.6 connection.6 client.6 test.6
LIB=.
TARGET=test
 
build: $(OBJS)
	6l -o $(TARGET) test.6
 
clean:
	rm -f $(OBJS) $(TARGET) 
 
.go.6:
	6g -I $(LIB) $<
