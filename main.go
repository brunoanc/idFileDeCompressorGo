/*
* This file is part of idFileDeCompressorGo (https://github.com/PowerBall253/idFileDeCompressorGo).
* Copyright (C) 2023 PowerBall253
*
* idFileDeCompressorGo is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* idFileDeCompressorGo is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with idFileDeCompressorGo. If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

// #cgo linux LDFLAGS: -L./ooz -looz_linux -lstdc++ -static
// #cgo windows LDFLAGS: -L./ooz -looz_windows -lstdc++ -static
// #include "ooz/ooz.h"
import "C"

// Decompress the .entities file and write it
func decompressEntities(compressedEntities []byte) []byte {
	// Read uncompressed size
	uncompressedSize := binary.LittleEndian.Uint64(compressedEntities[:8])
	compressedSize := binary.LittleEndian.Uint64(compressedEntities[8:16])

	// Decompress the data using Kraken
	uncompressedData := make([]byte, uncompressedSize+64)
	result := int(C.Kraken_Decompress((*C.uchar)(&compressedEntities[16]), C.size_t(compressedSize), (*C.uchar)(&uncompressedData[0]), C.size_t(uncompressedSize)))
	if uint64(result) != uncompressedSize {
		return nil
	}

	return uncompressedData[:result]
}

func compressEntities(uncompressedEntities []byte) []byte {
	// Compress the data using Kraken
	compressedEntities := make([]byte, 16+(len(uncompressedEntities)+274*((len(uncompressedEntities)+0x3FFFF)/0x40000)))
	result := C.Kraken_Compress((*C.uchar)(&uncompressedEntities[0]), C.size_t(len(uncompressedEntities)), (*C.uchar)(&compressedEntities[16]), C.int(4))
	if result <= 0 {
		return nil
	}

	// Copy compressed and uncompressed sizes to beginning of file
	binary.LittleEndian.PutUint64(compressedEntities[:8], uint64(len(uncompressedEntities)))
	binary.LittleEndian.PutUint64(compressedEntities[8:16], uint64(result))

	return compressedEntities[:(result + 16)]
}

func printHelp() {
	fmt.Print("Usage: idFileDeCompressor [options] <src> <dest>\n\n")
	fmt.Print("Options:\n")
	fmt.Print("\t-d, --decompress\t\tDecompress a compressed .entities file.\n")
	fmt.Print("\t-c, --compress\t\tCompress an uncompressed .entities file.\n\n")
	fmt.Print("Example: idFileDeCompressor D:\\e1m1.entities D:\\e1m1.dec\n\n")
	fmt.Print("If no option is provided, the tool will attempt to auto-detect the action to perform.\n")
	fmt.Println("If no destination path is provided, the tool will use the source path with an added extension.")
}

// Main function
func main() {
	fmt.Printf("idFileDeCompressor v1.0 by PowerBall253 :)\n\n")

	// Get options from arguments
	var inputFile string
	var outputFile string
	var compress *bool

	for _, v := range os.Args[1:] {
		switch v {
		case "-d", "--decompress":
			b := false
			compress = &b
		case "-c", "--compress":
			b := true
			compress = &b
		case "-h", "--help":
			printHelp()
			os.Exit(0)
		default:
			if inputFile == "" {
				inputFile = v
			} else if outputFile == "" {
				outputFile = v
			} else {
				printHelp()
				os.Exit(1)
			}
		}
	}

	if inputFile == "" {
		printHelp()
		os.Exit(1)
	}

	// Read input file
	fileBytes, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to read from %s: %s\n", inputFile, err.Error())
		os.Exit(1)
	}

	// If compression or decompression wasn't explicitly set, check file
	if compress == nil {
		// Check for uncompressed start string
		startString := []byte("Version 7")
		for i, v := range startString {
			if fileBytes[i] != v {
				b := false
				compress = &b
			}
		}
		if compress == nil {
			b := true
			compress = &b
		}
	}

	// If no output path was set, create one
	if outputFile == "" {
		if *compress {
			outputFile = inputFile + ".entities"
		} else {
			outputFile = inputFile + ".dec"
		}
	}

	// Compress or decompress
	var outData []byte
	if *compress {
		outData = compressEntities(fileBytes)
		if outData == nil {
			fmt.Fprintf(os.Stderr, "ERROR: Couldn't compress %s.\n", inputFile)
			os.Exit(1)
		}
	} else {
		outData = decompressEntities(fileBytes)
		if outData == nil {
			fmt.Fprintf(os.Stderr, "ERROR: Couldn't decompress %s, bad file?\n", inputFile)
			os.Exit(1)
		}
	}

	// Write output file
	if err := os.WriteFile(outputFile, outData, 0666); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to write to %s: %s\n", outputFile, err.Error())
		os.Exit(1)
	}

	if *compress {
		fmt.Printf("Successfully compressed %s into %s.\n", inputFile, outputFile)
	} else {
		fmt.Printf("Successfully decompressed %s into %s.\n", inputFile, outputFile)
	}
}
