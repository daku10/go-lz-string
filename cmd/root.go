package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-lz-string",
	Short: "compress/decompress using lz-string",
	Long: `go-lz-string is a CLI application to compress/decompress using lz-string algorithm[https://github.com/pieroxy/lz-string].
This application implements algorithm is compatible with lz-string@1.4.4`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
