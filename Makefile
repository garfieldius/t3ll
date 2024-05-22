
GPG_KEY ?= 4C29A601B8AD9DFAE9641C0F0D1F16703AB055AA
VERSION_TAG ?= $(shell git tag -l | sort | tail -n 1)
VERSION ?= $(shell git tag -l | sort | tail -n 1 | sed -e 's,^v,,g')
YEAR ?= $(shell date +%Y)
GIT_REV ?= $(shell git rev-parse --short HEAD)
DIR ?= /usr/local/bin/

NODE_ENV = production
BUILDFLAGS = -ldflags "-X main.Version=main@$(GIT_REV) -X main.Year=$(YEAR)"
GPG_CMD = gpg --sign --detach-sign --armor --local-user $(GPG_KEY)
TAR_CMD =  tar --numeric-owner --create --gzip --file

ifneq ($(findstring debug,$(MAKECMDGOALS)),)
    NODE_ENV = development
endif

ifneq ($(findstring dist,$(MAKECMDGOALS)),)
	BUILDFLAGS = -ldflags "-w -s -X main.Version=$(VERSION) -X main.Year=$(YEAR)"
endif

.PHONY: build
build: t3ll

.PHONY: debug
debug: t3ll

.PHONY: clean
clean:
	rm -f t3ll t3ll.exe
	rm -rf dist
	rm -rf frontend/build server/index.html

clobber: clean
	rm -rf frontend/node_modules

.PHONY: install
install: t3ll
	install -m 0755 t3ll $(DIR)

.PHONY: debug-run
debug-run:
	$(MAKE) debug
	./t3ll --debug test.xlf

.PHONY: fix
fix:
	go fmt ./...

.PHONY: dist
dist: \
    dist/release_body.txt \
    dist/t3ll_linux_x64 dist/t3ll_linux_x64.sig \
    dist/t3ll_macos_x64 dist/t3ll_macos_x64.sig \
    dist/t3ll_macos_arm64 dist/t3ll_macos_arm64.sig \
    dist/t3ll_windows_x64.exe dist/t3ll_windows_x64.exe.sig \
    dist/sha256sum dist/sha256sum.sig \
    dist/t3ll-$(VERSION).x86_64_linux.bottle.tar.gz.sha256.txt dist/t3ll-$(VERSION).x86_64_linux.bottle.tar.gz.sha256.txt \
    dist/t3ll-$(VERSION).sierra.bottle.tar.gz.sha256.txt dist/t3ll-$(VERSION).sierra.bottle.tar.gz.sha256.txt \
    dist/t3ll-$(VERSION).arm64_big_sur.bottle.tar.gz.sha256.txt dist/t3ll-$(VERSION).arm64_big_sur.bottle.tar.gz.sha256.txt

t3ll: server/index.html $(shell find . -type f -iname "*.go")
	CGO_ENABLED=0 go build $(BUILDFLAGS)
	touch t3ll

server/index.html: frontend/build/index.html
	cp frontend/build/index.html server/index.html
	touch server/index.html

frontend/build/index.html: frontend/node_modules/.yarn-integrity $(shell find frontend/js -type f) $(shell find frontend/scss -type f) $(shell find frontend/templates -type f)
	cd frontend; NODE_ENV=$(NODE_ENV) yarn run gulp
	touch frontend/build/index.html

frontend/node_modules/.yarn-integrity:
	yarn --cwd frontend install --prefer-offline --frozen-lockfile
	touch frontend/node_modules/.yarn-integrity

dist/release_body.txt:
	mkdir -p dist
	git tag -l --format='%(contents:body)' $(VERSION_TAG) >> dist/release_body.txt

dist/t3ll_linux_x64: server/index.html
	mkdir -p dist && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILDFLAGS) -o dist/t3ll_linux_x64

dist/t3ll_macos_x64: server/index.html
	mkdir -p dist && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(BUILDFLAGS) -o dist/t3ll_macos_x64

dist/t3ll_macos_arm64: server/index.html
	mkdir -p dist && CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(BUILDFLAGS) -o dist/t3ll_macos_arm64

dist/t3ll_windows_x64.exe: server/index.html
	mkdir -p dist && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(BUILDFLAGS) -o dist/t3ll_windows_x64.exe

dist/t3ll_windows_x64.exe.sig: dist/t3ll_windows_x64.exe
	cd dist && $(GPG_CMD) --output t3ll_windows_x64.exe.sig t3ll_windows_x64.exe

dist/t3ll_macos_x64.sig: dist/t3ll_macos_x64
	cd dist && $(GPG_CMD) --output t3ll_macos_x64.sig t3ll_macos_x64

dist/t3ll_macos_arm64.sig: dist/t3ll_macos_arm64
	cd dist && $(GPG_CMD) --output t3ll_macos_arm64.sig t3ll_macos_arm64

dist/t3ll_linux_x64.sig: dist/t3ll_linux_x64
	cd dist && $(GPG_CMD) --output t3ll_linux_x64.sig t3ll_linux_x64

dist/sha256sum.sig: dist/sha256sum
	cd dist && $(GPG_CMD) --output sha256sum.sig sha256sum

dist/sha256sum: dist/t3ll_linux_x64 dist/t3ll_macos_x64 dist/t3ll_macos_arm64 dist/t3ll_windows_x64.exe
	cd dist && sha256sum t3ll_linux_x64 t3ll_macos_x64 t3ll_macos_arm64 t3ll_windows_x64.exe > sha256sum

