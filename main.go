package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	inFile  string
	outFile string
)

func main() {
	flag.StringVar(&inFile, "inFile", "", "The input file to encode.")
	flag.StringVar(&outFile, "outFile", "out.h", "The source file output path.")
	flag.Parse()

	stat, err := os.Stat(inFile)
	if err != nil {
		panic(err)
	}

	in, err := os.Open(inFile)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	out, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	fileName  := filepath.Base(inFile)
	outBuffer := bufio.NewWriterSize(out, 64*1024)

	outBuffer.WriteString("#pragma once\n\n")
	outBuffer.WriteString(fmt.Sprintf(
		"static const unsigned char s_%s[] = {\n", 
		fileName[:strings.LastIndex(fileName, ".")],
	))

	inBuffer := make([]byte, 1)
	bytes    := int64(0)
	func() {
		for x := 0;; x++ {
			for i := 0; i < 8; i++ { // 8 items per line
				count, err := in.Read(inBuffer)
				if count == 0 {
					return
				}

				if err != nil {
					panic(err)
				}

				if i == 0 {
					if x > 0 {
						outBuffer.WriteString(",\n")
					}
					outBuffer.WriteString("    ")
				} else {
					outBuffer.WriteString(", ")
				}

				outBuffer.WriteString(fmt.Sprintf("0x%04x", inBuffer[:count]))
				bytes++
			}
		}
	}()

	if stat.Size() != bytes {
		panic(fmt.Errorf("Uh oh, sizes don't match! (%d != %d)", stat.Size(), bytes))
	}

	outBuffer.WriteString("\n};\n\n")
	outBuffer.WriteString(fmt.Sprintf(
		"static const unsigned s_%s_size = %d;\n", 
		fileName[:strings.LastIndex(fileName, ".")], 
		bytes,
	))
	outBuffer.Flush()
}
