COMPONENTS=proto sync proxy

all: $(COMPONENTS)

$(COMPONENTS):
	make -C $@

publish:
	for c in $(COMPONENTS); do make -C $$c publish; done

clean:
	for c in $(COMPONENTS); do make -C $$c clean; done

build_iso:
	make -C utils iso

proto:
	make -C proto

ci-test: proto
	make -C tools/compose common_test sync_test proxy_test

.PHONY: all $(COMPONENTS) publish build_iso ci-test clean
.SILENT:
