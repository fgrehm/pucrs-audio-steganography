package main

import (
	"io/ioutil"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "steganography-wav",
}

var encodeCmd = &cobra.Command{
	Use: "encode [input file] [output file] [text to write]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return fmt.Errorf("Invalid arguments provided")
		}
		return encode(args[0], args[1], LSBsToUse, []byte(args[2]))
	},
}

var encodeBinCmd = &cobra.Command{
	Use: "encode-bin [input file] [output file] [payload]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return fmt.Errorf("Invalid arguments provided")
		}
		fileContents, err := ioutil.ReadFile(args[2])
		if err != nil {
			return err
		}
		return encode(args[0], args[1], LSBsToUse, fileContents)
	},
}

var decodeCmd = &cobra.Command{
	Use: "decode [input file]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("No file name provided")
		}

		str, err := decode(args[0], LSBsToUse)
		if err != nil {
			return err
		}

		fmt.Println("String found:", string(str))
		return nil
	},
}

var decodeBinCmd = &cobra.Command{
	Use: "decode-bin [input file] [output file]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("Invalid arguments provided")
		}

		data, err := decode(args[0], LSBsToUse)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(args[1], data, 0644); err != nil {
			return err
		}
		fmt.Println("Payload wrote to", args[1])
		return nil
	},
}

var LSBsToUse = 1

func main() {
	rootCmd.PersistentFlags().IntVar(&LSBsToUse, "lsb-bits", LSBsToUse, "the amount of least significant bits to use")
	rootCmd.AddCommand(encodeCmd)
	rootCmd.AddCommand(decodeCmd)
	rootCmd.AddCommand(encodeBinCmd)
	rootCmd.AddCommand(decodeBinCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
