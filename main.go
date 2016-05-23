package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "steganography",
}

var encodeCmd = &cobra.Command{
	Use: "encode [input file] [output file] [text to write]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return fmt.Errorf("Invalid arguments provided")
		}
		if len(args[2]) > 255 {
			return fmt.Errorf("Text is too big")
		}
		return encode(args[0], args[1], lsbBitsToUse, []byte(args[2]))
	},
}

var decodeCmd = &cobra.Command{
	Use: "decode [input file]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("No file name provided")
		}

		return decode(args[0], lsbBitsToUse)
	},
}

var lsbBitsToUse = 1

func main() {
	rootCmd.PersistentFlags().IntVar(&lsbBitsToUse, "lsb-bits", lsbBitsToUse, "the amount of least significant bits to use")
	rootCmd.AddCommand(encodeCmd)
	rootCmd.AddCommand(decodeCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
