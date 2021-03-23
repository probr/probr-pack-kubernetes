# Config Engineering Notes

Internal code to manage config, including Cloud Driver parameters and Test Packs

## Log Filter Guidelines

Probr Log Levels:
- **ERROR** - Behavior that is a result of a definite misconfiguration or code failure
- **WARN** - Behavior that is likely due to a misconfiguration, but is not fatal
- **NOTICE** - (1) User config information to prevent confusion, or (2) behavior that could result from a misconfiguration but also may be intentional
- **INFO** - Non-verbose information that doesn't fit the above criteria
- **DEBUG** - Any potentially helpful information that doesn't fit the above criteria

Multi-line logs should be formatted prior to `//log.Printf(...)`. By using this command multiple times, each line will get a separate timestamp and will appear to be separate entries.

For example, `Results: ` could be read as if an empty string was being output.

However, by misusing `//log.Printf` we may cause a similar appearance:

```
//log.Printf("[NOTICE] Results:")
//log.Printf("[NOTICE] %s", myVar)
// Prints:
// 2020/09/28 11:18:01 [NOTICE] Results:
// 2020/09/28 11:18:01 [NOTICE] {"some": "information"}
```

## Config

Configuration docs are located in the README at the top level of the probr repository.

When creating new config vars, remember to do the following:

1. Add an entry to the struct `ConfigVars` in `internal/config/config.go`
1. Add an entry (matching the config vars struct) to `setEnvOrDefaults` in `internal/config/defaults.go`
1. If appropriate, add logic to `cmd/probr-cli/flags.go`

By following the above steps, you will have accomplished the following:
1. A new variable will be available across the entire probr codebase
1. That variable will have a default value
1. An environment variable can be set to override the default value
1. The env var can be overridden by a provided yaml config file
1. If set, a flag can be used to override the all other values
