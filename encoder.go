package main

import (
	"fmt"
	"log"
	"os"

	"github.com/youpy/go-wav"
)

func encode(inputPath, outputPath string, lsbsToUse int, filename string, data []byte) error {
	log.Println("Encoding", inputPath, "to", outputPath)

	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	samples, format := readSamples(inputFile)
	inputFile.Close()

	if lsbsToUse > int(format.BitsPerSample) {
		return fmt.Errorf("The file does not have enough bits to handle the amount of LSBs provided")
	}

	usableBytesCount := len(samples) * int(format.NumChannels) * lsbsToUse / 8
	log.Printf("%d samples read, %d kbytes available to write", len(samples), usableBytesCount/1024)

	headerLen := 4*2 + len(filename)
	if (len(data) - headerLen) > usableBytesCount {
		return fmt.Errorf("The input file is too small to handle the payload")
	}

	log.Println("Encoding", len(data), "bytes")
	encoder := &encoder{samples, format, lsbsToUse, encoderPtr{}}
	encoder.writeInt(len(filename))
	encoder.writeBytes([]byte(filename))
	encoder.writeInt(len(data))
	encoder.writeBytes(data)

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

type encoder struct {
	samples   []wav.Sample
	format    *wav.WavFormat
	lsbsToUse int
	ptr       encoderPtr
}

type encoderPtr struct {
	sample       int
	lsb, channel uint
}

func (e *encoder) writeInt(value int) {
	for i := 0; i < 4; i++ {
		e.writeByte(byte(value))
		value >>= 8
	}
}

func (e *encoder) writeBytes(bytes []byte) {
	for _, oneByte := range bytes {
		e.writeByte(oneByte)
	}
}

func (e *encoder) writeByte(oneByte byte) {
	for i := 0; i < 8; i++ {
		sample := e.samples[e.ptr.sample]
		value := sample.Values[e.ptr.channel]

		mask := byte(1 << uint(i))
		if oneByte&mask == mask {
			value |= 1 << e.ptr.lsb
		} else {
			value &= ^(1 << e.ptr.lsb)
		}

		e.samples[e.ptr.sample].Values[e.ptr.channel] = value

		e.ptr.lsb += 1
		if e.ptr.lsb < uint(e.lsbsToUse) {
			continue
		}

		e.ptr.lsb = 0
		e.ptr.channel += 1
		if e.ptr.channel == uint(e.format.NumChannels) {
			e.ptr.sample += 1
			e.ptr.channel = 0
		}
	}
}
