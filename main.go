package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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
		return encode(args[0], args[1], LSBsToUse, "__string__", []byte(args[2]))
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
		filename := filepath.Base(args[2])
		return encode(args[0], args[1], LSBsToUse, filename, fileContents)
	},
}

var decodeCmd = &cobra.Command{
	Use: "decode [input file]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("No file name provided")
		}

		filename, data, err := decode(args[0], LSBsToUse)
		if err != nil {
			return err
		}

		if filename == "__string__" {
			fmt.Println("String found:", string(data))
		} else {
			if err := ioutil.WriteFile(filename, data, 0644); err != nil {
				return err
			}
			fmt.Println("Payload wrote to", filename)
		}
		return nil
	},
}

var webCmd = &cobra.Command{
	Use: "web [port number]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("Invalid arguments provided")
		}

		port := "8080"
		if os.Getenv("PORT") != "" {
			port = os.Getenv("PORT")
		}
		if len(args) == 1 && args[0] != "" {
			port = args[0]
		}
		runServer(port)
		return nil
	},
}

var LSBsToUse = 1

func main() {
	rootCmd.PersistentFlags().IntVar(&LSBsToUse, "lsb", LSBsToUse, "the amount of least significant bits to use")
	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(encodeCmd)
	rootCmd.AddCommand(encodeBinCmd)
	rootCmd.AddCommand(decodeCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
