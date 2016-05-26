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

	// 8 bits split across 2 channels with lsbsToUse LSB bits modified
	samplesPerByte := 8 / 2 / lsbsToUse

	usableBytesCount := len(samples)/samplesPerByte - (32 / samplesPerByte)
	log.Printf("%d samples read, %d kbytes available to write", len(samples), usableBytesCount/1024)

	if len(data) > usableBytesCount {
		return fmt.Errorf("The input file is too small")
	}

	log.Println("Encoding", len(data), "bytes")
	encodeData(samples, data, lsbsToUse)

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

func encodeData(samples []wav.Sample, data []byte, lsbsToUse int) {
	// 8 bits split across 2 channels with lsbsToUse LSB bits modified
	samplesPerByte := 8 / 2 / lsbsToUse

	// dataLength is a 32bit int that takes up 4 bytes
	dataLength := len(data)
	for i := 0; i < 4; i++ {
		base := i * samplesPerByte
		encodeByte(samples[base:base+samplesPerByte], byte(dataLength), lsbsToUse)
		dataLength = dataLength >> 8
	}

	for i, oneByte := range data {
		// Each byte takes up X samples and we skip the first 4 usable bytes because
		// that's where we keep the length of the payload
		base := (i + 4) * samplesPerByte
		encodeByte(samples[base:base+samplesPerByte], oneByte, lsbsToUse)
	}
}

func encodeByte(samples []wav.Sample, oneByte byte, lsbsToUse int) {
	for i, sample := range samples {
		for channel := 0; channel < 2; channel++ {
			value := sample.Values[channel]
			value = value >> uint(lsbsToUse)
			for k := 0; k < lsbsToUse; k++ {
				value = value << 1
				value += int(oneByte & 1)
				oneByte = oneByte >> 1
			}
			samples[i].Values[channel] = value
		}
	}
}
