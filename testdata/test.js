#!/usr/bin/env node
const lzString = require("../third-party/lz-string/libs/lz-string.js");
const fs = require("fs");

if (process.argv.length < 3) {
  console.error("invalid input");
  process.exit(1);
}

const method = process.argv[2];

if (method === "invalid-utf16" && process.argv.length != 4) {
  console.error("output filepath must be specified");
  process.exit(1);
}

let func;

switch (method) {
  case "invalid-utf16":
    func = lzString.compress;
    break;
  case "utf16":
    func = lzString.compressToUTF16;
    break;
  case "base64":
    func = lzString.compressToBase64;
    break;
  case "uint8array":
    func = lzString.compressToUint8Array;
    break;
  case "encodedURIComponent":
    func = lzString.compressToEncodedURIComponent;
    break;
  default:
    console.error("invalid method is specified");
    process.exit(1);
}

const input = fs.readFileSync("/dev/stdin", "utf8");

if (method === "invalid-utf16") {
  fs.writeFileSync(process.argv[3], func(input), {
    encoding: "utf16le",
  });
} else {
  process.stdout.write(func(input));
}
