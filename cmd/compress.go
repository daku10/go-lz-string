package cmd

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"

	lzstring "github.com/daku10/go-lz-string"
	"github.com/spf13/cobra"
)

var flagMethodEnum methodEnum = methodInvalidUTF16

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:     "compress file",
	Short:   "compress input using lz-string",
	Long:    `Compress input using lz-string(compatible 1.4.4). If no file is specified, input is from standard input. Input format must be UTF-8 string. Output format depends on the compression method.`,
	Example: "  go-lz-string compress input.txt -m base64 -o output.txt",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var inputFilename string
		if len(args) > 0 {
			inputFilename = args[0]
		}
		outputFilename, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		var input string
		var reader io.Reader
		if inputFilename == "" {
			reader = os.Stdin
		} else {
			f, err := os.Open(inputFilename)
			if err != nil {
				return err
			}
			reader = f
			defer f.Close()
		}
		bytes, err := io.ReadAll(reader)
		if err != nil {
			return err
		}
		input = string(bytes)

		switch flagMethodEnum {
		case methodInvalidUTF16:
			result, err := lzstring.Compress(input)
			if err != nil {
				return err
			}
			var buf *bufio.Writer
			if outputFilename != "" {
				outputF, err := os.Create(outputFilename)
				if err != nil {
					return err
				}
				defer outputF.Close()
				buf = bufio.NewWriter(outputF)
			} else {
				buf = bufio.NewWriter(os.Stdout)
			}
			err = binary.Write(buf, binary.LittleEndian, result)
			if err != nil {
				return err
			}
			err = buf.Flush()
			if err != nil {
				return err
			}
			return nil
		case methodBase64:
			result, err := lzstring.CompressToBase64(input)
			if err != nil {
				return err
			}
			var buf *bufio.Writer
			if outputFilename != "" {
				outputF, err := os.Create(outputFilename)
				if err != nil {
					return err
				}
				defer outputF.Close()
				buf = bufio.NewWriter(outputF)
			} else {
				buf = bufio.NewWriter(os.Stdout)
			}
			_, err = buf.WriteString(result)
			if err != nil {
				return err
			}
			err = buf.Flush()
			if err != nil {
				return err
			}
			return nil
		case methodUTF16:
			result, err := lzstring.CompressToUTF16(input)
			if err != nil {
				return err
			}
			var buf *bufio.Writer
			if outputFilename != "" {
				outputF, err := os.Create(outputFilename)
				if err != nil {
					return err
				}
				defer outputF.Close()
				buf = bufio.NewWriter(outputF)
			} else {
				buf = bufio.NewWriter(os.Stdout)
			}
			err = binary.Write(buf, binary.LittleEndian, result)
			if err != nil {
				return err
			}
			err = buf.Flush()
			if err != nil {
				return err
			}
			return nil
		case methodUint8Array:
			result, err := lzstring.CompressToUint8Array(input)
			if err != nil {
				return err
			}
			var buf *bufio.Writer
			if outputFilename != "" {
				outputF, err := os.Create(outputFilename)
				if err != nil {
					return err
				}
				defer outputF.Close()
				buf = bufio.NewWriter(outputF)
			} else {
				buf = bufio.NewWriter(os.Stdout)
			}
			_, err = buf.Write(result)
			if err != nil {
				return err
			}
			err = buf.Flush()
			if err != nil {
				return err
			}
			return nil
		case methodEncodedURIComponent:
			result, err := lzstring.CompressToEncodedURIComponent(input)
			if err != nil {
				return err
			}
			var buf *bufio.Writer
			if outputFilename != "" {
				outputF, err := os.Create(outputFilename)
				if err != nil {
					return err
				}
				defer outputF.Close()
				buf = bufio.NewWriter(outputF)
			} else {
				buf = bufio.NewWriter(os.Stdout)
			}
			_, err = buf.WriteString(result)
			if err != nil {
				return err
			}
			err = buf.Flush()
			if err != nil {
				return err
			}
			return nil
		}
		return errors.New("invalid method is specified")
	},
}

func init() {
	rootCmd.AddCommand(compressCmd)

	compressCmd.Flags().StringP("output", "o", "", "Print the output to the output file instead of the standard output.")
	compressCmd.Flags().VarP(&flagMethodEnum, "method", "m", `Compression method.
invalid-utf16: invalid UTF-16(output format is UTF-16 Little Endian No BOM. Sometimes it contains invalid UTF-16 code unit)
base64: base64(output format is UTF-8)
utf16: valid UTF-16(output format is UTF-16 Little Endian No BOM)
uint8array: uint8 array(output format is []byte)
encodedURIComponent: URL safe strings like base64(output format is UTF-8)
`)
}
