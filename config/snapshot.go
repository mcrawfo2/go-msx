// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"fmt"
	"github.com/pkg/errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ResolvedEntry struct {
	ProviderEntry
	ResolvedValue Value
}

func (s ResolvedEntry) String() string {
	return fmt.Sprintf("%q => %q", s.Name, s.ResolvedValue)
}

func (s ResolvedEntry) HasPrefix(prefix string) bool {
	k := s.NormalizedName
	if len(k) < len(prefix) {
		return false
	}
	if len(k) == len(prefix) && k == prefix {
		return true
	}
	if !strings.HasPrefix(k, prefix) {
		return false
	}
	next := k[len(prefix)]
	return next == '.' || next == '['
}

type resolvedEntrySorter struct {
	entries []ResolvedEntry
	fn      func([]ResolvedEntry, int, int) bool
}

func (e resolvedEntrySorter) Less(i, j int) bool {
	return e.fn(e.entries, i, j)
}

func snapshotEntryLessByNormalizedName(entries []ResolvedEntry, i, j int) bool {
	return entries[i].NormalizedName < entries[j].NormalizedName
}

type resolvedEntryDelta struct {
	OldEntry ResolvedEntry
	NewEntry ResolvedEntry
}

func (e resolvedEntryDelta) IsSet() bool {
	return e.NewEntry.Source != nil
}

func (e resolvedEntryDelta) IsUnset() bool {
	return e.NewEntry.Source == nil
}

type SnapshotDelta []resolvedEntryDelta

type nodeName struct {
	NormalizedName string
	Prefix         string
	Name           string
	Suffix         string
	Index          int
}

type ResolvedEntries []ResolvedEntry

func (e ResolvedEntries) SortByNormalizedName() {
	sort.Slice(e,
		resolvedEntrySorter{
			entries: e,
			fn:      snapshotEntryLessByNormalizedName,
		}.Less)
}

func (e ResolvedEntries) ChildNodeNames(prefix string) []nodeName {
	children := make(map[string]nodeName)

	for _, entry := range e {
		if !entry.HasPrefix(prefix) {
			continue
		}
		if len(prefix) == len(entry.NormalizedName) {
			continue
		}

		normalizedName := ""
		suffix := entry.NormalizedName[len(prefix)+1:]
		name := ""
		index := 0
		sep := entry.NormalizedName[len(prefix)]
		if sep == '.' {
			suffixPrefixEnd := strings.Index(suffix, ".")
			if suffixPrefixEnd == -1 {
				suffixPrefixEnd = len(suffix)
			}
			suffix = entry.NormalizedName[len(prefix) : suffixPrefixEnd+len(prefix)+1]
			name = suffix[1:]
			normalizedName = entry.NormalizedName[:suffixPrefixEnd+len(prefix)+1]
		} else if sep == '[' {
			suffixPrefixEnd := strings.Index(suffix, "]")
			if suffixPrefixEnd == -1 {
				continue
			}
			suffix = entry.NormalizedName[len(prefix) : suffixPrefixEnd+len(prefix)+2]
			index, _ = strconv.Atoi(suffix[1 : len(suffix)-1])
			normalizedName = entry.NormalizedName[:suffixPrefixEnd+len(prefix)+2]
		}

		if _, ok := children[suffix]; !ok {
			children[suffix] = nodeName{
				NormalizedName: normalizedName,
				Prefix:         prefix,
				Name:           name,
				Suffix:         suffix,
				Index:          index,
			}
		}

	}

	var results []nodeName
	for _, v := range children {
		results = append(results, v)
	}
	return results
}

func (e ResolvedEntries) Compare(other ResolvedEntries) SnapshotDelta {
	e.SortByNormalizedName()
	other.SortByNormalizedName()

	var changes SnapshotDelta
	le, re := e, other
	li, ri := 0, 0

	lv, rv := li < len(le), ri < len(re)
	for lv || rv {

		switch {
		case lv && rv && le[li].NormalizedName == re[ri].NormalizedName:
			// Updated
			if le[li].ResolvedValue != re[ri].ResolvedValue {
				changes = append(changes, resolvedEntryDelta{
					OldEntry: le[li],
					NewEntry: re[ri],
				})
			}
			li++
			ri++

		case (lv && !rv) || (lv && rv && le[li].NormalizedName < re[ri].NormalizedName):
			// Removed
			changes = append(changes, resolvedEntryDelta{
				OldEntry: le[li],
				NewEntry: ResolvedEntry{
					ProviderEntry: ProviderEntry{
						NormalizedName: le[li].NormalizedName,
					},
				},
			})
			li++

		case (rv && !lv) || (lv && rv && le[li].NormalizedName > re[ri].NormalizedName):
			// Added
			changes = append(changes, resolvedEntryDelta{
				NewEntry: re[ri],
				OldEntry: ResolvedEntry{
					ProviderEntry: ProviderEntry{
						NormalizedName: re[ri].NormalizedName,
					},
				},
			})
			ri++

		}

		lv, rv = li < len(le), ri < len(re)
	}

	return changes
}

type SnapshotValues struct {
	index   map[string]int
	entries ResolvedEntries
}

func (s SnapshotValues) Empty() bool {
	return len(s.index) == 0
}

func (s SnapshotValues) String(key string) (string, error) {
	entry, err := s.ResolveByName(key)
	if err != nil {
		return "", err
	}

	return entry.ResolvedValue.String(), nil
}

