include $(GOROOT)/src/Make.inc

%:
	echo making $@
	$(GC) -o $@.$(O) $@.go
#	$(GC) -I ./../src/pkg/redis -o $@.$(O) $@.go
	$(LD) -o $@ $@.$(O)

