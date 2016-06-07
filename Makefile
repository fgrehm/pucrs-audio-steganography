GOPATH = "$(HOME)/gocode"

default: build

build: steganography-wav

serve:
	GOPATH=$(GOPATH) $(GOPATH)/bin/gin run web

steganography-wav: *.go
	GOPATH=$(GOPATH) go build -o $(@)

crosscompile:
	@mkdir -p build/
	@GOPATH=$(GOPATH) go get github.com/inconshreveable/mousetrap
	GOPATH=$(GOPATH) GOOS=linux GOARCH=amd64 go build -o build/steganography-wav-linux_amd64 .
	GOPATH=$(GOPATH) GOOS=linux GOARCH=386 go build -o build/steganography-wav-linux_386 .
	GOPATH=$(GOPATH) GOOS=darwin GOARCH=amd64 go build -o build/steganography-wav-darwin_amd64 .
	GOPATH=$(GOPATH) GOOS=windows GOARCH=386 go build -o build/steganography-wav-windows_386.exe .
	GOPATH=$(GOPATH) GOOS=windows GOARCH=amd64 go build -o build/steganography-wav-windows_amd64.exe .
