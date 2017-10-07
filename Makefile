BUILDCMD:=go build -a -tags netgo -ldflags "-w -s -extldflags '-static -s -pie'"
GOPKGS:= ./vendor/github.com/kr/pretty/formatter.go

build: ./t3ll

release: release/t3ll_linux_x64.tar.gz \
         release/t3ll_linux_i386.tar.gz \
         release/t3ll_macos_x64.zip \
         release/t3ll_macos_i386.zip \
         release/t3ll_windows_x64.zip \
         release/t3ll_windows_i386.zip

clean:
	rm -rf release
	rm -f t3ll t3ll.exe
	cd frontend; rm -rf _dev _live

clobber: clean
	cd frontend; rm -rf node_modules
	rm -rf vendor

godeps: $(GOPATH)/bin/go-bindata \
        $(GOPKGS)

install: $(GOPATH)/bin/t3ll

fmt:
	go fmt ./...

debug: godeps frontend/node_modules/.bin/gulp
	cd frontend && rm -rf _dev _live && yarn exec gulp html-debug && cd .. && \
	rm -f server/data.go && \
	go-bindata -pkg server -o server/data.go -prefix frontend/_dev/ frontend/_dev/editor.html && \
	go build -a -tags 'debug'

.PHONY: default build release deps clean godeps install

release/t3ll_linux_x64.tar.gz: release/linux_x64/t3ll
	cd release/linux_x64 && tar -cf ../t3ll_linux_x64.tar.gz t3ll

release/linux_x64/t3ll: godeps server/data.go
	mkdir -p release/linux_x64 && GOOS="linux" GOARCH="amd64" $(BUILDCMD) -o release/linux_x64/t3ll

release/t3ll_linux_i386.tar.gz: release/linux_i386/t3ll
	cd release/linux_i386 && tar -cf ../t3ll_linux_i386.tar.gz t3ll

release/linux_i386/t3ll: godeps server/data.go
	mkdir -p release/linux_i386 && GOOS="linux" GOARCH="386" $(BUILDCMD) -o release/linux_i386/t3ll

release/t3ll_macos_x64.zip: release/macos_x64/t3ll
	cd release/macos_x64 && zip -q9 ../t3ll_macos_x64.zip t3ll

release/macos_x64/t3ll: godeps server/data.go
	mkdir -p release/macos_x64 && GOOS="darwin" GOARCH="amd64" $(BUILDCMD) -o release/macos_x64/t3ll

release/t3ll_macos_i386.zip: release/macos_i386/t3ll
	cd release/macos_i386 && zip -q9 ../t3ll_macos_i386.zip t3ll

release/macos_i386/t3ll: godeps server/data.go
	mkdir -p release/macos_i386 && GOOS="darwin" GOARCH="386" $(BUILDCMD) -o release/macos_i386/t3ll

release/t3ll_windows_x64.zip: release/windows_x64/t3ll.exe
	cd release/windows_x64 && zip -q9 ../t3ll_windows_x64.zip t3ll.exe

release/windows_x64/t3ll.exe: godeps server/data.go
	mkdir -p release/windows_x64 && GOOS="windows" GOARCH="amd64" $(BUILDCMD) -o release/windows_x64/t3ll.exe

release/t3ll_windows_i386.zip: release/windows_i386/t3ll.exe
	cd release/windows_i386 && zip -q9 ../t3ll_windows_i386.zip t3ll.exe

release/windows_i386/t3ll.exe: godeps server/data.go
	mkdir -p release/windows_i386 && GOOS="windows" GOARCH="386" $(BUILDCMD) -o release/windows_i386/t3ll.exe

server/data.go: frontend/_live/editor.html
	go-bindata -nometadata -pkg server -o server/data.go \
		-prefix frontend/_live/ frontend/_live/editor.html

frontend/_live/editor.html: frontend/node_modules/.bin/gulp
	cd frontend; yarn exec gulp clean; yarn exec gulp live

frontend/node_modules/.bin/gulp:
	cd frontend; yarn install

./t3ll: godeps server/data.go
	$(BUILDCMD) -o t3ll main.go

$(GOPATH)/bin/t3ll: ./t3ll
	cp t3ll $(GOPATH)/bin/t3ll

$(GOPATH)/bin/go-bindata:
	go get github.com/jteeuwen/go-bindata/...

$(GOPKGS):
	dep ensure
