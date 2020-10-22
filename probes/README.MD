# Scenarios Engineering Notes

Each "scenario" is part of a Probr "probe" which provides topical organization. Probes are sorted within topical subdirectories within `probes/`.

## Features == Probes // true

The entry point for each **probe** is a behavior-driven development (BDD) _feature file_ that contains human-readable instructions for one or more probes.  Each of these feature files represents a separate Probr **probe**.

The term "feature" here harkens to the original intent of BDD and _cucumber_, the technology we currently use for BDD. Cucumber and BDD logic traditionally revolves around writing feature tests, which simulates user interactions with an application.

If it serves you to simply read _feature_ instead of **probe**, you would not be inaccurate for making that mental substitution when reading the Probr code, documentation, and logs. However, we have chosen to use the term "probes" because we believe it more clearly demonstrates that this is part of the core Probr functionality- and not simply a feature test meant to ensure system reliability.

In the summary logs you will see the term "Probes" as a top level object, which will contain each individual **probe**, which will then subsequently contain every individual probe.

For more reading on how a _feature file_ is written, please review the [Cucumber documentation](https://cucumber.io/docs/gherkin/reference/)

## Scenarios == Scenarios // true

Scenarios are the heart of Probr, for obvious reasons. Each individual **scenario** will attempt a behavior, and return a success or failure based on whether the result matches the expected result.

These **scenarios** are defined in _feature files_, and in Gherkin they would be considered _senarios_ or _scenario outlines_.

Similar to the section above, we believe that the most accurately descriptive terminology for our use-case is to substitute in the term **scenario** for "scenario."

## Probe/Feature Logic

Because each scenario is a _scenario_, it will have a corresponding golang function for each _step_. The primary logic for the steps is found in a `.go` file that is named according to the probe. Comments surrounding the logic will help distinguish which function belongs to which scenario. Some steps may be used across multiple probes, and may be found in `scenarios/scenarios.go`

Finally, each step may require more comprehensive logic. Currently, that supplemental logic lives in appropriately named directories within `internal/`.
