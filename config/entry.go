package config

import (
	"github.com/pkg/errors"
	"sort"
	"strconv"
	"strings"
)

type Layers []ProviderEntries

func (l Layers) Merge() ProviderEntries {
	var merged = make(map[string]ProviderEntry)
	// Later layers override earlier layers
	for i := len(l) - 1; i >= 0; i-- {
		layer := l[i]
		for _, entry := range layer {
			if _, ok := merged[entry.NormalizedName]; !ok {
				merged[entry.NormalizedName] = entry
			}
		}
	}

	var results = make(ProviderEntries, 0, len(merged))
	for _, entry := range merged {
		results = append(results, entry)
	}
	results.SortByNormalizedName()
	return results
}

type ProviderEntries []ProviderEntry

func (e ProviderEntries) SortByNormalizedName() {
	sort.Slice(e,
		ProviderEntrySorter{
			entries: e,
			fn:      entryLessByNormalizedName,
		}.Less)
}

func (e ProviderEntries) Validate() error {
	if len(e) == 0 {
		return nil
	}

	e.SortByNormalizedName()

	if e[0].NormalizedName == "" {
		return ErrEmptyKey
	}

	for i := 1; i < len(e); i++ {
		if e[i].NormalizedName == e[i-1].NormalizedName {
			return errors.Wrapf(ErrDuplicateKey,
				"Duplicate normalized name %q detected: %q vs %q",
				e[i].NormalizedName,
				e[i].Name,
				e[i-1].Name)
		}
	}

	return nil
}

func (e ProviderEntries) Clone() ProviderEntries {
	entries := make([]ProviderEntry, 0, len(e))
	entries = append(entries, e...)
	return entries
}

func (e ProviderEntries) Compare(other ProviderEntries) ProviderDelta {
	var delta ProviderDelta
	le, re := e, other
	li, ri := 0, 0

	lv, rv := li < len(le), ri < len(re)
	for lv || rv {

		switch {
		case lv && rv && le[li].NormalizedName == re[ri].NormalizedName:
			// Updated
			if le[li].Value != re[ri].Value {
				delta = append(delta, ProviderEntryDelta{
					OldEntry: le[li],
					NewEntry: re[ri],
				})
			}
			li++
			ri++

		case (lv && !rv) || (lv && rv && le[li].NormalizedName < re[ri].NormalizedName):
			// Removed
			delta = append(delta, ProviderEntryDelta{
				OldEntry: le[li],
				NewEntry: ProviderEntry{
					NormalizedName: le[li].NormalizedName,
				},
			})
			li++

		case (rv && !lv) || (lv && rv && le[li].NormalizedName > re[ri].NormalizedName):
			// Added
			delta = append(delta, ProviderEntryDelta{
				NewEntry: re[ri],
				OldEntry: ProviderEntry{
					NormalizedName: re[ri].NormalizedName,
				},
			})
			ri++

		}

		lv, rv = li < len(le), ri < len(re)
	}

	return delta
}

func (e ProviderEntries) Append(other ProviderEntries) ProviderEntries {
	entries := make([]ProviderEntry, 0, len(e) + len(other))
	entries = append(entries, e...)
	entries = append(entries, other...)
	return entries
}

type ProviderEntrySorter struct {
	entries []ProviderEntry
	fn      func([]ProviderEntry, int, int) bool
}

func (e ProviderEntrySorter) Less(i, j int) bool {
	return e.fn(e.entries, i, j)
}

func entryLessByNormalizedName(entries []ProviderEntry, i, j int) bool {
	return entries[i].NormalizedName < entries[j].NormalizedName
}

type ProviderEntry struct {
	NormalizedName string
	Name           string
	Value          string
	Source         Provider
}

func NewEntry(source Provider, name, value string) ProviderEntry {
	return ProviderEntry{
		NormalizedName: NormalizeKey(name),
		Name:           name,
		Value:          value,
		Source:         source,
	}
}

type ProviderEntryDelta struct {
	OldEntry ProviderEntry
	NewEntry ProviderEntry
}

func (e ProviderEntryDelta) IsSet() bool {
	return e.NewEntry.Source != nil
}

type ProviderDelta []ProviderEntryDelta

func NormalizeKey(key string) string {
	return strings.ToLower(
		strings.ReplaceAll(
			strings.ReplaceAll(
				key,
				"-",
				""),
			"_",
			"."))
}

func PrefixWithName(prefix string, key string) string {
	if prefix == "" {
		return key
	}
	if strings.HasSuffix(prefix, ".") {
		prefix = prefix[:len(prefix)-1]
	}
	return prefix + "." + key
}

func PrefixWithIndex(prefix string, index int) string {
	suffix := "[" + strconv.Itoa(index) + "]"
	if prefix == "" {
		return suffix
	}
	if strings.HasSuffix(prefix, ".") {
		prefix = prefix[:len(prefix)-1]
	}
	return prefix + suffix
}

func MapFromEntries(entries ProviderEntries) map[string]string {
	settings := map[string]string{}
	for _, entry := range entries {
		settings[entry.Name] = entry.Value
	}

	return settings
}
