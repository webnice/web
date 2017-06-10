DIR=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

GOPATH := $(DIR):$(GOPATH)
DATE=$(shell date -u +%Y%m%d.%H%M%S.%Z)
PACKETS=$(shell cat .testpackages)

default: lint test

## Generate code by go generate or other utilities
generate:
	#GOPATH=${GOPATH} go generate
	#GOPATH=${GOPATH} easyjson -output_filename configuration.go src/gopkg.in/webnice/web.v1/types.go
.PHONY: generate

## Dependence managers
dep:
	mkdir -p src/gopkg.in/webnice; cd src/gopkg.in/webnice && ln -s ../../.. web.v1 2>/dev/null; true
	GOPATH=${GOPATH} glide install
.PHONY: dep

test: dep
	clear
	echo "mode: set" > coverage.log
	for PACKET in $(PACKETS); do \
		touch coverage-tmp.log; \
		GOPATH=${GOPATH} go test -v -covermode=count -coverprofile=coverage-tmp.log $$PACKET; \
		if [ ! $$? == 0 ]; then exit $$?; fi; \
		tail -n +2 coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> coverage.log; \
		rm -f coverage-tmp.log; true; \
	done
.PHONY: test

cover: test
	GOPATH=${GOPATH} go tool cover -html=coverage.log
.PHONY: cover

bench: dep
	GOPATH=${GOPATH} go test -race -bench=. -benchmem ./...
.PHONY: bench

lint:
	gometalinter \
	--vendor \
	--deadline=15m \
	--cyclo-over=20 \
	--disable=aligncheck \
	--disable=gotype \
	--skip=src/vendor \
	--linter="vet:go tool vet -printf {path}/*.go:PATH:LINE:MESSAGE" \
	./...
.PHONY: lint

clean:
	rm -rf ${DIR}/src; true
	rm -rf ${DIR}/vendor; true
	rm -rf ${DIR}/bin/*; true
	rm -rf ${DIR}/pkg/*; true
	rm -rf ${DIR}/*.log; true
	rm -rf ${DIR}/*.lock; true
.PHONY: clean
