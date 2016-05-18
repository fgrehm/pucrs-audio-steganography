GOPATH = "$(HOME)/gocode"

default: steganography-wav

steganography-wav: main.go
	GOPATH=$(GOPATH) go build -o $(@)

