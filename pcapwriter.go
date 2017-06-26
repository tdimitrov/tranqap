package main

import (
	"fmt"
	"os"
)

type pcapwirter struct {
	fd *os.File
}

func (pw *pcapwirter) Init(fname string) (err error) {
	fd, err := os.OpenFile(fname, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("Error opening file: ", err)
	}
	pw.fd = fd
	return err
}

func (pw pcapwirter) Write(p []byte) (n int, err error) {
	n, err = pw.fd.Write(p)
	if err != nil {
		fmt.Println("Error writing to file: ", err)
	}
	return n, err
}

func (pw *pcapwirter) DeInit() {
	pw.fd.Close()
}
