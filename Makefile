##
# goimapnotify
#
# @file
# @version 0.1

# Definir las variables para la información de Git
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_TAG := $(shell git describe --tags)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

# Definir las flags del linker
LDFLAGS := -X main.commit=$(GIT_COMMIT) -X main.gittag=$(GIT_TAG) -X main.branch=$(GIT_BRANCH)

build:
	go build -ldflags "$(LDFLAGS)"

changelog:
	git-chglog -o README.md 2.3.14..

# end
