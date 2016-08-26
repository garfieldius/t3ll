BUILDCMD:=go build -a -tags netgo -ldflags "-w -s -extldflags '-static -s -pie'"

build: ./t3ll

release: release/t3ll_linux_amd64.tar.gz \
         release/t3ll_linux_386.tar.gz \
         release/t3ll_macos_amd64.zip \
         release/t3ll_macos_386.zip \
         release/t3ll_windows_amd64.zip \
         release/t3ll_windows_386.zip

clean:
	rm -rf release
	rm -f t3ll t3ll.exe
	cd frontend; rm -rf _dev _live


godeps: $(GOPATH)/bin/go-bindata \
        ./vendor/github.com/hydrogen18/stoppableListener/listener.go \
        ./vendor/github.com/kr/pretty/formatter.go

install: $(GOPATH)/bin/t3ll

.PHONY: default build release deps clean godeps install

release/t3ll_linux_amd64.tar.gz: release/linux_amd64/t3ll
	cd  release/linux_amd64 && tar -cf ../t3ll_linux_amd64.tar.gz t3ll

release/linux_amd64/t3ll: godeps server/data.go
	mkdir -p release/linux_amd64 && GOOS="linux" GOARCH="amd64" $(BUILDCMD) -o release/linux_amd64/t3ll

release/t3ll_linux_386.tar.gz: release/linux_386/t3ll
	cd  release/linux_amd64 && tar -cf ../t3ll_linux_x64.tar.gz t3ll

release/linux_386/t3ll: godeps server/data.go
	mkdir -p release/linux_386 && GOOS="linux" GOARCH="386" $(BUILDCMD) -o release/linux_386/t3ll

release/t3ll_macos_amd64.zip: release/macos_amd64/t3ll
	cd  release/macos_amd64 && zip -q9 ../t3ll_macos_amd64.zip t3ll

release/macos_amd64/t3ll: godeps server/data.go
	mkdir -p release/macos_amd64 && GOOS="darwin" GOARCH="amd64" $(BUILDCMD) -o release/macos_amd64/t3ll

release/t3ll_macos_386.zip: release/macos_386/t3ll
	cd  release/macos_386 && zip -q9 ../t3ll_macos_386.zip t3ll

release/macos_386/t3ll: godeps server/data.go
	mkdir -p release/macos_386 && GOOS="darwin" GOARCH="386" $(BUILDCMD) -o release/macos_386/t3ll

release/t3ll_windows_amd64.zip: release/windows_amd64/t3ll.exe
	cd  release/windows_amd64 && zip -q9 ../t3ll_windows_amd64.zip t3ll.exe

release/windows_amd64/t3ll.exe: godeps server/data.go
	mkdir -p release/windows_amd64 && GOOS="windows" GOARCH="amd64" $(BUILDCMD) -o release/windows_amd64/t3ll.exe

release/t3ll_windows_386.zip: release/windows_386/t3ll.exe
	cd  release/windows_386 && zip -q9 ../t3ll_windows_386.zip t3ll.exe

release/windows_386/t3ll.exe: godeps server/data.go
	mkdir -p release/windows_386 && GOOS="windows" GOARCH="386" $(BUILDCMD) -o release/windows_386/t3ll.exe

server/data.go: frontend/_live/editor.html
	go-bindata -nometadata -pkg server -o server/data.go \
		-prefix frontend/_live/ frontend/_live/editor.html

frontend/_live/editor.html: frontend/node_modules/.bin/gulp
	cd frontend; gulp clean; gulp live

frontend/node_modules/.bin/gulp:
	cd frontend; npm install

./t3ll: godeps server/data.go
	$(BUILDCMD) -o t3ll main.go

$(GOPATH)/bin/t3ll: ./t3ll
	cp t3ll $(GOPATH)/bin/t3ll

$(GOPATH)/bin/go-bindata:
	go get github.com/jteeuwen/go-bindata/...

./vendor/github.com/hydrogen18/stoppableListener/listener.go:
	glide install

./vendor/github.com/kr/pretty/formatter.go:
	glide install
