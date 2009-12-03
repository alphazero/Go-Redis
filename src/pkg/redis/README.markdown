To build and use from other go programs:  

    make clean && make install

This will install redis.a in your platform specific packages folder.  (For example, on my machine, it is under $GOROOT/pkg/darwin_amd64/redis.a).

Tests will be added soon and you can test the package by issuing gotest in this directory.
