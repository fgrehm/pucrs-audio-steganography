GOPATH = "$(HOME)/gocode"

default: steganography-wav

steganography-wav: *.go
	GOPATH=$(GOPATH) go build -o $(@)

