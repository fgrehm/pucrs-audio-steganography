package main

import (
	"log"
	"os"

	"github.com/youpy/go-wav"
)

func decode(inputPath string, lsbsToUse int) ([]byte, error) {
	log.Println("Decoding")

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	samples, format := readSamples(inputFile)

	decoder := &decoder{samples, format, lsbsToUse, decoderPtr{}}
	dataLength := decoder.readInt()
	log.Println("Payload size", dataLength, "bytes")
	return decoder.readBytes(dataLength), nil
}

type decoder struct {
	samples   []wav.Sample
	format    *wav.WavFormat
	lsbsToUse int
	ptr       decoderPtr
}

type decoderPtr struct {
	sample       int
	lsb, channel uint
}

func (d *decoder) readInt() int {
	ret := 0
	for i := 0; i < 4; i++ {
		ret += int(d.readByte()) << uint(8*i)
	}
	return ret
}

func (d *decoder) readBytes(count int) []byte {
	ret := []byte{}
	for i := 0; i < count; i++ {
		ret = append(ret, d.readByte())
	}
	return ret
}

func (d *decoder) readByte() byte {
	oneByte := byte(0)
	for i := 0; i < 8; i++ {
		sample := d.samples[d.ptr.sample]
		value := sample.Values[d.ptr.channel]

		mask := int(1 << d.ptr.lsb)
		if value&mask == mask {
			oneByte |= 1 << uint(i)
		}

		d.ptr.lsb += 1
		if d.ptr.lsb < uint(d.lsbsToUse) {
			continue
		}

		d.ptr.lsb = 0
		d.ptr.channel += 1
		if d.ptr.channel == uint(d.format.NumChannels) {
			d.ptr.sample += 1
			d.ptr.channel = 0
		}
	}
	return oneByte
}
