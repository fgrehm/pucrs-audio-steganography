package main

import (
	"log"
	"os"

	"github.com/youpy/go-wav"
)

func decode(inputPath string, lsbBitsToUse int) error {
	log.Println("Decoding")

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	samples, _ := readSamples(inputFile)
	inputFile.Close()

	countInBytes, err := decodeBits(samples[0:4])
	if err != nil {
		return err
	}
	count := int(countInBytes)

	log.Println(count, "bytes to read")

	data := []byte{}
	for i := 0; i < count; i++ {
		// Each byte takes up 4 samples and we skip the first 4 because that's where
		// we keep the length of the payload
		base := (i+1) * 4

		bits, err := decodeBits(samples[base:base+4])
		if err != nil {
			return err
		}
		data = append(data, bits)
	}

	log.Println("Data read:", string(data))

	return nil
}

func decodeBits(samples []wav.Sample) (byte, error) {
	base := 3
	bits := byte(0)
	for j := 7; j >= 0; j-- {
		channel := j % 2

		bits += uint8(samples[base].Values[channel] & 1)
		if j != 0 {
			bits = bits << 1
		}

		if channel == 0 {
			base -= 1
		}
	}
	return bits, nil
}
