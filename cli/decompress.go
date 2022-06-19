package cli

import (
	"bufio"
	"errors"
	"io"
	"os"

	lzstring "github.com/daku10/go-lz-string"
	"github.com/spf13/cobra"
)

func newDecompressCmd(config *Config) *cobra.Command {
	var flagMethodEnum methodEnum = methodInvalidUTF16
	var decompressCmd = &cobra.Command{
		Use:     "decompress",
		Short:   "decompress input using lz-string",
		Long:    `Decompress input using lz-string(compatible 1.4.4). If no file is specified, input is rom standard input. Input format is depends on the compression method. Output format is UTF-8 string`,
		Example: "  go-lz-string decompress input.txt -m base64 -o output.txt",
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
			var reader io.Reader
			if inputFilename == "" {
				reader = cmd.InOrStdin()
			} else {
				f, err := os.Open(inputFilename)
				if err != nil {
					return err
				}
				reader = f
				defer f.Close()
			}
			var buf *bufio.Writer
			if outputFilename == "" {
				buf = bufio.NewWriter(cmd.OutOrStdout())
			} else {
				outputF, err := os.Create(outputFilename)
				if err != nil {
					return err
				}
				buf = bufio.NewWriter(outputF)
				defer outputF.Close()
			}
			defer buf.Flush()
			return doDecompress(reader, buf, flagMethodEnum)
		},
	}
	decompressCmd.SetIn(config.In)
	decompressCmd.SetOut(config.Out)
	decompressCmd.SetErr(config.Err)
	decompressCmd.Flags().StringP("output", "o", "", "Print the output to the output file instead of the standard output.")
	decompressCmd.Flags().VarP(&flagMethodEnum, "method", "m", `Compression method.
invalid-utf16: invalid UTF-16(input format must be UTF-16 Little Endian No BOM. Sometimes it contains invalid UTF-16 code unit)
base64: base64(input format must be UTF-8)
utf16: valid UTF-16(input format must be UTF-16 Little Endian No BOM)
uint8array: uint8 array(input format must be []byte)
encodedURIComponent: URL safe strings like base64(input format must be UTF-8)
`)
	return decompressCmd
}

func readAsUint16Array(reader io.Reader) ([]uint16, error) {
	result := make([]uint16, 0)
	tmp := make([]byte, 2)
	for {
		n, err := reader.Read(tmp)
		if err == io.EOF {
			break
		}
		if n != 2 {
			return nil, errors.New("Not Uint16 Array")
		}
		result = append(result, uint16(tmp[0])|uint16(tmp[1])<<8)
	}
	return result, nil
}

func doDecompress(reader io.Reader, writer io.Writer, method methodEnum) error {
	var result string
	switch method {
	case methodInvalidUTF16:
		input, err := readAsUint16Array(reader)
		if err != nil {
			return err
		}
		result, err = lzstring.Decompress(input)
		if err != nil {
			return err
		}
	case methodUTF16:
		input, err := readAsUint16Array(reader)
		if err != nil {
			return err
		}
		result, err = lzstring.DecompressFromUTF16(input)
		if err != nil {
			return err
		}
	case methodBase64:
		input, err := io.ReadAll(reader)
		if err != nil {
			return err
		}
		inputString := string(input)
		result, err = lzstring.DecompressFromBase64(inputString)
		if err != nil {
			return err
		}
	case methodEncodedURIComponent:
		input, err := io.ReadAll(reader)
		if err != nil {
			return err
		}
		inputString := string(input)
		result, err = lzstring.DecompressFromEncodedURIComponent(inputString)
		if err != nil {
			return err
		}
	case methodUint8Array:
		input, err := io.ReadAll(reader)
		if err != nil {
			return err
		}
		result, err = lzstring.DecompressFromUint8Array(input)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid method is specified")
	}
	_, err := writer.Write([]byte(result))
	if err != nil {
		return err
	}
	return nil
}
