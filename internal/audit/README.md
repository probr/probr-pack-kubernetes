# Output (Internal)

## Audit

Audit is intended for development use to enable systematic event logging.

These logs are designed to be one-per-event. Each test may have multiple events, and multiple events will likely share similar event types. Events will overwrite _within_ each test, however, to ensure that only one event type is logged per test.

### Usage

A new AuditLog context is created each time `probr` is run, and is readily accessible anywhere in the code via `audit.AuditLog`.

Adding entries to the an Event's meta data requires the name of the test and a key-value pair to be inserted:

```
n := "name-of-the-current-test"
k = "arbitrary_key_name"
v = "string_value"
audit.AuditLog.AuditMeta(n, k, v)
```

Events and logs may be introspected:

```
r := o.Events["name-of-the-current-test"]["pods_created"]
fmt.Println(r) // "1"
```

Logs may be formatted to JSON:

```
json.MarshalIndent(audit.AuditLog.Events, "", "  ")
```

Here is an example of what the logs may contain:

```
{
  "account_manager": {
    "Meta": {
      "category": "General",
      "group": "clouddriver",
      "status": "CompleteSuccess"
    },
    "PodsCreated": 0,
    "PodsDestroyed": 0,
  },
  "container_registry_access": {
    "Meta": {
      "category": "Container Registry Access",
      "group": "kubernetes",
      "status": "CompleteFail"
    },
    "PodsCreated": 1,
    "PodsDestroyed": 1,
  }
}
```

A method exists to print the audit to the command line under a NOTICE loglevel:

```
audit.AuditLog.PrintAudit()
```
