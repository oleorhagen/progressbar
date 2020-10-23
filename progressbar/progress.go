// Copyright 2019 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
package progressbar

import (
	"io"
	"strings"
	// "time"
	"fmt"
)

type ProgressWriter struct {
	Wrap           io.Writer // wrapped writer
	Out            io.Writer // output device
	Size           uint64    // size of the input
	ProgressMarker string
	currentCount   uint64 // current count
	lastPercent    int
	over           bool // set to true of writes have gone over declared N bytes
}

func New(w io.Writer, size uint64) *ProgressWriter {
	// is a tty (?)
	return &ProgressWriter{
		Wrap:           w,
		Size:           size,
		ProgressMarker: ".",
	}
}

func (p *ProgressWriter) Write(data []byte) (int, error) {
	n := len(data)

	p.currentCount += uint64(n)
	p.rendera()
	return n, nil
}

func (p *ProgressWriter) Tick(n uint64) {
	p.currentCount += uint64(n)
	p.render()
}

func (p *ProgressWriter) Finish() {
	if p.Out == nil {
		return
	}
	str := strings.Repeat(p.ProgressMarker, 100)
	fmt.Fprint(p.Out, "\r" + str + " - 100 %%")
}

// reportGeneric prints the progressbar to the screen
func (p *ProgressWriter) render() {
	if p.Out == nil {
		return
	}
	// Get the percentage written
	percentNew := int((float64(p.currentCount) / float64(p.Size)) * 100)
	if percentNew > p.lastPercent {
		str := strings.Repeat(p.ProgressMarker, percentNew - p.lastPercent)
		p.Out.Write([]byte(str))
		p.lastPercent = percentNew
	}
}

func (p *ProgressWriter) rendera() {
	if p.Out == nil {
		return
	}
	percent := int((float64(p.currentCount) / float64(p.Size)) * 100)

	str := fmt.Sprintf("\r%s%s - %d ",
		strings.Repeat(p.ProgressMarker, percent),
		strings.Repeat(" ", 100 - percent),
		percent)
	fmt.Fprintf(p.Out, str)
}
