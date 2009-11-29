# Go-Redis

[Go][Go] Clients and Connectors for [Redis][Redis].  

The initial release provides the interface and implementation supporting the (~) full set of current Redis commands using synchrnous call semantics.  (Pipelines and asychronous goodness using the goroutines and channels is next.)

Hope to add rigorous tests as soon as I have a better understanding of the Go language tools.  Same applies to the makefile.  
Also am not sure regarding the efficiency of the implementation (for the obvious reasons), but definitely a goal is to make this a high performance connector.

## structure

The code is consolidated into a single 'redis' package and various elements of it are usable independently (for example if you wish to roll your own API but want to use the raw bytes protocol handling aspects).

## quickstart

Presuming you already have both Go and Redis installed and a running local Redis server:

	cd src/pkg/redis
	make
	cd ../../../bench
	./runbench synchclient
	./runbench gosynchclient


[Go]: http://golang.org/
[Redis]: http://github.com/antirez/redis
