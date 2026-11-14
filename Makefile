## Инструментарий проектов webnice - единый Makefile.
## (c)Alex Geer <monoflash@gmail.com>
## Версия: 14.11.2026.

## Вычисление текущей директории проекта.
export DIR := $(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

## Подключение переменных окружения проекта.
ifneq ("$(wildcard $(DIR)/.env.mk)","")
include $(DIR)/.env.mk
endif
## Подключение makefile части, содержащей специфичные для проекта команды.
ifneq ("$(wildcard $(DIR)/.prj.mk)","")
include $(DIR)/.prj.mk
endif

GOPATH := $(GOPATH)
DATE   := $(shell date -u +%Y%m%d.%H%M%S.%Z)

## Имя приложения по умолчанию - название корневой директории проекта.
APP ?= $(shell basename "$(DIR)")
## Имя go модуля проекта.
APP_PACKAGE := $(shell grep -m1 -E '^module[[:space:]]+' "$(DIR)/go.mod" 2>/dev/null | awk '{ print $$2 }')
## Директория размещения бинарных файлов.
export GOBIN := $(DIR)/bin
## Флаги для создания бинарника (производственный режим).
export LDFLAGS_MODE := $(shell if [ -n "$$GO_SERVER_URL" ]; then echo '-s -w'; fi)
## Флаги сборки с простановкой версии и даты сборки.
LDFLAGS := -X main.build=$(DATE) $(LDFLAGS_MODE)

## Пакеты кодогенерации (.gogenerate).
GOGENERATE = $(shell if [ -f "$(DIR)/.gogenerate" ]; then cat "$(DIR)/.gogenerate"; fi)
## Тестируемые пакеты (.testpackages).
TESTPACKETS = $(shell if [ -f "$(DIR)/.testpackages" ]; then cat "$(DIR)/.testpackages"; fi)
## Пакеты тестирования производительности (.benchpackages).
BENCHPACKETS = $(shell if [ -f "$(DIR)/.benchpackages" ]; then cat "$(DIR)/.benchpackages"; fi)
## Зависимости с зафиксированной версией (.packages_lock или .packages_lock_version).
PACKAGES_LOCK = $(shell for f in "$(DIR)/.packages_lock" "$(DIR)/.packages_lock_version"; do if [ -f $$f ]; then cat $$f; break; fi; done)
## Токен сервиса coveralls.io (.coveralls).
COVERALLS_REPO_TOKEN = $(shell if [ -f "$(DIR)/.coveralls" ]; then cat "$(DIR)/.coveralls"; fi)
## Локальные зависимости разработки (.localpackages).
LOCALPACKAGES = $(shell if [ -f "$(DIR)/.localpackages" ]; then cat "$(DIR)/.localpackages"; fi)

## Сценарий по умолчанию - отображение доступных команд.
default: help

## Подготовка проекта к сборке и тестированию.
## Переопределяется макросом PROJECT_LINK в файле .prj.mk.
link:
	$(call PROJECT_LINK)
.PHONY: link

## Создание директорий проекта и очистка устаревшей vendor директории.
dep-init:
	@rm -rf "$(DIR)/vendor" 2>/dev/null; true
	$(call PROJECT_FOLDERS)
.PHONY: dep-init

## Загрузка и обновление зависимостей проекта (продакшн режим).
dep: dep-init
	@for item in $(LOCALPACKAGES); do \
		PKGNAME=`echo "$${item}" | awk -F'=' '{ print $$1 }'`; \
		go mod edit -dropreplace "$${PKGNAME}" || true; \
	done
	@go mod download
	@go mod tidy
	@go mod vendor
	$(call PROJECT_DEPENDENCIES)
.PHONY: dep

## Загрузка и обновление зависимостей проекта (режим разработки, с локальными заменами).
dep-dev: dep-init
	@for item in $(LOCALPACKAGES); do \
		PKGNAME=`echo "$${item}" | awk -F'=' '{ print $$1 }'`; \
		REPLACE=`echo "$${item}" | awk -F'=' '{ print $$2 }'`; \
		go mod edit -replace "$${PKGNAME}=$${REPLACE}" || true; \
	done
	@go mod download
	@go mod tidy
	$(call PROJECT_DEPENDENCIES_DEVELOPMENT)
.PHONY: dep-dev

## Обновление всех зависимостей до последних версий.
update: upd
upd:
	@go clean -cache -modcache
	@go get -u ./...
	@for item in $(PACKAGES_LOCK); do \
		go get -u "$${item}" || true; \
	done
.PHONY: upd
.PHONY: update

## Кодогенерация (использовать только во время разработки).
## Все генерируемые пакеты перечислены в файле .gogenerate.
gen:
	@for PKGNAME in $(GOGENERATE); do \
		DB2STRUCT_DRV='$(MIGRATION_DRV_MYSQL)' \
		DB2STRUCT_DSN='$(MIGRATION_DSN_MYSQL)' \
		go generate -v "$${PKGNAME}" || exit 1; \
	done
.PHONY: gen

## Алиас команды gen.
generate: gen
.PHONY: generate

## Тестирование одного или нескольких пакетов с отчётом о проценте покрытия кода тестами.
## Все тестируемые пакеты перечислены в файле .testpackages.
test: link
	@echo "mode: count" > "$(DIR)/coverage.log"
	@for PACKET in $(TESTPACKETS); do \
		touch "$(DIR)/coverage-tmp.log"; \
		go test -v -covermode=count -coverprofile="$(DIR)/coverage-tmp.log" $${PACKET} || exit 1; \
		tail -n +2 "$(DIR)/coverage-tmp.log" | sort -r | awk '{ if ($$1 != last) { print $$0; last=$$1 } }' >> "$(DIR)/coverage.log"; \
		rm -f "$(DIR)/coverage-tmp.log"; true; \
	done
	@if [ -n "$(COVERALLS_REPO_TOKEN)" ] && command -v goveralls >/dev/null 2>&1; then \
		echo "Загрузка на coveralls.io процента покрытия кода тестами."; \
		goveralls -repotoken "$(COVERALLS_REPO_TOKEN)" >/dev/null 2>&1 || true; \
	fi
.PHONY: test

## Отображение в браузере покрытия тестируемого кода в HTML-отчёте (запускать только во время разработки).
cover: test
	@go tool cover -html="$(DIR)/coverage.log"
.PHONY: cover

## Тестирование производительности.
## Все тестируемые пакеты перечислены в файле .benchpackages.
bench:
	@for PACKET in $(BENCHPACKETS); do \
		go test -race -bench=. -benchmem $${PACKET}; \
	done
.PHONY: bench

## Очистка экрана консоли.
clear:
	@clear 2>/dev/null || true
.PHONY: clear

## Очистка директории проекта от временных файлов.
clean:
	@if [ -d "$(DIR)/pkg" ]; then \
		chown -R `whoami` "$(DIR)/pkg/" 2>/dev/null; true; \
		chmod -R 0777 "$(DIR)/pkg/" 2>/dev/null; true; \
	fi
	@rm -rf "$(DIR)/bin" 2>/dev/null; true
	@rm -rf "$(DIR)/pkg" 2>/dev/null; true
	@rm -rf "$(DIR)/rpmbuild" 2>/dev/null; true
	@rm -rf "$(DIR)/run"/*.pid 2>/dev/null; true
	@rm -rf "$(DIR)/log"/*.log 2>/dev/null; true
	@rm -rf "$(DIR)"/*.log 2>/dev/null; true
	@rm -rf "$(DIR)"/*.lock 2>/dev/null; true
	@rm -rf "$(DIR)"/*.db 2>/dev/null; true
.PHONY: clean

## Помощь по командам.
help:
	@echo "Usage: make [target]"
	@echo "  target is:"
	@echo "    dep-init             - Создание директорий проекта и очистка vendor."
	@echo "    dep                  - Загрузка и обновление зависимостей проекта."
	@echo "    dep-dev              - Загрузка зависимостей с локальными заменами (режим разработки)."
	@echo "    upd или update       - Обновление всех зависимостей до последних версий."
	@echo "    gen                  - Кодогенерация с использованием go generate."
	@echo "    test                 - Запуск тестов проекта с отчётом о покрытии кода."
	@echo "    cover                - Показ процента покрытия кода тестами в браузере."
	@echo "    bench                - Запуск тестов производительности."
	@echo "    clear                - Очистка экрана консоли."
	@echo "    clean                - Очистка директории проекта от временных файлов."
	@$(call PROJECT_HELP)
.PHONY: help
