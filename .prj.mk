## Project tooling makefile
##
## Проектные переопределения команд Makefile.
## Команды, не используемые проектом, оставлены пустыми.

## Встраивание статических ресурсов в приложение.
embed:
.PHONY: embed

## Компиляция проекта.
build:
.PHONY: build

## Запуск приложения в режиме разработки.
dev:
.PHONY: dev

## Запуск приложения в продакшн режиме.
run:
.PHONY: run

## Принудительное завершение процесса запущенного приложения.
kill:
.PHONY: kill

## Вывод на экран версии приложения.
version: v
.PHONY: version
v:
.PHONY: v

## Сборка RPM пакета.
rpm:
.PHONY: rpm

## Создание директорий миграций.
migration-mkdir:
.PHONY: migration-mkdir

## Команды инструментов миграций.
migration-commands:
.PHONY: migration-commands
