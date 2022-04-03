package cmd

import "io"

type Config struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}
