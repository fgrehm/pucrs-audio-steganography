package main

import (
	"fmt"
	"log"
	"os"

	"github.com/youpy/go-wav"
)

func encode(inputPath, outputPath string, lsbBitsToUse int, data []byte) error {
	log.Println("Encoding", inputPath, "to", outputPath)

	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	samples, format := readSamples(inputFile)
	inputFile.Close()

	if err := encodeData(samples, data, lsbBitsToUse); err != nil {
		return err
	}

	outputFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	fmt.Println("Writing new file")
	writer := wav.NewWriter(outputFile, uint32(len(samples)), format.NumChannels, format.SampleRate, format.BitsPerSample)
	return writer.WriteSamples(samples)
}

func encodeData(samples []wav.Sample, data []byte, lsbBitsToUse int) error {
	encodeBits(samples[0:4], byte(len(data)))

	for i, bits := range data {
		log.Printf("[%d] Encoding %s / %08b\n", i, string(bits), bits)

		// Each byte takes up 4 samples and we skip the first 4 because that's where
		// we keep the length of the payload
		base := (i + 1) * 4
		if err := encodeBits(samples[base:base+4], bits); err != nil {
			return err
		}
	}
	return nil
}

func encodeBits(samples []wav.Sample, bits byte) error {
	base := 0
	for j := uint(0); j < 8; j++ {
		channel := j % 2

		bitIsSet := (bits & (1 << j)) != 0
		if bitIsSet {
			samples[base].Values[channel] |= 1
		} else {
			samples[base].Values[channel] &= ^1
		}

		if channel == 1 {
			base += 1
		}
	}
	return nil
}
