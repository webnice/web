## Simple projects tooling for every day
## (c)Alex Geer <monoflash@gmail.com>
## Version: 2018.11.10

## Project source directory path
DIR         := $(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

GOPATH      := $(DIR):$(GOPATH)
DATE        := $(shell date -u +%Y%m%d.%H%M%S.%Z)
GOGENERATE   = $(shell if [ -f .gogenerate ]; then cat .gogenerate; fi)
TESTPACKETS  = $(shell if [ -f .testpackages ]; then cat .testpackages; fi)
BENCHPACKETS = $(shell if [ -f .benchpackages ]; then cat .benchpackages; fi)
GO111MODULE ?= $(GO111MODULE:on)

default: help

link:
	@mkdir -p src/gopkg.in/webnice; cd src/gopkg.in/webnice && ln -s ../../.. web.v1 2>/dev/null; true
.PHONY: link

## Code generation (run only during development)
# All generating files are included in a .gogenerate file
gen:
	@for PKGNAME in $(GOGENERATE); do GOPATH="$(DIR)" go generate $${PKGNAME}; done
.PHONY: gen

## Testing one or multiple packages as well as applications with reporting on the percentage of test coverage
# All testing files are included in a .testpackages file
test: link
	@echo "mode: set" > coverage.log
	@for PACKET in $(TESTPACKETS); do \
		touch coverage-tmp.log; \
		GOPATH=${GOPATH} go test -v -covermode=count -coverprofile=coverage-tmp.log $$PACKET; \
		if [ "$$?" -ne "0" ]; then exit $$?; fi; \
		tail -n +2 coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> coverage.log; \
		rm -f coverage-tmp.log; true; \
	done
.PHONY: test

	## Displaying in the browser coverage of tested code, on the html report (run only during development)
cover: test
	@GOPATH=${GOPATH} go tool cover -html=$(DIR)/coverage.log
.PHONY: cover

## Performance testing
# All testing files are included in a .benchpackages file
bench:
	@for PACKET in $(BENCHPACKETS); do GOPATH=${GOPATH} go test -race -bench=. -benchmem $$PACKET; done
.PHONY: bench

## Code quality testing
# https://github.com/alecthomas/gometalinter/
# install: curl -L https://git.io/vp6lP | sh
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

## Clearing project temporary files
clean:
	@if [ -d ${DIR}/pkg ]; then\
		chown -R `whoami` ${DIR}/pkg/; true; \
		chmod -R 0777 ${DIR}/pkg/; true; \
	fi
	@rm -rf ${DIR}/src; true
	@rm -rf ${DIR}/bin; true
	@rm -rf ${DIR}/pkg; true
	@rm -rf ${DIR}/rpmbuild; true
	@rm -rf ${DIR}/*.log; true
	@rm -rf ${DIR}/*.lock; true
.PHONY: clean

## Help for main targets
help:
	@echo "Usage: make [target]"
	@echo "  target is:"
	@echo "    gen                  - Кодогенерация с использованием go generate"
	@echo "    test                 - Запуск тестов проекта"
	@echo "    cover                - Запуск тестов проекта с отображением процента покрытия кода тестами"
	@echo "    bench                - Запуск тестов производительности"
	@echo "    lint                 - Запуск проверки кода с помощью gometalinter"
	@echo "    clean                - Очистка папки проекта от временных файлов"
.PHONY: help
