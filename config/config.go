package config

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var logger = log.NewLogger("msx.config")

var ErrNotLoaded = errors.New("Configuration not loaded")
var ErrNotFound = errors.New("Missing required setting")

type Config struct {
	Providers     []Provider
	Validate      func(map[string]string) error
	ReloadTimeout time.Duration
	ReloadContext context.Context
	Notify        func([]string)
	settings      map[string]string
}

func NewConfig(providers ...Provider) *Config {
	return &Config{
		Providers:     providers,
		settings:      nil,
		ReloadTimeout: time.Second * 90,
	}
}

func (c *Config) Loaded() bool {
	return c.settings != nil
}

func (c *Config) Load(ctx context.Context) error {
	if settings, err := c.reload(ctx); err != nil {
		return err
	} else {
		c.settings = settings
	}

	return nil
}

func (c *Config) reload(ctx context.Context) (map[string]string, error) {
	result := map[string]string{}

	// Load config from each provider, stacking appropriately
	for _, provider := range c.Providers {
		if settings, err := provider.Load(ctx); err != nil {
			return nil, errors.Wrap(err, "Failed to load config")
		} else {
			for key, val := range settings {
				result[key] = val
			}
		}
	}

	// Validate the config
	if c.Validate != nil {
		if err := c.Validate(result); err != nil {
			return nil, errors.Wrap(err, "Failed to validate config")
		}
	}

	// Resolve variables in the config
	if err := c.resolve(result); err != nil {
		return nil, errors.Wrap(err, "Failed to resolve variables")
	}

	if c.Loaded() && c.Notify != nil {
		c.compareAndNotify(result)
	}

	return result, nil
}

// Expand all references to ${variables} inside value
func (c *Config) resolveValue(resolved, settings map[string]string, value string) string {
	variableRegex, _ := regexp.Compile(`\${([\w._\-]+)(:([^}]*))?}`)

	if !strings.Contains(value, "${") {
		return value
	}

	stack := types.StringStack{"_"}
	defaults := make(map[string]string)
	for len(stack) > 0 {
		currentVariable := stack.Peek()
		currentValue := ""
		ok := false

		if currentValue, ok = resolved[currentVariable]; ok {
			// already resolved
		} else if currentVariable == "_" {
			// passed-in value
			currentValue = value
		} else if currentValue, ok = settings[currentVariable]; !ok {
			var defaultValue string
			if defaultValue, ok = defaults[currentVariable]; ok {
				currentValue = defaultValue
			} else {
				logger.Errorf("Failed to resolve variable %s", currentVariable)
				currentValue = ""
			}
		}

		unresolvedReferences := 0
		variables := variableRegex.FindAllStringSubmatch(currentValue, -1)
		for _, match := range variables {
			referenceVariableName := c.alias(match[1])
			if stack.Contains(referenceVariableName) {
				logger.Errorf("Circular variable reference detected: %s", referenceVariableName)
				resolved[referenceVariableName] = ""
			}
			if referenceVariableValue, ok := resolved[referenceVariableName]; ok {
				referenceRegex, _ := regexp.Compile(`\${` + strings.ReplaceAll(match[1], ".", "\\.") + `(:([^}]*))?}`)
				currentValue = referenceRegex.ReplaceAllLiteralString(currentValue, referenceVariableValue)
			} else {
				unresolvedReferences++
				stack = stack.Push(referenceVariableName)
				if len(match) == 4 && len(match[2]) > 0 {
					defaults[referenceVariableName] = match[3]
				}
			}
		}

		if unresolvedReferences == 0 {
			resolved[currentVariable] = currentValue
			stack = stack.Pop()
		}
	}

	value = resolved["_"]
	delete(resolved, "_")
	return value
}

// Expand all references to ${variables}
func (c *Config) resolve(settings map[string]string) error {
	resolved := map[string]string{}

	for k, v := range settings {
		resolved[k] = c.resolveValue(resolved, settings, v)
	}

	for k, v := range resolved {
		settings[k] = v
	}

	return nil
}

func (c *Config) compareAndNotify(newSettings map[string]string) {
	// Find changed variables
	changes := map[string]struct{}{}
	oldSettings := c.settings
	for k, v := range newSettings {
		if oldValue, ok := oldSettings[k]; ok {
			if v != oldValue {
				changes[k] = struct{}{}
			}
		} else {
			changes[k] = struct{}{}
		}
	}

	for k := range oldSettings {
		if _, ok := newSettings[k]; !ok {
			changes[k] = struct{}{}
		}
	}

	var changedVariables []string
	for k := range changes {
		changedVariables = append(changedVariables, k)
		// TODO: add reverse aliases
	}

	if len(changedVariables) > 0 {
		c.Notify(changedVariables)
	}
}

