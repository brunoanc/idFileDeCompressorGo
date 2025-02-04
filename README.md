# idFileDeCompressor
![Build Status](https://github.com/brunoanc/idFileDeCompressorGo/actions/workflows/test.yml/badge.svg)

Tool to decompress and recompress `.entities` files, allowing level editing.

## Usage

```
idFileDeCompressor [options] <src> <dest>
```

The tool will attempt to auto-detect the action to perform. You can override this behaviour with the `--decompress` and `--compress` flags.

If no destination path is provided, it will overwrite the source file.

## Compiling
The project requires the [go toolchain](https://go.dev/dl/) to be compiled. Additionally, a GCC toolchain such as MinGW is required on Windows.

To compile, run:

```
go build -o idFileDeCompressor -ldflags="-s -w" .
```

Additionally, you may use [UPX](https://upx.github.io/) to compress the binary:

```
upx --best idFileDeCompressor
```

## Credits
* proteh: For creating the original idFileDeCompressor.
