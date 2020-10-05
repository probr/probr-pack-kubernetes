# Audit Engineering Notes

Audit is intended for development use to enable systematic event logging.

These logs are designed to be one-per-event. Each test may have multiple events, and multiple events will likely share similar event types. Events will overwrite _within_ each test, however, to ensure that only one event type is logged per test.

### AuditLog

A new AuditLog context is created each time `probr` is run, and is readily accessible anywhere in the code via `audit.AuditLog`.


**AuditLog.AuditMeta**

Adding entries to the an Event's meta data requires the name of the test and a key-value pair to be inserted:

```
n := "name-of-the-current-test"
k = "arbitrary_key_name"
v = "string_value"
audit.AuditLog.AuditMeta(n, k, v)
```

**AuditLog.AuditPodName**

The names of all pods should be tracked, so users may identify whether Probr is the source of any unexpected pods in their cluster.

```
if pd != nil {
  s.PodName = pd.GetObjectMeta().GetName()
  audit.AuditLog.AuditPodName(s.PodName)
}
```

**AuditLog.GetEventLog**

Many audits will be made directly to events. In order to do so, the event must first be retrieved by name. In the example below, the event is being stored alongside other test context information for easy access during execution.

```
	ctx.BeforeScenario(func(s *godog.Scenario) {
		ps.name = s.Name
		ps.event = audit.AuditLog.GetEventLog(NAME)
		probes.LogScenarioStart(s)
	})
```

**AuditLog.EventComplete**

After an event has finished running every probe, we should audit the final outcome of the event.

```
s, o, err := g.Handler(g.Data)
audit.AuditLog.EventComplete(t.TestDescriptor.Name)
```

**AuditLog.SetProbrStatus**

After all events have completed, we should set the final probr status. This step may not always be relevant, as it may be possible to nest it within other methods such as `PrintAudit`. This should be reevaluated after more feedback has been gathered regarding how Probr is being used.

```
audit.AuditLog.SetProbrStatus()
```

**AuditLog.PrintAudit**

Instead of logging the final status of Probr within a particular _loglevel_ and formatting it using `log`, PrintAudit simply formats the status into JSON and prints it to the command line. This is currently the very last thing our CLI tool does prior to exiting.

```
audit.AuditLog.PrintAudit()
os.Exit(s)
```

### Events

**Event.CountPodCreated** and **Event.CountPodDestroyed**

Whenever Probr will create or destroy a pod, these counters should be called to update the event accordingly.

```
if pd != nil {
  s.PodName = pd.GetObjectMeta().GetName()
  e.CountPodCreated()
  audit.AuditLog.AuditPodName(s.PodName)
}
```

```
if wait {
  waitForDelete(c, ns, pname)
}
audit.AuditLog.GetEventLog(event).CountPodDestroyed()
```

**Event.AuditProbe**

This function should be used every time a step in a probe completes. `AuditProbe` will automatically form the name of the step from the name of the function that called it. The name value provided will establish which probe the step is a part of. The error (or nil) provided will dictate whether the test passes or fails.

```
p.event.AuditProbe(p.name, nil) // Success
```
