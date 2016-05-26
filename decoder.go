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

	samples, _ := readSamples(inputFile)
	return decodeData(samples, lsbsToUse), nil
}

func decodeData(samples []wav.Sample, lsbsToUse int) []byte {
	// 8 bits split across 2 channels with lsbsToUse LSB bits modified
	samplesPerByte := 8 / 2 / lsbsToUse

	count := 0
	for i := 3; i >= 0; i-- {
		count = count << 8

		base := i * samplesPerByte
		oneByte := decodeByte(samples[base:base+samplesPerByte], lsbsToUse)
		count += int(oneByte)
	}

	log.Println(count, "bytes to read")

	data := []byte{}
	for i := 0; i < count; i++ {
		// Each byte takes up 4 samples and we skip the first 4 because that's where
		// we keep the length of the payload
		base := (i + 4) * samplesPerByte

		oneByte := decodeByte(samples[base:base+samplesPerByte], lsbsToUse)
		data = append(data, oneByte)
	}

	return data
}

func decodeByte(samples []wav.Sample, lsbsToUse int) byte {
	oneByte := byte(0)
	for i := len(samples) - 1; i >= 0; i-- {
		sample := samples[i]
		for channel := 1; channel >= 0; channel-- {
			value := sample.Values[channel]
			for k := lsbsToUse; k > 0; k-- {
				oneByte = oneByte << 1
				oneByte += byte(value & 1)
				value = value >> 1
			}
		}
	}
	return oneByte
}
