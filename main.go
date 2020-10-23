package main

import (
	"io"
	"os"
	"time"

	"github.com/oleorhagen/progress/progressbar"
)

func main() {

	p := progressbar.NewWriter(100)
	f, _ := os.Open("/dev/random")
	for i := 0; i < 7; i++ {
		io.CopyN(p, f, 10)
		p.Tick(1)
		time.Sleep(1 * time.Second)
	}
	p.Finish()
}
