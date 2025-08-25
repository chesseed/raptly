VERSION=$(shell make -s version)
GOOS=$(shell go env GOHOSTOS)
GOARCH=$(shell go env GOHOSTARCH)

version:  ## Print aptly version
	if which dpkg-parsechangelog > /dev/null 2>&1; then \
		echo `dpkg-parsechangelog -S Version`$$ci; \
	else \
		echo `grep ^raptly -m1  debian/changelog | sed 's/.*(\([^)]\+\)).*/\1/'`$$ci ; \
	fi

raptly:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.Version=$(VERSION)" -o build/raptly ./cmd/raptly

pack: raptly
	@path="raptly_$(VERSION)_$(GOOS)_$(GOARCH)"; \
	rm -rf "build/$$path"; \
	rm -rf "build/$$path".zip; \
	cd build; \
	zip -r "$$path".zip "raptly" > /dev/null \
		&& echo "Built build/$${path}.zip"; \

.PHONY: pack version raptly