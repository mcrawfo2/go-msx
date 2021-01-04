package config

type Resolver struct {
	entries  map[string]ProviderEntry
	resolved map[string]ResolvedEntry
	active   []ProviderEntry
}

func (r *Resolver) isActive(key string) bool {
	for _, a := range r.active {
		if a.NormalizedName == key {
			return true
		}
	}
	return false
}

func (r *Resolver) setActive(e ProviderEntry) error {
	if r.isActive(e.NormalizedName) {
		logger.Errorf("Circular reference detected:")
		for _, entry := range r.active {
			logger.Errorf("- %q => %q", entry.Name, entry.Value)
		}
		return ErrCircularReference
	}

	r.active = append(r.active, e)
	return nil
}

func (r *Resolver) setInactive(e ProviderEntry) {
	for i := 0; i < len(r.active); i++ {
		if r.active[i].NormalizedName == e.NormalizedName {
			r.active = append(r.active[:i], r.active[i+1:]...)
			return
		}
	}
}

func (r *Resolver) Entries() (ResolvedEntries, error) {
	for _, entry := range r.entries {
		_, err := r.Resolve(entry)
		if err != nil {
			return nil, err
		}
	}

	var results = make(ResolvedEntries, 0, len(r.resolved))
	for _, snapshotEntry := range r.resolved {
		results = append(results, snapshotEntry)
	}

	results.SortByNormalizedName()

	return results, nil
}

func (r *Resolver) Resolve(entry ProviderEntry) (ResolvedEntry, error) {
	if resolvedEntry, ok := r.resolved[entry.NormalizedName]; ok {
		return resolvedEntry, nil
	}

	err := r.setActive(entry)
	if err != nil {
		return ResolvedEntry{}, err
	}
	defer r.setInactive(entry)


	expr, err := parseExpression(entry.Value)
	if err != nil {
		return ResolvedEntry{}, err
	}

	value, err := expr.Resolve(r)
	if err != nil {
		return ResolvedEntry{}, err
	}

	resolved := ResolvedEntry{
		ProviderEntry: entry,
		ResolvedValue: Value(value),
	}

	r.resolved[entry.NormalizedName] = resolved

	return resolved, nil
}

func (r *Resolver) ResolveByName(name string) (ResolvedEntry, error) {
	normalizedName := NormalizeKey(name)

	if entry, ok := r.entries[normalizedName]; ok {
		return r.Resolve(entry)
	}

	return ResolvedEntry{}, ErrNotFound
}

func NewResolver(entries ProviderEntries) Resolver {
	var entryIndex = make(map[string]ProviderEntry)
	for _, entry := range entries {
		entryIndex[entry.NormalizedName] = entry
	}

	return Resolver{
		entries:  entryIndex,
		resolved: make(map[string]ResolvedEntry),
	}
}

