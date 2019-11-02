prefix ?= /usr/local
exec_prefix ?= $(prefix)
bindir ?= $(exec_prefix)/bin

.PHONY: all install uninstall test update-vendor

all:
	env GO111MODULE=on go build -mod vendor -v ./cmd/...

install:
	env GO111MODULE=on GOBIN=$(bindir) go install -mod vendor -v ./cmd/...

uninstall:
	rm -f $(bindir)/scrit-engine $(bindir)/scrit-gov $(bindir)/scrit-mint $(bindir)/scrit-wallet

test:
	go get github.com/frankbraun/gocheck
	gocheck -g -c -v cmd

update-vendor:
	rm -rf vendor
	env GO111MODULE=on go get -u ./cmd/...
	env GO111MODULE=on go mod tidy -v
	env GO111MODULE=on go mod vendor
