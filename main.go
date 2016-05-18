package main

import (
	"flag"
	"fmt"
	"github.com/youpy/go-wav"
	"io"
	"os"
)

func main() {
	infile := flag.String("infile", "", "wav file to read")
	outfile := flag.String("outfile", "", "wav file to write")
	flag.Parse()

	samples, format := readSamples(*infile)

	if *outfile != "" {
		fmt.Println("Encoding")

		outputFile, err := os.OpenFile(*outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		numbers := []uint8{}
		for i := 0; i < 256; i++ {
			numbers = append(numbers, uint8(i))
		}
		writeUint8(samples, numbers...)

		fmt.Println("Writing new file")
		writer := wav.NewWriter(outputFile, uint32(len(samples)), format.NumChannels, format.SampleRate, format.BitsPerSample)
		writer.WriteSamples(samples)
	} else {
		fmt.Println("Decoding")
		numbers := readUint8(samples, 256)
		for _, num := range numbers {
			println(num)
		}
	}
}

// TODO: WRITE A SET OF BYTES
func writeUint8(samples []wav.Sample, nums ...uint8) {
	for i, num := range nums {
		fmt.Printf("[%d] Encoding %05d / %08b\n", i, num, num)

		// Each number takes up 4 samples
		base := i * 4

		for j := uint(0); j < 8; j++ {
			channel := j % 2

			bitIsSet := (num & (1 << j)) != 0
			if bitIsSet {
				samples[base].Values[channel] |= 1
			} else {
				samples[base].Values[channel] &= ^1
			}

			if channel == 1 {
				base += 1
			}
		}
	}
}

// TODO: READ A SET OF BYTES
func readUint8(samples []wav.Sample, count int) []uint8 {
	nums := []uint8{}
	for i := 0; i < count; i++ {
		num := uint8(0)

		// Each number takes up 4 samples
		base := i*4 + 3

		for j := 7; j >= 0; j-- {
			channel := j % 2

			num += uint8(samples[base].Values[channel] & 1)
			if j != 0 {
				num = num << 1
			}

			if channel == 0 {
				base -= 1
			}
		}

		nums = append(nums, num)
	}
	return nums
}

func readSamples(inputPath string) ([]wav.Sample, *wav.WavFormat) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}

	defer inputFile.Close()

	reader := wav.NewReader(inputFile)

	format, err := reader.Format()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", format)

	samplesRead := []wav.Sample{}
	for {
		samples, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}
		samplesRead = append(samplesRead, samples...)
	}
	return samplesRead, format
}
