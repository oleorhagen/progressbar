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
	"os"
)

type ProgressWriter struct {
	W io.Writer // wrapped writer
	*Renderer
}

type ProgressReader struct {
	R io.Reader // wrapped reader
	*Renderer
}

func NewWriter(size uint64) *ProgressWriter {
	return &ProgressWriter{
		Renderer: NewRenderer(size),
	}
}

func NewReader(size uint64) *ProgressReader {
	return &ProgressReader{
		Renderer: NewRenderer(size),
	}
}

func (p *ProgressWriter) Wrap(w io.Writer, size uint64) {
	p.W = w
	p.Renderer.Size = size
}

func (p *ProgressReader) Wrap(r io.Reader, size uint64) {
	p.R = r
	p.Renderer.Size = size
}

func (p *ProgressWriter) Write(data []byte) (int, error) {
	n := len(data)
	p.Renderer.Tick(uint64(n))
	return n, nil
}

type Renderer struct {
	Out            io.Writer // output device
	Size           uint64    // size of the input
	ProgressMarker string
	currentCount   uint64 // current count
}

func NewRenderer(size uint64) *Renderer {
	return &Renderer{
		Out:            os.Stderr,
		Size:           size,
		ProgressMarker: ".",
	}
}

func (p *Renderer) Tick(n uint64) {
	p.currentCount += uint64(n)
	p.render()
}


func (p *Renderer) render() {
	if p.Out == nil {
		return
	}
	percent := int((float64(p.currentCount) / float64(p.Size)) * 100)

	str := fmt.Sprintf("\r%s%s - %d ",
		strings.Repeat(p.ProgressMarker, percent),
		strings.Repeat(" ", 100-percent),
		percent)
	fmt.Fprintf(p.Out, str)
}

func (p *Renderer) renderNoTTY() {
	if p.Out == nil {
		return
	}
	// Get the percentage written
	// percentNew := int((float64(p.currentCount) / float64(p.Size)) * 100)
	// if percentNew > p.lastPercent {
	// 	str := strings.Repeat(p.ProgressMarker, percentNew-p.lastPercent)
	// 	p.Out.Write([]byte(str))
	// 	p.lastPercent = percentNew
	// }
}

func (p *Renderer) Finish() {
	if p.Out == nil {
		return
	}
	str := strings.Repeat(p.ProgressMarker, 100)
	fmt.Fprint(p.Out, "\r"+str+" - 100 %")
}
