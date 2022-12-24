package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func WriteToPty(w io.Writer) {
	for {
		buf := make([]byte, 1)
		var stdin io.Reader = os.Stdin
		stdin.Read(buf)

		byteReader := bytes.NewReader(buf)

		io.Copy(w, byteReader)
	}
}

func ReadFromPty(r io.Reader) {
	for {
		buf := make([]byte, 1)
		_, err := io.ReadFull(r, buf)

		if err != nil {
			fmt.Printf("Error: %s\n\r", err.Error())
			break
		}

		fmt.Print(string(buf))
	}
}
