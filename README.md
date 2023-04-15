# go-lz-string

### Go implementation of [lz-string](https://github.com/pieroxy/lz-string)
This is an implementation of lz-string written in Go. You can also embed go-lz-string as a library to your Go products.

## lz-string algorithm version
Currently it implements algorithm that is compatible with [lz-string@^1.4.4](https://github.com/pieroxy/lz-string/releases/tag/1.4.4)

## Usage as a CLI

### Installation

#### From GitHub Releases
See [GitHub Releases](https://github.com/daku10/go-lz-string/releases).  
Available binaries are:
- macOS
- Linux
- Windows

#### go install

Go v1.18 or later is required.

```sh
go install github.com/daku10/go-lz-string/cmd/go-lz-string@latest
```

### Example

```sh
$ go-lz-string compress <filename> -o <output-filename>
# use standard input/output
$ echo -n '🍎🍇🍌' | go-lz-string compress -m base64
jwbjl96cX3kGX2g=
$ echo -n 'jwbjl96cX3kGX2g=' | go-lz-string decompress -m base64
🍎🍇🍌
```
Compression and decompression methods are not only `base64` but also `invalid utf-16`, `utf-16`, `encodedURIComponent` and `byte array` are implemented as well as original javascript program. To use other methods, see `go-lz-string help`

## Usage as a library

### Installation
```sh
go get -u github.com/daku10/go-lz-string
```

### Example
```go
package main

import (
	"fmt"

	lzstring "github.com/daku10/go-lz-string"
)

func main() {
	var input string = "Hello, world"
	var compressed []uint16 = lzstring.Compress(input)
	// [1157 12342 24822 832 1038 59649 14720 9792]
	fmt.Println(compressed)
	var decompressed string = lzstring.Decompress(compressed)
	// Hello, world!
	fmt.Println(decompressed)
}
```

## Motivation

- Need CLI tool to use easily
- Some the go implementations of lz-string already exists, but there are lack of some functionality.
  - [pieroxy/lz-string-go](https://github.com/pieroxy/lz-string-go) support decompression encodedURIComponent only.
  - [lazarus/lz-string-go](https://github.com/lazarus/lz-string-go) support compression/decompresson, but base64 method only, and specific input like emoji can not be compressed correctly.

## Author
[daku10](https://github.com/daku10)

## License
This software is released under the MIT License, see LICENSE.

## Third-party Libraries

This repository contains third-party libraries in the "third-party" directory. Each library is stored as a Git submodule and is subject to the license terms of its respective library.

Please note that these third-party libraries are not developed, maintained, or supported by this repository's maintainers, and any issues or questions related to these libraries should be directed to their respective maintainers.

For more information on the licensing terms of each library, please see the LICENSE file located in each library's directory.
