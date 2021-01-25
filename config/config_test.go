package config

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	inMemoryProvider := NewInMemoryProvider("static", map[string]string{
		"alpha":   "a",
		"bravo":   "${alpha}b",
		"charlie": "${bravo}c",
		"delta":   "${charlie}d",
		"echo":    "${foxtrot:g}e",
	})

	config := NewConfig(inMemoryProvider)

	ctx, cancelCtx := context.WithCancel(context.Background())

	err := config.Load(ctx)
	assert.NoError(t, err)

	type expectedChange struct {
		Version       int
		Set           bool
		Name          string
		ResolvedValue string
	}

	expectedChanges := []expectedChange{
		{
			Version:       2,
			Set:           true,
			Name:          "echo",
			ResolvedValue: "fe",
		},
		{
			Version:       2,
			Set:           true,
			Name:          "foxtrot",
			ResolvedValue: "f",
		},
		{
			Version: 3,
			Set:     false,
			Name:    "charlie",
		},
		{
			Version:       3,
			Set:           true,
			Name:          "delta",
			ResolvedValue: "d",
		},
	}

	go func() {
		version := 1

		for {
			select {
			case <-ctx.Done():
				return

			case n := <-config.Notify():
				version++

				for _, c := range n.Delta {
					var expectedChange expectedChange
					assert.NotEmpty(t, expectedChanges)
					expectedChange, expectedChanges = expectedChanges[0], expectedChanges[1:]
					if c.IsSet() {
						assert.True(t, expectedChange.Set)
						assert.Equal(t, expectedChange.Version, version)
						assert.Equal(t, expectedChange.Name, c.NewEntry.Name)
						assert.Equal(t, expectedChange.ResolvedValue, c.NewEntry.ResolvedValue.String())
					} else {
						assert.False(t, expectedChange.Set)
						assert.Equal(t, expectedChange.Version, version)
						assert.Equal(t, expectedChange.Name, c.OldEntry.Name)
					}

				}

				if len(expectedChanges) == 0 {
					cancelCtx()
				}
			}
		}
	}()

	go func() {
		err := inMemoryProvider.SetValue("foxtrot", "f")
		assert.NoError(t, err)

		time.Sleep(1 * time.Second)

		err = inMemoryProvider.UnsetValue("charlie")
		assert.NoError(t, err)
	}()

	config.Watch(ctx)
}

func TestConfig_Values_Strings(t *testing.T) {
	inMemoryProvider := NewInMemoryProvider("static", map[string]string{
		"alpha": "a",
	})

	config := NewConfig(inMemoryProvider)

	// Errors before load
	alpha, err := config.String("alpha")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotLoaded))
	assert.Equal(t, "", alpha)

	alpha, err = config.StringOr("alpha", "z")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotLoaded))
	assert.Equal(t, "", alpha)

	// Load
	err = config.Load(context.Background())
	assert.NoError(t, err)

	// Loaded
	alpha, err = config.String("alpha")
	assert.NoError(t, err)
	assert.Equal(t, "a", alpha)

	alpha, err = config.StringOr("alpha", "z")
	assert.NoError(t, err)
	assert.Equal(t, "a", alpha)

	// Default
	charlie, err := config.StringOr("charlie", "z")
	assert.NoError(t, err)
	assert.Equal(t, "z", charlie)
}

func TestConfig_Values_Ints(t *testing.T) {
	inMemoryProvider := NewInMemoryProvider("static", map[string]string{
		"alpha": "42",
	})

	config := NewConfig(inMemoryProvider)

	// Errors before load
	alpha, err := config.Int("alpha")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotLoaded))
	assert.Equal(t, 0, alpha)

	alpha, err = config.IntOr("alpha", 21)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotLoaded))
	assert.Equal(t, 0, alpha)

	// Load
	err = config.Load(context.Background())
	assert.NoError(t, err)

	// Loaded
	alpha, err = config.Int("alpha")
	assert.NoError(t, err)
	assert.Equal(t, 42, alpha)

	alpha, err = config.IntOr("alpha", 21)
	assert.NoError(t, err)
	assert.Equal(t, 42, alpha)

	// Default
	charlie, err := config.IntOr("charlie", 21)
	assert.NoError(t, err)
	assert.Equal(t, 21, charlie)
}

func TestConfig_Values_Uints(t *testing.T) {
	inMemoryProvider := NewInMemoryProvider("static", map[string]string{
		"alpha": "42",
	})

	config := NewConfig(inMemoryProvider)

	// Errors before load
	alpha, err := config.Uint("alpha")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotLoaded))
	assert.Equal(t, uint(0), alpha)

	alpha, err = config.UintOr("alpha", 21)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotLoaded))
	assert.Equal(t, uint(0), alpha)

	// Load
	err = config.Load(context.Background())
	assert.NoError(t, err)

	// Loaded
	alpha, err = config.Uint("alpha")
	assert.NoError(t, err)
	assert.Equal(t, uint(42), alpha)

	alpha, err = config.UintOr("alpha", 21)
	assert.NoError(t, err)
	assert.Equal(t, uint(42), alpha)

	// Default
	charlie, err := config.UintOr("charlie", 21)
	assert.NoError(t, err)
	assert.Equal(t, uint(21), charlie)
}

