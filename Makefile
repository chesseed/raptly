VERSION=$(shell make -s version)
GOOS=$(shell go env GOHOSTOS)
GOARCH=$(shell go env GOHOSTARCH)

dpkg:
	@test -n "$(DEBARCH)" || (echo "please define DEBARCH"; exit 1)
	# Run dpkg-buildpackage
	@buildtype="any" ; \
	if [ "$(DEBARCH)" = "amd64" ]; then  \
	  buildtype="any,all" ; \
	fi ; \
	echo "\e[33m\e[1mBuilding: $$buildtype\e[0m" ; \
	cmd="dpkg-buildpackage -us -uc --build=$$buildtype -d --host-arch=$(DEBARCH)" ; \
	echo "$$cmd" ; \
	$$cmd
	lintian ../*_$(DEBARCH).changes || true
	# cleanup
	@test ! -f debian/changelog.dpkg-bak || mv debian/changelog.dpkg-bak debian/changelog; \
	mkdir -p build && mv ../*.deb build/ ; \
	cd build && ls -l *.deb


version:  ## Print raptly version
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

lint:
	$(go env GOPATH)/bin/golangci-lint run

.PHONY: pack version raptly lint