//go:generate mockery --name FileWatcher --structname MockFileWatcher --filename mock_FileWatcher_test.go --inpackage

package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-ini/ini"
	"github.com/magiconair/properties"
	"github.com/radovskyb/watcher"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"
)

type FileProvider struct {
	name   string
	path   string
	Format ContentFormat
	Reader ContentReader
	Notifier
}

func (p *FileProvider) Description() string {
	return fmt.Sprintf("%s: [%s]", p.name, p.path)
}

func (p *FileProvider) Load(ctx context.Context) (ProviderEntries, error) {
	return p.Format(p.Reader, p)
}

type ContentReader func() ([]byte, error)

func FileContentReader(fileName string) ContentReader {
	return func() (bytes []byte, err error) {
		return ioutil.ReadFile(fileName)
	}
}

func HttpFileContentReader(fs http.FileSystem, fileName string) ContentReader {
	return func() (bytes []byte, err error) {
		file, err := fs.Open(fileName)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(file)
	}
}

type FileWatcher interface {
	SetMaxEvents(int)
	FilterOps(...watcher.Op)
	Start(time.Duration) error
	Event() <-chan watcher.Event
	Error() <-chan error
	Closed() <-chan struct{}
	Close()
}

type FileWatcherFactory func(filename string) (FileWatcher, error)

type WatcherFileWatcher struct {
	*watcher.Watcher
}

func (w WatcherFileWatcher) Event() <-chan watcher.Event {
	return w.Watcher.Event
}

func (w WatcherFileWatcher) Error() <-chan error {
	return w.Watcher.Error
}

func (w WatcherFileWatcher) Closed() <-chan struct{} {
	return w.Watcher.Closed
}

func NewWatcherFileWatcher(filename string) (FileWatcher, error) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)
	err := w.Add(filename)
	if err != nil {
		return nil, err
	}
	return WatcherFileWatcher{
		Watcher: w,
	}, nil
}

type FileNotifier struct {
	filename       string
	notify         chan struct{}
	watcherFactory FileWatcherFactory
}

func (i *FileNotifier) Notify() <-chan struct{} {
	return i.notify
}

func (i *FileNotifier) Run(ctx context.Context) {
	w, err := i.watcherFactory(i.filename)
	if err != nil {
		logger.WithError(err).Errorf("Failed to watch config file %q", i.filename)
		return
	}

	go func() {
		err := w.Start(time.Millisecond * 1000)
		if err != nil {
			logger.WithError(err).Errorf("Failed to watch config file %q", i.filename)
		}
	}()

	for {
		select {
		case event := <-w.Event():
			logger.Infof("%s event received for config file %q", event.Op.String(), i.filename)
			i.notify <- struct{}{}
		case err := <-w.Error():
			logger.Error(err.Error())
		case <-w.Closed():
			return
		case <-ctx.Done():
			w.Close()
			return
		}
	}
}

func NewFileNotifier(filename string) *FileNotifier {
	return &FileNotifier{
		filename:       filename,
		notify:         make(chan struct{}),
		watcherFactory: NewWatcherFileWatcher,
	}
}

type ContentFormat func(reader ContentReader, provider Provider) (ProviderEntries, error)

func ParseIni(reader ContentReader, provider Provider) (ProviderEntries, error) {
	bytes, err := reader()
	if err != nil {
		return nil, err
	}

	file, err := ini.Load(bytes)
	if err != nil {
		return nil, err
	}

	var settings ProviderEntries

	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			name := fmt.Sprintf("%s.%s", section.Name(), key.Name())
			settings = append(settings, NewEntry(provider, name, key.String()))
		}
	}

	return settings, nil
}

func ParseProperties(reader ContentReader, provider Provider) (ProviderEntries, error) {
	l := &properties.Loader{
		Encoding:         properties.UTF8,
		DisableExpansion: true}

	bytes, err := reader()
	if err != nil {
		return nil, err
	}

	props, err := l.LoadBytes(bytes)
	if err != nil {
		return nil, err
	}

	var settings ProviderEntries

	for _, key := range props.Keys() {
		value, _ := props.Get(key)
		settings = append(settings, NewEntry(provider, key, value))
	}

	return settings, nil
}

func ParseYaml(reader ContentReader, provider Provider) (ProviderEntries, error) {
	yamlBytes, err := reader()
	if err != nil {
		return nil, err
	}

	jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
	if err != nil {
		return nil, err
	}

	jsonData := map[string]interface{}{}
	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		return nil, err
	}

	return visitJsonNode(provider, jsonData, "")
}

func ParseJson(reader ContentReader, provider Provider) (ProviderEntries, error) {
	jsonBytes, err := reader()
	if err != nil {
		return nil, err
	}

	jsonData := map[string]interface{}{}
	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		return nil, err
	}

	return visitJsonNode(provider, jsonData, "")
}

func visitJsonNode(source Provider, value interface{}, prefix string) (ProviderEntries, error) {
	switch v := value.(type) {
	case map[string]interface{}:
		return visitMap(source, v, prefix)

	case []interface{}:
		return visitArray(source, v, prefix)

	default:
		return visitScalar(source, v, prefix)
	}
}

func visitScalar(source Provider, input interface{}, prefix string) (ProviderEntries, error) {
	return ProviderEntries{
		NewEntry(source, prefix, fmt.Sprintf("%v", input)),
	}, nil
}

func visitMap(source Provider, input map[string]interface{}, prefix string) (ProviderEntries, error) {
	var results ProviderEntries

	for key, value := range input {
		token := PrefixWithName(prefix, key)

		settings, err := visitJsonNode(source, value, token)
		if err != nil {
			return nil, err
		}

		results = results.Append(settings)
	}

	return results, nil
}

func visitArray(source Provider, input []interface{}, prefix string) (ProviderEntries, error) {
	var results ProviderEntries

	for k, v := range input {
		token := PrefixWithIndex(prefix, k)

		settings, err := visitJsonNode(source, v, token)
		if err != nil {
			return nil, err
		}

		results = results.Append(settings)
	}
	return results, nil
}

func newFileProvider(name, fileName string, reader ContentReader, notifier Notifier) Provider {
	fileExt := strings.ToLower(path.Ext(fileName))
	var format ContentFormat

	switch fileExt {
	case ".yml", ".yaml":
		format = ParseYaml
	case ".ini":
		format = ParseIni
	case ".json", ".json5":
		format = ParseJson
	case ".properties":
		format = ParseProperties
	default:
		logger.Error("Unknown config file extension: ", fileExt)
		return nil
	}

	if notifier == nil {
		notifier = SilentNotifier{}
	}

	return &FileProvider{
		name:     name,
		path:     fileName,
		Format:   format,
		Reader:   reader,
		Notifier: notifier,
	}
}