func (s SnapshotValues) StringOr(key, alt string) (string, error) {
	value, err := s.String(key)
	if errors.Is(err, ErrNotFound) {
		return alt, nil
	} else if err != nil {
		return "", nil
	}
	return value, nil
}

func (s SnapshotValues) Int(key string) (int, error) {
	entry, err := s.ResolveByName(key)
	if err != nil {
		return 0, err
	}

	ival, err := entry.ResolvedValue.Int()
	if err != nil {
		return 0, err
	}

	return int(ival), nil
}

func (s SnapshotValues) IntOr(key string, alt int) (int, error) {
	value, err := s.Int(key)
	if errors.Is(err, ErrNotFound) {
		return alt, nil
	} else if err != nil {
		return 0, nil
	}
	return value, nil
}

func (s SnapshotValues) Uint(key string) (uint, error) {
	entry, err := s.ResolveByName(key)
	if err != nil {
		return 0, err
	}

	ival, err := entry.ResolvedValue.Uint()
	if err != nil {
		return 0, err
	}

	return uint(ival), nil
}

func (s SnapshotValues) UintOr(key string, alt uint) (uint, error) {
	value, err := s.Uint(key)
	if errors.Is(err, ErrNotFound) {
		return alt, nil
	} else if err != nil {
		return 0, nil
	}
	return value, nil
}

func (s SnapshotValues) Float(key string) (float64, error) {
	entry, err := s.ResolveByName(key)
	if err != nil {
		return 0, err
	}

	return entry.ResolvedValue.Float()
}

func (s SnapshotValues) FloatOr(key string, alt float64) (float64, error) {
	value, err := s.Float(key)
	if errors.Is(err, ErrNotFound) {
		return alt, nil
	} else if err != nil {
		return 0, nil
	}
	return value, nil
}

func (s SnapshotValues) Bool(key string) (bool, error) {
	entry, err := s.ResolveByName(key)
	if err != nil {
		return false, err
	}

	return entry.ResolvedValue.Bool()
}

func (s SnapshotValues) BoolOr(key string, alt bool) (bool, error) {
	value, err := s.Bool(key)
	if errors.Is(err, ErrNotFound) {
		return alt, nil
	} else if err != nil {
		return false, nil
	}
	return value, nil
}

func (s SnapshotValues) Duration(key string) (time.Duration, error) {
	entry, err := s.ResolveByName(key)
	if err != nil {
		return 0, err
	}

	ival, err := entry.ResolvedValue.Duration()
	if err != nil {
		return 0, err
	}

	return ival, nil
}

func (s SnapshotValues) DurationOr(key string, alt time.Duration) (time.Duration, error) {
	value, err := s.Duration(key)
	if errors.Is(err, ErrNotFound) {
		return alt, nil
	} else if err != nil {
		return 0, nil
	}
	return value, nil
}

func (s SnapshotValues) Settings() map[string]string {
	var results = make(map[string]string, len(s.entries))
	for _, entry := range s.entries {
		results[entry.NormalizedName] = entry.ResolvedValue.String()
	}
	return results
}

func (s SnapshotValues) Each(target func(string, string)) {
	for _, entry := range s.entries {
		target(entry.NormalizedName, entry.ResolvedValue.String())
	}
}

func (s SnapshotValues) Populate(target interface{}, prefix string) error {
	return Populate(target, prefix, s)
}

func (s SnapshotValues) Value(key string) (Value, error) {
	entry, err := s.ResolveByName(key)
	if err != nil {
		return "", err
	}

	return entry.ResolvedValue, nil
}

func (s SnapshotValues) Entries() ResolvedEntries {
	return append(ResolvedEntries{}, s.entries...)
}

func (s SnapshotValues) ValuesWithPrefix(prefix string) SnapshotValues {
	if prefix == "" {
		return s
	}

	prefix = NormalizeKey(prefix)

	// Find the first key with a prefix
	startIdx := 0
	for startIdx < len(s.entries) {
		if s.entries[startIdx].HasPrefix(prefix) {
			break
		}
		startIdx++
	}

	// Find the last key with a prefix
	endIdx := startIdx
	for endIdx < len(s.entries) {
		if !s.entries[endIdx].HasPrefix(prefix) {
			break
		}
		endIdx++
	}

	return newSnapshotValues(s.entries[startIdx:endIdx])
}

func (s SnapshotValues) ResolveByName(name string) (ResolvedEntry, error) {
	key := NormalizeKey(name)
	idx, ok := s.index[key]
	if !ok {
		return ResolvedEntry{}, errors.Wrap(ErrNotFound, name)
	}

	return s.entries[idx], nil
}

func newSnapshotValues(entries ResolvedEntries) SnapshotValues {
	index := make(map[string]int)
	for n, entry := range entries {
		index[entry.NormalizedName] = n
	}

	return SnapshotValues{
		index:   index,
		entries: entries,
	}
}

var emptySnapshotValues = SnapshotValues{
	index:   make(map[string]int),
	entries: ResolvedEntries{},
}

type Snapshot struct {
	SnapshotValues
	Delta SnapshotDelta
}

func newSnapshot(entries ProviderEntries, previous SnapshotValues) (Snapshot, error) {
	r := NewResolver(entries)
	resolvedEntries, err := r.Entries()
	if err != nil {
		return Snapshot{}, err
	}

	values := newSnapshotValues(resolvedEntries)

	delta := previous.entries.Compare(resolvedEntries)

	return Snapshot{
		SnapshotValues: values,
		Delta:          delta,
	}, nil
}
