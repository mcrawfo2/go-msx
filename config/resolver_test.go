package config

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestResolver(t *testing.T) {
	type wants struct {
		resolved map[string]string
		err      error
	}
	tests := []struct {
		name     string
		settings map[string]string
		want     wants
	}{
		{
			name: "Simple",
			settings: map[string]string{
				"alpha": "a",
				"bravo": "${alpha}b",
			},
			want: wants{
				resolved: map[string]string{
					"alpha": "a",
					"bravo": "ab",
				},
				err:      nil,
			},
		},
		{
			name: "Utf8",
			settings: map[string]string{
				"alpha": `a${bravo:Σ}`,
				"charlie": `dΞf`,
			},
			want: wants{
				resolved: map[string]string{
					"alpha": "aΣ",
					"charlie": "dΞf",
				},
			},
		},
		{
			name: "Circular",
			settings: map[string]string{
				"alpha": "${charlie}a",
				"bravo": "b${alpha}",
				"charlie": "${bravo}c",
			},
			want: wants{
				resolved: nil,
				err:      ErrCircularReference,
			},

		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inMem := NewInMemoryProvider("in-memory", tt.settings)
			entries, err := inMem.Load(context.Background())
			assert.NoError(t, err)
			entries.SortByNormalizedName()

			resolver := NewResolver(entries)

			resolvedEntries, err := resolver.Entries()
			if err != nil {
				if tt.want.err == nil {
					t.Errorf("Expected no error, got = %v", err)
				} else if !errors.Is(err, tt.want.err) {
					t.Errorf("Expected error = %v, got = %v", tt.want.err, err)
				}
				return
			}

			resolved := newSnapshotValues(resolvedEntries).Settings()
			if !reflect.DeepEqual(tt.want.resolved, resolved) {
				t.Errorf("Expected resolved = %v, got = %v", tt.want.resolved, resolved)
			}
		})
	}
}