dist/t3ll-$(VERSION).x86_64_linux.bottle.tar.gz.sha256.txt: dist/t3ll-$(VERSION).x86_64_linux.bottle.tar.gz
	printf $(shell sha256sum dist/t3ll-$(VERSION).x86_64_linux.bottle.tar.gz | cut -b-64) > dist/t3ll-$(VERSION).x86_64_linux.bottle.tar.gz.sha256.txt

dist/t3ll-$(VERSION).sierra.bottle.tar.gz.sha256.txt: dist/t3ll-$(VERSION).sierra.bottle.tar.gz
	printf $(shell sha256sum dist/t3ll-$(VERSION).sierra.bottle.tar.gz | cut -b-64) > dist/t3ll-$(VERSION).sierra.bottle.tar.gz.sha256.txt

dist/t3ll-$(VERSION).arm64_big_sur.bottle.tar.gz.sha256.txt: dist/t3ll-$(VERSION).arm64_big_sur.bottle.tar.gz
	printf $(shell sha256sum dist/t3ll-$(VERSION).arm64_big_sur.bottle.tar.gz | cut -b-64) > dist/t3ll-$(VERSION).arm64_big_sur.bottle.tar.gz.sha256.txt

dist/t3ll-$(VERSION).x86_64_linux.bottle.tar.gz: dist/.brew/linux/t3ll/$(VERSION)/bin/t3ll dist/.brew/linux/t3ll/$(VERSION)/README.md  dist/.brew/linux/t3ll/$(VERSION)/LICENSE
	$(TAR_CMD) dist/t3ll-$(VERSION).x86_64_linux.bottle.tar.gz -C dist/.brew/linux ./t3ll

dist/t3ll-$(VERSION).sierra.bottle.tar.gz: dist/.brew/macos-x64/t3ll/$(VERSION)/bin/t3ll dist/.brew/macos-x64/t3ll/$(VERSION)/README.md  dist/.brew/macos-x64/t3ll/$(VERSION)/LICENSE
	$(TAR_CMD) dist/t3ll-$(VERSION).sierra.bottle.tar.gz -C dist/.brew/macos-x64 ./t3ll

dist/t3ll-$(VERSION).arm64_big_sur.bottle.tar.gz: dist/.brew/macos-arm64/t3ll/$(VERSION)/bin/t3ll dist/.brew/macos-arm64/t3ll/$(VERSION)/README.md  dist/.brew/macos-arm64/t3ll/$(VERSION)/LICENSE
	$(TAR_CMD) dist/t3ll-$(VERSION).arm64_big_sur.bottle.tar.gz -C dist/.brew/macos-arm64 ./t3ll

dist/.brew/linux/t3ll/$(VERSION)/bin/t3ll: dist/t3ll_linux_x64
	mkdir -p dist/.brew/linux/t3ll/$(VERSION)/bin && cp dist/t3ll_linux_x64 dist/.brew/linux/t3ll/$(VERSION)/bin/t3ll

dist/.brew/macos-x64/t3ll/$(VERSION)/bin/t3ll: dist/t3ll_macos_x64
	mkdir -p dist/.brew/macos-x64/t3ll/$(VERSION)/bin && cp dist/t3ll_macos_x64 dist/.brew/macos-x64/t3ll/$(VERSION)/bin/t3ll

dist/.brew/macos-arm64/t3ll/$(VERSION)/bin/t3ll: dist/t3ll_macos_arm64
	mkdir -p dist/.brew/macos-arm64/t3ll/$(VERSION)/bin && cp dist/t3ll_macos_arm64 dist/.brew/macos-arm64/t3ll/$(VERSION)/bin/t3ll

dist/.brew/macos-x64/t3ll/$(VERSION)/README.md:
	mkdir -p dist/.brew/macos-x64/t3ll/$(VERSION) && cp -a README.md dist/.brew/macos-x64/t3ll/$(VERSION)/README.md

dist/.brew/macos-arm64/t3ll/$(VERSION)/README.md:
	mkdir -p dist/.brew/macos-arm64/t3ll/$(VERSION) && cp -a README.md dist/.brew/macos-arm64/t3ll/$(VERSION)/README.md

dist/.brew/linux/t3ll/$(VERSION)/README.md:
	mkdir -p dist/.brew/linux/t3ll/$(VERSION) && cp -a README.md dist/.brew/linux/t3ll/$(VERSION)/README.md

dist/.brew/macos-x64/t3ll/$(VERSION)/LICENSE:
	mkdir -p dist/.brew/macos-x64/t3ll/$(VERSION) && cp -a LICENSE dist/.brew/macos-x64/t3ll/$(VERSION)/LICENSE

dist/.brew/macos-arm64/t3ll/$(VERSION)/LICENSE:
	mkdir -p dist/.brew/macos-arm64/t3ll/$(VERSION) && cp -a LICENSE dist/.brew/macos-arm64/t3ll/$(VERSION)/LICENSE

dist/.brew/linux/t3ll/$(VERSION)/LICENSE:
	mkdir -p dist/.brew/linux/t3ll/$(VERSION) && cp -a LICENSE dist/.brew/linux/t3ll/$(VERSION)/LICENSE
