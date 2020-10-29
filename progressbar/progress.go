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

// TODO -- Add terminal width respect

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
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

type Render interface {
	render(uint64, uint64, string) string
}

type TTYWriter struct {
}

func (p *TTYWriter) render(currentCount uint64, size uint64, marker string) string {
	percent := int((float64(currentCount) / float64(size)) * 100)

	str := fmt.Sprintf("\r%s%s - %d ",
		strings.Repeat(marker, percent),
		strings.Repeat(" ", 100-percent),
		percent)
	return str
}

type NoTTYWriter struct {
	lastPercent int
}

func (p *NoTTYWriter) render(currentCount uint64, size uint64, marker string) string {
	// Get the percentage written
	// percentNew := int((float64(p.currentCount) / float64(p.Size)) * 100)
	// if percentNew > p.lastPercent {
	// 	str := strings.Repeat(p.ProgressMarker, percentNew-p.lastPercent)
	// 	p.Out.Write([]byte(str))
	// 	p.lastPercent = percentNew
	// }
	return ""
}

type Renderer struct {
	Out            io.Writer // output device
	Type           Render
	Size           uint64 // size of the input
	ProgressMarker string
	currentCount   uint64 // current count
}

func NewRenderer(size uint64) *Renderer {
	var rt Render
	if isatty.IsTerminal(os.Stdout.Fd()) {
		rt = &TTYWriter{}
	} else {
		rt = &NoTTYWriter{}
	}
	return &Renderer{
		Out:            os.Stderr,
		Size:           size,
		ProgressMarker: ".",
		Type:           rt,
	}
}

func (p *Renderer) Tick(n uint64) {
	p.currentCount += uint64(n)
	str := p.Type.render(p.currentCount, p.Size, p.ProgressMarker)
	fmt.Fprintf(p.Out, str)
}

func (p *Renderer) Finish() {
	if p.Out == nil {
		return
	}
	str := strings.Repeat(p.ProgressMarker, 100)
	fmt.Fprint(p.Out, "\r"+str+" - 100 %")
}
