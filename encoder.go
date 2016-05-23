package main

import (
	"fmt"
	"log"
	"os"

	"github.com/youpy/go-wav"
)

func encode(inputPath, outputPath string, lsbsToUse int, data []byte) error {
	log.Println("Encoding", inputPath, "to", outputPath)

	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	samples, format := readSamples(inputFile)
	inputFile.Close()

	if format.NumChannels != 2 {
		return fmt.Errorf("Mono audio files are not supported")
	}

	usableBytesCount := len(samples)/4-4
	log.Printf("%d samples read, %d kbytes available to write", len(samples), usableBytesCount/1024)

	log.Println("Encoding", len(data), "bytes")
	if err := encodeData(samples, data, lsbsToUse); err != nil {
		return err
	}

	outputFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	log.Println("Writing new file")
	writer := wav.NewWriter(outputFile, uint32(len(samples)), format.NumChannels, format.SampleRate, format.BitsPerSample)
	if err := writer.WriteSamples(samples); err != nil {
		return err
	}
	log.Println("Done")
	return nil
}

func encodeData(samples []wav.Sample, data []byte, lsbsToUse int) error {
	dataLength := len(data)
	for i := 0; i < 4; i++ {
		// Each byte takes up 4 samples
		base := i * 4
		encodeBits(samples[base:base+4], byte(dataLength))
		dataLength = dataLength >> 8
	}

	for i, bits := range data {
		// Each byte takes up 4 samples and we skip the first 4 sets because that's
		// where we keep the length of the payload
		base := (i + 4) * 4
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
