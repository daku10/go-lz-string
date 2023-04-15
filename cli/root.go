package cli

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
func newRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "go-lz-string",
		Short: "compress/decompress using lz-string",
		Long: `go-lz-string is a CLI application to compress/decompress using lz-string algorithm[https://github.com/pieroxy/lz-string].
This application implements algorithm is compatible with lz-string@^1.4.4`,
	}
	config := &Config{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	rootCmd.AddCommand(
		newCompressCmd(config),
		newDecompressCmd(config),
		newVersionCmd(config),
	)
	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := newRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
