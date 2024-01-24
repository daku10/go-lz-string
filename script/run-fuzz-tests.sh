#!/bin/bash

set -eux

cd "$(dirname "$0")/.."

files=$(grep -r --include='**_test.go' --files-with-matches 'func Fuzz' .)

for file in ${files}
do
	funcs=$(grep -o 'func Fuzz\w*' "$file" | sed 's/func //')
	for func in ${funcs}
	do
		echo "Fuzzing $func in $file"
		go test -fuzz="^${func}\$" -fuzztime="10s"
	done
done