func TestConfig_Values_Floats(t *testing.T) {
	inMemoryProvider := NewInMemoryProvider("static", map[string]string{
		"alpha": "12.5",
	})

	config := NewConfig(inMemoryProvider)

	// Errors before load
	alpha, err := config.Float("alpha")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotLoaded))
	assert.Equal(t, float64(0), alpha)

	alpha, err = config.FloatOr("alpha", 21)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotLoaded))
	assert.Equal(t, float64(0), alpha)

	// Load
	err = config.Load(context.Background())
	assert.NoError(t, err)

	// Loaded
	alpha, err = config.Float("alpha")
	assert.NoError(t, err)
	assert.Equal(t, 12.5, alpha)

	alpha, err = config.FloatOr("alpha", 6.25)
	assert.NoError(t, err)
	assert.Equal(t, 12.5, alpha)

	// Default
	charlie, err := config.FloatOr("charlie", 6.25)
	assert.NoError(t, err)
	assert.Equal(t, 6.25, charlie)
}

func TestConfig_Values_Bools(t *testing.T) {
	inMemoryProvider := NewInMemoryProvider("static", map[string]string{
		"alpha": "true",
	})

	config := NewConfig(inMemoryProvider)

	// Errors before load
	alpha, err := config.Bool("alpha")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotLoaded))
	assert.Equal(t, false, alpha)

	alpha, err = config.BoolOr("alpha", true)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotLoaded))
	assert.Equal(t, false, alpha)

	// Load
	err = config.Load(context.Background())
	assert.NoError(t, err)

	// Loaded
	alpha, err = config.Bool("alpha")
	assert.NoError(t, err)
	assert.Equal(t, true, alpha)

	alpha, err = config.BoolOr("alpha", false)
	assert.NoError(t, err)
	assert.Equal(t, true, alpha)

	// Default
	charlie, err := config.BoolOr("charlie", true)
	assert.NoError(t, err)
	assert.Equal(t, true, charlie)
}

func TestConfig_Values_Durations(t *testing.T) {
	inMemoryProvider := NewInMemoryProvider("static", map[string]string{
		"alpha": "10m",
	})

	config := NewConfig(inMemoryProvider)

	// Errors before load
	alpha, err := config.Duration("alpha")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotLoaded))
	assert.Equal(t, time.Duration(0), alpha)

	alpha, err = config.DurationOr("alpha", 15*time.Second)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotLoaded))
	assert.Equal(t, time.Duration(0), alpha)

	// Load
	err = config.Load(context.Background())
	assert.NoError(t, err)

	// Loaded
	alpha, err = config.Duration("alpha")
	assert.NoError(t, err)
	assert.Equal(t, 10*time.Minute, alpha)

	alpha, err = config.DurationOr("alpha", 15*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 10*time.Minute, alpha)

	// Default
	charlie, err := config.DurationOr("charlie", 15*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 15*time.Second, charlie)
}

func TestConfig_Values_Settings(t *testing.T) {
	inMemoryProvider := NewInMemoryProvider("static", map[string]string{
		"alpha": "true",
	})

	config := NewConfig(inMemoryProvider)

	// Empty before load
	settings := config.Settings()
	assert.Empty(t, settings)

	// Load
	err := config.Load(context.Background())
	assert.NoError(t, err)

	// Loaded
	settings = config.Settings()
	assert.Len(t, settings, 1)
}

func TestConfig_Value(t *testing.T) {
	inMemoryProvider := NewInMemoryProvider("static", map[string]string{
		"alpha": "abc",
	})

	config := NewConfig(inMemoryProvider)
	err := config.Load(context.Background())
	assert.NoError(t, err)

	v, err := config.Value("alpha")
	assert.NoError(t, err)
	assert.Equal(t, "abc", v.String())

	v, err = config.Value("beta")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestConfig_Each(t *testing.T) {
	inMemoryProvider := NewInMemoryProvider("static", map[string]string{
		"alpha": "abc",
		"delta": "def",
	})

	config := NewConfig(inMemoryProvider)
	err := config.Load(context.Background())
	assert.NoError(t, err)

	var count = 0
	config.Each(func(k string, v string) {
		assert.Equal(t, inMemoryProvider.settings[k], v)
		count++
	})
	assert.Equal(t, 2, count)
}
