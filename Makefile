##
# goimapnotify
#
# @file
# @version 0.1

build:
	go build -o goimapnotify

test: install-deps gen-mocks
	go test -v ./... -coverprofile coverage.out
	uncover -min 80.0 coverage.out

install-deps:
	go install github.com/gregoryv/uncover/cmd/uncover@v0.7.0
	go install github.com/subtle-byte/mockigo/cmd/mockigo@latest

gen-mocks:
	mockigo

# end
