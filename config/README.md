# MSX Configuration Module

MSX configuration is a modified version of the [github.com/zpatrick/go-config](README.orig.md) package with added support for remote configuration stores, JSON5 files, key normalization, and structure population.

## Model

MSX Configuration has three main components: **providers**, **settings**, and the **config** object. 
* **Providers** load settings for your application. This could be from a file, environment variables, or some other source of configuration.
* **Settings** represent the configuration options for your application. Settings are represented as key/value pairs. 
* **Config** holds all of the providers and loaded settings. This object allows you to load, watch, retrieve, apply and convert your settings.

## Quick Start

### Instantiation

When using MSX Configuration inside the MSX Application context, you can retrieve the configuration object from the `ctx context.Context`:

```go
cfg := config.MustFromContext(ctx)
```

When using MSX Configuration outside of the MSX Application context, you can instantiate your own providers.  For example, to consume the environment variables from the current process:

```go
environmentProvider := config.NewEnvironment("env")
cfg := config.NewConfig(environmentProvider)
```

### Value Retrieval

Using one of the above `cfg` objects, you can retrieve the user's home directory from the `HOME` environment variable.  Note that all config keys are normalized to be lowercase, no hyphens, period-separated.  This means the `HOME` environment variable will be mapped to `home`:

```go
homePath, err := cfg.String("home")
```

The `cfg` object presents a number of functions to return a strongly-typed value:
- `String(key string)`
- `Int(key string)`
- `Float(key string)`
- `Bool(key string)`

These functions will look up the specified key in the configuration, and if found, will attempt to convert the value to the specified type.  If the key is not found or configuration has not yet been loaded, an appropriate error will be returned.

If you wish to use an alternative (default) value in case the lookup fails, you can use the `Or` functions:
- `StringOr(key string, other string)`
- `Int(key string, other int)`
- `Float(key string, other float)`
- `Bool(key string, other bool)`

The specified `other` value will be returned if the config has been loaded, but lookup fails:
```go
buildPath, err := cfg.StringOr("build.path", "./build")
```

### Structure Population

You can also populate appropriately defined structures:

```go
type ConnectionConfig struct {
    Name string
    Skipped bool `config:-`
    AnotherName int `config:somethingelse`
}

var connectionConfig ConnectionConfig
err := cfg.Populate(&connectionConfig, "some.connection")
```

Each structure field is treated a little differently based on the contents/existence of the `config` struct tag:
- `Name`: populated from `some.connection.name` (default behaviour)
- `Skipped`: not populated due to the `config:-` (omit when source name is a hyphen)
- `AnotherName`: populated from `some.connection.somethingelse` (overridden field name)

## Spring Compatibility

One of the primary goals for MSX Configuration is close compatibility with Spring-style configuration.
Several known incompatibilities and limitations currently exist:
  1. Key Normalization
     - Configuration keys in MSX Configuration are simply normalized to be lowercase, no hyphens, period-separated. As of Spring 2.0, configuration keys are expected to be snake-case, period-separated. MSX Configuration cannot distinguish between the `app.some-data` and `app.somedata` keys, and normalizes them both to `app.somedata`.
  1. Arbitrary Population
     - MSX Configuration currently supports `@ConfigurationProperties` style structure population.  As a consequence, all data used to populate a structure must be direct descendants of the key used to populate the structure.  We intend to support arbitrary key specification for structures in the future.
  1. Nested Defaults
     - MSX Configuration does not currently support nested/chained default values.

## Built-In Providers 

MSX Configuration has many built-in providers, allowing the application to unify configuration from a wide variety of sources:

* `INIFile` - Loads settings from a `.ini` file
* `JSONFile`  - Loads settings from a `.json` or `.json5` file
* `YAMLFile` - Loads settings from a `.yaml` or `.yml` file
* `TOMLFile` - Loads settings from a `.toml` file
* `PropertiesFile` - Loads settings from a `.properties` file
* `CobraProvider` - Loads settings from a Cobra command context
* `PFlagProvider` - Loads settings from PFlag flagset
* `GoFlagProvider` - Loads settings from a go flag flagset
* `ConsulProvider` - Loads settings from Consul
* `VaultProvider` - Loads settings from Vault
* `Environment` - Loads settings from environment variables 
* `Static` - Loads settings from an in-memory map

### Helpers
Along with the above providers, there are some wrappers for managing config lifecycle:
* `CachedLoader` - Caches settings in memory until flushed by `Invalidate()`
* `OnceLoader` - Caches settings in memory permanently
* `Resolver` - Remaps settings from one key to another

