// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testEmit(generator Generator) string {
	emitter := NewEmitter(skel.FileFormatOther)
	generator.Generate(emitter)
	return emitter.String()
}

func TestSnippet_GetName(t *testing.T) {
	snippet := Snippet{Name: "snippet"}
	assert.Equal(t, "snippet", snippet.GetName())
}

func TestSnippet_Generate(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{
			name:  "OneLine",
			lines: []string{"abc"},
			want:  "abc\n",
		},
		{
			name:  "TwoLines",
			lines: []string{"abc", "def"},
			want:  "abc\ndef\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Snippet{Lines: tt.lines}
			got := testEmit(s)
			assert.Equal(t, tt.want, got)
		})
	}
}
