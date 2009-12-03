# Go-Redis

[Go][Go] Clients and Connectors for [Redis][Redis].  

The initial release provides the interface and implementation supporting the (~) full set of current Redis commands using synchrnous call semantics.  (Pipelines and asychronous goodness using the goroutines and channels is next.)

Hope to add rigorous tests as soon as I have a better understanding of the Go language tools.  Same applies to the makefile.  
Also am not sure regarding the efficiency of the implementation (for the obvious reasons), but definitely a goal is to make this a high performance connector.

## structure

The code is consolidated into a single 'redis' package and various elements of it are usable independently (for example if you wish to roll your own API but want to use the raw bytes protocol handling aspects).

## build and install

To build and install Go-Redis, from the root directory of the git clone:

	cd src/pkg/redis
	make clean && make install
	cd -

Confirm the install:

	ls -l $GOROOT/pkg/"$GOOS"_"$GOARCH"/redis.a


## run the benchmarks
	
After installing Go-Redis (per above), try (again from the root dir of Go-Redis):

	cd bench
	./runbench synchclient
	./runbench gosynchclient

## examples

Ciao.go is a sort of hello world for redis and should get you started for the barebones necessities of getting a client and issuing commands.

	cd examples
	./run ciao

[Go]: http://golang.org/
[Redis]: http://github.com/antirez/redis
