package main

import (
	"io"
	"log"
	"os"

	"github.com/youpy/go-wav"
)

func readSamples(inputFile *os.File) ([]wav.Sample, *wav.WavFormat) {
	log.Println("Reading samples")

	reader := wav.NewReader(inputFile)

	format, err := reader.Format()
	if err != nil {
		panic(err)
	}
	log.Printf("%+v\n", format)

	samplesRead := []wav.Sample{}
	for {
		samples, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}
		samplesRead = append(samplesRead, samples...)
	}

	log.Printf("%d samples read", len(samplesRead))

	return samplesRead, format
}
