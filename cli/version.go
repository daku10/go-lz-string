package cli

import (
	"bufio"
	"errors"
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
			fmt.Fprintln(bufWriter, version.Version)
			if err := bufWriter.Flush(); err != nil {
				return errors.New("version command failed")
			}
			return nil
		},
	}
	versionCmd.SetIn(config.In)
	versionCmd.SetOut(config.Out)
	versionCmd.SetErr(config.Err)
	return versionCmd
}