func (c *Config) Watch(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	notifier := make(chan struct{}, 1)
	go func() {
		defer close(notifier)

		for _, provider := range c.Providers {
			if watcher, ok := provider.(Watcher); ok {
				watcher.Watch(notifier, ctx)
			}
		}

		var err error

		for {
			select {
			case <-notifier:
				// Something was invalidated
				err = func() error {
					subctx, cancel := context.WithTimeout(c.ReloadContext, c.ReloadTimeout)
					defer cancel()
					return c.Load(subctx)
				}()

				if err != nil {
					logger.Error(errors.Wrap(err, "Failed to load configuration").Error())
				}

			case <-ctx.Done():
				return
			}
		}

	}()

	return ctx.Err()
}

func (c *Config) String(key string) (string, error) {
	if !c.Loaded() {
		return "", ErrNotLoaded
	}

	targetKey := c.alias(key)
	if val, ok := c.settings[targetKey]; !ok {
		return "", errors.Wrap(ErrNotFound, key)
	} else {
		return val, nil
	}
}

func (c *Config) StringOr(key, alt string) (string, error) {
	if !c.Loaded() {
		return "", ErrNotLoaded
	}

	targetKey := c.alias(key)
	if val, ok := c.settings[targetKey]; !ok {
		return alt, nil
	} else {
		return val, nil
	}
}

func (c *Config) Int(key string) (int, error) {
	if !c.Loaded() {
		return 0, ErrNotLoaded
	}

	targetKey := c.alias(key)
	if val, ok := c.settings[targetKey]; !ok {
		return 0, errors.Wrap(ErrNotFound, key)
	} else {
		return strconv.Atoi(val)
	}
}

func (c *Config) IntOr(key string, alt int) (int, error) {
	if !c.Loaded() {
		return 0, ErrNotLoaded
	}

	targetKey := c.alias(key)
	if val, ok := c.settings[targetKey]; !ok {
		return alt, nil
	} else {
		return strconv.Atoi(val)
	}
}

func (c *Config) Float(key string) (float64, error) {
	if !c.Loaded() {
		return 0, ErrNotLoaded
	}

	targetKey := c.alias(key)
	if val, ok := c.settings[targetKey]; !ok {
		return 0, errors.Wrap(ErrNotFound, key)
	} else {
		return strconv.ParseFloat(val, 64)
	}
}

func (c *Config) FloatOr(key string, alt float64) (float64, error) {
	if !c.Loaded() {
		return 0, ErrNotLoaded
	}

	targetKey := c.alias(key)
	if val, ok := c.settings[targetKey]; !ok {
		return alt, nil
	} else {
		return strconv.ParseFloat(val, 64)
	}
}

func (c *Config) Bool(key string) (bool, error) {
	if !c.Loaded() {
		return false, ErrNotLoaded
	}

	targetKey := c.alias(key)
	if val, ok := c.settings[targetKey]; !ok {
		return false, errors.Wrap(ErrNotFound, key)
	} else {
		return strconv.ParseBool(val)
	}
}

func (c *Config) BoolOr(key string, alt bool) (bool, error) {
	if !c.Loaded() {
		return false, ErrNotLoaded
	}

	targetKey := c.alias(key)
	if val, ok := c.settings[targetKey]; !ok {
		return alt, nil
	} else {
		return strconv.ParseBool(val)
	}
}

func (c *Config) Settings() map[string]string {
	result := map[string]string{}
	for k, v := range c.settings {
		result[k] = v
	}
	return result
}

func (c *Config) Each(target func(string, string)) {
	for name, value := range c.Settings() {
		target(name, value)
	}
}

func (c *Config) Populate(target interface{}, prefix string) error {
	// Wrap the properties map in a partial config
	partialConfig := NewPartialConfig(c.settings, c)

	// Filter by prefix
	partialConfig = partialConfig.FilterStripPrefix(NormalizeKey(prefix))

	// Populate the object from the properties map
	return partialConfig.Populate(target)
}

func (c *Config) alias(key string) string {
	// Return the actual target key name mapping through aliases
	return NormalizeKey(key)
}

func (c *Config) reloadContext() context.Context {
	if c.ReloadContext == nil {
		return context.Background()
	}
	return c.ReloadContext
}
