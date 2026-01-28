package cli

import (
	"bufio"
	"fmt"

	"github.com/daku10/go-lz-string/version"
	"github.com/spf13/cobra"
)

func newVersionCmd(config *Config) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print version",
		Long:  "print the version of go-lz-string",
		RunE: func(cmd *cobra.Command, args []string) error {
			bufWriter := bufio.NewWriter(cmd.OutOrStdout())
			_, err := fmt.Fprintln(bufWriter, version.Version)
			if err != nil {
				return fmt.Errorf("version command failed: %w", err)
			}
			if err := bufWriter.Flush(); err != nil {
				return fmt.Errorf("version command failed: %w", err)
			}
			return nil
		},
	}
	versionCmd.SetIn(config.In)
	versionCmd.SetOut(config.Out)
	versionCmd.SetErr(config.Err)
	return versionCmd
}
