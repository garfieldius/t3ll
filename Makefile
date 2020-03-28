
NODE_ENV = production
BUILDFLAGS = -ldflags "-w -s -X main.Version=$(shell git tag -l | sort | tail -n 1)"

ifneq ($(findstring debug,$(MAKECMDGOALS)),)
	NODE_ENV = development
	BUILDFLAGS = -tags debug -ldflags "-X main.Version=master@$(shell git rev-parse --short HEAD) -X github.com/garfieldius/t3ll/server.Dir=$(shell pwd)"
endif

.PHONY: build
build: t3ll

.PHONY: debug
debug: t3ll

.PHONY: clean
clean:
	rm -f t3ll t3ll.exe
	rm -rf frontend/build dist
	rm -f server/html.go

clobber: clean
	rm -rf frontend/node_modules

.PHONY: install
install: t3ll
	install -m 0755 t3ll /usr/local/bin/

.PHONY: dist
dist: dist/t3ll_linux_x64.sig dist/t3ll_macosx_x64.sig dist/t3ll_windows_x64.exe.sig dist/sha256sum.sig

t3ll: server/html.go frontend/build/index.html
	go build $(BUILDFLAGS)

frontend/build/index.html: frontend/node_modules/.bin/gulp
	cd frontend; NODE_ENV=$(NODE_ENV) yarn run gulp

frontend/node_modules/.bin/gulp:
	cd frontend; yarn install --prefer-offline --frozen-lockfile

server/html.go: frontend/build/index.html
	node embed-html.js

dist/t3ll_linux_x64: server/html.go
	mkdir -p dist && GOOS=linux GOARCH=amd64 go build $(BUILDFLAGS) -o dist/t3ll_linux_x64

dist/t3ll_macosx_x64: server/html.go
	mkdir -p dist && GOOS=darwin GOARCH=amd64 go build $(BUILDFLAGS) -o dist/t3ll_macosx_x64

dist/t3ll_windows_x64.exe: server/html.go
	mkdir -p dist && GOOS=windows GOARCH=amd64 go build $(BUILDFLAGS) -o dist/t3ll_windows_x64.exe

dist/t3ll_windows_x64.exe.sig: dist/t3ll_windows_x64.exe
	cd dist && gpg -b -a -o t3ll_windows_x64.exe.sig t3ll_windows_x64.exe

dist/t3ll_macosx_x64.sig: dist/t3ll_macosx_x64
	cd dist && gpg -b -a -o t3ll_macosx_x64.sig t3ll_macosx_x64

dist/t3ll_linux_x64.sig: dist/t3ll_linux_x64
	cd dist && gpg -b -a -o t3ll_linux_x64.sig t3ll_linux_x64

dist/sha256sum.sig: dist/sha256sum
	cd dist && gpg -b -a -o sha256sum.sig sha256sum

dist/sha256sum: dist/t3ll_linux_x64 dist/t3ll_macosx_x64 dist/t3ll_windows_x64.exe
	cd dist && sha256sum t3ll_linux_x64 t3ll_macosx_x64 t3ll_windows_x64.exe > sha256sum
