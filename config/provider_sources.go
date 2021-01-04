package config

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type SourcesProvider struct {
	Named
	SilentNotifier
}

func (p *SourcesProvider) Load(ctx context.Context) (entries ProviderEntries, err error) {
	commandDir, err := types.FindEntryPointDirFromStack()
	if err == types.ErrSourceDirUnavailable {
		logger.WithContext(ctx).WithError(err).Warningf("Did not detect source directory.")
		return entries, nil
	}
	entries = append(entries, NewEntry(p, "fs.roots.command", commandDir))

	sourceDir, err := types.FindSourceDirFromStack()
	if err == types.ErrSourceDirUnavailable {
		logger.WithContext(ctx).WithError(err).Warningf("Did not detect source directory.")
		return entries, nil
	} else if err != nil {
		return nil, err
	}

	entries = append(entries, NewEntry(p, "fs.sources", sourceDir))

	return entries, nil
}

func NewSourcesProvider(name string) *SourcesProvider {
	return &SourcesProvider{
		Named: NewNamed(name),
	}
}
