include $(GOROOT)/src/Make.$(GOARCH)

%:
	echo making $@
	$(GC) -I ./../src/pkg/redis -o $@.$(O) $@.go
	$(LD) -o $@ $@.$(O)

