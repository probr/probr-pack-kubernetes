# Summary Engineering Notes

Summary is intended for development use to enable systematic probe logging.

These logs are designed to be one-per-probe. Each test may have multiple probes, and multiple probes will likely share similar probe types. Probes will overwrite _within_ each test, however, to ensure that only one probe type is logged per test.

### State

A new State context is created each time `probr` is run, and is readily accessible anywhere in the code via `audit.State`.


**SummaryStateStruct.LogProbeMeta**

Adding entries to the an Probe's meta data requires the name of the test and a key-value pair to be inserted. 

```
n := "name-of-the-current-test"
k = "arbitrary_key_name"
v = "string_value"
audit.State.LogProbeMeta(n, k, v)
```

**SummaryStateStruct.LogPodName**

The names of all pods should be tracked, so users may identify whether Probr is the source of any unexpected pods in their cluster.

```
if pd != nil {
  s.PodName = pd.GetObjectMeta().GetName()
  audit.State.LogPodName(s.PodName)
}
```

**SummaryStateStruct.GetProbeLog**

Many summaries will be made directly to probes. In order to do so, the probe must first be retrieved by name. In the example below, the probe is being stored alongside other test context information for easy access during execution.

```
	ctx.BeforeScenario(func(s *godog.Scenario) {
		ps.name = s.Name
		ps.probe = audit.State.GetProbeLog(NAME)
		coreengine.LogScenarioStart(s)
	})
```

**SummaryStateStruct.ProbeComplete**

After an probe has finished running every scenario, we should audit the final outcome of the probe.

```
s, o, err := g.Handler(g.Data)
audit.State.ProbeComplete(t.ProbeDescriptor.Name)
```

**SummaryStateStruct.SetProbrStatus**

After all probes have completed, we should set the final probr status. This step may not always be relevant, as it may be possible to nest it within other methods such as `PrintSummary`. This should be reevaluated after more feedback has been gathered regarding how Probr is being used.

```
audit.State.SetProbrStatus()
```

**SummaryStateStruct.PrintSummary**

Instead of logging the final status of Probr within a particular _loglevel_ and formatting it using `log`, PrintSummary simply formats the status into JSON and prints it to the command line. This is currently the very last thing our CLI tool does prior to exiting.

```
audit.State.PrintSummary()
os.Exit(s)
```

### Probes

**Probe.CountPodCreated** and **Probe.CountPodDestroyed**

Whenever Probr will create or destroy a pod, these counters should be called to update the probe accordingly.

**Probe.AuditScenarioStep**

This function should be used every time a step in a scenario completes. `AuditScenarioStep` will automatically form the name of the step from the name of the function that called it. The name value provided will establish which scenario the step is a part of. The error (or nil) provided will dictate whether the test passes or fails.

The description and payload values are arbitrary and are used only to assist auditors in their evaluation. A `nil` error will be recorded as a successful step.

```
s.audit.AuditScenarioStep( "description string", payloadObject, err) 
```
