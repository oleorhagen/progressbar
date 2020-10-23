// Copyright 2017 Northern.tech AS
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
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func writeZeros(out io.Writer, cnt int64) {
	for i := int64(0); i < cnt; i++ {
		out.Write([]byte{0})
	}
}

func TestProgress(t *testing.T) {
	b := &bytes.Buffer{}
	p := &ProgressWriter{
		Out: b,
		Size:   100,
		ProgressMarker: ".",
	}

	n, err := p.Write([]byte{})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)

	writeZeros(p, 100)
	assert.Equal(t, "", b.String())


}
