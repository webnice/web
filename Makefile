DIR=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

GOPATH := $(DIR):$(GOPATH)
DATE=$(shell date -u +%Y%m%d.%H%M%S.%Z)
GOGENERATE=$(shell if [ -f .gogenerate ]; then cat .gogenerate; fi)
PACKETS=$(shell cat .testpackages)

default: lint test

link:
	mkdir -p src/gopkg.in/webnice; cd src/gopkg.in/webnice && ln -s ../../.. web.v1 2>/dev/null; true
.PHONY: link

## Generate code by go generate or other utilities
generate: link
	for PKGNAME in $(GOGENERATE); do GOPATH="$(DIR)" go generate $${PKGNAME}; done
.PHONY: generate

## Dependence managers
dep: link
	if command -v "gvt"; then GOPATH="$(DIR)" gvt update -all; fi
	rm -rf vendor/golang.org/x/text/cmd 2>/dev/null; true
	rm -rf vendor/golang.org/x/text/collate/tools 2>/dev/null; true
	rm -rf vendor/golang.org/x/net/http2/h2i 2>/dev/null; true
.PHONY: dep

test: link
	echo "mode: set" > coverage.log
	for PACKET in $(PACKETS); do \
		touch coverage-tmp.log; \
		GOPATH=${GOPATH} go test -v -covermode=count -coverprofile=coverage-tmp.log $$PACKET; \
		if [ "$$?" -ne "0" ]; then exit $$?; fi; \
		tail -n +2 coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> coverage.log; \
		rm -f coverage-tmp.log; true; \
	done
.PHONY: test

cover: test
	GOPATH=${GOPATH} go tool cover -html=coverage.log
.PHONY: cover

bench: link
	GOPATH=${GOPATH} go test -race -bench=. -benchmem ./...
.PHONY: bench

lint: link
	gometalinter \
	--vendor \
	--deadline=15m \
	--cyclo-over=20 \
	--disable=aligncheck \
	--disable=gotype \
	--skip=vendor \
	--skip=src/vendor \
	--linter="vet:go tool vet -printf {path}/*.go:PATH:LINE:MESSAGE" \
	./...
.PHONY: lint

clean:
	rm -rf ${DIR}/src; true
	rm -rf ${DIR}/bin/*; true
	rm -rf ${DIR}/pkg/*; true
	rm -rf ${DIR}/*.log; true
	rm -rf ${DIR}/*.lock; true
.PHONY: clean
