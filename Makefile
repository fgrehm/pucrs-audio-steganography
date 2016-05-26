GOPATH = "$(HOME)/gocode"

default: build

build: steganography-wav

serve:
	GOPATH=$(GOPATH) $(GOPATH)/bin/gin run web

steganography-wav: *.go
	GOPATH=$(GOPATH) go build -o $(@)
