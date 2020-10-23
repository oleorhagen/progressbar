package main

import (
	"io"
	"os"
	"time"

	"github.com/oleorhagen/progress/progressbar"
)

func main() {

	p := &progressbar.ProgressWriter{
		Out:            os.Stdout,
		Size:           100,
		ProgressMarker: ".",
	}

	f, _ := os.Open("/dev/random")
	for i := 0; i < 7; i++ {
		io.CopyN(p, f, 10)
		time.Sleep(1 * time.Second)
	}
	p.Finish()
}
